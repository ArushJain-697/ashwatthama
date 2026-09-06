// Package fidelity holds the intentionally small, provider-independent core of
// the Fidelity intent-versus-implementation verifier.  Keeping it outside cli
// means its policy and orchestration can be tested without executing an Entire
// command or reading a repository.
package fidelity

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is the stable policy seam used by reconciliation.  It deliberately
// describes confidence thresholds rather than claiming that Entire Graph emits
// an extracted/inferred/ambiguous tri-state; graph relations carry numeric
// confidence, resolution, warning codes, and evidence.
type Config struct {
	BlastRadiusHops   int
	EdgeClasses       EdgeClasses
	Extraction        Extraction
	ModuleGranularity bool
	VerdictThresholds VerdictThresholds
	// Coverage governs the Curveball response (Track 2, #33-#37): whether a
	// changed entity's absence of graph edges means "confirmed disconnected"
	// or "the graph never tried to look here." Orthogonal to EdgeClasses,
	// which classifies confidence in an edge that DOES exist.
	Coverage   Coverage
	Databricks Databricks
}

type Coverage struct {
	// DetectVia selects the detection mechanism: "capabilities_api" reads the
	// provider's own completeness signals (language tier, partial failures,
	// file presence) before falling back to edge counting; "heuristic_zero_edge"
	// uses only the zero-edge count. #7 confirmed capabilities_api is real for
	// this provider, so it is the default.
	DetectVia string
	// TreatZeroEdgeAs is the coverage value assigned when an entity has code
	// presence but zero extracted-class edges. "partial" is the conservative
	// default per the curveball's "never present incomplete as certain."
	TreatZeroEdgeAs string
	// FallbackVerificationPath is stamped on every unverifiable_coverage entry
	// as the baseline, always-present safety guarantee (#35): it must hold
	// with Databricks fully disabled.
	FallbackVerificationPath string
}

type EdgeClasses struct {
	DeterministicConfidenceGTE float64
	AdvisoryBelowConfidence    float64
}

type Extraction struct {
	FallbackLLM         bool
	LLMVocabularySource string
}

type VerdictThresholds struct {
	ReviewRequiredIfScopeCreepGTE           int
	ReviewRequiredIfUnverifiableCoverageGTE int
}

type Databricks struct {
	CalibrationEnabled bool
	ScoreSource        string
	Gate               CalibrationGate
	AISearch           AISearch
	// CoverageCorroborationEnabled turns on §12.15's enhancement of the
	// unverifiable_coverage baseline: a second query against the same AI
	// Search index used for advisory-edge corroboration. Never a substitute
	// for the baseline (#35), which must hold with this false.
	CoverageCorroborationEnabled bool
}

type CalibrationGate struct {
	Mode        string
	AlphaTarget float64
}

type AISearch struct {
	Enabled   bool
	IndexName string
}

func DefaultConfig() Config {
	return Config{
		BlastRadiusHops: 2,
		EdgeClasses: EdgeClasses{
			DeterministicConfidenceGTE: 0.9,
			AdvisoryBelowConfidence:    0.9,
		},
		Extraction:        Extraction{FallbackLLM: true, LLMVocabularySource: "graph_snapshot"},
		ModuleGranularity: true,
		VerdictThresholds: VerdictThresholds{
			ReviewRequiredIfScopeCreepGTE:           1,
			ReviewRequiredIfUnverifiableCoverageGTE: 1,
		},
		Coverage: Coverage{
			DetectVia:                "capabilities_api",
			TreatZeroEdgeAs:          "partial",
			FallbackVerificationPath: "manual_review_recommended",
		},
		Databricks: Databricks{
			ScoreSource: "batch",
			Gate:        CalibrationGate{Mode: "fallback", AlphaTarget: 0.05},
		},
	}
}

// LoadConfig reads Fidelity's deliberately narrow YAML policy file.  The
// parser accepts the scalar/mapping subset used by fidelity.config.yaml and
// rejects unknown keys rather than silently applying a typo as a default.  A
// full YAML dependency is unnecessary until configuration needs YAML features
// such as anchors, multiline scalars, or sequences.
func LoadConfig(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read Fidelity config: %w", err)
	}
	config := DefaultConfig()
	section := ""
	subsection := ""
	for lineNumber, line := range strings.Split(string(content), "\n") {
		raw := strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		if raw == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		parts := strings.SplitN(raw, ":", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("Fidelity config line %d: expected key: value", lineNumber+1)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if value == "" {
			switch indent {
			case 0:
				section, subsection = key, ""
			case 2:
				subsection = key
			default:
				return Config{}, fmt.Errorf("Fidelity config line %d: unsupported nesting", lineNumber+1)
			}
			continue
		}
		// A scalar at an enclosing indentation level is a sibling of any
		// earlier mapping (for example the tiers after expected_blast_radius),
		// not a member of that mapping.
		if indent == 0 {
			section, subsection = "", ""
		} else if indent == 2 {
			subsection = ""
		}
		if err := assignConfigValue(&config, section, subsection, key, value); err != nil {
			return Config{}, fmt.Errorf("Fidelity config line %d: %w", lineNumber+1, err)
		}
	}
	return config, validateConfig(config)
}

