package solve2018

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 1}
}

func (d Day1) Part1(input string) (string, error) {
	answer := 0
	for _, line := range strings.Fields(input) {
		num, _ := strconv.Atoi(line)
		answer += num
	}
	return strconv.Itoa(answer), nil
}

func (d Day1) Part2(input string) (string, error) {
	freq := 0
	seen := map[int]bool{0: true}

	for {
		for _, line := range strings.Fields(input) {
			num, _ := strconv.Atoi(line)
			freq += num
			if seen[freq] {
				return strconv.Itoa(freq), nil
			}
			seen[freq] = true
		}
	}
}

func init() {
	solve.Register(Day1{})
}
