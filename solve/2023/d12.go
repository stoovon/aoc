package solve2023

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day12 struct {
}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 12}
}

func permutations(spring string, counts []int) int64 {
	// Pad with '.' at the start, trim trailing '.'
	spring = "." + strings.TrimRight(spring, ".")
	s := []rune(spring)
	n := len(s)

	possible := make([]int64, n+1)
	possible[0] = 1
	for i := 0; i < n && s[i] != '#'; i++ {
		possible[i+1] = 1
	}

	for _, count := range counts {
		newPossible := make([]int64, n+1)
		chunk := 0
		for i, c := range s {
			if c != '.' {
				chunk++
			} else {
				chunk = 0
			}
			if c != '#' {
				newPossible[i+1] += newPossible[i]
			}
			if chunk >= count && s[i-count] != '#' {
				newPossible[i+1] += possible[i-count]
			}
		}
		possible = newPossible
	}
	return possible[n]
}

func (d Day12) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var sum int64
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return "", errors.New("invalid input")
		}
		spring := parts[0]
		countStrs := strings.Split(parts[1], ",")
		counts := make([]int, len(countStrs))
		for i, cs := range countStrs {
			n, _ := strconv.Atoi(cs)
			counts[i] = n
		}
		sum += permutations(spring, counts)
	}
	return strconv.FormatInt(sum, 10), nil
}

func (d Day12) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var sum int64
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return "", errors.New("invalid input")
		}
		spring := parts[0]
		countStrs := strings.Split(parts[1], ",")
		counts := make([]int, len(countStrs))
		for i, cs := range countStrs {
			n, _ := strconv.Atoi(cs)
			counts[i] = n
		}
		// Repeat spring and counts 5 times, separated by '?'
		fullSpring := spring
		for i := 0; i < 4; i++ {
			fullSpring += "?" + spring
		}
		fullCounts := make([]int, 0, 5*len(counts))
		for i := 0; i < 5; i++ {
			fullCounts = append(fullCounts, counts...)
		}
		sum += permutations(fullSpring, fullCounts)
	}
	return strconv.FormatInt(sum, 10), nil
}

func init() {
	solve.Register(Day12{})
}
