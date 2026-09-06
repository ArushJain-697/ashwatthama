package fidelity

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/entireio/entire-graph/internal/sem"
)

// Adapter is the only seam through which Fidelity's capture, declaration, and
// observation stages may obtain external data.  The graph plugin can provide
// snapshots and semantic diffs, while checkpoint transcript retrieval belongs
// to the parent Entire/Brain layer; keeping both behind one interface prevents
// that distinction leaking into reconciliation and rendering.
type Adapter interface {
	CheckpointTranscript(context.Context, string) (string, error)
	SymbolDictionary(context.Context) ([]Symbol, error)
	ChangedEntities(context.Context, string, string) ([]ChangedEntity, error)
}

type Symbol struct {
	ID            string
	QualifiedName string
	FilePath      string
	StartLine     int
}

type ChangedEntity struct {
	ID        string
	Name      string
	FilePath  string
	StartLine int
}

// NativeGraphAdapter is the implementation used by the Fidelity CLI. It uses
// the in-process semantic provider for graph facts and Entire's documented
// checkpoint command for stored transcript bytes.
type NativeGraphAdapter struct {
	Repo            string
	ProviderVersion string
}

func (adapter NativeGraphAdapter) CheckpointTranscript(ctx context.Context, checkpointID string) (string, error) {
	if strings.TrimSpace(checkpointID) == "" {
		return "", fmt.Errorf("checkpoint transcript requires a checkpoint ID")
	}
	// Do not use a shell here. Checkpoint IDs are external input, and passing
	// them as a single argument keeps them data even if malformed. --transcript
	// is the purpose-built CLI export: unlike --full it emits stored JSONL only,
	// without a presentation layer that Fidelity would have to scrape.
	command := exec.CommandContext(ctx, "entire", "checkpoint", "explain", "--checkpoint", checkpointID, "--transcript", "--no-pager")
	command.Dir = adapter.Repo
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("read checkpoint transcript %q with Entire CLI: %w", checkpointID, err)
	}
	if len(output) == 0 {
		return "", fmt.Errorf("checkpoint transcript %q was empty", checkpointID)
	}
	return string(output), nil
}

func (adapter NativeGraphAdapter) SymbolDictionary(ctx context.Context) ([]Symbol, error) {
	snapshot, err := sem.BuildProviderSnapshotWithOptions(ctx, adapter.Repo, adapter.ProviderVersion, sem.ProviderSnapshotOptions{
		NoNetwork: true,
		Profile:   sem.ProfileFull,
	})
	if err != nil {
		return nil, fmt.Errorf("build graph snapshot: %w", err)
	}
	symbols := make([]Symbol, 0, len(snapshot.Symbols))
	for _, symbol := range snapshot.Symbols {
		symbols = append(symbols, Symbol{
			ID: symbol.ID, QualifiedName: symbol.QualifiedName,
			FilePath: symbol.FilePath, StartLine: symbol.StartLine,
		})
	}
	return symbols, nil
}

func (adapter NativeGraphAdapter) ChangedEntities(ctx context.Context, base, head string) ([]ChangedEntity, error) {
	result, err := sem.AnalyzeGitRange(ctx, adapter.Repo, base, head, nil)
	if err != nil {
		return nil, fmt.Errorf("analyze graph diff: %w", err)
	}
	var entities []ChangedEntity
	for _, file := range result.Files {
		for _, change := range file.Changes {
			line := change.AfterStartLine
			if line == 0 {
				line = change.BeforeStartLine
			}
			entities = append(entities, ChangedEntity{
				ID:   fmt.Sprintf("%s:%s:%d", file.Path, change.Name, line),
				Name: change.Name, FilePath: file.Path, StartLine: line,
			})
		}
	}
	return entities, nil
}

// UnavailableAdapter is the safe default for the command skeleton.  It makes
// the missing transcript capability explicit instead of pretending that a Git
// trailer contains the agent's stated intent.
type UnavailableAdapter struct{}

func (UnavailableAdapter) CheckpointTranscript(_ context.Context, checkpointID string) (string, error) {
	return "", fmt.Errorf("checkpoint transcript %q is unavailable: Entire Graph resolves checkpoint trailers to commits but does not expose transcript content", checkpointID)
}

func (UnavailableAdapter) SymbolDictionary(context.Context) ([]Symbol, error) {
	return nil, fmt.Errorf("graph adapter is not wired yet")
}

func (UnavailableAdapter) ChangedEntities(_ context.Context, _, _ string) ([]ChangedEntity, error) {
	return nil, fmt.Errorf("graph adapter is not wired yet")
}

// Stage identifies a stable pipeline boundary.  It allows the CLI skeleton to
// prove ordering today and lets later tickets replace one stage at a time.
type Stage string

const (
	StageCapture   Stage = "capture"
	StageDeclare   Stage = "declare"
	StageObserve   Stage = "observe"
	StageReconcile Stage = "reconcile"
	StageManifest  Stage = "manifest"
)

type Request struct {
	CheckpointID     string
	BaseCheckpointID string
	Config           Config
	OutputDirectory  string
}

type StageRunner interface {
	RunStage(context.Context, Stage, Request) error
}

// RunSkeleton is deliberately useful even before transcript access is wired:
// it pins the command's five-stage order and makes the external boundary
// injectable in tests.
func RunSkeleton(ctx context.Context, runner StageRunner, request Request) error {
	for _, stage := range []Stage{StageCapture, StageDeclare, StageObserve, StageReconcile, StageManifest} {
		if err := runner.RunStage(ctx, stage, request); err != nil {
			return fmt.Errorf("Fidelity %s stage: %w", stage, err)
		}
	}
	return nil
}
