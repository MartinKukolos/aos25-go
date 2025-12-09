package main

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
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
	fmt.Printf("Part 2: %s\n", part2.String())
}

func resolveInputPath(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	if _, err := os.Stat("Day7/input.txt"); err == nil {
		return "Day7/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int64, *big.Int, error) {
	grid, err := readGrid(r)
	if err != nil {
		return 0, nil, err
	}
	if len(grid) == 0 {
		return 0, nil, fmt.Errorf("empty grid")
	}
	startRow, startCol, err := findStart(grid)
	if err != nil {
		return 0, nil, err
	}

	part1, err := simulatePart1(grid, startRow, startCol)
	if err != nil {
		return 0, nil, err
	}
	part2, err := simulatePart2(grid, startRow, startCol)
	if err != nil {
		return 0, nil, err
	}
	return part1, part2, nil
}

func readGrid(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func findStart(grid []string) (int, int, error) {
	for row, line := range grid {
		for col := 0; col < len(line); col++ {
			if line[col] == 'S' {
				return row, col, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("no start position found")
}

func simulatePart1(grid []string, startRow, startCol int) (int64, error) {
	cols := maxLen(grid)
	if cols == 0 {
		return 0, fmt.Errorf("invalid grid width")
	}
	active := make([]bool, cols)
	next := make([]bool, cols)
	if startCol >= cols {
		return 0, fmt.Errorf("start column outside grid")
	}
	active[startCol] = true

	var splits int64
	for row := startRow + 1; row < len(grid); row++ {
		line := grid[row]
		for i := range next {
			next[i] = false
		}
		for col := 0; col < cols; col++ {
			if !active[col] {
				continue
			}
			cell := charAt(line, col)
			if cell == '^' {
				splits++
				if col-1 >= 0 {
					next[col-1] = true
				}
				if col+1 < cols {
					next[col+1] = true
				}
			} else {
				next[col] = true
			}
		}
		active, next = next, active
	}

	return splits, nil
}

func simulatePart2(grid []string, startRow, startCol int) (*big.Int, error) {
	cols := maxLen(grid)
	if cols == 0 {
		return nil, fmt.Errorf("invalid grid width")
	}
	active := make([]*big.Int, cols)
	next := make([]*big.Int, cols)
	if startCol >= cols {
		return nil, fmt.Errorf("start column outside grid")
	}
	active[startCol] = big.NewInt(1)

	completed := big.NewInt(0)
	for row := startRow + 1; row < len(grid); row++ {
		line := grid[row]
		for i := range next {
			next[i] = nil
		}
		for col := 0; col < cols; col++ {
			count := active[col]
			if count == nil {
				continue
			}
			cell := charAt(line, col)
			if cell == '^' {
				if col-1 >= 0 {
					addCount(next, col-1, count)
				} else {
					completed.Add(completed, count)
				}
				if col+1 < cols {
					addCount(next, col+1, count)
				} else {
					completed.Add(completed, count)
				}
			} else {
				addCount(next, col, count)
			}
		}
		active, next = next, active
	}

	total := new(big.Int).Set(completed)
	for _, count := range active {
		if count != nil {
			total.Add(total, count)
		}
	}
	return total, nil
}

func addCount(dst []*big.Int, idx int, value *big.Int) {
	if idx < 0 || idx >= len(dst) || value == nil {
		return
	}
	if dst[idx] == nil {
		dst[idx] = new(big.Int).Set(value)
	} else {
		dst[idx].Add(dst[idx], value)
	}
}

func maxLen(lines []string) int {
	max := 0
	for _, line := range lines {
		if len(line) > max {
			max = len(line)
		}
	}
	return max
}

func charAt(line string, idx int) byte {
	if idx < 0 || idx >= len(line) {
		return '.'
	}
	return line[idx]
}
