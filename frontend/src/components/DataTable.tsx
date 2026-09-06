import { useMemo, useState, type ReactNode } from 'react'

export interface Column<T> {
  header: string
  width?: string // a grid-template-columns track, e.g. "1fr" or "120px"
  align?: 'left' | 'right'
  render: (row: T) => ReactNode
}

const PAGE_SIZE = 25

// Generic paginated data table -- mirrors Dwarpal's AttackFeed.jsx
// .thead/.trow/.tbody pattern, parameterized by column definitions instead
// of being rewritten per tier. Every tier page (confirmed/advisory/
// unverifiable/scope-creep) is a thin wrapper choosing columns + rows.
export function DataTable<T>({
  title,
  meta,
  columns,
  rows,
  emptyText,
}: {
  title: string
  meta?: ReactNode
  columns: Column<T>[]
  rows: T[]
  emptyText: string
}) {
  const [page, setPage] = useState(0)
  const pageCount = Math.max(1, Math.ceil(rows.length / PAGE_SIZE))
  const clampedPage = Math.min(page, pageCount - 1)
  const visible = useMemo(() => rows.slice(clampedPage * PAGE_SIZE, clampedPage * PAGE_SIZE + PAGE_SIZE), [rows, clampedPage])
  const gridTemplate = columns.map((c) => c.width ?? '1fr').join(' ')

  return (
    <div className="card">
      <div className="card-core">
        <div className="card-head">
          <h2>{title}</h2>
          <span className="meta">{meta ?? `${rows.length} ${rows.length === 1 ? 'entry' : 'entries'}`}</span>
        </div>

        {rows.length === 0 ? (
          <div className="empty">{emptyText}</div>
        ) : (
          <>
            <div className="thead" style={{ gridTemplateColumns: gridTemplate }}>
              {columns.map((c) => (
                <span key={c.header} style={{ textAlign: c.align ?? 'left' }}>
                  {c.header}
                </span>
              ))}
            </div>
            <div className="tbody">
              {visible.map((row, index) => (
                <div className="trow" key={clampedPage * PAGE_SIZE + index} style={{ gridTemplateColumns: gridTemplate }}>
                  {columns.map((c) => (
                    <span key={c.header} style={{ textAlign: c.align ?? 'left' }}>
                      {c.render(row)}
                    </span>
                  ))}
                </div>
              ))}
            </div>
            {pageCount > 1 && (
              <div className="pager">
                <button onClick={() => setPage((p) => Math.max(0, p - 1))} disabled={clampedPage === 0}>
                  ← prev
                </button>
                <span className="pg-n">
                  page {clampedPage + 1} of {pageCount}
                </span>
                <button onClick={() => setPage((p) => Math.min(pageCount - 1, p + 1))} disabled={clampedPage === pageCount - 1}>
                  next →
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
