package fidelity

import (
	"context"
	"fmt"
)

// Verification is the structured result produced before renderer-specific
// verdict.json fields are added. Its inputs all pass through Adapter, keeping
// the CLI, parent Entire checkpoint store, and semantic provider separated.
type Verification struct {
	CheckpointID     string               `json:"checkpoint_id"`
	BaseCheckpointID string               `json:"base_checkpoint_id,omitempty"`
	HeadCommit       string               `json:"head_commit"`
	BaseCommit       string               `json:"base_commit"`
	NoBaseline       bool                 `json:"no_baseline"`
	Intent           Intent               `json:"intent"`
	Changes          []ChangedEntity      `json:"changes"`
	Reconciliation   DirectReconciliation `json:"reconciliation"`
}

// Pipeline runs the five stable Fidelity stages against a real adapter. A
// StageRunner remains optional instrumentation: it lets callers record or
// render stage progress without allowing those concerns to read external data.
type Pipeline struct {
	Adapter Adapter
	Runner  StageRunner
}

func (pipeline Pipeline) Run(ctx context.Context, request Request) (Verification, error) {
	if pipeline.Adapter == nil {
		return Verification{}, fmt.Errorf("Fidelity pipeline requires an adapter")
	}
	runStage := func(stage Stage) error {
		if pipeline.Runner == nil {
			return nil
		}
		if err := pipeline.Runner.RunStage(ctx, stage, request); err != nil {
			return fmt.Errorf("Fidelity %s stage: %w", stage, err)
		}
		return nil
	}

	if err := runStage(StageCapture); err != nil {
		return Verification{}, err
	}
	transcript, err := pipeline.Adapter.CheckpointTranscript(ctx, request.CheckpointID)
	if err != nil {
		return Verification{}, fmt.Errorf("Fidelity capture stage: %w", err)
	}

	if err := runStage(StageDeclare); err != nil {
		return Verification{}, err
	}
	symbols, err := pipeline.Adapter.SymbolDictionary(ctx)
	if err != nil {
		return Verification{}, fmt.Errorf("Fidelity declare stage: %w", err)
	}
	noBaseline := request.BaseCheckpointID == ""
	var baselineTranscript string
	if !noBaseline {
		baselineTranscript, err = pipeline.Adapter.CheckpointTranscript(ctx, request.BaseCheckpointID)
		if err != nil {
			return Verification{}, fmt.Errorf("Fidelity declare stage: %w", err)
		}
	}
	intent := DeclareIntent(transcript, baselineTranscript, symbols, noBaseline)

	if err := runStage(StageObserve); err != nil {
		return Verification{}, err
	}
	head, err := pipeline.Adapter.CheckpointCommit(ctx, request.CheckpointID)
	if err != nil {
		return Verification{}, fmt.Errorf("Fidelity observe stage: %w", err)
	}
	base := emptyTreeObjectID
	if !noBaseline {
		base, err = pipeline.Adapter.CheckpointCommit(ctx, request.BaseCheckpointID)
		if err != nil {
			return Verification{}, fmt.Errorf("Fidelity observe stage: %w", err)
		}
	}
	changes, err := pipeline.Adapter.ChangedEntities(ctx, base, head)
	if err != nil {
		return Verification{}, fmt.Errorf("Fidelity observe stage: %w", err)
	}

	if err := runStage(StageReconcile); err != nil {
		return Verification{}, err
	}
	reconciliation := ReconcileDirectClaims(intent, changes)

	if err := runStage(StageManifest); err != nil {
		return Verification{}, err
	}
	return Verification{
		CheckpointID: request.CheckpointID, BaseCheckpointID: request.BaseCheckpointID,
		HeadCommit: head, BaseCommit: base, NoBaseline: noBaseline,
		Intent: intent, Changes: changes, Reconciliation: reconciliation,
	}, nil
}

// emptyTreeObjectID is Git's canonical empty tree. Using it for a first
// checkpoint makes the missing baseline explicit and avoids pretending a
// repository's historical parent commit was an earlier agent declaration.
const emptyTreeObjectID = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
