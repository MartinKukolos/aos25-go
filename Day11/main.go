package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	if _, err := os.Stat("Day11/input.txt"); err == nil {
		return "Day11/input.txt"
	}
	return "input.txt"
}

// Solve reads a directed graph specification and returns:
// - number of distinct simple paths from "you" to "out"
// - number of distinct simple paths from "svr" to "out" that visit both "dac" and "fft"
func Solve(r io.Reader) (int64, int64, error) {
	graph, err := parseGraph(r)
	if err != nil {
		return 0, 0, err
	}

	// Part 1: paths from "you" to "out"
	var p1 int64
	if _, ok := graph["you"]; ok {
		// Try optimized DAG DP; if cycle detected, fall back to DFS
		pruned, order, ok := prunedTopo(graph, "you", "out")
		if ok {
			p1 = countPathsDAG(pruned, order, "you", "out")
		} else {
			p1 = countPathsSimple(graph, "you", "out")
		}
	} else {
		p1 = 0
	}

	// Part 2: paths from "svr" to "out" that visit both dac and fft
	var p2 int64
	if _, ok := graph["svr"]; ok {
		// Try optimized DAG DP with 2-bit mask for dac/fft; fallback to DFS on cycles
		pruned, order, ok := prunedTopo(graph, "svr", "out")
		if ok {
			p2 = countPathsMustVisitDAG(pruned, order, "svr", "out", "dac", "fft")
		} else {
			p2 = countPathsWithMustVisit(graph, "svr", "out", "dac", "fft")
		}
	} else {
		p2 = 0
	}

	return p1, p2, nil
}

func parseGraph(r io.Reader) (map[string][]string, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 1024)
	scanner.Buffer(buf, 1<<20)
	graph := make(map[string][]string)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// allow comments in input
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Expect format: name: a b c
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			// If the line has no colon, skip it gracefully
			continue
		}
		src := strings.TrimSpace(parts[0])
		targets := strings.Fields(strings.TrimSpace(parts[1]))
		// ensure node exists even with no targets
		if _, exists := graph[src]; !exists {
			graph[src] = nil
		}
		if len(targets) > 0 {
			graph[src] = append(graph[src], targets...)
			// ensure target nodes appear in map too (even if no outgoing list given elsewhere)
			for _, t := range targets {
				if _, ok := graph[t]; !ok {
					graph[t] = nil
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return graph, nil
}

func countPathsSimple(graph map[string][]string, start, goal string) int64 {
	var total int64
	visited := make(map[string]bool)
	var dfs func(string)
	dfs = func(u string) {
		if u == goal {
			total++
			return
		}
		visited[u] = true
		for _, v := range graph[u] {
			if !visited[v] {
				dfs(v)
			}
		}
		visited[u] = false
	}
	dfs(start)
	return total
}

func countPathsWithMustVisit(graph map[string][]string, start, goal, mustA, mustB string) int64 {
	var total int64
	visited := make(map[string]bool)
	var dfs func(string, bool, bool)
	dfs = func(u string, seenA, seenB bool) {
		if u == mustA {
			seenA = true
		}
		if u == mustB {
			seenB = true
		}
		if u == goal {
			if seenA && seenB {
				total++
			}
			return
		}
		visited[u] = true
		for _, v := range graph[u] {
			if !visited[v] {
				dfs(v, seenA, seenB)
			}
		}
		visited[u] = false
	}
	dfs(start, false, false)
	return total
}

// prunedTopo builds a subgraph containing only nodes reachable from start and
// that can also reach goal. It then computes a topological order on that
// subgraph. Returns (subgraph, order, true) if DAG; if a cycle is detected in
// the pruned subgraph, returns (nil, nil, false).
func prunedTopo(graph map[string][]string, start, goal string) (map[string][]string, []string, bool) {
	// 1) Reachable from start (forward BFS)
	reach := make(map[string]bool)
	var stack []string
	stack = append(stack, start)
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if reach[u] {
			continue
		}
		reach[u] = true
		for _, v := range graph[u] {
			if !reach[v] {
				stack = append(stack, v)
			}
		}
	}

	// 2) Can reach goal (reverse BFS)
	rev := make(map[string][]string)
	for u, outs := range graph {
		for _, v := range outs {
			rev[v] = append(rev[v], u)
		}
		// ensure keys exist
		if _, ok := rev[u]; !ok {
			rev[u] = rev[u]
		}
	}
	canReachGoal := make(map[string]bool)
	stack = stack[:0]
	stack = append(stack, goal)
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if canReachGoal[u] {
			continue
		}
		canReachGoal[u] = true
		for _, p := range rev[u] {
			if !canReachGoal[p] {
				stack = append(stack, p)
			}
		}
	}

	// 3) Build pruned subgraph
	sub := make(map[string][]string)
	for u := range graph {
		if !(reach[u] && canReachGoal[u]) {
			continue
		}
		for _, v := range graph[u] {
			if reach[v] && canReachGoal[v] {
				sub[u] = append(sub[u], v)
			}
		}
		// ensure node exists even if no outgoing
		if _, ok := sub[u]; !ok {
			sub[u] = nil
		}
	}
	if _, ok := sub[start]; !ok {
		// Start not connected to goal
		return sub, nil, true // treat as DAG with zero ways; empty order okay
	}

	// 4) Topological sort (Kahn's algorithm)
	indeg := make(map[string]int)
	for u := range sub {
		indeg[u] = 0
	}
	for _, outs := range sub {
		for _, v := range outs {
			indeg[v]++
		}
	}
	// queue of zero in-degree nodes
	q := make([]string, 0, len(indeg))
	for u, d := range indeg {
		if d == 0 {
			q = append(q, u)
		}
	}
	order := make([]string, 0, len(indeg))
	for i := 0; i < len(q); i++ {
		u := q[i]
		order = append(order, u)
		for _, v := range sub[u] {
			indeg[v]--
			if indeg[v] == 0 {
				q = append(q, v)
			}
		}
	}
	if len(order) != len(indeg) {
		// cycle detected in pruned subgraph
		return nil, nil, false
	}
	return sub, order, true
}

