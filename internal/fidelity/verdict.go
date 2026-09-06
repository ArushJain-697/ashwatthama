package fidelity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Verdict is the locked verdict.json shape (#26, extended by #36 for the
// Curveball's coverage-confidence response). Every renderer (#27, #38-#40) is
// a pure function of this struct; nothing here reads the adapter or repo.
type Verdict struct {
	CheckpointID              string         `json:"checkpoint_id"`
	BaseCheckpointID          string         `json:"base_checkpoint_id,omitempty"`
	GeneratedAt               string         `json:"generated_at"`
	GraphCapabilitiesVerified []string       `json:"graph_capabilities_verified"`
	Intent                    VerdictIntent  `json:"intent"`
	Reconciliation            Reconciliation `json:"reconciliation"`
	Summary                   VerdictSummary `json:"summary"`
}

type VerdictIntent struct {
	RawPromptExcerpt    string            `json:"raw_prompt_excerpt"`
	ExtractedEntities   []ExtractedEntity `json:"extracted_entities"`
	UnresolvedFragments []string          `json:"unresolved_fragments"`
}

type VerdictSummary struct {
	TotalChangedEntities      int    `json:"total_changed_entities"`
	ScopeCreepCount           int    `json:"scope_creep_count"`
	UnimplementedCount        int    `json:"unimplemented_count"`
	AdvisoryEdgeCount         int    `json:"advisory_edge_count"`
	UnverifiableCoverageCount int    `json:"unverifiable_coverage_count"`
	VerdictLabel              string `json:"verdict_label"`
}

const rawPromptExcerptLimit = 200

// BuildVerdict composes Stage 5's manifest from Stage 1-4's Verification, the
// preflight capability report, and the request's config-driven thresholds.
// It performs no graph or transcript access of its own — the §13.1 firewall
// keeps rendering a pure function of already-computed facts.
func BuildVerdict(verification Verification, preflight PreflightReport, config Config, now time.Time) Verdict {
	reconciliation := verification.Reconciliation
	scopeCreepCount := len(reconciliation.UndeclaredScopeCreep)
	unverifiableCount := len(reconciliation.UnverifiableCoverage)
	label := "CLEAN"
	if scopeCreepCount >= config.VerdictThresholds.ReviewRequiredIfScopeCreepGTE && config.VerdictThresholds.ReviewRequiredIfScopeCreepGTE > 0 ||
		unverifiableCount >= config.VerdictThresholds.ReviewRequiredIfUnverifiableCoverageGTE && config.VerdictThresholds.ReviewRequiredIfUnverifiableCoverageGTE > 0 {
		label = "REVIEW_REQUIRED"
	}

	unresolved := verification.Intent.UnresolvedFragments
	if unresolved == nil {
		unresolved = []string{}
	}
	entities := verification.Intent.ExtractedEntities
	if entities == nil {
		entities = []ExtractedEntity{}
	}

	return Verdict{
		CheckpointID: verification.CheckpointID, BaseCheckpointID: verification.BaseCheckpointID,
		GeneratedAt:               now.UTC().Format(time.RFC3339),
		GraphCapabilitiesVerified: verifiedRelationTypes(preflight),
		Intent: VerdictIntent{
			RawPromptExcerpt:  excerpt(verification.Intent.RawText, rawPromptExcerptLimit),
			ExtractedEntities: entities, UnresolvedFragments: unresolved,
		},
		Reconciliation: reconciliation,
		Summary: VerdictSummary{
			TotalChangedEntities:      len(verification.Changes),
			ScopeCreepCount:           scopeCreepCount,
			UnimplementedCount:        len(reconciliation.DeclaredUnimplemented),
			AdvisoryEdgeCount:         len(reconciliation.AdvisoryLowConfidence),
			UnverifiableCoverageCount: unverifiableCount,
			VerdictLabel:              label,
		},
	}
}

func verifiedRelationTypes(preflight PreflightReport) []string {
	types := make([]string, 0, len(preflight.RelationTypes))
	for relation, present := range preflight.RelationTypes {
		if present {
			types = append(types, relation)
		}
	}
	sort.Strings(types)
	return types
}

func excerpt(text string, limit int) string {
	text = strings.TrimSpace(text)
	if len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}

// WriteVerdict writes verdict.json into outputDirectory, creating it if
// needed. This is the file the mandatory fixture test (#42) and every
// renderer read back.
func WriteVerdict(outputDirectory string, verdict Verdict) (string, error) {
	if err := os.MkdirAll(outputDirectory, 0o755); err != nil {
		return "", fmt.Errorf("create Fidelity output directory: %w", err)
	}
	path := filepath.Join(outputDirectory, "verdict.json")
	content, err := json.MarshalIndent(verdict, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode verdict.json: %w", err)
	}
	if err := os.WriteFile(path, append(content, '\n'), 0o644); err != nil {
		return "", fmt.Errorf("write verdict.json: %w", err)
	}
	return path, nil
}
