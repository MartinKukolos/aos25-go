# Day 10: Factory Startup Log

Today’s machines needed two entirely different tuning modes: binary indicator lights driven by toggle buttons, and integer joltage counters driven by the same buttons acting as incrementers. Both parts boiled down to solving small linear systems, but over different number systems and with different objective functions.

## Puzzle Recap

- Each line describes a machine: an indicator pattern, a list of buttons (each toggles specific indices), and a list of joltage targets.
- Part 1 (lights): starting from all-off, find the minimum number of button presses that produces the desired on/off pattern. Buttons toggle their listed bits modulo 2.
- Part 2 (joltage): ignore the indicator pattern; pressing a button now adds 1 to every listed counter. Starting from all zeros, reach the exact target vector with the fewest total presses.

## Parsing

I split the line into the indicator block `[...]`, the `(…)` button schematics, and the `{…}` joltage list. Each button becomes a slice of indices, the indicators become a boolean array, and the jolts become integers. Indices are validated to stay within the number of lights/counters (≤10 in the real input), which makes the linear algebra convenient.

## Part 1: GF(2) Linear System + Hamming Minimization

Let `x_j` be the number of times button `j` is pressed modulo 2. Each light is a linear equation over GF(2): the XOR of the relevant button variables must equal the target bit. I build an augmented matrix (`lights × buttons+1`), perform Gaussian elimination over GF(2), and obtain:

- A particular solution (derived from the reduced system).
- A basis for the nullspace (each free column yields one basis vector).

The answer is the minimum Hamming weight across all vectors in the affine space `x_particular + span(basis)`. Nullity is at most 3 in the dataset, so I can brute force all combinations of the basis vectors (`2^nullity <= 8`) and track the smallest weight. That gives the minimal number of presses needed.

**Complexity:** For each machine with `L` lights and `B` buttons, elimination costs `O(L * B^2)` because I manipulate bitsets, but with `B ≤ 13` this is negligible. Enumerating the nullspace is `O(2^nullity * B)`.

## Part 2: Integer Linear System with Branch & Bound

Here each button adds +1 to its listed counters, so the system is `A * x = target`, with non-negative integer variables and a cost function `minimize sum(x_j)`. Over the rationals I perform standard Gaussian elimination (exact arithmetic via `big.Rat`) to produce row-echelon form and identify pivot/free columns.

- If every column has a pivot, there’s a unique rational solution; I verify it’s integral, non-negative, and doesn’t demand pressing a button more times than the smallest counter it touches (an easy upper bound).
- If there are free variables (nullity ≤ 3 in the input), I enumerate them using DFS with pruning. Partial assignments are checked with interval reasoning: remaining free vars can’t push any pivot variable outside `[0, maxCounter]`. Once all free variables are set, I back-solve for the pivot variables and accumulate the total presses.

This strategy preserves exact arithmetic and avoids integer linear programming. The pruning quickly trims the search space because bounds are tight (targets ≤ 200 and buttons always touch at least one counter).

**Complexity:** Elimination is `O(N^3)` with `N ≤ 10`. The search explores at most `∏ (maxBound_i + 1)` states, but the constraints and pruning mean only a few hundred combinations are visited per machine.

## Testing

`Day10/main_test.go` feeds the sample three-machine input, asserting the totals `Part1 = 7` and `Part2 = 33`. Running `go test ./...` covers every day’s solver and ensures no regressions.

## Takeaways

Day 10 boiled down to solving tiny linear systems twice: once over GF(2) with a “minimum weight solution” objective, once over the integers with a bounded cost. Keeping the matrices small allowed exact arithmetic and exhaustive search without heavy tooling, and the shared parsing logic made both modes easy to reason about.
