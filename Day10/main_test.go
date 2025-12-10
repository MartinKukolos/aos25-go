package main

import (
	"os"
	"testing"
)

func TestSolveSample(t *testing.T) {
	f, err := os.Open("sample.txt")
	if err != nil {
		t.Fatalf("open sample: %v", err)
	}
	defer f.Close()

	part1, part2, err := Solve(f)
	if err != nil {
		t.Fatalf("Solve(sample) error = %v", err)
	}
	if part1 != 7 {
		t.Fatalf("part1 = %d, want 7", part1)
	}
	if part2 != 33 {
		t.Fatalf("part2 = %d, want 33", part2)
	}
}
