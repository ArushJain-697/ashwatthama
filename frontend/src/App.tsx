import { useEffect, useState } from 'react'
import { useVerdict } from './hooks/useVerdict'
import { Sidebar } from './components/Sidebar'
import { Overview } from './components/pages/Overview'
import { Confirmed } from './components/pages/Confirmed'
import { Advisory } from './components/pages/Advisory'
import { Unverifiable } from './components/pages/Unverifiable'
import { ScopeCreep } from './components/pages/ScopeCreep'
import { GraphView } from './components/GraphView'
import { Landing } from './landing/Landing'
import type { Verdict } from './types'

// Console routes. Mirrors Dwarpal's Dashboard.jsx PAGES map: an id, a page
// title, a one-line subtitle. State-only routing inside the console (no URL) —
// the app is a single-file report bundle, deep links aren't a requirement.
export type Page = 'overview' | 'confirmed' | 'advisory' | 'unverifiable' | 'scope-creep' | 'graph'

// Hash-based top-level routing (Dwarpal uses pathname + an nginx rewrite; a
// hash works from a plain static bundle with no server config and survives
// refresh). Anything under #/dashboard is the console; everything else is the
// marketing landing.
export function navigate(to: string) {
  window.location.hash = to === '/' ? '' : `#${to}`
  window.scrollTo(0, 0)
}

function useHash() {
  const [hash, setHash] = useState(window.location.hash)
  useEffect(() => {
    const onHash = () => setHash(window.location.hash)
    window.addEventListener('hashchange', onHash)
    return () => window.removeEventListener('hashchange', onHash)
  }, [])
  return hash
}

const META: Record<Page, { title: string; sub: string }> = {
  overview: { title: 'Overview', sub: 'Intent vs. diff, reconciled by the code graph' },
  confirmed: { title: 'Confirmed', sub: 'Named in intent, matched by the diff — with its real blast radius' },
  advisory: { title: 'Advisory', sub: 'Low-confidence matches resting on an inferred edge alone' },
  unverifiable: { title: 'Unverifiable', sub: 'Changes the graph could not see clearly — need source or tests' },
  'scope-creep': { title: 'Scope Creep', sub: 'Drift from stated intent — undeclared changes and unimplemented claims' },
  graph: { title: 'Graph View', sub: 'The blast radius in 3D — distance is a real hop count inside the core' },
}

function countsOf(verdict: Verdict) {
  const r = verdict.reconciliation
  return {
    confirmed: r.confirmed.length + r.expected_blast_radius.length,
    advisory: r.advisory_low_confidence.length,
    unverifiable: r.unverifiable_coverage.length,
    scopeCreep: r.undeclared_scope_creep.length + r.declared_unimplemented.length,
  }
}

function Console() {
  const state = useVerdict()
  const [page, setPage] = useState<Page>('overview')

  if (state.status === 'loading') {
    return (
      <div className="shell">
        <div className="aurora" />
        <main className="page">Loading verdict.json…</main>
      </div>
    )
  }

  if (state.status === 'error') {
    return (
      <div className="shell">
        <div className="aurora" />
        <main className="page">
          <div className="card">
            <div className="card-core">
              <div className="card-head">
                <h2>Couldn't load verdict.json</h2>
              </div>
              <div className="empty">
                {state.message}
                <p style={{ marginTop: 12 }}>
                  In dev this reads <code>public/verdict.json</code>. In a deployment, drop a real verdict.json
                  (from <code>verify-intent --out</code>) next to the built <code>dist/</code> bundle before serving it.
                </p>
              </div>
            </div>
          </div>
        </main>
      </div>
    )
  }

  const verdict = state.verdict
  const meta = META[page]

  return (
    <div className="shell">
      <div className="aurora" />
      <Sidebar
        page={page}
        onNavigate={setPage}
        counts={countsOf(verdict)}
        verdictLabel={verdict.summary.verdict_label}
      />

      <div className="content">
        <header className="pagebar">
          <div>
            <h1>{meta.title}</h1>
            <div className="sub">{meta.sub}</div>
          </div>
        </header>
        <main className="page" key={page}>
          {page === 'overview' && (
            <Overview verdict={verdict} onNavigate={(key) => setPage(key === 'drift' ? 'scope-creep' : key)} />
          )}
          {page === 'confirmed' && <Confirmed verdict={verdict} />}
          {page === 'advisory' && <Advisory verdict={verdict} />}
          {page === 'unverifiable' && <Unverifiable verdict={verdict} />}
          {page === 'scope-creep' && <ScopeCreep verdict={verdict} />}
          {page === 'graph' && <GraphView verdict={verdict} />}
        </main>
      </div>
    </div>
  )
}

export default function App() {
  const hash = useHash()
  return hash.startsWith('#/dashboard') ? <Console /> : <Landing />
}
