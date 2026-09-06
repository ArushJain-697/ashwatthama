import { Suspense, useMemo } from 'react'
import { Canvas } from '@react-three/fiber'
import { OrbitControls, Sparkles, Stars } from '@react-three/drei'
import { EffectComposer, Bloom } from '@react-three/postprocessing'
import type { GraphEdge, GraphNode } from '../../graph'
import { RING_FOG, RING_UNCERTAIN } from '../../graph'
import { Nodes } from './Nodes'
import { Edges } from './Edges'
import { HopRings } from './HopRings'

export function BlastRadiusScene({
  nodes,
  edges,
  visibleZones,
  autoRotate,
  onInteractStart,
  onInteractEnd,
  onSelect,
  resetSignal,
}: {
  nodes: GraphNode[]
  edges: GraphEdge[]
  visibleZones: Set<string>
  autoRotate: boolean
  onInteractStart: () => void
  onInteractEnd: () => void
  onSelect: (node: GraphNode) => void
  resetSignal: number
}) {
  const nodeById = useMemo(() => new Map(nodes.map((node) => [node.id, node])), [nodes])
  const maxHops = useMemo(
    () => Math.max(1, ...nodes.filter((node) => node.zone === 'blast').map((node) => node.hops || 1)),
    [nodes],
  )

  // Reset-view toolbar button: re-keying the Canvas via resetSignal is the
  // simplest way to force the camera and OrbitControls back to their
  // initial pose -- a fresh mount, not a manual reset of internal state.
  const cameraKey = resetSignal

  return (
    <Canvas
      key={cameraKey}
      camera={{ position: [0, 160, 520], fov: 55, near: 1, far: 3000 }}
      gl={{ antialias: true }}
      dpr={[1, 2]}
    >
      <color attach="background" args={['#020305']} />
      <fog attach="fog" args={['#020305', RING_UNCERTAIN, RING_FOG + 140]} />
      <ambientLight intensity={0.5} />
      <pointLight position={[250, 300, 300]} intensity={1.2} />
      <pointLight position={[-300, -200, -200]} intensity={0.6} color="#58a6ff" />

      <Suspense fallback={null}>
        <Stars radius={900} depth={200} count={2500} factor={2} saturation={0} fade speed={0.2} />
        {/* ponytail: drei's Sparkles fills a uniform box volume rather than
            the fog-radius SHELL the Go-rendered version scatters points in
            (a hand-rolled BufferGeometry there); reusing the library
            primitive costs some precision in where particles land, gained
            back in code not written. The scale is sized to roughly bracket
            the fog ring; upgrade to a custom shell sampler if the box-fill
            look reads wrong at demo time. */}
        <Sparkles
          count={700}
          size={14}
          scale={[(RING_FOG + 90) * 2, (RING_FOG + 90) * 2, (RING_FOG + 90) * 2]}
          position={[0, 0, 0]}
          color="#bc8cff"
          opacity={0.35}
          speed={0.1}
        />
      </Suspense>

      <HopRings maxHops={maxHops} />
      <Nodes nodes={nodes} visibleZones={visibleZones} onSelect={onSelect} />
      <Edges edges={edges} nodeById={nodeById} />

      <OrbitControls
        enableDamping
        dampingFactor={0.08}
        minDistance={60}
        maxDistance={1400}
        autoRotate={autoRotate}
        autoRotateSpeed={0.6}
        onStart={onInteractStart}
        onEnd={onInteractEnd}
      />

      <EffectComposer>
        <Bloom intensity={1.1} luminanceThreshold={0.15} luminanceSmoothing={0.5} />
      </EffectComposer>
    </Canvas>
  )
}
