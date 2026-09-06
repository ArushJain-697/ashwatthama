import { useMemo, useRef, useState } from 'react'
import type { Verdict } from '../types'
import { buildGraph, type GraphNode, type Zone } from '../graph'
import { BlastRadiusScene } from './graph/BlastRadiusScene'
import { Inspector } from './Inspector'
import { Legend } from './Legend'

const ALL_ZONES: Zone[] = ['core', 'blast', 'uncertain', 'fog', 'outside']
const IDLE_RESUME_MS = 2500

export function GraphView({ verdict }: { verdict: Verdict }) {
  const { nodes, edges } = useMemo(() => buildGraph(verdict), [verdict])
  const [visibleZones, setVisibleZones] = useState<Set<string>>(new Set(ALL_ZONES))
  const [autoRotate, setAutoRotate] = useState(true)
  const [autoRotateWanted, setAutoRotateWanted] = useState(true)
  const [selected, setSelected] = useState<GraphNode | null>(null)
  const [resetSignal, setResetSignal] = useState(0)
  const idleTimer = useRef<number | null>(null)

  function toggleZones(zones: Zone[], checked: boolean) {
    setVisibleZones((current) => {
      const next = new Set(current)
      for (const zone of zones) (checked ? next.add(zone) : next.delete(zone))
      return next
    })
  }

  function handleInteractStart() {
    setAutoRotate(false)
    if (idleTimer.current) window.clearTimeout(idleTimer.current)
  }

  function handleInteractEnd() {
    idleTimer.current = window.setTimeout(() => setAutoRotate(autoRotateWanted), IDLE_RESUME_MS)
  }

  function toggleSpin() {
    const next = !autoRotateWanted
    setAutoRotateWanted(next)
    setAutoRotate(next)
  }

  return (
    <section className="px-7 pb-16 relative z-[1]">
      <div className="relative w-full h-[78vh] border border-line rounded-[20px] bg-[#020305] overflow-hidden shadow-[0_0_60px_rgba(255,69,0,0.06)_inset]">
        <div className="absolute top-3 left-3 flex gap-1.5 z-20">
          <button
            className="bg-black/85 border border-line text-fg rounded-full px-3.5 py-1.5 text-[11px] backdrop-blur hover:border-accent transition-colors"
            onClick={toggleSpin}
          >
            {autoRotateWanted ? '⏸ pause spin' : '▶ resume spin'}
          </button>
          <button
            className="bg-black/85 border border-line text-fg rounded-full px-3.5 py-1.5 text-[11px] backdrop-blur hover:border-accent transition-colors"
            onClick={() => setResetSignal((n) => n + 1)}
          >
            ⟲ reset view
          </button>
        </div>
        <div className="absolute top-3 right-3 text-[11px] text-muted bg-black/85 border border-line rounded-full px-3.5 py-1.5 z-20 font-mono">
          {nodes.length} nodes · {edges.length} real edges
        </div>

        <BlastRadiusScene
          nodes={nodes}
          edges={edges}
          visibleZones={visibleZones}
          autoRotate={autoRotate}
          onInteractStart={handleInteractStart}
          onInteractEnd={handleInteractEnd}
          onSelect={setSelected}
          resetSignal={resetSignal}
        />

        <Inspector node={selected} verdict={verdict} onClose={() => setSelected(null)} />
      </div>

      <Legend visibleZones={visibleZones} onToggle={toggleZones} />

      <p className="text-xs text-muted mt-3 leading-relaxed max-w-3xl">
        Distance from center is a real graph fact ONLY inside the green blast radius (hop count from a confirmed
        change, connected by real traversed edges, shown as pulsing lines). The uncertain/fog/outside rings are
        fixed layout zones, not measured distances — advisory evidence sits in a ring, unverifiable-coverage
        entities are literally rendered inside the scene's particle fog (the graph could not see clearly there),
        and scope-creep sits outside it all: confirmed drift, not uncertainty. Drag to orbit, scroll to zoom, click
        a node to inspect it.
      </p>
    </section>
  )
}
