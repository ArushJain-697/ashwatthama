package fidelity

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/entireio/entire-graph/internal/gitutil"
	"github.com/entireio/entire-graph/internal/sem"
)

// Adapter is the only seam through which Fidelity's capture, declaration, and
// observation stages may obtain external data.  The graph plugin can provide
// snapshots and semantic diffs, while checkpoint transcript retrieval belongs
// to the parent Entire/Brain layer; keeping both behind one interface prevents
// that distinction leaking into reconciliation and rendering.
type Adapter interface {
	CheckpointTranscript(context.Context, string) (string, error)
	CheckpointCommit(context.Context, string) (string, error)
	// Graph returns the whole-repo symbol dictionary, its relations, and the
	// provider's own completeness signals, from a single snapshot build. Stage
	// 2 (extraction vocabulary) and Stage 4 (reachability, coverage-confidence)
	// both read it, so the graph is built once, not once per stage.
	Graph(context.Context) (GraphSnapshot, error)
	ChangedEntities(context.Context, string, string) ([]ChangedEntity, error)
}

type Symbol struct {
	ID            string
	QualifiedName string
	FilePath      string
	StartLine     int
}

// GraphSnapshot is the whole-repo graph facts Stage 4 reconciles against.
// LanguageTiers and PartialFailureFiles are the provider's own completeness
// signals (§7, §33) — the "capabilities_api" detection mechanism reads them
// instead of relying solely on a zero-edge heuristic.
type GraphSnapshot struct {
	Symbols []Symbol
	// Relations is the whole-repo relation set, used for N-hop reachability
	// (#21) and coverage-confidence's extracted-edge count (#33).
	Relations []GraphRelation
	// FileLanguages maps a repository-relative file path to the language the
	// provider recorded for it.
	FileLanguages map[string]string
	// LanguageTiers classifies each language as "semantic" or
	// "inventory-only" (the provider's own tri-state-adjacent signal: an
	// inventory-only file was never a candidate for relation extraction at
	// all, which is exactly coverage-confidence's "partial" case).
	LanguageTiers map[string]string
	// PartialFailureFiles names files with a real parse/read gap (not an
	// intentional skip such as E_FILE_TOO_LARGE) — the region the provider
	// itself flagged as incompletely analyzed.
	PartialFailureFiles map[string]bool
	// CompletenessLevel is the provider's own ok/degraded/unsafe verdict for
	// this snapshot, carried through unchanged for the verdict manifest.
	CompletenessLevel string
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

// CheckpointCommit resolves the checkpoint trailer before Stage 3 passes a
// concrete, local Git range to the semantic diff engine. The checkpoint ID is
// intentionally never treated as a Git revision directly.
func (adapter NativeGraphAdapter) CheckpointCommit(ctx context.Context, checkpointID string) (string, error) {
	if strings.TrimSpace(checkpointID) == "" {
		return "", fmt.Errorf("checkpoint commit requires a checkpoint ID")
	}
	commit, err := gitutil.FindCommitWithCheckpoint(ctx, adapter.Repo, checkpointID)
	if err != nil {
		return "", fmt.Errorf("resolve checkpoint %q to commit: %w", checkpointID, err)
	}
	return commit, nil
}

// intentionalSkipCodes are partial-failure codes that mean the provider
// deliberately declined to parse a file (it is still recorded, still
// coverage: full-eligible if it turns out to have edges) rather than failed
// to. Mirrors internal/sem's own completenessFailureCount distinction so
// Fidelity's "unknown" coverage tracks the provider's own "degraded" signal,
// not a stricter or looser one.
var intentionalSkipCodes = map[string]bool{
	"E_FILE_TOO_LARGE": true,
	"E_MINIFIED":       true,
}

// Graph runs one panic-guarded snapshot build. The guard is Fidelity's answer
// to Build Map #20 at the seam it actually owns: entire-graph's own DATA_FLOWS
// boundary case (Issue #32) was not reproduced against this repository (see
// BUILDATHON.md), so this does not claim a targeted fix for it — it is a
// general defensive wrapper around the one call site Fidelity controls, so
// ANY panic from the underlying analyzer (that one included, if it ever
// fires) degrades to an error and an "unknown" coverage reading for the
// affected region instead of crashing verify-intent.
func (adapter NativeGraphAdapter) Graph(ctx context.Context) (snapshot GraphSnapshot, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("graph snapshot panicked (degraded to error, not crash): %v", r)
		}
	}()
	provider, buildErr := sem.BuildProviderSnapshotWithOptions(ctx, adapter.Repo, adapter.ProviderVersion, sem.ProviderSnapshotOptions{
		NoNetwork: true,
		Profile:   sem.ProfileFull,
	})
	if buildErr != nil {
		return GraphSnapshot{}, fmt.Errorf("build graph snapshot: %w", buildErr)
	}
	symbols := make([]Symbol, 0, len(provider.Symbols))
	for _, symbol := range provider.Symbols {
		symbols = append(symbols, Symbol{
			ID: symbol.ID, QualifiedName: symbol.QualifiedName,
			FilePath: symbol.FilePath, StartLine: symbol.StartLine,
		})
	}
	relations := make([]GraphRelation, 0, len(provider.Relations))
	for _, relation := range provider.Relations {
		relations = append(relations, GraphRelation{
			FromID: relation.FromID, ToID: relation.ToID, Type: relation.Type,
			Confidence: relation.Confidence, Resolution: relation.Resolution,
		})
	}
	fileLanguages := make(map[string]string, len(provider.Files))
	for _, file := range provider.Files {
		fileLanguages[file.Path] = file.Language
	}
	partialFailureFiles := map[string]bool{}
	for _, failure := range provider.Header.PartialFailures {
		if failure.FilePath != "" && !intentionalSkipCodes[failure.Code] {
			partialFailureFiles[failure.FilePath] = true
		}
	}
	return GraphSnapshot{
		Symbols: symbols, Relations: relations,
		FileLanguages: fileLanguages, LanguageTiers: provider.Header.LanguageTiers,
		PartialFailureFiles: partialFailureFiles, CompletenessLevel: provider.Header.Stats.CompletenessLevel,
	}, nil
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

func (UnavailableAdapter) CheckpointCommit(_ context.Context, checkpointID string) (string, error) {
	return "", fmt.Errorf("checkpoint commit %q is unavailable: Entire Graph needs the checkpoint trailer mapping", checkpointID)
}

func (UnavailableAdapter) Graph(context.Context) (GraphSnapshot, error) {
	return GraphSnapshot{}, fmt.Errorf("graph adapter is not wired yet")
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
