package main

import (
	"strings"
	"testing"
)

const sampleInput = "162,817,812\n57,618,57\n906,360,560\n592,479,940\n352,342,300\n466,668,158\n542,29,236\n431,825,988\n739,650,466\n52,470,668\n216,146,977\n819,987,18\n117,168,530\n805,96,715\n346,949,466\n970,615,88\n941,993,340\n862,61,35\n984,92,344\n425,690,689\n"

func TestSolveSample(t *testing.T) {
	part1, part2, err := SolveWithLimit(strings.NewReader(sampleInput), 10)
	if err != nil {
		t.Fatalf("SolveWithLimit error = %v", err)
	}
	if part1 != 40 {
		t.Fatalf("part1 = %d, want 40", part1)
	}
	if part2 != 25272 {
		t.Fatalf("part2 = %d, want 25272", part2)
	}
}
