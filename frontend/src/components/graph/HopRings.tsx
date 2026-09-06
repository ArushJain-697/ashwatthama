import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'
import { RING_UNIT } from '../../graph'

// Glowing torus rings at each real hop distance -- a torus stays visible
// from any camera angle, unlike a flat disc that vanishes edge-on.
export function HopRings({ maxHops }: { maxHops: number }) {
  const groupRef = useRef<THREE.Group>(null)
  useFrame(({ clock }) => {
    if (groupRef.current) groupRef.current.rotation.z = clock.getElapsedTime() * 0.03
  })

  const hops = Array.from({ length: maxHops }, (_, i) => i + 1)
  return (
    <group ref={groupRef} rotation-x={Math.PI / 2}>
      {hops.map((hop) => (
        <mesh key={hop}>
          <torusGeometry args={[hop * RING_UNIT, 0.6, 8, 96]} />
          <meshBasicMaterial color="#3fb950" transparent opacity={0.25} />
        </mesh>
      ))}
    </group>
  )
}
