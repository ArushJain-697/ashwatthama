# Fidelity

**BTW Buildathon 2026 · Entire Track E2 — Build with Graph Intelligence · + Best Use of Databricks (opt-in)**

> This file follows the exact `BUILDATHON.md` outline from the Participant Guide. Sections marked
> `🔲 FILL DURING BUILD` contain the live, on-the-day facts (links, commit SHA, verified CLI flags,
> test results) that don't exist until the event actually runs — everything else is the locked
> design decided ahead of time. Update the flagged sections as the day progresses; don't leave any
> unresolved at submission.

---

## 🔲 Pre-Flight Verification Findings (run first, 9:00–9:30)

Record the ground truth here before writing any pipeline code — several later sections assume it.

| Check                                       | Command                                                       | Result                                                                                                                                             |
| ------------------------------------------- | ------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| Actual edge-type enum                       | `entire graph capabilities --json`                            | 🔲 *(confirm `CALLS` / `CONTAINS` / `DATA_FLOWS` / `type`; do not assume `TESTS`, `HANDLES_ROUTE`, `HANDLES_GRPC` exist)*                           |
| Snapshot shape + tri-state key              | `entire graph snapshot`                                       | 🔲                                                                                                                                                  |
| Diff output shape                           | `entire graph diff --base <ref> --head <ref>`                 | 🔲                                                                                                                                                  |
| Depth/direction flags on impact/neighbors   | `entire graph impact --help`, `entire graph neighbors --help` | 🔲 *(if absent, N-hop reachability is computed by iterating `neighbors` ourselves)*                                                                 |
| Checkpoint ID format                        | —                                                             | 🔲 *(confirm ULID and legacy hex are both handled)*                                                                                                 |
| Issue #32 panic reproduction                | malformed/boundary file through `graph search`                | 🔲 *(if it reproduces, wrap the call defensively from the start)*                                                                                   |
| Native "partial/unresolved" coverage marker | `entire graph capabilities --json` / `entire graph snapshot`  | 🔲 *(if none exists, coverage-confidence falls back to the zero-edge heuristic — see Curveball section)*                                            |
| Lakebase provisionable under Free Edition   | live `lakebase` project-creation attempt                      | 🔲 *(sources disagree; if unavailable, the Databricks module runs on the already-verified AI Search + logistic-regression/CRC design, no Lakebase)* |
| Confirmed submission deadline               | ask at kickoff                                                | 🔲 *(Participant Guide states 3:00 PM IST; a separate event-site source has said 4:00 PM — get the live, spoken answer and write it here)*          |

---

## One-Sentence Summary

Fidelity is an intent-vs-implementation verifier for AI coding agents that deterministically reconciles a stated Checkpoint's intent against the Entire Graph's structural evidence of what actually changed.

---

## Problem, Intended User, and Why It Matters

**Intended user:** the person who has to approve an AI agent's work before it merges — a tech lead reviewing a teammate's agent session, or an engineer reviewing their own agent's output before pushing.

**Problem:** agents are fast, but reviewers have no reliable way to know whether an agent did *exactly* what was asked, *more* than was asked (silent scope creep), or *less* than was asked (task left half-done). Today the only options are re-reading the entire diff by hand, or trusting an LLM's narrative summary of the transcript — a judgment call, not a verification.

**Why it matters:** fast AI generation needs equally fast, trustable verification. Fidelity turns code review from a semantic impression into a deterministic audit trail.

**Why it needs Entire specifically:** Checkpoints are the only ground-truth stated intent tied to a specific atomic change. Graph diff is the only structural ground truth of what actually moved. Neither one alone can answer the reviewer's question — Fidelity exists at their intersection.

---

## Selected Entire Track and Why Entire Is Essential

**Track: E2 — Build with Graph Intelligence.**

Entire is essential because it uniquely supplies both halves of the verification equation: Checkpoints provide the stated intent tied to an atomic change, and Graph diff provides the structural ground truth of what changed. Fidelity doesn't just call Entire commands — the reconciliation logic *is* a function of Graph output; remove Entire and there is no product left, just an unverifiable narrative.

### Differentiation from Entire's own built-in skills

