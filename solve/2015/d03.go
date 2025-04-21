package solve2015

import (
	"strconv"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 3}
}

func (d Day3) solve(input string) (part1, part2 int) {
	step := grids.DirectionsComplex("^>v<")

	var z complex128
	seen := map[complex128]bool{z: true}

	for _, c := range input {
		z += step[c]
		seen[z] = true
	}
	part1 = len(seen)

	seen = map[complex128]bool{0: true}

	z = 0
	for i := 0; i < len(input); i += 2 {
		z += step[rune(input[i])]
		seen[z] = true
	}

	z = 0
	for i := 1; i < len(input); i += 2 {
		z += step[rune(input[i])]
		seen[z] = true
	}

	part2 = len(seen)

	return part1, part2
}

func (d Day3) Part1(input string) (string, error) {
	part1, _ := d.solve(input)
	return strconv.Itoa(part1), nil
}

func (d Day3) Part2(input string) (string, error) {
	_, part2 := d.solve(input)
	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day3{})
}
