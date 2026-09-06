package fidelity

import (
	"fmt"
	"strconv"
	"strings"
)

// render.go implements #27 (terminal), retrofitted for the Curveball's
// mandatory three-way split (#38), plus #39's Markdown renderer built
// three-way-aware from the start. Both are pure functions of Verdict — no
// adapter or repo access — per the §13.1 firewall.
//
// The three visually distinct states the card requires, verbatim:
//   - confirmed structural evidence      (confirmed, expected_blast_radius)
//   - heuristic/incomplete evidence      (advisory_low_confidence)
//   - claims needing source/test verification (unverifiable_coverage)
// undeclared_scope_creep and declared_unimplemented are drift/gap signals,
// not one of these three states, and get their own (fourth) treatment.

const (
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiMagenta = "\033[35m"
	ansiRed     = "\033[31m"
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
)

// Visual category labels, shared verbatim between the terminal and Markdown
// renderers so the distinctness test (#41) can assert both say the same thing.
const (
	labelConfirmed    = "CONFIRMED STRUCTURAL EVIDENCE"
	labelHeuristic    = "HEURISTIC / INCOMPLETE EVIDENCE"
	labelUnverifiable = "NEEDS SOURCE/TEST VERIFICATION"
	labelDrift        = "DRIFT FROM STATED INTENT"
)

// RenderTerminal is a colored, tiered, file:line-annotated report.
func RenderTerminal(verdict Verdict) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%sFidelity verdict%s — checkpoint %s (%s)\n", ansiBold, ansiReset, verdict.CheckpointID, verdict.Summary.VerdictLabel)
	fmt.Fprintf(&b, "%d changed entities: %d confirmed, %d expected blast radius, %d scope creep, %d advisory, %d unverifiable coverage, %d declared unimplemented\n\n",
		verdict.Summary.TotalChangedEntities, len(verdict.Reconciliation.Confirmed), len(verdict.Reconciliation.ExpectedBlastRadius),
		verdict.Summary.ScopeCreepCount, verdict.Summary.AdvisoryEdgeCount, verdict.Summary.UnverifiableCoverageCount, verdict.Summary.UnimplementedCount)

	terminalSection(&b, ansiGreen, "✓", labelConfirmed, "confirmed", len(verdict.Reconciliation.Confirmed), func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.Confirmed {
			fmt.Fprintf(w, "  ✓ %s (%s) [coverage: %s]\n", entry.Entity, location(entry.File, entry.Line), entry.CoverageConfidence)
		}
	})
	terminalSection(&b, ansiGreen, "✓", labelConfirmed, "expected_blast_radius", len(verdict.Reconciliation.ExpectedBlastRadius), func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.ExpectedBlastRadius {
			fmt.Fprintf(w, "  ✓ %s (%s) — %s hop(s) via %s [coverage: %s]\n",
				entry.Entity, location(entry.File, entry.Line), strconv.Itoa(entry.Hops), strings.Join(entry.EdgeTypesTraversed, ","), entry.CoverageConfidence)
		}
	})
	terminalSection(&b, ansiYellow, "~", labelHeuristic, "advisory_low_confidence", len(verdict.Reconciliation.AdvisoryLowConfidence), func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.AdvisoryLowConfidence {
			fmt.Fprintf(w, "  ~ %s (%s) — %s [coverage: %s]\n", entry.Entity, location(entry.File, entry.Line), entry.Reason, entry.CoverageConfidence)
		}
	})
	terminalSection(&b, ansiMagenta, "?", labelUnverifiable, "unverifiable_coverage", len(verdict.Reconciliation.UnverifiableCoverage), func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.UnverifiableCoverage {
			fmt.Fprintf(w, "  ? %s (%s) — %s [%s]\n", entry.Entity, location(entry.File, entry.Line), entry.CoverageReason, entry.VerificationPath)
		}
	})
	terminalSection(&b, ansiRed, "✗", labelDrift, "undeclared_scope_creep", len(verdict.Reconciliation.UndeclaredScopeCreep), func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.UndeclaredScopeCreep {
			fmt.Fprintf(w, "  ✗ %s (%s) — %s [coverage: %s]\n", entry.Entity, location(entry.File, entry.Line), entry.Reason, entry.CoverageConfidence)
		}
	})
	terminalSection(&b, ansiRed, "✗", labelDrift, "declared_unimplemented", len(verdict.Reconciliation.DeclaredUnimplemented), func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.DeclaredUnimplemented {
			fmt.Fprintf(w, "  ✗ %s — %s\n", entry.Entity, entry.Reason)
		}
	})
	return b.String()
}

