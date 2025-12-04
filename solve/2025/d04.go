package solve2025

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 4}
}

func (d Day4) Part1(input string) (string, error) {
	grid := d.parseGrid(input)
	count := countAccessibleRolls(grid)
	return strconv.FormatInt(int64(count), 10), nil
}

func (d Day4) parseGrid(input string) [][]rune {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}
	return grid
}

type vec2 struct{ X, Y int }

// 8 directions: up, down, left, right, and 4 diagonals
var dir8 = []vec2{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // up, down, left, right
	{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, // diagonals
}

func countAdjacentRolls(grid [][]rune, r, c int) int {
	rows := len(grid)
	cols := len(grid[0])
	count := 0

	for _, dir := range dir8 {
		nr, nc := r+dir.X, c+dir.Y
		if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '@' {
			count++
		}
	}

	return count
}

func countAccessibleRolls(grid [][]rune) int {
	if len(grid) == 0 {
		return 0
	}

	rows := len(grid)
	cols := len(grid[0])
	count := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == '@' && countAdjacentRolls(grid, r, c) < 4 {
				count++
			}
		}
	}

	return count
}

func (d Day4) Part2(input string) (string, error) {
	grid := d.parseGrid(input)
	count := removeAccessibleRolls(grid)
	return strconv.FormatInt(int64(count), 10), nil
}

func removeAccessibleRolls(grid [][]rune) int {
	if len(grid) == 0 {
		return 0
	}

	rows := len(grid)
	cols := len(grid[0])
	totalRemoved := 0

	for {
		accessible := make([][2]int, 0)

		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				if grid[r][c] == '@' && countAdjacentRolls(grid, r, c) < 4 {
					accessible = append(accessible, [2]int{r, c})
				}
			}
		}

		if len(accessible) == 0 {
			break
		}

		for _, pos := range accessible {
			grid[pos[0]][pos[1]] = '.'
		}

		totalRemoved += len(accessible)
	}

	return totalRemoved
}

func init() {
	solve.Register(Day4{})
}
