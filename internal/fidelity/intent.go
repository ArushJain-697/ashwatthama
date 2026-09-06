package fidelity

import (
	"sort"
	"strings"
)

// ExtractedEntity is a named intent target. Method and confidence remain part
// of the schema even for deterministic direct mentions so later constrained
// fallback extraction cannot be mistaken for an exact transcript claim.
type ExtractedEntity struct {
	Name       string  `json:"name"`
	Method     string  `json:"method"`
	Confidence float64 `json:"confidence"`
}

// Intent is Stage 2's provider-independent result. UnresolvedFragments is
// intentionally empty until a constrained fallback is installed; the initial
// direct matcher never invents a symbol for prose it cannot ground.
type Intent struct {
	RawText             string            `json:"raw_text"`
	ExtractedEntities   []ExtractedEntity `json:"extracted_entities"`
	UnresolvedFragments []string          `json:"unresolved_fragments"`
	InheritedIntent     bool              `json:"inherited_intent"`
	NoBaseline          bool              `json:"no_baseline"`
}

// FallbackResolver is the only permitted extension point for resolving prose
// that direct graph-grounded matching cannot identify. Implementations receive
// the fragment and the complete allowed vocabulary; they must return either an
// exact vocabulary member or an empty string to abstain.
type FallbackResolver interface {
	Resolve(fragment string, vocabulary []string) (name string, confidence float64, err error)
}

// ResolveFallback fails closed. A model/provider answer outside the supplied
// graph vocabulary is represented as an unresolved fragment, never as an
// invented entity in the manifest.
func ResolveFallback(fragment string, vocabulary []string, resolver FallbackResolver) (ExtractedEntity, string, error) {
	if resolver == nil {
		return ExtractedEntity{}, fragment, nil
	}
	name, confidence, err := resolver.Resolve(fragment, vocabulary)
	if err != nil {
		return ExtractedEntity{}, fragment, err
	}
	for _, candidate := range vocabulary {
		if name == candidate && confidence >= 0 && confidence <= 1 {
			return ExtractedEntity{Name: name, Method: "llm_fallback", Confidence: confidence}, "", nil
		}
	}
	return ExtractedEntity{}, fragment, nil
}

// ExtractDirectMentions finds exact, token-bounded qualified-symbol mentions
// in checkpoint text. This is deliberately conservative: "Verify" is not
// allowed to match "VerifyIntent", and a match is emitted once even when a
// transcript repeats it. Fuzzy matching and LLM fallback belong in later
// tickets because they need an explicit, testable abstention policy.
func ExtractDirectMentions(text string, symbols []Symbol) Intent {
	seen := map[string]bool{}
	var entities []ExtractedEntity
	for _, symbol := range symbols {
		name := strings.TrimSpace(symbol.QualifiedName)
		if name == "" || seen[name] || !containsSymbolMention(text, name) {
			continue
		}
		seen[name] = true
		entities = append(entities, ExtractedEntity{Name: name, Method: "regex", Confidence: 1})
	}
	sort.Slice(entities, func(i, j int) bool { return entities[i].Name < entities[j].Name })
	return Intent{RawText: text, ExtractedEntities: entities, UnresolvedFragments: []string{}}
}

// DeclareIntent recognizes the deliberately narrow class of follow-up prompts
// that contain no actionable declaration by themselves. When a caller has
// supplied the immediately preceding checkpoint transcript, it can preserve
// the graph-grounded declaration while recording that it was inherited. It
// never silently searches arbitrary older sessions for a plausible claim.
func DeclareIntent(text, baselineText string, symbols []Symbol, noBaseline bool) Intent {
	intent := ExtractDirectMentions(text, symbols)
	intent.NoBaseline = noBaseline
	if len(intent.ExtractedEntities) != 0 || !isVagueFollowUp(text) || baselineText == "" {
		return intent
	}
	inherited := ExtractDirectMentions(baselineText, symbols)
	if len(inherited.ExtractedEntities) == 0 {
		return intent
	}
	inherited.RawText = text
	inherited.InheritedIntent = true
	inherited.NoBaseline = noBaseline
	return inherited
}

func isVagueFollowUp(text string) bool {
	normalized := strings.Trim(strings.ToLower(strings.TrimSpace(text)), ".!?")
	return normalized == "continue" || normalized == "fix it"
}

func containsSymbolMention(text, symbol string) bool {
	for offset := 0; ; {
		index := strings.Index(text[offset:], symbol)
		if index < 0 {
			return false
		}
		index += offset
		end := index + len(symbol)
		beforeOK := index == 0 || !symbolByte(text[index-1])
		afterOK := end == len(text) || !symbolByte(text[end]) ||
			(text[end] == '.' && (end+1 == len(text) || !symbolByte(text[end+1])))
		if beforeOK && afterOK {
			return true
		}
		offset = end
	}
}

func symbolByte(value byte) bool {
	return value == '_' || value == '.' ||
		(value >= '0' && value <= '9') ||
		(value >= 'A' && value <= 'Z') ||
		(value >= 'a' && value <= 'z')
}
