# Day 1: Safe-Cracking Notes

My holiday journey started with a supposedly simple task: spin a safe dial per a list of rotations and figure out how often the pointer hits `0`. Two different password policies meant I needed to count zeros both at the end of each rotation and during the intermediate clicks.

## Puzzle Recap

- The dial shows `0-99`, starts at `50`, and each instruction is `L` or `R` plus a distance.
- Part 1 counts how many rotations finish with the dial pointing at `0`.
- Part 2 counts *every* click that lands on `0`, even mid-rotation.

## Parsing and Simulation

I streamed the instructions with a `bufio.Scanner`, trimming blank lines and validating each row for a direction and integer distance. The dial position always stays in `[0,99]` thanks to a `mod` helper that wraps negatives correctly.

For part 1 the simulation is literal: update the position after each rotation and increment a counter whenever the landing spot is `0`.

Part 2 needs arithmetic rather than per-click simulation. Given a direction and `steps`, I compute how many full 100-click laps the dial makes and whether the path crosses `0` before the final position. The `countZeroHits` helper handles that math in `O(1)` by finding the first time the rotation would wrap past `0` and then adding one hit every additional 100 clicks.

## Complexity Discussion

Let `m` be the number of rotations.

- Input parsing and part 1 simulation are both `O(m)`.
- Part 2 uses constant-time math per instruction, so it is also `O(m)`.
- Memory stays `O(1)` beyond the input buffer.

## Testing and Validation

Unit tests feed the sample rotation list through `Solve` and assert both password counts (`3` for part 1, `6` for part 2). I also spot-check extreme single rotations like `R1000` to ensure the per-lap counting works.

## Final Thoughts

This day is all about being careful with modular arithmetic. By doing the per-rotation math analytically I avoided simulating thousands of clicks and kept the code short and predictable.
