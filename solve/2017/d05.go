package solve2017

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 5}
}

func (d Day5) Part1(input string) (string, error) {
	lines := strings.Fields(strings.TrimSpace(input))
	jumps := make([]int, len(lines))
	for i, line := range lines {
		n, err := strconv.Atoi(line)
		if err != nil {
			return "", err
		}
		jumps[i] = n
	}
	steps := 0
	for i := 0; i >= 0 && i < len(jumps); steps++ {
		offset := jumps[i]
		jumps[i]++
		i += offset
	}
	return fmt.Sprintf("%d", steps), nil
}

func (d Day5) Part2(input string) (string, error) {
	lines := strings.Fields(strings.TrimSpace(input))
	jumps := make([]int, len(lines))
	for i, line := range lines {
		n, err := strconv.Atoi(line)
		if err != nil {
			return "", err
		}
		jumps[i] = n
	}
	steps := 0
	for i := 0; i >= 0 && i < len(jumps); steps++ {
		offset := jumps[i]
		if offset >= 3 {
			jumps[i]--
		} else {
			jumps[i]++
		}
		i += offset
	}
	return fmt.Sprintf("%d", steps), nil
}

func init() {
	solve.Register(Day5{})
}
