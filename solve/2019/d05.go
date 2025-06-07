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
	prog, err := parseIntCode(input)
	if err != nil {
		return "", err
	}
	ic := NewIntCode(prog, []int{1})
	outputs := ic.Run()
	return strconv.Itoa(outputs[len(outputs)-1]), nil
}

func (d Day5) Part2(input string) (string, error) {
	prog, err := parseIntCode(input)
	if err != nil {
		return "", err
	}
	ic := NewIntCode(prog, []int{5})
	outputs := ic.Run()
	return strconv.Itoa(outputs[len(outputs)-1]), nil
}

func init() {
	solve.Register(Day5{})
}
