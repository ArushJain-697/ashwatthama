package fidelity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
)

// dashboard.go implements #40: dashboard.html, three-way-aware from the
// start. Per the user's explicit design choices: a tier-grouped summary tab
// (default) plus a Graph View tab for the entities that actually carry real
// path/edge data (expected_blast_radius), and verdict.json's content
// embedded inline at generation time rather than fetch()ed at runtime — no
// server, no CORS, opens by double-click. A pure function of Verdict, like
// every other renderer; no adapter or repo access.
//
// The Graph View draws an edge ONLY where a real path_from_confirmed exists
// (expected_blast_radius). Advisory/scope-creep/unverifiable-coverage nodes
// are placed and colored by tier but never connected by a fabricated edge —
// consistent with the whole project's refusal to present an absent
// relationship as a real one.

func RenderDashboard(verdict Verdict) (string, error) {
	verdictJSON, err := json.Marshal(verdict)
	if err != nil {
		return "", fmt.Errorf("encode verdict for dashboard: %w", err)
	}
	graphNodes, graphEdges := dashboardGraph(verdict)
	nodesJSON, err := json.Marshal(graphNodes)
	if err != nil {
		return "", fmt.Errorf("encode dashboard graph nodes: %w", err)
	}
	edgesJSON, err := json.Marshal(graphEdges)
	if err != nil {
		return "", fmt.Errorf("encode dashboard graph edges: %w", err)
	}

	data := dashboardTemplateData{
		VerdictJSON:    template.JS(verdictJSON),
		GraphNodes:     template.JS(nodesJSON),
		GraphEdges:     template.JS(edgesJSON),
		CheckpointID:   verdict.CheckpointID,
		VerdictLabel:   verdict.Summary.VerdictLabel,
		LabelConfirmed: labelConfirmed, LabelHeuristic: labelHeuristic,
		LabelUnverifiable: labelUnverifiable, LabelDrift: labelDrift,
		ConfirmedCount:    len(verdict.Reconciliation.Confirmed) + len(verdict.Reconciliation.ExpectedBlastRadius),
		HeuristicCount:    len(verdict.Reconciliation.AdvisoryLowConfidence),
		UnverifiableCount: len(verdict.Reconciliation.UnverifiableCoverage),
		DriftCount:        len(verdict.Reconciliation.UndeclaredScopeCreep) + len(verdict.Reconciliation.DeclaredUnimplemented),
	}
	var out bytes.Buffer
	if err := dashboardTemplate.Execute(&out, data); err != nil {
		return "", fmt.Errorf("render dashboard: %w", err)
	}
	return out.String(), nil
}

type dashboardNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Tier  string `json:"tier"`  // one of the four visual categories
	Color string `json:"color"` // matches the terminal/Markdown renderer palette
}

type dashboardEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

const (
	colorConfirmed    = "#3fb950"
	colorHeuristic    = "#d29922"
	colorUnverifiable = "#bc8cff"
	colorDrift        = "#f85149"
)

// dashboardGraph builds the Graph View's nodes/edges. Real edges exist only
// for expected_blast_radius (the only tier carrying path_from_confirmed);
// every other tier contributes unconnected, tier-colored nodes.
func dashboardGraph(verdict Verdict) ([]dashboardNode, []dashboardEdge) {
	var nodes []dashboardNode
	var edges []dashboardEdge
	seen := map[string]bool{}
	addNode := func(id, label, tier, color string) {
		if id == "" || seen[id] {
			return
		}
		seen[id] = true
		nodes = append(nodes, dashboardNode{ID: id, Label: label, Tier: tier, Color: color})
	}

	for _, entry := range verdict.Reconciliation.Confirmed {
		addNode(entry.Entity, entry.Entity, "confirmed", colorConfirmed)
	}
	for _, entry := range verdict.Reconciliation.ExpectedBlastRadius {
		addNode(entry.Entity, entry.Entity, "expected_blast_radius", colorConfirmed)
		for i := 0; i < len(entry.PathFromConfirmed)-1; i++ {
			from, to := lastSegment(entry.PathFromConfirmed[i]), lastSegment(entry.PathFromConfirmed[i+1])
			addNode(from, from, "confirmed", colorConfirmed)
			addNode(to, to, "expected_blast_radius", colorConfirmed)
			edgeType := ""
			if i < len(entry.EdgeTypesTraversed) {
				edgeType = entry.EdgeTypesTraversed[i]
			}
			edges = append(edges, dashboardEdge{Source: from, Target: to, Type: edgeType})
		}
	}
	for _, entry := range verdict.Reconciliation.AdvisoryLowConfidence {
		addNode(entry.Entity, entry.Entity, "advisory_low_confidence", colorHeuristic)
	}
	for _, entry := range verdict.Reconciliation.UnverifiableCoverage {
		addNode(entry.Entity, entry.Entity, "unverifiable_coverage", colorUnverifiable)
	}
	for _, entry := range verdict.Reconciliation.UndeclaredScopeCreep {
		addNode(entry.Entity, entry.Entity, "undeclared_scope_creep", colorDrift)
	}
	if nodes == nil {
		nodes = []dashboardNode{}
	}
	if edges == nil {
		edges = []dashboardEdge{}
	}
	return nodes, edges
}

