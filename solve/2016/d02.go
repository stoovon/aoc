package solve2016

import (
	"errors"
	"strings"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 2}
}

const off = '.'

func decode(keypad []string, instructions []string, x, y int) string {
	var result strings.Builder
	for _, line := range instructions {
		for _, command := range line {
			x, y = move(keypad, command, x, y)
		}
		result.WriteByte(keypad[y][x])
	}
	return result.String()
}

func move(keypad []string, direction rune, x, y int) (int, int) {
	switch direction {
	case 'L':
		if keypad[y][x-1] != off {
			x--
		}
	case 'R':
		if keypad[y][x+1] != off {
			x++
		}
	case 'U':
		if keypad[y-1][x] != off {
			y--
		}
	case 'D':
		if keypad[y+1][x] != off {
			y++
		}
	}
	return x, y
}

func (d Day2) Part1(input string) (string, error) {
	data := strings.TrimSpace(input)

	keypad := []string{
		".....",
		".123.",
		".456.",
		".789.",
		".....",
	}

	if keypad[2][2] != '5' {
		return "", errors.New("invalid keypad configuration")
	}

	instructions := strings.Split(data, "\n")
	return decode(keypad, instructions, 2, 2), nil
}

func (d Day2) Part2(input string) (string, error) {
	data := strings.TrimSpace(input)

	keypad := []string{
		".......",
		"...1...",
		"..234..",
		".56789.",
		"..ABC..",
		"...D...",
		".......",
	}

	instructions := strings.Split(data, "\n")
	return decode(keypad, instructions, 1, 3), nil
}

func init() {
	solve.Register(Day2{})
}