func assignConfigValue(config *Config, section, subsection, key, value string) error {
	value = strings.Trim(value, "\"'")
	parseBool := func() (bool, error) { return strconv.ParseBool(value) }
	parseInt := func() (int, error) { return strconv.Atoi(value) }
	parseFloat := func() (float64, error) { return strconv.ParseFloat(value, 64) }
	switch {
	case section == "" && key == "blast_radius_hops":
		v, err := parseInt()
		config.BlastRadiusHops = v
		return err
	case section == "" && key == "module_granularity":
		v, err := parseBool()
		config.ModuleGranularity = v
		return err
	case section == "edge_classes" && key == "deterministic_confidence_gte":
		v, err := parseFloat()
		config.EdgeClasses.DeterministicConfidenceGTE = v
		return err
	case section == "edge_classes" && key == "advisory_below_confidence":
		v, err := parseFloat()
		config.EdgeClasses.AdvisoryBelowConfidence = v
		return err
	case section == "tiers" && subsection == "" &&
		(key == "confirmed" || key == "declared_unimplemented" || key == "undeclared_scope_creep" || key == "advisory_low_confidence"):
		if value != "{}" {
			return fmt.Errorf("tier %q must be an empty mapping", key)
		}
		return nil
	case section == "tiers" && subsection == "expected_blast_radius" && key == "max_hops":
		v, err := parseInt()
		if err != nil {
			return err
		}
		if v != config.BlastRadiusHops {
			return fmt.Errorf("tiers.expected_blast_radius.max_hops must match blast_radius_hops")
		}
		return nil
	case section == "extraction" && key == "fallback_llm":
		v, err := parseBool()
		config.Extraction.FallbackLLM = v
		return err
	case section == "extraction" && key == "llm_vocabulary_source":
		config.Extraction.LLMVocabularySource = value
		return nil
	case section == "verdict_thresholds" && key == "review_required_if_scope_creep_gte":
		v, err := parseInt()
		config.VerdictThresholds.ReviewRequiredIfScopeCreepGTE = v
		return err
	case section == "verdict_thresholds" && key == "review_required_if_unverifiable_coverage_gte":
		v, err := parseInt()
		config.VerdictThresholds.ReviewRequiredIfUnverifiableCoverageGTE = v
		return err
	case section == "coverage" && key == "detect_via":
		config.Coverage.DetectVia = value
		return nil
	case section == "coverage" && key == "treat_zero_edge_as":
		config.Coverage.TreatZeroEdgeAs = value
		return nil
	case section == "coverage" && key == "fallback_verification_path":
		config.Coverage.FallbackVerificationPath = value
		return nil
	case section == "databricks" && subsection == "" && key == "coverage_corroboration_enabled":
		v, err := parseBool()
		config.Databricks.CoverageCorroborationEnabled = v
		return err
	case section == "databricks" && subsection == "" && key == "calibration_enabled":
		v, err := parseBool()
		config.Databricks.CalibrationEnabled = v
		return err
	case section == "databricks" && subsection == "" && key == "score_source":
		config.Databricks.ScoreSource = value
		return nil
	case section == "databricks" && subsection == "gate" && key == "mode":
		config.Databricks.Gate.Mode = value
		return nil
	case section == "databricks" && subsection == "gate" && key == "alpha_target":
		v, err := parseFloat()
		config.Databricks.Gate.AlphaTarget = v
		return err
	case section == "databricks" && subsection == "ai_search" && key == "enabled":
		v, err := parseBool()
		config.Databricks.AISearch.Enabled = v
		return err
	case section == "databricks" && subsection == "ai_search" && key == "index_name":
		if value != "null" {
			config.Databricks.AISearch.IndexName = value
		}
		return nil
	default:
		return fmt.Errorf("unknown key %q in %s.%s", key, section, subsection)
	}
}

func validateConfig(config Config) error {
	if config.BlastRadiusHops < 1 || config.BlastRadiusHops > 2 {
		return fmt.Errorf("blast_radius_hops must be 1 or 2; Entire Graph neighbors supports at most 2")
	}
	if config.EdgeClasses.DeterministicConfidenceGTE < 0 || config.EdgeClasses.DeterministicConfidenceGTE > 1 ||
		config.EdgeClasses.AdvisoryBelowConfidence < 0 || config.EdgeClasses.AdvisoryBelowConfidence > 1 {
		return fmt.Errorf("confidence thresholds must be between 0 and 1")
	}
	if config.Extraction.LLMVocabularySource != "graph_snapshot" {
		return fmt.Errorf("llm_vocabulary_source must be graph_snapshot")
	}
	if config.Databricks.ScoreSource != "batch" && config.Databricks.ScoreSource != "serving_endpoint" {
		return fmt.Errorf("score_source must be batch or serving_endpoint")
	}
	if config.Databricks.Gate.Mode != "fallback" && config.Databricks.Gate.Mode != "icalm_crc" {
		return fmt.Errorf("gate.mode must be fallback or icalm_crc")
	}
	if config.Databricks.Gate.AlphaTarget <= 0 || config.Databricks.Gate.AlphaTarget >= 1 {
		return fmt.Errorf("gate.alpha_target must be between 0 and 1")
	}
	if config.Coverage.DetectVia != "capabilities_api" && config.Coverage.DetectVia != "heuristic_zero_edge" {
		return fmt.Errorf("coverage.detect_via must be capabilities_api or heuristic_zero_edge")
	}
	if config.Coverage.TreatZeroEdgeAs != "partial" && config.Coverage.TreatZeroEdgeAs != "full" {
		return fmt.Errorf("coverage.treat_zero_edge_as must be partial or full")
	}
	if config.Coverage.FallbackVerificationPath == "" {
		return fmt.Errorf("coverage.fallback_verification_path must not be empty")
	}
	return nil
}
