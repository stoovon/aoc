package solve2020

import (
	"fmt"
	"strings"

	"aoc/solve"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 3}
}

func (d Day3) Lines(input string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func (d Day3) Part1(input string) (string, error) {
	lines := d.Lines(input)
	if len(lines) == 0 {
		return "0", nil
	}
	width := len(lines[0])
	x, y, trees := 0, 0, 0
	for y < len(lines) {
		if lines[y][x%width] == '#' {
			trees++
		}
		x += 3
		y += 1
	}
	return fmt.Sprintf("%d", trees), nil
}

func countTrees(lines []string, right, down int) int {
	width := len(lines[0])
	x, y, trees := 0, 0, 0
	for y < len(lines) {
		if lines[y][x%width] == '#' {
			trees++
		}
		x += right
		y += down
	}
	return trees
}

func (d Day3) Part2(input string) (string, error) {
	lines := d.Lines(input)
	slopes := [][2]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	result := 1
	for _, slope := range slopes {
		result *= countTrees(lines, slope[0], slope[1])
	}
	return fmt.Sprintf("%d", result), nil
}

func init() {
	solve.Register(Day3{})
}
