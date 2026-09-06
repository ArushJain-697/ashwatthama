import { navigate } from '../App'
import type { Page } from '../App'

// The console's left rail. Mirrors Dwarpal's Sidebar.jsx structure: brand
// mark, a nav list with an active filled-charcoal item, a footer status pill.

const NAV: { id: Page; label: string; count?: (c: Counts) => number }[] = [
  { id: 'overview', label: 'Overview' },
  { id: 'confirmed', label: 'Confirmed', count: (c) => c.confirmed },
  { id: 'advisory', label: 'Advisory', count: (c) => c.advisory },
  { id: 'unverifiable', label: 'Unverifiable', count: (c) => c.unverifiable },
  { id: 'scope-creep', label: 'Scope Creep', count: (c) => c.scopeCreep },
  { id: 'graph', label: 'Graph View' },
]

export interface Counts {
  confirmed: number
  advisory: number
  unverifiable: number
  scopeCreep: number
}

export function Sidebar({
  page,
  onNavigate,
  counts,
  verdictLabel,
}: {
  page: Page
  onNavigate: (page: Page) => void
  counts: Counts
  verdictLabel: string
}) {
  const isReview = verdictLabel === 'REVIEW_REQUIRED'
  return (
    <aside className="sidebar">
      <div className="sb-brand">
        <span className="sb-mark">F</span>
        <span className="sb-word">
          Fidelity
          <span className="sb-sub">Verdict</span>
        </span>
      </div>

      <nav className="sb-nav">
        {NAV.map((item) => (
          <button key={item.id} className={`sb-link ${page === item.id ? 'on' : ''}`} onClick={() => onNavigate(item.id)}>
            <span>{item.label}</span>
            {item.count && <span className="sb-count">{item.count(counts)}</span>}
          </button>
        ))}
      </nav>

      <div className="sb-foot">
        <span className={`live ${isReview ? 'off' : 'on'}`}>
          <span className="live-dot" />
          {verdictLabel}
        </span>
        <button className="sb-back" onClick={() => navigate('/')}>
          ← Back to site
        </button>
      </div>
    </aside>
  )
}
