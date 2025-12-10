package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

type machine struct {
	lights  []bool
	buttons [][]int
	jolts   []int
}

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
	if _, err := os.Stat("Day10/input.txt"); err == nil {
		return "Day10/input.txt"
	}
	return "input.txt"
}

func Solve(r io.Reader) (int64, int64, error) {
	machines, err := parseMachines(r)
	if err != nil {
		return 0, 0, err
	}
	var sumPart1 int64
	var sumPart2 int64
	for idx, m := range machines {
		p1, err := minIndicatorPresses(m)
		if err != nil {
			return 0, 0, fmt.Errorf("machine %d indicators: %w", idx+1, err)
		}
		sumPart1 += int64(p1)
		p2, err := minJoltagePresses(m)
		if err != nil {
			return 0, 0, fmt.Errorf("machine %d jolts: %w", idx+1, err)
		}
		sumPart2 += p2
	}
	return sumPart1, sumPart2, nil
}

func parseMachines(r io.Reader) ([]machine, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 1024)
	scanner.Buffer(buf, 1<<20)
	var machines []machine
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		m, err := parseMachineLine(line)
		if err != nil {
			return nil, err
		}
		machines = append(machines, m)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return machines, nil
}

func parseMachineLine(line string) (machine, error) {
	var result machine
	open := strings.Index(line, "[")
	closeIdx := strings.Index(line, "]")
	if open == -1 || closeIdx == -1 || closeIdx <= open {
		return result, fmt.Errorf("invalid indicator diagram: %q", line)
	}
	pattern := line[open+1 : closeIdx]
	lights := make([]bool, len(pattern))
	for i, ch := range pattern {
		switch ch {
		case '.':
			lights[i] = false
		case '#':
			lights[i] = true
		default:
			return result, fmt.Errorf("invalid indicator char %q", string(ch))
		}
	}
	rest := strings.TrimSpace(line[closeIdx+1:])
	var buttons [][]int
	for len(rest) > 0 && rest[0] == '(' {
		end := strings.Index(rest, ")")
		if end == -1 {
			return result, errors.New("missing closing parenthesis")
		}
		inside := rest[1:end]
		idxs, err := parseIndexList(inside)
		if err != nil {
			return result, err
		}
		buttons = append(buttons, idxs)
		rest = strings.TrimSpace(rest[end+1:])
	}
	if len(rest) == 0 || rest[0] != '{' {
		return result, errors.New("missing joltage requirements")
	}
	end := strings.Index(rest, "}")
	if end == -1 {
		return result, errors.New("missing closing brace")
	}
	targetList := rest[1:end]
	jolts, err := parseIntList(targetList)
	if err != nil {
		return result, err
	}
	rest = strings.TrimSpace(rest[end+1:])
	if rest != "" {
		return result, fmt.Errorf("unexpected trailing data %q", rest)
	}
	if len(jolts) == 0 && len(lights) != 0 {
		return result, errors.New("joltage requirements missing entries")
	}
	if len(lights) != len(jolts) && len(jolts) != 0 {
		if len(lights) == 0 {
			lights = make([]bool, len(jolts))
		} else {
			return result, fmt.Errorf("indicator lights (%d) and joltage counters (%d) mismatch", len(lights), len(jolts))
		}
	}
	// Validate indices
	for _, btn := range buttons {
		for _, idx := range btn {
			if idx < 0 || idx >= len(lights) {
				return result, fmt.Errorf("button references invalid index %d for %d lights", idx, len(lights))
			}
			if idx >= len(jolts) {
				return result, fmt.Errorf("button index %d beyond joltage counters %d", idx, len(jolts))
			}
		}
	}
	result = machine{lights: lights, buttons: buttons, jolts: jolts}
	return result, nil
}