Entire already ships skills that sound adjacent to Fidelity — `review` (audits changes by reading checkpoint transcripts to produce intent-aware findings), `explainskill` (looks up the session behind a function to explain why it exists), and `what-happened` (traces the latest change via git blame + checkpoint context for debugging). All three are **retrieval-and-narration** tools: they answer "why is this code the way it is" by surfacing and summarizing session context, largely through an LLM reading a transcript.

Fidelity answers a different question — **"did the change match what was claimed"** — and answers it with a deterministic, tiered, graph-verified reconciliation, not an LLM's narrative judgment. Every tier in `verdict.json` is backed by an actual graph reachability computation over classified edges, not a semantic impression. `review` could use Fidelity's `verdict.json` as a structured input to sharpen its own narrative — the two are complementary layers, not competitors.

### Differentiation from the wider graph-tooling landscape

| Tool                     | What it actually does                                                                                                            | How Fidelity differs                                                                                                                                                                                                                                   |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Graphify**             | Builds a persistent, queryable multimodal knowledge graph for agent context/exploration; tags edges EXTRACTED/INFERRED/AMBIGUOUS | The tri-state idea isn't Fidelity's differentiator — it converges on the same taxonomy independently. Fidelity's differentiator is what it *does* with the signal: a deterministic post-hoc verdict, not a context-retrieval aid for an agent mid-task |
| **Serena**               | Precision code editing/symbol resolution via compiler-grade LSP                                                                  | Higher raw precision, but no checkpoint/intent linkage at all — it edits correctly, it verifies nothing against stated intent                                                                                                                          |
| **Sourcegraph/Cody**     | Code search and assistant context                                                                                                | A retrieval tool, not a verification/reconciliation tool                                                                                                                                                                                               |
| **CodexGraph/RepoGraph** | Agent queries its own targeted subgraph via LLM-written Cypher                                                                   | Agent-directed retrieval, not an independent auditor of what the agent already claimed and did                                                                                                                                                         |

**Ready answer if asked "isn't this just Graphify's tri-state model?"** — Graphify's tagging routes an agent's *attention* during exploration; Fidelity's drives a deterministic *verdict* after the fact. Tagging an edge INFERRED helps an agent decide how much to trust it while working. Fidelity uses the same kind of signal to decide, after the work is done, whether a specific claim holds — reconciliation, not retrieval.

We also deliberately never adopt Graphify itself as a dependency — doing so would blur the exact positioning above rather than strengthen it.

---

## Architecture and Main Workflow

```
┌─────────────────┐   ┌──────────────────┐   ┌────────────────────┐
│ 1. CAPTURE       │   │ 2. DECLARE        │   │ 3. OBSERVE          │
│ Read checkpoint  │──▶│ Extract intended  │   │ entire graph diff   │
│ transcript       │   │ entities          │   │ + snapshot edges    │
└─────────────────┘   │ (regex → LLM fbk) │   │ (+ tri-state class) │
                       └──────────────────┘   └────────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌───────────────────────────────────────┐
                       │ 4. RECONCILE                           │
                       │ Tier every changed entity via config   │
                       │ rules + N-hop reachability over        │
                       │ CONFIRMED edge types; flag anything    │
                       │ resting on inferred/ambiguous edges    │
                       └───────────────────────────────────────┘
                                         │
                                         ▼
                       ┌───────────────────────────────────────┐
                       │ 5. MANIFEST                            │
                       │ verdict.json, wrapped in a Checkpoint  │
                       └───────────────────────────────────────┘
                           │             │             │
                           ▼             ▼             ▼
                     Terminal report  VERIFY_REPORT.md  dashboard.html

                 [ Optional bolt-on, never touches 1–4: ]
                 ┌───────────────────────────────────────┐
                 │ 6. CALIBRATE (Databricks)               │
                 │ Advisory edges → conformal gate →       │
                 │ AI Search corroboration → calibrated_   │
                 │ confidence written back into verdict    │
                 └───────────────────────────────────────┘
```

**Command surface:**

```
entire graph verify-intent <checkpoint-id> \
    [--base <checkpoint-id>] \
    [--config fidelity.config.yaml] \
    [--out ./fidelity-out/]
```

