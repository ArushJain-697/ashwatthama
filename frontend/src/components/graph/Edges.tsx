import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'
import type { GraphEdge, GraphNode } from '../../graph'
import { positionFor } from '../../graph'

// Real blast-radius edges only, drawn as pulsing dashed lines -- a moving
// dash pattern reads as "evidence flowing along a real traversed path,"
// which is exactly what an expected_blast_radius edge is. Same technique as
// the Go-rendered scene: LineDashedMaterial's dashOffset animated per frame.
function Edge({ from, to }: { from: [number, number, number]; to: [number, number, number] }) {
  const lineRef = useRef<THREE.Line>(null)
  const geometry = useRef(new THREE.BufferGeometry().setFromPoints([new THREE.Vector3(...from), new THREE.Vector3(...to)]))

  useFrame(({ clock }) => {
    const line = lineRef.current
    // dashOffset is a real, animatable LineDashedMaterial property (feeds
    // the built-in dashed-line shader) that this version of @types/three
    // doesn't declare -- cast rather than wait on upstream typings.
    if (line) (line.material as THREE.LineDashedMaterial & { dashOffset: number }).dashOffset = -clock.getElapsedTime() * 12
  })

  return (
    <primitive
      object={
        new THREE.Line(
          geometry.current,
          new THREE.LineDashedMaterial({ color: 0x56d364, dashSize: 6, gapSize: 3, transparent: true, opacity: 0.9 }),
        )
      }
      ref={lineRef}
      onUpdate={(line: THREE.Line) => line.computeLineDistances()}
    />
  )
}

export function Edges({ edges, nodeById }: { edges: GraphEdge[]; nodeById: Map<string, GraphNode> }) {
  return (
    <group>
      {edges.map((edge, index) => {
        const from = nodeById.get(edge.source)
        const to = nodeById.get(edge.target)
        if (!from || !to) return null
        return <Edge key={index} from={positionFor(from)} to={positionFor(to)} />
      })}
    </group>
  )
}
