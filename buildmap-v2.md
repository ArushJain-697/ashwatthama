# Fidelity — Detailed Build Map & Step-by-Step Implementation Guide (v2)

**For anyone building Fidelity who has NOT read the full Bible.** This file is self-contained and ordered so each ticket can be built and tested on its own before moving to the next. The companion `FIDELITY_FINAL_BIBLE_v6.md` has the deep "why" — read it if a step's rationale is unclear.

## What Changed Since v1

- **The Noon Curveball is no longer hypothetical.** v1's generic "receive constraint, implement smallest response" placeholder (old #30/#33) is replaced with concrete tickets against the actual received card: **Track 2 — "Graph Is Evidence, Not an Oracle"** (#30), and a full **coverage-confidence** implementation chain (#33–#37) that didn't exist in v1 at all.
- **A mandatory three-way visual split** is now a scored renderer requirement, not a design nicety — the terminal renderer gets retrofitted (#38) and the Markdown/dashboard renderers are built three-way-aware from the start (#39–#40), with a dedicated distinctness test (#41).
- **A mandatory fixture test** (#42) is new — the curveball card requires at least one test against a supplied partial-analysis fixture, with four specific assertions.
- **Checkpoint 3** (#44) now has real content requirements: name the card, name the invalidated assumption, name the fix, and — per house style — honestly state that none of the five ranked predictions called this exact pivot.
- **A new open-source augmentation phase** (#45–#46, plus #49/#52 inline in the Databricks phase) — Bible §20 triages these explicitly; **Semgrep (#45) is the single recommended addition if time allows for only one.**
- **§12.15's Databricks-enhanced coverage corroboration** (#56) is new — it reuses the existing AI Search index, never replaces the mandatory baseline (#35).
- **Community-aware scope-creep clustering** (#57, Leiden via `leidenalg`+`python-igraph`) is restored as an explicit stretch ticket.
- **Fidelity-as-MCP-tool (old §11.1) is demoted out of the ticket list entirely** — Bible §20.3 explicitly places it in "not-now" roadmap territory; it now only appears in Future Work.
- Pre-flight verification (#7), the config skeleton split (#9 vs. the new #37), the differentiation pitch (#3, now covering Graphify/§1.5), the final rehearsal (#63), and `BUILDATHON.md` (#64) are all updated for v6 content.

**Build model:** a single continuous work session — the **BTW Buildathon 2026 · Track E2, Build with Graph Intelligence · + Best Use of Databricks (opt-in)**. The Noon Curveball has already happened by the time this revision was written: **Track 2 — "Graph Is Evidence, Not an Oracle"** (§13.5), received at the mandatory 12:00–1:00 gate. There is still a flagged, **unresolved deadline conflict** (§0: 3:00 PM per the Participant Guide vs. 4:00 PM per the event website). Build order targets a fully working core pipeline (Stages 1–5) locked before noon, then spends the remaining window on the concrete coverage-confidence response, the optional Databricks module, and — time permitting — the open-source augmentation pass.

This map is granular on purpose: bigger pieces are split so each ticket is buildable and independently testable in a single sitting, and every ticket has an explicit **Test** step, not just a "done when" description.

---

## How to Read This Map

Tickets are numbered in strict dependency order — `Blocked by` only ever points to a lower ticket number, so the logical build sequence never leaves you stuck on something undone.

### Priority tag

| Tag | Meaning |
|---|---|
| 🟥 **CORE** | The Entire-track spine. Every point on the 100-pt Entire rubric traces to one of these. If you build nothing else, build these. |
| 🟨 **FIREWALL** | Architecture decisions made before the Curveball specifically so it could be absorbed via config/adapter, not a rewrite (§13.1–13.2). |
| 🟦 **DATABRICKS** | The opt-in Best-of-Databricks bolt-on (§12). Attempt only once the CORE stages work end-to-end — it never touches Stages 1–4's interfaces. |
| ⬜ **STRETCH** | Genuinely optional extras — Lakebase v3 (§12.10–§12.13) and community clustering (§10.1). Safe to cut entirely. |
| 🔴 **AUGMENT-CRITICAL** | NEW v6 (§20.1) — the single highest-leverage open-source addition, directly strengthening the mechanism the Curveball is scoring. Do this first if only one addition is made. |
| 🟡 **AUGMENT-OPTIONAL** | NEW v6 (§20.2) — real, legitimate upgrades to already-cited mechanisms. Genuinely optional against the clock. |

### Light tag (which window a ticket belongs in)

| Tag | Meaning |
|---|---|
| ⚪ **PRE-EVENT** | Decide/verify before or right at the 9:00 kickoff — no `entire` commands yet. |
| 🌅 **MORNING** | The 9:00–11:45 core build window, before the pre-noon lockdown. |
| 🛑 **NOON-GATE** | 11:45–1:00 — the mandatory stop/commit/curveball-receipt window. Nothing gets *built* here except the required process steps themselves. |
| 🌇 **AFTERNOON** | 1:00 until the (unresolved, see #1) deadline — the concrete Curveball response, remaining renderers, Databricks module, augmentation pass. |
| 🏁 **FINAL** | The last ~20 minutes — packaging and submission only, no new features. |

### Per-ticket fields

Each ticket has: **Build**, **Watch out**, **Test**, **Done when**, plus a **Milestone** line (the build-sequence window it targets — kept distinct from "Checkpoint," which is Entire's own literal artifact and a ticket in its own right: #12, #29, #44, #61).

**Golden rule:** don't start a ticket until its blockers show green on their own Test step. A ticket that "mostly works" is not done — the next ticket silently inherits its bugs.

**Scheduling rule:** three tickets are now load-bearing for the whole afternoon — **the adapter layer (#10)**, **the locked pre-noon `verdict.json` schema (#26)**, and **the coverage-confidence signal (#33)**. Everything from the three-way renderer split through both required post-noon graph demonstrations depends on one of these holding still. If time is genuinely short in the afternoon, the single recommended addition beyond the mandatory response is **Semgrep (#45)** — per §20.5, build that before any other optional item in this map. Don't let augmentation tickets (#45–#46, #49, #52, #57) compete for time against anything still red in Phases 8–9; those are the scored core of the Curveball response.

---

## Phase 0 — Pre-Event Decisions

### #1: Timeline Conflict Flagged for Kickoff Verification
Priority: 🟥 CORE
Type: Discuss
Light: ⚪ PRE-EVENT
Milestone: Before kickoff (8:00–9:00)
Blocked by: —

**Build:** Write down both candidate deadlines (3:00 PM per the Participant Guide, 4:00 PM per the event website) and assign one person to get the live, spoken answer at kickoff and confirm it back to the team in writing.
**Watch out:** these are two different documents from the same organizers, an hour apart — don't quietly pick one and hope.
**Test:** after kickoff, confirm the whole team has heard the same stated deadline out loud.
**Done when:** one deadline is confirmed and communicated before 9:30.

### #2: Track, Databricks Opt-In & Ownership Locked
Priority: 🟥 CORE
Type: Discuss
Light: ⚪ PRE-EVENT
Milestone: Before kickoff (8:00–9:00)
Blocked by: —

**Build:** Finalize the track (**E2, Build with Graph Intelligence**), the **Databricks opt-in decision**, the **GitHub fork target** and India-region mirror choice, and **one named demo/deployment owner** who can sign in and run the critical path on judging day (§3.4, §18).
**Watch out:** deciding the demo owner after the build is finished is a scramble, not a plan.
**Test:** each item has a single written, agreed answer with no open "we'll decide later."
**Done when:** track, opt-in decision, fork target, and owner are settled in writing before you arrive.

### #3: Differentiation Pitch Drafted
Priority: 🟥 CORE
Type: Discuss
Light: ⚪ PRE-EVENT
Milestone: Before kickoff (8:00–9:00)
Blocked by: —

**Build:** Draft, in writing, the three answers a judge is most likely to probe: "isn't this just `review`/`explainskill`/`what-happened`?" (§1.4 — retrieval-and-narration vs. deterministic graph-verified reconciliation), "why not Genie Ontology?" (§1.2 — 200-snippet cap, SQL-only traversal), and **NEW — "isn't this just Graphify's tri-state model?"** (§1.5: Graphify's EXTRACTED/INFERRED/AMBIGUOUS tagging routes an agent's *attention* during exploration; Fidelity's signal drives a deterministic *verdict* after the fact — reconciliation, not retrieval). Also prepare the one-line rejection of the "Thousand-Graph Hypothesis" for a technically literate judge: Fidelity deliberately keeps persisted, classified edges over per-query LLM latent materialization, because an audit-trail product needs a verdict a reviewer can trust as reproducible, not a fresh inference each time.
**Watch out:** there are now three near-guaranteed gotcha questions, not two — don't let the newest one (Graphify) be the one nobody rehearsed; it's the closest taxonomic lookalike in the whole competitive landscape.
**Test:** say all three answers out loud, cold, in under 30 seconds each.
**Done when:** all three answers are written down and rehearsed once.

### #4: Curveball Intelligence Briefing
Priority: 🟨 FIREWALL
Type: Research
Light: ⚪ PRE-EVENT
Milestone: Before kickoff (8:00–9:00)
Blocked by: —

**Build:** Read the ranked pivot predictions (§13.2) as a team before the event: unified-query-engine convergence (strong confidence), peer-evaluation handoff, subagent PII redaction, multi-repo DATA_FLOWS tracing, real-time incremental sync under partition.
**Watch out:** this briefing exists so the adapter layer (#10) isn't a mystery decision on the morning of. It's also the honest baseline that Checkpoint 3 (#44) gets measured against — the actual card (#30) turned out not to match any of these exactly, and that gap needs to be stated plainly, not smoothed over, when the time comes.
**Test:** each team member can name the #1-ranked prediction and its safeguard without looking it up.
**Done when:** the whole team has read §13.2 once before kickoff.

---

## Phase 1 — Kickoff Setup & Pre-Flight Verification

### #5: Mirror, Clone & Checkpoints Enabled
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Kickoff (9:00–9:30)
Blocked by: #2

**Build:** Run the exact §3.1 sequence: `entire login`, `entire repo mirror create`, `entire repo clone /gh/YOUR-GITHUB-HANDLE/REPOSITORY` (select the fork and **India region**), then `entire enable -y --agent YOUR-BUILD-AGENT` and `entire status` — before any agent-assisted development begins.
**Watch out:** work only inside the clone created through the mirror flow, and enable Checkpoints before development starts.
**Test:** `entire status` reports Checkpoints enabled against the correct mirrored repo.
**Done when:** the clone exists, is mirrored to the India region, and Checkpoints are confirmed live.

### #6: Graph Plugin Installed + Fresh Agent Session
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Kickoff (9:00–9:30)
Blocked by: #5

**Build:** Run `entire plugin install graph`, `entire graph version`, `entire graph init-agents --repo .`, then **immediately start a fresh agent session**.
**Watch out:** an agent session started before `init-agents` never gets the graph context and every downstream graph call degrades silently.
**Test:** confirm `entire graph version` returns successfully and the new agent session acknowledges graph tooling.
**Done when:** the graph plugin is installed and a fresh, graph-aware agent session is running.

### #7: Pre-Flight Graph Capability & Edge-Field Verification
Priority: 🟥 CORE
Type: Research
Light: 🌅 MORNING
Milestone: Kickoff (9:00–9:30)
Blocked by: #6

**Build:** Run the §5 checklist against the live graph: `entire graph capabilities --json` for the actual edge-type enum (do not assume `TESTS`/`HANDLES_ROUTE`/`HANDLES_GRPC` exist), one real `entire graph snapshot` for node/edge fields and the exact tri-state classification key, one real `entire graph diff --base <ref> --head <ref>` for output shape, `graph impact --help`/`graph neighbors --help` for depth/direction flags, a malformed/boundary-file sanity test against `graph search` for Issue #32's panic, and checkpoint transcript reads for both ULID and legacy hex ID formats. **NEW (v5 item):** also check whether `entire graph capabilities --json` or `entire graph snapshot` exposes any native "partial/unresolved" or "unindexed region" marker before assuming a zero-edge heuristic is the only option for detecting incomplete coverage.
**Watch out:** the new coverage-marker check matters later, not just now — the afternoon's coverage-confidence work (#33) needs to know which detection mechanism is real before it's built. Recording the wrong assumption here costs a rebuild later, not just a wrong checklist box.
**Test:** each of the seven checklist items produces a recorded, real output — not an assumption.
**Done when:** all findings are recorded at the top of `BUILDATHON.md`, including whether Issue #32 reproduces and whether a native coverage marker exists.

### #8: Lakebase Free-Edition Availability Verified
Priority: ⬜ STRETCH (gate)
Type: Research
Light: 🌅 MORNING
Milestone: Kickoff (9:00–9:30)
Blocked by: #5

**Build:** Attempt a live `lakebase` project-creation call inside the Databricks Free Edition account, in parallel with #7, before allocating any build time to §12.10–§12.13.
**Watch out:** sources disagree — only a live attempt resolves this for your account.
**Test:** the project-creation call either succeeds (clean yes) or fails with a clear unsupported-feature error (clean no).
**Done when:** the yes/no answer is recorded, and it's the single gate deciding whether #58–#60 are attempted at all.

### #9: `fidelity.config.yaml` Skeleton Written
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Kickoff (9:00–9:30)
Blocked by: #7

**Build:** Write the initial config file per the v1–v4 shape of §8.2: `blast_radius_hops: 2`, `edge_classes`, the five original tier keys, extraction settings, `verdict_thresholds.review_required_if_scope_creep_gte`, and the `databricks` block defaulted to `calibration_enabled: false` / `score_source: batch` / `gate.mode: fallback`.
**Watch out:** don't pre-build the `coverage:` block or the `unverifiable_coverage` threshold here — those don't exist yet at this point in the real timeline; they're added later, in #37, once the Curveball response has actually been designed. Default `score_source` to `batch`, not `serving_endpoint`.
**Test:** the config parses cleanly and every downstream stage can read its section without a schema error.
**Done when:** a valid, five-tier config skeleton exists and is checked in.

---

## Phase 2 — Architecture Foundation & Checkpoint 1

### #10: Repository/Adapter Data-Access Layer
Priority: 🟨 FIREWALL
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #7

**Build:** Put a thin repository/adapter layer between Fidelity's Stage 1–3 logic and the actual `entire` CLI calls, so every graph/checkpoint interaction goes through one seam.
**Watch out:** this defends against the #1-ranked curveball prediction (§13.2) — build it regardless of what the actual curveball turned out to be; it's good architecture on its own merits and belongs in Checkpoint 1.
**Test:** swap one real `entire` call behind the adapter for a stub and confirm nothing above the adapter layer needs to change.
**Done when:** every Stage 1–3 call to `entire` goes through the adapter, not a direct CLI invocation.

### #11: `verify-intent` Command Surface Skeleton
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #10

**Build:** Stand up the command surface from §7: `entire graph verify-intent <checkpoint-id> [--base <checkpoint-id>] [--config fidelity.config.yaml] [--out ./fidelity-out/]`, wired to no-op stage functions.
**Test:** run with placeholder args and confirm it parses flags and reaches each stub stage in order.
**Done when:** the CLI entry point exists, argument-parses correctly, and calls stubs for all five stages.

### #12: Checkpoint 1 — Initial Understanding & Architecture Recorded
Priority: 🟥 CORE
Type: Discuss
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #2, #9, #10, #11

**Build:** Write the first required Entire Checkpoint (§3.2 milestone 1): the problem/user framing, the chosen architecture, and the decision to build the adapter layer (#10) as a pre-emptive curveball defense.
**Watch out:** `entire-judge` may programmatically inspect checkpoint content (§13.4) — real decision content, not a commit message with a longer body.
**Test:** re-read cold; confirm a stranger could reconstruct why the adapter layer exists from it alone.
**Done when:** Checkpoint 1 is committed with real decision content.

---

## Phase 3 — Stage 1–2: Capture & Declare Intent

### #13: Stage 1 — Checkpoint Transcript Capture
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #11

**Build:** Implement Stage 1: read a checkpoint's transcript via the adapter, given a checkpoint ID.
**Test:** point it at a real checkpoint from #12; confirm the raw transcript comes back correctly.
**Done when:** any valid checkpoint ID returns its transcript through the adapter.

### #14: Stage 2 — Symbol Dictionary Load (`entire graph snapshot`)
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #13

**Build:** Load the full symbol dictionary via `entire graph snapshot` — the **first required graph demonstration** (§3.3).
**Watch out:** this snapshot must also expose the tri-state classification key confirmed in #7 — a wrong field name here silently breaks every later tier.
**Test:** run the snapshot on the demo repo; confirm symbol names and the tri-state field parse correctly.
**Done when:** a full, correctly-keyed symbol dictionary loads on demand.

### #15: Stage 2 — Regex/Fuzzy Match Extraction
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #14

**Build:** Fuzzy-match transcript tokens against the symbol dictionary to extract entities with `method: regex`/`confidence: 1.0` where matched.
**Test:** feed a transcript naming a real symbol verbatim; confirm it's extracted with high confidence.
**Done when:** direct symbol mentions are reliably extracted.

### #16: Stage 2 — LLM Fallback, Fail-Closed
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #15

**Build:** For fragments #15 can't resolve, send only the fragment plus the constrained symbol dictionary to an LLM, instructed to pick from the exact list or return "no confident match" — never invent an entity.
**Watch out:** this respects the graph's own anti-forgery posture (PR #143). If asked why a particular open-weight model was chosen for this step, **Qwen2.5-Coder** (1.5B/7B) is the Bible's citable, code-tuned, cheaply-run reference (§20.2) — worth naming if it's actually the model in use, not worth switching to mid-build otherwise.
**Test:** feed a fragment with no real match; confirm the output is `unresolved_fragment`, never a fabricated entity name.
**Done when:** every extracted entity carries `method`/`confidence`, and unmatched fragments fail closed.

### #17: Stage 2 — Inherited-Intent & No-Baseline Tagging
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #16

**Build:** Vague prompts inherit intent from the nearest resolvable ancestor checkpoint (`inherited_intent: true`); the first checkpoint diffs against an empty tree (`no_baseline: true`).
**Test:** feed a vague follow-up; confirm correct inheritance and tagging. Run the very first checkpoint; confirm the `no_baseline` tag appears.
**Done when:** both edge cases are tagged correctly.

---

## Phase 4 — Stage 3: Observe (Graph Diff & Tri-State Classification)

### #18: Stage 3 — `entire graph diff` Changed-Entity Extraction
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #7

**Build:** Call `entire graph diff --base <base> --head <checkpoint-id>` through the adapter to get changed entities for a checkpoint pair.
**Test:** diff two real checkpoints with a known code change; confirm the changed entity list matches.
**Done when:** changed entities come back correctly for any valid checkpoint pair.

### #19: Stage 3 — Tri-State Edge Classification Wired
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #18

**Build:** Attach the tri-state class (`extracted`/`inferred`/`ambiguous`) to every edge touched by the diff.
**Watch out:** don't collapse the tri-state into a single confidence number — each state has a distinct downstream meaning.
**Test:** feed a diff with mixed edge types; confirm each edge tags correctly.
**Done when:** every edge in a diff carries its real tri-state classification.

### #20: Stage 3 — Issue #32 Panic Guard
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #19

**Build:** Wrap the DATA_FLOWS expansion call defensively so a slice-bounds panic (`[-1:0]`, Issue #32) degrades to an explicit warning, not a crash.
**Test:** feed the known boundary case from #7; confirm a warning, not a stack trace.
**Done when:** the specific reproduced panic case degrades gracefully every time.

---

## Phase 5 — Stage 4: Reconcile (Base Tiers)

### #21: N-Hop Reachability Helper (`neighbors`/`impact`)
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #19

**Build:** Implement reachability from a confirmed entity out to `blast_radius_hops` using `entire graph impact`/`neighbors` if depth flags exist, or by iterating single-hop `neighbors`.
**Watch out:** this is the mechanism behind the **second required graph demonstration** (§3.3) — it needs to actually run before a change is tiered as safe, not just exist as a library function.
**Test:** run from a known entity with a known 2-hop dependency; confirm path, hop count, and edge types traversed.
**Done when:** N-hop reachability is available and records `edge_types_traversed`/`min_edge_class` per path.

### #22: Confirmed / Declared-Unimplemented Tiering
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #16, #21

**Build:** Tier extracted entities: `confirmed` when the diff shows a matching change; `declared_unimplemented` when a prompt-named entity has no corresponding graph change.
**Test:** feed one implemented and one named-but-untouched intent; confirm correct tiering.
**Done when:** both tiers populate correctly.

### #23: Expected Blast-Radius Tiering
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #21, #22

**Build:** Tier a changed entity as `expected_blast_radius` when reachable from a `confirmed` entity within `blast_radius_hops` via deterministic edges only.
**Watch out:** deterministic-edge reachability only — never blend in `inferred`/`ambiguous` paths here.
**Test:** feed a confirmed change with a 1-hop deterministic dependency; confirm correct tiering.
**Done when:** blast-radius entities tier correctly and never rest on a non-deterministic edge.

### #24: Undeclared Scope-Creep Detection
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #18, #22

**Build:** Tier a changed entity as `undeclared_scope_creep` when unreachable from any confirmed entity within the hop limit.
**Test:** feed an unrelated, unreachable change; confirm it's flagged with a stated reason.
**Done when:** genuinely unreachable changes are reliably flagged.

### #25: Advisory Low-Confidence Tier
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #19, #23

**Build:** Tier anything whose blast-radius reasoning only holds via an `inferred`/`ambiguous` edge as `advisory_low_confidence`, with the Databricks-reserved fields left `null`.
**Watch out:** never silently promote one of these to `expected_blast_radius`.
**Done when:** advisory entries never get promoted, and the schema's optional fields are correctly stubbed.

---

## Phase 6 — Stage 5: Manifest & Terminal Render (locked before noon)

### #26: `verdict.json` Schema Locked (five-tier shape)
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #22, #23, #24, #25

**Build:** Assemble `verdict.json` per the v1–v4 shape of §8.1: checkpoint IDs, timestamp, verified graph capabilities, the intent block, the five original reconciliation tiers, and the summary block.
**Watch out:** this is load-bearing — every renderer and the entire Databricks module are pure functions of this shape. It will be *extended*, not replaced, once the coverage-confidence work lands (#36) — lock this version cleanly now so the extension is additive.
**Test:** run the full pipeline on a real transcript+diff pair; confirm the output matches the schema exactly.
**Done when:** the five-tier `verdict.json` shape is final and every upstream tier writes into it correctly.

### #27: Terminal Renderer (five-tier shape)
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #26

**Build:** A colored, tiered, `file:line`-annotated terminal report as a pure function of `verdict.json`.
**Test:** run against a `verdict.json` with entries in all five tiers; confirm every tier renders distinctly.
**Done when:** the terminal report is legible and complete from `verdict.json` alone.

### #28: Golden-Manifest + Fail-Closed + No-Baseline + Graceful-Panic Tests
Priority: 🟥 CORE
Type: Prototype
Light: 🌅 MORNING
Milestone: Morning build (by 11:45)
Blocked by: #27

**Build:** 2–3 hand-built golden-manifest fixtures; the fail-closed test (#16); the no-baseline test (#17); the graceful-panic test (#20).
**Test:** run the full suite; confirm all four classes pass.
**Done when:** the core pipeline has a real, passing regression suite before the noon lockdown.

---

## Phase 7 — Pre-Noon Lockdown & Curveball Receipt

### #29: Checkpoint 2 — Pre-Noon Stable State Committed
Priority: 🟥 CORE
Type: Discuss
Light: 🛑 NOON-GATE
Milestone: Pre-noon lock (11:45)
Blocked by: #28

**Build:** Get to a runnable state, commit, and write the second required Entire Checkpoint: the last stable state before the Noon Curveball, with real decision content on what's done and what's deliberately deferred.
**Watch out:** "runnable" means #28's suite passes clean at the moment of this checkpoint.
**Test:** run #28's suite one more time immediately before committing; confirm green.
**Done when:** Checkpoint 2 is committed with a genuinely passing, runnable state behind it.

### #30: Noon Curveball Received: Track 2 — "Graph Is Evidence, Not an Oracle"
Priority: 🟥 CORE
Type: Discuss
Light: 🛑 NOON-GATE
Milestone: Noon Curveball window (12:00–1:00)
Blocked by: #29

**Build:** Read the received card in full (§13.5): the repo under review now contains dynamic dispatch, generated code, or reflection that static analysis can't fully resolve. Requirements — must not present incomplete Graph relationships as certain; must identify when analysis may be partial; must provide a safe fallback/verification path; existing behavior for fully-resolved code must continue unchanged; at least one test must use the supplied partial-analysis fixture; users/agents must be able to tell apart confirmed evidence, heuristic/incomplete evidence, and claims needing source/test verification. Close the current agent session. **Do not implement yet.**
**Watch out:** the card is numbered "TRACK 2" — the mapping from Entire's own Track E2 to this card number is an inference (two signals: the number match, and the content being written for "your Graph-powered experience"), not a stated fact. Confirm with a mentor at the earliest opportunity rather than building three hours of response against an unconfirmed guess.
**Test:** confirm the card's six requirements are written down verbatim, and a mentor has been asked to confirm the track mapping.
**Done when:** the constraint is fully read, recorded, and (ideally) mentor-confirmed, with zero implementation started.

---

## Phase 8 — Curveball Response: Coverage-Confidence (Track 2)

### #31: Fresh-Session Reconstruction via Pre-Noon Checkpoint
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #30

**Build:** Start a fresh agent session and have it reconstruct project state specifically from the pre-noon Checkpoint (#29) — intent, architecture, completed work, and open risks — plus current Graph context, before any editing begins.
**Watch out:** this is the guide's mandatory step before touching code post-Curveball; skipping straight to implementation from the old session's context defeats the reconstruction requirement.
**Test:** confirm the fresh session can describe intent/architecture/completed-work/open-risks correctly without being fed the old session's raw context directly.
**Done when:** the fresh session has genuinely reconstructed context from Checkpoint/Graph alone.

### #32: Pre-Edit `graph impact` Run — Scored Gate
Priority: 🟥 CORE
Type: Research
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #31

**Build:** Run `graph impact` (via #21's helper) on the dynamic-dispatch/reflection/codegen area Track 2's constraint affects, before making any change — the **second required graph demonstration** (§3.3), fired for real.
**Watch out:** per §13.5, the organizers' rubric explicitly downgrades to the partial band if graph use happens *after* implementation, or if the adaptation is only explained verbally rather than captured in Checkpoint context. This is a scored, non-negotiable gate now, not just good practice.
**Test:** confirm the impact run produced a concrete, reviewed set of affected entities before any code was touched in response.
**Done when:** the impact analysis is run, reviewed, and recorded as having happened strictly before implementation.

### #33: Coverage-Confidence Signal Implemented (full / partial / unknown)
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #7, #32

**Build:** Implement the §10.2 coverage-confidence signal, attached to every entity reconciliation touches, not just changed ones: `full` (default for any entity with at least one `extracted`-class edge, incoming or outgoing, in the snapshot); `partial` (zero `extracted` edges despite non-trivial code presence, or a region known to degrade under static analysis); `unknown` (the snapshot/diff call errored, timed out, or returned partial for that entity's file/module). Use whichever detection mechanism #7 confirmed is real — the native capabilities/snapshot marker if it exists, otherwise the zero-edge heuristic — with `treat_zero_edge_as: partial` as the conservative default per the curveball's "never present incomplete as certain" requirement.
**Watch out:** this signal is **orthogonal to the existing edge tri-state** and must never be conflated with it — tri-state classifies confidence in a relationship the graph recorded; coverage-confidence classifies whether the graph's silence about an entity means anything at all.
**Test:** feed an entity with a real extracted edge (expect `full`), an entity with zero edges despite real code (expect `partial`), and a forced snapshot error (expect `unknown`); confirm each classifies correctly.
**Done when:** every entity reconciliation touches carries a correctly-computed coverage value.

### #34: `unverifiable_coverage` Tier & Never-Silently-Upgrade Rule Wired
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #24, #33

**Build:** Wire Stage 4 to check coverage-confidence *before* assigning a tier: an entity with `coverage: partial` or `unknown` can never land in `undeclared_scope_creep` on the strength of an absent path alone — route it instead to a new `unverifiable_coverage` tier, carrying the same file/line/reason data plus the coverage reason. An entity with `coverage: full` proceeds through the existing tiers exactly as before.
**Watch out:** this is squarely a "what counts as drift" change, which the §13.1 firewall already scoped to Stage 4 + config — if this touches Stage 1–3 extraction or Stage 5's pure-function property, something has gone wrong.
**Test:** feed an entity that would previously have been flagged as scope creep but has `coverage: partial`; confirm it now lands in `unverifiable_coverage`. Separately confirm a fully-resolved scope-creep case from #24 still tiers exactly as before.
**Done when:** the never-silently-upgrade rule holds in both directions.

### #35: Baseline Manual-Review Fallback Flag Wired
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #34

**Build:** Flag every `unverifiable_coverage` entry with `verification_path: "manual_review_recommended"` as the baseline, always-present state — must work correctly with or without the Databricks module.
**Watch out:** this baseline **is** the actual safety guarantee the curveball requires — it must pass its test with `calibration_enabled: false` and `coverage_corroboration_enabled: false`.
**Test:** run the pipeline with Databricks fully disabled; confirm `unverifiable_coverage` entries still carry a populated `verification_path`.
**Done when:** the fallback flag works with zero Databricks dependency.

### #36: Schema Extended — `coverage_confidence` on Every Tier + New Tier + Summary Count
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #26, #35

**Build:** Extend `verdict.json` per §8.1's v5 fields: `coverage_confidence` on every entity across all six tiers (`confirmed`, `declared_unimplemented`, `expected_blast_radius`, `undeclared_scope_creep`, `advisory_low_confidence`, and the new `unverifiable_coverage`), the new tier itself with `coverage_reason`/`verification_path`/`corroboration`, and a new `summary.unverifiable_coverage_count`.
**Watch out:** never conflate `unverifiable_coverage` with `advisory_low_confidence` — one is "the graph has an edge but isn't sure about it," the other is "the graph may not have tried to record an edge here at all."
**Test:** run the full pipeline on a fixture with entities in all six tiers; confirm every entity carries `coverage_confidence` and the new fields populate correctly.
**Done when:** the extended schema is live and every downstream renderer/test consumes it correctly.

### #37: Config Extended — `coverage:` Block + New Verdict Threshold
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #9, #33

**Build:** Add the §8.2 v5 config additions: `verdict_thresholds.review_required_if_unverifiable_coverage_gte` (on by default), the new `coverage:` block (`detect_via: capabilities_api | heuristic_zero_edge` — set to whichever #7 confirmed is real; `treat_zero_edge_as: partial`; `fallback_verification_path: manual_review_recommended`), and `databricks.coverage_corroboration_enabled: false` as the new default.
**Test:** confirm the config parses and Stage 4 correctly reads `detect_via` to choose its detection mechanism.
**Done when:** the extended config is checked in and live.

---

## Phase 9 — Renderers, Fixture Test & Final Graph Demo (three-way split, mandatory)

### #38: Terminal Renderer Updated for Three-Way Split
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #27, #36

**Build:** Update the existing terminal renderer so confirmed structural evidence, heuristic/incomplete evidence (the `inferred`/`ambiguous` tri-state), and `unverifiable_coverage` claims each get a visually distinct treatment.
**Watch out:** this is a direct, scored requirement from the Track 2 card — a shared "risky" bucket with a tooltip does not satisfy it.
**Test:** render a fixture with entries in all three visual categories; confirm a reviewer can tell them apart at a glance.
**Done when:** all three states are visually distinct in the terminal output.

### #39: Markdown Renderer Built (Three-Way Aware from Start)
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #36

**Build:** Build `VERIFY_REPORT.md` as a pure function of the extended `verdict.json`, implementing the three-way split from the start — this renderer didn't exist before the curveball, so there's no retrofit.
**Test:** run against the same fixture as #38; confirm tier counts and the three-way split match the terminal renderer.
**Done when:** `VERIFY_REPORT.md` generates correctly with the three-way split built in.

### #40: Dashboard Renderer Built (Three-Way Aware from Start)
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #36

**Build:** Build `dashboard.html` — the single static file that `fetch()`es `verdict.json` and renders a radial/force graph, color-coded by edge class — with the three-way visual split built in from the start.
**Watch out:** keep it a pure function of `verdict.json` with no server dependency, per the clean-checkout requirement in the submission checklist.
**Test:** load the same fixture; confirm all three visual states render distinctly with no other running service required.
**Done when:** the dashboard renders correctly and satisfies the three-way split with no server dependency.

### #41: Renderer Purity + Three-Way Distinctness Test
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #38, #39, #40

**Build:** A test confirming all three renderers agree on tier counts (including the new `unverifiable_coverage` count) from the same `verdict.json`, and each independently satisfies the three-way visual distinctness requirement.
**Test:** run all three against several fixtures spanning all six tiers; diff reported counts and confirm each visually separates the three states.
**Done when:** all three renderers agree on counts and each passes the distinctness check.

### #42: Mandatory Partial-Coverage Fixture Test
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #34, #36

**Build:** Run the supplied partial-analysis repository fixture through the full pipeline and assert: (1) no entity in the affected region is misfiled into `undeclared_scope_creep`; (2) every affected entity carries `coverage_confidence: partial` (or `unknown` if the fixture triggers a snapshot error) and lands in `unverifiable_coverage`; (3) `verification_path` is populated at minimum with `manual_review_recommended`; (4) a fully-resolved entity elsewhere in the *same* fixture still reconciles normally.
**Watch out:** this test is mandatory per the curveball card, not optional.
**Test:** run it and confirm all four assertions pass against the *actual supplied* fixture, not a hand-built substitute.
**Done when:** all four assertions pass.

### #43: Final Semantic-Diff Analysis Wired
Priority: 🟥 CORE
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #36, #32

**Build:** Run `entire graph diff --base <base> --head <checkpoint-id>` against the final post-Curveball implementation, feeding the final `verdict.json` — the **third required graph demonstration**, now including coverage data.
**Test:** run against the actual final checkpoint; confirm the resulting `verdict.json` reflects the complete post-Curveball pipeline.
**Done when:** the submitted implementation's final `verdict.json` is generated from a real end-to-end run against the extended schema.

### #44: Checkpoint 3 — Response to Curveball Recorded
Priority: 🟥 CORE
Type: Discuss
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #42, #43

**Build:** Write the third required Entire Checkpoint: name the received card (**Track 2 — "Graph Is Evidence, Not an Oracle"**), the assumption it invalidated (a zero-edge result was being treated as confident scope creep regardless of *why* the graph was silent), the fix (coverage-confidence signal + `unverifiable_coverage` tier + baseline fallback), and state honestly that **none of the five §13.2 ranked predictions named this exact constraint** — noting rank 4's partial relevance (the Issue #32 graceful-degradation instinct generalized usefully, but the coverage-confidence signal itself is new work).
**Watch out:** house style is to flag prediction misses rather than retrofit them to look foreseen — overclaiming foresight is a bigger credibility risk than admitting the miss.
**Test:** re-read cold; confirm it names the card, the invalidated assumption, the fix, and the honest prediction-accuracy note.
**Done when:** Checkpoint 3 is committed with all four elements present.

---

## Phase 10 — Open-Source Augmentation: Curveball Tier (v6, do first if time allows only one addition)

### #45: Semgrep Integration for Coverage Detection
Priority: 🔴 AUGMENT-CRITICAL
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #33

**Build:** Run Semgrep's prebuilt rule packs for dynamic-dispatch/reflection/codegen signatures (`getattr`, `eval`, `__getattr__`, reflection APIs, decorator-generated methods, per language) against entities the zero-edge heuristic (#33) flags, promoting `coverage: partial` from an inferred proxy to a **detected finding** when a pattern actually matches.
**Watch out:** per §20.5, this is the single highest-leverage addition in the whole augmentation pass because it strengthens the exact mechanism the curveball is scoring, not an adjacent one — build this before any other optional item if time only allows one.
**Test:** run Semgrep against a known dynamic-dispatch snippet and a known genuinely-disconnected snippet; confirm only the former gets a detected-finding upgrade.
**Done when:** partial-coverage entities carry a real detected pattern match where one exists, not just an absence heuristic.

### #46: Joern Evaluation (heavier alternative — only if Semgrep isn't convincing)
Priority: 🟡 AUGMENT-OPTIONAL
Type: Research
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline) — cut-safe
Blocked by: #45

**Build:** Only if Semgrep's output doesn't feel convincing in the demo and 30+ minutes remain, evaluate Joern's code-property-graph modeling of reflective/dynamic call edges as a heavier alternative.
**Watch out:** do not attempt both Semgrep and Joern — Semgrep is the better time-to-value trade in almost every case; this is a fallback for a specific "still not convincing" scenario.
**Done when:** either Joern replaces Semgrep for a stronger demo, or the team confirms Semgrep alone is convincing and skips this entirely.

---

## Phase 11 — Databricks Module: Core Chain (opt-in — only if #2's decision was yes, and only once the core response is solid)

### #47: Calibration Set + I-CALM Prompt Structure
Priority: 🟦 DATABRICKS
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #26

**Build:** Compile a small calibration set of known true/false heuristic edges, and run the extraction prompt with the I-CALM structure — full reward for a correct edge, heavy penalty for incorrect, partial credit for abstaining. Return each proposed edge with a verbal confidence score in `[0,1]`.
**Watch out:** document that calibration labels are synthetic/injected, not real production revert data (§12.8) — state this honestly in `BUILDATHON.md`.
**Test:** run the I-CALM prompt on the calibration set; confirm every edge returns a verbal confidence score.
**Done when:** a labeled calibration set exists and the prompt reliably returns scored edges.

### #48: Conformal Risk Control Threshold
Priority: 🟦 DATABRICKS
Type: Research
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #47

**Build:** Using Conformal Risk Control on #47's calibration set, compute the smallest threshold `λ̂` such that the expected false-positive rate stays below the target `α` (default 0.05).
**Test:** verify the computed `λ̂` actually holds the false-positive rate under `α` on a held-out slice.
**Done when:** `λ̂` is computed, verified, and written into the runtime config.

### #49: MAPIE Conformal Prediction Library Swap-In
Priority: 🟡 AUGMENT-OPTIONAL
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline) — cut-safe
Blocked by: #48

**Build:** Swap the hand-rolled CRC threshold computation for **MAPIE**, a tested open-source conformal prediction library, as a near drop-in replacement (§20.2).
**Watch out:** only do this if #48's hand-rolled version is already working and time is genuinely spare — the payoff is defensibility ("show me the math") against a Databricks judge, not new functionality.
**Test:** confirm MAPIE produces a threshold with the same statistical guarantee on the same calibration set.
**Done when:** either MAPIE or the hand-rolled version is live, and it's clearly documented which.

### #50: Online Gating — Bypass + Epistemic Abstention
Priority: 🟦 DATABRICKS
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #48

**Build:** Deterministic edges bypass the gate entirely; advisory edges compared against `λ̂` — pass sets `conformal_gate_passed: true`, fail routes to `pending_verification` rather than silently trusting or dropping.
**Watch out:** epistemic abstention is the correct behavior below threshold — never a silent drop.
**Test:** feed one high-confidence and one low-confidence advisory edge; confirm correct routing.
**Done when:** gating behaves correctly for deterministic, passing-advisory, and abstaining-advisory cases.

### #51: AI Search Corroboration Index + Query
Priority: 🟦 DATABRICKS
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #50

**Build:** Build a Databricks AI Search index over the repo's own docstrings, comments, and commit messages, and query it in hybrid mode for every `pending_verification` edge with the edge's hypothesis as a natural-language question.
**Watch out:** treat this as a capped, occasional call — only `pending_verification` edges reach it — which is also what keeps it inside Free Edition's limits (§3.4).
**Test:** run one query for a real pending edge; confirm RRF-ranked snippets return with Unity Catalog governance carrying through.
**Done when:** corroboration queries reliably return ranked evidence for a pending edge.

### #52: Local Corroboration Fallback (sentence-transformers/FAISS or rank_bm25)
Priority: 🟡 AUGMENT-OPTIONAL
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline) — cut-safe
Blocked by: #51

**Build:** If Databricks AI Search is proving flaky, stand up a local hybrid-retrieval fallback using **sentence-transformers (`all-MiniLM-L6-v2`) + FAISS or `rank_bm25`**, working identically with or without Databricks (§20.2).
**Watch out:** don't build this for redundancy alone — only reach for it if AI Search is actually causing problems; a second retrieval stack under time pressure costs more than it protects.
**Test:** confirm the local fallback returns comparably relevant snippets on the same corroboration queries.
**Done when:** either AI Search alone is working, or the local fallback is live as a genuine substitute — not both maintained in parallel.

### #53: Calibrated Fields Written Back to `verdict.json`
Priority: 🟦 DATABRICKS
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #51

**Build:** Populate `verbal_confidence`, `conformal_gate_passed`, and `corroboration.verified`+evidence snippet on every `advisory_low_confidence` entry.
**Test:** run the full calibration chain on a `verdict.json` with advisory entries; confirm every field populates when enabled.
**Done when:** the schema's Databricks-specific fields are genuinely live end-to-end.

### #54: Fallback Tier — Logistic Regression in MLflow
Priority: 🟦 DATABRICKS
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline) — cut-safe if #48–#53 are running
Blocked by: #47

**Build:** If time runs out before the full CRC/AI-Search chain works, train a plain logistic regression in MLflow on calibration features (edge class, edge type, language, churn size) against outcome labels, writing `calibrated_confidence` back into each advisory edge. Set `gate.mode: fallback`.
**Watch out:** still a legitimate, schema-compatible Databricks entry — don't treat it as a failure state if it's what ships.
**Test:** confirm the fallback path populates `calibrated_confidence` without touching the schema shape from #26.
**Done when:** either the full chain or this fallback is live — never neither.

### #55: Dashboard Ranks Advisory Edges by Calibrated Signal
Priority: 🟦 DATABRICKS
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #53, #40

**Build:** Update the dashboard renderer to rank advisory edges by the real calibrated signal instead of an undifferentiated flat list.
**Test:** load a `verdict.json` with several advisory edges at different calibrated confidence levels; confirm correct ordering.
**Done when:** the dashboard visibly reflects calibration quality.

### #56: Coverage Corroboration — Databricks-Enhanced Tier for `unverifiable_coverage`
Priority: 🟦 DATABRICKS
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Afternoon (by deadline)
Blocked by: #34, #51

**Build:** If `coverage_corroboration_enabled: true`, query the same AI Search index already built for #51 a second way (§12.15): for each `unverifiable_coverage` entity, ask whether the repo's own documentation or commit history describes what it dispatches to or is dispatched from. A hit populates `corroboration.verified`/`evidence_snippet` exactly as the advisory-edge path does; a miss leaves `verification_path` unchanged.
**Watch out:** this is an enhancement of the mandatory baseline (#35), **never a substitute for it** — the baseline must keep working with this disabled. No new index, no new infrastructure — state this distinction explicitly if asked; it's more honest than implying Databricks is required for curveball compliance.
**Test:** run one `unverifiable_coverage` entity through corroboration with a real doc/comment hit, and one with no hit; confirm the former gets verified evidence and the latter keeps the manual-review flag.
**Done when:** the enhancement works and is provably non-required for the baseline safety guarantee.

---

## Phase 12 — Community-Aware Scope-Creep (optional enhancement)

### #57: Community-Aware Scope-Creep Framing (Leiden Clustering)
Priority: ⬜ STRETCH
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Stretch
Blocked by: #24

**Build:** Apply Leiden-style community detection over the snapshot graph — via **`leidenalg` + `python-igraph`**, the actual pinned implementation (§20.2) — so `undeclared_scope_creep` entries can be reported by subsystem ("this change also touched 3 entities in the Billing community") rather than by unrelated-looking file paths (§10.1).
**Watch out:** this needs a graph large and connected enough for clustering to produce meaningful, non-trivial communities — a small hackathon demo repo may not have one. Attempt only once the base pipeline (through #44) is solid, and state the caveat plainly if the demo repo is too small, rather than forcing a result that isn't real.
**Test:** run clustering on the demo repo's snapshot; confirm communities are non-trivial before using this in the demo.
**Done when:** either community-labeled scope-creep entries are genuinely informative, or the caveat is stated and this is skipped.

---

## Phase 13 — Databricks Stretch: Lakebase v3 (only if #8 cleared clean, otherwise skip entirely)

### #58: Zero-Copy Counterfactual Sandbox — What-If Verdicts
Priority: ⬜ STRETCH
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Stretch
Blocked by: #8, #26

**Build:** Run Stage 3–4 reconciliation against a disposable Lakebase branch instead of the canonical store, and let a reviewer replay the same diff with a different `fidelity.config.yaml` (e.g. `blast_radius_hops: 3`) to get a second `verdict.json` in seconds without re-running graph extraction.
**Watch out:** this requires materializing the dependency graph inside Lakebase/Postgres — a real change to Stage 3's architecture, not a Stage 6 bolt-on. Name this to a judge as a deliberate exception to the firewall, not a hidden one.
**Test:** change `blast_radius_hops` on a live branch and confirm a new verdict appears within seconds, with the canonical store untouched.
**Done when:** a live "widen the blast radius" demo works end-to-end on a disposable branch.

### #59: Live Verdict Ledger
Priority: ⬜ STRETCH
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Stretch
Blocked by: #8, #53

**Build:** Write every `verdict.json` produced during the day as a row into a `verdicts` table in Lakebase; LTAP transcodes it into Delta automatically for cross-session trend queries.
**Watch out:** this is the strongest version of "why does this need Databricks" — only claim it if actually live.
**Test:** produce several verdicts across sessions; confirm a trend query returns correct historical data.
**Done when:** the ledger accumulates real verdicts and supports at least one cross-session query.

### #60: Lakebase Curveball-Reconstruction Memory
Priority: ⬜ STRETCH
Type: Prototype
Light: 🌇 AFTERNOON
Milestone: Stretch
Blocked by: #8, #31

**Build:** Persist a structured, queryable snapshot of each Checkpoint's key decisions into Lakebase, so the fresh-session reconstruction step (#31) pulls from a fast structured store instead of re-parsing a raw transcript cold.
**Watch out:** only counts if it's genuinely what #31 used, not built after the fact to describe #31 retroactively.
**Test:** clear the fresh session's context and confirm it can reconstruct decisions correctly by querying Lakebase alone.
**Done when:** #31's reconstruction step demonstrably reads from this store.

---

## Phase 14 — Proof, Checkpoint 4 & Submission

### #61: Checkpoint 4 — Final Implementation & Verification Recorded
Priority: 🟥 CORE
Type: Discuss
Light: 🏁 FINAL
Milestone: Final 20 min
Blocked by: #41, #43

**Build:** Write the fourth required Entire Checkpoint: final implementation state, verification evidence (tests passing, all three renderers agreeing, the final semantic-diff run), and known limitations.
**Test:** re-read cold; confirm it references real test results and the actual final `verdict.json`.
**Done when:** Checkpoint 4 is committed with real verification evidence attached.

### #62: Three Graph Demonstrations Rehearsed by Name
Priority: 🟥 CORE
Type: Discuss
Light: 🏁 FINAL
Milestone: Final 20 min
Blocked by: #14, #32, #43

**Build:** Rehearse narrating, by name, the exact moment each required graph demonstration happened: the snapshot lookup (#14), the pre-change impact analysis (#32), and the final semantic-diff (#43).
**Watch out:** the rubric line literally asks for these to be named explicitly.
**Test:** say all three, in order, with which ticket/stage each corresponds to, in under a minute.
**Done when:** every team member can name and point to all three live.

### #63: Challenge Questions & Positioning Rehearsed Cold
Priority: 🟥 CORE
Type: Discuss
Light: 🏁 FINAL
Milestone: Final 20 min
Blocked by: #3, #44

**Build:** Rehearse the final versions of §3's differentiation answers (including the Graphify/§1.5 answer) plus §19's positioning lines, updated with real specifics: what the curveball actually was, that none of the five predictions called it, and — if #45 landed — the line that coverage detection is a real Semgrep pattern match, not just an absence heuristic.
**Test:** every team member answers "isn't this just `review`?", "why not Genie Ontology?", "isn't this just Graphify?", and "what did the curveball actually change?" cleanly, cold.
**Done when:** the answers are rehearsed with the final, real details.

### #64: `BUILDATHON.md` Written
Priority: 🟥 CORE
Type: Prototype
Light: 🏁 FINAL
Milestone: Final 20 min
Blocked by: #12, #44, #61, #62

**Build:** Fill the exact §17 outline: project name, one-sentence summary, problem/user/track rationale (with §1.4/§1.5 differentiation folded in), architecture, the three named graph demonstrations, the Curveball section naming **Track 2 — "Graph Is Evidence, Not an Oracle"**, honestly stating none of §13.2's five predictions called it, and naming the invalidated assumption (zero-edge ⇒ scope creep) and the fix; links to all four Checkpoints; setup/run/test instructions; Databricks use and limitations if opted in; and known limitations/next steps — including which §20 open-source additions actually made it in (e.g. Semgrep-backed coverage detection) versus which stayed roadmap-only.
**Watch out:** free of secrets, readable by someone who wasn't in the room.
**Test:** have someone outside the build team read it and confirm they could set up and run the project from it alone.
**Done when:** `BUILDATHON.md` is complete against the exact §17 outline, with v6's curveball and augmentation content included.

### #65: Fallback Demo Recording Captured
Priority: 🟥 CORE
Type: Prototype
Light: 🏁 FINAL
Milestone: Final 20 min
Blocked by: #40

**Build:** Record a full successful run — dashboard open, a real verdict rendering with all three visual states visible, the Curveball response demonstrable — as a fallback if the live demo fails during judging.
**Test:** play the recording start to finish; confirm it clearly shows the pipeline working end-to-end, including the three-way split.
**Done when:** a usable, local backup recording exists.

### #66: Final Submission Checklist Run
Priority: 🟥 CORE
Type: Discuss
Light: 🏁 FINAL
Milestone: Deadline
Blocked by: #63, #64, #65

**Build:** Run the §18 checklist verbatim: final commit SHA matches the submission; clean-checkout launch; all four Checkpoints open and explain their milestone; graph evidence and final semantic-diff recorded and named; `BUILDATHON.md` complete and secret-free; critical + Curveball tests passing (including #42's mandatory fixture test); Databricks links/data notes included if opted in; demo owner (#2) can sign in and run the critical path; fallback recording exists; submitted before the actual confirmed deadline (#1).
**Watch out:** confirm the real deadline from #1 one more time — don't let an augmentation ticket (#46, #49, #52, #57) eat the last minutes before it.
**Test:** go through every checklist line as a literal pass/fail, not a skim.
**Done when:** every line is checked and the submission is in before the confirmed deadline.

---

## Future Work (out of scope for this buildathon build)

These aren't rejected ideas — they're real extensions that don't fit this build's time/complexity budget, or that only conditionally got built above. Worth naming if asked, worth returning to after the event:

1. **Plan-Trust Ledger** — persist (claim, outcome) pairs across sessions to learn how much to trust a given agent/model's plans over time. Deferred unless #59 was built, in which case a live version already exists.
2. **Checkpoint Court** — a second, adversarial Agent Bricks agent that tries to falsify each verdict claim before acceptance, logged via Unity AI Gateway as a replayable trial record. Never build here.
3. **Forensic Replay** — Delta time-travel on versioned, ULID-keyed snapshots to reconstruct what an agent actually saw at a past decision. #58's what-if verdicts are a narrower, live-buildable version if #8 cleared.
4. **Subagent Reasoning Capture** — intercept subagent reasoning before condensation (Issue #2058) and reconcile each subagent's sub-intent against its scoped graph diff.
5. **Compliance reconciliation via Unity Catalog lineage** — resolve a code change through data lineage to "which regulated fields/attestations does this touch."
6. **SCIP export alignment** — once Entire's RFD 0006 SCIP-export convergence ships, expose `verdict.json` findings through that standard channel instead of a bespoke format.
7. **Fidelity as an MCP-exposed self-check tool** (§11.1) — letting another coding agent call `verify-intent` on itself before handing a diff to a human. Real and low-effort once `verdict.json` is stable, but explicitly placed in "not-now" territory by §20.3: additive scope this late in the day, not curveball-scoring.
8. **CodeQL / Kythe** — both real, both heavier integrations than Semgrep for the identical §10.2 coverage-detection payoff (§20.3). Good `BUILDATHON.md` roadmap lines, not same-day adds.

**Explicitly out of scope, not just deferred:** fleet/cross-repo *structural* graph analytics (`entire-graph` enforces single-store isolation, PR #169); any feature depending on `TESTS`, `HANDLES_ROUTE`, or `HANDLES_GRPC` edges (unverified in the source engine); Genie Ontology as the graph store (§1.2). **And, stated plainly per §20.3: never adopt Graphify itself as a dependency** — §1.5 exists specifically to differentiate Fidelity's reconciliation-verdict model from Graphify's context-retrieval model, and folding it in as a component would blur that positioning rather than strengthen it.

---

*This build map is a companion to `FIDELITY_FINAL_BIBLE_v6.md`. The Bible explains WHY; this map tells you WHAT to build, in WHAT order, how to test each piece, which window and priority tag each belongs to, and what's genuinely future work. When a step and the Bible disagree, the Bible is the source of truth on intent — but follow this map's ticket order, Light tags, and Test steps for execution.*