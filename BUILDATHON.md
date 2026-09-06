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
| Issue #32 boundary case | Not yet run against a known reproducer | Open verification item. Fidelity must not claim a panic guard until a real reproducer is confirmed. |
| Transcript identifiers | Adapter unit coverage for a placeholder identifier | Open verification item: test real ULID and legacy hexadecimal checkpoint IDs once checkpoint history exists. |

## Current build boundary

Fidelity currently has a real local pipeline for transcript capture, graph
symbol grounding, checkpoint-to-commit resolution, semantic changed-entity
observation, and direct confirmed/unimplemented reconciliation. N-hop
reachability, scope tiers, locked `verdict.json` rendering, UI renderers, and
all Databricks work remain separate build tickets.
