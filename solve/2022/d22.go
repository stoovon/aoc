package solve2022

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

var pathTokenRE = regexp.MustCompile(`(\d+|[LR])`)

type Day22 struct{}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 22}
}

type step struct {
	steps int
	turn  byte // 'L', 'R', or 0
}

var dirDeltas = [4][2]int{
	{0, 1},  // right
	{1, 0},  // down
	{0, -1}, // left
	{-1, 0}, // up
}

func parseInput22(input string) ([][]byte, []step) {
	parts := strings.SplitN(input, "\n\n", 2)
	lines := strings.Split(parts[0], "\n")
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	grid := make([][]byte, len(lines))
	for i, l := range lines {
		row := make([]byte, maxLen)
		for j := range row {
			if j < len(l) {
				row[j] = l[j]
			} else {
				row[j] = ' '
			}
		}
		grid[i] = row
	}
	path := parts[1]
	tokens := pathTokenRE.FindAllString(path, -1)
	var steps []step
	for _, tok := range tokens {
		if tok == "L" || tok == "R" {
			steps = append(steps, step{0, tok[0]})
		} else {
			n, _ := strconv.Atoi(tok)
			steps = append(steps, step{n, 0})
		}
	}
	return grid, steps
}

func findStart(grid [][]byte) (int, int) {
	for j, c := range grid[0] {
		if c == '.' {
			return 0, j
		}
	}
	return -1, -1
}

func wrapFlat(grid [][]byte, r, c, dir int) (int, int) {
	rows, cols := len(grid), len(grid[0])
	nr, nc := r, c
	for {
		nr = (nr + dirDeltas[dir][0] + rows) % rows
		nc = (nc + dirDeltas[dir][1] + cols) % cols
		if grid[nr][nc] != ' ' {
			return nr, nc
		}
	}
}

func (d Day22) Part1(input string) (string, error) {
	grid, steps := parseInput22(input)
	r, c := findStart(grid)
	dir := 0 // right
	for _, s := range steps {
		if s.steps > 0 {
			for i := 0; i < s.steps; i++ {
				nr, nc := r+dirDeltas[dir][0], c+dirDeltas[dir][1]
				if nr < 0 || nr >= len(grid) || nc < 0 || nc >= len(grid[0]) || grid[nr][nc] == ' ' {
					nr, nc = wrapFlat(grid, r, c, dir)
				}
				if grid[nr][nc] == '#' {
					break
				}
				r, c = nr, nc
			}
		} else if s.turn != 0 {
			if s.turn == 'R' {
				dir = (dir + 1) % 4
			} else {
				dir = (dir + 3) % 4
			}
		}
	}
	// Rows and columns are 1-based
	return fmt.Sprintf("%d", 1000*(r+1)+4*(c+1)+dir), nil
}

func wrapCube(grid [][]byte, r, c, dir int) (int, int, int) {
	const face = 50
	faceX, faceY := c/face, r/face

	switch {
	// Face 1 (faceX=1, faceY=0)
	case dir == 3 && faceX == 1 && faceY == 0: // NORTH from face 1 -> face 6, EAST
		return 150 + (c - 50), 0, 0
	case dir == 2 && faceX == 1 && faceY == 0: // WEST from face 1 -> face 4, EAST
		return 149 - r, 0, 0

	// Face 2 (faceX=2, faceY=0)
	case dir == 3 && faceX == 2 && faceY == 0: // NORTH from face 2 -> face 6, NORTH
		return 199, c - 100, 3
	case dir == 0 && faceX == 2 && faceY == 0: // EAST from face 2 -> face 5, WEST
		return 149 - r, 99, 2
	case dir == 1 && faceX == 2 && faceY == 0: // SOUTH from face 2 -> face 3, WEST
		return 50 + (c - 100), 99, 2

	// Face 3 (faceX=1, faceY=1)
	case dir == 0 && faceX == 1 && faceY == 1: // EAST from face 3 -> face 2, NORTH
		return 49, 100 + (r - 50), 3
	case dir == 2 && faceX == 1 && faceY == 1: // WEST from face 3 -> face 4, SOUTH
		return 100, r - 50, 1

	// Face 4 (faceX=0, faceY=2)
	case dir == 3 && faceX == 0 && faceY == 2: // NORTH from face 4 -> face 3, EAST
		return 50 + (c), 50, 0
	case dir == 2 && faceX == 0 && faceY == 2: // WEST from face 4 -> face 1, EAST
		return 149 - r, 50, 0
	case dir == 1 && faceX == 0 && faceY == 2: // SOUTH from face 4 -> face 6, EAST
		return 150 + (c - 0), 49, 0

	// Face 5 (faceX=1, faceY=2)
	case dir == 0 && faceX == 1 && faceY == 2: // EAST from face 5 -> face 2, WEST
		return 149 - r, 149, 2
	case dir == 1 && faceX == 1 && faceY == 2: // SOUTH from face 5 -> face 6, WEST
		return 150 + (c - 50), 49, 2

	// Face 6 (faceX=0, faceY=3)
	case dir == 0 && faceX == 0 && faceY == 3: // EAST from face 6 -> face 5, NORTH
		return 149, 50 + (r - 150), 3
	case dir == 1 && faceX == 0 && faceY == 3: // SOUTH from face 6 -> face 2, SOUTH
		return 0, 100 + c, 1
	case dir == 2 && faceX == 0 && faceY == 3: // WEST from face 6 -> face 1, SOUTH
		return 0, 50 + (r - 150), 1
	}

	panic(fmt.Sprintf("Invalid cube wrap: r=%d c=%d dir=%d", r, c, dir))
}

func (d Day22) Part2(input string) (string, error) {
	grid, steps := parseInput22(input)
	r, c := findStart(grid)
	dir := 0 // right
	for _, s := range steps {
		if s.steps > 0 {
			for i := 0; i < s.steps; i++ {
				nr, nc, ndir := r+dirDeltas[dir][0], c+dirDeltas[dir][1], dir
				if nr < 0 || nr >= len(grid) || nc < 0 || nc >= len(grid[0]) || grid[nr][nc] == ' ' {
					nr, nc, ndir = wrapCube(grid, r, c, dir)
				}
				if grid[nr][nc] == '#' {
					break
				}
				r, c, dir = nr, nc, ndir
			}
		} else if s.turn != 0 {
			if s.turn == 'R' {
				dir = (dir + 1) % 4
			} else {
				dir = (dir + 3) % 4
			}
		}
	}
	return fmt.Sprintf("%d", 1000*(r+1)+4*(c+1)+dir), nil
}

func init() {
	solve.Register(Day22{})
}
