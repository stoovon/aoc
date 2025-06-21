package solve2017

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day17 struct {
}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 17}
}

func parseSteps(input string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(input))
}

func (d Day17) Part1(input string) (string, error) {
	steps, err := parseSteps(input)
	if err != nil {
		return "", err
	}

	buffer := []int{0}
	pos := 0

	for i := 1; i <= 2017; i++ {
		pos = (pos + steps) % len(buffer)
		pos++
		buffer = append(buffer[:pos], append([]int{i}, buffer[pos:]...)...)
	}

	for idx, v := range buffer {
		if v == 2017 {
			return strconv.Itoa(buffer[(idx+1)%len(buffer)]), nil
		}
	}
	return "", errors.New("2017 not found in buffer")
}

func (d Day17) Part2(input string) (string, error) {
	steps, err := parseSteps(input)
	if err != nil {
		return "", err
	}

	pos := 0
	result := 0
	for i := 1; i <= 50000000; i++ {
		pos = (pos + steps) % i
		if pos == 0 {
			result = i
		}
		pos++
	}
	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day17{})
}
