package solve2019

import (
	"errors"
	"image"
	"strconv"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day11 struct {
}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 11}
}

func (d Day11) Part1(input string) (string, error) {
	code := parseIntcode(input)
	// Directions: up, right, down, left
	dirs := []image.Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	dir := 0 // start facing up
	pos := image.Pt(0, 0)
	panels := make(map[image.Point]int)
	painted := make(map[image.Point]bool)

	for {
		// Input: current panel color (default black/0)
		color := panels[pos]
		outputs, halted := code.Step([]int64{int64(color)}, 2)
		if halted {
			break
		}
		if len(outputs) < 2 {
			return "", errors.New("Intcode did not output two values")
		}
		// Paint
		panels[pos] = int(outputs[0])
		painted[pos] = true
		// Turn
		if outputs[1] == 0 {
			dir = (dir + 3) % 4 // left
		} else {
			dir = (dir + 1) % 4 // right
		}
		// Move
		pos = pos.Add(dirs[dir])
	}

	return strconv.Itoa(len(painted)), nil
}

func (d Day11) Part2(input string) (string, error) {
	code := parseIntcode(input)
	dirs := []image.Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	dir := 0
	pos := image.Pt(0, 0)
	panels := make(map[image.Point]int)
	panels[pos] = 1 // Start on a white panel

	for {
		color := panels[pos]
		outputs, halted := code.Step([]int64{int64(color)}, 2)
		if halted {
			break
		}
		if len(outputs) < 2 {
			return "", errors.New("Intcode did not output two values")
		}
		panels[pos] = int(outputs[0])
		if outputs[1] == 0 {
			dir = (dir + 3) % 4
		} else {
			dir = (dir + 1) % 4
		}
		pos = pos.Add(dirs[dir])
	}

	// Find bounds
	var minX, minY, maxX, maxY int
	first := true
	for pt, color := range panels {
		if color == 0 {
			continue
		}
		if first {
			minX, maxX, minY, maxY = pt.X, pt.X, pt.Y, pt.Y
			first = false
		} else {
			if pt.X < minX {
				minX = pt.X
			}
			if pt.X > maxX {
				maxX = pt.X
			}
			if pt.Y < minY {
				minY = pt.Y
			}
			if pt.Y > maxY {
				maxY = pt.Y
			}
		}
	}

	width := maxX - minX + 1
	height := maxY - minY + 1
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			if panels[image.Pt(x+minX, y+minY)] == 1 {
				grid[y][x] = 1
			} else {
				grid[y][x] = 0
			}
		}
	}

	return grids.OCR(grid), nil
}

func init() {
	solve.Register(Day11{})
}
