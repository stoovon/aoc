package solve2021

import (
	"aoc/solve"
	"errors"
	"fmt"
)

type Day11 struct{}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 11}
}

func (d Day11) Part1(input string) (string, error) {
	grid, err := d.parseGrid(input)
	if err != nil {
		return "", err
	}
	totalFlashes := 0
	for range 100 {
		totalFlashes += step(grid)
	}
	return fmt.Sprintf("%d", totalFlashes), nil
}

func (d Day11) parseGrid(input string) ([][]int, error) {
	lines := make([]string, 0)
	for _, line := range splitLines(input) {
		if len(line) == 0 {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) != 10 {
		return nil, errors.New("input must have 10 lines")
	}
	grid := make([][]int, 10)
	for r, line := range lines {
		if len(line) != 10 {
			return nil, errors.New("each line must have 10 digits")
		}
		grid[r] = make([]int, 10)
		for c, ch := range line {
			grid[r][c] = int(ch - '0')
		}
	}
	return grid, nil
}

func splitLines(s string) []string {
	lines := make([]string, 0)
	start := 0
	for i, ch := range s {
		if ch == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func step(grid [][]int) int {
	flashed := make([][]bool, 10)
	for i := range flashed {
		flashed[i] = make([]bool, 10)
	}
	var flashCount int
	// Increase all by 1
	for r := range 10 {
		for c := range 10 {
			grid[r][c]++
		}
	}
	// Process flashes
	var flash func(int, int)
	flash = func(r, c int) {
		if r < 0 || r >= 10 || c < 0 || c >= 10 {
			return
		}
		if flashed[r][c] || grid[r][c] <= 9 {
			return
		}
		flashed[r][c] = true
		flashCount++
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				nr, nc := r+dr, c+dc
				if nr < 0 || nr >= 10 || nc < 0 || nc >= 10 {
					continue
				}
				grid[nr][nc]++
				if !flashed[nr][nc] && grid[nr][nc] > 9 {
					flash(nr, nc)
				}
			}
		}
	}
	for r := range 10 {
		for c := range 10 {
			if grid[r][c] > 9 && !flashed[r][c] {
				flash(r, c)
			}
		}
	}
	// Set flashed cells to 0
	for r := range 10 {
		for c := range 10 {
			if flashed[r][c] {
				grid[r][c] = 0
			}
		}
	}
	return flashCount
}

func (d Day11) Part2(input string) (string, error) {
	grid, err := d.parseGrid(input)
	if err != nil {
		return "", err
	}
	stepNum := 1
	for {
		flashes := step(grid)
		if flashes == 100 {
			return fmt.Sprintf("%d", stepNum), nil
		}
		stepNum++
	}
}

func init() {
	solve.Register(Day11{})
}
