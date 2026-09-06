import type { ReactNode } from 'react'
import { useReveal } from '../hooks/useReveal'

export function TierCard({
  count,
  label,
  color,
  delayMs = 0,
}: {
  count: number
  label: ReactNode
  color: string
  delayMs?: number
}) {
  const { ref, active } = useReveal<HTMLDivElement>()
  return (
    <div
      ref={ref}
      data-active={active}
      className="reveal bg-white border-2 border-black shadow-hard-sm p-5 transition-transform duration-200"
      style={{ transitionDelay: `${delayMs}ms` }}
    >
      <div className="w-4 h-4 border-2 border-black mb-3" style={{ background: color }} />
      <div className="font-serif font-extrabold text-4xl leading-none text-black">{count}</div>
      <div className="text-xs font-sans font-medium text-black/60 mt-2.5 uppercase tracking-wide">{label}</div>
    </div>
  )
}
