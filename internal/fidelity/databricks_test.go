package fidelity

import (
	"context"
	"testing"
)

type fakeDatabricksClient struct {
	confidence float64
	verified   bool
	snippet    string
}

func (f fakeDatabricksClient) ScoreEdge(context.Context, string, []string) (float64, error) {
	return f.confidence, nil
}

func (f fakeDatabricksClient) Corroborate(context.Context, string) (bool, string, error) {
	return f.verified, f.snippet, nil
}

func TestCalibrateAdvisoryEdgesPassesHighConfidenceGate(t *testing.T) {
	reconciliation := &Reconciliation{AdvisoryLowConfidence: []AdvisoryEntity{{Entity: "X", Reason: "maybe calls Y"}}}
	CalibrateAdvisoryEdges(t.Context(), reconciliation, fakeDatabricksClient{confidence: 0.95}, 0.9)
	entry := reconciliation.AdvisoryLowConfidence[0]
	if entry.VerbalConfidence == nil || *entry.VerbalConfidence != 0.95 {
		t.Fatalf("verbal_confidence = %#v", entry.VerbalConfidence)
	}
	if entry.ConformalGatePassed == nil || !*entry.ConformalGatePassed {
		t.Fatalf("conformal_gate_passed = %#v, want true", entry.ConformalGatePassed)
	}
	if entry.Corroboration.Attempted {
		t.Fatal("a passing gate must not attempt corroboration")
	}
}

func TestCalibrateAdvisoryEdgesCorroboratesBelowThreshold(t *testing.T) {
	reconciliation := &Reconciliation{AdvisoryLowConfidence: []AdvisoryEntity{{Entity: "X", Reason: "maybe calls Y"}}}
	CalibrateAdvisoryEdges(t.Context(), reconciliation, fakeDatabricksClient{confidence: 0.3, verified: true, snippet: "docs say X calls Y"}, 0.9)
	entry := reconciliation.AdvisoryLowConfidence[0]
	if entry.ConformalGatePassed == nil || *entry.ConformalGatePassed {
		t.Fatalf("conformal_gate_passed = %#v, want false", entry.ConformalGatePassed)
	}
	if !entry.Corroboration.Attempted || entry.Corroboration.Verified == nil || !*entry.Corroboration.Verified {
		t.Fatalf("corroboration = %#v, want attempted+verified", entry.Corroboration)
	}
	if entry.Corroboration.EvidenceSnippet == nil || *entry.Corroboration.EvidenceSnippet != "docs say X calls Y" {
		t.Fatalf("evidence_snippet = %#v", entry.Corroboration.EvidenceSnippet)
	}
}

func TestUnavailableDatabricksClientFailsClosedWithoutCrashing(t *testing.T) {
	reconciliation := &Reconciliation{AdvisoryLowConfidence: []AdvisoryEntity{{Entity: "X", Reason: "maybe calls Y"}}}
	CalibrateAdvisoryEdges(t.Context(), reconciliation, UnavailableDatabricksClient{}, 0.9)
	entry := reconciliation.AdvisoryLowConfidence[0]
	if entry.VerbalConfidence != nil || entry.ConformalGatePassed != nil {
		t.Fatalf("an unconfigured client must leave fields null, got %#v", entry)
	}
}

func TestCorroborateUnverifiableCoverageLeavesBaselineIntactOnMiss(t *testing.T) {
	reconciliation := &Reconciliation{UnverifiableCoverage: []UnverifiableCoverageEntity{
		{Entity: "X", VerificationPath: "manual_review_recommended"},
	}}
	CorroborateUnverifiableCoverage(t.Context(), reconciliation, fakeDatabricksClient{verified: false})
	entry := reconciliation.UnverifiableCoverage[0]
	if entry.VerificationPath != "manual_review_recommended" {
		t.Fatalf("a corroboration miss must never change the baseline verification_path: %#v", entry)
	}
	if !entry.Corroboration.Attempted || entry.Corroboration.Verified == nil || *entry.Corroboration.Verified {
		t.Fatalf("corroboration = %#v, want attempted+not-verified", entry.Corroboration)
	}
}
