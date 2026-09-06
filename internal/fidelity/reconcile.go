package fidelity

import (
	"strconv"
	"strings"
)

type EntityLocation struct {
	Entity string `json:"entity"`
	File   string `json:"file,omitempty"`
	Line   int    `json:"line,omitempty"`
	// CoverageConfidence is stamped even on success: a confirmed match still
	// tells a reviewer how much of the surrounding graph it rests on (#33).
	CoverageConfidence CoverageConfidence `json:"coverage_confidence"`
}

type UnimplementedEntity struct {
	Entity string `json:"entity"`
	Reason string `json:"reason"`
}

// BlastRadiusEntity is one expected_blast_radius entry (#23): a changed
// entity reachable from a confirmed one within blast_radius_hops via
// deterministic edges only. Never rests on an inferred/ambiguous hop.
type BlastRadiusEntity struct {
	Entity             string             `json:"entity"`
	File               string             `json:"file,omitempty"`
	Line               int                `json:"line,omitempty"`
	PathFromConfirmed  []string           `json:"path_from_confirmed"`
	Hops               int                `json:"hops"`
	EdgeTypesTraversed []string           `json:"edge_types_traversed"`
	MinEdgeClass       string             `json:"min_edge_class"`
	CoverageConfidence CoverageConfidence `json:"coverage_confidence"`
}

// ScopeCreepEntity is one undeclared_scope_creep entry (#24): unreachable
// from any confirmed entity within the hop limit, AND its coverage is full —
// #34's never-silently-upgrade rule means this tier is reachable only when
// the absence of a path is itself trustworthy evidence.
type ScopeCreepEntity struct {
	Entity             string             `json:"entity"`
	File               string             `json:"file,omitempty"`
	Line               int                `json:"line,omitempty"`
	Reason             string             `json:"reason"`
	CoverageConfidence CoverageConfidence `json:"coverage_confidence"`
}

// AdvisoryEntity is one advisory_low_confidence entry (#25): reachable from a
// confirmed entity only via an inferred/ambiguous-equivalent (advisory) edge.
// The Databricks-reserved fields stay nil/zero-value until the calibration
// module (§12) populates them; never silently promoted to expected_blast_radius.
type AdvisoryEntity struct {
	Entity               string             `json:"entity"`
	File                 string             `json:"file,omitempty"`
	Line                 int                `json:"line,omitempty"`
	EdgeClass            string             `json:"edge_class"`
	Reason               string             `json:"reason"`
	CoverageConfidence   CoverageConfidence `json:"coverage_confidence"`
	VerbalConfidence     *float64           `json:"verbal_confidence"`
	ConformalGatePassed  *bool              `json:"conformal_gate_passed"`
	CalibratedConfidence *float64           `json:"calibrated_confidence,omitempty"`
	Corroboration        Corroboration      `json:"corroboration"`
}

// UnverifiableCoverageEntity is the new sixth tier (#34, #36): an entity that
// would otherwise have landed in undeclared_scope_creep on the strength of an
// absent path alone, but whose coverage is partial or unknown. Never
// conflated with AdvisoryEntity: that tier means "the graph has an edge but
// isn't sure about it"; this one means "the graph may not have tried to
// record an edge here at all."
type UnverifiableCoverageEntity struct {
	Entity             string             `json:"entity"`
	File               string             `json:"file,omitempty"`
	Line               int                `json:"line,omitempty"`
	CoverageConfidence CoverageConfidence `json:"coverage_confidence"`
	CoverageReason     string             `json:"coverage_reason"`
	// VerificationPath is the baseline safety guarantee (#35): populated with
	// or without Databricks. "manual_review_recommended" is the fallback.
	VerificationPath string        `json:"verification_path"`
	Corroboration    Corroboration `json:"corroboration"`
}

// Corroboration is the Databricks AI Search corroboration result (§12.4,
// §12.15), shared by advisory edges and unverifiable_coverage entries.
// Attempted is false and Verified/EvidenceSnippet are nil until the
// Databricks module runs; the schema shape does not change either way.
type Corroboration struct {
	Attempted       bool    `json:"attempted"`
	Verified        *bool   `json:"verified"`
	EvidenceSnippet *string `json:"evidence_snippet"`
}

// Reconciliation is Stage 4's full six-tier output (#22-#25, #33-#34).
type Reconciliation struct {
	Confirmed             []EntityLocation             `json:"confirmed"`
	DeclaredUnimplemented []UnimplementedEntity        `json:"declared_unimplemented"`
	ExpectedBlastRadius   []BlastRadiusEntity          `json:"expected_blast_radius"`
	UndeclaredScopeCreep  []ScopeCreepEntity           `json:"undeclared_scope_creep"`
	AdvisoryLowConfidence []AdvisoryEntity             `json:"advisory_low_confidence"`
	UnverifiableCoverage  []UnverifiableCoverageEntity `json:"unverifiable_coverage"`
}

// DirectReconciliation is the original two-tier result (confirmed,
// declared_unimplemented) kept for its own tests and as the input to the
// fuller Reconcile below; it performs no graph reachability itself.
type DirectReconciliation struct {
	Confirmed             []EntityLocation      `json:"confirmed"`
	DeclaredUnimplemented []UnimplementedEntity `json:"declared_unimplemented"`
}

