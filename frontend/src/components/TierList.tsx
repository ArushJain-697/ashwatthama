import type { ReactNode } from 'react'
import { useReveal } from '../hooks/useReveal'

export function TierList<T>({ title, entries, format }: { title: string; entries: T[]; format: (entry: T) => ReactNode }) {
  const { ref, active } = useReveal<HTMLDivElement>()
  if (entries.length === 0) return null
  return (
    <div ref={ref} data-active={active} className="reveal bg-white border-2 border-black p-5 mb-5">
      <h3 className="text-xs font-sans font-medium uppercase tracking-wider text-black/60 mb-2.5">
        {title} ({entries.length})
      </h3>
      {entries.map((entry, index) => (
        <div key={index} className="text-[13px] py-1.5 border-b border-black/10 font-mono text-black">
          {format(entry)}
        </div>
      ))}
    </div>
  )
}
