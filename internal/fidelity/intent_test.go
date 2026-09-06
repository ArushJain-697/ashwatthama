package fidelity

import "testing"

type staticFallback struct{ name string }

func (fallback staticFallback) Resolve(string, []string) (string, float64, error) {
	return fallback.name, 0.7, nil
}

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

func TestExtractDirectMentionsAcceptsPeriodAfterQualifiedSymbol(t *testing.T) {
	intent := ExtractDirectMentions("Implement service.VerifyIntent.", []Symbol{{QualifiedName: "service.VerifyIntent"}})
	if len(intent.ExtractedEntities) != 1 || intent.ExtractedEntities[0].Name != "service.VerifyIntent" {
		t.Fatalf("entities = %#v", intent.ExtractedEntities)
	}
}

func TestDeclareIntentInheritsOnlyExplicitVagueFollowUp(t *testing.T) {
	intent := DeclareIntent("continue", "Implement service.VerifyIntent.", []Symbol{{QualifiedName: "service.VerifyIntent"}}, false)
	if !intent.InheritedIntent || intent.NoBaseline || len(intent.ExtractedEntities) != 1 || intent.ExtractedEntities[0].Name != "service.VerifyIntent" {
		t.Fatalf("intent = %#v", intent)
	}

	notInherited := DeclareIntent("please continue", "Implement service.VerifyIntent.", []Symbol{{QualifiedName: "service.VerifyIntent"}}, false)
	if notInherited.InheritedIntent || len(notInherited.ExtractedEntities) != 0 {
		t.Fatalf("non-vague intent = %#v", notInherited)
	}
}

func TestDeclareIntentMarksNoBaseline(t *testing.T) {
	intent := DeclareIntent("Implement service.VerifyIntent", "", []Symbol{{QualifiedName: "service.VerifyIntent"}}, true)
	if !intent.NoBaseline || intent.InheritedIntent {
		t.Fatalf("intent = %#v", intent)
	}
}

func TestResolveFallbackRejectsInventedEntity(t *testing.T) {
	entity, unresolved, err := ResolveFallback("login flow", []string{"AuthRouter"}, staticFallback{name: "ImaginaryRouter"})
	if err != nil {
		t.Fatal(err)
	}
	if entity.Name != "" || unresolved != "login flow" {
		t.Fatalf("fallback invented entity: %#v, %q", entity, unresolved)
	}
}
