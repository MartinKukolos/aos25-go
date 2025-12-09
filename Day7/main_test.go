package main

import (
	"math/big"
	"strings"
	"testing"
)

const sampleInput = ".......S.......\n...............\n.......^.......\n...............\n......^.^......\n...............\n.....^.^.^.....\n...............\n....^.^...^....\n...............\n...^.^...^.^...\n...............\n..^...^.....^..\n...............\n.^.^.^.^.^...^.\n...............\n"

func TestSolveSample(t *testing.T) {
	part1, part2, err := Solve(strings.NewReader(sampleInput))
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}
	if part1 != 21 {
		t.Fatalf("part1 = %d, want 21", part1)
	}
	want := big.NewInt(40)
	if part2 == nil || part2.Cmp(want) != 0 {
		t.Fatalf("part2 = %v, want %v", part2, want)
	}
}
