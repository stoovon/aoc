package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day25 struct {
}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 25}
}

func codeNumber(row, col int) int {
	codeLevel := row + (col - 1)
	return col + (codeLevel*(codeLevel-1))/2
}

func code(row, col int) int {
	const code1 = 20151125
	const mult = 252533
	const modulus = 33554393

	count := codeNumber(row, col)
	code := code1
	for i := 1; i < count; i++ {
		code = (code * mult) % modulus
	}
	return code
}

func (d Day25) Part1(input string) (string, error) {
	parts := strings.Split(strings.TrimSpace(input), " ")
	row, err := strconv.Atoi(strings.TrimSuffix(parts[16], ","))
	if err != nil {
		return "", err
	}

	col, err := strconv.Atoi(strings.TrimSuffix(parts[18], "."))
	if err != nil {
		return "", err
	}

	result := code(row, col)
	return strconv.Itoa(result), nil
}

func (d Day25) Part2(_ string) (string, error) {
	return "", nil
}

func init() {
	solve.Register(Day25{})
}
