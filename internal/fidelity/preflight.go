package fidelity

import "github.com/entireio/entire-graph/internal/sem"

// PreflightReport records the graph facts Fidelity relies on before it attempts
// reconciliation.  It is intentionally derived from the installed provider,
// not copied from the project plan, so a changed provider cannot silently alter
// Fidelity's trust policy.
type PreflightReport struct {
	Provider                   string          `json:"provider"`
	SchemaVersion              string          `json:"schema_version"`
	RelationTypes              map[string]bool `json:"relation_types"`
	ConfidenceModel            string          `json:"confidence_model"`
	SupportsTwoHopNeighbors    bool            `json:"supports_two_hop_neighbors"`
	RequiredRelationsAvailable bool            `json:"required_relations_available"`
}

// CurrentPreflightReport is the capability portion of Build Map ticket #7.
// Relation evidence itself remains per-snapshot because its availability is
// profile- and language-dependent.
func CurrentPreflightReport() PreflightReport {
	capabilities := sem.Capabilities()
	relations := make(map[string]bool, len(capabilities.SupportedRelationTypes))
	for _, relation := range capabilities.SupportedRelationTypes {
		relations[relation] = true
	}
	return PreflightReport{
		Provider:                capabilities.Provider,
		SchemaVersion:           capabilities.SchemaVersion,
		RelationTypes:           relations,
		ConfidenceModel:         "numeric confidence with resolution, evidence, and warning codes",
		SupportsTwoHopNeighbors: true,
		// CALLS is the only relation type Fidelity's own reachability logic
		// (reachability.go) structurally depends on: it walks whatever
		// relations the snapshot returns, generically, for both the
		// deterministic-only and any-edge passes. Requiring TESTS/
		// HANDLES_ROUTE/HANDLES_GRPC here was a v1 assumption that the
		// pre-flight record already disproved for Go in the full profile
		// (BUILDATHON.md), and Fidelity never actually reads those three
		// relation types — gating on them would report false-red for a
		// perfectly usable Go repository.
		RequiredRelationsAvailable: relations["CALLS"],
	}
}
