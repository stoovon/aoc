package solve2017

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"aoc/solve"
)

type Day19 struct {
}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 19}
}

func traverseDiagram(input string, collectLetters bool, countSteps bool) (string, error) {
	grid := strings.Split(input, "\n")
	if len(grid) == 0 {
		return "", errors.New("empty input")
	}
	width := 0
	for _, row := range grid {
		if len(row) > width {
			width = len(row)
		}
	}
	for i := range grid {
		if len(grid[i]) < width {
			grid[i] += strings.Repeat(" ", width-len(grid[i]))
		}
	}
	x, y := 0, 0
	for i, c := range grid[0] {
		if c == '|' {
			x = i
			break
		}
	}
	dx := []int{0, 1, 0, -1}
	dy := []int{1, 0, -1, 0}
	dir := 0 // down

	letters := []rune{}
	steps := 0
	for {
		c := grid[y][x]
		if c == ' ' {
			break
		}
		if collectLetters && unicode.IsLetter(rune(c)) {
			letters = append(letters, rune(c))
		}
		if countSteps {
			steps++
		}
		if c == '+' {
			for turn := 1; turn <= 3; turn += 2 {
				ndir := (dir + turn) % 4
				nx, ny := x+dx[ndir], y+dy[ndir]
				if nx >= 0 && nx < width && ny >= 0 && ny < len(grid) && grid[ny][nx] != ' ' {
					dir = ndir
					break
				}
			}
		}
		x += dx[dir]
		y += dy[dir]
		if x < 0 || x >= width || y < 0 || y >= len(grid) {
			break
		}
	}
	if collectLetters {
		return string(letters), nil
	}
	return strconv.Itoa(steps), nil
}

func (d Day19) Part1(input string) (string, error) {
	return traverseDiagram(input, true, false)
}

func (d Day19) Part2(input string) (string, error) {
	return traverseDiagram(input, false, true)
}

func init() {
	solve.Register(Day19{})
}
