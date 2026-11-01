package solve2024

import (
	"errors"
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day10 struct {
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

func (day Day10) solve(input string, distinctPeaks bool) (string, error) {
	grid, err := parseTrailMap(input)
	if err != nil {
		return "", err
	}
	total := 0
	for p, height := range grid {
		if height == '0' {
			var seen map[image.Point]bool
			if distinctPeaks {
				seen = make(map[image.Point]bool)
			}
			total += day.dfs(grid, p, seen)
		}
	}
	return strconv.Itoa(total), nil
}

func (day Day10) Part1(input string) (string, error) {
	return day.solve(input, true)
}

func (day Day10) Part2(input string) (string, error) {
	return day.solve(input, false)
}

func parseTrailMap(input string) (map[image.Point]rune, error) {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.TrimSuffix(input, "\n")
	rows := strings.Split(input, "\n")
	width := len(rows[0])
	if width == 0 {
		return nil, errors.New("trail map must not be empty")
	}
	grid := make(map[image.Point]rune)
	for y, row := range rows {
		if len(row) != width {
			return nil, errors.New("trail map must have nonempty rows of equal width")
		}
		for x := range row {
			if row[x] < '0' || row[x] > '9' {
				return nil, errors.New("trail map must contain only decimal digits")
			}
			grid[image.Pt(x, y)] = rune(row[x])
		}
	}
	if len(grid) < 10 {
		return nil, errors.New("trail map must contain at least ten cells")
	}
	return grid, nil
}

func init() {
	solve.Register(Day10{})
}
