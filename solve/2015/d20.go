package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day20 struct {
}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 20}
}

func findFirstHouse(goal int, part2 bool) int {
	const BIG_NUM = 1000000
	houses := make([]int, BIG_NUM)

	for elf := 1; elf < BIG_NUM; elf++ {
		limit := BIG_NUM
		if part2 {
			limit = elf * 50 // Elves stop after 50 houses in Part 2
		}
		for house := elf; house < limit && house < BIG_NUM; house += elf {
			if part2 {
				houses[house] += 11 * elf
			} else {
				houses[house] += 10 * elf
			}
		}
	}

	for i, presents := range houses {
		if presents >= goal {
			return i
		}
	}
	return -1
}

func (d Day20) Part1(input string) (string, error) {
	goal, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(findFirstHouse(goal, false)), nil
}

func (d Day20) Part2(input string) (string, error) {
	goal, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(findFirstHouse(goal, true)), nil
}

func init() {
	solve.Register(Day20{})
}
