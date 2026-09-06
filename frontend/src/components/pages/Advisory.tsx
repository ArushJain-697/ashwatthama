import type { Verdict } from '../../types'
import { DataTable, type Column } from '../DataTable'

export function Advisory({ verdict }: { verdict: Verdict }) {
  const rows = verdict.reconciliation.advisory_low_confidence

  const columns: Column<(typeof rows)[number]>[] = [
    { header: 'Entity', width: '1fr', render: (e) => e.entity },
    { header: 'Reason', width: '2fr', render: (e) => <span className="muted">{e.reason}</span> },
    { header: 'Verbal conf.', width: '110px', align: 'right', render: (e) => e.verbal_confidence ?? '—' },
    { header: 'Gate', width: '90px', align: 'right', render: (e) => (e.conformal_gate_passed === null ? '—' : e.conformal_gate_passed ? 'pass' : 'abstain') },
    { header: 'Corroborated', width: '110px', align: 'right', render: (e) => (e.corroboration.verified === null ? '—' : e.corroboration.verified ? 'yes' : 'no') },
  ]

  return (
    <DataTable
      title="Advisory Low Confidence"
      columns={columns}
      rows={rows}
      emptyText="No advisory entities — nothing here rests on an inferred/ambiguous edge alone."
    />
  )
}
