package fidelity

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runGit(t *testing.T, repo string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = repo
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
	return string(output)
}

func TestNativeGraphAdapterLoadsSymbolsAndChanges(t *testing.T) {
	repo := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		command := exec.Command("git", args...)
		command.Dir = repo
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, output)
		}
	}
	git("init")
	git("config", "user.email", "fidelity@example.test")
	git("config", "user.name", "Fidelity Test")
	path := filepath.Join(repo, "service.go")
	if err := os.WriteFile(path, []byte("package service\n\nfunc Verify() bool { return true }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	git("add", "service.go")
	git("commit", "-m", "base")
	git("rev-parse", "HEAD")
	if err := os.WriteFile(path, []byte("package service\n\nfunc Verify() bool { return false }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	git("commit", "-am", "change")

	adapter := NativeGraphAdapter{Repo: repo, ProviderVersion: "test"}
	graph, err := adapter.Graph(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(graph.Symbols) == 0 || graph.Symbols[0].QualifiedName == "" {
		t.Fatalf("symbols = %#v", graph.Symbols)
	}
	if graph.CompletenessLevel == "" {
		t.Fatalf("Graph did not report a completeness level: %#v", graph)
	}
	changes, err := adapter.ChangedEntities(t.Context(), "HEAD~1", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 || changes[0].Name != "Verify" {
		t.Fatalf("changes = %#v", changes)
	}
}

func TestNativeGraphAdapterDiffsAgainstEmptyTree(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "fidelity@example.test")
	runGit(t, repo, "config", "user.name", "Fidelity Test")
	if err := os.WriteFile(filepath.Join(repo, "service.go"), []byte("package service\n\nfunc First() {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "service.go")
	runGit(t, repo, "commit", "-m", "first")
	changes, err := (NativeGraphAdapter{Repo: repo, ProviderVersion: "test"}).ChangedEntities(t.Context(), emptyTreeObjectID, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 || changes[0].Name != "First" {
		t.Fatalf("changes = %#v", changes)
	}
}

func TestNativeGraphAdapterReadsCheckpointTranscriptWithEntireCLI(t *testing.T) {
	bin := t.TempDir()
	command := filepath.Join(bin, "entire")
	if err := os.WriteFile(command, []byte("#!/bin/sh\nprintf '%s\\n' '{\"type\":\"message\",\"text\":\"fix timeout\"}'\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	transcript, err := (NativeGraphAdapter{Repo: t.TempDir()}).CheckpointTranscript(t.Context(), "checkpoint-1")
	if err != nil {
		t.Fatal(err)
	}
	if transcript != "{\"type\":\"message\",\"text\":\"fix timeout\"}\n" {
		t.Fatalf("transcript = %q", transcript)
	}
}

func TestNativeGraphAdapterResolvesCheckpointCommit(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.com")
	runGit(t, repo, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(repo, "file.txt"), []byte("content\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "file.txt")
	runGit(t, repo, "commit", "-m", "checkpoint\n\nEntire-Checkpoint: checkpoint-1")
	want := strings.TrimSpace(runGit(t, repo, "rev-parse", "HEAD"))
	got, err := (NativeGraphAdapter{Repo: repo}).CheckpointCommit(t.Context(), "checkpoint-1")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("CheckpointCommit = %q, want %q", got, want)
	}
}
