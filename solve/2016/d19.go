package solve2016

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day19 struct {
}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 19}
}

func (d Day19) Part1(input string) (string, error) {
	input = strings.TrimSpace(input)

	n, err := strconv.Atoi(input)
	if err != nil {
		return "", err
	}

	// Josephus problem solution for k=2
	winner := 0
	for i := 1; i <= n; i++ {
		winner = (winner + 2) % i
	}

	return strconv.Itoa(winner + 1), nil
}

func (d Day19) Part2(input string) (string, error) {
	input = strings.TrimSpace(input)

	n, err := strconv.Atoi(input)
	if err != nil {
		return "", err
	}

	// Josephus problem for k=3 (stealing across)
	power := 1
	for power*3 <= n {
		power *= 3
	}

	// Calculate the winner
	winner := n - power + max(0, n-2*power)
	return strconv.Itoa(winner), nil
}

func init() {
	solve.Register(Day19{})
}
