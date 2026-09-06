import { useRef } from 'react'
import { useFrame, type ThreeEvent } from '@react-three/fiber'
import * as THREE from 'three'
import type { GraphNode } from '../../graph'
import { positionFor } from '../../graph'

const sphereGeometry = new THREE.SphereGeometry(6, 24, 24)
const wireGeometry = new THREE.IcosahedronGeometry(7, 0)

function Node({ node, onSelect }: { node: GraphNode; onSelect: (node: GraphNode) => void }) {
  const meshRef = useRef<THREE.Mesh>(null)
  const position = positionFor(node)
  const isFog = node.zone === 'fog'

  // Pop-in growth on mount, same as the Go-rendered scene: converges to 1
  // and then costs nothing to leave running.
  useFrame(() => {
    const mesh = meshRef.current
    if (mesh && mesh.scale.x < 1) mesh.scale.setScalar(Math.min(1, mesh.scale.x + 0.06))
  })

  function handleClick(event: ThreeEvent<MouseEvent>) {
    event.stopPropagation()
    onSelect(node)
  }

  return (
    <mesh
      ref={meshRef}
      position={position}
      scale={0.001}
      geometry={isFog ? wireGeometry : sphereGeometry}
      onClick={handleClick}
    >
      <meshStandardMaterial
        color={node.color}
        emissive={node.color}
        emissiveIntensity={isFog ? 0.6 : 1.1}
        roughness={0.35}
        metalness={0.15}
        transparent={isFog}
        opacity={isFog ? 0.65 : 1}
        wireframe={isFog}
      />
    </mesh>
  )
}

export function Nodes({
  nodes,
  visibleZones,
  onSelect,
}: {
  nodes: GraphNode[]
  visibleZones: Set<string>
  onSelect: (node: GraphNode) => void
}) {
  return (
    <group>
      {nodes
        .filter((node) => visibleZones.has(node.zone))
        .map((node) => (
          <Node key={node.id} node={node} onSelect={onSelect} />
        ))}
    </group>
  )
}
