package solve2015

import (
	"aoc/solve"
	"strconv"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 1}
}

func (d Day1) solve(input string) (floor, basement int) {
	direction := map[rune]int{'(': 1, ')': -1}
	basement = -1
	floor = 0

	for i, c := range input {
		floor += direction[c]
		if basement == -1 && floor == -1 {
			basement = i + 1
		}
	}

	return floor, basement
}

func (d Day1) Part1(input string) (string, error) {
	answer, _ := d.solve(input)
	return strconv.Itoa(answer), nil
}

func (d Day1) Part2(input string) (string, error) {
	_, answer := d.solve(input)
	return strconv.Itoa(answer), nil
}

func init() {
	solve.Register(Day1{})
}
