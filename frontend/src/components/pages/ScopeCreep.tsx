import type { Verdict } from '../../types'
import { DataTable, type Column } from '../DataTable'
import { Badge } from '../Badge'

export function ScopeCreep({ verdict }: { verdict: Verdict }) {
  const r = verdict.reconciliation

  const creepColumns: Column<(typeof r.undeclared_scope_creep)[number]>[] = [
    { header: 'Entity', width: '1fr', render: (e) => e.entity },
    { header: 'Location', width: '220px', render: (e) => <span className="path">{e.file ? `${e.file}:${e.line ?? ''}` : '—'}</span> },
    { header: 'Reason', width: '2fr', render: (e) => <span className="muted">{e.reason}</span> },
    { header: 'Community', width: '140px', render: (e) => <span className="muted">{e.community_label ?? '—'}</span> },
    { header: 'Coverage', width: '90px', align: 'right', render: (e) => e.coverage_confidence },
  ]

  const unimplColumns: Column<(typeof r.declared_unimplemented)[number]>[] = [
    { header: 'Entity', width: '1fr', render: (e) => e.entity },
    { header: 'Reason', width: '2fr', render: (e) => <span className="muted">{e.reason}</span> },
  ]

  return (
    <div className="stack">
      <div className="enter">
        <div className="card-head" style={{ border: 'none', padding: '0 0 8px' }}>
          <h2>
            <Badge tone="drift">undeclared scope creep</Badge> Changed, never mentioned in intent
          </h2>
        </div>
        <DataTable
          title="Undeclared Scope Creep"
          columns={creepColumns}
          rows={r.undeclared_scope_creep}
          emptyText="No undeclared scope creep — every change traces back to stated intent."
        />
      </div>
      <div className="enter" style={{ animationDelay: '60ms' }}>
        <div className="card-head" style={{ border: 'none', padding: '0 0 8px' }}>
          <h2>
            <Badge tone="drift">declared unimplemented</Badge> Named in intent, absent from the diff
          </h2>
        </div>
        <DataTable
          title="Declared Unimplemented"
          columns={unimplColumns}
          rows={r.declared_unimplemented}
          emptyText="Nothing declared was left unimplemented."
        />
      </div>
    </div>
  )
}
