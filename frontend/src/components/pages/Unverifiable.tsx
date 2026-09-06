import type { Verdict } from '../../types'
import { DataTable, type Column } from '../DataTable'

export function Unverifiable({ verdict }: { verdict: Verdict }) {
  const rows = verdict.reconciliation.unverifiable_coverage

  const columns: Column<(typeof rows)[number]>[] = [
    { header: 'Entity', width: '1fr', render: (e) => e.entity },
    { header: 'Location', width: '220px', render: (e) => <span className="path">{e.file ? `${e.file}:${e.line ?? ''}` : '—'}</span> },
    { header: 'Coverage', width: '90px', align: 'right', render: (e) => e.coverage_confidence },
    { header: 'Why', width: '2fr', render: (e) => <span className="muted">{e.coverage_reason}</span> },
    { header: 'Verify via', width: '150px', render: (e) => <span className="muted">{e.verification_path}</span> },
  ]

  return (
    <DataTable
      title="Unverifiable Coverage"
      columns={columns}
      rows={rows}
      emptyText="No unverifiable-coverage entities — the graph could see every change clearly."
    />
  )
}
