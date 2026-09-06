import type { GraphNode, Zone } from '../graph'
import { ZONE_COLOR } from '../graph'
import type { Verdict } from '../types'

// Looks up the full reconciliation entry behind a clicked graph node, by
// entity name -- the graph itself only ever carries id/label/tier/zone/hops,
// everything else stays in one place, the fetched verdict.
function findEntry(verdict: Verdict, label: string): Record<string, unknown> | null {
  const tiers = [
    'confirmed',
    'expected_blast_radius',
    'advisory_low_confidence',
    'unverifiable_coverage',
    'undeclared_scope_creep',
  ] as const
  for (const tier of tiers) {
    const entries = verdict.reconciliation[tier] as unknown as Array<Record<string, unknown>>
    const hit = entries.find(
      (entry) => entry.entity === label || (typeof entry.entity === 'string' && entry.entity.endsWith('.' + label)),
    )
    if (hit) return hit
  }
  return null
}

export function Inspector({ node, verdict, onClose }: { node: GraphNode | null; verdict: Verdict; onClose: () => void }) {
  const open = node !== null
  const entry = node ? findEntry(verdict, node.label) : null
  const tierColor = node ? ZONE_COLOR[node.zone as Zone] : '#ffffff'

  return (
    <div className={`gv-inspector ${open ? 'open' : ''}`}>
      <button className="gv-inspector-close" onClick={onClose} aria-label="Close inspector">
        ✕
      </button>
      {node && (
        <>
          <h4>{node.label}</h4>
          <span className="badge" style={{ background: tierColor }}>
            {node.tier}
          </span>
          <dl>
            {node.zone === 'blast' && (
              <>
                <dt>Hops from confirmed</dt>
                <dd>{node.hops}</dd>
              </>
            )}
            {entry &&
              Object.entries(entry)
                .filter(([, value]) => value !== null && value !== undefined && value !== '')
                .map(([key, value]) => (
                  <div key={key}>
                    <dt>{key}</dt>
                    <dd>{typeof value === 'object' ? JSON.stringify(value) : String(value)}</dd>
                  </div>
                ))}
          </dl>
        </>
      )}
    </div>
  )
}
