package fidelity

import "testing"

func TestReconcileDirectClaimsSeparatesConfirmedAndMissing(t *testing.T) {
	intent := Intent{ExtractedEntities: []ExtractedEntity{
		{Name: "AuthRouter.handleTimeout", Method: "regex", Confidence: 1},
		{Name: "TokenValidator", Method: "regex", Confidence: 1},
	}}
	result := ReconcileDirectClaims(intent, []ChangedEntity{{Name: "handleTimeout", FilePath: "auth/router.go", StartLine: 88}})
	if len(result.Confirmed) != 1 || result.Confirmed[0].File != "auth/router.go" {
		t.Fatalf("confirmed = %#v", result.Confirmed)
	}
	if len(result.DeclaredUnimplemented) != 1 || result.DeclaredUnimplemented[0].Entity != "TokenValidator" {
		t.Fatalf("declared unimplemented = %#v", result.DeclaredUnimplemented)
	}
}
