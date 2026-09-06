# Fidelity — Phase 0 Brief

This is the pre-event decision record for Fidelity. It separates confirmed
design decisions from event-specific choices that must be made by the team at
kickoff.

## 1. Timeline confirmation — pending kickoff

The supplied planning material contains a deadline conflict:

- Participant Guide: submit by 3:00 PM IST.
- Event website: code freeze/submission at 4:00 PM IST.

The live spoken kickoff announcement is authoritative. The team must record
the confirmed deadline here before 9:30 AM IST.

**Confirmed deadline:** 3:00 PM IST (confirmed by the team at the Noon Curveball gate)

**Person who confirmed it:** team, verbally, at the 12:00 curveball briefing

## 2. Track, Databricks, repository, and demo ownership — pending team decision

| Decision | Current value |
| --- | --- |
| Entire track | E2 — Build with Graph Intelligence |
| Databricks opt-in | **Yes** — a special prize is offered for the best Databricks use. See the honest environment caveat below. |
| GitHub fork target | Deferred — working entirely locally for the remainder of the build per team direction; no fork/mirror work was done post-curveball. |
| Entire mirror region | India |
| Demo/deployment owner | Not yet named |

**Honest Databricks environment caveat, recorded here rather than discovered
at demo time:** this machine has no `databricks` CLI, no `~/.databrickscfg`,
and no `mlflow` installed, and no workspace credentials were available inside
the remaining build window. The opt-in decision is real, but the mechanism
actually reachable in this window is §12.7's fallback tier (a plain, local
scoring computation, schema-compatible with the full I-CALM/CRC/AI-Search
chain) — attempted only after the mandatory Curveball response (Track 2) is
solid, per house priority. If it does not land before the deadline, that is
stated plainly in `BUILDATHON.md` as environment-blocked, not silently
dropped.

## 3. Differentiation answers

### “Isn’t this just the review/explainskill/what-happened skills?”

Those tools retrieve and narrate why a change exists by reading session
context. Fidelity checks whether the claimed intent matches the structural code
change. Its output is a tiered, inspectable reconciliation backed by graph
evidence, rather than a narrative judgment. The two can work together: a
review tool can use Fidelity’s verdict as structured evidence.

### “Why not Genie Ontology?”

Fidelity should not force a code graph into Genie Ontology. The supplied plan
identifies its snippet limit and SQL-oriented traversal as a poor fit for code
graph structure. Entire Graph remains the structural source of truth;
Databricks, if selected, is an optional confidence/corroboration and reporting
layer rather than the graph store.

### "Isn't this just Graphify's tri-state model?"

Graphify's EXTRACTED/INFERRED/AMBIGUOUS tagging routes an agent's *attention*
during exploration — it helps an agent decide how much to trust an edge while
it works. Fidelity uses the same kind of signal to drive a deterministic
*verdict* after the work is done: reconciliation against a stated claim, not
retrieval to support one. The tri-state idea itself is not Fidelity's
differentiator — the two tools are converging on the same taxonomy
independently — the differentiator is what happens after: a persisted,
graph-verified reconciliation a reviewer can trust as reproducible, versus a
context aid for further exploration. (Bible v6 §1.5.)

## 4. Curveball briefing

The highest-probability planned pivot is a unified query layer joining graph
and checkpoint/session access. Fidelity’s response is the adapter seam in
`internal/fidelity/adapter.go`: capture, graph symbol loading, and semantic
change retrieval are all behind one interface.

Other anticipated pivots are foreign-checkpoint review, subagent PII
redaction, multi-repository data-flow requests, and incremental syncing. They
are not part of the initial core implementation. Any new data source must enter
through the adapter and preserve the graph provider’s no-egress, single-repo
contract.

## Phase 0 completion checklist

- [x] Record the official deadline after kickoff. (3:00 PM IST)
- [x] Record the Databricks opt-in decision. (Yes, with an honest environment caveat above.)
- [ ] Record the fork target and India mirror setup owner. (Deferred — local-only for the remainder of the build.)
- [ ] Name the demo/deployment owner.
- [x] Record the judge-facing differentiation answers, including §1.5 (Graphify).
- [x] Record the curveball briefing and adapter safeguard.
