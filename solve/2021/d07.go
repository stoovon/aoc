package solve2021

import (
	"aoc/solve"
	"aoc/utils/maths"
	"errors"
	"sort"
	"strconv"
	"strings"
)

type Day7 struct{}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 7}
}

func (d Day7) parse(input string) ([]int, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, errors.New("empty input")
	}
	parts := strings.Split(input, ",")
	positions := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		positions = append(positions, n)
	}
	return positions, nil
}

func (d Day7) Part1(input string) (string, error) {
	positions, err := d.parse(input)
	if err != nil {
		return "", err
	}
	sort.Ints(positions)
	median := positions[len(positions)/2]
	fuel := 0
	for _, pos := range positions {
		fuel += maths.Abs(pos - median)
	}
	return strconv.Itoa(fuel), nil
}

func (d Day7) Part2(input string) (string, error) {
	positions, err := d.parse(input)
	if err != nil {
		return "", err
	}
	minPos, maxPos := positions[0], positions[0]
	for _, pos := range positions {
		if pos < minPos {
			minPos = pos
		}
		if pos > maxPos {
			maxPos = pos
		}
	}
	minFuel := -1
	for align := minPos; align <= maxPos; align++ {
		fuel := 0
		for _, pos := range positions {
			dist := maths.Abs(pos - align)
			fuel += dist * (dist + 1) / 2
		}
		if minFuel == -1 || fuel < minFuel {
			minFuel = fuel
		}
	}
	return strconv.Itoa(minFuel), nil
}

func init() {
	solve.Register(Day7{})
}
