package solve2020

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day11 struct{}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 11}
}

func simulateSeating(input string, visible bool, tolerance int) int {
	lines := strings.Fields(strings.ReplaceAll(input, "\r\n", "\n"))
	rows := len(lines)
	cols := len(lines[0])
	grid := make([][]byte, rows)
	for i := range lines {
		grid[i] = []byte(lines[i])
	}
	dirs := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	changed := true
	for changed {
		changed = false
		next := make([][]byte, rows)
		for i := range grid {
			next[i] = make([]byte, cols)
			copy(next[i], grid[i])
		}
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				if grid[r][c] == '.' {
					continue
				}
				occ := 0
				for _, d := range dirs {
					nr, nc := r+d[0], c+d[1]
					if visible {
						for nr >= 0 && nr < rows && nc >= 0 && nc < cols {
							if grid[nr][nc] == 'L' {
								break
							}
							if grid[nr][nc] == '#' {
								occ++
								break
							}
							nr += d[0]
							nc += d[1]
						}
					} else {
						if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '#' {
							occ++
						}
					}
				}
				if grid[r][c] == 'L' && occ == 0 {
					next[r][c] = '#'
					changed = true
				} else if grid[r][c] == '#' && occ >= tolerance {
					next[r][c] = 'L'
					changed = true
				}
			}
		}
		grid = next
	}
	count := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == '#' {
				count++
			}
		}
	}
	return count
}

func (d Day11) Part1(input string) (string, error) {
	return strconv.Itoa(simulateSeating(input, false, 4)), nil
}

func (d Day11) Part2(input string) (string, error) {
	return strconv.Itoa(simulateSeating(input, true, 5)), nil
}

func init() {
	solve.Register(Day11{})
}
