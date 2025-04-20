package solve2024

import (
	"image"
	"strconv"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day10 struct {
}

type day10DS struct {
	part1 int
	part2 int
}

func (day Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 10}
}

func (day Day10) dfs(grid map[image.Point]rune, p image.Point, seen map[image.Point]bool) (score int) {
	if grid[p] == '9' {
		if seen[p] {
			return 0
		} else if seen != nil {
			seen[p] = true
		}
		return 1
	}
	for _, d := range grids.URDL() {
		if n := p.Add(d); grid[n] == grid[p]+1 {
			score += day.dfs(grid, n, seen)
		}
	}
	return score
}

func (day Day10) solve(grid map[image.Point]rune) day10DS {
	part1, part2 := 0, 0
	for p := range grid {
		if grid[p] == '0' {
			part1 += day.dfs(grid, p, map[image.Point]bool{})
			part2 += day.dfs(grid, p, nil)
		}
	}

	return day10DS{
		part1: part1,
		part2: part2,
	}
}

func (day Day10) Part1(input string) (string, error) {
	grid := grids.ParseCharsFromTopLeft(input)

	return strconv.Itoa(day.solve(grid).part1), nil
}

func (day Day10) Part2(input string) (string, error) {
	grid := grids.ParseCharsFromTopLeft(input)

	return strconv.Itoa(day.solve(grid).part2), nil
}

func init() {
	solve.Register(Day10{})
}
