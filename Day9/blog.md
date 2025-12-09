# Day 9: Theater Tiling Report

At the movie theater I had to evaluate every pair of red tiles to find the biggest rectangle, then enforce a “red or green only” constraint that turns the outline into a filled polygon.

## Puzzle Recap

- Input: coordinates of red tiles; successive points are connected by straight lines of green tiles, and the list wraps around.
- Part 1: choose any two red tiles as opposite corners; maximize rectangle area regardless of interior tiles.
- Part 2: the rectangle’s interior must consist solely of red or green tiles, meaning it must lie entirely inside the polygon traced by the red/green loop.

## Part 1: Brute-Force Pairs

With only the red points to consider, the simplest strategy is best: iterate all `n(n-1)/2` pairs, compute the axis-aligned rectangle defined by those corners, and track the maximum area. Because coordinates are modest and `n` isn’t huge, this direct approach is plenty fast.

## Modeling the Polygon Interior

For part 2 I needed a grid of which unit squares are “inside” the red/green region. Steps:

1. Collect the unique sorted `x` and `y` coordinates from the red tiles; these define grid cell boundaries.
2. For each horizontal strip between consecutive `y` values, perform a scanline: cast a horizontal ray across the polygon edges (which are guaranteed to be axis-aligned) and toggle inside/outside whenever you cross a vertical edge. Fill every cell between paired intersections.
3. Build a 2D prefix-sum array over this boolean grid so interior queries become `O(1)`.

## Part 2: Valid Rectangles Only

When evaluating a red pair for part 2, I map each coordinate to its compressed grid index and compute the area as before. Before accepting the rectangle, I query the prefix sum to ensure every cell beneath it is marked inside; if not, the rectangle would include forbidden tiles and must be skipped.

## Complexity Discussion

Let `n` be the number of red tiles, `X` the number of unique x-values, and `Y` the number of unique y-values.

- Part 1: `O(n^2)` pair checks.
- Building the interior grid: scanline intersect passes touch every polygon edge per strip, roughly `O(n * Y)`; prefix sums are `O(XY)`.
- Part 2 rectangle checks: still `O(n^2)`, but each includes an `O(1)` prefix query.
- Memory: `O(XY)` for the inside grid and prefix sums.

## Testing and Validation

The prompt’s sample (max area 50 unrestricted, 24 with the green constraint) anchors my tests. I also feed in thin rectangles and degenerate polygons to make sure the scanline logic handles shared edges correctly.

## Final Thoughts

Axis-aligned geometry problems reward coordinate compression and prefix sums. Once the polygon interior was rasterized, validating “all tiles are allowed” became a fast lookup, letting the brute-force pair search continue to shine.