func parseIndexList(s string) ([]int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	values := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

func parseIntList(s string) ([]int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	values := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

func minIndicatorPresses(m machine) (int, error) {
	nLights := len(m.lights)
	nButtons := len(m.buttons)
	if nLights == 0 {
		return 0, nil
	}
	words := ((nButtons + 1) + 63) / 64
	rows := make([][]uint64, nLights)
	for i := 0; i < nLights; i++ {
		row := make([]uint64, words)
		for j, btn := range m.buttons {
			for _, idx := range btn {
				if idx == i {
					setBit(row, j)
					break
				}
			}
		}
		if m.lights[i] {
			setBit(row, nButtons)
		}
		rows[i] = row
	}
	pivotColForRow, pivotRowForCol, err := rrefGF2(rows, nButtons)
	if err != nil {
		return 0, err
	}
	particular := makeBitset(nButtons)
	for r, col := range pivotColForRow {
		if col == -1 {
			continue
		}
		if getBit(rows[r], nButtons) {
			flipBit(particular, col)
		}
	}
	var basis [][]uint64
	for col := 0; col < nButtons; col++ {
		if pivotRowForCol[col] != -1 {
			continue
		}
		vec := makeBitset(nButtons)
		setBit(vec, col)
		for r, pivotCol := range pivotColForRow {
			if pivotCol == -1 {
				continue
			}
			if getBit(rows[r], col) {
				flipBit(vec, pivotCol)
			}
		}
		basis = append(basis, vec)
	}
	best := minWeightInAffine(particular, basis)
	return best, nil
}

func rrefGF2(rows [][]uint64, nCols int) ([]int, []int, error) {
	nRows := len(rows)
	pivotColForRow := make([]int, nRows)
	for i := range pivotColForRow {
		pivotColForRow[i] = -1
	}
	pivotRowForCol := make([]int, nCols)
	for i := range pivotRowForCol {
		pivotRowForCol[i] = -1
	}
	row := 0
	for col := 0; col < nCols && row < nRows; col++ {
		pivot := -1
		for r := row; r < nRows; r++ {
			if getBit(rows[r], col) {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			continue
		}
		rows[row], rows[pivot] = rows[pivot], rows[row]
		pivotColForRow[row] = col
		pivotRowForCol[col] = row
		for r := 0; r < nRows; r++ {
			if r == row {
				continue
			}
			if getBit(rows[r], col) {
				xorRow(rows[r], rows[row])
			}
		}
		row++
	}
	for r := 0; r < nRows; r++ {
		if pivotColForRow[r] == -1 && getBit(rows[r], nCols) {
			return nil, nil, errors.New("no solution for indicators")
		}
	}
	return pivotColForRow, pivotRowForCol, nil
}

func minWeightInAffine(part []uint64, basis [][]uint64) int {
	best := bitCount(part)
	if len(basis) == 0 {
		return best
	}
	total := 1 << len(basis)
	for mask := 1; mask < total; mask++ {
		vec := cloneBits(part)
		m := mask
		idx := 0
		for m != 0 {
			if m&1 == 1 {
				xorRow(vec, basis[idx])
			}
			m >>= 1
			idx++
		}
		w := bitCount(vec)
		if w < best {
			best = w
			if best == 0 {
				return 0
			}
		}
	}
	return best
}

func minJoltagePresses(m machine) (int64, error) {
	rows := len(m.jolts)
	cols := len(m.buttons)
	if rows == 0 {
		return 0, nil
	}
	matrix := make([][]*big.Rat, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]*big.Rat, cols+1)
		for j := 0; j < cols; j++ {
			coeff := int64(0)
			for _, idx := range m.buttons[j] {
				if idx == i {
					coeff = 1
					break
				}
			}
			matrix[i][j] = big.NewRat(coeff, 1)
		}
		matrix[i][cols] = big.NewRat(int64(m.jolts[i]), 1)
	}
	pivotColForRow := make([]int, rows)
	for i := range pivotColForRow {
		pivotColForRow[i] = -1
	}
	row := 0
	for col := 0; col < cols && row < rows; col++ {
		pivot := -1
		for r := row; r < rows; r++ {
			if matrix[r][col].Sign() != 0 {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			continue
		}
		matrix[row], matrix[pivot] = matrix[pivot], matrix[row]
		pivotVal := new(big.Rat).Set(matrix[row][col])
		for c := col; c <= cols; c++ {
			matrix[row][c] = new(big.Rat).Quo(matrix[row][c], pivotVal)
		}
		for r := 0; r < rows; r++ {
			if r == row {
				continue
			}
			factor := new(big.Rat).Set(matrix[r][col])
			if factor.Sign() == 0 {
				continue
			}
			for c := col; c <= cols; c++ {
				term := new(big.Rat).Mul(factor, matrix[row][c])
				matrix[r][c].Sub(matrix[r][c], term)
			}
		}
		pivotColForRow[row] = col
		row++
	}
	rhsIdx := cols
	for r := 0; r < rows; r++ {
		zero := true
		for c := 0; c < cols; c++ {
			if matrix[r][c].Sign() != 0 {
				zero = false
				break
			}
		}
		if zero && matrix[r][rhsIdx].Sign() != 0 {
			return 0, errors.New("no integer solution for joltage")
		}
	}
	pivotRowForCol := make([]int, cols)
	for i := range pivotRowForCol {
		pivotRowForCol[i] = -1
	}
	for r, col := range pivotColForRow {
		if col >= 0 {
			pivotRowForCol[col] = r
		}
	}
	var freeCols []int
	for col := 0; col < cols; col++ {
		if pivotRowForCol[col] == -1 {
			freeCols = append(freeCols, col)
		}
	}
	varMax := buttonUpperBounds(m)
	if len(freeCols) == 0 {
		total, err := evaluateJoltageSolution(matrix, pivotColForRow, varMax)
		return total, err
	}
	best := int64(math.MaxInt64)
	assigned := make([]int64, len(freeCols))
	varMaxRat := make([]*big.Rat, len(varMax))
	for i, v := range varMax {
		varMaxRat[i] = big.NewRat(v, 1)
	}
	var dfs func(int, int64)
	dfs = func(idx int, sum int64) {
		if sum >= best {
			return
		}
		if idx == len(freeCols) {
			total := sum
			values := make([]*big.Rat, cols)
			for i, col := range freeCols {
				values[col] = big.NewRat(assigned[i], 1)
			}
			for r, col := range pivotColForRow {
				if col == -1 {
					continue
				}
				val := new(big.Rat).Set(matrix[r][rhsIdx])
				for _, freeCol := range freeCols {
					coeff := matrix[r][freeCol]
					if coeff.Sign() == 0 {
						continue
					}
					term := new(big.Rat).Mul(coeff, values[freeCol])
					val.Sub(val, term)
				}
				if val.Sign() < 0 || !val.IsInt() {
					return
				}
				intVal := val.Num().Int64()
				if intVal > varMax[col] {
					return
				}
				total += intVal
				if total >= best {
					return
				}
				values[col] = val
			}
			if total < best {
				best = total
			}
			return
		}
		col := freeCols[idx]
		maxVal := varMax[col]
		for v := int64(0); v <= maxVal; v++ {
			assigned[idx] = v
			newSum := sum + v
			if newSum >= best {
				continue
			}
			if !partialFeasible(assigned[:idx+1], freeCols[:idx+1], freeCols, matrix, pivotColForRow, varMaxRat) {
				continue
			}
			dfs(idx+1, newSum)
		}
	}
	dfs(0, 0)
	if best == int64(math.MaxInt64) {
		return 0, errors.New("no feasible joltage configuration")
	}
	return best, nil
}

func evaluateJoltageSolution(matrix [][]*big.Rat, pivotColForRow []int, varMax []int64) (int64, error) {
	rhsIdx := len(matrix[0]) - 1
	var total int64
	for r, col := range pivotColForRow {
		if col == -1 {
			continue
		}
		val := matrix[r][rhsIdx]
		if val.Sign() < 0 || !val.IsInt() {
			return 0, errors.New("non-integer solution")
		}
		intVal := val.Num().Int64()
		if intVal > varMax[col] {
			return 0, errors.New("button press exceeds counter bound")
		}
		total += intVal
	}
	return total, nil
}

func partialFeasible(assignedVals []int64, assignedCols []int, allFreeCols []int, matrix [][]*big.Rat, pivotColForRow []int, varMaxRat []*big.Rat) bool {
	rhsIdx := len(matrix[0]) - 1
	assignCount := len(assignedVals)
	for r, pivotCol := range pivotColForRow {
		if pivotCol == -1 {
			if matrix[r][rhsIdx].Sign() != 0 {
				return false
			}
			continue
		}
		val := new(big.Rat).Set(matrix[r][rhsIdx])
		for i := 0; i < assignCount; i++ {
			coeff := matrix[r][assignedCols[i]]
			if coeff.Sign() == 0 {
				continue
			}
			term := new(big.Rat).Mul(coeff, big.NewRat(assignedVals[i], 1))
			val.Sub(val, term)
		}
		minVal := new(big.Rat).Set(val)
		maxVal := new(big.Rat).Set(val)
		for i := assignCount; i < len(allFreeCols); i++ {
			coeff := matrix[r][allFreeCols[i]]
			if coeff.Sign() == 0 {
				continue
			}
			limit := varMaxRat[allFreeCols[i]]
			term := new(big.Rat).Mul(coeff, limit)
			if coeff.Sign() >= 0 {
				minVal.Sub(minVal, term)
			} else {
				maxVal.Sub(maxVal, term)
			}
		}
		maxAllowed := varMaxRat[pivotCol]
		if minVal.Cmp(maxAllowed) > 0 {
			return false
		}
		if maxVal.Sign() < 0 {
			return false
		}
		if minVal.Cmp(maxVal) > 0 {
			return false
		}
	}
	return true
}

func buttonUpperBounds(m machine) []int64 {
	bounds := make([]int64, len(m.buttons))
	for i, btn := range m.buttons {
		if len(btn) == 0 {
			bounds[i] = 0
			continue
		}
		minVal := math.MaxInt64
		for _, idx := range btn {
			if m.jolts[idx] < minVal {
				minVal = m.jolts[idx]
			}
		}
		bounds[i] = int64(minVal)
	}
	return bounds
}

func makeBitset(bits int) []uint64 {
	return make([]uint64, (bits+63)/64)
}

func setBit(row []uint64, idx int) {
	row[idx/64] |= 1 << (idx & 63)
}

func flipBit(row []uint64, idx int) {
	row[idx/64] ^= 1 << (idx & 63)
}

func getBit(row []uint64, idx int) bool {
	return (row[idx/64]>>(idx&63))&1 == 1
}

func xorRow(dst, src []uint64) {
	for i := range dst {
		dst[i] ^= src[i]
	}
}

func cloneBits(src []uint64) []uint64 {
	dup := make([]uint64, len(src))
	copy(dup, src)
	return dup
}

func bitCount(row []uint64) int {
	total := 0
	for _, w := range row {
		total += bits.OnesCount64(w)
	}
	return total
}
