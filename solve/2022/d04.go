package solve2022

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 4}
}

func (d Day4) parseInput(input string) (part1, part2 string, err error) {
	lines := strings.Fields(strings.TrimSpace(input))

	countPart1 := 0
	countPart2 := 0

	for _, line := range lines {
		var a1, a2, b1, b2 int
		_, err := fmt.Sscanf(line, "%d-%d,%d-%d", &a1, &a2, &b1, &b2)
		if err != nil {
			return "", "", err
		}
		// Check if one range fully contains the other
		if (a1 <= b1 && a2 >= b2) || (b1 <= a1 && b2 >= a2) {
			countPart1++
		}

		// Check if ranges overlap at all
		if a1 <= b2 && b1 <= a2 {
			countPart2++
		}
	}

	return strconv.Itoa(countPart1), strconv.Itoa(countPart2), nil
}

func (d Day4) Part1(input string) (string, error) {
	count, _, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	return count, nil
}

func (d Day4) Part2(input string) (string, error) {
	_, count, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	return count, nil
}

func init() {
	solve.Register(Day4{})
}