**Extraction logic (Stage 2):** load the full symbol dictionary from `entire graph snapshot`, fuzzy-match transcript tokens against it, and send only unresolved fragments plus the constrained symbol dictionary to an LLM fallback that must pick from the exact list or return "no confident match" — it fails closed, never inventing an entity. Vague prompts ("continue", "fix it") inherit intent from the nearest resolvable ancestor checkpoint (`inherited_intent: true`); a first checkpoint diffs against an empty tree (`no_baseline: true`).

**Reconciliation logic (Stage 4):** every changed entity is tiered as `confirmed`, `declared_unimplemented`, `expected_blast_radius` (via deterministic-edge N-hop reachability), `undeclared_scope_creep`, or `advisory_low_confidence` (anything that would only hold via an `inferred`/`ambiguous` edge — never silently promoted to "expected"). See the Curveball section below for the sixth tier added mid-event.

**Rendering (Stage 5):** all three renderers — terminal report, `VERIFY_REPORT.md`, and `dashboard.html` (a single static file that `fetch()`es `verdict.json`, radial/force graph color-coded by edge class) — are pure functions of `verdict.json`. The dashboard is the screen open during judging.

**Open-source and AI components used (transparency):**
- **Qwen2.5-Coder** (open-weight) — the fail-closed LLM extraction fallback in Stage 2.
- **Semgrep** (open source) — pattern-based dynamic-dispatch/reflection/codegen detection feeding the coverage-confidence signal (Curveball response, see below).
- Optional, evaluated but not load-bearing: **MAPIE** (conformal prediction library, §12.2 upgrade), **leidenalg + python-igraph** (Leiden community detection, stretch), **sentence-transformers + FAISS/`rank_bm25`** (local hybrid-retrieval fallback if AI Search proves flaky).
- **Never adopted:** Graphify itself, CodeQL/Kythe (heavier alternatives to Semgrep for the same payoff — named as roadmap, not built).

---

## Entire Graph Findings and Verification

Three required graph demonstrations, explicitly named:

1. **Graph Definition Lookup** — the full symbol dictionary is loaded via `entire graph snapshot` to ground the LLM extraction vocabulary and capture edge-class keys (Stage 1/2).
2. **Impact Analysis Before a High-Risk Change** — before editing the code affected by the Noon Curveball, we ran a mandatory `entire graph impact` (or iterated `neighbors`, per pre-flight findings above) to identify the affected dynamic-dispatch entities *before* implementation, not after.
3. **Final Semantic-Diff Analysis** — `entire graph diff --base <base> --head <checkpoint-id>` feeds the final pipeline run that produced the submitted `verdict.json`.

Per the guide's own instruction, graph results are evidence, not an oracle: every finding above is cross-checked against source and tests before being trusted, and the reconciliation tiers are designed so the product itself never asserts more certainty than the graph actually has (see the tri-state model in Stage 3/4, and the coverage-confidence signal below).

---

## Noon Curveball: What Changed and How We Adapted

**Constraint received:** Track 2 — **"Graph Is Evidence, Not an Oracle."** The repository under review contains dynamic dispatch, generated code, or reflection that static analysis cannot fully resolve. Requirements: never present incomplete graph relationships as certain; identify when analysis may be partial; provide a safe fallback/verification path; keep existing behavior for fully-resolved code unchanged; test against a supplied partial-analysis fixture; let users tell apart confirmed evidence, heuristic/incomplete evidence, and claims needing source/test verification.

