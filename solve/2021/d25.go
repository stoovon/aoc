package solve2021

import (
	"aoc/solve"
	"fmt"
	"strconv"
	"strings"
)

type Day25 struct{}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 25}
}

func (d Day25) Part1(input string) (string, error) {
	grid := parseCucumbers(input)
	steps := numberOfSteps(grid)
	return strconv.Itoa(steps), nil
}

func (d Day25) Part2(input string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func parseCucumbers(input string) [][]rune {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}
	return grid
}

func numberOfSteps(grid [][]rune) int {
	steps := 0
	for moveStep(grid) {
		steps++
	}
	return steps + 1
}

func moveStep(grid [][]rune) bool {
	movedEast := moveHerd(grid, '>', 1, 0)
	movedSouth := moveHerd(grid, 'v', 0, 1)
	return movedEast || movedSouth
}

func moveHerd(grid [][]rune, herd rune, dx, dy int) bool {
	moved := false
	rows, cols := len(grid), len(grid[0])
	toMove := [][2][2]int{}

	for y := range rows {
		for x := range cols {
			if grid[y][x] == herd {
				nx, ny := (x+dx)%cols, (y+dy)%rows
				if grid[ny][nx] == '.' {
					toMove = append(toMove, [2][2]int{{x, y}, {nx, ny}})
				}
			}
		}
	}

	for _, move := range toMove {
		from, to := move[0], move[1]
		grid[to[1]][to[0]] = herd
		grid[from[1]][from[0]] = '.'
		moved = true
	}

	return moved
}

func init() {
	solve.Register(Day25{})
}
