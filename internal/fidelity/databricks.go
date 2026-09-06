package fidelity

import (
	"context"
	"fmt"
)

// databricks.go is the seam §12's Databricks module calls through, mirroring
// adapter.go's own pattern: every real network/workspace call goes through
// one small interface, so wiring real credentials in later touches this file
// only — CalibrateAdvisoryEdges and the pure math in calibration.go never
// change.
type DatabricksClient interface {
	// ScoreEdge is §12.2's I-CALM extraction call: given an edge's hypothesis
	// in natural language and the constrained vocabulary it must ground to,
	// return a verbal confidence in [0,1]. Never invents a score outside that
	// range — same fail-closed posture as ResolveFallback (intent.go).
	ScoreEdge(ctx context.Context, hypothesis string, vocabulary []string) (float64, error)
	// Corroborate is §12.4's AI Search hybrid-retrieval call for a
	// pending_verification edge or an unverifiable_coverage entity (§12.15):
	// ask whether the repo's own docs/comments/commits confirm the
	// hypothesis. verified=false with no error means "no corroborating
	// evidence found," not a failure.
	Corroborate(ctx context.Context, hypothesis string) (verified bool, evidenceSnippet string, err error)
}

// UnavailableDatabricksClient is the safe default when calibration_enabled
// is false or no workspace is configured — every advisory/unverifiable
// entry's Databricks-reserved fields stay null, exactly as the schema
// promises (#26's "New fields ... null when it's off").
type UnavailableDatabricksClient struct{}

func (UnavailableDatabricksClient) ScoreEdge(context.Context, string, []string) (float64, error) {
	return 0, fmt.Errorf("Databricks client is not configured")
}

func (UnavailableDatabricksClient) Corroborate(context.Context, string) (bool, string, error) {
	return false, "", fmt.Errorf("Databricks client is not configured")
}

// CalibrateAdvisoryEdges implements §12.3/§12.5: for every advisory_low_confidence
// entry, score it via I-CALM, gate it against the conformal threshold, and — for
// anything the gate abstains on — attempt corroboration. Never touches
// unverifiable_coverage's baseline verification_path (#35's guarantee holds
// independent of this function ever running at all).
func CalibrateAdvisoryEdges(ctx context.Context, reconciliation *Reconciliation, client DatabricksClient, threshold float64) {
	for i := range reconciliation.AdvisoryLowConfidence {
		entry := &reconciliation.AdvisoryLowConfidence[i]
		confidence, err := client.ScoreEdge(ctx, entry.Reason, nil)
		if err != nil {
			continue // fail closed: leave verbal_confidence/conformal_gate_passed null, exactly as an unconfigured client would
		}
		entry.VerbalConfidence = &confidence
		decision := GateAdvisoryEdge(confidence, threshold)
		entry.ConformalGatePassed = &decision.ConformalGatePassed
		if !decision.PendingVerification {
			continue
		}
		entry.Corroboration.Attempted = true
		verified, snippet, err := client.Corroborate(ctx, entry.Reason)
		if err != nil {
			continue // corroboration attempted, verified stays nil — a failed lookup is not a "no"
		}
		entry.Corroboration.Verified = &verified
		if verified {
			entry.Corroboration.EvidenceSnippet = &snippet
		}
	}
}

// CorroborateUnverifiableCoverage implements §12.15: the Databricks-enhanced
// tier for the mandatory Curveball baseline (#35), never a substitute for
// it. A miss leaves VerificationPath unchanged — the baseline still holds.
func CorroborateUnverifiableCoverage(ctx context.Context, reconciliation *Reconciliation, client DatabricksClient) {
	for i := range reconciliation.UnverifiableCoverage {
		entry := &reconciliation.UnverifiableCoverage[i]
		entry.Corroboration.Attempted = true
		verified, snippet, err := client.Corroborate(ctx, entry.CoverageReason)
		if err != nil {
			continue
		}
		entry.Corroboration.Verified = &verified
		if verified {
			entry.Corroboration.EvidenceSnippet = &snippet
		}
	}
}
