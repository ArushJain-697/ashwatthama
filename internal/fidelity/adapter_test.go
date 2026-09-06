package fidelity

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

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
	symbols, err := adapter.SymbolDictionary(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(symbols) == 0 || symbols[0].QualifiedName == "" {
		t.Fatalf("symbols = %#v", symbols)
	}
	changes, err := adapter.ChangedEntities(t.Context(), "HEAD~1", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 || changes[0].Name != "Verify" {
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
