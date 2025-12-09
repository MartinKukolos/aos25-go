package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	path := resolveInputPath(os.Args)

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open input %q: %v\n", path, err)
		os.Exit(1)
	}
	defer file.Close()

	part1, part2, err := Solve(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "solve error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}

func resolveInputPath(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	if _, err := os.Stat("Day8/input.txt"); err == nil {
		return "Day8/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int64, int64, error) {
	return SolveWithLimit(r, 1000)
}

func SolveWithLimit(r io.Reader, limit int) (int64, int64, error) {
	points, err := parsePoints(r)
	if err != nil {
		return 0, 0, err
	}
	if len(points) == 0 {
		return 0, 0, fmt.Errorf("no junction boxes found")
	}

	pairs := allPairs(points)
	if len(pairs) == 0 {
		return 0, 0, fmt.Errorf("need at least two junction boxes")
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].dist == pairs[j].dist {
			if pairs[i].i == pairs[j].i {
				return pairs[i].j < pairs[j].j
			}
			return pairs[i].i < pairs[j].i
		}
		return pairs[i].dist < pairs[j].dist
	})

	uf := newUnionFind(len(points))
	var part1 int64
	part1Computed := false
	var part2 int64
	part2Computed := false
	connectionsProcessed := 0

	for _, p := range pairs {
		if connectionsProcessed < limit {
			connectionsProcessed++
		}
		merged := uf.union(p.i, p.j)

		if !part1Computed && connectionsProcessed == limit {
			part1 = productOfTopThree(uf)
			part1Computed = true
			if part2Computed {
				break
			}
		}

		if merged && !part2Computed && uf.components == 1 {
			part2 = points[p.i].x * points[p.j].x
			part2Computed = true
			if part1Computed {
				break
			}
		}
	}

	if !part1Computed {
		part1 = productOfTopThree(uf)
		part1Computed = true
	}

	if !part2Computed {
		if uf.components != 1 {
			return 0, 0, fmt.Errorf("unable to connect all junction boxes")
		}
		return part1, 0, fmt.Errorf("missing final connection information")
	}

	return part1, part2, nil
}

type point struct {
	x int64
	y int64
	z int64
}

type pair struct {
	i    int
	j    int
	dist int64
}

func parsePoints(r io.Reader) ([]point, error) {
	scanner := bufio.NewScanner(r)
	var points []point
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		coords := strings.Split(line, ",")
		if len(coords) != 3 {
			return nil, fmt.Errorf("invalid coordinate line %q", line)
		}
		x, err := strconv.ParseInt(strings.TrimSpace(coords[0]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse X from %q: %w", line, err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(coords[1]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse Y from %q: %w", line, err)
		}
		z, err := strconv.ParseInt(strings.TrimSpace(coords[2]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse Z from %q: %w", line, err)
		}
		points = append(points, point{x: x, y: y, z: z})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return points, nil
}

func allPairs(points []point) []pair {
	var pairs []pair
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			pairs = append(pairs, pair{i: i, j: j, dist: squaredDistance(points[i], points[j])})
		}
	}
	return pairs
}

func squaredDistance(a, b point) int64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}

type unionFind struct {
	parent     []int
	size       []int
	components int
}

func newUnionFind(n int) *unionFind {
	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &unionFind{parent: parent, size: size, components: n}
}

func (uf *unionFind) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *unionFind) union(a, b int) bool {
	ra := uf.find(a)
	rb := uf.find(b)
	if ra == rb {
		return false
	}
	if uf.size[ra] < uf.size[rb] {
		ra, rb = rb, ra
	}
	uf.parent[rb] = ra
	uf.size[ra] += uf.size[rb]
	uf.components--
	return true
}

func productOfTopThree(uf *unionFind) int64 {
	var sizes []int
	for i := range uf.parent {
		if uf.parent[i] == i {
			sizes = append(sizes, uf.size[i])
		}
	}
	sort.Slice(sizes, func(i, j int) bool { return sizes[i] > sizes[j] })
	for len(sizes) < 3 {
		sizes = append(sizes, 1)
	}
	prod := int64(1)
	for i := 0; i < 3; i++ {
		prod *= int64(sizes[i])
	}
	return prod
}
