# Day 8: Playground Wiring Diary

Today I left the teleporter labs behind and walked straight into the Day 8 puzzle: wiring up a playground full of floating junction boxes. Here is how I modeled the problem, implemented the solver, and reasoned about its complexity.

## Puzzle Recap

- Every junction box sits at an integer coordinate `(x, y, z)`.
- I need to keep connecting the closest pair of boxes that are not yet directly linked.
- For part 1, after exactly 1000 connections I must multiply the sizes of the three largest connected components.
- For part 2, I keep wiring until everything forms one component; the answer is the product of the `x` coordinates of the final pair that completes the single circuit.

The sample input only asks for the first 10 connections and is perfect for verifying the implementation.

## Parsing and Data Model

Each line in `input.txt` is just three comma-separated integers. I map those into a simple `point` struct:

```go
 type point struct {
     x, y, z int64
 }
```

While reading the file I ignore blank lines and stash all points in a slice. Any malformed coordinate line immediately returns an error so I never run the solver on corrupted data.

## Generating Candidate Connections

The problem says “always pick the closest unconnected pair”, which is exactly Kruskal’s algorithm on a complete graph:

1. Generate every pair of distinct points. For `n` points there are `n(n-1)/2` pairs.
2. Compute squared Euclidean distance `dx*dx + dy*dy + dz*dz` to avoid floating point math.
3. Sort the list of pairs by distance (breaking ties deterministically by the indices of the points).

This is the most expensive step; for the real input the quadratic blowup is unavoidable because the puzzle explicitly requires knowledge of the global ordering of edges.

## Tracking Circuits with Union–Find

I rely on a classic disjoint-set union data structure:

- Path compression and union-by-size deliver near-constant-time merges (`α(n)` inverse Ackermann).
- Every successful union decreases the component count.
- Component sizes are stored so I can quickly compute circuit sizes for part 1.

## Simulation Loop

With sorted pairs and union–find ready, the main loop becomes straightforward:

```go
connections := 0
for _, edge := range pairs {
    if connections < limit {
        connections++
    }
    merged := uf.union(edge.i, edge.j)

    if connections == limit && !gotPart1 {
        part1 = productOfTopThree(uf)
        gotPart1 = true
    }

    if merged && uf.components == 1 && !gotPart2 {
        part2 = points[edge.i].x * points[edge.j].x
        gotPart2 = true
    }

    if gotPart1 && gotPart2 {
        break
    }
}
```

- `limit` is 1000 for the real problem but I expose `SolveWithLimit` so the unit test can use `limit=10`.
- When I reach the limit, I scan the union–find roots, collect their sizes, sort descending, and multiply the largest three counts to produce part 1.
- For part 2 I watch for the first union that drops the component count to exactly 1; the ordered pair responsible for that merge gives me the final product of `x` coordinates.

## Complexity Discussion

Let `n` be the number of junction boxes. The dominating work is in generating and sorting all candidate edges:

- **Pair generation:** `Θ(n²)` pairs.
- **Sorting:** `O(n² log n²)` which is `O(n² log n)`.
- **Union–Find loop:** processes each pair once, so `Θ(n² α(n))`, effectively linear in the number of pairs.

So, overall time complexity is `O(n² log n)` and space complexity is `O(n²)` for storing the pair list. This matches the requirement to always consider the globally shortest remaining connection; there is no cheaper exact alternative under those rules.

## Testing and Validation

- The sample from the puzzle statement uses only the first 10 connections. My unit test runs `SolveWithLimit` with `limit=10` and verifies `part1=40` and `part2=25272`.
- The main solver calls `Solve` with the default `limit=1000` for the full input. Running `go test ./...` covers all existing days plus the new logic.

## Final Thoughts

Day 8 reads like minimum-spanning-tree Kruskal, but with a twist: instead of halting once the MST is complete, we snapshot the system after a fixed number of edges and again at the moment connectivity becomes global. Once the pair list is sorted, both answers fall out of the same simulation. The quadratic cost is hefty, yet perfectly manageable for the given input sizes and gives a clear, reliable solution.
