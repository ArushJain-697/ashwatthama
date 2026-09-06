package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/entireio/entire-graph/internal/fidelity"
	"github.com/entireio/entire-graph/internal/termsafe"
)

// verifyIntentFlags intentionally keeps the initial command surface small and
// stable. Transcript capture is not guessed here: the later Entire/Brain
// adapter will provide that capability behind internal/fidelity.Adapter.
type verifyIntentFlags struct {
	BaseCheckpointID string
	ConfigPath       string
	OutputDirectory  string
	Repo             string
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
	if err := fidelity.RunSkeleton(ctx, runner, request); err != nil {
		return err
	}
	response := verifyIntentResponse{
		FormatVersion: 1, Status: "skeleton", CheckpointID: checkpointID,
		BaseCheckpointID: flags.BaseCheckpointID, ConfigPath: configPath,
		OutputDirectory: flags.OutputDirectory, Stages: runner.stages,
		Preflight: fidelity.CurrentPreflightReport(),
		Message:   "command surface and pipeline order are ready; checkpoint transcript and graph adapters are intentionally not wired yet",
	}
	encoder := json.NewEncoder(termsafe.NewJSONWriter(opts.Stdout))
	encoder.SetEscapeHTML(false)
	return encoder.Encode(response)
}
