package main

import (
	"strings"
	"testing"
)

const sampleInput = "987654321111111\n811111111111119\n234234234234278\n818181911112111\n"

func TestSolveSample(t *testing.T) {
	part1, part2, err := Solve(strings.NewReader(sampleInput))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	var wantPart1 int64 = 357
	var wantPart2 int64 = 3121910778619

	if part1 != wantPart1 {
		t.Fatalf("Solve() part1 = %d, want %d", part1, wantPart1)
	}
	if part2 != wantPart2 {
		t.Fatalf("Solve() part2 = %d, want %d", part2, wantPart2)
	}
}
