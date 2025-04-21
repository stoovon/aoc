package solve2019

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 1}
}

func fuel(mass int) int {
	return mass/3 - 2
}

func (d Day1) run(input string) (int, int) {
	totalPart1, totalPart2 := 0, 0
	lines := strings.Split(strings.TrimSpace(input), "\n")

	for _, line := range lines {
		mass, _ := strconv.Atoi(line)

		totalPart1 += mass/3 - 2

		for f := fuel(mass); f > 0; f = fuel(f) {
			totalPart2 += f
		}
	}

	return totalPart1, totalPart2
}

func (d Day1) Part1(input string) (string, error) {
	part1, _ := d.run(input)
	return strconv.Itoa(part1), nil
}

func (d Day1) Part2(input string) (string, error) {
	_, part2 := d.run(input)
	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day1{})
}
