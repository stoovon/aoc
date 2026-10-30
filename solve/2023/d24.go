package solve2023

import (
	"math/big"
	"strconv"

	"aoc/solve"
)

type Day24 struct{}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 24}
}

func (d Day24) Part1(input string) (string, error) {
	count, err := hailCountCrossings(input, big.NewInt(200000000000000), big.NewInt(400000000000000))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(count), nil
}

func (d Day24) Part2(input string) (string, error) {
	return hailExactPart2(input)
}

func init() { solve.Register(Day24{}) }
