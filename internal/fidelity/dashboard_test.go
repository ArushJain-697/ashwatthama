package fidelity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestRenderDashboardEmbedsValidVerdictAndSatisfiesThreeWaySplit extends
// #41's distinctness requirement to the dashboard renderer: it must embed a
// verdict.json a viewer's own JS can parse back out, with no server, and
// must independently satisfy the mandatory three-way visual split.
func TestRenderDashboardEmbedsValidVerdictAndSatisfiesThreeWaySplit(t *testing.T) {
	verdict := BuildVerdict(fixtureVerification(), CurrentPreflightReport(), DefaultConfig(), time.Now())
	html, err := RenderDashboard(verdict)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(html, `id="verdict-data"`) {
		t.Fatal("dashboard did not embed verdict-data script tag")
	}
	if strings.Contains(html, "fetch(") {
		t.Fatal("dashboard must embed verdict.json inline, not fetch() it at runtime")
	}

	// Extract and round-trip the embedded JSON exactly as a browser's JS would.
	start := strings.Index(html, `id="verdict-data" type="application/json">`) + len(`id="verdict-data" type="application/json">`)
	end := strings.Index(html[start:], "</script>")
	embedded := html[start : start+end]
	var roundTrip Verdict
	if err := json.Unmarshal([]byte(embedded), &roundTrip); err != nil {
		t.Fatalf("embedded verdict JSON did not parse: %v\n%s", err, embedded)
	}
	if roundTrip.CheckpointID != verdict.CheckpointID {
		t.Fatalf("round-tripped checkpoint_id = %q, want %q", roundTrip.CheckpointID, verdict.CheckpointID)
	}
	if roundTrip.Summary.UnverifiableCoverageCount != verdict.Summary.UnverifiableCoverageCount {
		t.Fatalf("round-tripped unverifiable_coverage_count = %d, want %d", roundTrip.Summary.UnverifiableCoverageCount, verdict.Summary.UnverifiableCoverageCount)
	}

	for _, label := range []string{labelConfirmed, labelHeuristic, labelUnverifiable} {
		if !strings.Contains(html, label) {
			t.Fatalf("dashboard missing distinct visual label %q", label)
		}
	}
	if !strings.Contains(html, "Graph View") || !strings.Contains(html, "Summary") {
		t.Fatal("dashboard missing the summary/graph-view tab structure")
	}
}

// TestDashboardGraphNeverFabricatesEdgesOutsideBlastRadius pins the honesty
// property: an advisory or scope-creep entity must appear as a node but
// never gain an edge the graph never actually computed.
func TestDashboardGraphNeverFabricatesEdgesOutsideBlastRadius(t *testing.T) {
	verdict := BuildVerdict(fixtureVerification(), CurrentPreflightReport(), DefaultConfig(), time.Now())
	nodes, edges := dashboardGraph(verdict)

	nodeIDs := map[string]bool{}
	for _, node := range nodes {
		nodeIDs[node.ID] = true
	}
	for _, entry := range verdict.Reconciliation.AdvisoryLowConfidence {
		if !nodeIDs[entry.Entity] {
			t.Fatalf("advisory entity %q missing from graph nodes", entry.Entity)
		}
	}
	if len(verdict.Reconciliation.ExpectedBlastRadius) == 0 && len(edges) != 0 {
		t.Fatalf("no expected_blast_radius entries in this fixture, but %d edges were drawn: %#v", len(edges), edges)
	}
}
