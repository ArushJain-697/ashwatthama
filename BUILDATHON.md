# Fidelity Buildathon Record

## Phase 1 pre-flight evidence

Recorded locally on 2026-09-06. This is a factual setup record, not a claim
that the unreached build tickets are complete.

| Check | Evidence | Finding |
| --- | --- | --- |
| Entire mirror | `git remote -v` | `origin` uses `entire://aws-ap-south-1.entire.io/gh/arushjain-697/ashwatthama` for fetch and push. |
| Checkpoints and hooks | `entire status --json`, `entire doctor` | Entire is enabled and Codex hook approvals are present. At the time of this record, no agent session was active, so the first new session is still needed before making a checkpoint-backed commit. |
| Graph provider | `entire graph version` | Installed provider is `entire-graph v0.4.0`. |
| Agent guide | `entire graph init-agents --repo .` | The generated graph guide was refreshed. A fresh agent session is still required for the lifecycle requirement. |
| Capabilities | `entire graph capabilities --json` | The provider advertises `CALLS`, `DATA_FLOWS`, `TESTS`, `HANDLES_ROUTE`, and `HANDLES_GRPC`; for Go in the full profile, the supported set includes `CALLS` and `DATA_FLOWS`, but not `TESTS`, `HANDLES_ROUTE`, or `HANDLES_GRPC`. Fidelity must not rely on those three for Go. |
| Snapshot shape | `entire graph snapshot --repo . --format ndjson` | The header declares `relation_evidence` and `relation_resolution`. Relation records use numeric `confidence`, plus `resolution`, `evidence`, and `warning_codes`; there is no native extracted/inferred/ambiguous tri-state field. Fidelity derives its deterministic/advisory classes from these facts. |
| Semantic diff shape | `entire graph diff --base 3a2a715 --head e078008 --json` | The result contains files and named entity changes, locations, kinds, and dependent counts. The Fidelity adapter consumes these changed entities. |
| Relationship controls | `entire graph impact --help`; `entire graph neighbors --help` | Both commands support bounded depth of one or two hops. `neighbors` supports `in`, `out`, and `both` directions. |
| Issue #32 boundary case | Not yet run against a known reproducer | Open verification item. Fidelity must not claim a targeted panic FIX until a real reproducer is confirmed; a general defensive wrapper was added anyway (see below) because it is cheap and correct regardless. |
| Transcript identifiers | Adapter unit coverage for a placeholder identifier | Open verification item: test real ULID and legacy hexadecimal checkpoint IDs once checkpoint history exists. |
| Native partial/incomplete coverage marker (v5 pre-flight item) | `entire graph snapshot`/`capabilities --json` inspection, this session | **A native marker exists, at five layers**: `stats.completeness_level` (ok/degraded/unsafe), `completeness.languages{files,symbols}` and `completeness.relations{TYPE:count}`, `partial_failures[]` with a machine-readable code and an `effect_on_semantic_completeness` string, `language_tiers{lang: semantic\|inventory-only}`, and `capabilities --json`'s `relation_support_by_language`/`relation_support_by_profile`. `fidelity.config.yaml`'s `coverage.detect_via: capabilities_api` (the default) reads these; `heuristic_zero_edge` is the explicit fallback if a future provider build removes them. |
| Tri-state classification key (v1–v4 assumption) | Snapshot header/relation inspection | **The tri-state does not exist as a native field.** Relations carry numeric `confidence` + `resolution` (`exact`/`import_resolved`/`package`/`type_inferred`/`name_only`/`pattern`) instead. `internal/fidelity/classify.go`'s `deterministic`/`advisory` classes (confidence ≥ 0.9 AND a strict-resolution allowlist) are Fidelity's own derived mapping onto that reality, not a passthrough of a provider tri-state. This correction predates the Curveball and is orthogonal to it — see the coverage-confidence signal below for the actually-new axis. |
| Structural edges inflate naive edge-counting | Empirical, via the mandatory fixture test below | `DEFINES`/`CONTAINS` relations exist for **every** symbol regardless of whether anything references it (a method always `CONTAINS`-relates to its type). An early coverage-confidence draft that counted all relation types toward "has an edge" was always `full`, defeating the whole signal. `internal/fidelity/coverage.go` and `reachability.go` both explicitly exclude `DEFINES`/`CONTAINS` from their edge counts and reachability adjacency. |

## Current build boundary

