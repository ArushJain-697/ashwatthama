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
    <div
      className={`absolute top-0 right-0 h-full w-[300px] bg-white border-l-2 border-black p-5 box-border overflow-y-auto z-30 transition-transform duration-300 ease-[cubic-bezier(.22,1,.36,1)] ${
        open ? 'translate-x-0' : 'translate-x-full'
      }`}
    >
      <button
        className="absolute top-3.5 right-3.5 text-black text-lg font-sans"
        onClick={onClose}
        aria-label="Close inspector"
      >
        ✕
      </button>
      {node && (
        <>
          <h4 className="font-serif font-extrabold text-base mb-2 break-words text-black">{node.label}</h4>
          <span
            className="inline-block text-[10px] uppercase tracking-wide px-2.5 py-0.5 border-2 border-black font-sans font-medium mb-3.5"
            style={{ background: tierColor }}
          >
            {node.tier}
          </span>
          <dl className="text-xs font-sans text-black">
            {node.zone === 'blast' && (
              <>
                <dt className="text-black/50 mt-3">Hops from confirmed</dt>
                <dd className="mt-0.5 font-mono break-words">{node.hops}</dd>
              </>
            )}
            {entry &&
              Object.entries(entry)
                .filter(([, value]) => value !== null && value !== undefined && value !== '')
                .map(([key, value]) => (
                  <div key={key}>
                    <dt className="text-black/50 mt-3">{key}</dt>
                    <dd className="mt-0.5 font-mono break-words">
                      {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                    </dd>
                  </div>
                ))}
          </dl>
        </>
      )}
    </div>
  )
}
