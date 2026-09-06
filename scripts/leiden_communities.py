#!/usr/bin/env python3
"""Build Map #57: community-aware scope-creep framing.

Reads an NDJSON relation-edge stream from `entire graph edges` and runs Leiden
community detection over it (leidenalg + python-igraph, the pinned Bible
implementation), so undeclared_scope_creep entries can be reported by
subsystem ("this change also touched 3 entities in the Billing community")
rather than by unrelated-looking file paths.

Structural DEFINES/CONTAINS edges are excluded, same rule as
internal/fidelity/coverage.go and reachability.go: every symbol has one
regardless of whether anything references it, so counting them would collapse
the whole graph into one trivial "same file" community.

Usage:
    entire-graph symbols --repo . --format ndjson > symbols.ndjson
    entire-graph edges --repo . --format ndjson \
        --relation CALLS,DATA_FLOWS,USES_TYPE,PARAM_TYPE,RETURNS_TYPE,\
READS_FIELD,WRITES_FIELD,ACCESSES,EXTENDS,IMPLEMENTS,INHERITS,OVERRIDES,\
CONSTRUCTS,ASYNC_CALLS \
        | python3 scripts/leiden_communities.py symbols.ndjson > communities.json

Output: {"communities": {symbol_id: community_index, ...},
         "community_labels": {community_index: "dominant/path/prefix"},
         "stats": {"nodes": N, "edges": M, "communities": K,
                   "largest_community_fraction": F, "singleton_fraction": S}}

ponytail: this is a standalone script invoked by a human or a later CI step,
not wired into the Go pipeline's hot path — a robust subprocess bridge with
missing-python fallback is real engineering that this stretch ticket's own
buildmap entry explicitly permits skipping ("state the caveat plainly ...
rather than forcing a result that isn't real"). Add wiring when this output
is actually load-bearing for a demo.
"""
import json
import sys
from collections import Counter

import igraph
import leidenalg

STRUCTURAL_TYPES = {"DEFINES", "CONTAINS"}


def load_file_paths(symbols_path):
    """id -> file_path, read from the symbol stream's own field (never the ID
    positionally — see snapshot-format.md's "fields are not escaped" note)."""
    paths = {}
    with open(symbols_path) as handle:
        for line in handle:
            line = line.strip()
            if not line:
                continue
            record = json.loads(line)
            if record.get("record_type") == "symbol" and record.get("file_path"):
                paths[record["id"]] = record["file_path"]
    return paths


def main():
    if len(sys.argv) != 2:
        print("usage: leiden_communities.py <symbols.ndjson> < edges.ndjson", file=sys.stderr)
        sys.exit(2)
    file_paths = load_file_paths(sys.argv[1])

    node_ids = {}  # id -> igraph vertex index
    edges = []
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        record = json.loads(line)
        if record.get("record_type") != "relation":
            continue
        if record.get("type") in STRUCTURAL_TYPES:
            continue
        from_id, to_id = record.get("from_id"), record.get("to_id")
        if not from_id or not to_id or from_id == to_id:
            continue
        for node in (from_id, to_id):
            if node not in node_ids:
                node_ids[node] = len(node_ids)
        edges.append((node_ids[from_id], node_ids[to_id]))

    if not node_ids:
        print(json.dumps({"error": "no non-structural relation edges found in input"}), file=sys.stderr)
        sys.exit(1)

    graph = igraph.Graph(n=len(node_ids), edges=edges, directed=False)
    graph.simplify(multiple=True, loops=True)
    partition = leidenalg.find_partition(graph, leidenalg.ModularityVertexPartition)

    id_by_index = {index: node_id for node_id, index in node_ids.items()}
    communities = {}
    community_members = {}
    for vertex_index, community_index in enumerate(partition.membership):
        node_id = id_by_index[vertex_index]
        communities[node_id] = community_index
        community_members.setdefault(community_index, []).append(node_id)

    # Label each community by its most common directory prefix, read from
    # each member symbol's own file_path field -- never sliced out of the
    # compound ID positionally, per snapshot-format.md's explicit warning
    # that ID fields are unescaped and a path segment can itself contain ':'.
    community_labels = {}
    for community_index, members in community_members.items():
        directories = Counter()
        for member in members:
            file_path = file_paths.get(member)
            if not file_path:
                continue
            directory = file_path.rsplit("/", 1)[0] if "/" in file_path else "(root)"
            directories[directory] += 1
        community_labels[community_index] = directories.most_common(1)[0][0] if directories else f"community-{community_index}"

    sizes = sorted((len(m) for m in community_members.values()), reverse=True)
    total = sum(sizes)
    singleton_count = sum(1 for s in sizes if s == 1)

    result = {
        "communities": communities,
        "community_labels": {str(k): v for k, v in community_labels.items()},
        "stats": {
            "nodes": len(node_ids),
            "edges": len(edges),
            "communities": len(community_members),
            "largest_community_fraction": round(sizes[0] / total, 4) if total else 0,
            "singleton_fraction": round(singleton_count / len(community_members), 4) if community_members else 0,
        },
    }
    json.dump(result, sys.stdout, indent=2)
    print()


if __name__ == "__main__":
    main()
