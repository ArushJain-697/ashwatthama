package fidelity

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const validConfig = `blast_radius_hops: 2
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

func TestLoadConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fidelity.config.yaml")
	if err := os.WriteFile(path, []byte(validConfig), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.BlastRadiusHops != 2 || config.EdgeClasses.DeterministicConfidenceGTE != 0.9 {
		t.Fatalf("unexpected config: %#v", config)
	}
	if config.Databricks.ScoreSource != "batch" || config.Databricks.Gate.AlphaTarget != 0.05 {
		t.Fatalf("unexpected Databricks config: %#v", config.Databricks)
	}
}

func TestLoadConfigRejectsUnsupportedGraphDepth(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fidelity.config.yaml")
	bad := "blast_radius_hops: 3\n"
	if err := os.WriteFile(path, []byte(bad), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("LoadConfig accepted a graph depth that neighbors cannot support")
	}
}

type stageRecorder struct{ stages []Stage }

func (r *stageRecorder) RunStage(_ context.Context, stage Stage, _ Request) error {
	r.stages = append(r.stages, stage)
	return nil
}

func TestRunSkeletonOrdersEveryStage(t *testing.T) {
	recorder := &stageRecorder{}
	if err := RunSkeleton(t.Context(), recorder, Request{}); err != nil {
		t.Fatal(err)
	}
	want := []Stage{StageCapture, StageDeclare, StageObserve, StageReconcile, StageManifest}
	if !reflect.DeepEqual(recorder.stages, want) {
		t.Fatalf("stages = %#v, want %#v", recorder.stages, want)
	}
}

func TestCurrentPreflightReportReflectsGraphCapabilities(t *testing.T) {
	report := CurrentPreflightReport()
	if report.Provider != "entire-graph" || report.ConfidenceModel == "" {
		t.Fatalf("preflight report = %#v", report)
	}
	for _, relation := range []string{"CALLS", "DATA_FLOWS", "TESTS", "HANDLES_ROUTE", "HANDLES_GRPC"} {
		if !report.RelationTypes[relation] {
			t.Fatalf("preflight did not report supported relation %q: %#v", relation, report.RelationTypes)
		}
	}
}
