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
	// Zone drives the 3D radial layout: "core" (confirmed, center),
	// "blast" (expected_blast_radius, distance = real hop count),
	// "uncertain" (advisory, a fixed ring beyond the farthest real hop),
	// "fog" (unverifiable_coverage, further out still, inside the scene's
	// literal THREE.Fog — visually "the graph can't see clearly here"), and
	// "outside" (undeclared_scope_creep, the farthest ring — confirmed
	// drift, not uncertainty, so it stays sharp instead of foggy).
	Zone string `json:"zone"`
	// Hops is the real graph-reachability distance for "blast" nodes; 0 for
	// every other zone, where distance is a fixed layout constant, not a
	// graph fact — never conflate the two.
	Hops int `json:"hops,omitempty"`
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
	addNode := func(id, label, tier, color, zone string, hops int) {
		if id == "" || seen[id] {
			return
		}
		seen[id] = true
		nodes = append(nodes, dashboardNode{ID: id, Label: label, Tier: tier, Color: color, Zone: zone, Hops: hops})
	}

	for _, entry := range verdict.Reconciliation.Confirmed {
		addNode(entry.Entity, entry.Entity, "confirmed", colorConfirmed, "core", 0)
	}
	for _, entry := range verdict.Reconciliation.ExpectedBlastRadius {
		// The path's own last element IS this entity's graph symbol ID
		// (reachability.go appends the target itself), so labeling it with
		// entry.Entity (the diff's short name) instead of the same
		// lastSegment used for every other path node would silently create
		// two disconnected nodes for one entity. Use the path consistently.
		if len(entry.PathFromConfirmed) == 0 {
			addNode(entry.Entity, entry.Entity, "expected_blast_radius", colorConfirmed, "blast", entry.Hops)
			continue
		}
		for i := 0; i < len(entry.PathFromConfirmed)-1; i++ {
			from, to := lastSegment(entry.PathFromConfirmed[i]), lastSegment(entry.PathFromConfirmed[i+1])
			fromTier, toTier := "confirmed", "confirmed"
			fromZone, toZone := "core", "core"
			fromHops, toHops := 0, 0
			if i+1 == len(entry.PathFromConfirmed)-1 {
				toTier, toZone, toHops = "expected_blast_radius", "blast", entry.Hops
			}
			if i > 0 {
				fromTier, fromZone, fromHops = "expected_blast_radius", "blast", i
			}
			addNode(from, from, fromTier, colorConfirmed, fromZone, fromHops)
			addNode(to, to, toTier, colorConfirmed, toZone, toHops)
			edgeType := ""
			if i < len(entry.EdgeTypesTraversed) {
				edgeType = entry.EdgeTypesTraversed[i]
			}
			edges = append(edges, dashboardEdge{Source: from, Target: to, Type: edgeType})
		}
	}
	for _, entry := range verdict.Reconciliation.AdvisoryLowConfidence {
		addNode(entry.Entity, entry.Entity, "advisory_low_confidence", colorHeuristic, "uncertain", 0)
	}
	for _, entry := range verdict.Reconciliation.UnverifiableCoverage {
		addNode(entry.Entity, entry.Entity, "unverifiable_coverage", colorUnverifiable, "fog", 0)
	}
	for _, entry := range verdict.Reconciliation.UndeclaredScopeCreep {
		addNode(entry.Entity, entry.Entity, "undeclared_scope_creep", colorDrift, "outside", 0)
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
  #graph-canvas-wrap { position:relative; width:100%; height:70vh; border:1px solid #30363d; border-radius:8px; background:#050709; overflow:hidden; }
  #graph-canvas-wrap canvas { display:block; cursor:grab; }
  #graph-canvas-wrap canvas:active { cursor:grabbing; }
  #node-tooltip { position:absolute; pointer-events:none; background:#161b22; border:1px solid #30363d; border-radius:6px; padding:6px 10px; font-size:12px; display:none; max-width:280px; }
  .legend { display:flex; gap:16px; margin-top:10px; font-size:12px; color:#8b949e; flex-wrap:wrap; }
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
  <div id="graph-canvas-wrap"><div id="node-tooltip"></div></div>
  <div class="legend">
    <span><i class="dot" style="background:#3fb950"></i>Core / Blast Radius — {{.LabelConfirmed}} (distance = real hop count)</span>
    <span><i class="dot" style="background:#d29922"></i>Uncertain ring — {{.LabelHeuristic}}</span>
    <span><i class="dot" style="background:#bc8cff"></i>Fog zone — {{.LabelUnverifiable}}</span>
    <span><i class="dot" style="background:#f85149"></i>Outside the blast radius — {{.LabelDrift}}</span>
  </div>
  <p class="note">Distance from center is a real graph fact ONLY inside the green blast radius (hop count from a confirmed change, connected by real traversed edges). The uncertain/fog/outside rings are fixed layout zones, not measured distances — advisory evidence sits in a ring, unverifiable-coverage entities are literally rendered inside the scene's fog (the graph could not see clearly there), and scope-creep sits outside it all: confirmed drift, not uncertainty. Drag to rotate, scroll to zoom.</p>
</section>
<script id="verdict-data" type="application/json">{{.VerdictJSON}}</script>
<script id="graph-nodes-data" type="application/json">{{.GraphNodes}}</script>
<script id="graph-edges-data" type="application/json">{{.GraphEdges}}</script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r128/three.min.js"></script>
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

// Radial "blast radius" layout: distance from the origin is a real graph
// fact (hop count) inside the blast zone, and a fixed zone ring everywhere
// else -- this scene never implies a measured distance it doesn't have.
const RING_UNIT = 55;      // world units per hop, inside the real blast radius
const RING_UNCERTAIN = 260; // advisory: fixed ring beyond the farthest plausible hop
const RING_FOG = 340;       // unverifiable_coverage: inside the literal scene fog
const RING_OUTSIDE = 430;   // undeclared_scope_creep: farthest, sharp (not foggy)

function hashAngle(id) {
  let h = 0;
  for (let i = 0; i < id.length; i++) h = (h * 31 + id.charCodeAt(i)) >>> 0;
  return h;
}

function positionFor(node) {
  const h = hashAngle(node.id);
  const theta = (h % 1000) / 1000 * Math.PI * 2;
  const phi = ((h >> 10) % 1000) / 1000 * Math.PI;
  let radius;
  if (node.zone === 'core') radius = 12 + (h % 15);
  else if (node.zone === 'blast') radius = node.hops * RING_UNIT + (h % 20);
  else if (node.zone === 'uncertain') radius = RING_UNCERTAIN + (h % 40);
  else if (node.zone === 'fog') radius = RING_FOG + (h % 50);
  else radius = RING_OUTSIDE + (h % 60);
  return new THREE.Vector3(
    radius * Math.sin(phi) * Math.cos(theta),
    radius * Math.sin(phi) * Math.sin(theta),
    radius * Math.cos(phi)
  );
}

let graphDrawn = false;
function drawGraph() {
  if (graphDrawn) return;
  graphDrawn = true;
  const wrap = document.getElementById('graph-canvas-wrap');
  const tooltip = document.getElementById('node-tooltip');
  const width = wrap.clientWidth || 800, height = wrap.clientHeight || 500;

  const scene = new THREE.Scene();
  scene.background = new THREE.Color(0x050709);
  // The fog IS the visualization of coverage uncertainty: anything at or
  // beyond the fog ring fades into haze, exactly like an entity the graph
  // could not see clearly -- a real renderer feature standing in for a real
  // property of the data, not decoration.
  scene.fog = new THREE.Fog(0x050709, RING_UNCERTAIN, RING_FOG + 120);

  const camera = new THREE.PerspectiveCamera(55, width / height, 1, 2000);
  camera.position.set(0, 120, 480);

  const renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setSize(width, height);
  renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
  wrap.insertBefore(renderer.domElement, tooltip);

  scene.add(new THREE.AmbientLight(0xffffff, 0.65));
  const point = new THREE.PointLight(0xffffff, 0.8);
  point.position.set(200, 300, 300);
  scene.add(point);

  const group = new THREE.Group();
  scene.add(group);

  // Faint reference rings at each real hop distance, so "farther = more
  // hops" reads as a measurement, not just a vibe.
  const maxHops = Math.max(1, ...graphNodes.filter(n => n.zone === 'blast').map(n => n.hops || 1));
  for (let hop = 1; hop <= maxHops; hop++) {
    const ringGeo = new THREE.RingGeometry(hop * RING_UNIT - 0.5, hop * RING_UNIT + 0.5, 64);
    const ringMat = new THREE.MeshBasicMaterial({ color: 0x3fb950, transparent: true, opacity: 0.12, side: THREE.DoubleSide });
    const ring = new THREE.Mesh(ringGeo, ringMat);
    ring.rotation.x = Math.PI / 2;
    group.add(ring);
  }

  const positions = new Map();
  const meshByID = new Map();
  const sphereGeo = new THREE.SphereGeometry(6, 16, 16);
  graphNodes.forEach(n => {
    const pos = positionFor(n);
    positions.set(n.id, pos);
    const mat = new THREE.MeshStandardMaterial({
      color: n.color,
      transparent: n.zone === 'fog', opacity: n.zone === 'fog' ? 0.55 : 1,
      wireframe: n.zone === 'fog',
    });
    const mesh = new THREE.Mesh(sphereGeo, mat);
    mesh.position.copy(pos);
    mesh.userData = n;
    group.add(mesh);
    meshByID.set(n.id, mesh);
  });

  graphEdges.forEach(e => {
    const from = positions.get(e.source), to = positions.get(e.target);
    if (!from || !to) return;
    const geo = new THREE.BufferGeometry().setFromPoints([from, to]);
    const mat = new THREE.LineBasicMaterial({ color: 0x3fb950 });
    group.add(new THREE.Line(geo, mat));
  });

  // Hand-rolled orbit/zoom: drag to rotate the whole scene group, wheel to
  // dolly the camera. No extra script beyond three.js core itself.
  let dragging = false, lastX = 0, lastY = 0;
  renderer.domElement.addEventListener('mousedown', e => { dragging = true; lastX = e.clientX; lastY = e.clientY; });
  window.addEventListener('mouseup', () => { dragging = false; });
  renderer.domElement.addEventListener('mousemove', e => {
    if (dragging) {
      group.rotation.y += (e.clientX - lastX) * 0.005;
      group.rotation.x += (e.clientY - lastY) * 0.005;
      lastX = e.clientX; lastY = e.clientY;
    }
    const rect = renderer.domElement.getBoundingClientRect();
    pointerNDC.x = ((e.clientX - rect.left) / rect.width) * 2 - 1;
    pointerNDC.y = -((e.clientY - rect.top) / rect.height) * 2 + 1;
    tooltip.style.left = (e.clientX - rect.left + 12) + 'px';
    tooltip.style.top = (e.clientY - rect.top + 12) + 'px';
  });
  renderer.domElement.addEventListener('wheel', e => {
    e.preventDefault();
    camera.position.z = Math.max(80, Math.min(1200, camera.position.z + e.deltaY * 0.4));
  }, { passive: false });

  const pointerNDC = new THREE.Vector2(-10, -10);
  const raycaster = new THREE.Raycaster();
  function animate() {
    requestAnimationFrame(animate);
    if (!dragging) group.rotation.y += 0.0015; // slow idle spin so the shape reads even before a viewer touches it
    raycaster.setFromCamera(pointerNDC, camera);
    const hits = raycaster.intersectObjects(Array.from(meshByID.values()));
    if (hits.length > 0) {
      const n = hits[0].object.userData;
      tooltip.style.display = 'block';
      tooltip.textContent = n.label + ' — ' + n.tier + (n.zone === 'blast' ? ' (' + n.hops + ' hop' + (n.hops === 1 ? '' : 's') + ')' : '');
    } else {
      tooltip.style.display = 'none';
    }
    renderer.render(scene, camera);
  }
  animate();

  window.addEventListener('resize', () => {
    const w = wrap.clientWidth, h = wrap.clientHeight;
    camera.aspect = w / h; camera.updateProjectionMatrix();
    renderer.setSize(w, h);
  });
}
</script>
</body>
</html>
`
