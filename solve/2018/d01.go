package solve2018

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day1 struct{}

func (d Day1) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2018, Day: 1} }

func parseFrequencyChanges(input string) ([]int64, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("puzzle input is empty")
	}
	input = strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	changes := make([]int64, 0, len(lines))
	for index, line := range lines {
		value := strings.TrimSpace(line)
		digits := value
		if len(value) > 0 && (value[0] == '-' || value[0] == '+') {
			digits = value[1:]
		}
		valid := digits != ""
		for _, c := range digits {
			if c < '0' || c > '9' {
				valid = false
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("line %d: expected one signed decimal frequency change", index+1)
		}
		change, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", index+1, err)
		}
		changes = append(changes, change)
	}
	return changes, nil
}
func (d Day1) Part1(input string) (string, error) {
	changes, err := parseFrequencyChanges(input)
	if err != nil {
		return "", err
	}
	var total, change big.Int
	for _, value := range changes {
		total.Add(&total, change.SetInt64(value))
	}
	return total.String(), nil
}

type frequencyPrefix struct {
	normalized, remainder *big.Int
	offset                int
}

func firstRepeatedFrequency(changes []int64) (*big.Int, error) {
	if len(changes) == 0 {
		return nil, fmt.Errorf("puzzle input is empty")
	}
	seen := map[string]bool{"0": true}
	prefixes := make([]frequencyPrefix, 0, len(changes))
	var total, change big.Int
	for offset, value := range changes {
		prefixes = append(prefixes, frequencyPrefix{normalized: new(big.Int).Set(&total), offset: offset})
		total.Add(&total, change.SetInt64(value))
		key := total.String()
		if seen[key] {
			return new(big.Int).Set(&total), nil
		}
		seen[key] = true
	}
	// Zero drift necessarily repeats zero during the first-pass check above.
	shift := new(big.Int).Abs(&total)
	for i := range prefixes {
		if total.Sign() < 0 {
			prefixes[i].normalized.Neg(prefixes[i].normalized)
		}
		prefixes[i].remainder = new(big.Int).Mod(prefixes[i].normalized, shift)
	}
	sort.Slice(prefixes, func(i, j int) bool {
		if compare := prefixes[i].remainder.Cmp(prefixes[j].remainder); compare != 0 {
			return compare < 0
		}
		return prefixes[i].normalized.Cmp(prefixes[j].normalized) < 0
	})
	var bestCycles, answer *big.Int
	bestOffset := 0
	for i := 1; i < len(prefixes); i++ {
		low, high := prefixes[i-1], prefixes[i]
		if low.remainder.Cmp(high.remainder) != 0 {
			continue
		}
		// The next higher prefix is the earliest target for this phase.
		cycles := new(big.Int).Sub(high.normalized, low.normalized)
		cycles.Quo(cycles, shift)
		if bestCycles == nil || cycles.Cmp(bestCycles) < 0 || (cycles.Cmp(bestCycles) == 0 && low.offset < bestOffset) {
			bestCycles = cycles
			bestOffset = low.offset
			answer = new(big.Int).Set(high.normalized)
			if total.Sign() < 0 {
				answer.Neg(answer)
			}
		}
	}
	if answer == nil {
		return nil, fmt.Errorf("frequency never repeats")
	}
	return answer, nil
}
func (d Day1) Part2(input string) (string, error) {
	changes, err := parseFrequencyChanges(input)
	if err != nil {
		return "", err
	}
	frequency, err := firstRepeatedFrequency(changes)
	if err != nil {
		return "", err
	}
	return frequency.String(), nil
}
func init() { solve.Register(Day1{}) }
