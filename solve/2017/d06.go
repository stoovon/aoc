package solve2017

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day6 struct {
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 6}
}

func (d Day6) Part1(input string) (string, error) {
	banks := parseBanks(input)
	seen := make(map[string]bool)
	steps := 0

	for {
		key := banksKey(banks)
		if seen[key] {
			return strconv.Itoa(steps), nil
		}
		seen[key] = true
		redistribute(banks)
		steps++
	}
}

func parseBanks(input string) []int {
	fields := strings.Fields(input)
	banks := make([]int, len(fields))
	for i, f := range fields {
		banks[i], _ = strconv.Atoi(f)
	}
	return banks
}

func banksKey(banks []int) string {
	parts := make([]string, len(banks))
	for i, v := range banks {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}

func redistribute(banks []int) {
	n := len(banks)
	// Find max
	maxIdx, maxVal := 0, banks[0]
	for i, v := range banks {
		if v > maxVal || (v == maxVal && i < maxIdx) {
			maxIdx, maxVal = i, v
		}
	}
	banks[maxIdx] = 0
	for i := 1; i <= maxVal; i++ {
		banks[(maxIdx+i)%n]++
	}
}

func (d Day6) Part2(input string) (string, error) {
	banks := parseBanks(input)
	seen := make(map[string]int)
	steps := 0

	for {
		key := banksKey(banks)
		if prevStep, found := seen[key]; found {
			return strconv.Itoa(steps - prevStep), nil
		}
		seen[key] = steps
		redistribute(banks)
		steps++
	}
}

func init() {
	solve.Register(Day6{})
}
