package fidelity

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// community.go implements Build Map #57: community-aware scope-creep
// framing. Leiden community detection itself runs out-of-process
// (scripts/leiden_communities.py, leidenalg + python-igraph, the pinned
// Bible implementation) rather than being reimplemented in Go — this is a
// stretch enhancement the buildmap explicitly permits skipping when the
// repository is too small for meaningful communities, and a robust
// in-process Go/Python bridge is real engineering this ticket does not ask
// for. What's wired here is the consuming half: loading that script's
// output and annotating undeclared_scope_creep entries with it, entirely
// optionally — no community map means no labels, unchanged behavior.
//
// Verified non-trivial on this repository (see BUILDATHON.md): 149
// communities, largest holding 12.9% of the graph, 0% singletons.

// CommunityMap is scripts/leiden_communities.py's output, loaded as-is.
type CommunityMap struct {
	Communities     map[string]int    `json:"communities"`
	CommunityLabels map[string]string `json:"community_labels"`
}

// LoadCommunityMap reads a communities.json produced by
// scripts/leiden_communities.py. A missing or malformed file is the caller's
// signal to skip annotation, not a pipeline failure — this stretch
// enhancement must never block the mandatory reconciliation it decorates.
func LoadCommunityMap(path string) (CommunityMap, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return CommunityMap{}, fmt.Errorf("read community map: %w", err)
	}
	var communityMap CommunityMap
	if err := json.Unmarshal(content, &communityMap); err != nil {
		return CommunityMap{}, fmt.Errorf("parse community map: %w", err)
	}
	return communityMap, nil
}

// AnnotateScopeCreepCommunities fills in CommunityLabel on every
// undeclared_scope_creep entry it can resolve to a graph symbol, turning "3
// unrelated-looking file paths" into "3 entities in the Billing community."
// Entries it cannot resolve (or when communityMap is empty) are left
// unchanged — this never removes or reorders reconciliation results.
func AnnotateScopeCreepCommunities(reconciliation *Reconciliation, changes []ChangedEntity, graph GraphSnapshot, communityMap CommunityMap) {
	if len(communityMap.Communities) == 0 {
		return
	}
	for i, entry := range reconciliation.UndeclaredScopeCreep {
		for _, change := range changes {
			if change.FilePath != entry.File || change.StartLine != entry.Line {
				continue
			}
			symbolID, found := symbolIDForEntity(change, graph.Symbols)
			if !found {
				continue
			}
			communityIndex, inCommunity := communityMap.Communities[symbolID]
			if !inCommunity {
				continue
			}
			if label, hasLabel := communityMap.CommunityLabels[strconv.Itoa(communityIndex)]; hasLabel {
				reconciliation.UndeclaredScopeCreep[i].CommunityLabel = label
			}
			break
		}
	}
}
