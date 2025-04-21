package solve2020

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 1}
}

func (d Day1) parseInput(input string) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	numbers := make([]int, len(lines))

	for i, line := range lines {
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		numbers[i] = num
	}

	return numbers, nil
}

func (d Day1) Part1(input string) (string, error) {
	numbers, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(numbers); i++ {
		for j := i + 1; j < len(numbers); j++ {
			if numbers[i]+numbers[j] == 2020 {
				return strconv.Itoa(numbers[i] * numbers[j]), nil
			}
		}
	}

	return "", nil
}

func (d Day1) Part2(input string) (string, error) {
	numbers, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(numbers); i++ {
		for j := i + 1; j < len(numbers); j++ {
			for k := j + 1; k < len(numbers); k++ {
				if numbers[i]+numbers[j]+numbers[k] == 2020 {
					return strconv.Itoa(numbers[i] * numbers[j] * numbers[k]), nil
				}
			}
		}
	}

	return "", nil
}

func init() {
	solve.Register(Day1{})
}
