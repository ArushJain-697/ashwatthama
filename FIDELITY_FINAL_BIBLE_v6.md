# FIDELITY — Project Bible v6
### An Intent-vs-Implementation Verifier for AI Coding Agents
**BTW Buildathon 2026 · Entire Track E2 — Build with Graph Intelligence · + Best Use of Databricks (opt-in)**
**⚡ LIVE BUILD — Noon Curveball received and in effect. See §0 (v5) and §13.5. Open-source augmentation pass added, §0 (v6) and new §20.**

> Single source of truth. This revision folds in the official Participant Guide's operational requirements (exact CLI commands, exact checkpoint/graph deliverables, exact submission format, Databricks Free Edition constraints), two research passes (trust-but-verify calibration design; curveball-prediction intelligence), **the actual Noon Curveball as received**, with the concrete pipeline change it forces, and now **a prioritized open-source model/tool augmentation pass (§20)** scoped explicitly against the time remaining before the 3:00 PM deadline.

---

## 0. What Changed in This Revision

- **Operational grounding.** The exact `entire` CLI commands, the exact four required checkpoints, the exact three required graph demonstrations, and the exact submission package now come directly from the Participant Guide, not inference. See §3.
- **⚠️ Timeline conflict, flagged not resolved.** The Buildathon website (fetched earlier) says code freeze / submission at **4:00 PM**, event ends 5:00 PM. The uploaded Participant Guide says submission deadline **3:00 PM**, judging 3:00–5:00. These are different documents from the same organizers, one hour apart on the second half of the day. **Verify the live schedule the morning of the event** — this changes how much time you budget post-curveball (either ~3.5 hrs or ~2 hrs). Everything in §14's build sequence below is written against the Participant Guide's 12:00–3:00 window since it's the more detailed operating document, but treat the actual kickoff announcement as authoritative.
- **Databricks module upgraded substantially.** Section 12 now specifies a real, research-grounded trust-but-verify pipeline (I-CALM prompting + conformal risk control thresholding + Databricks AI Search hybrid retrieval as a corroboration engine) instead of a simple logistic-regression sketch. This is a materially stronger "innovation" and "meaningful use" story. A lighter fallback tier is kept for time pressure.
- **Curveball intelligence added.** Section 13 now includes ranked, evidence-based predictions of the actual pivot themes Entire is likely to use, sourced from their real open issues/PRs, with a specific architectural safeguard (a repository/adapter data-access layer) that defends against the single highest-confidence candidate.
- **Differentiation vs. Entire's own built-in skills sharpened.** Entire ships `review`, `explainskill`, `what-happened`, and `session-handoff` skills that already do intent-aware, checkpoint-grounded analysis. Fidelity must be positioned precisely against these — see §1.4. This is a question a judge will very plausibly ask; have the answer ready cold.
- **`entire-judge` discovered.** Entire's own hackathon-evaluation tooling appears to programmatically inspect a team's checkpoint history ("judge lenses over a submission's brain"). This raises the stakes on checkpoint quality — see §13.4.
- **Lakebase/LTAP research pass added (v3, ideation-stage, feasibility deliberately ignored per instruction).** A dedicated technical read of Databricks Lakebase (serverless Postgres, reached GA Feb 3 2026) and its LTAP dual OLTP/OLAP architecture surfaced three genuinely new Databricks differentiators, not reframings of v2 material: **zero-copy counterfactual branching** for instant "what-if" verdicts, a **dual-write-free Verdict Ledger** that promotes the roadmap's Plan-Trust Ledger from deferred to potentially-buildable, and a **Lakebase-backed memory pattern for the mandatory noon curveball reconstruction step** — see new §12.10–§12.14. This reopens the Databricks module (§12 only) for revision; §§1–11 (core Entire mechanics) remain untouched and still locked.
  **⚠️ New conflict flagged, unresolved:** sources disagree on whether Lakebase is provisionable under Databricks Free Edition. Databricks' own "What's coming to Free Edition" announcement states Lakebase is now included for every Free Edition user "from day one." A separate Free Edition limitations reference page lists **"Lakebase database instances"** explicitly under *unsupported* features. This is not resolved here — see §12.14 for the required verification step before committing any build time to §12.10–§12.13. Everything else in this document stands regardless of the outcome.
