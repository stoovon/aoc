package solve2017

import (
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 3}
}

func (d Day3) Part1(input string) (string, error) {
	n, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	if n == 1 {
		return "0", nil
	}
	// Find the layer (ring) where n is located
	layer := 0
	for {
		if (2*layer+1)*(2*layer+1) >= n {
			break
		}
		layer++
	}
	// The maximum value in this layer
	maxVal := (2*layer + 1) * (2*layer + 1)
	// The side length of the current layer
	sideLen := 2 * layer
	// Find the closest axis position
	stepsFromAxis := 0
	for i := 0; i < 4; i++ {
		axis := maxVal - layer - i*sideLen
		dist := maths.Abs(n - axis)
		if i == 0 || dist < stepsFromAxis {
			stepsFromAxis = dist
		}
	}
	distance := layer + stepsFromAxis
	return strconv.Itoa(distance), nil
}

func (d Day3) Part2(input string) (string, error) {
	n, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}

	type point struct{ x, y int }
	dirs := []point{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
	grid := map[point]int{{0, 0}: 1}

	// Spiral movement: right, up, left, down
	moveDirs := []point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	pos := point{0, 0}
	steps := 1
	for {
		for dirIdx := 0; dirIdx < 4; dirIdx++ {
			for i := 0; i < steps; i++ {
				pos.x += moveDirs[dirIdx].x
				pos.y += moveDirs[dirIdx].y
				sum := 0
				for _, d := range dirs {
					adj := point{pos.x + d.x, pos.y + d.y}
					sum += grid[adj]
				}
				grid[pos] = sum
				if sum > n {
					return strconv.Itoa(sum), nil
				}
			}
			// Increase steps after left and right moves
			if dirIdx == 1 || dirIdx == 3 {
				steps++
			}
		}
	}
}

func init() {
	solve.Register(Day3{})
}
