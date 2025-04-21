package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 5}
}

func vowelCount(s string) int {
	vowels := "aeiou"
	count := 0
	for _, c := range s {
		if strings.ContainsRune(vowels, c) {
			count++
		}
	}
	return count
}

func hasDouble(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			return true
		}
	}
	return false
}

func blocklisted(s string) bool {
	blocklist := []string{"ab", "cd", "pq", "xy"}
	for _, b := range blocklist {
		if strings.Contains(s, b) {
			return true
		}
	}
	return false
}

func hasPair(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		pair := s[i : i+2]
		if strings.Count(s, pair) > 1 {
			return true
		}
	}
	return false
}

func hasSkipRepeat(s string) bool {
	for i := 0; i < len(s)-2; i++ {
		if s[i] == s[i+2] {
			return true
		}
	}
	return false
}

func isNiceA(s string) bool {
	return vowelCount(s) >= 3 && hasDouble(s) && !blocklisted(s)
}

func isNiceB(s string) bool {
	return hasPair(s) && hasSkipRepeat(s)
}

func (d Day5) solve(input string) (part1, part2 int) {
	words := strings.Split(strings.TrimSpace(input), "\n")
	for _, word := range words {
		if isNiceA(word) {
			part1++
		}
		if isNiceB(word) {
			part2++
		}
	}
	return part1, part2
}

func (d Day5) Part1(input string) (string, error) {
	part1, _ := d.solve(input)
	return strconv.Itoa(part1), nil
}

func (d Day5) Part2(input string) (string, error) {
	_, part2 := d.solve(input)
	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day5{})
}
