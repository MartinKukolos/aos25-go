# Day 2: Gift Shop Forensics

Today’s stop was the gift shop, where I needed to flag “cute pattern” IDs hiding in huge numeric ranges. It felt more like number theory than inventory control.

## Puzzle Recap

- Input is a single comma-separated list of inclusive ranges `a-b`.
- Part 1: sum every ID that consists of some block of digits repeated **exactly twice** (examples: `11`, `6464`).
- Part 2: broaden to any repetition count ≥ 2 (`123123123`, `565656`, etc.).

## Parsing Ranges

I read the full line, split on commas, and turn each chunk into an `idRange{start,end}`. Validation catches missing dashes and reversed endpoints. Having everything in memory makes the later math loops straightforward.

## Counting Exact Doubles (Part 1)

Instead of iterating every value, I iterate block lengths `k` (1–9 because IDs are ≤ 18 digits). Any repeated-twice number looks like `base * (10^k + 1)` where `base` is a `k`-digit integer. For each range I intersect the allowable `base` interval with `[ceil(start/multiplier), floor(end/multiplier)]` and sum the resulting arithmetic progression in `O(1)`.

## Counting Any Repetition (Part 2)

For part 2 I still limit by total digit length, but now each length can have multiple divisors. For a target length `L` I:

1. Enumerate proper divisors `d` of `L` (meaning a `d`-digit block repeats `L/d` times).
2. Compute the multiplier `repeatMultiplier(d, repeats)` using geometric-series math.
3. Sum every `d`-digit base whose repeated form stays inside the `[start,end]` segment.
4. Apply inclusion–exclusion so that patterns built from smaller blocks aren’t double-counted.

Each length contributes at most the number of its divisors, so the work stays tiny.

## Complexity Discussion

Let `R` be the number of input ranges.

- Parsing is `O(R)`.
- Part 1 iterates at most 9 block sizes per range → `O(R)`.
- Part 2 iterates up to 18 total digits, and each length touches only its divisors (≤ the divisor count of 18). So still effectively `O(R)` with a very small constant.
- Memory usage is `O(1)` besides the range slice.

## Testing and Validation

The provided example ranges drive my unit test: part 1 must total `1227775554`, part 2 must reach `4174379265`. I also throw in a few synthetic ranges covering single blocks to ensure the arithmetic boundaries are correct.

## Final Thoughts

This puzzle rewarded algebraic thinking. By modeling repeating numbers parametrically, I avoided scanning massive ranges and instead solved everything with interval arithmetic and a sprinkle of inclusion–exclusion.
