package fidelity

import (
	"context"
	"reflect"
	"testing"
)

type pipelineAdapter struct {
	transcripts map[string]string
	commits     map[string]string
	symbols     []Symbol
	changes     []ChangedEntity
	base        string
	head        string
}

func (adapter *pipelineAdapter) CheckpointTranscript(_ context.Context, id string) (string, error) {
	return adapter.transcripts[id], nil
}

func (adapter *pipelineAdapter) CheckpointCommit(_ context.Context, id string) (string, error) {
	return adapter.commits[id], nil
}

func (adapter *pipelineAdapter) Graph(context.Context) (GraphSnapshot, error) {
	return GraphSnapshot{Symbols: adapter.symbols}, nil
}

func (adapter *pipelineAdapter) ChangedEntities(_ context.Context, base, head string) ([]ChangedEntity, error) {
	adapter.base, adapter.head = base, head
	return adapter.changes, nil
}

type pipelineStages struct{ stages []Stage }

func (stages *pipelineStages) RunStage(_ context.Context, stage Stage, _ Request) error {
	stages.stages = append(stages.stages, stage)
	return nil
}

func TestPipelineGroundsTranscriptAndReconcilesCheckpointRange(t *testing.T) {
	adapter := &pipelineAdapter{
		transcripts: map[string]string{"current": "Implement service.VerifyIntent and leave service.Missing for later."},
		commits:     map[string]string{"current": "head", "base": "base"},
		symbols: []Symbol{
			{QualifiedName: "service.VerifyIntent"},
			{QualifiedName: "service.Missing"},
		},
		changes: []ChangedEntity{{Name: "VerifyIntent", FilePath: "service/intent.go", StartLine: 24}},
	}
	stages := &pipelineStages{}
	result, err := (Pipeline{Adapter: adapter, Runner: stages}).Run(t.Context(), Request{CheckpointID: "current", BaseCheckpointID: "base"})
	if err != nil {
		t.Fatal(err)
	}
	wantStages := []Stage{StageCapture, StageDeclare, StageObserve, StageReconcile, StageManifest}
	if !reflect.DeepEqual(stages.stages, wantStages) {
		t.Fatalf("stages = %#v, want %#v", stages.stages, wantStages)
	}
	if adapter.base != "base" || adapter.head != "head" {
		t.Fatalf("diff range = %q..%q, want base..head", adapter.base, adapter.head)
	}
	if len(result.Reconciliation.Confirmed) != 1 || result.Reconciliation.Confirmed[0].Entity != "service.VerifyIntent" {
		t.Fatalf("confirmed = %#v", result.Reconciliation.Confirmed)
	}
	if len(result.Reconciliation.DeclaredUnimplemented) != 1 || result.Reconciliation.DeclaredUnimplemented[0].Entity != "service.Missing" {
		t.Fatalf("unimplemented = %#v", result.Reconciliation.DeclaredUnimplemented)
	}
}

func TestPipelineMarksMissingBaseCheckpoint(t *testing.T) {
	adapter := &pipelineAdapter{
		transcripts: map[string]string{"current": ""},
		commits:     map[string]string{"current": "head"},
	}
	result, err := (Pipeline{Adapter: adapter}).Run(t.Context(), Request{CheckpointID: "current"})
	if err != nil {
		t.Fatal(err)
	}
	if !result.NoBaseline || adapter.base != emptyTreeObjectID || adapter.head != "head" {
		t.Fatalf("result/base/head = %#v, %q, %q", result, adapter.base, adapter.head)
	}
}

func TestPipelineInheritsVagueIntentFromExplicitBaseline(t *testing.T) {
	adapter := &pipelineAdapter{
		transcripts: map[string]string{
			"current": "fix it",
			"base":    "Implement service.VerifyIntent.",
		},
		commits: map[string]string{"current": "head", "base": "base"},
		symbols: []Symbol{{QualifiedName: "service.VerifyIntent"}},
	}
	result, err := (Pipeline{Adapter: adapter}).Run(t.Context(), Request{CheckpointID: "current", BaseCheckpointID: "base"})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Intent.InheritedIntent || len(result.Intent.ExtractedEntities) != 1 {
		t.Fatalf("intent = %#v", result.Intent)
	}
}