- **Competitive-landscape research pass added (v4).** A broad market/academic survey of graph-augmented coding-agent tooling (Graphify, Serena, Sourcegraph/Cody, CodeQL, Kythe, CodexGraph/RepoGraph, and a 2026 "Thousand-Graph Hypothesis" pre-print) was read for genuinely new differentiators, not just corroboration of what's already here. Two are integrated as new build candidates: **community-aware scope-creep framing** via graph clustering (new §10.1) and **Fidelity as an MCP-exposed self-check tool** other agents can call on themselves (new §11.1). A new §1.5 sharpens Fidelity's position against Graphify — whose EXTRACTED/INFERRED/AMBIGUOUS tagging is close enough to Entire's own tri-state model that a judge may ask about it directly — and states an explicit, reasoned rejection of the Thousand-Graph Hypothesis's latent-materialization paradigm for Fidelity's core reconciliation tier. As with v3, this round is ideation-stage with feasibility deliberately ignored per instruction; §§1–11's core mechanics otherwise stand as locked. No new source conflicts surfaced this round — the deadline conflict (§0/v0) and the Lakebase Free Edition conflict (§0/v3, §12.14) both remain open.
- **⚡ v5 — THE ACTUAL NOON CURVEBALL, received and now binding.** Working assumption, flagged as inference not a stated mapping: the received card is **Track 2 — "Graph Is Evidence, Not an Oracle"** (see §13.5 for the reasoning; confirm with a mentor at the earliest opportunity — the doc does not spell out the track→card-number correspondence explicitly). The gap it exposes is real and distinct from anything already in the tri-state model: `extracted`/`inferred`/`ambiguous` only classifies edges that exist; it says nothing about relationships that were never recorded at all (dynamic dispatch, generated code, reflection — a silent hole, not a flagged guess). Without a fix, a changed entity sitting in such a hole would be reported as confidently reachable-or-not, which is exactly the "incomplete evidence presented as certain" failure the curveball forbids. This revision adds a new **coverage-confidence** signal (§10.2), orthogonal to edge tri-state, threaded through reconciliation (§10.2), the schema (§8.1), config (§8.2), the mandatory-vs-Databricks-enhanced fallback split (§12.4/§12.15), the build sequence (§14), and testing (§15) against the supplied partial-analysis fixture. **Honest note, stated plainly per house style: none of the five ranked pivot predictions in §13.2 called this exact pivot.** Rank 4 (multi-repo DATA_FLOWS tracing, evidence: Issue #32 boundary panics) shares the *spirit* of graceful degradation but not this specific mandatory-disclosure mechanism — the safeguard architecture built for it (defensive wrapping around graph search) is reusable, but the coverage-confidence signal itself is new work, not something the repository/adapter firewall already absorbed for free.
- **⚡ v6 — Open-source model/tool augmentation pass, prioritized against remaining build time.** New §20 surveys concrete open-source models and tools that could strengthen specific pipeline stages, explicitly triaged into three tiers rather than presented as an undifferentiated wishlist: **curveball-critical** (do now — real static-detection evidence for §10.2's `partial` coverage tier, currently only a zero-edge heuristic), **real-but-optional upgrades** to already-cited mechanisms (a proper CRC library for §12.2, the actual Leiden implementation for §10.1, an open-source local corroboration fallback for §12.4/§12.15, a named open-weight model for the §9/§12.2 extraction/I-CALM step), and **explicitly not-now** (MCP SDK exposure for §11.1, CodeQL/Kythe as heavier alternatives to the curveball fix, and — stated plainly — never adopting Graphify itself as a dependency, since §1.5 exists specifically to differentiate against it). This is scoped as a menu against the clock, not a mandate: §20.5 states outright that only the single top-tier item is recommended if time allows for just one addition. No change to §§1–11's locked core mechanics; this bolts onto Stage 4 (§10.2) and the Databricks module (§12) exactly where those sections already declared open extension points.

---

## 1. Why Fidelity Is the Core

### 1.1 What `entire-graph` actually is (verified from source/PR analysis)

- **Confirmed edge types:** `CALLS`, `CONTAINS`, `DATA_FLOWS`, and type/data-flow relations across ~10 languages — actively developed (PRs #148, #161, #163, #168, #174, #177, #184).
- **Unverified in the core engine:** `TESTS`, `HANDLES_ROUTE`, `HANDLES_GRPC`. Do not design a feature that depends on these existing.
- **Edge confidence is a tri-state**, not a float or boolean: `extracted` (deterministic, from explicit syntax — trust as fact), `inferred` (heuristic — validate before destructive use), `ambiguous` (multiple valid targets — present as possible, never definitive). The engine's own benchmark philosophy counts **honest ambiguity as correct** — it rewards graceful degradation over confident hallucination. Fidelity's whole design mirrors this.
- **Cross-repository graph stitching is explicitly not supported** (PR #169: "the graph is scoped to a single primary store"). **Correction from the newer research:** a *narrower* cross-repo capability does exist one layer up — `entire explain --repo` can drill into a specific foreign repo's checkpoints over HTTP for session/search purposes — but this is checkpoint/session retrieval, not structural graph stitching. It does not change the verdict on fleet-style graph analytics (§16).
- **O(N²) relation-building phase** with a wall-clock ceiling (PR #142) — no guaranteed real-time on huge commits.
- **Freshness is explicit**: queries use the last indexed commit; working-tree changes are never silently folded in.
- **Partial analysis is diagnostic, not authoritative** — only fully indexed commits are trustworthy.
- **Known fragility:** Issue #32 documents a slice-bounds panic (`[-1:0]`) in graph search on certain malformed/boundary syntax during DATA_FLOWS expansion — a real, currently-open bug. Fidelity should catch and degrade gracefully here rather than crash (this is free credit on "failures handled safely").
- **Agent-payload record forgery is a known risk** (PR #143) — the graph quarantines strings that mimic graph records because it parses agent-written code. Our own LLM extraction step must fail closed for the same reason.
- **Checkpoint IDs are ULIDs stored as Git commit trailers**, binding session metadata directly to commits — this is the temporal/ordering backbone the whole product relies on.
- **The RFD 0006 / SCIP export effort** (PRs #152, #153) signals Entire is standardizing the graph's export format for cross-tool interoperability — worth a one-line mention in the roadmap as a natural place for `verdict.json` to plug into later.

### 1.2 What Genie Ontology is NOT (unchanged from v1, still holds)

Genie Ontology is a closed, auto-inferred system: no direct graph write API, no Cypher/Gremlin, natural-language queries resolve to **SQL** on a SQL warehouse (so deep multi-hop traversal depends on an LLM writing flawless recursive CTEs — it won't), and each agent is capped at **200 knowledge-store snippets** shared across tables/joins/metrics. A code graph's schema blows that instantly. **Decision, unchanged: do not force the code graph into Genie Ontology.** State this trade-off to the Databricks judge explicitly — it reads as platform fluency, not a gap.

### 1.3 What this kills, keeps, and defers (unchanged verdicts from v1)

| Prior concept | Verdict | Reason |
|---|---|---|
| **Fidelity / Intent-vs-Implementation** | **CORE** | Single-repo, single-session; uses only confirmed mechanics + the real tri-state model |
| Edge-confidence calibration | **KEPT, substantially upgraded** | Now the full I-CALM/CRC/AI-Search design — see §12 |
| Plan-Trust Ledger, Forensic Replay, Checkpoint Court, Subagent Reasoning Capture, Compliance-lineage reconciliation | **Deferred → roadmap** | Real extensions of the same core; see §16 |
| Fleet Intelligence, Ghost Dependency, Graph Weather, Org-Topology Mirror | **CUT** | Cross-repo graph stitching is unsupported; edges relied on are unverified |
| Genie Ontology as the graph store | **CUT** | See §1.2 |

### 1.4 Differentiation from Entire's own built-in skills (new — critical for the pitch)

Entire already ships skills that sound adjacent to Fidelity, and a judge who knows the ecosystem may ask about this directly:

- **`review`** — "audits code changes on the current branch by reading checkpoint transcripts to understand developer intent, producing intent-aware findings."
- **`explainskill`** — looks up the session behind a function or line to explain why it exists.
- **`what-happened`** — traces the latest change with git blame + checkpoint context to debug regressions.

**The precise distinction to hold onto:** these are all *retrieval-and-narration* tools — they answer "why is this code the way it is" by surfacing and summarizing session context, largely through an LLM reading a transcript. **Fidelity answers a different question — "did the change match what was claimed" — and answers it with a deterministic, tiered, graph-verified reconciliation, not an LLM's narrative judgment.** Every tier in `verdict.json` is backed by an actual graph reachability computation over classified edges, not a semantic impression. `review` could *use* Fidelity's `verdict.json` as a structured input to make its own narrative sharper — they're complementary layers, not competitors. Have this framing ready verbatim; it's the single most likely "gotcha" question in the room.

### 1.5 NEW (v4) — Competitive Landscape & Positioning Beyond Entire's Own Skills

Fidelity's design should also be read against the wider graph-augmented coding-agent tooling landscape, not only Entire's own skills (§1.4):

| Tool | Extraction | Confidence model | What it actually does | How Fidelity differs |
|---|---|---|---|---|
| **Graphify** | Tree-sitter (code) + LLM (docs/PDF/video), multimodal | Explicit **EXTRACTED / INFERRED / AMBIGUOUS** tri-state tagging | Builds a persistent, queryable multimodal knowledge graph for agent context/exploration | **The tri-state idea is not Fidelity's differentiator — it's converging on the same taxonomy independently.** Fidelity's differentiator is what it *does* with the signal: a deterministic verdict reconciling stated intent against structural change, not a context-retrieval aid. Likely judge question — verbatim answer below. |
| **Serena** | Native LSP (compiler-grade) | Absolute, compiler-derived precision | Precision code editing/symbol resolution via MCP tools | Higher raw precision than tree-sitter, but no checkpoint/intent linkage at all — it edits correctly, it verifies *nothing* against stated intent |
| **Sourcegraph/Cody** | Hybrid static analysis + embeddings | Implicit, via search ranking | Code search and assistant context | A retrieval tool, not a verification/reconciliation tool |
| **CodexGraph/RepoGraph** | LLM-written Cypher over a Neo4j schema | N/A (agent trusts its own queries) | Lets the agent query its own targeted subgraph for debugging/localization | Agent-directed retrieval, not an independent, deterministic auditor of what the agent already claimed and did |

**Ready-verbatim answer if asked "isn't this just Graphify's tri-state model?"** — "Graphify's tagging routes an agent's *attention* during exploration; ours drives a deterministic *verdict* after the fact. Tagging an edge INFERRED helps an agent decide how much to trust it while working. Fidelity uses the same kind of signal to decide, after the work is done, whether a specific claim holds — that's reconciliation, not retrieval."

**Explicit rejection of the "Thousand-Graph Hypothesis" (2026 pre-print) for Fidelity's core tier:** this research argues persistent graph edges should be abandoned in favor of an LLM's attention mechanism latently materializing task-specific relationships from raw entities at inference time. Fidelity's entire value proposition is the opposite bet — a verdict a reviewer can trust *because* it rests on persisted, classified, non-hallucinated edges, not a fresh, unreproducible LLM inference per query. Worth stating this explicitly if a technically literate judge raises it: Fidelity isn't unaware of the paradigm, it deliberately rejects it for an audit-trail product — though the extraction step's LLM fallback (§9) is itself a small, tightly fenced concession to this same kind of task-conditioned inference, and is fail-closed for exactly that reason.

**Two named, unsolved industry problems (per this landscape research) that validate design choices already in this bible, not new features:**
- *"The field currently lacks a standardized mathematical framework for bounding the error rate of LLM-inferred graph edges"* — this is precisely the gap §12.2's Conformal Risk Control step exists to close, and is worth citing to the Databricks judge as grounding beyond Entire's own tri-state model.
- Dynamic/late-bound language resolution (duck typing, monkey-patching, dynamic dispatch) is called out as a genuinely unsolved problem for static parsers — this is exactly why `ambiguous` is a first-class, never-silently-resolved tier in Fidelity's reconciliation (§10), not an edge case to explain away.

**Caveat on cited figures:** Graphify's own headline claims (e.g. a reported token-reduction multiple) are self-reported and explicitly flagged as unverified even in the source survey — nothing here repeats those figures as fact; only the taxonomy and market positioning are used.

---

## 2. Problem & User

**Who:** the person who must approve an AI agent's work before it merges — a tech lead reviewing a teammate's agent session, or an engineer reviewing their own agent's output before pushing.

**Problem:** agents are fast, but reviewers have no reliable way to know if the agent did *exactly* what was asked, *more* (silent scope creep), or *less* (task left half-done). Today the only answer is re-reading the whole diff by hand, or trusting an LLM's narrative summary of the transcript (see §1.4 — that's what the existing skills give you, and it's judgment, not verification).

**Why it needs Entire specifically:** Checkpoints are the only ground-truth stated intent tied to a specific atomic change. Graph diff is the only structural ground truth. Neither alone solves it.

---

## 3. Operational Compliance (from the Participant Guide — verbatim commands and deliverables)

### 3.1 Required setup sequence (run after the official 9:00 AM start, not before)

```bash
entire login
entire repo mirror create
entire repo clone /gh/YOUR-GITHUB-HANDLE/REPOSITORY
```
Select your fork and the **India region** when prompted. Work only inside the clone created through this mirror flow.

```bash
entire enable -y --agent YOUR-BUILD-AGENT
entire status
```
Enable Checkpoints **before** any agent-assisted development begins.

```bash
entire plugin install graph
entire graph version
entire graph init-agents --repo .
```
**Start a fresh agent session immediately after this** so the agent actually receives the graph instructions — this is an explicit requirement, not a suggestion.

### 3.2 Required checkpoint milestones (exactly four, per the guide)

1. Initial understanding and intended architecture.
2. The last stable state before the Noon Curveball.
3. Your response to the Noon Curveball.
4. Final implementation and verification.

Checkpoint **quality** is explicitly weighted over quantity: capture decisions, rejected options, failures, assumptions, open risks, and the evidence that changed your approach — not just diffs.

### 3.3 Required graph demonstrations (exactly three, per the guide) — mapped onto Fidelity's own stages

| Guide's requirement | Where Fidelity satisfies it |
|---|---|
| A graph search or definition lookup | Stage 1/2 — loading the full symbol dictionary via `entire graph snapshot` to ground the LLM extraction vocabulary |
| A relationship or impact analysis before a high-risk change | Stage 4 — `graph neighbors`/`impact` N-hop reachability, run *before* tiering any change as safe |
| A final semantic-diff analysis of the submitted implementation | Stage 3/5 — `entire graph diff --base <base> --head <checkpoint-id>` feeding the final `verdict.json` |

Make sure the demo explicitly narrates these three moments by name — the rubric line literally asks for them.

### 3.4 Databricks Free Edition constraints (design around these explicitly if opting in)

| Area | Constraint | Fidelity's fit |
|---|---|---|
| SQL/jobs | One serverless 2X-Small warehouse; ≤5 concurrent job tasks | Trivial — our calibration job is a small batch write, not a cluster workload |
| Models/AI Search | No GPU/provisioned throughput; AI Search is limited | Our calibration model is a **logistic regression** by design (§12) — never needed GPU. AI Search corroboration queries should be treated as a capped, occasional call, not a per-edge hot path |
| Apps/database | Up to 3 Apps, idle auto-stop; one Lakebase project where available | Fits a single dashboard app + a scoring endpoint comfortably |
| Workspace | One workspace, one metastore per account | Plan the shared final demo workspace in advance; name one deployment owner (per the guide's own team-readiness requirement) |
| Quota | Exhausted quota can take out compute for the rest of the day | **Prefer the batch-score fallback over a live Serving endpoint during the demo** — this was already the config default (`score_source: batch`) in v1, and this constraint is the concrete reason why. Test the calibration path early in the morning, not near noon. |

### 3.5 Submission package (exact fields, from the guide)

Selected track; GitHub fork URL + final commit SHA; Entire mirror/project URL; links to the four required Checkpoints; setup/run/test instructions; working demo or a reliable fallback recording; a complete `BUILDATHON.md` at repo root; **if opting into Databricks:** capabilities used + why essential, workspace/app/endpoint URL, relevant repo paths, reproduction steps, data provenance, and how the Curveball affected the Databricks workflow.

---

## 4. Rubric Alignment (confirmed identical point weights across both source documents)

### Entire (100 pts)
| Line | Pts | How Fidelity earns it |
|---|---|---|
| Problem & innovation | 20 | Named, low-saturation gap; explicitly differentiated from Entire's own `review`/`explainskill` (§1.4) |
| Technical implementation | 25 | Real subcommand, working pipeline, meaningful tests, graceful handling of known bugs (Issue #32) |
| Response to Curveball | 15 | ⚡ v5: the actual received constraint (§13.5) is answered concretely — coverage-confidence signal (§10.2), new `unverifiable_coverage` tier, mandatory fallback path, three-way renderer split (§11.1), pre-edit `graph impact` run (scored explicitly), fixture test (§15); honestly noted where prediction (§13.2) fell short rather than retrofitted to look foreseen |
| Use of Checkpoints | 15 | Transcript **is** the intent input; verdict is itself wrapped in a Checkpoint; all 4 required milestones captured with real decision content |
| Use of Graph | 15 | All 3 required demonstrations are load-bearing, not decorative (§3.3) |
| Demo & future potential | 10 | Three renderers of one manifest; explicit roadmap (§16); if time permits, an MCP-exposed self-check (§11.1) shows the same manifest already has a second, agent-facing consumer, not just human dashboards |

### Best Use of Databricks (100 pts, opt-in, scored separately)
| Line | Pts | How the calibration layer earns it |
|---|---|---|
| Meaningful use | 30 | Removing it reverts a calibrated, corroborated confidence signal to a static advisory list — measurable functional loss. **v3, if §12.14 clears:** removing it also deletes the cross-session Verdict Ledger (§12.11) and the structured checkpoint-memory the fresh session reconstructs from (§12.13) — a stronger, more literal loss than a confidence score alone |
| Working implementation & reliability | 25 | Delta table + MLflow run + AI Search corroboration call, reproducible, with a documented batch fallback |
| User value & product decisions | 20 | Reviewer sees exactly which advisory edges have been independently corroborated vs. merely flagged |
| Data quality, provenance, responsible use | 15 | Synthetic/injected calibration labels **documented honestly**; AI Search index built only from the repo's own docs/comments/commit messages, governed by Unity Catalog |
| Curveball response | 10 | ⚡ v5: §12.15's coverage-corroboration reuses the existing §12.4 AI Search index — genuinely low-marginal-cost, and explicitly framed as an enhancement of a fallback that already works without Databricks, not a dependency introduced by the curveball |

---

## 5. Pre-Flight Verification Checklist (run in the first setup window)

- [ ] `entire graph capabilities --json` → confirm the **actual edge-type enum**. Do not assume `TESTS`/`HANDLES_ROUTE`/`HANDLES_GRPC` exist.
- [ ] Capture one real `entire graph snapshot` → confirm node/edge fields and the **exact key** for tri-state classification.
- [ ] Capture one real `entire graph diff --base <ref> --head <ref>` → confirm the output shape.
- [ ] Check `entire graph impact --help` / `entire graph neighbors --help` for depth/direction flags. If absent, implement N-hop reachability by iterating `neighbors` ourselves.
- [ ] Confirm the checkpoint transcript read path, and that both ULID and legacy hex checkpoint ID formats are handled.
- [ ] Sanity-test a deliberately malformed/edge-case file against `graph search` to see if Issue #32's panic reproduces locally — if so, wrap the call defensively from the start.
- [ ] ⚡ **v5 — Curveball item.** Check whether `entire graph capabilities --json` or `entire graph snapshot` exposes any native "partial/unresolved" or "unindexed region" marker before assuming the zero-edge heuristic (§10.2, `coverage.detect_via`) is the only option. Record whichever is real at the top of `BUILDATHON.md` alongside the other pre-flight findings.

Record findings at the top of `BUILDATHON.md`.

---

## 6. Architecture

```
┌─────────────────┐   ┌──────────────────┐   ┌────────────────────┐
│ 1. CAPTURE      │   │ 2. DECLARE        │   │ 3. OBSERVE          │
│ Read checkpoint │──▶│ Extract intended  │   │ entire graph diff   │
│ transcript      │   │ entities          │   │ + snapshot edges    │
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

Each arrow is a clean interface boundary — the curveball firewall (§13). All renderers are pure functions of `verdict.json`.

---

## 7. Command Surface

```
entire graph verify-intent <checkpoint-id> \
    [--base <checkpoint-id>] \
    [--config fidelity.config.yaml] \
    [--out ./fidelity-out/]
```

Internal calls (exact syntax pinned after §5's verification): checkpoint transcript read → intent; `entire graph snapshot` → symbol dictionary + tri-state edges; `entire graph diff --base <base> --head <checkpoint-id>` → changed entities; `entire graph neighbors`/`impact` → N-hop reachability.

---

## 8. Data Schemas

### 8.1 `verdict.json`

```json
{
  "checkpoint_id": "01K9TQ8ZP7X3F5M2WVJ4CNRB6D",
  "base_checkpoint_id": "01K9TQ8YAB12CDEF34GH5IJK6L",
  "generated_at": "2026-09-06T10:42:00Z",
  "graph_capabilities_verified": ["CALLS", "CONTAINS", "DATA_FLOWS", "type"],
  "intent": {
    "raw_prompt_excerpt": "Fix the timeout bug in AuthRouter...",
    "extracted_entities": [
      { "name": "AuthRouter", "method": "regex", "confidence": 1.0 },
      { "name": "TokenValidator", "method": "llm_fallback", "confidence": 0.62 }
    ],
    "unresolved_fragments": ["the login flow"]
  },
  "reconciliation": {
    "confirmed": [
      { "entity": "AuthRouter.handleTimeout", "file": "src/auth/router.go", "line": 88 }
    ],
    "declared_unimplemented": [
      { "entity": "TokenValidator", "reason": "named in prompt, no graph change detected" }
    ],
    "expected_blast_radius": [
      { "entity": "SessionCache.invalidate", "file": "src/auth/cache.go", "line": 40,
        "path_from_confirmed": ["AuthRouter.handleTimeout", "SessionCache.invalidate"],
        "hops": 1, "edge_types_traversed": ["CALLS"], "min_edge_class": "extracted" }
    ],
    "undeclared_scope_creep": [
      { "entity": "BillingController.chargeRetry", "file": "src/billing/controller.go", "line": 210,
        "reason": "changed, not reachable from any confirmed entity within 2 hops",
        "coverage_confidence": "full" }
    ],
    "advisory_low_confidence": [
      { "entity": "PaymentGateway.route", "edge_class": "inferred",
        "reason": "blast-radius reasoning here rests on a non-deterministic edge",
        "coverage_confidence": "full",
        "verbal_confidence": null,
        "conformal_gate_passed": null,
        "calibrated_confidence": null,
        "corroboration": { "attempted": false, "verified": null, "evidence_snippet": null }
      }
    ],
    "unverifiable_coverage": [
      { "entity": "PluginRegistry.dispatch", "file": "src/plugins/registry.go", "line": 57,
        "coverage_confidence": "partial",
        "coverage_reason": "zero extracted-class edges detected; entity resolved via reflection-style dynamic dispatch that static analysis cannot fully resolve",
        "verification_path": "manual_review_recommended",
        "corroboration": { "attempted": false, "verified": null, "evidence_snippet": null }
      }
    ]
  },
  "summary": {
    "total_changed_entities": 5,
    "scope_creep_count": 1,
    "unimplemented_count": 1,
    "advisory_edge_count": 1,
    "unverifiable_coverage_count": 1,
    "verdict_label": "REVIEW_REQUIRED"
  }
}
```

**New fields vs v1:** `verbal_confidence`, `conformal_gate_passed`, `corroboration` — populated by the upgraded Databricks module (§12) when enabled; `null` when it's off, so the schema is stable either way.

**⚡ New fields vs v4 (Curveball response, §10.2):** `coverage_confidence` (`full`/`partial`/`unknown`) now appears on **every** entity in **every** tier, including `confirmed` and `expected_blast_radius` (omitted above only for brevity — it must be populated on all six tiers in the real implementation). A new sixth reconciliation tier, `unverifiable_coverage`, holds entities that would otherwise have fallen into `undeclared_scope_creep` but sit in a `partial`/`unknown` coverage region — carrying `coverage_reason` and `verification_path` (`"manual_review_recommended"` baseline, or a populated `corroboration` block if the Databricks enhancement in §12.15 is enabled). `summary.unverifiable_coverage_count` is new alongside it. **Never conflate this tier with `advisory_low_confidence`** — the latter is "the graph has an edge but isn't sure about it," the former is "the graph may not have tried to record an edge here at all."

### 8.2 `fidelity.config.yaml`

```yaml
blast_radius_hops: 2
edge_classes:
  deterministic: [extracted]
  advisory: [inferred, ambiguous]
tiers:
  confirmed: {}
  declared_unimplemented: {}
  expected_blast_radius: { max_hops: 2 }
  undeclared_scope_creep: {}
  advisory_low_confidence: {}
extraction:
  fallback_llm: true
  llm_vocabulary_source: graph_snapshot
module_granularity: true
verdict_thresholds:
  review_required_if_scope_creep_gte: 1
  review_required_if_unverifiable_coverage_gte: 1   # NEW v5 — Curveball requirement, on by default
coverage:                       # NEW v5 — §10.2, orthogonal to edge_classes above
  detect_via: capabilities_api  # capabilities_api | heuristic_zero_edge — verify which is real at kickoff (§5)
  treat_zero_edge_as: partial   # partial | full — conservative default per curveball ("never present incomplete as certain")
  fallback_verification_path: manual_review_recommended
databricks:
  calibration_enabled: false
  score_source: batch          # batch | serving_endpoint — prefer batch (see §3.4 quota risk)
  gate:
    mode: fallback              # fallback | icalm_crc — see §12
    alpha_target: 0.05          # only used in icalm_crc mode
  ai_search:
    enabled: false
    index_name: null            # e.g. engineering_catalog.code_graph.repo_context_index
  coverage_corroboration_enabled: false   # NEW v5 — §12.15, reuses ai_search index for unverifiable_coverage entities
```

---

## 9. Extraction Logic (Stage 2) — unchanged from v1

1. Load the full symbol dictionary from `entire graph snapshot`.
2. Fuzzy-match transcript tokens against it.
3. For unresolved fragments, send **only those fragments + the constrained symbol dictionary** to an LLM with an explicit instruction to pick from the exact list or return "no confident match." **Fail closed** — never invent an entity (this also respects the graph's own anti-forgery posture, PR #143).
4. Every extracted entity carries `method` and `confidence`.

Vague prompts ("continue", "fix it") inherit intent from the nearest resolvable ancestor checkpoint, tagged `inherited_intent: true`. No-baseline (first checkpoint) diffs against an empty tree, tagged `no_baseline: true`.

---

## 10. Reconciliation Logic (Stage 4) — core unchanged from v1; ⚡ v5 adds §10.2 (mandatory Curveball response)

`confirmed` / `declared_unimplemented` / `expected_blast_radius` (via deterministic-edge N-hop reachability, recording `edge_types_traversed` and `min_edge_class`) / `undeclared_scope_creep` / `advisory_low_confidence` (anything that would only hold via an `inferred`/`ambiguous` edge — never silently promoted to "expected"). If `impact`/`neighbors` lack depth flags, compute N-hop reachability by iterating single-hop `neighbors`.

### 10.1 NEW (v4) — Community-Aware Scope-Creep Framing (optional enhancement)

Raw reachability tiering above tells a reviewer *which* entities changed and how many hops away they are — it doesn't say what part of the system they belong to. Applying a community-detection pass (Leiden-style clustering, the current standard for surfacing architectural boundaries that folder structure often hides) over the snapshot graph lets `undeclared_scope_creep` entries be reported by subsystem, not just by file — e.g. "this change also touched 3 entities in the **Billing** community" rather than three unrelated-looking file paths. This turns a flat list into an architectural statement a reviewer can act on faster, directly answering the Graph track's requirement that raw graph output "produce a useful decision... not raw graph output."

**Honest caveat:** this needs a graph large and connected enough for clustering to produce meaningful, non-trivial communities — a small hackathon demo repo may not have one. Treat this as a stretch enhancement to attempt only once the base pipeline (§6–11) is solid, and be ready to state the caveat plainly if the demo repo is too small to show it convincingly, rather than forcing a clustering result that isn't real.

### 10.2 ⚡ NEW (v5) — Coverage-Confidence: the mandatory Curveball response (Track 2)

This is orthogonal to the existing edge tri-state and must not be conflated with it. Edge tri-state (`extracted`/`inferred`/`ambiguous`) classifies confidence in a relationship **that the graph recorded**. Coverage-confidence classifies confidence in whether the graph's silence about an entity's relationships **means anything at all**.

**The signal — attached to every entity that reconciliation touches, not just changed ones:**

| Coverage value | Meaning | Detection (heuristic, pending live verification) |
|---|---|---|
| `full` | Entity's relationships are backed by normal static extraction; a zero-hop result is a real "not connected" | Default state for any entity with at least one `extracted`-class edge of any kind (incoming or outgoing) in the snapshot |
| `partial` | Entity sits in a region where static analysis is known to degrade — dynamic dispatch, generated code, reflection, or any zero-edge entity with no `extracted` edges at all despite non-trivial code presence | **Needs live verification (add to §5's pre-flight checklist):** check whether `entire graph capabilities --json` or `entire graph snapshot` exposes a native "partial/unresolved" marker first. If not, use the heuristic proxy: an entity with zero edges of any class is *not* automatically `full`-coverage-confirmed-empty; it is provisionally `partial` unless a stronger signal (e.g. the entity is a leaf data type with no calls by design) rules that out. **⚡ v6 upgrade path (§20.1):** run Semgrep's dynamic-dispatch/reflection/codegen rule packs against the same entity before falling back to the zero-edge heuristic alone — a real pattern match promotes `partial` from an inferred proxy to a detected finding, which is a materially stronger answer if a judge asks "how do you actually know this is dynamic dispatch and not just disconnected code?" |
| `unknown` | Coverage itself could not be determined this run (e.g. `graph snapshot` failed or partially failed for that region) | Any snapshot/diff call that errors, times out, or returns a partial result for the entity's file/module |

**The rule that actually satisfies the curveball — never silently upgrade:** Stage 4 reconciliation must check coverage-confidence *before* assigning a tier, not after. Concretely:
- An entity with `coverage: partial` or `coverage: unknown` **can never be assigned to `undeclared_scope_creep` on the strength of an absent path alone.** It is instead routed to a new tier, `unverifiable_coverage` (see §8.1), carrying the same file/line/reason data `undeclared_scope_creep` would have, plus the coverage reason.
- An entity with `coverage: full` proceeds through the existing tiers exactly as in v1–v4 — **this preserves existing behavior for fully-resolved code, unchanged, per the curveball's explicit requirement.**
- `confirmed` and `expected_blast_radius` entries are unaffected in the success case (a real `extracted`-edge path found is real regardless of what else is unresolved nearby) but must still carry their coverage value in the schema so a reviewer can see *how much of the surrounding graph* the confirmation rests on.

**Safe fallback / verification path (mandatory, must work with or without Databricks):**
- **Baseline (always present, no Databricks required):** any `unverifiable_coverage` entry is flagged `verification_path: "manual_review_recommended"` in the schema, surfaced distinctly in every renderer (§11) — never merged visually with `undeclared_scope_creep`, never given a severity color that implies certainty.
- **Enhanced (Databricks opted-in):** the existing §12.4 AI Search corroboration step is the natural upgrade path — the same corroboration query mechanism can be pointed at `unverifiable_coverage` entities instead of (or in addition to) `advisory_low_confidence` edges, since both are "the deterministic path can't resolve this, ask the docs/comments/commits instead." See §12.15 for the concrete wiring.

**Interface requirement (per the curveball, verbatim):** all three renderers must let a user visually distinguish, at a glance, three states — confirmed structural evidence, heuristic/incomplete evidence (the existing `inferred`/`ambiguous` tri-state), and now claims that require source/test verification (`unverifiable_coverage`). Three visually distinct states, not two with a footnote.

**What does not change:** Stage 1–3 extraction and Stage 5 rendering's *pure-function-of-`verdict.json`* property both hold — this is a Stage 4 reconciliation-logic change plus a schema addition, not a rewrite of the pipeline. The §13.1 structural firewall holds: this is squarely a "what counts as drift" change, which the firewall already scoped to Stage 4 + config.

---

## 11. Rendering (Stage 5) — unchanged from v1

Terminal report (colored, tiered, `file:line`) → `VERIFY_REPORT.md` (linked from `BUILDATHON.md`) → `dashboard.html` (single static file, `fetch()`es `verdict.json`, radial/force graph color-coded by edge class). All pure functions of `verdict.json`. Dashboard is the screen open during judging.

**⚡ v5 requirement, binding on all three renderers:** each must visually separate `unverifiable_coverage` (§10.2) from `undeclared_scope_creep` and from `advisory_low_confidence` — three distinct visual treatments, not a shared "risky" bucket with a tooltip. This is a direct, scored requirement from the Track 2 curveball card, not a nice-to-have.

### 11.1 NEW (v4) — Fidelity as an MCP-exposed self-check tool for other agents

The competitive-landscape research confirms MCP is consolidating as the standard interface for agentic tool access (Graphify and Serena both ship MCP servers as first-class citizens, not an afterthought). Exposing `verify-intent` and `verdict.json` querying as MCP tools — not just a CLI + dashboard for humans — lets another coding agent call Fidelity on itself *before* ever presenting a diff to a human reviewer: "run your own verification, and don't hand this back to me unless every tier is `confirmed` or `expected_blast_radius`." This reframes Fidelity from a purely human-facing audit tool into something an agent can use to self-gate its own output, and it's a low-effort extension once `verdict.json` is stable (§6–8) — the MCP tool surface is a thin wrapper over the same command already built.

**Ties into an existing thread:** this is a good present-day complement to the RFD 0006 SCIP-export mention already in §1.1/§16 item 6 — both are about making `verdict.json` consumable by something other than a human reading a dashboard.

---

## 12. Databricks Module — Calibrated Trust-But-Verify Gate (substantially upgraded)

**Attempt only after §6–11 work end-to-end. Bolts onto Stage 5's output; never touches Stages 1–4's interfaces.**

### 12.1 The idea, in one paragraph

The tri-state edge model tells us an edge is `inferred` or `ambiguous`, but not *how* uncertain, and gives no path to independent corroboration. This module adds three real, cited research mechanisms on top: **I-CALM** for eliciting a calibrated self-reported confidence at extraction time, **Conformal Risk Control (CRC)** for turning that into a mathematically bounded accept/reject threshold, and **Databricks AI Search hybrid retrieval** as the corroboration engine for whatever the gate can't resolve on its own.

### 12.2 Phase 1 — Offline calibration (before the demo runs)

Compile a small calibration set of known true/false heuristic edges (hand-labeled or injected). Run the extraction prompt with an explicit **I-CALM** structure: the model is told, in the prompt, that a correct edge earns full reward, an incorrect one is heavily penalized, and abstaining ("no confident match") earns partial credit — plus a short normative instruction not to assume unverified relationships. The model returns each proposed edge with a **verbal confidence score** in `[0,1]`. Using **Conformal Risk Control**, compute the smallest threshold `λ̂` such that the *expected false-positive rate* on the calibration set stays below a target `α` (e.g. 0.05). This gives a threshold with an actual statistical guarantee, not an arbitrary cutoff.

### 12.3 Phase 2 — Online gating (per advisory edge, at verify-intent time)

- **Deterministic edges bypass the gate entirely** — always trusted.
- **Advisory (`inferred`/`ambiguous`) edges** get the I-CALM verbal-confidence score.
  - If `verbal_confidence ≥ λ̂`: gate passes; edge is used in reconciliation with `conformal_gate_passed: true`.
  - If `verbal_confidence < λ̂`: **epistemic abstention** — the edge is held `pending_verification` and routed to corroboration (§12.4) rather than silently trusted or silently dropped.

*(Optional, time-permitting extra rigor: for a small set of especially critical edges, sample the extraction prompt N=5 times at high temperature and compute semantic entropy over the resulting (source, target) tuples — exact-match clustering is enough since output is structured JSON, no NLI needed — as a second signal alongside verbal confidence. Skip this tier under time pressure; the single-pass I-CALM score alone is a legitimate, citable mechanism on its own.)*

### 12.4 Corroboration via Databricks AI Search

For every `pending_verification` edge, query a Databricks AI Search index built over the repo's own docstrings, inline comments, and commit messages (landed in Delta, indexed via AI Search — vector + full-text + hybrid via Reciprocal Rank Fusion, one API call with `query_type="hybrid"`). The query is simply the edge's hypothesis in natural language ("does X trigger/route to/call Y?"). The retrieved, RRF-ranked snippets are shown to the reviewer (and can be evaluated by a lightweight LLM check) to either **confirm** the edge (`corroboration.verified: true`, with the evidence snippet attached) or **reject** it (discarded from reconciliation, logged). Unity Catalog governs the index automatically — any access control on the source Delta tables carries through to retrieval, which is worth stating explicitly for the "responsible use" rubric line.

### 12.5 What lands back in `verdict.json`

Every `advisory_low_confidence` entry gets populated: `verbal_confidence`, `conformal_gate_passed`, and — for anything that went to corroboration — `corroboration.verified` plus the evidence snippet. The dashboard renderer ranks advisory edges by this real signal instead of showing an undifferentiated flat list.

### 12.6 Free Edition design fit (see §3.4)

This design was chosen specifically because it fits the Free Edition's real limits: the model is a lightweight scoring mechanism (no GPU need), AI Search calls are capped and occasional (not a per-request hot path — only `pending_verification` edges reach it), and the batch-score default avoids exposing the live demo to Serving-endpoint quota risk.

### 12.7 Fallback tier (if time runs out before the full chain works)

Skip CRC and AI Search. Train a plain logistic regression in MLflow on the same calibration features (edge class, edge type, language, churn size) against outcome labels, and write a single `calibrated_confidence` float back into each advisory edge. This is the v1 design — still legitimate, still qualifies as meaningful Databricks use, just less differentiated. Set `gate.mode: fallback` in the config to signal which tier is live; the schema doesn't change either way (§8.1).

### 12.8 Honest limitations to state in `BUILDATHON.md`

Calibration labels are synthetic/injected, not real production revert data, given the time window. The AI Search index is built only from the demo repo's own docs/comments — a real deployment would need a much larger, curated corpus. This is a demonstration of the *mechanism*, not a production-calibrated system.

### 12.9 Why not Genie Ontology — see §1.2. Keep this answer ready verbatim for the Databricks judge.

### 12.10 NEW (v3) — Zero-Copy Counterfactual Sandbox (Lakebase branching)

Lakebase's copy-on-write branching (seconds to create, near-zero storage cost until modified, fully isolated, its own connection string) enables two genuinely new capabilities, not reframings of anything in v1/v2:

1. **Verification isolation.** Run the entire Stage 3–4 reconciliation against a disposable branch, never the canonical store — a bug in reconciliation logic can never corrupt live graph state. Discard the branch on completion regardless of outcome.
2. **What-if verdicts (new, judge-facing feature).** Let a reviewer replay the *same* diff against the branch with a different `fidelity.config.yaml` — e.g. `blast_radius_hops: 3` instead of `2`, or a different `edge_classes` mapping — and get a second `verdict.json` in seconds, without re-running graph extraction. This is directly demonstrable live during judging: "watch the verdict change if we widen the blast radius."

**Honest trade-off, not smoothed over:** this requires materializing the dependency graph itself inside Lakebase/Postgres (recursive CTEs supplementing or replacing `entire graph neighbors`/`impact` calls), which is a change to Stage 3's architecture, not a Stage 6 bolt-on. That's a bigger integration than the "never touches Stages 1–4" firewall promised in §13.1 — worth naming to a judge as an ambitious, deliberate exception if pursued, not hidden.

### 12.11 NEW (v3) — Live Verdict Ledger (promotes Roadmap §16 item 1 from deferred to potentially-buildable)

LTAP's core mechanic — write once to Lakebase (OLTP), get an always-fresh Delta/Parquet analytical copy automatically, no CDC pipeline engineering — removes the reason §16 originally deferred the "Plan-Trust Ledger" to the roadmap. Every `verdict.json` produced during the day gets an additional row written to a `verdicts` table in Lakebase; LTAP transcodes it into Delta with no extra work. Live at the buildathon, this gives:

- Cross-session trend queries ("this agent's `undeclared_scope_creep` rate over its last N verdicts") with zero pipeline engineering.
- A more literal answer to "why does this need Databricks" than v2's single confidence gate: removing Databricks now removes not just an advisory score but the entire longitudinal trust history — see the strengthened §4 rubric line.

**Depends on the same Lakebase-availability question as §12.10 — see §12.14.**

### 12.12 NEW (v3) — Structural + semantic corroboration collapse (`lakebase_text` / `lakebase_vector`) — optional, not a forced replacement

`lakebase_text` (BM25) and `lakebase_vector` (ANN) are Postgres extensions living inside Lakebase itself, meaning the §12.4 corroboration retrieval could in principle run in the same engine and the same query as graph traversal — one round trip instead of two systems. **Stated honestly as a trade-off, not a strict win:** Databricks AI Search's hybrid retrieval (vector + full-text + Reciprocal Rank Fusion in one call, `query_type="hybrid"`) is a more mature, purpose-built primitive than manually fusing `lakebase_text` and `lakebase_vector` results yourself. If Lakebase is already in the stack for §12.10/12.11, this is worth mentioning to a judge as a "could simplify to one engine" option — not a recommendation to replace §12.4's design.

### 12.13 NEW (v3) — Lakebase-backed curveball reconstruction memory

Databricks documents an "agent state and memory" reference pattern for Lakebase: short-term session memory plus long-term cross-session memory, with LangGraph "time travel" to fork or resume from a stored checkpoint. This maps closely onto the Noon Curveball's mandatory step — "start a fresh session and reconstruct the project from checkpoint context before editing" (Participant Guide, 12:00 PM). Concretely: persist a structured, queryable snapshot of each Entire Checkpoint's key decisions into Lakebase, so the fresh-session reconstruction step pulls from a fast, structured store instead of the agent re-parsing a potentially long raw transcript cold.

This is the single strongest "meaningful use" argument of the four new ideas, because it ties the Databricks module into the **core required workflow** (checkpoint reconstruction) rather than an advisory side-channel — it speaks directly to a rubric line judges are explicitly told to look for: *"Checkpoints preserve useful intent and decision context; the fresh session can resume from them"* (Use of Entire Checkpoints, 15 pts).

### 12.14 Open verification item (flag, unresolved — do this first)

Whether Lakebase is actually provisionable inside a Databricks Free Edition account is contested in current sources (see §0). §12.10–§12.13 all depend on it. **Verify with a live `lakebase` project-creation attempt in the first setup window (§5) before allocating any build time to these.** If unavailable, the v2 design (§12.1–§12.9 — AI Search + logistic regression/CRC, no Lakebase) is the already-verified, safe fallback, and none of §12.10–§12.13 are load-bearing to the core Entire submission or to the Best-of-Databricks entry.

### 12.15 ⚡ NEW (v5) — Coverage Corroboration (Databricks-enhanced tier for §10.2's mandatory fallback)

**This is an enhancement of the mandatory baseline, never a substitute for it.** The curveball's safe-fallback requirement must be satisfied whether or not Databricks is opted into — §10.2's `manual_review_recommended` flag is that baseline and stands alone. If `coverage_corroboration_enabled: true`, the same AI Search index already built for §12.4 (repo docstrings, inline comments, commit messages) is queried a second way: for each `unverifiable_coverage` entity, ask whether the repo's own documentation or commit history describes what it dispatches to or is dispatched from. A hit populates `corroboration.verified`/`evidence_snippet` exactly as §12.4 does for advisory edges; a miss leaves `verification_path: "manual_review_recommended"` unchanged. **No new index, no new infrastructure** — this is a second query shape against an index that already has to exist for §12.4, which is a genuine "meaningful use, low marginal cost" story for the Databricks judge: removing Databricks here doesn't remove the safety guarantee (the baseline still holds), it removes a *chance* at resolving what would otherwise stay a manual-review flag. State this distinction explicitly if asked — it's a more honest framing than implying Databricks is required for curveball compliance, which it is not.

---

## 13. Curveball Intelligence & Resilience

### 13.1 The structural firewall (unchanged principle from v1)

- Changes **what counts as intent** → Stage 2 only.
- Changes **what counts as drift** → Stage 4 + config only.
- Changes **how it's shown** → Stage 5 renderers only.
- Changes **the edge model / risk weighting** → config `edge_classes` + the Databricks gate, no core rewrite.

### 13.2 Ranked, evidence-based pivot predictions (new — from real Entire issue/PR mining)

| Rank | Candidate pivot | Evidence | Confidence | Architectural safeguard |
|---|---|---|---|---|
| 1 | **Unified query engine convergence** — forced to query AST structure and Checkpoint/session context through one API instead of separate calls | PR #146 (RFD 11, "converge code search on one engine"); Entire's own Sept 2026 blog post explicitly promises "one API for your code and the reasoning that produced it" | **Strong** | **Build now, regardless of the curveball:** put a thin repository/adapter layer between Fidelity's Stage 1–3 and the actual `entire` calls, so swapping two calls for one unified call is a config/adapter change, not a rewrite |
| 2 | Deterministic peer-evaluation handoff — ingest a foreign checkpoint (e.g. another team's) and critique it | `entireio/entire-judge` ("judge lenses over a submission's brain"); `session-handoff` and `review` skills | Suggestive | Keep checkpoint ingestion decoupled from the local object store — accept a ULID/HTTP-fetched checkpoint the same way as a local one |
| 3 | Subagent PII redaction under the OPF gate | Issue #2058 (redaction bypass on the `git-refs` backend); Aug 2026 changelog on a second secret-scanning engine | Suggestive | Never ingest raw subagent transcript blobs directly — go only through the documented checkpoint-read path; note this awareness explicitly in `BUILDATHON.md` limitations even if no dedicated redaction layer is built |
| 4 | Multi-repo DATA_FLOWS tracing | PR #148/#168 (DATA_FLOWS expansion); Issue #1439 (cross-repo adoption risk); Issue #32 (boundary panics) | Speculative | Parameterize the repo context on every graph call rather than hardcoding a single origin, even though cross-repo stitching itself stays out of scope (§1.1) |
| 5 | Real-time incremental sync under network partition | `entireio/git-sync` `fb287a6` (`--best-effort` exit-code tightening) | Speculative | Lower priority; only relevant if the product does live streaming ingestion, which Fidelity doesn't |

**Practical takeaway:** build the repository/adapter data-access layer as a first-class architectural decision from 9:30 AM regardless of what the actual curveball turns out to be — it's good architecture on its own merits, it directly defends the single highest-confidence prediction, and it's exactly the kind of decision worth recording explicitly in the initial-architecture Checkpoint.

### 13.3 Mandatory process (per the guide, verbatim)

At noon: stop, get to a runnable state, commit, confirm the pre-noon Checkpoint. Close the current agent session. Receive the actual constraint. **Start a fresh session** and reconstruct via Checkpoint + Graph context *before* editing. Run `graph impact` on the affected area before making changes. Implement the smallest complete response, test it, record it in a new Checkpoint.

### 13.4 `entire-judge` — a real reason to over-invest in checkpoint quality

Entire appears to run its own deterministic plugin over a submission's checkpoint history as part of evaluation — not just a human skim of your `BUILDATHON.md`. This means the actual **content** of your checkpoints (documented decisions, rejected options, stated risks) may be programmatically inspected, not just read casually. Treat every one of the four required checkpoints as a real, structured artifact — not a commit message with a slightly longer body. If `entire-judge` (or an equivalent lens) is installable, running it against your own history before submission as a self-audit is worth the ten minutes it costs.

### 13.5 ⚡ NEW (v5) — The Actual Noon Curveball, As Received

**Track determination (inference, not a stated fact — confirm with a mentor):** the card is numbered "TRACK 2" and titled **"Graph Is Evidence, Not an Oracle."** Two independent signals point to this being our card: the number (Entire Track E2 ↔ "TRACK 2"), and the content — it's written for "your Graph-powered experience," which is Fidelity verbatim. Tracks 1 (Privacy Boundary) and 3 (agent transcript-format change) read as generic/integration constraints that don't presuppose a graph-centric product. Building against Track 2 as the working assumption; the mapping isn't spelled out explicitly in the organizers' material, so treat this as high-confidence, not confirmed.

**The constraint, verbatim substance:** the repository under review now contains dynamic dispatch, generated code, or reflection that static analysis cannot fully resolve. Requirements:
- The product must not present incomplete Graph relationships as certain.
- It must identify when analysis may be partial.
- It must provide a safe fallback or verification path.
- Existing behavior for fully-resolved code must continue to work, unchanged.
- At least one test must use the supplied partial-analysis fixture.
- Users and agents must be able to tell apart: confirmed structural evidence, heuristic/incomplete evidence, and claims that require source/test verification.

**Mandatory process, per the organizers' own scoring note — graph impact analysis must run *before* editing, not after.** The rubric explicitly downgrades to the partial band if graph use happens post-implementation, or if the adaptation is only explained verbally rather than captured in Checkpoint context. This reprioritizes the existing §13.3 process step ("run `graph impact` on the affected area before making changes") from good practice to a scored, non-negotiable gate.

**What was already aligned, stated honestly:** Fidelity's foundational philosophy — tri-state edges, "honest ambiguity counts as correct" per the engine's own benchmark philosophy (§1.1), never silently promoting `inferred`/`ambiguous` to a confirmed fact (§10) — is spiritually the same value this curveball is testing for. Nothing about that needs reversing.

**The actual gap, stated plainly:** the tri-state model only classifies edges that *exist in the graph*. It has no mechanism at all for the case where a relationship was never extracted in the first place because the surrounding code is dynamically dispatched, generated, or reflective. Today, a changed entity sitting in one of these unresolved regions falls through Stage 4's reachability check exactly like any other entity with no path from a confirmed node — and gets filed as confident `undeclared_scope_creep`. That is precisely the failure mode the curveball names: presenting an absence of evidence as evidence of absence, with certainty the graph never actually had. See §10.2 for the fix, §8.1/§8.2 for the schema and config changes, §12.15 for the fallback-path split, §14 for the build-sequence change, and §15 for the required fixture test.

**Ranked-prediction accuracy, stated honestly (per house style — flag misses, don't smooth over them):** none of the five ranked pivots in §13.2 named this exact constraint. The closest partial relevance is rank 4 (multi-repo `DATA_FLOWS` tracing, evidence: Issue #32's boundary panic) — it shares the general instinct of "graph search can hit an unresolved edge case and must degrade gracefully instead of crashing or asserting confidently," and the defensive-wrapping habit built for Issue #32 (§1.1, §15's graceful-panic test) generalizes usefully here. But the coverage-confidence signal itself — a first-class, always-computed field distinct from crash-handling — is new work this revision adds, not something the existing architecture already absorbed for free.

---

## 14. Build Sequence (per the Participant Guide's timeline — see the §0 flag on the conflicting version)

| Window (IST) | Target |
|---|---|
| 8:00–9:00 | Breakfast/check-in; confirm GitHub access, agent tooling, optional Databricks Free Edition account |
| 9:00–9:30 | Fork (post-kickoff only), mirror (India region), clone; run §3.1's exact command sequence; run §5's pre-flight checklist; write `fidelity.config.yaml` skeleton |
| 9:30–11:45 | Stages 1–5 working end-to-end on a small real repo; terminal renderer first; **lock `verdict.json` schema**; write meaningful tests (§15); build the repository/adapter layer (§13.2) as part of the base architecture, not as a curveball reaction |
| 11:45–12:00 | Get to runnable state; commit; confirm the pre-noon Checkpoint exists with real decision content |
| 12:00–1:00 | Noon Curveball + lunch; close session; receive constraint; **do not implement yet** |
| 1:00–3:00 (or later, per §0) | ⚡ **v5, now the live path:** fresh session; reconstruct via pre-noon Checkpoint (intent, architecture, completed work, open risks); **run `entire graph impact` on the affected area BEFORE any edit — this is scored, not optional, per §13.5**; implement §10.2's coverage-confidence signal + `unverifiable_coverage` tier + baseline fallback flag; update all three renderers for the three-way visual split (§11.1); add the partial-analysis fixture test (§15); wire §12.15 only if the Databricks module (§12) is already solid; write the final Checkpoint stating explicitly which assumption the curveball invalidated and why the revised design is safe |
| Final 20 min before deadline | Run the §17 checklist verbatim; write `BUILDATHON.md`; confirm demo owner can access every resource |
| Deadline | Submit — do not miss it chasing a stretch feature |
| Post-deadline | Judging — be ready to explain your own code without the agent, per the guide's explicit accountability requirement |

---

## 15. Testing

- **Golden-manifest tests:** 2–3 hand-built (transcript, diff) fixtures with an exact expected `verdict.json`.
- **Fail-closed test:** an LLM fallback given a non-matching fragment must produce `unresolved_fragment`, never an invented entity.
- **No-baseline test:** first-checkpoint case yields only `confirmed`/`undeclared_scope_creep`, tagged `no_baseline: true`.
- **Graceful-panic test:** feed a known Issue #32–style boundary case through the pipeline and confirm it degrades to an explicit warning, not a crash.
- **Renderer purity test:** all three renderers agree on tier counts from the same `verdict.json`.
- **Reproducible setup:** one script that builds a tiny sample repo and runs `verify-intent` to a known verdict.
- **⚡ Partial-coverage test (v5, mandatory per the Curveball card, not optional):** run the supplied partial-analysis repository fixture through the full pipeline. Assert: (1) no entity in that region is misfiled into `undeclared_scope_creep`; (2) every affected entity carries `coverage_confidence: partial` (or `unknown` if the fixture triggers a snapshot error) and lands in `unverifiable_coverage`; (3) `verification_path` is populated (`manual_review_recommended` at minimum); (4) a fully-resolved entity elsewhere in the *same* fixture still reconciles normally — proving existing behavior for resolved code is unchanged, per the curveball's explicit requirement.
- **Optional:** run `entire-judge` (if available) against your own checkpoint history as a final self-audit (§13.4).

---

## 16. Roadmap — Everything Parked (nothing lost, refined reasoning)

1. **Plan-Trust Ledger** — persist (claim, outcome) pairs across sessions in Delta; learn how much to trust a given agent/model's *plans* over time. Fidelity's own verdicts are the training data. **v3: potentially promotable to a live build via the Verdict Ledger, §12.11 — pending the Free Edition Lakebase verification in §12.14.**
2. **Checkpoint Court** — a second, adversarial Agent Bricks agent tries to falsify each verdict claim before acceptance, logged via Unity AI Gateway as a replayable trial record.
3. **Forensic Replay** — Delta time-travel on versioned, ULID-keyed snapshots + trace capture, to reconstruct what an agent actually saw at a past decision and re-judge it against the current graph. **v3: the "what-if verdict" capability in §12.10 is a narrower, live-buildable version of this same idea — pending §12.14.**
4. **Subagent Reasoning Capture** — intercept subagent reasoning before condensation (a real, open bug otherwise loses it — Issue #2058) and reconcile each subagent's sub-intent against its scoped graph diff.
5. **Compliance reconciliation via Unity Catalog lineage** — resolve a code change through data lineage to "which regulated fields/attestations does this touch." Viable because lineage is real relational metadata (unlike Genie Ontology as a graph store, §1.2).
6. **SCIP export alignment** — once Entire's RFD 0006 SCIP-export convergence ships, expose `verdict.json` findings through that standard channel for cross-tool interoperability, rather than a bespoke format.

**Explicitly out of scope, with the corrected reasoning:** fleet/cross-repo *structural* graph analytics — `entire-graph` itself enforces single-store isolation (PR #169); the *checkpoint/session*-level cross-repo retrieval that does exist (`entire explain --repo`) is a different capability and doesn't change this. Test-selection on `TESTS` edges and cross-service routing on `HANDLES_ROUTE`/`HANDLES_GRPC` — unverified in the source. Genie Ontology as the graph store — §1.2.

---

## 17. `BUILDATHON.md` — exact outline (per the Participant Guide)

```markdown
Project name
One-sentence summary
Problem, intended user and why it matters
Selected Entire track and why Entire is essential
  [include the §1.4 differentiation from review/explainskill/what-happened here]
Architecture and main workflow
Entire Graph findings and verification
  [explicitly name the 3 required demonstrations from §3.3]
Noon Curveball: what changed and how we adapted
  [Track 2 — "Graph Is Evidence, Not an Oracle" (§13.5); state honestly that none of §13.2's five ranked predictions called it; name the assumption invalidated (zero-edge ⇒ scope creep) and the fix (coverage-confidence signal + unverifiable_coverage tier, §10.2)]
Checkpoint links and what each checkpoint proves
  [all 4 required milestones from §3.2, each with real decision content]
Setup, run and test instructions
Databricks use, data sources and limitations (if applicable)
  [capabilities used, why essential, honest provenance per §12.8]
Known limitations and next steps
  [pull from §16's roadmap; note which §20 open-source additions (if any) actually made it into the build — e.g. Semgrep-backed coverage-confidence detection — versus which remain roadmap-only]
```

---

## 18. Final Submission Checklist (verbatim from the guide, applied to Fidelity)

- [ ] Final commit pushed; SHA matches the submission.
- [ ] Project launches from a clean checkout / documented setup path.
- [ ] All four required checkpoints open and clearly explain their milestone.
- [ ] Entire Graph evidence and the final semantic-diff analysis are recorded and named.
- [ ] `BUILDATHON.md` complete, readable, free of secrets.
- [ ] Tests covering critical + Curveball behavior pass.
- [ ] Databricks resource links and data notes included, if opted in.
- [ ] Demo owner can sign in, open every resource, run the critical path.
- [ ] A fallback screenshot/recording exists locally.
- [ ] Submitted before the deadline — **confirm the actual time at kickoff** (§0).

---

## 19. Positioning (for the pitch)

> *"Every AI coding agent tells you what it did. None of them prove it did only what you asked. Fidelity reads the Checkpoint's stated intent, reads the Graph's actual structural change, and produces a verdict you can check against source in seconds — not a chat answer, an audit trail. And when the graph itself isn't certain, Fidelity isn't certain — it degrades honestly instead of asserting a dependency that might not exist."*

Two lines to have ready cold:

- **If asked "isn't this just the `review` skill?"** — "`review` narrates why a change makes sense by reading the transcript with an LLM. Fidelity verifies whether it matches what was claimed, tier by tier, backed by actual graph reachability over classified edges — a judgment call versus a computed reconciliation. They'd compose well together."
- **For the Databricks judge:** "We deliberately didn't force the code graph into Genie Ontology — the 200-snippet cap and SQL-only traversal make it the wrong tool for this. Instead we used Databricks for a real trust-but-verify gate: a conformally-calibrated confidence score on heuristic edges, corroborated against the repo's own documentation via AI Search's hybrid retrieval — governed by Unity Catalog the whole way through."
- **⚡ If asked "what did the curveball actually change?":** "Our tri-state edge model already handled uncertainty *about a recorded relationship*. What it couldn't handle was a relationship that was never recorded at all — dynamic dispatch, generated code, reflection. We added a second, orthogonal signal — coverage-confidence — so a changed entity in one of those blind spots is never filed as confident scope creep. It gets its own tier, a mandatory manual-review flag, and an optional Databricks-corroboration upgrade that reuses infrastructure we'd already built — not a new dependency introduced under time pressure."
- **⚡ v6, if Semgrep made it into the build (§20.1):** "We didn't stop at a zero-edge heuristic for `partial` coverage — we ran a real static pattern-match for dynamic dispatch, reflection, and codegen signatures underneath it. So when we say an entity's coverage is uncertain, that's a detected finding, not just an absence we're guessing about."
- **If Lakebase cleared §12.14 and made it into the build:** "We went further than an advisory score — every verdict is a row in a live Lakebase ledger, so the Databricks layer is also the reason a fresh session after the Noon Curveball can reconstruct trust history instead of re-reading a raw transcript cold, and a reviewer can branch the graph zero-copy to ask 'what if the blast radius were 3 hops' and get a real answer in seconds, not a re-run."

---

## 20. ⚡ NEW (v6) — Open-Source Model & Tool Augmentation, Prioritized Against Remaining Time

**Framing, stated plainly:** this section is a menu, not a mandate. Everything in it is a real, genuinely-fitting addition to a section already open for extension elsewhere in this bible — nothing here is invented scope. It is explicitly triaged by whether it strengthens the *live, scored* Curveball response (§13.5) versus whether it's a legitimate but optional upgrade to an already-working mechanism versus whether it's roadmap-only given the clock. **If only one thing from this section gets built, §20.1 is the recommendation** — see §20.5.

### 20.1 🔴 Curveball-critical — strengthens §10.2 directly, do this first

**Semgrep** (open source, pattern-based static analysis). Today §10.2's `partial` coverage detection is a heuristic proxy — "zero edges ⇒ maybe dynamic dispatch, maybe just genuinely disconnected code." Honest, but weak under judge scrutiny. Semgrep ships prebuilt rule packs that pattern-match real dynamic-dispatch/reflection/codegen signatures per language (`getattr`, `eval`, `__getattr__`, reflection APIs, decorator-generated methods, etc.), so an entity can be promoted to `coverage: partial` on **detected evidence**, not absence alone. Same-afternoon integration — no training, no model hosting. This is the single highest-leverage addition in this section because it upgrades the exact mechanism the curveball is scoring, not an adjacent one.

**Joern** (open source code property graph) — a heavier alternative to Semgrep for the same gap, purpose-built for modeling reflective/dynamic call edges. Only worth it if Semgrep alone doesn't feel convincing in the demo and 30+ minutes remain; otherwise Semgrep is the better time-to-value trade. Do not attempt both.

### 20.2 🟡 Real, lower-priority upgrades to already-cited mechanisms

| Tool | Section it upgrades | Why |
|---|---|---|
| **MAPIE** (open-source conformal prediction library) | §12.2 (Conformal Risk Control) | The CRC threshold `λ̂` is currently hand-rolled. MAPIE is a tested, citable CRC implementation — a near drop-in swap that makes the "mathematically bounded" claim in §1.5 defensible against a Databricks judge who asks to see the math. |
| **leidenalg + python-igraph** | §10.1 (community-aware scope-creep) | The actual Leiden implementation the bible references by name but doesn't pin. Only relevant if the demo repo is large/connected enough for §10.1 to produce non-trivial communities (§10.1's own honest caveat). |
| **sentence-transformers (`all-MiniLM-L6-v2`) + FAISS or `rank_bm25`** | §12.4 / §12.15 (AI Search corroboration) | A local, open-source hybrid-retrieval fallback that works identically with or without Databricks cooperating on time. Worth naming to the Databricks judge as "this corroboration step isn't platform-locked" — but only add it if AI Search is already proving flaky; don't build a second retrieval stack for redundancy alone under time pressure. |
| **Qwen2.5-Coder (1.5B/7B, open weights)** | §9 / §12.2 (extraction fallback + I-CALM scoring) | A strong, code-tuned, cheaply-run open-weight model, well-behaved for structured "pick from this list or abstain" prompting — a good citation if asked why that specific model. Not worth switching to mid-build unless the current model choice is actively causing problems. |

### 20.3 🟢 Explicitly not-now (roadmap only — fold into §16)

- **Official MCP Python/TypeScript SDK** for §11.1's agent-facing self-check tool — low-effort once `verdict.json` is stable, but additive scope this late in the day, not curveball-scoring.
- **CodeQL / Kythe** — both real, both already named in §1.5's competitive landscape, both heavier integrations than Semgrep for the identical §10.2 payoff. A good `BUILDATHON.md` §16 roadmap line, not a same-day add.
- **Graphify itself** — never integrate as a dependency. §1.5 exists specifically to differentiate Fidelity's reconciliation-verdict model from Graphify's context-retrieval model; adding it as a component would blur that positioning, not strengthen it.

### 20.4 What does not change

This section adds implementation options underneath mechanisms §§9–12 already specify — it does not reopen the Stage 1–4 interfaces, the locked `verdict.json` schema (§8.1), or the structural firewall (§13.1). Semgrep's output lands as an input to §10.2's existing `coverage` field computation, the same way the zero-edge heuristic already does; nothing about the schema shape or the tier logic changes.

### 20.5 Recommendation, stated once

Given the live build clock (§14, §0/v5): build **§20.1's Semgrep integration** if only one addition is made — it directly strengthens the mechanism the noon Curveball is actively scoring. Everything in §20.2 is a genuine improvement but optional against the deadline; §20.3 is honest future work, not something to reach for today.
