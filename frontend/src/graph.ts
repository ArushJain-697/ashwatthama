import type { Verdict } from './types'

// Ports internal/fidelity/dashboard.go's dashboardGraph faithfully: real
// edges exist ONLY for expected_blast_radius (the only tier carrying a real
// path_from_confirmed); every other tier contributes an unconnected,
// tier-colored node. Keep this in sync with that function by hand.

export type Zone = 'core' | 'blast' | 'uncertain' | 'fog' | 'outside'

export interface GraphNode {
  id: string
  label: string
  tier: string
  color: string
  zone: Zone
  hops: number
}

export interface GraphEdge {
  source: string
  target: string
  type: string
}

export const ZONE_COLOR: Record<Zone, string> = {
  core: '#3fb950',
  blast: '#3fb950',
  uncertain: '#d29922',
  fog: '#bc8cff',
  outside: '#f85149',
}

// lastSegment reads a compound symbol ID the one positionally-safe way
// snapshot-format.md documents for a trailing field: anchor on the LAST
// separator. It is a display label, not an identity.
function lastSegment(id: string): string {
  const index = id.lastIndexOf(':')
  return index === -1 ? id : id.slice(index + 1)
}

export function buildGraph(verdict: Verdict): { nodes: GraphNode[]; edges: GraphEdge[] } {
  const nodes: GraphNode[] = []
  const edges: GraphEdge[] = []
  const seen = new Set<string>()

  function addNode(id: string, label: string, tier: string, zone: Zone, hops: number) {
    if (!id || seen.has(id)) return
    seen.add(id)
    nodes.push({ id, label, tier, color: ZONE_COLOR[zone], zone, hops })
  }

  for (const entry of verdict.reconciliation.confirmed) {
    addNode(entry.entity, entry.entity, 'confirmed', 'core', 0)
  }

  for (const entry of verdict.reconciliation.expected_blast_radius) {
    if (entry.path_from_confirmed.length === 0) {
      addNode(entry.entity, entry.entity, 'expected_blast_radius', 'blast', entry.hops)
      continue
    }
    for (let i = 0; i < entry.path_from_confirmed.length - 1; i++) {
      const from = lastSegment(entry.path_from_confirmed[i])
      const to = lastSegment(entry.path_from_confirmed[i + 1])
      const isLastHop = i + 1 === entry.path_from_confirmed.length - 1
      addNode(from, from, i > 0 ? 'expected_blast_radius' : 'confirmed', i > 0 ? 'blast' : 'core', i)
      addNode(to, to, isLastHop ? 'expected_blast_radius' : 'confirmed', isLastHop ? 'blast' : 'core', isLastHop ? entry.hops : 0)
      edges.push({ source: from, target: to, type: entry.edge_types_traversed[i] ?? '' })
    }
  }

  for (const entry of verdict.reconciliation.advisory_low_confidence) {
    addNode(entry.entity, entry.entity, 'advisory_low_confidence', 'uncertain', 0)
  }
  for (const entry of verdict.reconciliation.unverifiable_coverage) {
    addNode(entry.entity, entry.entity, 'unverifiable_coverage', 'fog', 0)
  }
  for (const entry of verdict.reconciliation.undeclared_scope_creep) {
    addNode(entry.entity, entry.entity, 'undeclared_scope_creep', 'outside', 0)
  }

  return { nodes, edges }
}

// Same radial "blast radius" layout as the Go-rendered version: distance
// from the origin is a real graph fact (hop count) inside the blast zone,
// and a fixed zone ring everywhere else.
export const RING_UNIT = 55
export const RING_UNCERTAIN = 260
export const RING_FOG = 340
export const RING_OUTSIDE = 430

function hashAngle(id: string): number {
  let h = 0
  for (let i = 0; i < id.length; i++) h = (Math.imul(h, 31) + id.charCodeAt(i)) >>> 0
  return h
}

export function positionFor(node: GraphNode): [number, number, number] {
  const h = hashAngle(node.id)
  const theta = ((h % 1000) / 1000) * Math.PI * 2
  const phi = (((h >>> 10) % 1000) / 1000) * Math.PI
  let radius: number
  if (node.zone === 'core') radius = 14 + (h % 15)
  else if (node.zone === 'blast') radius = node.hops * RING_UNIT + (h % 20)
  else if (node.zone === 'uncertain') radius = RING_UNCERTAIN + (h % 40)
  else if (node.zone === 'fog') radius = RING_FOG + (h % 50)
  else radius = RING_OUTSIDE + (h % 60)
  return [
    radius * Math.sin(phi) * Math.cos(theta),
    radius * Math.sin(phi) * Math.sin(theta),
    radius * Math.cos(phi),
  ]
}