Fidelity now has a real local pipeline for all five stages: transcript
capture, graph symbol/relation/completeness grounding (one snapshot build,
reused by both extraction and reconciliation), checkpoint-to-commit
resolution, semantic changed-entity observation, and full six-tier
reconciliation (`confirmed`, `declared_unimplemented`, `expected_blast_radius`,
`undeclared_scope_creep`, `advisory_low_confidence`, `unverifiable_coverage`)
with N-hop reachability and the coverage-confidence signal wired in per §10.2.
`verdict.json` and `VERIFY_REPORT.md` are written to `--out` on every run, and
a terminal report renders the mandatory three-way visual split. Databricks
work is not yet implemented — see `docs/fidelity-phase-0.md` §2's honest
environment caveat.

## The Noon Curveball: Track 2 — "Graph Is Evidence, Not an Oracle"

**The assumption it invalidated.** Before this response, a changed entity with
zero graph relations was tiered as confident `undeclared_scope_creep`
regardless of *why* the graph was silent — a genuinely disconnected change and
a change resolved only through dynamic dispatch, generated code, or reflection
looked identical to Stage 4. That is precisely "presenting an absence of
evidence as evidence of absence."

**The fix.** A new, orthogonal signal — `coverage_confidence`
(`full`/`partial`/`unknown`) — is computed for every entity reconciliation
touches, using the provider's own completeness signals first (file presence,
language tier, per-file partial failures) before falling back to a zero-edge
heuristic. A new sixth tier, `unverifiable_coverage`, catches anything that
would otherwise have been filed as scope creep on the strength of an absent
path alone, carrying a `verification_path: "manual_review_recommended"`
baseline that holds with Databricks fully disabled. All three renderers make
the three required states visually distinct: confirmed structural evidence,
heuristic/incomplete evidence (the pre-existing deterministic/advisory
classes), and claims needing verification (the new tier). Fully-resolved code
paths are provably unaffected — the mandatory fixture test below asserts it
directly, not just by code inspection.

**Graph demonstration, run before any edit (the scored gate).** `entire graph
impact --repo . --symbol ReconcileDirectClaims --format text` was run at
~13:05 IST, before this session made any change, and surfaced a live,
concrete instance of the exact failure mode the card names: it reported
`grepTreePaths`, `grepFixedStringMatches`, `runWithStderr`, and
`treeBlobMembersBatch` (all in `internal/gitutil`) as transitive callers of
`ReconcileDirectClaims` via `Pipeline.Run` — none of which actually call it.
That is the graph asserting a dependency with more confidence than the
evidence supports, the mirror image of the silent-graph case the card
describes, found in this repository's own graph output rather than assumed.

**Mandatory fixture test.** No organizer-supplied partial-analysis fixture was
found on this machine or attached to the received card text. Per the card's
own instruction not to fabricate confidence, this gap is stated plainly rather
than smoothed over: `internal/fidelity/curveball_fixture_test.go` builds a
disposable, real Git repository with a genuinely unresolvable region (a
function reached only through a runtime string-keyed function-value registry
— verified empirically against the real engine to produce zero non-structural
relations, not assumed) alongside a fully-resolved region in the same commit,
and runs it through the real `NativeGraphAdapter` and the real engine — no
mocked adapter, no hand-built graph. All four required assertions pass. If the
organizer's actual fixture surfaces later, it should replace this one; the
test's assertions are written to the required shape, not to this fixture's
specifics.

**Ranked-prediction accuracy, stated honestly.** None of the five ranked
pivots in the pre-noon curveball briefing (unified query engine, peer-review
handoff, PII redaction, multi-repo DATA_FLOWS, real-time sync) named this
constraint. The closest partial relevance is the multi-repo DATA_FLOWS
prediction, whose evidence included Issue #32's boundary panic — it shares the
general instinct that graph queries can hit an unresolved edge case and must
degrade gracefully rather than assert confidently, and that instinct
generalizes usefully here (the defensive `recover()` wrapper added to
`NativeGraphAdapter.Graph` is a direct descendant of it). But the
coverage-confidence signal itself — a first-class, always-computed field
distinct from crash handling — is new work this response adds, not something
the pre-curveball architecture already covered.

**What did not change.** Stage 1–3 extraction and Stage 5's
pure-function-of-`verdict.json` property both hold. This was a Stage 4
reconciliation-logic change plus a schema addition, exactly where the
pre-noon structural firewall (adapter layer, config-driven tiers) said a
"what counts as drift" change should land.
