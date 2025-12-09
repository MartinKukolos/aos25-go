package main

import (
	"strings"
	"testing"
)

const sampleInput = "123 328  51 64 \n 45 64  387 23 \n  6 98  215 314\n*   +   *   +  \n"

func TestSolveSample(t *testing.T) {
	part1, part2, err := Solve(strings.NewReader(sampleInput))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}
	if part1 != 4277556 {
		t.Fatalf("part1 = %d, want 4277556", part1)
	}
	if part2 != 3263827 {
		t.Fatalf("part2 = %d, want 3263827", part2)
	}
}

