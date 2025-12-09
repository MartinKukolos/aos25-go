# Day 6: Cephalopod Math Journal

Waiting inside the trash compactor turned into tutoring time. The worksheet stacked multiple arithmetic problems side-by-side, so solving it meant reading vertical slices carefully—twice.

## Puzzle Recap

- Input is a grid: digits stacked in columns with a `+` or `*` operator on the last row of each problem, separated by blank columns.
- Part 1: interpret problems left-to-right, treating each contiguous block of columns as one problem whose numbers are written top-down.
- Part 2: reinterpret the same sheet right-to-left; now every column holds exactly one number (digits stacked vertically), and problems are delineated by spaces when scanning from the right edge.

## Column Span Detection

I first locate “problem spans” by scanning columns and grouping contiguous non-blank columns. Each span knows its `[start,end)` indexes, making later slicing deterministic regardless of ragged row lengths.

## Part 1 Evaluation

For each span:

1. Slice every row above the operator row across `[start,end)` and trim whitespace to extract a number (if the slice is empty, skip that row).
2. Parse the operator from the last row within the same span.
3. Accumulate either a sum or product of the collected numbers and add it to the worksheet total.

This effectively transposes each vertical block into a traditional list of operands.

## Part 2 Evaluation

Right-to-left reading treats each column as an independent number. I iterate spans from the rightmost to leftmost and, inside a span, iterate columns from right to left. For every column I build the vertical number by concatenating digits from top to one row above the operator. Once I gather all column numbers for a span I apply the same sum/product logic as part 1.

## Complexity Discussion

Let `R` be the number of rows and `C` the maximum row length.

- Building spans takes `O(R*C)` to detect blank columns.
- Part 1 and part 2 each visit every cell at most a constant number of times while slicing strings, giving `O(R*C)` overall time and `O(C)` extra space for buffers.

## Testing and Validation

The provided worksheet sample yields `4,277,556` for part 1 and `3,263,827` for part 2; both values are asserted in `main_test.go`. I also validated edge cases such as empty rows and single-column problems to ensure span detection doesn’t panic.

## Final Thoughts

The trick was to focus on column spans first; once the grid is partitioned cleanly, both interpretations fall out naturally. Explicit slicing kept the parser robust even when numbers were misaligned within their columns.
