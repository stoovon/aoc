package solve2021

import (
	"aoc/solve"
	"aoc/utils/grids"
	"fmt"
	"strings"
)

type Day13 struct{}

// foldDots applies a fold instruction to a set of dots
func foldDots(dots map[[2]int]struct{}, axis string, value int) map[[2]int]struct{} {
	newDots := make(map[[2]int]struct{})
	for dot := range dots {
		x, y := dot[0], dot[1]
		if axis == "x" {
			if x > value {
				x = value - (x - value)
			}
		} else if axis == "y" {
			if y > value {
				y = value - (y - value)
			}
		}
		newDots[[2]int{x, y}] = struct{}{}
	}
	return newDots
}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 13}
}

func (d Day13) Part1(input string) (string, error) {
	// Split input into dots and folds sections
	sections := strings.SplitN(input, "\n\n", 2)
	if len(sections) < 2 {
		return "0", nil
	}
	dotLines := strings.Split(sections[0], "\n")
	foldLines := strings.Split(sections[1], "\n")

	dots := make(map[[2]int]struct{})
	for _, line := range dotLines {
		if len(line) == 0 {
			continue
		}
		var x, y int
		_, err := fmt.Sscanf(line, "%d,%d", &x, &y)
		if err == nil {
			dots[[2]int{x, y}] = struct{}{}
		}
	}

	folds := make([]string, 0)
	for _, line := range foldLines {
		if len(line) == 0 {
			continue
		}
		folds = append(folds, line)
	}

	// Only apply the first fold
	if len(folds) == 0 {
		return "0", nil
	}
	fold := folds[0]
	var axis string
	var value int
	_, err := fmt.Sscanf(fold, "fold along %1s=%d", &axis, &value)
	if err != nil {
		return "", err
	}

	newDots := foldDots(dots, axis, value)
	return fmt.Sprintf("%d", len(newDots)), nil
}

func (d Day13) Part2(input string) (string, error) {
	// Split input into dots and folds sections
	sections := strings.SplitN(input, "\n\n", 2)
	if len(sections) < 2 {
		return "", nil
	}
	dotLines := strings.Split(sections[0], "\n")
	foldLines := strings.Split(sections[1], "\n")

	dots := make(map[[2]int]struct{})
	for _, line := range dotLines {
		if len(line) == 0 {
			continue
		}
		var x, y int
		_, err := fmt.Sscanf(line, "%d,%d", &x, &y)
		if err == nil {
			dots[[2]int{x, y}] = struct{}{}
		}
	}

	folds := make([][2]interface{}, 0)
	for _, line := range foldLines {
		if len(line) == 0 {
			continue
		}
		var axis string
		var value int
		_, err := fmt.Sscanf(line, "fold along %1s=%d", &axis, &value)
		if err == nil {
			folds = append(folds, [2]interface{}{axis, value})
		}
	}

	// Apply all folds
	for _, f := range folds {
		axis := f[0].(string)
		value := f[1].(int)
		dots = foldDots(dots, axis, value)
	}

	// Find grid size
	maxX, maxY := 0, 0
	for dot := range dots {
		if dot[0] > maxX {
			maxX = dot[0]
		}
		if dot[1] > maxY {
			maxY = dot[1]
		}
	}

	// Build grid for OCR
	grid := make([][]int, maxY+1)
	for y := 0; y <= maxY; y++ {
		grid[y] = make([]int, maxX+1)
		for x := 0; x <= maxX; x++ {
			if _, ok := dots[[2]int{x, y}]; ok {
				grid[y][x] = 1
			} else {
				grid[y][x] = 0
			}
		}
	}
	code := grids.OCR(grid)
	return code, nil
}

func init() {
	solve.Register(Day13{})
}
