package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/entireio/entire-graph/internal/fidelity"
)

const verifyIntentConfig = `blast_radius_hops: 2
edge_classes:
  deterministic_confidence_gte: 0.9
  advisory_below_confidence: 0.9
tiers:
  confirmed: {}
  declared_unimplemented: {}
  expected_blast_radius:
    max_hops: 2
  undeclared_scope_creep: {}
  advisory_low_confidence: {}
extraction:
  fallback_llm: true
  llm_vocabulary_source: graph_snapshot
module_granularity: true
verdict_thresholds:
  review_required_if_scope_creep_gte: 1
databricks:
  calibration_enabled: false
  score_source: batch
  gate:
    mode: fallback
    alpha_target: 0.05
  ai_search:
    enabled: false
    index_name: null
`

func TestParseVerifyIntentFlags(t *testing.T) {
	flags, checkpoint, err := parseVerifyIntentFlags([]string{"checkpoint-1", "--base", "checkpoint-0", "--config", "policy.yaml", "--out", "out", "--repo", "repo"})
	if err != nil {
		t.Fatal(err)
	}
	if checkpoint != "checkpoint-1" || flags.BaseCheckpointID != "checkpoint-0" || flags.ConfigPath != "policy.yaml" || flags.OutputDirectory != "out" || flags.Repo != "repo" {
		t.Fatalf("parseVerifyIntentFlags = %#v, %q", flags, checkpoint)
	}
}

func TestVerifyIntentSkeletonReportsOrderedStages(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "fidelity.config.yaml"), []byte(verifyIntentConfig), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	err := Run(t.Context(), Options{Env: EntireEnv{RepoRoot: repo}, Stdout: &output}, []string{"verify-intent", "checkpoint-1"})
	if err != nil {
		t.Fatal(err)
	}
	var response verifyIntentResponse
	if err := json.Unmarshal(output.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	want := []fidelity.Stage{fidelity.StageCapture, fidelity.StageDeclare, fidelity.StageObserve, fidelity.StageReconcile, fidelity.StageManifest}
	if response.Status != "skeleton" || !reflect.DeepEqual(response.Stages, want) {
		t.Fatalf("response = %#v, want ordered skeleton stages %#v", response, want)
	}
	if !response.Preflight.RequiredRelationsAvailable || response.Preflight.ConfidenceModel == "" {
		t.Fatalf("response did not include a usable preflight report: %#v", response.Preflight)
	}
}
