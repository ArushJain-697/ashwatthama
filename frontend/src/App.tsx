import { useState } from 'react'
import { useVerdict } from './hooks/useVerdict'
import { Header } from './components/Header'
import { TabNav, type Tab } from './components/TabNav'
import { SummaryTab } from './components/SummaryTab'
import { GraphView } from './components/GraphView'

export default function App() {
  const state = useVerdict()
  const [tab, setTab] = useState<Tab>('summary')

  return (
    <div className="min-h-screen selection:bg-accent selection:text-white">
      <div className="noise-overlay" />

      {state.status === 'loading' && <div className="p-7 text-muted">Loading verdict.json…</div>}

      {state.status === 'error' && (
        <div className="p-7 text-[#f85149] max-w-2xl">
          <h2 className="font-serif text-xl mb-2">Couldn't load verdict.json</h2>
          <p className="text-sm text-muted">{state.message}</p>
          <p className="text-sm text-muted mt-3">
            In dev, this reads <code className="font-mono">public/verdict.json</code>. In a deployment, drop a real
            verdict.json (from <code className="font-mono">verify-intent --out</code>) next to the built{' '}
            <code className="font-mono">dist/</code> bundle before serving it.
          </p>
        </div>
      )}

      {state.status === 'ready' && (
        <>
          <Header checkpointId={state.verdict.checkpoint_id} verdictLabel={state.verdict.summary.verdict_label} />
          <TabNav tab={tab} onChange={setTab} />
          {tab === 'summary' ? <SummaryTab verdict={state.verdict} /> : <GraphView verdict={state.verdict} />}
        </>
      )}
    </div>
  )
}
