# Fidelity

**Fidelity is an intent-versus-implementation verifier for AI coding agents.**
It reads an Entire Checkpoint's stated intent, reads Entire Graph's actual
structural change, and produces a tiered, `file:line`-cited reconciliation a
reviewer can check in seconds — a computed verdict, not a chat summary.

## Problem, user, and why Entire

**User:** the person who has to approve an AI agent's work before it
merges — a tech lead reviewing a teammate's agent session, or an engineer
reviewing their own agent's output before pushing.

**Problem:** agents are fast, but a reviewer has no reliable way to know
whether an agent did *exactly* what was asked, *more* (silent scope creep), or
*less* (task left half-done) without re-reading the whole diff by hand or
trusting an LLM's narrative summary of the transcript.

**Why Entire specifically:** a Checkpoint is the only ground-truth stated
intent tied to one atomic change; Entire Graph's semantic diff is the only
structural ground truth for what actually changed. Neither alone answers the
question; Fidelity is the reconciliation between them.

**Track:** E2 — Build with Graph Intelligence.

**"Isn't this just `review`/`explainskill`/`what-happened`?"** Those tools
retrieve and narrate why a change exists by reading session context with an
LLM. Fidelity checks whether the claimed intent matches the structural code
change — a computed, tiered reconciliation backed by graph reachability over
classified edges, not a narrative judgment. They compose: `review` could use
Fidelity's `verdict.json` as structured evidence.

**"Why not Genie Ontology?"** Its 200-snippet cap and SQL-only traversal are a
poor fit for code graph structure. Entire Graph stays the structural source of
truth; Databricks, where used, is an optional confidence/corroboration layer,
never the graph store.

**"Isn't this just Graphify's tri-state model?"** See
`docs/fidelity-phase-0.md` — Graphify's tagging routes an agent's *attention*
during exploration; Fidelity's signal drives a deterministic *verdict* after
the fact. Reconciliation, not retrieval.

## Architecture

Five stages, each a clean interface boundary (`internal/fidelity/adapter.go`'s
`Stage` enum), with every `entire`/graph call routed through one `Adapter`
seam so a later API convergence is a config/adapter change, not a rewrite:

```
1. CAPTURE    read checkpoint transcript (via Entire's checkpoint CLI)
2. DECLARE    ground direct symbol mentions against one graph snapshot
3. OBSERVE    entire graph diff -> changed entities
4. RECONCILE  six-tier reachability + coverage-confidence tiering
5. MANIFEST   verdict.json, terminal report, VERIFY_REPORT.md
```

Stage 4 produces six tiers: `confirmed`, `declared_unimplemented`,
`expected_blast_radius` (N-hop reachability over deterministic edges only),
`undeclared_scope_creep`, `advisory_low_confidence` (reachable only via a
low-confidence edge), and `unverifiable_coverage` — the Curveball response,
see below. All three renderers are pure functions of `verdict.json`.

## Setup, run, and test

```sh
mise exec -- go build -o entire-graph ./cmd/entire-graph
mise exec -- go test ./internal/fidelity/... ./internal/cli/...

# Run against a real checkpoint pair:
./entire-graph verify-intent <checkpoint-id> --base <base-checkpoint-id> \
    --repo . --config fidelity.config.yaml --out fidelity-out
# writes fidelity-out/verdict.json and fidelity-out/VERIFY_REPORT.md,
# prints the terminal report to stderr, and the JSON response to stdout.
```

## Checkpoints

1. **Initial understanding / architecture** — not captured as a dedicated
   Checkpoint; see the honest note in "Checkpoint content, stated honestly"
   below.
2. **Pre-noon stable state** — `caefd603ffcd` (commit `bb61095`).
3. **Response to the Curveball** — `af04126413fd` (commit `1e334ab`): the
   fresh-session reconstruction, the pre-edit `impact` run, and the full
   six-tier/coverage-confidence implementation.
4. **Final implementation and verification** — `9ae11dfa25d3` (commit
   `6cd3d05`): the whole-repo test pass, the live end-to-end
   `verify-intent` run, the third graph demonstration, and this record.

