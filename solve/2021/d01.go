package solve2021

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 1}
}

func (d Day1) run(input string) (int, int) {
	var depths []int
	for _, s := range strings.Fields(string(input)) {
		i, _ := strconv.Atoi(s)
		depths = append(depths, i)
	}

	part1, part2 := 0, 0
	for i := range depths {
		if i >= 1 && depths[i] > depths[i-1] {
			part1++
		}
		if i >= 3 && depths[i] > depths[i-3] {
			part2++
		}
	}

	return part1, part2
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
