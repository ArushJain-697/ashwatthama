import type { Verdict } from '../types'
import { LABEL_CONFIRMED, LABEL_DRIFT, LABEL_HEURISTIC, LABEL_UNVERIFIABLE } from '../types'
import { TierCard } from './TierCard'
import { TierList } from './TierList'

export function SummaryTab({ verdict }: { verdict: Verdict }) {
  const r = verdict.reconciliation
  const confirmedCount = r.confirmed.length + r.expected_blast_radius.length
  const driftCount = r.undeclared_scope_creep.length + r.declared_unimplemented.length

  return (
    <section className="px-7 pb-16 bg-charcoal">
      <div className="grid grid-cols-4 gap-4 mb-7">
        <TierCard count={confirmedCount} label={LABEL_CONFIRMED} color="#3fb950" />
        <TierCard count={r.advisory_low_confidence.length} label={LABEL_HEURISTIC} color="#d29922" delayMs={80} />
        <TierCard count={r.unverifiable_coverage.length} label={LABEL_UNVERIFIABLE} color="#bc8cff" delayMs={160} />
        <TierCard count={driftCount} label={LABEL_DRIFT} color="#f85149" delayMs={240} />
      </div>

      <TierList
        title="Confirmed"
        entries={r.confirmed}
        format={(e) => (
          <>
            {e.entity} — {e.file ?? ''}:{e.line ?? ''} [coverage: {e.coverage_confidence}]
          </>
        )}
      />
      <TierList
        title="Expected Blast Radius"
        entries={r.expected_blast_radius}
        format={(e) => (
          <>
            {e.entity} — {e.file ?? ''}:{e.line ?? ''} — {e.hops} hop(s)
          </>
        )}
      />
      <TierList
        title="Advisory Low Confidence"
        entries={r.advisory_low_confidence}
        format={(e) => (
          <>
            {e.entity} — {e.reason}
          </>
        )}
      />
      <TierList
        title="Unverifiable Coverage"
        entries={r.unverifiable_coverage}
        format={(e) => (
          <>
            {e.entity} — {e.coverage_reason} [{e.verification_path}]
          </>
        )}
      />
      <TierList
        title="Undeclared Scope Creep"
        entries={r.undeclared_scope_creep}
        format={(e) => (
          <>
            {e.entity} — {e.reason}
            {e.community_label ? ` [community: ${e.community_label}]` : ''}
          </>
        )}
      />
      <TierList
        title="Declared Unimplemented"
        entries={r.declared_unimplemented}
        format={(e) => (
          <>
            {e.entity} — {e.reason}
          </>
        )}
      />
    </section>
  )
}
