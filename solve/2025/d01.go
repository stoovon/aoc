package solve2025

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 1}
}

func (d Day1) Part1(input string) (string, error) {
	return d.solve(input, false)
}

func (d Day1) Part2(input string) (string, error) {
	return d.solve(input, true)
}

func (d Day1) solve(input string, countAnyZero bool) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	position := 50
	zerosCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		direction := line[0]
		distance, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}

		if countAnyZero {
			wraps := distance / 100
			zerosCount += wraps

			remainingDistance := distance % 100
			if direction == 'L' {
				// Moving left: we cross 0 if we go below 0, but not if we start at 0
				if position > 0 && position <= remainingDistance {
					zerosCount++
				}
			} else {
				// Moving right: we cross or land on 0 if we reach or pass 100
				if position+remainingDistance >= 100 {
					zerosCount++
				}
			}
		}

		switch direction {
		case 'L':
			position = (position - distance) % 100
			if position < 0 {
				position += 100
			}
		case 'R':
			position = (position + distance) % 100
		}

		if !countAnyZero && position == 0 {
			zerosCount++
		}
	}

	return strconv.Itoa(zerosCount), nil
}

func init() {
	solve.Register(Day1{})
}
