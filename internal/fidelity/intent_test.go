package fidelity

import "testing"

func TestExtractDirectMentionsGroundsOnlyWholeSymbols(t *testing.T) {
	intent := ExtractDirectMentions(
		"Update AuthRouter.handleTimeout, then inspect AuthRouter.handleTimeout again. Do not touch VerifyIntent.",
		[]Symbol{{QualifiedName: "AuthRouter.handleTimeout"}, {QualifiedName: "Verify"}},
	)
	if len(intent.ExtractedEntities) != 1 {
		t.Fatalf("entities = %#v", intent.ExtractedEntities)
	}
	entity := intent.ExtractedEntities[0]
	if entity.Name != "AuthRouter.handleTimeout" || entity.Method != "regex" || entity.Confidence != 1 {
		t.Fatalf("entity = %#v", entity)
	}
}
