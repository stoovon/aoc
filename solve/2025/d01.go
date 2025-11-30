package solve2025

import (
	"errors"

	"aoc/solve"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 1}
}

func (d Day1) Part1(input string) (string, error) {
	return "", errors.New("not implemented")
}

func (d Day1) Part2(input string) (string, error) {
	return "", errors.New("not implemented")
}

func init() {
	solve.Register(Day1{})
}