**Checkpoint content, stated honestly.** `caefd603ffcd`'s actual transcript
content is a single mechanical instruction ("create a checkpoint and commit")
executed by an unrelated agent session — it has no architecture decision
content. No dedicated Checkpoint 1 was ever written. The real architectural
reasoning (the adapter-layer firewall, the tri-state correction, the
coverage-confidence design) lives in this session's transcript
(`af04126413fd`) and in `docs/fidelity-phase-0.md`, not in a separate
checkpoint per milestone. Recorded here plainly rather than smoothed over,
per house style.

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

## Third required graph demonstration: final semantic diff

`entire graph diff --repo . --base bb61095 --head 1e334ab` was run against
this response's own final commit (checkpoint `af04126413fd`) versus the
pre-noon stable commit. 202 lines of entity-level changes across 21 files,
including `Pipeline.Run body changed (565 dependents)` — the tool's own
heuristic dependent count flagging exactly the kind of signature change that
should not ship untested. It did not: `go test ./internal/fidelity/...
./internal/cli/...` passed in full both before this commit and after,
including the mandatory fixture test and the six-tier reconciliation test
against a real synthetic graph. Full diff output committed at
`docs/evidence/2026-09-06-final-semantic-diff.txt`.

## Live end-to-end run

`entire graph verify-intent af04126413fd --base caefd603ffcd --repo . --out
fidelity-out` was run against this response's own real checkpoints — not a
unit-test fixture — reading this whole session's actual transcript as intent
and this repository's real graph as evidence. Committed at `fidelity-out/`.

```
158 changed entities: 65 confirmed, 0 expected blast radius, 1 scope creep,
19 advisory, 73 unverifiable coverage, 582 declared unimplemented
verdict_label: REVIEW_REQUIRED
```

One `unverifiable_coverage` entry worth pointing at directly in a demo: a
Markdown section header was correctly tiered `coverage_confidence: partial`
with reason *"this file's language (Markdown) is inventory-only: the graph
records file/symbol structure but does not attempt relationship extraction
for it"* — the coverage-confidence signal firing correctly on a real
inventory-only-language case, not a synthetic one.

## Three required graph demonstrations, named

1. **Symbol dictionary load** — Stage 2's `adapter.Graph()` calling `entire
   graph snapshot` (via `sem.BuildProviderSnapshotWithOptions`), now also
   grounding Stage 4's reachability and coverage-confidence from the same
   build.
2. **Pre-edit impact analysis** — `entire graph impact --symbol
   ReconcileDirectClaims`, run before any Curveball-response edit, surfacing
   the false-positive transitive-caller finding recorded above.
3. **Final semantic diff** — the `entire graph diff` run immediately above,
   against this response's own final checkpoint.

## Lakebase stretch phase (#8, #57-#60)

**#8 — Lakebase Free Edition availability: not clearable, recorded honestly.**
No live project-creation call could be attempted: no `databricks` CLI, no
`~/.databrickscfg`, no `DATABRICKS_*` env vars, and no `databricks-sdk`
installed on this machine — there is no workspace to attempt it against. Per
the ticket's own gate ("the single gate deciding whether #58-60 are attempted
at all"), **#58, #59, and #60 are skipped entirely**, exactly as the buildmap
instructs when #8 does not clear clean — not silently dropped, but the
ticket's own specified outcome for this case.

**#57 — Community-aware scope-creep framing: built, and genuinely
non-trivial on this repository.** Leiden clustering (`leidenalg` +
`python-igraph`, the pinned Bible implementation, installed via `pip3
install python-igraph leidenalg` — no Databricks or network dependency of
its own) ran over this repo's own real relation graph (36,574 non-structural
edges — `DEFINES`/`CONTAINS` excluded for the same reason `coverage.go` and
`reachability.go` exclude them):

```
149 communities, largest = 12.9% of the graph, 0% singletons
top labels: internal/sem, internal/cli, internal/fidelity, internal/gitutil,
            internal/termsafe, scripts, bench/memory/benchmarks/common
