package solve2019

import (
	"strconv"

	"aoc/solve"
)

type Day9 struct {
}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 9}
}

func (d Day9) Part1(input string) (string, error) {
	code := parseIntcode(input)
	out := code.Run(1)
	return strconv.FormatInt(out[len(out)-1], 10), nil
}
func (d Day9) Part2(input string) (string, error) {
	code := parseIntcode(input)
	out := code.Run(2)
	return strconv.FormatInt(out[len(out)-1], 10), nil

}

func init() {
	solve.Register(Day9{})
}
