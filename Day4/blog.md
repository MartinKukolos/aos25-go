# Day 4: Forklift Field Notes

Optimizing forklift access meant deciding which paper rolls had fewer than four occupied neighbors, and then repeatedly peeling away accessible rolls to see how much paper we could free.

## Puzzle Recap

- Input: rectangular grid of `@` (roll) and `.` (empty) cells.
- Part 1: count rolls that currently have < 4 occupied neighbors among the 8 surrounding cells.
- Part 2: iteratively remove every accessible roll, then re-evaluate until no more can be removed; total how many rolls disappear.

## Grid Representation

I read the grid line-by-line, enforcing constant row width, and store it as `[][]bool`. That makes neighbor checks cheap and keeps memory predictable.

## Part 1: Static Accessibility

For each `@`, I count its eight neighbors using a small offset table. Any roll with a neighbor count below 4 contributes to the answer. This is a single pass over the grid.

## Part 2: Cascading Removals

To simulate forklifts removing rolls I:

1. Copy the grid so part 1’s result remains untouched.
2. Loop until no `@` satisfies the accessibility rule.
3. In each loop, collect all removable coordinates, flip them to empty, and add the batch size to the total.

This mimics a breadth-first wave of removals but doesn’t require a queue because each round re-scans the grid.

## Complexity Discussion

Let `R` and `C` be grid dimensions and `N = R*C`.

- Part 1 is `O(N)` with constant extra space.
- Part 2 worst-case scans the grid once per removal layer. In the extreme, each pass removes a single roll, giving `O(N^2)` time, though the actual puzzle grids shrink quickly in practice. Extra memory is `O(N)` for the working copy.

## Testing and Validation

The 10×10 example from the prompt serves as the regression test: part 1 must report 13 accessible rolls, and part 2 must total 43 removals. I also checked degenerate grids (all empty, all full) to be sure the loops terminate cleanly.

## Final Thoughts

This day’s lesson: sometimes the simplest simulation—scan, flag, remove—is the most maintainable. With clear neighbor counting and a copy of the grid, it’s easy to reason about both the static and dynamic views of the warehouse.