```

Non-trivial by the ticket's own bar (no giant single cluster, no
all-singleton degenerate case), and the labels are real subsystems, not
noise — this repo clears the caveat rather than needing it invoked.

Wired as an explicitly optional enhancement, not a Stage 4 core dependency:
`scripts/leiden_communities.py` runs out-of-process against
`entire graph symbols`/`edges` NDJSON output; `internal/fidelity/community.go`
loads its `communities.json` and annotates `undeclared_scope_creep` entries
(`community_label`) via the new `verify-intent --community-map <path>` flag.
No flag, no Python run — unchanged behavior; this never touches the
mandatory tiers. Reproduce with:

```sh
./entire-graph symbols --repo . --format ndjson > /tmp/symbols.ndjson
./entire-graph edges --repo . --format ndjson \
    --relation CALLS,DATA_FLOWS,USES_TYPE,PARAM_TYPE,RETURNS_TYPE,READS_FIELD,WRITES_FIELD,ACCESSES,EXTENDS,IMPLEMENTS,INHERITS,OVERRIDES,CONSTRUCTS,ASYNC_CALLS \
    | python3 scripts/leiden_communities.py /tmp/symbols.ndjson > /tmp/communities.json
./entire-graph verify-intent <checkpoint-id> --community-map /tmp/communities.json --repo . --out fidelity-out
```

## Databricks use

**Opted in** (a special prize was offered for the best use). **Not
implemented in this build window.** This machine has no `databricks` CLI, no
`~/.databrickscfg`, and no `mlflow` — no workspace credentials were reachable
inside the remaining time, and the honest choice was to say so rather than
write calibration/AI-Search glue code that cannot actually reach a workspace
and call it "meaningful use." The schema (`fidelity.config.yaml`'s
`databricks:` block, `verdict.json`'s `verbal_confidence` /
`conformal_gate_passed` / `calibrated_confidence` / `corroboration` fields) is
already wired end-to-end and stable either way — every advisory and
`unverifiable_coverage` entry carries those fields as `null`/`false`, exactly
as the schema promises when Databricks is disabled. The nearest reachable
tier if a workspace becomes available is §12.7's fallback (a local logistic
regression writing `calibrated_confidence`, `gate.mode: fallback`) — legitimate
per the Bible, and schema-compatible with the fuller I-CALM/CRC/AI-Search
chain if that is built later.

## Known limitations and next steps

- **No organizer-supplied partial-analysis fixture was available** for the
  mandatory Track 2 test; a real, empirically-verified substitute was built
  instead (`internal/fidelity/curveball_fixture_test.go`). Replace it if the
  actual fixture surfaces.
- **Databricks is schema-ready but not implemented**, per the section above.
- **Issue #32's boundary panic was never reproduced** against this
  repository; a general defensive `recover()` wrapper was added at the one
  seam Fidelity controls regardless, since it is correct with or without a
  reproducer, but no targeted fix is claimed.
- **The LLM fallback extraction path (`ResolveFallback`) has no model wired
  up** — it fails closed (returns `unresolved_fragment`) by construction, so
  disabling it costs recall, not correctness. Only direct, exact symbol
  mentions are currently extracted.
- **No fork/mirror work was done post-curveball** — the remaining window was
  spent entirely on the local pipeline per team direction; a demo/deployment
  owner has not yet been named.
- **`entire-judge` self-audit (§13.4) was not run** against this checkpoint
  history — worth doing before submission if time allows.
- **Lakebase (#58-#60) was skipped**, per #8's own gate — no Databricks
  workspace was reachable to attempt the availability check against. Leiden
  clustering (#57) WAS attempted and is genuinely non-trivial on this repo —
  see the Lakebase stretch phase section above.
- **The dashboard renderer (#40)** is in progress; a design-decision question
  was raised with the user before building it, per their explicit request.
- **A fallback demo recording does not exist yet** — this needs a human to
  actually run the live demo and capture it; recorded here as an open item
  for submission, not something this session can produce.
