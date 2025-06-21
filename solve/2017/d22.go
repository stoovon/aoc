package solve2017

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day22 struct {
}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 22}
}

func (d Day22) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	infected := make(map[[2]int]bool)
	n := len(lines)
	for y, row := range lines {
		for x, c := range row {
			if c == '#' {
				infected[[2]int{x, y}] = true
			}
		}
	}
	// Start in the middle
	x, y := n/2, n/2
	// Directions: up, right, down, left
	dirs := [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	dir := 0 // 0=up

	infections := 0
	for burst := 0; burst < 10000; burst++ {
		pos := [2]int{x, y}
		if infected[pos] {
			// Turn right
			dir = (dir + 1) % 4
			// Clean node
			delete(infected, pos)
		} else {
			// Turn left
			dir = (dir + 3) % 4
			// Infect node
			infected[pos] = true
			infections++
		}
		// Move forward
		x += dirs[dir][0]
		y += dirs[dir][1]
	}
	return strconv.Itoa(infections), nil
}

func (d Day22) Part2(input string) (string, error) {
	const (
		Clean    = 0
		Weakened = 1
		Infected = 2
		Flagged  = 3
	)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make(map[[2]int]int)
	n := len(lines)
	for y, row := range lines {
		for x, c := range row {
			if c == '#' {
				grid[[2]int{x, y}] = Infected
			}
		}
	}
	x, y := n/2, n/2
	dirs := [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	dir := 0 // up
	infections := 0
	for burst := 0; burst < 10000000; burst++ {
		pos := [2]int{x, y}
		state := grid[pos]
		switch state {
		case Clean:
			dir = (dir + 3) % 4 // left
			grid[pos] = Weakened
		case Weakened:
			// no turn
			grid[pos] = Infected
			infections++
		case Infected:
			dir = (dir + 1) % 4 // right
			grid[pos] = Flagged
		case Flagged:
			dir = (dir + 2) % 4 // reverse
			delete(grid, pos)   // clean
		}
		x += dirs[dir][0]
		y += dirs[dir][1]
	}
	return strconv.Itoa(infections), nil
}

func init() {
	solve.Register(Day22{})
}
