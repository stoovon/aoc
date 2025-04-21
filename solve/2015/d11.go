package solve2015

import (
	"strings"

	"aoc/solve"
)

type Day11 struct {
}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 11}
}

func valid(s string) bool {
	// Check for invalid characters
	if strings.ContainsAny(s, "iol") {
		return false
	}

	count := 0
	flag := false
	seen := make(map[rune]bool)

	// Check for two non-overlapping pairs
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] && !seen[rune(s[i])] {
			count++
			seen[rune(s[i])] = true
		}
	}

	// Check for a sequence of three consecutive increasing characters
	for i := 0; i < len(s)-2; i++ {
		if s[i]+1 == s[i+1] && s[i+1]+1 == s[i+2] {
			flag = true
		}
	}

	return count >= 2 && flag
}

func next(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == 'z' {
			runes[i] = 'a'
		} else {
			runes[i]++
			break
		}
	}
	return string(runes)
}

func (d Day11) findNext(s string, changes int) string {
	s = strings.TrimSpace(s)

	for s != "zzzzzzzz" {
		s = next(s)
		if valid(s) {
			if changes == 0 {
				return s
			}
			changes--
		}
	}
	return s
}

func (d Day11) Part1(input string) (string, error) {
	return d.findNext(input, 0), nil
}

func (d Day11) Part2(input string) (string, error) {
	return d.findNext(input, 1), nil
}

func init() {
	solve.Register(Day11{})
}
