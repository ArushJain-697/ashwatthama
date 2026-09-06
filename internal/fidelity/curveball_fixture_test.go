package fidelity

import (
	"os"
	"path/filepath"
	"testing"
)

// curveball_fixture_test.go is Build Map #42: the mandatory partial-analysis
// fixture test for Track 2, "Graph Is Evidence, Not an Oracle."
//
// HONEST DISCLOSURE (see BUILDATHON.md): no organizer-supplied fixture
// repository was found on this machine or attached to the received card. The
// card explicitly asks for a run against "the actual supplied fixture, not a
// hand-built substitute" — that requirement cannot be met without the file.
// This test is a substitute, built to the same shape: a real, disposable Git
// repository, analyzed by the REAL NativeGraphAdapter and the real
// entire-graph engine (no mocked Adapter, no hand-built GraphSnapshot), with
// one region reflection cannot resolve and one fully-resolved region in the
// same commit. If the real fixture arrives, this test's fixture repo should
// be replaced with it; the four assertions below should not need to change
// shape to do that.
func writeCurveballFixtureRepo(t *testing.T) (repo, base, head string) {
	t.Helper()
	repo = t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "fidelity@example.test")
	runGit(t, repo, "config", "user.name", "Fidelity Test")

	// dispatch.go: handlePing is reached only through a runtime function-value
	// registry keyed by string, the classic command/plugin-registry dynamic
	// dispatch pattern. Verified empirically against the real engine (not
	// assumed): entire-graph attaches no CALLS/DATA_FLOWS/any reference edge
	// to handlePing here — only the structural DEFINES edge every symbol
	// gets — so it is a genuine zero-edge case, not a hand-picked one.
	dispatch := `package dispatch

var handlers map[string]func() string

func init() {
	handlers = map[string]func() string{"ping": handlePing}
}

func handlePing() string { return "pong" }

func Dispatch(name string) string { return handlers[name]() }
`
	// plain.go: Caller calls Callee directly — a fully-resolved region that
	// must keep reconciling exactly as it did before this Curveball response.
	plain := `package plain

func Callee() int { return 1 }

func Caller() int { return Callee() }
`
	mustWrite(t, filepath.Join(repo, "dispatch.go"), dispatch)
	mustWrite(t, filepath.Join(repo, "plain.go"), plain)
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-m", "base")
	base = trimmed(t, repo, "rev-parse", "HEAD")

	dispatchChanged := `package dispatch

var handlers map[string]func() string

func init() {
	handlers = map[string]func() string{"ping": handlePing}
}

func handlePing() string { return "PONG" }

func Dispatch(name string) string { return handlers[name]() }
`
	plainChanged := `package plain

func Callee() int { return 2 }

func Caller() int { return Callee() }
`
	mustWrite(t, filepath.Join(repo, "dispatch.go"), dispatchChanged)
	mustWrite(t, filepath.Join(repo, "plain.go"), plainChanged)
	runGit(t, repo, "commit", "-am", "change handlePing and Callee")
	head = trimmed(t, repo, "rev-parse", "HEAD")
	return repo, base, head
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func trimmed(t *testing.T, repo string, args ...string) string {
	t.Helper()
	out := runGit(t, repo, args...)
	for len(out) > 0 && (out[len(out)-1] == '\n' || out[len(out)-1] == '\r') {
		out = out[:len(out)-1]
	}
	return out
}

// TestPartialAnalysisFixtureNeverMisfilesReflectionDispatchAsScopeCreep is
// #42's mandatory test, run against the substitute fixture above (see the
// honest disclosure at the top of this file). It asserts all four required
// properties.
func TestPartialAnalysisFixtureNeverMisfilesReflectionDispatchAsScopeCreep(t *testing.T) {
	repo, base, head := writeCurveballFixtureRepo(t)
	adapter := NativeGraphAdapter{Repo: repo, ProviderVersion: "test"}

	graph, err := adapter.Graph(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	changes, err := adapter.ChangedEntities(t.Context(), base, head)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("fixture produced no changed entities; nothing to assert against")
	}

	result := Reconcile(Intent{}, changes, graph, DefaultConfig())

	// (1) No entity in the affected (reflection-dispatched) region is
	// misfiled into undeclared_scope_creep.
	for _, entry := range result.UndeclaredScopeCreep {
		if entry.Entity == "handlePing" {
			t.Fatalf("handlePing (reflection-dispatched, unresolvable) was misfiled as confident scope creep: %#v", entry)
		}
	}

	// (2) Every affected entity carries coverage_confidence partial (or
	// unknown) and lands in unverifiable_coverage.
	var handlePingEntry *UnverifiableCoverageEntity
	for i, entry := range result.UnverifiableCoverage {
		if entry.Entity == "handlePing" {
			handlePingEntry = &result.UnverifiableCoverage[i]
		}
	}
	if handlePingEntry == nil {
		t.Fatalf("handlePing did not land in unverifiable_coverage: %#v", result.UnverifiableCoverage)
	}
	if handlePingEntry.CoverageConfidence != CoveragePartial && handlePingEntry.CoverageConfidence != CoverageUnknown {
		t.Fatalf("handlePing coverage_confidence = %q, want partial or unknown", handlePingEntry.CoverageConfidence)
	}

	// (3) verification_path is populated, at minimum with
	// manual_review_recommended — the baseline safety guarantee (#35), which
	// this test runs with Databricks fully disabled (DefaultConfig's zero
	// value: calibration_enabled/ai_search/coverage_corroboration all false).
	if handlePingEntry.VerificationPath == "" {
		t.Fatal("handlePing's verification_path was not populated")
	}
	if handlePingEntry.VerificationPath != "manual_review_recommended" {
		t.Fatalf("verification_path = %q, want the manual_review_recommended baseline", handlePingEntry.VerificationPath)
	}

	// (4) A fully-resolved entity elsewhere in the SAME fixture (Callee, a
	// direct static call from Caller) still reconciles normally: it must not
	// appear in unverifiable_coverage, and its coverage must be full.
	for _, entry := range result.UnverifiableCoverage {
		if entry.Entity == "Callee" {
			t.Fatalf("Callee is fully resolved and must not land in unverifiable_coverage: %#v", entry)
		}
	}
	found := false
	for _, entry := range result.UndeclaredScopeCreep {
		if entry.Entity == "Callee" {
			found = true
			if entry.CoverageConfidence != CoverageFull {
				t.Fatalf("Callee coverage_confidence = %q, want full", entry.CoverageConfidence)
			}
		}
	}
	if !found {
		t.Fatalf("Callee (fully resolved, unreachable from any confirmed entity here) was not found in undeclared_scope_creep as expected: %#v", result)
	}
}
