package solve2016

import (
	"errors"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 2}
}

func (d Day2) Part1(input string) (string, error) {
	return "", errors.New("Not implemented")
}

func (d Day2) Part2(input string) (string, error) {
	return "", errors.New("Not implemented")
}

func init() {
	solve.Register(Day2{})
}
