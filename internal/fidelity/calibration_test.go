package fidelity

import "testing"

// calibrationFixture is a synthetic set where higher verbal confidence
// genuinely correlates with correctness (the monotonicity CRC assumes),
// deliberately labeled synthetic per BUILDATHON.md's honesty requirement —
// not real production revert data.
func calibrationFixture() []CalibrationEdge {
	var edges []CalibrationEdge
	// High confidence, mostly correct (9/10 correct at 0.9+).
	for i := 0; i < 9; i++ {
		edges = append(edges, CalibrationEdge{VerbalConfidence: 0.95, CorrectlyPredicted: true})
	}
	edges = append(edges, CalibrationEdge{VerbalConfidence: 0.95, CorrectlyPredicted: false})
	// Mid confidence, roughly half correct.
	for i := 0; i < 5; i++ {
		edges = append(edges, CalibrationEdge{VerbalConfidence: 0.6, CorrectlyPredicted: i%2 == 0})
	}
	// Low confidence, mostly wrong.
	for i := 0; i < 8; i++ {
		edges = append(edges, CalibrationEdge{VerbalConfidence: 0.2, CorrectlyPredicted: i == 0})
	}
	return edges
}

func TestConformalThresholdHoldsTargetRiskOnCalibrationSet(t *testing.T) {
	calibration := calibrationFixture()
	threshold, err := ConformalThreshold(calibration, 0.05)
	if err != nil {
		t.Fatal(err)
	}
	if !ConformalRiskHolds(calibration, threshold, 0.05) {
		t.Fatalf("threshold %v does not hold alpha=0.05 on its own calibration set", threshold)
	}
	// The whole point of CRC: threshold must reject the low-confidence,
	// mostly-wrong tier and accept the high-confidence, mostly-right tier.
	if threshold < 0.9 {
		t.Fatalf("threshold = %v, want >= 0.9 (only the high-confidence tier meets 5%% risk here)", threshold)
	}
}

func TestConformalThresholdHoldsOnHeldOutSlice(t *testing.T) {
	// #48's required test: verify the computed threshold holds on data it
	// was NOT computed from, not just the set that produced it.
	trainSet := calibrationFixture()
	threshold, err := ConformalThreshold(trainSet, 0.05)
	if err != nil {
		t.Fatal(err)
	}
	heldOut := []CalibrationEdge{
		{VerbalConfidence: 0.97, CorrectlyPredicted: true},
		{VerbalConfidence: 0.93, CorrectlyPredicted: true},
		{VerbalConfidence: 0.91, CorrectlyPredicted: true},
		{VerbalConfidence: 0.6, CorrectlyPredicted: false}, // below threshold, must not count against it
	}
	if !ConformalRiskHolds(heldOut, threshold, 0.05) {
		t.Fatalf("threshold %v did not hold alpha=0.05 on a held-out slice", threshold)
	}
}

func TestConformalThresholdRejectsEmptyCalibrationSet(t *testing.T) {
	if _, err := ConformalThreshold(nil, 0.05); err == nil {
		t.Fatal("ConformalThreshold accepted an empty calibration set")
	}
}

func TestGateAdvisoryEdgeAbstainsRatherThanDropsBelowThreshold(t *testing.T) {
	pass := GateAdvisoryEdge(0.95, 0.9)
	if !pass.ConformalGatePassed || pass.PendingVerification {
		t.Fatalf("high-confidence edge = %#v, want gate passed", pass)
	}
	abstain := GateAdvisoryEdge(0.5, 0.9)
	if abstain.ConformalGatePassed || !abstain.PendingVerification {
		t.Fatalf("low-confidence edge = %#v, want epistemic abstention (pending_verification), not a silent drop", abstain)
	}
}
