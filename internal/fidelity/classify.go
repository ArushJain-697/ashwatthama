package fidelity

// EdgeClass is Fidelity's consumer-side trust category. It is derived from
// the provider's numeric confidence and resolution fields; it is not claimed
// to be a native Entire Graph tri-state.
type EdgeClass string

const (
	EdgeDeterministic EdgeClass = "deterministic"
	EdgeAdvisory      EdgeClass = "advisory"
)

type GraphRelation struct {
	FromID     string
	ToID       string
	Type       string
	Confidence float64
	Resolution string
}

func ClassifyRelation(relation GraphRelation, config Config) EdgeClass {
	if relation.Confidence >= config.EdgeClasses.DeterministicConfidenceGTE &&
		(relation.Resolution == "exact" || relation.Resolution == "import_resolved" || relation.Resolution == "package") {
		return EdgeDeterministic
	}
	return EdgeAdvisory
}
