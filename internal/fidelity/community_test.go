package fidelity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAnnotateScopeCreepCommunitiesLabelsResolvedEntity(t *testing.T) {
	graph := fixtureGraph() // C (BillingController.chargeRetry) is the unreachable scope-creep case
	changes := []ChangedEntity{
		{Name: "chargeRetry", FilePath: "billing/controller.go", StartLine: 30},
	}
	reconciliation := Reconcile(Intent{}, changes, graph, DefaultConfig())
	if len(reconciliation.UndeclaredScopeCreep) != 1 {
		t.Fatalf("fixture setup: undeclared_scope_creep = %#v", reconciliation.UndeclaredScopeCreep)
	}
	if reconciliation.UndeclaredScopeCreep[0].CommunityLabel != "" {
		t.Fatalf("community label populated before annotation: %#v", reconciliation.UndeclaredScopeCreep[0])
	}

	communityMap := CommunityMap{
		Communities:     map[string]int{"C": 3},
		CommunityLabels: map[string]string{"3": "billing"},
	}
	AnnotateScopeCreepCommunities(&reconciliation, changes, graph, communityMap)

	if got := reconciliation.UndeclaredScopeCreep[0].CommunityLabel; got != "billing" {
		t.Fatalf("community_label = %q, want %q", got, "billing")
	}
}

func TestAnnotateScopeCreepCommunitiesNoOpWithoutMap(t *testing.T) {
	graph := fixtureGraph()
	changes := []ChangedEntity{{Name: "chargeRetry", FilePath: "billing/controller.go", StartLine: 30}}
	reconciliation := Reconcile(Intent{}, changes, graph, DefaultConfig())
	AnnotateScopeCreepCommunities(&reconciliation, changes, graph, CommunityMap{})
	if reconciliation.UndeclaredScopeCreep[0].CommunityLabel != "" {
		t.Fatalf("an empty community map must be a no-op: %#v", reconciliation.UndeclaredScopeCreep[0])
	}
}

func TestLoadCommunityMapRoundTrips(t *testing.T) {
	path := filepath.Join(t.TempDir(), "communities.json")
	content, err := json.Marshal(CommunityMap{
		Communities:     map[string]int{"local/1:Go:a.go:function:F": 2},
		CommunityLabels: map[string]string{"2": "internal/fidelity"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadCommunityMap(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.CommunityLabels["2"] != "internal/fidelity" {
		t.Fatalf("loaded = %#v", loaded)
	}
}

func TestLoadCommunityMapMissingFileFailsClosed(t *testing.T) {
	if _, err := LoadCommunityMap(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatal("LoadCommunityMap accepted a missing file")
	}
}
