import type { ReactNode } from 'react'
import { useReveal } from '../hooks/useReveal'

export function TierList<T>({ title, entries, format }: { title: string; entries: T[]; format: (entry: T) => ReactNode }) {
  const { ref, active } = useReveal<HTMLDivElement>()
  if (entries.length === 0) return null
  return (
    <div ref={ref} data-active={active} className="reveal mb-5.5">
      <h3 className="text-xs uppercase tracking-wider text-muted mb-2 font-medium">
        {title} ({entries.length})
      </h3>
      {entries.map((entry, index) => (
        <div key={index} className="text-[13px] py-1.5 border-b border-line font-mono">
          {format(entry)}
        </div>
      ))}
    </div>
  )
}
