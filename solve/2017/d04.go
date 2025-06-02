package solve2017

import (
	"fmt"
	"sort"
	"strings"

	"aoc/solve"
	"aoc/utils/acstrings"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 4}
}

func (d Day4) Part1(input string) (string, error) {
	lines := acstrings.Lines(strings.TrimSpace(input))
	valid := 0
	for _, line := range lines {
		words := make(map[string]bool)
		duplicate := false
		for _, word := range strings.Fields(line) {
			if words[word] {
				duplicate = true
				break
			}
			words[word] = true
		}
		if !duplicate {
			valid++
		}
	}
	return fmt.Sprintf("%d", valid), nil
}

func sortString(s string) string {
	letters := strings.Split(s, "")
	sort.Strings(letters)
	return strings.Join(letters, "")
}

func (d Day4) Part2(input string) (string, error) {
	lines := acstrings.Lines(strings.TrimSpace(input))
	valid := 0
	for _, line := range lines {
		seen := make(map[string]bool)
		duplicate := false
		for _, word := range strings.Fields(line) {
			sorted := sortString(word)
			if seen[sorted] {
				duplicate = true
				break
			}
			seen[sorted] = true
		}
		if !duplicate {
			valid++
		}
	}
	return fmt.Sprintf("%d", valid), nil
}

func init() {
	solve.Register(Day4{})
}
