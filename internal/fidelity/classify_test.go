package fidelity

import "testing"

func TestClassifyRelationUsesRealProviderFields(t *testing.T) {
	config := DefaultConfig()
	if got := ClassifyRelation(GraphRelation{Confidence: 1, Resolution: "exact"}, config); got != EdgeDeterministic {
		t.Fatalf("exact relation = %q", got)
	}
	if got := ClassifyRelation(GraphRelation{Confidence: 1, Resolution: "pattern"}, config); got != EdgeAdvisory {
		t.Fatalf("pattern relation = %q", got)
	}
}
