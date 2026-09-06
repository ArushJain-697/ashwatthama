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

func containsSymbolMention(text, symbol string) bool {
	for offset := 0; ; {
		index := strings.Index(text[offset:], symbol)
		if index < 0 {
			return false
		}
		index += offset
		end := index + len(symbol)
		beforeOK := index == 0 || !symbolByte(text[index-1])
		afterOK := end == len(text) || !symbolByte(text[end])
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
