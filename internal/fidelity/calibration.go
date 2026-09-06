package fidelity

import (
	"fmt"
	"sort"
)

// calibration.go implements Bible §12.1-12.3: the Databricks calibration
// gate's PURE MATH — Conformal Risk Control threshold selection and online
// gating. Both run with zero network access and zero Databricks dependency,
// so they are built and tested now, independent of workspace credentials.
// The only Databricks-dependent pieces are (1) the I-CALM extraction call
// that PRODUCES a VerbalConfidence for a real edge, and (2) AI Search
// corroboration for whatever the gate can't resolve — both wired through a
// small client interface (see databricks.go) so this file never needs to
// change once credentials arrive.

// CalibrationEdge is one labeled example for I-CALM calibration (§12.2): a
// heuristic edge with a verbal confidence score and a known ground-truth
// outcome, used to compute a statistically-bounded acceptance threshold.
type CalibrationEdge struct {
	EdgeType           string
	VerbalConfidence   float64 // the model's own reported confidence in [0,1]
	CorrectlyPredicted bool    // ground truth: was this edge actually correct?
}

// ConformalThreshold computes λ̂ (#48): the smallest threshold such that
// accepting only edges with VerbalConfidence >= λ̂ keeps the empirical
// false-positive rate on the calibration set at or below alpha.
//
// This assumes the standard Conformal Risk Control monotonicity property:
// risk (false-positive rate among ACCEPTED edges) is non-increasing as the
// threshold rises. Under that assumption, walking thresholds from strictest
// (1.0, accept nothing, risk 0) down to more permissive ones and stopping at
// the first violation gives the smallest threshold that still holds the
// target risk. Real math, no external service — MAPIE (§20.2) would replace
// this hand-rolled selection with a citable library implementation without
// changing the config surface (gate.mode) or verdict.json schema.
func ConformalThreshold(calibration []CalibrationEdge, alpha float64) (float64, error) {
	if len(calibration) == 0 {
		return 0, fmt.Errorf("conformal threshold requires a non-empty calibration set")
	}
	if alpha <= 0 || alpha >= 1 {
		return 0, fmt.Errorf("alpha must be between 0 and 1")
	}
	candidates := make([]float64, 0, len(calibration)+1)
	seen := map[float64]bool{1.0: true}
	candidates = append(candidates, 1.0)
	for _, edge := range calibration {
		if !seen[edge.VerbalConfidence] {
			seen[edge.VerbalConfidence] = true
			candidates = append(candidates, edge.VerbalConfidence)
		}
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(candidates))) // strictest (1.0) first

	best := 1.0
	for _, threshold := range candidates {
		if conformalRisk(calibration, threshold) > alpha {
			break // this and every more-permissive threshold violates; keep the last good `best`
		}
		best = threshold
	}
	return best, nil
}

func conformalRisk(calibration []CalibrationEdge, threshold float64) float64 {
	accepted, falsePositives := 0, 0
	for _, edge := range calibration {
		if edge.VerbalConfidence >= threshold {
			accepted++
			if !edge.CorrectlyPredicted {
				falsePositives++
			}
		}
	}
	if accepted == 0 {
		return 0
	}
	return float64(falsePositives) / float64(accepted)
}

// ConformalRiskHolds verifies λ̂ actually holds the target false-positive
// rate on a held-out slice (the buildmap's own required test for #48).
func ConformalRiskHolds(heldOut []CalibrationEdge, threshold, alpha float64) bool {
	return conformalRisk(heldOut, threshold) <= alpha
}

// GateDecision is #50's online gating result for one advisory edge.
type GateDecision struct {
	ConformalGatePassed bool
	PendingVerification bool // epistemic abstention: below threshold, routed to corroboration, never silently trusted or dropped
}

// GateAdvisoryEdge implements §12.3's online gating. Deterministic edges
// never reach this function at all — classify.go's EdgeDeterministic bypass
// happens upstream, in reconciliation, not here. This only ever sees
// advisory edges.
func GateAdvisoryEdge(verbalConfidence, threshold float64) GateDecision {
	if verbalConfidence >= threshold {
		return GateDecision{ConformalGatePassed: true}
	}
	return GateDecision{ConformalGatePassed: false, PendingVerification: true}
}
