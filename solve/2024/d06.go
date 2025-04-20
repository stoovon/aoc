package solve2024

import (
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day6 struct {
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 6}
}

type GridSpec struct {
	grid  map[image.Point]rune
	start image.Point
}

func (d Day6) parseGrid(input string) GridSpec {
	grid := make(map[image.Point]rune)
	var start image.Point

	for y, line := range strings.Split(strings.TrimSpace(input), "\n") {
		for x, char := range line {
			point := image.Point{X: x, Y: y}
			grid[point] = char
			if char == '^' {
				start = point
			}
		}
	}

	return GridSpec{grid: grid, start: start}
}

func (d Day6) visitPoints(origin image.Point, spec GridSpec) int {
	directions := grids.URDL()
	position, direction := spec.start, 0
	visited := make(map[image.Point]int)

	for {
		if _, exists := spec.grid[position]; !exists {
			return len(visited)
		}
		if visited[position]&(1<<direction) != 0 {
			return -1
		}

		visited[position] |= 1 << direction
		next := position.Add(directions[direction])

		if spec.grid[next] == '#' || next == origin {
			direction = (direction + 1) % len(directions)
		} else {
			position = next
		}
	}
}

func (d Day6) Part1(input string) (string, error) {
	spec := d.parseGrid(input)
	result := d.visitPoints(image.Point{X: -1, Y: -1}, spec)
	return strconv.Itoa(result), nil
}

func (d Day6) Part2(input string) (string, error) {
	spec := d.parseGrid(input)
	blockedCount := 0

	for point := range spec.grid {
		if d.visitPoints(point, spec) == -1 {
			blockedCount++
		}
	}

	return strconv.Itoa(blockedCount), nil
}

func init() {
	solve.Register(Day6{})
}
