package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	if _, err := os.Stat("Day6/input.txt"); err == nil {
		return "Day6/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int64, int64, error) {
	grid, err := readGrid(r)
	if err != nil {
		return 0, 0, err
	}
	if len(grid) == 0 {
		return 0, 0, fmt.Errorf("empty grid")
	}

	part1, err := evaluateLeftToRight(grid)
	if err != nil {
		return 0, 0, err
	}
	part2, err := evaluateRightToLeft(grid)
	if err != nil {
		return 0, 0, err
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

type run struct {
	values []int64
	add    bool
}

func evaluateLeftToRight(grid []string) (int64, error) {
	runs, err := parseRunsLeftToRight(grid)
	if err != nil {
		return 0, err
	}
	return sumRuns(runs)
}

type span struct {
	start int
	end   int
}

func parseRunsLeftToRight(grid []string) ([]run, error) {
	spans := problemSpans(grid)
	runs := make([]run, 0, len(spans))
	rows := len(grid)
	for _, sp := range spans {
		values := []int64{}
		for r := 0; r < rows-1; r++ {
			segment := strings.TrimSpace(sliceRow(grid[r], sp.start, sp.end))
			if segment == "" {
				continue
			}
			value, err := strconv.ParseInt(segment, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse value %q at row %d columns [%d,%d): %w", segment, r, sp.start, sp.end, err)
			}
			values = append(values, value)
		}
		if len(values) == 0 {
			return nil, fmt.Errorf("no values found in columns [%d,%d)", sp.start, sp.end)
		}
		op, err := readOperator(grid, sp)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run{values: values, add: op == '+'})
	}
	return runs, nil
}

func problemSpans(grid []string) []span {
	cols := maxLen(grid)
	spans := []span{}
	for c := 0; c < cols; {
		if isBlankColumn(grid, c) {
			c++
			continue
		}
		start := c
		for c < cols && !isBlankColumn(grid, c) {
			c++
		}
		spans = append(spans, span{start: start, end: c})
	}
	return spans
}

func sliceRow(row string, start, end int) string {
	if start >= len(row) {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > len(row) {
		end = len(row)
	}
	if start >= end {
		return ""
	}
	return row[start:end]
}

func readOperator(grid []string, sp span) (byte, error) {
	if len(grid) == 0 {
		return 0, fmt.Errorf("empty grid")
	}
	line := grid[len(grid)-1]
	segment := strings.TrimSpace(sliceRow(line, sp.start, sp.end))
	if segment == "" {
		return 0, fmt.Errorf("missing operator in columns [%d,%d)", sp.start, sp.end)
	}
	op := segment[0]
	if op != '+' && op != '*' {
		return 0, fmt.Errorf("invalid operator %q in columns [%d,%d)", op, sp.start, sp.end)
	}
	return op, nil
}

func readColumnValue(grid []string, col int) (int64, bool, error) {
	rows := len(grid)
	if rows == 0 {
		return 0, false, fmt.Errorf("empty grid")
	}
	var sb strings.Builder
	for r := 0; r < rows-1; r++ {
		ch := charAt(grid[r], col)
		if ch == ' ' {
			continue
		}
		if ch < '0' || ch > '9' {
			return 0, false, fmt.Errorf("invalid digit %q at row %d column %d", ch, r, col)
		}
		sb.WriteByte(ch)
	}
	if sb.Len() == 0 {
		return 0, false, nil
	}
	text := sb.String()
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("parse column %d value %q: %w", col, text, err)
	}
	return value, true, nil
}

func sumRuns(runs []run) (int64, error) {
	var total int64
	for _, r := range runs {
		if len(r.values) == 0 {
			continue
		}
		var acc int64
		if r.add {
			for _, v := range r.values {
				acc += v
			}
		} else {
			acc = 1
			for _, v := range r.values {
				acc *= v
			}
		}
		total += acc
	}
	return total, nil
}

func evaluateRightToLeft(grid []string) (int64, error) {
	runs, err := parseRunsRightToLeft(grid)
	if err != nil {
		return 0, err
	}
	return sumRuns(runs)
}

func parseRunsRightToLeft(grid []string) ([]run, error) {
	spans := problemSpans(grid)
	runs := make([]run, 0, len(spans))
	for i := len(spans) - 1; i >= 0; i-- {
		sp := spans[i]
		values := []int64{}
		for c := sp.end - 1; c >= sp.start; c-- {
			value, ok, err := readColumnValue(grid, c)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			values = append(values, value)
		}
		if len(values) == 0 {
			return nil, fmt.Errorf("no column values found in columns [%d,%d)", sp.start, sp.end)
		}
		op, err := readOperator(grid, sp)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run{values: values, add: op == '+'})
	}
	return runs, nil
}

func maxLen(grid []string) int {
	max := 0
	for _, row := range grid {
		if len(row) > max {
			max = len(row)
		}
	}
	return max
}

func isBlankColumn(grid []string, c int) bool {
	for _, row := range grid {
		if c < len(row) && row[c] != ' ' {
			return false
		}
	}
	return true
}

func charAt(line string, idx int) byte {
	if idx < 0 || idx >= len(line) {
		return ' '
	}
	return line[idx]
}