// ReconcileDirectClaims performs the first two deterministic tiers. Matching
// permits a graph diff's short entity name to satisfy a transcript's qualified
// symbol name, but never performs fuzzy equivalence; that remains an explicit
// later policy decision.
func ReconcileDirectClaims(intent Intent, changes []ChangedEntity) DirectReconciliation {
	result := DirectReconciliation{
		Confirmed:             []EntityLocation{},
		DeclaredUnimplemented: []UnimplementedEntity{},
	}
	for _, claim := range intent.ExtractedEntities {
		matched := false
		for _, change := range changes {
			if !sameEntityName(claim.Name, change.Name) {
				continue
			}
			matched = true
			result.Confirmed = append(result.Confirmed, EntityLocation{
				Entity: claim.Name, File: change.FilePath, Line: change.StartLine,
			})
		}
		if !matched {
			result.DeclaredUnimplemented = append(result.DeclaredUnimplemented, UnimplementedEntity{
				Entity: claim.Name, Reason: "named in intent, no matching graph change detected",
			})
		}
	}
	return result
}

// Reconcile is the full Stage 4: direct claims, then N-hop reachability
// (#21) tiering every remaining changed entity, with coverage-confidence
// (#33) checked BEFORE any entity can be filed as undeclared_scope_creep
// (#34) — the never-silently-upgrade rule the Curveball requires.
func Reconcile(intent Intent, changes []ChangedEntity, graph GraphSnapshot, config Config) Reconciliation {
	direct := ReconcileDirectClaims(intent, changes)

	confirmedIDs := map[string]bool{}
	for _, change := range changes {
		if !isConfirmedChange(change, direct.Confirmed) {
			continue
		}
		if id, ok := symbolIDForEntity(change, graph.Symbols); ok {
			confirmedIDs[id] = true
		}
	}
	var confirmedIDList []string
	for id := range confirmedIDs {
		confirmedIDList = append(confirmedIDList, id)
	}

	deterministicHits := reachable(confirmedIDList, graph.Relations, config, EdgeDeterministic, config.BlastRadiusHops)
	anyEdgeHits := reachable(confirmedIDList, graph.Relations, config, "", config.BlastRadiusHops)

	result := Reconciliation{
		Confirmed:             stampConfirmedCoverage(direct.Confirmed, changes, graph, config),
		DeclaredUnimplemented: direct.DeclaredUnimplemented,
		ExpectedBlastRadius:   []BlastRadiusEntity{},
		UndeclaredScopeCreep:  []ScopeCreepEntity{},
		AdvisoryLowConfidence: []AdvisoryEntity{},
		UnverifiableCoverage:  []UnverifiableCoverageEntity{},
	}

	for _, change := range changes {
		if isConfirmedChange(change, direct.Confirmed) {
			continue
		}
		symbolID, symbolFound := symbolIDForEntity(change, graph.Symbols)
		coverage, coverageReason := ComputeCoverage(change, symbolID, symbolFound, graph, config)

		if symbolFound {
			if hit, ok := deterministicHits[symbolID]; ok {
				result.ExpectedBlastRadius = append(result.ExpectedBlastRadius, BlastRadiusEntity{
					Entity: change.Name, File: change.FilePath, Line: change.StartLine,
					PathFromConfirmed: hit.Path, Hops: hit.Hops,
					EdgeTypesTraversed: hit.EdgeTypesTraversed, MinEdgeClass: "deterministic",
					CoverageConfidence: coverage,
				})
				continue
			}
		}

		// #34: coverage is checked BEFORE undeclared_scope_creep can be
		// assigned on the strength of an absent path alone.
		if coverage != CoverageFull {
			result.UnverifiableCoverage = append(result.UnverifiableCoverage, UnverifiableCoverageEntity{
				Entity: change.Name, File: change.FilePath, Line: change.StartLine,
				CoverageConfidence: coverage, CoverageReason: coverageReason,
				VerificationPath: config.Coverage.FallbackVerificationPath,
				Corroboration:    Corroboration{Attempted: false},
			})
			continue
		}

		if symbolFound {
			if hit, ok := anyEdgeHits[symbolID]; ok {
				result.AdvisoryLowConfidence = append(result.AdvisoryLowConfidence, AdvisoryEntity{
					Entity: change.Name, File: change.FilePath, Line: change.StartLine,
					EdgeClass:          "advisory",
					Reason:             "reachable from a confirmed entity within " + hopsWord(hit.Hops) + " only via an inferred/ambiguous edge",
					CoverageConfidence: coverage,
					Corroboration:      Corroboration{Attempted: false},
				})
				continue
			}
		}

		result.UndeclaredScopeCreep = append(result.UndeclaredScopeCreep, ScopeCreepEntity{
			Entity: change.Name, File: change.FilePath, Line: change.StartLine,
			Reason:             "changed, not reachable from any confirmed entity within " + hopsWord(config.BlastRadiusHops),
			CoverageConfidence: coverage,
		})
	}
	return result
}

func stampConfirmedCoverage(confirmed []EntityLocation, changes []ChangedEntity, graph GraphSnapshot, config Config) []EntityLocation {
	stamped := make([]EntityLocation, len(confirmed))
	for i, entry := range confirmed {
		stamped[i] = entry
		for _, change := range changes {
			if change.FilePath != entry.File || change.StartLine != entry.Line {
				continue
			}
			symbolID, symbolFound := symbolIDForEntity(change, graph.Symbols)
			coverage, _ := ComputeCoverage(change, symbolID, symbolFound, graph, config)
			stamped[i].CoverageConfidence = coverage
			break
		}
	}
	return stamped
}

func isConfirmedChange(change ChangedEntity, confirmed []EntityLocation) bool {
	for _, entry := range confirmed {
		if entry.File == change.FilePath && entry.Line == change.StartLine {
			return true
		}
	}
	return false
}

func hopsWord(hops int) string {
	if hops == 1 {
		return "1 hop"
	}
	return strconv.Itoa(hops) + " hops"
}

func sameEntityName(claim, changed string) bool {
	if claim == changed {
		return true
	}
	lastDot := strings.LastIndex(claim, ".")
	return lastDot >= 0 && claim[lastDot+1:] == changed
}
