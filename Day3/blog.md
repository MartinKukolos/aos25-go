# Day 3: Escalator Battery Log

The printing department needed maximum joltage from each battery bank. That translated into picking the lexicographically largest subsequence of digits under two different pick counts.

## Puzzle Recap

- Input: each line is a bank—just a string of digits.
- Part 1: select exactly 2 digits per bank (in order) to form the largest possible two-digit number; sum across banks.
- Part 2: same idea, but select 12 digits per bank.

## Greedy Digit Selection

This is the classic “keep the biggest subsequence” problem. For each line I run a monotonic stack that keeps the best digits seen so far while ensuring I can still fill the remaining slots:

1. Iterate the digits with index `i`.
2. While the top of the stack is smaller than the current digit and there are enough remaining digits to fill `pick`, pop the stack.
3. Push the current digit if the stack isn’t full yet.
4. Convert the stack to a number at the end.

Because the digits can repeat heavily, this greedy approach ensures we always keep the best possible prefix while respecting order.

## Complexity Discussion

Let `n` be the number of digits in a bank and `k` the pick size (2 or 12).

- Each digit is pushed and popped at most once → `O(n)` per bank.
- Summed over all banks, the runtime is linear in the total input size.
- Memory per bank is `O(k)` for the stack.

## Testing and Validation

I copied the four-bank example into a unit test. It asserts both the part 1 sum (`357`) and the much larger part 2 total (`3121910778619`). Additional tests check edge cases where the line length equals the pick count to ensure the algorithm keeps everything.

## Final Thoughts

This day showed how far a small, well-known greedy trick can go. With just a stack and awareness of the characters left to process, the solution handles both tiny and huge pick sizes gracefully.
