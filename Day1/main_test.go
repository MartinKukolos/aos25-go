package main

import (
	"strings"
	"testing"
)

func TestSolveSample(t *testing.T) {
	const input = `L68
L30
R48
L5
R60
L55
L1
L99
R14
L82
`

	part1, part2, err := Solve(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	const wantPart1 = 3
	const wantPart2 = 6
	if part1 != wantPart1 {
		t.Fatalf("Solve() part1 = %d, want %d", part1, wantPart1)
	}
	if part2 != wantPart2 {
		t.Fatalf("Solve() part2 = %d, want %d", part2, wantPart2)
	}
}

func TestSolveMultiRevolution(t *testing.T) {
	const input = "R1000\n"

	part1, part2, err := Solve(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	const wantPart1 = 0
	const wantPart2 = 10
	if part1 != wantPart1 {
		t.Fatalf("Solve() part1 = %d, want %d", part1, wantPart1)
	}
	if part2 != wantPart2 {
		t.Fatalf("Solve() part2 = %d, want %d", part2, wantPart2)
	}
}
