# Day 7: Tachyon Lab Log

The teleporter lab handed me a beam-splitting puzzle: trace every tachyon path, count splitter events, and then embrace the many-worlds interpretation where every split duplicates timelines.

## Puzzle Recap

- Input is a grid with `S` at the top and `^` cells acting as splitters; everything else is empty space the beam can pass through.
- Beams start at `S` and travel downward.
- Part 1: whenever a beam hits `^`, it stops and spawns two new beams exiting left and right of the splitter. Count how many splitters actually fire.
- Part 2: quantum rules say a single particle follows every branch simultaneously. We need to know how many timelines remain when all branches finish (summing beam counts that exit the grid).

## Grid Handling

I read the grid into a slice of strings, locate `S`, and note the maximum width so I can track beams column by column even if rows are ragged.

## Part 1 Simulation

I maintain two boolean slices `active`/`next`, indicating which columns currently contain a beam just above the current row. As I step row by row:

- If the cell below a beam is empty, the beam keeps moving straight down.
- If it’s a splitter, I increment the split counter and activate the left and right neighbor columns for the next row.

Beams never move upward, so a simple double-buffer suffices.

## Part 2 Timeline Counting

The quantum version replaces booleans with `*big.Int` counts per column. When a beam hits a splitter, its count is added to both side branches (or to a running `completed` sum if the branch would exit the grid). After processing the entire grid, I add any remaining counts in the last `active` row to the total timeline count.

## Complexity Discussion

Let `R` be the number of rows and `C` the maximum width.

- Both simulations examine each cell at most once, so runtime is `O(R*C)`.
- Memory usage is `O(C)` for the active column slices, plus whatever BigInts require for the timeline counts in part 2.

## Testing and Validation

The example worksheet drives a unit test that asserts part 1’s split count (`21`) and the part 2 timeline total (`40`). I also tested grids without splitters and with splitters on the edges to confirm branches that fall off the grid are counted correctly.

## Final Thoughts

This day showcased how a modest state machine can scale from classical to quantum bookkeeping. Once the column-based DP was in place, swapping booleans for `*big.Int` gave me part 2 almost for free.
