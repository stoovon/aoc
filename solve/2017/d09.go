package solve2017

import (
	"errors"
	"strconv"

	"aoc/solve"
)

type Day9 struct {
}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 9}
}

func processStream(input string) (score int, garbageCount int, err error) {
	depth := 0
	inGarbage := false
	skipNext := false

	for _, c := range input {
		if skipNext {
			skipNext = false
			continue
		}
		switch c {
		case '!':
			skipNext = true
		case '>':
			if inGarbage {
				inGarbage = false
			}
		case '<':
			if !inGarbage {
				inGarbage = true
			} else {
				garbageCount++
			}
		case '{':
			if !inGarbage {
				depth++
			} else {
				garbageCount++
			}
		case '}':
			if !inGarbage {
				score += depth
				depth--
			} else {
				garbageCount++
			}
		default:
			if inGarbage {
				garbageCount++
			}
		}
	}
	if depth != 0 {
		return 0, 0, errors.New("unbalanced groups")
	}
	return score, garbageCount, nil
}

func (d Day9) Part1(input string) (string, error) {
	score, _, err := processStream(input)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(score), nil
}

func (d Day9) Part2(input string) (string, error) {
	_, garbageCount, err := processStream(input)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(garbageCount), nil
}

func init() {
	solve.Register(Day9{})
}
