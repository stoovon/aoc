package solve2017

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day13 struct {
}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 13}
}

func parseFirewall(input string) map[int]int {
	layers := make(map[int]int)
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		depth, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		rng, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		layers[depth] = rng
	}
	return layers
}

func isCaught(layers map[int]int, delay int) bool {
	for depth, rng := range layers {
		period := 2 * (rng - 1)
		if (depth+delay)%period == 0 {
			return true
		}
	}
	return false
}

func (d Day13) Part1(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("empty input")
	}
	layers := parseFirewall(input)
	severity := 0
	for depth, rng := range layers {
		period := 2 * (rng - 1)
		if depth%period == 0 {
			severity += depth * rng
		}
	}
	return strconv.Itoa(severity), nil
}

func (d Day13) Part2(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("empty input")
	}
	layers := parseFirewall(input)
	delay := 0
	for {
		if !isCaught(layers, delay) {
			return strconv.Itoa(delay), nil
		}
		delay++
	}
}

func init() {
	solve.Register(Day13{})
}
