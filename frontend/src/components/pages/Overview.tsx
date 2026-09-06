import type { CSSProperties } from 'react'
import { LABEL_CONFIRMED, LABEL_DRIFT, LABEL_HEURISTIC, LABEL_UNVERIFIABLE, type Verdict } from '../../types'

const TILES = [
  { key: 'confirmed', label: LABEL_CONFIRMED, tone: 'var(--confirmed)' },
  { key: 'advisory', label: LABEL_HEURISTIC, tone: 'var(--heuristic)' },
  { key: 'unverifiable', label: LABEL_UNVERIFIABLE, tone: 'var(--unverifiable)' },
  { key: 'drift', label: LABEL_DRIFT, tone: 'var(--drift)' },
] as const

export function Overview({ verdict, onNavigate }: { verdict: Verdict; onNavigate: (key: (typeof TILES)[number]['key']) => void }) {
  const r = verdict.reconciliation
  const counts = {
    confirmed: r.confirmed.length + r.expected_blast_radius.length,
    advisory: r.advisory_low_confidence.length,
    unverifiable: r.unverifiable_coverage.length,
    drift: r.undeclared_scope_creep.length + r.declared_unimplemented.length,
  }
  const total = counts.confirmed + counts.advisory + counts.unverifiable + counts.drift

  return (
    <>
      <div className="enter stats">
        {TILES.map((tile) => (
          <div
            className="stat"
            key={tile.key}
            style={{ '--tone': tile.tone } as CSSProperties}
            onClick={() => onNavigate(tile.key)}
          >
            <div className="stat-core">
              <div className="stat-label">{tile.label}</div>
              <div className="stat-num" style={{ color: tile.tone }}>
                {counts[tile.key]}
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="enter" style={{ animationDelay: '80ms', marginTop: 20 }}>
        <div className="card">
          <div className="card-core">
            <div className="card-head">
              <h2>Tier mix</h2>
              <span className="meta">{verdict.summary.total_changed_entities} changed entities</span>
            </div>
            <div className="cats">
              {TILES.map((tile) => {
                const count = counts[tile.key]
                const pct = total > 0 ? (count / total) * 100 : 0
                return (
                  <div className="cat" key={tile.key}>
                    <div className="cat-top">
                      <span className="cat-name" style={{ color: tile.tone }}>
                        {tile.label}
                      </span>
                      <span className="cat-count">{count}</span>
                    </div>
                    <div className="cat-track">
                      <div className="cat-fill" style={{ width: `${pct}%`, background: tile.tone }} />
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        </div>
      </div>

      <div className="enter" style={{ animationDelay: '140ms', marginTop: 20 }}>
        <div className="card">
          <div className="card-core">
            <div className="card-head">
              <h2>Checkpoint</h2>
            </div>
            <div className="cats">
              <div className="cat">
                <div className="cat-top">
                  <span className="cat-name">checkpoint_id</span>
                  <span className="cat-count">{verdict.checkpoint_id}</span>
                </div>
              </div>
              {verdict.base_checkpoint_id && (
                <div className="cat">
                  <div className="cat-top">
                    <span className="cat-name">base_checkpoint_id</span>
                    <span className="cat-count">{verdict.base_checkpoint_id}</span>
                  </div>
                </div>
              )}
              <div className="cat">
                <div className="cat-top">
                  <span className="cat-name">generated_at</span>
                  <span className="cat-count">{verdict.generated_at}</span>
                </div>
              </div>
              <div className="cat">
                <div className="cat-top">
                  <span className="cat-name">graph_capabilities_verified</span>
                  <span className="cat-count">{verdict.graph_capabilities_verified.length} relation types</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