// lastSegment reads a compound symbol ID the one positionally-safe way
// snapshot-format.md documents for a trailing field: anchor on the LAST
// separator. It is a display label, not an identity.
func lastSegment(id string) string {
	for i := len(id) - 1; i >= 0; i-- {
		if id[i] == ':' {
			return id[i+1:]
		}
	}
	return id
}

type dashboardTemplateData struct {
	VerdictJSON, GraphNodes, GraphEdges                           template.JS
	CheckpointID, VerdictLabel                                    string
	LabelConfirmed, LabelHeuristic, LabelUnverifiable, LabelDrift string
	ConfirmedCount, HeuristicCount, UnverifiableCount, DriftCount int
}

var dashboardTemplate = template.Must(template.New("dashboard").Parse(dashboardHTML))

const dashboardHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Fidelity Verdict Dashboard</title>
<style>
  :root { color-scheme: dark; }
  body { margin:0; font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif; background:#0d1117; color:#c9d1d9; }
  header { padding:16px 24px; border-bottom:1px solid #30363d; }
  header h1 { margin:0; font-size:18px; }
  header .label { color:#8b949e; font-size:13px; }
  nav.tabs { display:flex; gap:4px; padding:0 24px; border-bottom:1px solid #30363d; }
  nav.tabs button { background:none; border:none; color:#8b949e; padding:10px 16px; cursor:pointer; font-size:14px; border-bottom:2px solid transparent; }
  nav.tabs button.active { color:#c9d1d9; border-bottom-color:#58a6ff; }
  section { padding:20px 24px; }
  section[hidden] { display:none; }
  .cards { display:grid; grid-template-columns:repeat(4,1fr); gap:12px; margin-bottom:20px; }
  .card { border:1px solid #30363d; border-radius:8px; padding:14px; }
  .card .count { font-size:28px; font-weight:600; }
  .card .name { font-size:12px; color:#8b949e; margin-top:4px; }
  .tier { margin-bottom:18px; }
  .tier h3 { font-size:13px; text-transform:uppercase; letter-spacing:.04em; color:#8b949e; margin:0 0 6px; }
  .row { font-size:13px; padding:4px 0; border-bottom:1px solid #21262d; font-family:ui-monospace,monospace; }
  .row .loc { color:#8b949e; }
  #graph-svg { width:100%; height:70vh; border:1px solid #30363d; border-radius:8px; background:#0d1117; }
  .legend { display:flex; gap:16px; margin-top:10px; font-size:12px; color:#8b949e; }
  .legend span { display:inline-flex; align-items:center; gap:6px; }
  .dot { width:10px; height:10px; border-radius:50%; display:inline-block; }
  .note { font-size:12px; color:#8b949e; margin-top:8px; }
</style>
</head>
<body>
<header>
  <h1>Fidelity Verdict — {{.CheckpointID}}</h1>
  <div class="label">{{.VerdictLabel}}</div>
</header>
<nav class="tabs">
  <button data-tab="summary" class="active">Summary</button>
  <button data-tab="graph">Graph View</button>
</nav>
<section id="summary-tab">
  <div class="cards">
    <div class="card" style="border-color:#3fb950"><div class="count" style="color:#3fb950">{{.ConfirmedCount}}</div><div class="name">{{.LabelConfirmed}}</div></div>
    <div class="card" style="border-color:#d29922"><div class="count" style="color:#d29922">{{.HeuristicCount}}</div><div class="name">{{.LabelHeuristic}}</div></div>
    <div class="card" style="border-color:#bc8cff"><div class="count" style="color:#bc8cff">{{.UnverifiableCount}}</div><div class="name">{{.LabelUnverifiable}}</div></div>
    <div class="card" style="border-color:#f85149"><div class="count" style="color:#f85149">{{.DriftCount}}</div><div class="name">{{.LabelDrift}}</div></div>
  </div>
  <div id="tier-lists"></div>
</section>
<section id="graph-tab" hidden>
  <svg id="graph-svg"></svg>
  <div class="legend">
    <span><i class="dot" style="background:#3fb950"></i>{{.LabelConfirmed}}</span>
    <span><i class="dot" style="background:#d29922"></i>{{.LabelHeuristic}}</span>
    <span><i class="dot" style="background:#bc8cff"></i>{{.LabelUnverifiable}}</span>
    <span><i class="dot" style="background:#f85149"></i>{{.LabelDrift}}</span>
  </div>
  <p class="note">Edges are drawn only where a real graph path exists (expected_blast_radius). Every other node is placed by tier and colored accordingly, but is never connected by a relationship the graph did not actually compute.</p>
</section>
<script id="verdict-data" type="application/json">{{.VerdictJSON}}</script>
<script id="graph-nodes-data" type="application/json">{{.GraphNodes}}</script>
<script id="graph-edges-data" type="application/json">{{.GraphEdges}}</script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/d3/7.9.0/d3.min.js"></script>
<script>
const verdict = JSON.parse(document.getElementById('verdict-data').textContent);
const graphNodes = JSON.parse(document.getElementById('graph-nodes-data').textContent);
const graphEdges = JSON.parse(document.getElementById('graph-edges-data').textContent);

document.querySelectorAll('nav.tabs button').forEach(btn => {
  btn.addEventListener('click', () => {
    document.querySelectorAll('nav.tabs button').forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    document.getElementById('summary-tab').hidden = btn.dataset.tab !== 'summary';
    document.getElementById('graph-tab').hidden = btn.dataset.tab !== 'graph';
    if (btn.dataset.tab === 'graph') drawGraph();
  });
});

function escapeHTML(s) {
  return String(s).replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
}

function renderTierLists() {
  const container = document.getElementById('tier-lists');
  const tiers = [
    ['Confirmed', verdict.reconciliation.confirmed, e => e.entity + ' — ' + (e.file||'') + ':' + (e.line||'') + ' [coverage: ' + e.coverage_confidence + ']'],
    ['Expected Blast Radius', verdict.reconciliation.expected_blast_radius, e => e.entity + ' — ' + (e.file||'') + ':' + (e.line||'') + ' — ' + e.hops + ' hop(s)'],
    ['Advisory Low Confidence', verdict.reconciliation.advisory_low_confidence, e => e.entity + ' — ' + e.reason],
    ['Unverifiable Coverage', verdict.reconciliation.unverifiable_coverage, e => e.entity + ' — ' + e.coverage_reason + ' [' + e.verification_path + ']'],
    ['Undeclared Scope Creep', verdict.reconciliation.undeclared_scope_creep, e => e.entity + ' — ' + e.reason + (e.community_label ? ' [community: ' + e.community_label + ']' : '')],
    ['Declared Unimplemented', verdict.reconciliation.declared_unimplemented, e => e.entity + ' — ' + e.reason],
  ];
  container.innerHTML = tiers.map(([title, entries, fmt]) => {
    if (!entries || entries.length === 0) return '';
    const rows = entries.map(e => '<div class="row">' + escapeHTML(fmt(e)) + '</div>').join('');
    return '<div class="tier"><h3>' + title + ' (' + entries.length + ')</h3>' + rows + '</div>';
  }).join('');
}
renderTierLists();

let graphDrawn = false;
function drawGraph() {
  if (graphDrawn) return;
  graphDrawn = true;
  const svg = d3.select('#graph-svg');
  const width = svg.node().clientWidth || 800;
  const height = svg.node().clientHeight || 500;
  svg.attr('viewBox', [0, 0, width, height]);

  const simulation = d3.forceSimulation(graphNodes)
    .force('link', d3.forceLink(graphEdges).id(d => d.id).distance(80))
    .force('charge', d3.forceManyBody().strength(-120))
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('collide', d3.forceCollide(24));

  const link = svg.append('g').selectAll('line').data(graphEdges).join('line')
    .attr('stroke', '#484f58').attr('stroke-width', 1.5);

  const node = svg.append('g').selectAll('circle').data(graphNodes).join('circle')
    .attr('r', 8).attr('fill', d => d.color)
    .call(d3.drag()
      .on('start', (event, d) => { if (!event.active) simulation.alphaTarget(0.3).restart(); d.fx = d.x; d.fy = d.y; })
      .on('drag', (event, d) => { d.fx = event.x; d.fy = event.y; })
      .on('end', (event, d) => { if (!event.active) simulation.alphaTarget(0); d.fx = null; d.fy = null; }));
  node.append('title').text(d => d.label);

  const label = svg.append('g').selectAll('text').data(graphNodes).join('text')
    .text(d => d.label).attr('font-size', 10).attr('fill', '#8b949e').attr('dx', 12).attr('dy', 4);

  simulation.on('tick', () => {
    link.attr('x1', d => d.source.x).attr('y1', d => d.source.y).attr('x2', d => d.target.x).attr('y2', d => d.target.y);
    node.attr('cx', d => d.x).attr('cy', d => d.y);
    label.attr('x', d => d.x).attr('y', d => d.y);
  });
}
</script>
</body>
</html>
`
