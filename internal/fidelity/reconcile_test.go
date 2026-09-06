package fidelity

import "testing"

// fixtureGraph is a small synthetic snapshot exercising all six tiers:
//   - AuthRouter.handleTimeout (A): named in intent, matched by diff -> confirmed
//   - SessionCache.invalidate (B): 1-hop deterministic callee of A -> expected_blast_radius
//   - BillingController.chargeRetry (C): unreachable from A, but has its own
//     deterministic edge elsewhere (to F) so coverage is full -> undeclared_scope_creep
//   - PaymentGateway.route (D): reachable from A only via a low-confidence
//     pattern-resolved edge -> advisory_low_confidence, coverage full (the
//     edge exists, it's just uncertain — orthogonal to coverage)
//   - PluginRegistry.dispatch (E): zero graph edges despite a known file ->
//     unverifiable_coverage, coverage partial
//   - Ghost (G): its file is entirely absent from the snapshot ->
//     unverifiable_coverage, coverage unknown
func fixtureGraph() GraphSnapshot {
	return GraphSnapshot{
		Symbols: []Symbol{
			{ID: "A", QualifiedName: "AuthRouter.handleTimeout", FilePath: "auth/router.go"},
			{ID: "B", QualifiedName: "SessionCache.invalidate", FilePath: "auth/cache.go"},
			{ID: "C", QualifiedName: "BillingController.chargeRetry", FilePath: "billing/controller.go"},
			{ID: "D", QualifiedName: "PaymentGateway.route", FilePath: "gateway/route.go"},
			{ID: "E", QualifiedName: "PluginRegistry.dispatch", FilePath: "plugins/registry.go"},
			{ID: "F", QualifiedName: "Ledger.append", FilePath: "billing/ledger.go"},
		},
		Relations: []GraphRelation{
			{FromID: "A", ToID: "B", Type: "CALLS", Confidence: 1.0, Resolution: "exact"},
			{FromID: "A", ToID: "D", Type: "CALLS", Confidence: 0.6, Resolution: "pattern"},
			{FromID: "C", ToID: "F", Type: "CALLS", Confidence: 1.0, Resolution: "exact"},
		},
		FileLanguages: map[string]string{
			"auth/router.go": "Go", "auth/cache.go": "Go", "billing/controller.go": "Go",
			"gateway/route.go": "Go", "plugins/registry.go": "Go", "billing/ledger.go": "Go",
		},
		LanguageTiers:       map[string]string{"Go": "semantic"},
		PartialFailureFiles: map[string]bool{},
	}
}

func TestReconcileSixTiers(t *testing.T) {
	intent := Intent{ExtractedEntities: []ExtractedEntity{{Name: "AuthRouter.handleTimeout", Method: "regex", Confidence: 1}}}
	changes := []ChangedEntity{
		{Name: "handleTimeout", FilePath: "auth/router.go", StartLine: 10},
		{Name: "invalidate", FilePath: "auth/cache.go", StartLine: 20},
		{Name: "chargeRetry", FilePath: "billing/controller.go", StartLine: 30},
		{Name: "route", FilePath: "gateway/route.go", StartLine: 40},
		{Name: "dispatch", FilePath: "plugins/registry.go", StartLine: 50},
		{Name: "Ghost", FilePath: "ghost/file.go", StartLine: 1},
	}
	result := Reconcile(intent, changes, fixtureGraph(), DefaultConfig())

	if len(result.Confirmed) != 1 || result.Confirmed[0].Entity != "AuthRouter.handleTimeout" || result.Confirmed[0].CoverageConfidence != CoverageFull {
		t.Fatalf("confirmed = %#v", result.Confirmed)
	}
	if len(result.ExpectedBlastRadius) != 1 || result.ExpectedBlastRadius[0].Entity != "invalidate" ||
		result.ExpectedBlastRadius[0].Hops != 1 || result.ExpectedBlastRadius[0].MinEdgeClass != "deterministic" {
		t.Fatalf("expected_blast_radius = %#v", result.ExpectedBlastRadius)
	}
	if len(result.UndeclaredScopeCreep) != 1 || result.UndeclaredScopeCreep[0].Entity != "chargeRetry" ||
		result.UndeclaredScopeCreep[0].CoverageConfidence != CoverageFull {
		t.Fatalf("undeclared_scope_creep = %#v", result.UndeclaredScopeCreep)
	}
	if len(result.AdvisoryLowConfidence) != 1 || result.AdvisoryLowConfidence[0].Entity != "route" ||
		result.AdvisoryLowConfidence[0].CoverageConfidence != CoverageFull {
		t.Fatalf("advisory_low_confidence = %#v", result.AdvisoryLowConfidence)
	}
	if len(result.UnverifiableCoverage) != 2 {
		t.Fatalf("unverifiable_coverage = %#v", result.UnverifiableCoverage)
	}
	byEntity := map[string]UnverifiableCoverageEntity{}
	for _, entry := range result.UnverifiableCoverage {
		byEntity[entry.Entity] = entry
	}
	if byEntity["dispatch"].CoverageConfidence != CoveragePartial {
		t.Fatalf("dispatch coverage = %#v", byEntity["dispatch"])
	}
	if byEntity["Ghost"].CoverageConfidence != CoverageUnknown {
		t.Fatalf("Ghost coverage = %#v", byEntity["Ghost"])
	}
	for _, entry := range result.UnverifiableCoverage {
		if entry.VerificationPath != "manual_review_recommended" {
			t.Fatalf("verification_path = %#v", entry)
		}
	}
}

// TestReconcileNeverUpgradesPartialCoverageToScopeCreep pins #34: an entity
// that would previously have been flagged scope creep must route to
// unverifiable_coverage when its coverage is not full, in both directions —
// this does not touch a fully-resolved scope-creep case.
func TestReconcileNeverUpgradesPartialCoverageToScopeCreep(t *testing.T) {
	graph := fixtureGraph()
	changes := []ChangedEntity{
		{Name: "dispatch", FilePath: "plugins/registry.go", StartLine: 50},      // partial coverage
		{Name: "chargeRetry", FilePath: "billing/controller.go", StartLine: 30}, // full coverage
	}
	result := Reconcile(Intent{}, changes, graph, DefaultConfig())
	if len(result.UndeclaredScopeCreep) != 1 || result.UndeclaredScopeCreep[0].Entity != "chargeRetry" {
		t.Fatalf("a fully-resolved scope-creep case must still tier as before: %#v", result.UndeclaredScopeCreep)
	}
	if len(result.UnverifiableCoverage) != 1 || result.UnverifiableCoverage[0].Entity != "dispatch" {
		t.Fatalf("a partial-coverage entity must never land in undeclared_scope_creep: %#v", result)
	}
}

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
