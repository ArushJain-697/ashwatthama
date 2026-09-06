package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/entireio/entire-graph/internal/fidelity"
	"github.com/entireio/entire-graph/internal/termsafe"
)

// verifyIntentFlags intentionally keeps the command surface small and stable.
type verifyIntentFlags struct {
	BaseCheckpointID string
	ConfigPath       string
	OutputDirectory  string
	Repo             string
	// CommunityMapPath is #57's optional enhancement input: a communities.json
	// produced by scripts/leiden_communities.py. Empty means skip annotation.
	CommunityMapPath string
}

type verifyIntentResponse struct {
	FormatVersion    int                      `json:"format_version"`
	Status           string                   `json:"status"`
	CheckpointID     string                   `json:"checkpoint_id"`
	BaseCheckpointID string                   `json:"base_checkpoint_id,omitempty"`
	ConfigPath       string                   `json:"config_path"`
	OutputDirectory  string                   `json:"output_directory"`
	Stages           []fidelity.Stage         `json:"stages"`
	Preflight        fidelity.PreflightReport `json:"preflight"`
	Verification     fidelity.Verification    `json:"verification"`
	Verdict          fidelity.Verdict         `json:"verdict"`
	VerdictPath      string                   `json:"verdict_path"`
	ReportPath       string                   `json:"report_path"`
	Message          string                   `json:"message"`
}

func parseVerifyIntentFlags(args []string) (verifyIntentFlags, string, error) {
	flags := verifyIntentFlags{ConfigPath: "fidelity.config.yaml", OutputDirectory: "fidelity-out"}
	var checkpointID string
	for index := 0; index < len(args); index++ {
		argument := args[index]
		value := func() (string, error) {
			index++
			if index >= len(args) {
				return "", fmt.Errorf("%s requires a value", argument)
			}
			return args[index], nil
		}
		var err error
		switch argument {
		case "--base":
			flags.BaseCheckpointID, err = value()
		case "--config":
			flags.ConfigPath, err = value()
		case "--out":
			flags.OutputDirectory, err = value()
		case "--repo":
			flags.Repo, err = value()
		case "--community-map":
			flags.CommunityMapPath, err = value()
		default:
			if strings.HasPrefix(argument, "-") {
				return flags, "", fmt.Errorf("verify-intent received unexpected argument %q", argument)
			}
			if checkpointID != "" {
				return flags, "", fmt.Errorf("verify-intent accepts exactly one checkpoint ID")
			}
			checkpointID = argument
		}
		if err != nil {
			return flags, "", err
		}
	}
	if checkpointID == "" {
		return flags, "", fmt.Errorf("verify-intent requires a checkpoint ID")
	}
	return flags, checkpointID, nil
}

type recordingFidelityStages struct{ stages []fidelity.Stage }

func (runner *recordingFidelityStages) RunStage(_ context.Context, stage fidelity.Stage, _ fidelity.Request) error {
	runner.stages = append(runner.stages, stage)
	return nil
}

func runVerifyIntent(ctx context.Context, opts Options, args []string) error {
	return runVerifyIntentWithAdapter(ctx, opts, args, fidelity.NativeGraphAdapter{})
}

func runVerifyIntentWithAdapter(ctx context.Context, opts Options, args []string, adapter fidelity.Adapter) error {
	flags, checkpointID, err := parseVerifyIntentFlags(args)
	if err != nil {
		return err
	}
	repo, err := resolveRepo(ctx, opts.Env, flags.Repo)
	if err != nil {
		return err
	}
	configPath := flags.ConfigPath
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(repo, configPath)
	}
	config, err := fidelity.LoadConfig(configPath)
	if err != nil {
		return err
	}
	runner := &recordingFidelityStages{}
	request := fidelity.Request{
		CheckpointID: checkpointID, BaseCheckpointID: flags.BaseCheckpointID,
		Config: config, OutputDirectory: flags.OutputDirectory,
	}
	if native, ok := adapter.(fidelity.NativeGraphAdapter); ok && native.Repo == "" {
		native.Repo = repo
		native.ProviderVersion = opts.Version
		adapter = native
	}
	verification, err := (fidelity.Pipeline{Adapter: adapter, Runner: runner}).Run(ctx, request)
	if err != nil {
		return err
	}
	if flags.CommunityMapPath != "" {
		// ponytail: this rebuilds the graph snapshot a second time rather than
		// threading it out of Pipeline.Run — community annotation is opt-in
		// and off by default, so the extra build cost is paid only when asked
		// for. Thread it through Verification instead if this flag sees
		// regular use and the rebuild cost starts to matter.
		communityMapPath := flags.CommunityMapPath
		if !filepath.IsAbs(communityMapPath) {
			communityMapPath = filepath.Join(repo, communityMapPath)
		}
		communityMap, err := fidelity.LoadCommunityMap(communityMapPath)
		if err != nil {
			return fmt.Errorf("--community-map: %w", err)
		}
		graph, err := adapter.Graph(ctx)
		if err != nil {
			return fmt.Errorf("--community-map: %w", err)
		}
		fidelity.AnnotateScopeCreepCommunities(&verification.Reconciliation, verification.Changes, graph, communityMap)
	}

	preflight := fidelity.CurrentPreflightReport()
	verdict := fidelity.BuildVerdict(verification, preflight, config, time.Now())

	outputDirectory := flags.OutputDirectory
	if !filepath.IsAbs(outputDirectory) {
		outputDirectory = filepath.Join(repo, outputDirectory)
	}
	verdictPath, err := fidelity.WriteVerdict(outputDirectory, verdict)
	if err != nil {
		return err
	}
	reportPath := filepath.Join(outputDirectory, "VERIFY_REPORT.md")
	if err := writeTextFile(reportPath, fidelity.RenderMarkdown(verdict)); err != nil {
		return err
	}
	// opts.Stdout is the JSON response contract; the human-readable terminal
	// report goes to stderr so scripted callers still get clean JSON on stdout.
	fmt.Fprint(opts.Stderr, fidelity.RenderTerminal(verdict))

	response := verifyIntentResponse{
		FormatVersion: 1, Status: "ok", CheckpointID: checkpointID,
		BaseCheckpointID: flags.BaseCheckpointID, ConfigPath: configPath,
		OutputDirectory: flags.OutputDirectory, Stages: runner.stages,
		Preflight: preflight, Verification: verification,
		Verdict: verdict, VerdictPath: verdictPath, ReportPath: reportPath,
		Message: "capture, declaration, semantic-diff observation, six-tier reconciliation (including coverage-confidence), and verdict/report rendering are complete",
	}
	encoder := json.NewEncoder(termsafe.NewJSONWriter(opts.Stdout))
	encoder.SetEscapeHTML(false)
	return encoder.Encode(response)
}

func writeTextFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
