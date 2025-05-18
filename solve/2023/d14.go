package solve2023

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day14 struct {
}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 14}
}

func (d Day14) parseGrid(input string) [][]rune {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}
	return grid
}

func (d Day14) tiltNorth(grid [][]rune) {
	rows, cols := len(grid), len(grid[0])
	for col := 0; col < cols; col++ {
		for row := 1; row < rows; row++ {
			if grid[row][col] == 'O' {
				newRow := row
				for newRow > 0 && grid[newRow-1][col] == '.' {
					grid[newRow][col], grid[newRow-1][col] = '.', 'O'
					newRow--
				}
			}
		}
	}
}

func (d Day14) tiltSouth(grid [][]rune) {
	rows, cols := len(grid), len(grid[0])
	for col := 0; col < cols; col++ {
		for row := rows - 2; row >= 0; row-- {
			if grid[row][col] == 'O' {
				newRow := row
				for newRow < rows-1 && grid[newRow+1][col] == '.' {
					grid[newRow][col], grid[newRow+1][col] = '.', 'O'
					newRow++
				}
			}
		}
	}
}

func (d Day14) tiltWest(grid [][]rune) {
	rows, cols := len(grid), len(grid[0])
	for row := 0; row < rows; row++ {
		for col := 1; col < cols; col++ {
			if grid[row][col] == 'O' {
				newCol := col
				for newCol > 0 && grid[row][newCol-1] == '.' {
					grid[row][newCol], grid[row][newCol-1] = '.', 'O'
					newCol--
				}
			}
		}
	}
}

func (d Day14) tiltEast(grid [][]rune) {
	rows, cols := len(grid), len(grid[0])
	for row := 0; row < rows; row++ {
		for col := cols - 2; col >= 0; col-- {
			if grid[row][col] == 'O' {
				newCol := col
				for newCol < cols-1 && grid[row][newCol+1] == '.' {
					grid[row][newCol], grid[row][newCol+1] = '.', 'O'
					newCol++
				}
			}
		}
	}
}

func (d Day14) spinCycle(grid [][]rune) {
	d.tiltNorth(grid)
	d.tiltWest(grid)
	d.tiltSouth(grid)
	d.tiltEast(grid)
}

func (d Day14) northLoad(grid [][]rune) int {
	load := 0
	rows := len(grid)
	for row := 0; row < rows; row++ {
		for col := 0; col < len(grid[row]); col++ {
			if grid[row][col] == 'O' {
				load += rows - row
			}
		}
	}
	return load
}

func (d Day14) gridFingerprint(grid [][]rune) string {
	var sb strings.Builder
	for _, row := range grid {
		sb.WriteString(string(row))
	}
	return sb.String()
}

func (d Day14) Part1(input string) (string, error) {
	grid := d.parseGrid(input)
	d.tiltNorth(grid)
	return strconv.Itoa(d.northLoad(grid)), nil
}

func (d Day14) Part2(input string) (string, error) {
	grid := d.parseGrid(input)
	cache := make(map[string]int)
	var cycleFoundAt int
	const totalCycles = 1000000000

	for j := 0; j < 1000; j++ {
		d.spinCycle(grid)
		fp := d.gridFingerprint(grid)
		if _, ok := cache[fp]; ok {
			if cycleFoundAt == 0 {
				cycleFoundAt = j
				cache = map[string]int{fp: j}
				continue
			}
			remaining := totalCycles - j - 1
			remaining %= j - cycleFoundAt
			for i := 0; i < remaining; i++ {
				d.spinCycle(grid)
			}
			return strconv.Itoa(d.northLoad(grid)), nil
		}
		cache[fp] = j
	}
	return "", errors.New("no cycle found")
}

func init() {
	solve.Register(Day14{})
}
