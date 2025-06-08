package solve2019

import (
	"strconv"

	"aoc/solve"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 5}
}

func (d Day5) Part1(input string) (string, error) {
	ic := parseIntcode(input)
	outputs := ic.Run(1)
	return strconv.FormatInt(outputs[len(outputs)-1], 10), nil
}

func (d Day5) Part2(input string) (string, error) {
	ic := parseIntcode(input)
	outputs := ic.Run(5)
	return strconv.FormatInt(outputs[len(outputs)-1], 10), nil
}

func init() {
	solve.Register(Day5{})
}
