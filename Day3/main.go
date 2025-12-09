package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const part1Digits = 2
const part2Digits = 12

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
	if _, err := os.Stat("Day3/input.txt"); err == nil {
		return "Day3/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int64, int64, error) {
	scanner := bufio.NewScanner(r)
	var totalPart1 int64
	var totalPart2 int64

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		maxP1, err := maxValueForDigits(line, part1Digits)
		if err != nil {
			return 0, 0, err
		}
		maxP2, err := maxValueForDigits(line, part2Digits)
		if err != nil {
			return 0, 0, err
		}
		totalPart1 += maxP1
		totalPart2 += maxP2
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return totalPart1, totalPart2, nil
}

func maxValueForDigits(line string, pick int) (int64, error) {
	if pick <= 0 {
		return 0, fmt.Errorf("invalid pick %d", pick)
	}
	digits := []byte(line)
	if len(digits) < pick {
		return 0, fmt.Errorf("line %q shorter than %d digits", line, pick)
	}

	stack := make([]byte, 0, pick)
	for i, d := range digits {
		remaining := len(digits) - i
		for len(stack) > 0 && len(stack)+remaining > pick && stack[len(stack)-1] < d {
			stack = stack[:len(stack)-1]
		}
		if len(stack) < pick {
			stack = append(stack, d)
		}
	}

	var val int64
	for _, d := range stack {
		val = val*10 + int64(d-'0')
	}
	return val, nil
}
