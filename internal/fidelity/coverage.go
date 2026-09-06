package fidelity

// coverage.go implements Build Map #33: the coverage-confidence signal, the
// mandatory Curveball (Track 2, "Graph Is Evidence, Not an Oracle") response.
//
// This is orthogonal to EdgeClass (classify.go): EdgeClass classifies
// confidence in a relation the graph DID record. CoverageConfidence
// classifies whether the graph's SILENCE about an entity means anything at
// all — the gap the Curveball names: a changed entity with zero relations was
// being filed as confident undeclared_scope_creep regardless of WHY the graph
// was silent (genuinely disconnected code, vs. dynamic dispatch / generated
// code / reflection that static analysis cannot follow).
type CoverageConfidence string

const (
	// CoverageFull: the entity's relationships are backed by normal static
	// extraction; a zero-hop result from here is a real "not connected."
	CoverageFull CoverageConfidence = "full"
	// CoveragePartial: the entity sits in a region static analysis is known
	// to degrade in, or has zero extracted-class edges despite real code
	// presence with no stronger signal ruling that out.
	CoveragePartial CoverageConfidence = "partial"
	// CoverageUnknown: coverage itself could not be determined this run — the
	// snapshot returned no data at all for this entity's file.
	CoverageUnknown CoverageConfidence = "unknown"
)

// ComputeCoverage implements both detection mechanisms named in config
// coverage.detect_via. "capabilities_api" reads the provider's own
// completeness signals (file presence, language tier, partial failures)
// before falling back to edge counting; "heuristic_zero_edge" is the edge
// count alone, exactly as named, with no unknown case — #7/#37 confirmed
// capabilities_api is real for this provider, so it is the config default.
func ComputeCoverage(entity ChangedEntity, symbolID string, symbolFound bool, graph GraphSnapshot, config Config) (CoverageConfidence, string) {
	// Coverage counts ANY relation the graph recorded for this symbol,
	// deterministic or advisory — an inferred/pattern-matched edge is still
	// evidence the parser looked at this code and found something, just
	// something uncertain. That uncertainty is EdgeClass's job (orthogonal,
	// per §10.2). Coverage's only question is whether the graph recorded
	// nothing about this entity at all.
	edgeCount := countGraphEdges(symbolID, symbolFound, graph.Relations)

	if config.Coverage.DetectVia == "heuristic_zero_edge" {
		if edgeCount > 0 {
			return CoverageFull, ""
		}
		return zeroEdgeCoverage(config), "zero graph edges detected for this symbol (heuristic_zero_edge: code presence not otherwise checked)"
	}

	// capabilities_api: consult the provider's own signals first.
	language, fileKnown := graph.FileLanguages[entity.FilePath]
	if !fileKnown {
		return CoverageUnknown, "file absent from the graph snapshot; the snapshot returned no data for this region, so coverage could not be determined"
	}
	if graph.PartialFailureFiles[entity.FilePath] {
		return CoverageUnknown, "the graph snapshot reported a partial failure for this file, so its relations may be incomplete"
	}
	if tier := graph.LanguageTiers[language]; tier == "inventory-only" {
		return CoveragePartial, "this file's language (" + language + ") is inventory-only: the graph records file/symbol structure but does not attempt relationship extraction for it"
	}
	if !symbolFound {
		return CoveragePartial, "entity could not be grounded to a symbol in the graph snapshot; it may be dynamically dispatched, generated, or otherwise unresolvable to static analysis"
	}
	if edgeCount > 0 {
		return CoverageFull, ""
	}
	return zeroEdgeCoverage(config), "zero graph edges detected for this symbol despite real code presence; static analysis may not have resolved dynamic dispatch, generated code, or reflection here"
}

func zeroEdgeCoverage(config Config) CoverageConfidence {
	if config.Coverage.TreatZeroEdgeAs == "full" {
		return CoverageFull
	}
	return CoveragePartial
}

// structuralRelationTypes are containment/declaration facts the provider
// attaches to every symbol regardless of whether anything ever references it
// (a method always CONTAINS-relates to its type). Counting them toward
// coverage would make every symbol "full" by construction — coverage asks
// whether the graph resolved a REFERENCE to/from this entity, not whether the
// entity exists in the tree at all.
var structuralRelationTypes = map[string]bool{
	"DEFINES":  true,
	"CONTAINS": true,
}

func countGraphEdges(symbolID string, symbolFound bool, relations []GraphRelation) int {
	if !symbolFound {
		return 0
	}
	count := 0
	for _, relation := range relations {
		if structuralRelationTypes[relation.Type] {
			continue
		}
		if relation.FromID == symbolID || relation.ToID == symbolID {
			count++
		}
	}
	return count
}