**Honest note on our own predictions:** none of our five pre-event ranked pivot predictions named this exact constraint. The closest partial relevance was our rank-4 candidate (multi-repo `DATA_FLOWS` tracing, evidence: Issue #32's boundary panic) — it shared the general instinct that graph search can hit an unresolved edge case and must degrade gracefully, and the defensive-wrapping habit built for Issue #32 generalized usefully here. But the coverage-confidence signal itself is new work this event added, not something our existing architecture already absorbed for free.

**The invalidated assumption:** we had assumed a zero-edge result was confident proof of scope creep. The curveball invalidated this — dynamic dispatch means the graph's silence can just mean the code is unresolvable by static analysis, not that it's genuinely disconnected.

**The fix:** a new `coverage_confidence` signal (`full` / `partial` / `unknown`), orthogonal to the existing edge tri-state (`extracted`/`inferred`/`ambiguous`). Edge tri-state classifies confidence in a relationship *the graph recorded*; coverage-confidence classifies confidence in whether the graph's *silence* about an entity means anything at all. An entity with `coverage: partial` or `coverage: unknown` can never be filed into `undeclared_scope_creep` on the strength of an absent path alone — it is instead routed to a new sixth reconciliation tier, `unverifiable_coverage`, flagged `verification_path: "manual_review_recommended"` as a baseline that works with or without Databricks. Entities with `coverage: full` proceed through the existing tiers exactly as before — fully-resolved code's behavior is unchanged, per the curveball's explicit requirement.

**Open-source augmentation:** we integrated **Semgrep** to run prebuilt dynamic-dispatch/reflection/codegen rule packs against zero-edge entities, so `partial` coverage is a **detected finding**, not just an absence heuristic — the single highest-leverage addition against the mechanism this curveball is actually scoring.

**Interface requirement:** all three renderers visually separate three states at a glance — confirmed structural evidence, heuristic/incomplete evidence, and claims requiring source/test verification — never a shared "risky" bucket with a tooltip.

**Mandatory test (against the supplied partial-analysis fixture):** asserts (1) no entity in the fixture's dynamic-dispatch region is misfiled into `undeclared_scope_creep`; (2) every affected entity carries `coverage_confidence: partial` (or `unknown`) and lands in `unverifiable_coverage`; (3) `verification_path` is populated; (4) a fully-resolved entity elsewhere in the same fixture still reconciles normally. Result: 🔲 *(pass/fail, fill in after running)*.

---

## Checkpoint Links and What Each Proves

| #   | Link              | What it proves                                                                                                                                                                                                              |
| --- | ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | 🔲 `[INSERT LINK]` | Initial understanding and intended architecture — including the proactive decision to build a repository/adapter firewall layer so an unrelated API-convergence pivot could be absorbed via config, not a rewrite.          |
| 2   | 🔲 `[INSERT LINK]` | The last stable state before the Noon Curveball — a runnable pipeline with a passing regression suite.                                                                                                                      |
| 3   | 🔲 `[INSERT LINK]` | Response to the Noon Curveball — reconstructed intent in a fresh session, explicitly names Track 2, logs the pre-edit `graph impact` run, and documents the `coverage_confidence` fix and the invalidated assumption above. |
| 4   | 🔲 `[INSERT LINK]` | Final implementation and verification — the final semantic-diff run, the passing partial-analysis fixture test, and the three-way visual split across all renderers.                                                        |

---

## Setup, Run and Test Instructions

```bash
# 1. Sign in and mirror (after the official 9:00 AM start — fork first, choose the India region)
entire login
entire repo mirror create
entire repo clone /gh/YOUR-GITHUB-HANDLE/REPOSITORY

# 2. Enable checkpoints before any agent-assisted development
entire enable -y --agent YOUR-BUILD-AGENT
entire status

# 3. Activate the graph plugin, then start a FRESH agent session immediately after
entire plugin install graph
entire graph version
entire graph init-agents --repo .

# 4. Run Fidelity
entire graph verify-intent <checkpoint-id> --config fidelity.config.yaml
```

**Testing:**
- Golden-manifest tests — hand-built (transcript, diff) fixtures against an exact expected `verdict.json`.
- Fail-closed test — a non-matching fragment through the LLM fallback must produce `unresolved_fragment`, never an invented entity.
- No-baseline test — a first-checkpoint case yields only `confirmed`/`undeclared_scope_creep`, tagged `no_baseline: true`.
- Graceful-panic test — a known Issue #32–style boundary case degrades to an explicit warning, not a crash.
- Renderer purity test — all three renderers agree on tier counts from the same `verdict.json`.
- Reproducible setup — one script builds a tiny sample repo and runs `verify-intent` to a known verdict.
- Partial-coverage fixture test — the mandatory Curveball test described above.

Test suite result: 🔲 *(fill in — e.g. "12/12 passing as of commit <SHA>")*.

---

## Databricks Use, Data Sources and Limitations

**Opted in:** yes — a Calibrated Trust-But-Verify Gate on top of Stage 5's output. Never touches Stages 1–4's interfaces.

**Capabilities used and why essential:** the tri-state edge model tells us an edge is `inferred` or `ambiguous`, but not *how* uncertain, and gives no path to independent corroboration.
- **I-CALM prompting** elicits a calibrated verbal confidence score from the extraction model at prompt time.
- **Conformal Risk Control (CRC)** turns that score into a threshold `λ̂` with an actual bounded false-positive-rate guarantee on a calibration set, rather than an arbitrary cutoff.
- **Databricks AI Search hybrid retrieval** (vector + full-text + Reciprocal Rank Fusion, one `query_type="hybrid"` call) corroborates whatever the gate can't resolve on its own, querying an index built over the repo's own docstrings, comments, and commit messages — governed automatically by Unity Catalog.
- **Curveball-enhanced reuse (§12.15):** the same AI Search index is queried a second way for `unverifiable_coverage` entities — a genuine low-marginal-cost extension of the mandatory baseline, never a substitute for it. The `manual_review_recommended` flag stands alone whether or not Databricks is enabled.

**Why not Genie Ontology:** it's a closed, auto-inferred system with no direct graph write API — natural-language queries resolve to SQL on a SQL warehouse, and each agent is capped at 200 knowledge-store snippets shared across tables. A code graph's schema exceeds that instantly, so the code graph was never forced into it.

**Free Edition fit:** the calibration model is a lightweight scoring mechanism (logistic regression as fallback tier — no GPU need); AI Search calls are capped and occasional, only reaching `pending_verification` edges; the demo defaults to batch scoring (`score_source: batch`) to avoid live Serving-endpoint quota risk during judging.

**Fallback tier (if the full I-CALM/CRC/AI-Search chain doesn't land in time):** a plain logistic regression trained in MLflow on calibration features (edge class, edge type, language, churn size), writing a single `calibrated_confidence` float back into each advisory edge — still legitimate, still qualifies as meaningful use, just less differentiated. `gate.mode: fallback` in config signals which tier is live; the schema doesn't change either way.

**Data provenance and honest limitations:** calibration labels are synthetic/injected, not real production revert data, given the time window. The AI Search index is built only from the demo repo's own docs/comments/commit messages — a real deployment would need a much larger, curated corpus. This demonstrates the *mechanism*, not a production-calibrated system.

**Lakebase stretch features (§12.10–§12.13 — zero-copy counterfactual "what-if" verdicts, a live cross-session Verdict Ledger, curveball-reconstruction memory):** contingent on the Free Edition provisioning question flagged in Pre-Flight Verification above. Status: 🔲 *(built / evaluated and skipped — record which, and why, here)*.

**Workspace/app/endpoint URL:** 🔲 `[INSERT LINK]`
**Relevant repo paths:** 🔲 `[INSERT PATHS]`
**Reproduction steps:** 🔲 `[INSERT — or point to Setup section above if identical]`

---

## Known Limitations and Next Steps

**Limitations, stated honestly:**
- Calibration labels for the Databricks gate are synthetic/injected, not real revert history.
- The AI Search corroboration index is built only from the small hackathon demo repo's own docs and comments.
- Community-aware scope-creep clustering (Leiden) needs a graph large and connected enough to produce meaningful architectural boundaries — a small demo repo may not show this convincingly.
- Coverage-confidence detection beyond Semgrep's pattern packs is still a heuristic where no native "partial/unresolved" marker exists in the graph API.

**What actually made it into the build vs. stayed roadmap-only:** 🔲 *(fill in at the end — e.g. "Semgrep-backed coverage detection: built. MAPIE, Leiden clustering, local hybrid-retrieval fallback, Lakebase stretch features: evaluated, not built, see reasoning above.")*

**Next steps (practical continuation path):**
1. **Plan-Trust Ledger** — persist (claim, outcome) pairs across sessions to learn how much to trust a given agent/model's plans over time; promotable via the Lakebase Verdict Ledger if that cleared Free Edition provisioning.
2. **Checkpoint Court** — a second, adversarial agent that tries to falsify each verdict claim before acceptance, logged as a replayable trial record.
3. **Forensic Replay** — time-travel on versioned snapshots to reconstruct what an agent actually saw at a past decision; the Lakebase "what-if verdict" is a narrower, live-buildable version of this.
4. **Subagent Reasoning Capture** — reconcile each subagent's sub-intent against its own scoped graph diff.
5. **Compliance reconciliation via Unity Catalog lineage** — resolve a code change to which regulated fields/attestations it touches.
6. **SCIP export alignment** — once Entire's RFD 0006 SCIP-export effort ships, expose `verdict.json` through that standard channel for cross-tool interoperability.
7. **Fidelity as an MCP-exposed self-check tool** — let another coding agent call `verify-intent` on itself before ever handing a diff to a human; low-effort once the schema is stable, deliberately not built same-day.
8. **CodeQL / Kythe** — heavier alternatives to Semgrep for the same coverage-detection payoff; worth naming as a roadmap line, not a same-day add.

**Explicitly out of scope:** cross-repo structural graph analytics (`entire-graph` enforces single-store isolation); any feature depending on `TESTS`, `HANDLES_ROUTE`, or `HANDLES_GRPC` edges (unverified in the source engine); Genie Ontology as the graph store; Graphify as a dependency.

---

## Positioning (for the pitch)

> *"Every AI coding agent tells you what it did. None of them prove it did only what you asked. Fidelity reads the Checkpoint's stated intent, reads the Graph's actual structural change, and produces a verdict you can check against source in seconds — not a chat answer, an audit trail. And when the graph itself isn't certain, Fidelity isn't certain — it degrades honestly instead of asserting a dependency that might not exist."*

Answers rehearsed cold:
- **"Isn't this just `review`?"** — `review` narrates why a change makes sense by reading the transcript with an LLM. Fidelity verifies whether it matches what was claimed, tier by tier, backed by actual graph reachability over classified edges — a judgment call versus a computed reconciliation.
- **"Why not Genie Ontology?"** — the 200-snippet cap and SQL-only traversal make it the wrong tool for a code graph's schema.
- **"Isn't this just Graphify?"** — Graphify's tagging routes an agent's attention during exploration; Fidelity's drives a deterministic verdict after the fact.
- **"What did the curveball actually change?"** — see the Curveball section above, rehearsed verbatim.

---

## 🔲 Final Demo Readiness

- [ ] State the user and problem in one sentence (above).
- [ ] Show the working product and the critical path live, not a slide-only walkthrough.
- [ ] Explain why Entire is essential to the solution.
- [ ] Show one useful checkpoint and one graph finding that changed or verified a decision.
- [ ] Explain the Noon Curveball, the behavior that changed, and the test that proves it.
- [ ] Show the essential Databricks function and evidence it's working.
- [ ] Close with known limitations and the next step toward production readiness.

## 🔲 Final Submission Checklist

- [ ] Final commit pushed; SHA matches the submission: `[INSERT SHA]`
- [ ] Project launches from a clean checkout / documented setup path.
- [ ] All four required checkpoints open and clearly explain their milestone.
- [ ] Entire Graph evidence and the final semantic-diff analysis are recorded and named.
- [ ] `BUILDATHON.md` complete, readable, free of secrets.
- [ ] Tests covering critical + Curveball behavior pass.
- [ ] Databricks resource links and data notes included.
- [ ] Demo owner (`[NAME]`) can sign in, open every resource, run the critical path.
- [ ] A fallback screenshot/recording exists locally: `[INSERT PATH/LINK]`
- [ ] Submitted before the confirmed deadline (see Pre-Flight Verification Findings above).

**Submission package fields:**
- Selected track: **E2 — Build with Graph Intelligence**
- GitHub fork URL: 🔲 `[INSERT]` — Final commit SHA: 🔲 `[INSERT]`
- Entire mirror/project URL: 🔲 `[INSERT]`
- Working demo or fallback recording: 🔲 `[INSERT]`
