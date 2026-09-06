package fidelity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func fixtureVerification() Verification {
	graph := fixtureGraph()
	changes := []ChangedEntity{
		{Name: "handleTimeout", FilePath: "auth/router.go", StartLine: 10},
		{Name: "invalidate", FilePath: "auth/cache.go", StartLine: 20},
		{Name: "chargeRetry", FilePath: "billing/controller.go", StartLine: 30},
		{Name: "route", FilePath: "gateway/route.go", StartLine: 40},
		{Name: "dispatch", FilePath: "plugins/registry.go", StartLine: 50},
	}
	intent := Intent{
		RawText:           "Fix the timeout bug in AuthRouter.handleTimeout",
		ExtractedEntities: []ExtractedEntity{{Name: "AuthRouter.handleTimeout", Method: "regex", Confidence: 1}},
	}
	return Verification{
		CheckpointID: "checkpoint-1", HeadCommit: "head", BaseCommit: "base",
		Intent: intent, Changes: changes,
		Reconciliation: Reconcile(intent, changes, graph, DefaultConfig()),
	}
}

func TestBuildVerdictLabelsReviewRequiredOnScopeCreep(t *testing.T) {
	verdict := BuildVerdict(fixtureVerification(), CurrentPreflightReport(), DefaultConfig(), time.Unix(0, 0))
	if verdict.Summary.VerdictLabel != "REVIEW_REQUIRED" {
		t.Fatalf("verdict_label = %q, want REVIEW_REQUIRED (one scope-creep entity present)", verdict.Summary.VerdictLabel)
	}
	if verdict.Summary.UnverifiableCoverageCount != 1 {
		t.Fatalf("unverifiable_coverage_count = %d, want 1", verdict.Summary.UnverifiableCoverageCount)
	}
	if len(verdict.GraphCapabilitiesVerified) == 0 {
		t.Fatal("graph_capabilities_verified was empty")
	}
}

func TestWriteVerdictProducesValidJSON(t *testing.T) {
	dir := t.TempDir()
	verdict := BuildVerdict(fixtureVerification(), CurrentPreflightReport(), DefaultConfig(), time.Now())
	path, err := WriteVerdict(dir, verdict)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(dir, "verdict.json") {
		t.Fatalf("path = %q", path)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var roundTrip Verdict
	if err := json.Unmarshal(content, &roundTrip); err != nil {
		t.Fatalf("verdict.json did not parse: %v", err)
	}
	if roundTrip.CheckpointID != "checkpoint-1" {
		t.Fatalf("round-tripped checkpoint_id = %q", roundTrip.CheckpointID)
	}
}
