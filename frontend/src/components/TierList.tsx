import { useMemo, useState, type ReactNode } from 'react'
import { useReveal } from '../hooks/useReveal'

const PAGE_SIZE = 25

export function TierList<T>({ title, entries, format }: { title: string; entries: T[]; format: (entry: T) => ReactNode }) {
  const { ref, active } = useReveal<HTMLDivElement>()
  const [page, setPage] = useState(0)

  const pageCount = Math.max(1, Math.ceil(entries.length / PAGE_SIZE))
  const clampedPage = Math.min(page, pageCount - 1)
  const visible = useMemo(
    () => entries.slice(clampedPage * PAGE_SIZE, clampedPage * PAGE_SIZE + PAGE_SIZE),
    [entries, clampedPage],
  )

  if (entries.length === 0) return null

  return (
    <div ref={ref} data-active={active} className="reveal bg-white border-2 border-black p-5 mb-5">
      <div className="flex items-center justify-between mb-2.5 flex-wrap gap-2">
        <h3 className="text-xs font-sans font-medium uppercase tracking-wider text-black/60">
          {title} ({entries.length})
        </h3>
        {pageCount > 1 && (
          <div className="flex items-center gap-2">
            <button
              className="brutal-btn bg-white border-2 border-black shadow-hard-sm px-2.5 py-0.5 text-xs font-sans font-medium disabled:opacity-30 disabled:pointer-events-none"
              onClick={() => setPage((p) => Math.max(0, p - 1))}
              disabled={clampedPage === 0}
              aria-label="Previous page"
            >
              ← prev
            </button>
            <span className="text-xs font-sans font-medium text-black/60 whitespace-nowrap">
              page {clampedPage + 1} of {pageCount}
            </span>
            <button
              className="brutal-btn bg-white border-2 border-black shadow-hard-sm px-2.5 py-0.5 text-xs font-sans font-medium disabled:opacity-30 disabled:pointer-events-none"
              onClick={() => setPage((p) => Math.min(pageCount - 1, p + 1))}
              disabled={clampedPage === pageCount - 1}
              aria-label="Next page"
            >
              next →
            </button>
          </div>
        )}
      </div>
      {visible.map((entry, index) => (
        <div key={clampedPage * PAGE_SIZE + index} className="text-[13px] py-1.5 border-b border-black/10 font-mono text-black">
          {format(entry)}
        </div>
      ))}
    </div>
  )
}
