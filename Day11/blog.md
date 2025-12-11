# Day 11: Reactor Routing Notes

Day 11 felt like a network-diagnostics sprint: enumerate every path through a directed graph, first from `you` to `out`, then from `svr` to `out` while forcing the walk to touch both `dac` and `fft`. Here is how I reasoned about the problem and built the solver.

## Puzzle Recap

- Each line describes a device and every outbound connection it fans out to.
- Part 1 counts all simple paths from `you` to `out`.
- Part 2 counts paths from `svr` to `out` that also visit `dac` and `fft` (order does not matter, but they must both appear).
- Inputs may contain comments, blank lines, and orphan nodes, so the parser needs to be forgiving.

## Parsing and Graph Shape

The parser tokenizes `name: targets` pairs, ignores comments, and stores a dense adjacency list. Every mentioned node is inserted into the map even if it lacks outgoing edges to avoid missing leaf vertices. I also raise the scanner buffer cap to cope with large files.

## DAG-First Strategy

Simple path counting on arbitrary directed graphs is exponential, but the puzzle inputs appear to form DAGs. I still guard against cycles by first pruning and topologically sorting the relevant subgraph:

1. DFS/BFS to mark nodes reachable from the chosen start.
2. Reverse BFS to keep only nodes that can also reach the goal.
3. Kahn's algorithm on that intersection. If the queue exhausts all nodes we get a topological order; otherwise we detected a cycle.

With an order in hand, counting paths is just dynamic programming:

```
ways[start] = 1
for node in topoOrder:
    for each neighbor:
        ways[neighbor] += ways[node]
```

If the pruned subgraph has a cycle, I fall back to classic DFS with a `visited` set to enumerate simple paths. That worst case is exponential, but it only triggers when the optimized route is impossible.

## Must-Visit Constraints (Part 2)

For the `svr`→`out` request I extend the DAG DP to track whether we have hit `dac` and `fft`. Each node carries a 2-bit mask; state `(node, mask)` stores how many paths reach that node having already visited the required devices indicated by the mask. Propagating through the topological order is the same as part 1, just with four masks instead of one count. The fallback DFS also threads two booleans through the recursion when cycles force enumeration.

## Complexity and Testing

Let `n` be the number of relevant nodes and `m` their edges.

- Forward/backward reachability plus Kahn's algorithm: `O(n + m)` time, `O(n + m)` space.
- DAG DP for either part: `O(n + m)` time, `O(n)` space.
- Cycle fallback: exponential time in the number of nodes on a cyclic component, but only activated if the input truly requires it.

Unit tests run the two samples from the puzzle statement to lock in both the `you` paths and the constrained `svr` paths. `go test ./...` covers all days, including this graph-heavy one.

The end result is a solver that is fast on the intended DAG inputs, yet still correct if a stray cycle sneaks into the wiring diagram.
