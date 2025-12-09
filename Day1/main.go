package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	dialSize      = 100
	startPosition = 50
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
	if _, err := os.Stat("Day1/input.txt"); err == nil {
		return "Day1/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int, int, error) {
	scanner := bufio.NewScanner(r)
	position := startPosition
	zeroHitsEnd := 0
	zeroHitsAll := 0
	lineNumber := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineNumber++

		if len(line) < 2 {
			return 0, 0, fmt.Errorf("line %d: rotation too short", lineNumber)
		}

		dir := line[0]
		steps, err := strconv.Atoi(line[1:])
		if err != nil {
			return 0, 0, fmt.Errorf("line %d: invalid distance: %w", lineNumber, err)
		}

		zeroHitsAll += countZeroHits(position, steps, dir)

		switch dir {
		case 'L':
			position = mod(position - steps)
		case 'R':
			position = mod(position + steps)
		default:
			return 0, 0, fmt.Errorf("line %d: invalid direction %q", lineNumber, dir)
		}

		if position == 0 {
			zeroHitsEnd++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return zeroHitsEnd, zeroHitsAll, nil
}

func mod(value int) int {
	value %= dialSize
	if value < 0 {
		value += dialSize
	}
	return value
}

func countZeroHits(position, steps int, dir byte) int {
	if steps <= 0 {
		return 0
	}

	var first int
	switch dir {
	case 'L':
		first = position % dialSize
		if first == 0 {
			first = dialSize
		}
	case 'R':
		first = (dialSize - (position % dialSize)) % dialSize
		if first == 0 {
			first = dialSize
		}
	default:
		return 0
	}

	if steps < first {
		return 0
	}

	return 1 + (steps-first)/dialSize
}
