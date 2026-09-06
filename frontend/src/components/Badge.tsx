import type { ReactNode } from 'react'

export type Tone = 'confirmed' | 'heuristic' | 'unverifiable' | 'drift'

export function Badge({ tone, children }: { tone: Tone; children: ReactNode }) {
  return <span className={`badge b-${tone}`}>{children}</span>
}
