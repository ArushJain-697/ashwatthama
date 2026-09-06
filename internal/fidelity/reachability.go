package fidelity

// reachability.go implements Build Map #21: N-hop reachability over the
// relation graph, computed twice — once restricted to deterministic edges
// (for expected_blast_radius, #23), once over every edge (to detect the
// advisory-only case, #25). The two are kept as separate BFS passes rather
// than one pass with a "best edge class seen" merge, because #23 requires
// deterministic-edge-only reasoning never to rest on an inferred/ambiguous
// hop — mixing them in one traversal risks a deterministic-looking path that
// silently used an advisory edge partway through.

// reachabilityHit is one entity's shortest path back to a confirmed entity.
type reachabilityHit struct {
	Hops               int
	Path               []string // symbol IDs, confirmed entity first
	EdgeTypesTraversed []string
	MinEdgeClass       string // "deterministic" or "advisory"
}

// symbolIndex resolves a ChangedEntity to the graph symbol ID reachability
// walks over. Matching permits a diff's short entity name to satisfy a
// snapshot's qualified name (sameEntityName, reused from reconcile.go) plus a
// same-file check, since two files can each define a same-named symbol.
func symbolIDForEntity(entity ChangedEntity, symbols []Symbol) (string, bool) {
	for _, symbol := range symbols {
		if symbol.FilePath == entity.FilePath && sameEntityName(symbol.QualifiedName, entity.Name) {
			return symbol.ID, true
		}
	}
	return "", false
}

// reachable runs bounded BFS from every ID in from, over relations restricted
// to edgeClass (or every relation when edgeClass is empty), up to hops. The
// graph is walked as undirected: a change to either a caller or a callee of a
// confirmed entity is a real blast-radius candidate, matching the CLI
// `impact` command's own both-directions framing.
func reachable(from []string, relations []GraphRelation, config Config, edgeClass EdgeClass, hops int) map[string]reachabilityHit {
	type step struct {
		id, edgeType string
	}
	adjacency := map[string][]step{}
	for _, relation := range relations {
		if structuralRelationTypes[relation.Type] {
			// DEFINES/CONTAINS are containment, not reference: two methods of
			// the same type would otherwise look "reachable" from each other
			// through no actual call/data/type relationship at all.
			continue
		}
		if edgeClass != "" && ClassifyRelation(GraphRelation{Confidence: relation.Confidence, Resolution: relation.Resolution}, config) != edgeClass {
			continue
		}
		adjacency[relation.FromID] = append(adjacency[relation.FromID], step{relation.ToID, relation.Type})
		adjacency[relation.ToID] = append(adjacency[relation.ToID], step{relation.FromID, relation.Type})
	}
	result := map[string]reachabilityHit{}
	visited := map[string]bool{}
	type queued struct {
		id                 string
		depth              int
		path               []string
		edgeTypesTraversed []string
	}
	var queue []queued
	for _, id := range from {
		visited[id] = true
		queue = append(queue, queued{id: id, depth: 0, path: []string{id}})
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.depth >= hops {
			continue
		}
		for _, next := range adjacency[current.id] {
			if visited[next.id] {
				continue
			}
			visited[next.id] = true
			path := append(append([]string{}, current.path...), next.id)
			edgeTypes := append(append([]string{}, current.edgeTypesTraversed...), next.edgeType)
			minClass := "deterministic"
			if edgeClass == "" {
				// Whole-graph pass: report the weakest class actually crossed,
				// since a caller of the whole-graph result cares whether ANY
				// hop was advisory.
				minClass = string(EdgeAdvisory)
				for _, relation := range relations {
					if (relation.FromID == current.id && relation.ToID == next.id) || (relation.ToID == current.id && relation.FromID == next.id) {
						if ClassifyRelation(GraphRelation{Confidence: relation.Confidence, Resolution: relation.Resolution}, config) == EdgeDeterministic {
							minClass = string(EdgeDeterministic)
						}
					}
				}
			}
			result[next.id] = reachabilityHit{
				Hops: current.depth + 1, Path: path,
				EdgeTypesTraversed: edgeTypes, MinEdgeClass: minClass,
			}
			queue = append(queue, queued{id: next.id, depth: current.depth + 1, path: path, edgeTypesTraversed: edgeTypes})
		}
	}
	return result
}
