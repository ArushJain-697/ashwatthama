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
		Provider:                   capabilities.Provider,
		SchemaVersion:              capabilities.SchemaVersion,
		RelationTypes:              relations,
		ConfidenceModel:            "numeric confidence with resolution, evidence, and warning codes",
		SupportsTwoHopNeighbors:    true,
		RequiredRelationsAvailable: relations["CALLS"] && relations["DATA_FLOWS"] && relations["TESTS"] && relations["HANDLES_ROUTE"] && relations["HANDLES_GRPC"],
	}
}
