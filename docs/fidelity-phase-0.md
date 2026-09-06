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

**Confirmed deadline:** pending

**Person who confirmed it:** pending

## 2. Track, Databricks, repository, and demo ownership — pending team decision

| Decision | Current value |
| --- | --- |
| Entire track | E2 — Build with Graph Intelligence |
| Databricks opt-in | Pending team decision |
| GitHub fork target | Pending team decision |
| Entire mirror region | India |
| Demo/deployment owner | Pending team decision |

Databricks stays disabled in `fidelity.config.yaml` until the team explicitly
opts in. Its Lakebase stretch work is additionally gated on a live Free Edition
availability check.

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

- [ ] Record the official deadline after kickoff.
- [ ] Record the Databricks opt-in decision.
- [ ] Record the fork target and India mirror setup owner.
- [ ] Name the demo/deployment owner.
- [x] Record the judge-facing differentiation answers.
- [x] Record the curveball briefing and adapter safeguard.
