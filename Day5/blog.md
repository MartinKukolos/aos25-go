# Day 5: Kitchen Inventory Journal

The cafeteria had overlapping ranges of “fresh” ingredient IDs and a second list of IDs currently on hand. Time to tidy up their interval math.

## Puzzle Recap

- Input has two sections: a list of inclusive ranges `a-b`, a blank line, then explicit ingredient IDs.
- Part 1: count how many listed IDs fall inside at least one fresh range.
- Part 2: ignore the explicit IDs and report how many integers are covered by the union of all ranges.

## Parsing Strategy

A single `bufio.Scanner` pass splits the file into “ranges” and “IDs” sections. Each range line becomes an `interval{start,end}` (swapping endpoints if needed), and each ID line becomes an `int64`. Blank lines flip the section counter.

## Interval Merging

Overlapping or adjacent ranges should act like a single continuous span. I sort the intervals by start/end and perform a standard merge: extend the current span whenever the next range begins before `current.end + 1`; otherwise, push the span and start a new one. The merged list feeds both parts.

## Freshness Checks

For part 1 I binary-search the merged intervals for each ID. Specifically, I find the first interval whose `end >= id` and then check whether `start <= id`. This keeps each lookup `O(log M)` where `M` is the number of merged intervals.

Part 2 simply sums `(end - start + 1)` across every merged interval—no need to enumerate IDs individually.

## Complexity Discussion

Let `R` be the number of raw ranges and `K` the number of explicit IDs.

- Parsing is `O(R + K)`.
- Merging sorts the ranges: `O(R log R)` time, `O(R)` space.
- Part 1 performs `K` binary searches → `O(K log R)`.
- Part 2 is `O(R)` for the merged sweep.

## Testing and Validation

The sample database (three fresh IDs and fourteen total fresh integers) anchors the unit tests. I also added quick checks for back-to-back ranges to ensure the `+1` merge rule keeps them continuous.

## Final Thoughts

Good interval hygiene goes a long way. Once the ranges were normalized, both questions became tiny: part 1 is “membership queries against sorted spans” and part 2 is “measure the union length.”
