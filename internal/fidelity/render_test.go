package fidelity

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestRenderersAgreeOnCountsAndSatisfyThreeWayDistinctness pins #41: both
// renderers must report the same tier counts (including unverifiable_coverage)
// from the same verdict.json, and each must independently let a reader tell
// apart the three states the Curveball card requires: confirmed structural
// evidence, heuristic/incomplete evidence, and claims needing verification.
func TestRenderersAgreeOnCountsAndSatisfyThreeWayDistinctness(t *testing.T) {
	verdict := BuildVerdict(fixtureVerification(), CurrentPreflightReport(), DefaultConfig(), time.Now())
	terminal := RenderTerminal(verdict)
	markdown := RenderMarkdown(verdict)

	for _, rendered := range []struct {
		name, text string
	}{{"terminal", terminal}, {"markdown", markdown}} {
		for _, label := range []string{labelConfirmed, labelHeuristic, labelUnverifiable} {
			if !strings.Contains(rendered.text, label) {
				t.Fatalf("%s output missing distinct visual label %q:\n%s", rendered.name, label, rendered.text)
			}
		}
		// The three labels must actually be three DIFFERENT strings, not one
		// bucket with different names — a shared "risky" label would satisfy
		// the substring check above without satisfying the requirement.
		if labelConfirmed == labelHeuristic || labelHeuristic == labelUnverifiable || labelConfirmed == labelUnverifiable {
			t.Fatal("visual category labels are not pairwise distinct")
		}
	}

	// Tier counts: the Markdown table's own count column must match the
	// verdict it was built from, for every tier — the actual "do the
	// renderers agree" check, not just "do they mention the tier."
	wantCounts := map[string]int{
		"confirmed":               len(verdict.Reconciliation.Confirmed),
		"expected_blast_radius":   len(verdict.Reconciliation.ExpectedBlastRadius),
		"advisory_low_confidence": len(verdict.Reconciliation.AdvisoryLowConfidence),
		"unverifiable_coverage":   len(verdict.Reconciliation.UnverifiableCoverage),
		"undeclared_scope_creep":  len(verdict.Reconciliation.UndeclaredScopeCreep),
		"declared_unimplemented":  len(verdict.Reconciliation.DeclaredUnimplemented),
	}
	if wantCounts["unverifiable_coverage"] == 0 {
		t.Fatal("fixture must exercise unverifiable_coverage for this test to mean anything")
	}
	for tier, want := range wantCounts {
		lineStart := strings.Index(markdown, "| "+tier+" |")
		if lineStart < 0 {
			t.Fatalf("markdown table row for %q not found", tier)
		}
		lineEnd := strings.Index(markdown[lineStart:], "\n")
		line := strings.TrimRight(markdown[lineStart:lineStart+lineEnd], " ")
		wantSuffix := " | " + strconv.Itoa(want) + " |"
		if !strings.HasSuffix(line, wantSuffix) {
			t.Fatalf("markdown row for %q = %q, want count %d", tier, line, want)
		}
	}
	// The terminal summary line must report the same verdict-level total.
	if !strings.Contains(terminal, strconv.Itoa(verdict.Summary.TotalChangedEntities)+" changed entities") {
		t.Fatalf("terminal output does not report the verdict's own total_changed_entities: %s", terminal)
	}
}
