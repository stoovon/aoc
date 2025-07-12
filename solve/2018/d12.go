package solve2018

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day12 struct{}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 12}
}

func parseInitialStateAndRules(input string) (map[int]byte, map[string]byte) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	init := strings.TrimPrefix(lines[0], "initial state: ")
	rules := map[string]byte{}
	for _, line := range lines[1:] {
		if len(line) < 9 {
			continue
		}
		parts := strings.Split(line, " => ")
		if len(parts) == 2 {
			rules[parts[0]] = parts[1][0]
		}
	}
	state := map[int]byte{}
	for i, c := range init {
		if c == '#' {
			state[i] = '#'
		}
	}
	return state, rules
}

func simulateGenerations(state map[int]byte, rules map[string]byte, generations int, detectStable bool) (int, bool) {
	prevSum := 0
	prevDelta := 0
	lastPattern := ""
	for gen := 1; gen <= generations; gen++ {
		newState := map[int]byte{}
		min, max := 0, 0
		for k := range state {
			if k < min {
				min = k
			}
			if k > max {
				max = k
			}
		}
		for i := min - 2; i <= max+2; i++ {
			pat := ""
			for d := -2; d <= 2; d++ {
				if state[i+d] == '#' {
					pat += "#"
				} else {
					pat += "."
				}
			}
			if rules[pat] == '#' {
				newState[i] = '#'
			}
		}
		state = newState
		sum := 0
		min, max = 0, 0
		for k := range state {
			if k < min {
				min = k
			}
			if k > max {
				max = k
			}
		}
		pat := ""
		for i := min; i <= max; i++ {
			if state[i] == '#' {
				pat += "#"
			} else {
				pat += "."
			}
		}
		for k := range state {
			if state[k] == '#' {
				sum += k
			}
		}
		if detectStable && prevDelta == sum-prevSum && pat == lastPattern {
			return sum, true
		}
		prevDelta = sum - prevSum
		prevSum = sum
		lastPattern = pat
	}
	sum := 0
	for k := range state {
		if state[k] == '#' {
			sum += k
		}
	}
	return sum, false
}

func (d Day12) Part1(input string) (string, error) {
	state, rules := parseInitialStateAndRules(input)
	sum, _ := simulateGenerations(state, rules, 20, false)
	return strconv.Itoa(sum), nil
}

func (d Day12) Part2(input string) (string, error) {
	state, rules := parseInitialStateAndRules(input)
	generations := int64(50000000000)

	prevSum := 0
	prevDelta := 0
	stableCount := 0

	for gen := 1; gen <= 200; gen++ { // Check for stability in first 200 generations
		newState := map[int]byte{}
		min, max := 0, 0
		for k := range state {
			if k < min {
				min = k
			}
			if k > max {
				max = k
			}
		}

		for i := min - 2; i <= max+2; i++ {
			pat := ""
			for d := -2; d <= 2; d++ {
				if state[i+d] == '#' {
					pat += "#"
				} else {
					pat += "."
				}
			}
			if rules[pat] == '#' {
				newState[i] = '#'
			}
		}
		state = newState

		sum := 0
		for k := range state {
			if state[k] == '#' {
				sum += k
			}
		}

		currentDelta := sum - prevSum
		if currentDelta == prevDelta && gen > 1 {
			stableCount++
			if stableCount >= 5 { // Confirm stability over 5 generations
				// Pattern is stable, calculate final result
				remaining := generations - int64(gen)
				finalSum := int64(sum) + remaining*int64(currentDelta)
				return strconv.FormatInt(finalSum, 10), nil
			}
		} else {
			stableCount = 0
		}

		prevDelta = currentDelta
		prevSum = sum
	}

	// If we get here, pattern didn't stabilize in 200 generations
	// This shouldn't happen for valid AoC inputs
	return "0", nil
}

func init() {
	solve.Register(Day12{})
}
