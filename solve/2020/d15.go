package solve2020

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day15 struct{}

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 15}
}

func memoryGame(input string, target int) (string, error) {
	input = strings.TrimSpace(input)
	parts := strings.Split(input, ",")
	lastSpoken := make([]int, target)
	var prev int

	for i, p := range parts[:len(parts)-1] {
		num, err := strconv.Atoi(p)
		if err != nil {
			return "", err
		}
		lastSpoken[num] = i + 1
	}
	prev, _ = strconv.Atoi(parts[len(parts)-1])

	for turn := len(parts) + 1; turn <= target; turn++ {
		last := lastSpoken[prev]
		lastSpoken[prev] = turn - 1
		if last == 0 {
			prev = 0
		} else {
			prev = (turn - 1) - last
		}
	}
	return strconv.Itoa(prev), nil
}

func (d Day15) Part1(input string) (string, error) {
	return memoryGame(input, 2020)
}

func (d Day15) Part2(input string) (string, error) {
	return memoryGame(input, 30000000)
}

func init() {
	solve.Register(Day15{})
}
