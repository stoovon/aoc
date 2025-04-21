package solve2017

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 1}
}

func (d Day1) run(input string, part int) (string, error) {
	input = strings.TrimSpace(input)

	digits := input
	sum := 0
	length := len(digits)
	offset := 1
	if part == 2 {
		offset = length / 2
	}

	for i := 0; i < length; i++ {
		if digits[i] == digits[(i+offset)%length] {
			digit, _ := strconv.Atoi(string(digits[i]))
			sum += digit
		}
	}

	return strconv.Itoa(sum), nil
}

func (d Day1) Part1(input string) (string, error) {
	return d.run(input, 1)
}

func (d Day1) Part2(input string) (string, error) {
	return d.run(input, 2)
}

func init() {
	solve.Register(Day1{})
}
