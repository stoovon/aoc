package solve2020

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day6 struct{}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 6}
}

func countGroupAnswers(input string, allMustAgree bool) int {
	groups := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n\n")
	sum := 0
	for _, group := range groups {
		lines := strings.Fields(group)
		if len(lines) == 0 {
			continue
		}
		counts := make(map[rune]int)
		for _, line := range lines {
			for _, c := range line {
				if c >= 'a' && c <= 'z' {
					counts[c]++
				}
			}
		}
		if allMustAgree {
			for _, v := range counts {
				if v == len(lines) {
					sum++
				}
			}
		} else {
			sum += len(counts)
		}
	}
	return sum
}

func (d Day6) Part1(input string) (string, error) {
	return strconv.Itoa(countGroupAnswers(input, false)), nil
}

func (d Day6) Part2(input string) (string, error) {
	return strconv.Itoa(countGroupAnswers(input, true)), nil
}

func init() {
	solve.Register(Day6{})
}
