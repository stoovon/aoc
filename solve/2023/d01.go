package solve2023

import (
	"aoc/solve"
	"errors"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 1}
}

func (d Day1) Part1(input string) (string, error) {
	return "", errors.New("Not implemented")
}

func (d Day1) Part2(input string) (string, error) {
	return "", errors.New("Not implemented")
}

func init() {
	solve.Register(Day1{})
}
