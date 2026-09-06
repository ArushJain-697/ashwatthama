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
    <div className="min-h-screen bg-charcoal">
      {state.status === 'loading' && <div className="p-7 text-sage font-sans">Loading verdict.json…</div>}

      {state.status === 'error' && (
        <div className="p-7 max-w-2xl">
          <div className="bg-white border-2 border-black shadow-hard p-6">
            <h2 className="font-serif font-extrabold text-xl mb-2 text-black">Couldn't load verdict.json</h2>
            <p className="text-sm text-black/70 font-sans">{state.message}</p>
            <p className="text-sm text-black/70 font-sans mt-3">
              In dev, this reads <code className="font-mono bg-black/5 px-1">public/verdict.json</code>. In a
              deployment, drop a real verdict.json (from{' '}
              <code className="font-mono bg-black/5 px-1">verify-intent --out</code>) next to the built{' '}
              <code className="font-mono bg-black/5 px-1">dist/</code> bundle before serving it.
            </p>
          </div>
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