// countPathsDAG counts number of paths from start to goal in a DAG using a
// topological order.
func countPathsDAG(graph map[string][]string, order []string, start, goal string) int64 {
	ways := make(map[string]int64, len(order))
	ways[start] = 1
	// process in topological order
	for _, u := range order {
		w := ways[u]
		if w == 0 {
			continue
		}
		for _, v := range graph[u] {
			ways[v] += w
		}
	}
	return ways[goal]
}

// countPathsMustVisitDAG counts number of paths from start to goal that visit
// both mustA and mustB (any order) in a DAG using DP over (node, mask).
func countPathsMustVisitDAG(graph map[string][]string, order []string, start, goal, mustA, mustB string) int64 {
	// nodeMask marks if being at node sets a bit
	bitA := 1
	bitB := 2
	nodeMask := func(name string) int {
		m := 0
		if name == mustA {
			m |= bitA
		}
		if name == mustB {
			m |= bitB
		}
		return m
	}
	// dp[node][mask]
	dp := make(map[string][4]int64, len(order))
	// initialize
	initMask := nodeMask(start)
	arr := dp[start]
	arr[initMask] = arr[initMask] + 1
	dp[start] = arr

	for _, u := range order {
		cur := dp[u]
		// propagate for all masks present
		for mask := 0; mask < 4; mask++ {
			val := cur[mask]
			if val == 0 {
				continue
			}
			for _, v := range graph[u] {
				nm := mask | nodeMask(v)
				nxt := dp[v]
				nxt[nm] += val
				dp[v] = nxt
			}
		}
	}
	return dp[goal][bitA|bitB]
}