func terminalSection(b *strings.Builder, color, marker, visualLabel, tierName string, count int, body func(*strings.Builder)) {
	if count == 0 {
		return
	}
	fmt.Fprintf(b, "%s%s %s%s — %s (%d)\n", color, marker, visualLabel, ansiReset, tierName, count)
	body(b)
	b.WriteString("\n")
}

func location(file string, line int) string {
	if file == "" {
		return "no location"
	}
	if line == 0 {
		return file
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// RenderMarkdown produces VERIFY_REPORT.md, three-way-aware from the start —
// this renderer did not exist before the Curveball, so there is no retrofit.
func RenderMarkdown(verdict Verdict) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Fidelity Verdict Report\n\n")
	fmt.Fprintf(&b, "Checkpoint `%s`", verdict.CheckpointID)
	if verdict.BaseCheckpointID != "" {
		fmt.Fprintf(&b, " (base `%s`)", verdict.BaseCheckpointID)
	}
	fmt.Fprintf(&b, " — generated %s — **%s**\n\n", verdict.GeneratedAt, verdict.Summary.VerdictLabel)

	fmt.Fprintf(&b, "| Tier | Visual category | Count |\n|---|---|---|\n")
	fmt.Fprintf(&b, "| confirmed | %s | %d |\n", labelConfirmed, len(verdict.Reconciliation.Confirmed))
	fmt.Fprintf(&b, "| expected_blast_radius | %s | %d |\n", labelConfirmed, len(verdict.Reconciliation.ExpectedBlastRadius))
	fmt.Fprintf(&b, "| advisory_low_confidence | %s | %d |\n", labelHeuristic, len(verdict.Reconciliation.AdvisoryLowConfidence))
	fmt.Fprintf(&b, "| unverifiable_coverage | %s | %d |\n", labelUnverifiable, len(verdict.Reconciliation.UnverifiableCoverage))
	fmt.Fprintf(&b, "| undeclared_scope_creep | %s | %d |\n", labelDrift, len(verdict.Reconciliation.UndeclaredScopeCreep))
	fmt.Fprintf(&b, "| declared_unimplemented | %s | %d |\n\n", labelDrift, len(verdict.Reconciliation.DeclaredUnimplemented))

	markdownSection(&b, "✅ "+labelConfirmed, "Confirmed", func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.Confirmed {
			fmt.Fprintf(w, "- `%s` at %s — coverage: `%s`\n", entry.Entity, location(entry.File, entry.Line), entry.CoverageConfidence)
		}
		for _, entry := range verdict.Reconciliation.ExpectedBlastRadius {
			fmt.Fprintf(w, "- `%s` at %s — %d hop(s) via %s — coverage: `%s`\n",
				entry.Entity, location(entry.File, entry.Line), entry.Hops, strings.Join(entry.EdgeTypesTraversed, ", "), entry.CoverageConfidence)
		}
	})
	markdownSection(&b, "🟡 "+labelHeuristic, "advisory_low_confidence", func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.AdvisoryLowConfidence {
			fmt.Fprintf(w, "- `%s` at %s — %s — coverage: `%s`\n", entry.Entity, location(entry.File, entry.Line), entry.Reason, entry.CoverageConfidence)
		}
	})
	markdownSection(&b, "❓ "+labelUnverifiable, "unverifiable_coverage", func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.UnverifiableCoverage {
			fmt.Fprintf(w, "- `%s` at %s — %s — `%s`\n", entry.Entity, location(entry.File, entry.Line), entry.CoverageReason, entry.VerificationPath)
		}
	})
	markdownSection(&b, "🔴 "+labelDrift, "undeclared_scope_creep / declared_unimplemented", func(w *strings.Builder) {
		for _, entry := range verdict.Reconciliation.UndeclaredScopeCreep {
			fmt.Fprintf(w, "- `%s` at %s — %s — coverage: `%s`\n", entry.Entity, location(entry.File, entry.Line), entry.Reason, entry.CoverageConfidence)
		}
		for _, entry := range verdict.Reconciliation.DeclaredUnimplemented {
			fmt.Fprintf(w, "- `%s` — %s\n", entry.Entity, entry.Reason)
		}
	})
	return b.String()
}

func markdownSection(b *strings.Builder, heading, tierNames string, body func(*strings.Builder)) {
	fmt.Fprintf(b, "## %s\n\n", heading)
	var section strings.Builder
	body(&section)
	if section.Len() == 0 {
		fmt.Fprintf(b, "_none in %s_\n\n", tierNames)
		return
	}
	b.WriteString(section.String())
	b.WriteString("\n")
}
