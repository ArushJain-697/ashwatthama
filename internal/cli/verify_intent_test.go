package cli

import (
	"bytes"
	"context"
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

type testFidelityAdapter struct{}

func (testFidelityAdapter) CheckpointTranscript(_ context.Context, _ string) (string, error) {
	return "Implement cli.VerifyIntent and cli.Missing.", nil
}

func (testFidelityAdapter) CheckpointCommit(_ context.Context, checkpoint string) (string, error) {
	if checkpoint == "checkpoint-0" {
		return "base", nil
	}
	return "head", nil
}

func (testFidelityAdapter) Graph(context.Context) (fidelity.GraphSnapshot, error) {
	return fidelity.GraphSnapshot{
		Symbols: []fidelity.Symbol{{QualifiedName: "cli.VerifyIntent"}, {QualifiedName: "cli.Missing"}},
	}, nil
}

func (testFidelityAdapter) ChangedEntities(_ context.Context, _, _ string) ([]fidelity.ChangedEntity, error) {
	return []fidelity.ChangedEntity{{Name: "VerifyIntent", FilePath: "internal/cli/verify_intent.go", StartLine: 1}}, nil
}

func TestVerifyIntentReportsPartialPipeline(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "fidelity.config.yaml"), []byte(verifyIntentConfig), 0o600); err != nil {
		t.Fatal(err)
	}
	var output, errOutput bytes.Buffer
	err := runVerifyIntentWithAdapter(t.Context(), Options{Env: EntireEnv{RepoRoot: repo}, Stdout: &output, Stderr: &errOutput}, []string{"checkpoint-1", "--base", "checkpoint-0"}, testFidelityAdapter{})
	if err != nil {
		t.Fatal(err)
	}
	var response verifyIntentResponse
	if err := json.Unmarshal(output.Bytes(), &response); err != nil {
		t.Fatalf("stdout was not clean JSON: %v\n%s", err, output.String())
	}
	want := []fidelity.Stage{fidelity.StageCapture, fidelity.StageDeclare, fidelity.StageObserve, fidelity.StageReconcile, fidelity.StageManifest}
	if response.Status != "ok" || !reflect.DeepEqual(response.Stages, want) {
		t.Fatalf("response = %#v, want ordered pipeline stages %#v", response, want)
	}
	if !response.Preflight.RequiredRelationsAvailable || response.Preflight.ConfidenceModel == "" {
		t.Fatalf("response did not include a usable preflight report: %#v", response.Preflight)
	}
	if len(response.Verification.Reconciliation.Confirmed) != 1 || len(response.Verification.Reconciliation.DeclaredUnimplemented) != 1 {
		t.Fatalf("response did not include direct reconciliation: %#v", response.Verification)
	}
	if response.VerdictPath == "" || response.ReportPath == "" {
		t.Fatalf("response did not report where verdict.json/VERIFY_REPORT.md were written: %#v", response)
	}
	if _, err := os.Stat(response.VerdictPath); err != nil {
		t.Fatalf("verdict.json was not written: %v", err)
	}
	if _, err := os.Stat(response.ReportPath); err != nil {
		t.Fatalf("VERIFY_REPORT.md was not written: %v", err)
	}
	if errOutput.Len() == 0 {
		t.Fatal("terminal report was not printed to stderr")
	}
}
