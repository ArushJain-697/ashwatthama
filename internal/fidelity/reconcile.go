package fidelity

import "strings"

type EntityLocation struct {
	Entity string `json:"entity"`
	File   string `json:"file,omitempty"`
	Line   int    `json:"line,omitempty"`
}

type UnimplementedEntity struct {
	Entity string `json:"entity"`
	Reason string `json:"reason"`
}

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

func sameEntityName(claim, changed string) bool {
	if claim == changed {
		return true
	}
	lastDot := strings.LastIndex(claim, ".")
	return lastDot >= 0 && claim[lastDot+1:] == changed
}
