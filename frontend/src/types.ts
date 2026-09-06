// Mirrors internal/fidelity/verdict.go's Verdict struct and reconcile.go's
// per-tier entity types, field-for-field, by their JSON tags. Keep this in
// sync with the Go source by hand -- there is no generator wiring this up.

export type CoverageConfidence = 'full' | 'partial' | 'unknown'

export interface ExtractedEntity {
  name: string
  method: string
  confidence: number
}

export interface VerdictIntent {
  raw_prompt_excerpt: string
  extracted_entities: ExtractedEntity[]
  unresolved_fragments: string[]
}

export interface EntityLocation {
  entity: string
  file?: string
  line?: number
  coverage_confidence: CoverageConfidence
}

export interface UnimplementedEntity {
  entity: string
  reason: string
}

export interface BlastRadiusEntity {
  entity: string
  file?: string
  line?: number
  path_from_confirmed: string[]
  hops: number
  edge_types_traversed: string[]
  min_edge_class: string
  coverage_confidence: CoverageConfidence
}

export interface ScopeCreepEntity {
  entity: string
  file?: string
  line?: number
  reason: string
  coverage_confidence: CoverageConfidence
  community_label?: string
}

export interface Corroboration {
  attempted: boolean
  verified: boolean | null
  evidence_snippet: string | null
}

export interface AdvisoryEntity {
  entity: string
  file?: string
  line?: number
  edge_class: string
  reason: string
  coverage_confidence: CoverageConfidence
  verbal_confidence: number | null
  conformal_gate_passed: boolean | null
  calibrated_confidence?: number
  corroboration: Corroboration
}

export interface UnverifiableCoverageEntity {
  entity: string
  file?: string
  line?: number
  coverage_confidence: CoverageConfidence
  coverage_reason: string
  verification_path: string
  corroboration: Corroboration
}

export interface Reconciliation {
  confirmed: EntityLocation[]
  declared_unimplemented: UnimplementedEntity[]
  expected_blast_radius: BlastRadiusEntity[]
  undeclared_scope_creep: ScopeCreepEntity[]
  advisory_low_confidence: AdvisoryEntity[]
  unverifiable_coverage: UnverifiableCoverageEntity[]
}

export interface VerdictSummary {
  total_changed_entities: number
  scope_creep_count: number
  unimplemented_count: number
  advisory_edge_count: number
  unverifiable_coverage_count: number
  verdict_label: string
}

export interface Verdict {
  checkpoint_id: string
  base_checkpoint_id?: string
  generated_at: string
  graph_capabilities_verified: string[]
  intent: VerdictIntent
  reconciliation: Reconciliation
  summary: VerdictSummary
}

// The four visual categories the Curveball card requires be distinguishable
// at a glance -- confirmed structural evidence, heuristic/incomplete
// evidence, claims needing verification, and drift from stated intent.
export const LABEL_CONFIRMED = 'CONFIRMED STRUCTURAL EVIDENCE'
export const LABEL_HEURISTIC = 'HEURISTIC / INCOMPLETE EVIDENCE'
export const LABEL_UNVERIFIABLE = 'NEEDS SOURCE/TEST VERIFICATION'
export const LABEL_DRIFT = 'DRIFT FROM STATED INTENT'
