import type { Verdict } from '../../types'
import { DataTable, type Column } from '../DataTable'
import { Badge } from '../Badge'

export function Confirmed({ verdict }: { verdict: Verdict }) {
  const r = verdict.reconciliation

  const confirmedColumns: Column<(typeof r.confirmed)[number]>[] = [
    { header: 'Entity', width: '1fr', render: (e) => e.entity },
    { header: 'Location', width: '260px', render: (e) => <span className="path">{e.file ? `${e.file}:${e.line ?? ''}` : '—'}</span> },
    { header: 'Coverage', width: '120px', align: 'right', render: (e) => e.coverage_confidence },
  ]

  const blastColumns: Column<(typeof r.expected_blast_radius)[number]>[] = [
    { header: 'Entity', width: '1fr', render: (e) => e.entity },
    { header: 'Location', width: '220px', render: (e) => <span className="path">{e.file ? `${e.file}:${e.line ?? ''}` : '—'}</span> },
    { header: 'Hops', width: '70px', align: 'right', render: (e) => e.hops },
    { header: 'Via', width: '160px', render: (e) => <span className="muted">{e.edge_types_traversed.join(', ') || '—'}</span> },
    { header: 'Coverage', width: '110px', align: 'right', render: (e) => e.coverage_confidence },
  ]

  return (
    <div className="stack">
      <div className="enter">
        <div className="card-head" style={{ border: 'none', padding: '0 0 8px' }}>
          <h2>
            <Badge tone="confirmed">confirmed</Badge> Named in intent, matched by the diff
          </h2>
        </div>
        <DataTable title="Confirmed" columns={confirmedColumns} rows={r.confirmed} emptyText="No confirmed entities in this verdict." />
      </div>
      <div className="enter" style={{ animationDelay: '60ms' }}>
        <div className="card-head" style={{ border: 'none', padding: '0 0 8px' }}>
          <h2>
            <Badge tone="confirmed">expected blast radius</Badge> Reachable from a confirmed entity, real hop count
          </h2>
        </div>
        <DataTable
          title="Expected Blast Radius"
          columns={blastColumns}
          rows={r.expected_blast_radius}
          emptyText="No blast-radius entities in this verdict — nothing changed was reachable from a confirmed entity within the configured hop limit."
        />
      </div>
    </div>
  )
}
