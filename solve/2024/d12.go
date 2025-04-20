package solve2024

import (
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day12 struct {
}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 12}
}

func (d Day12) parseGrid(input string) map[image.Point]rune {
	grid := map[image.Point]rune{}
	for y, s := range strings.Fields(input) {
		for x, r := range s {
			grid[image.Point{X: x, Y: y}] = r
		}
	}

	return grid
}

func (d Day12) solve(grid map[image.Point]rune) (int, int) {
	seen := map[image.Point]bool{}
	part1, part2 := 0, 0
	for p := range grid {
		if seen[p] {
			continue
		}
		seen[p] = true

		area := 1
		perimeter, sides := 0, 0
		queue := []image.Point{p}
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]

			for _, d := range grids.URDL() {
				if n := p.Add(d); grid[n] != grid[p] {
					perimeter++
					r := p.Add(image.Point{X: -d.Y, Y: d.X})
					if grid[r] != grid[p] || grid[r.Add(d)] == grid[p] {
						sides++
					}
				} else if !seen[n] {
					seen[n] = true
					queue = append(queue, n)
					area++
				}
			}
		}
		part1 += area * perimeter
		part2 += area * sides
	}

	return part1, part2
}

func (d Day12) Part1(input string) (string, error) {
	part1, _ := d.solve(d.parseGrid(input))

	return strconv.Itoa(part1), nil
}

func (d Day12) Part2(input string) (string, error) {
	_, part2 := d.solve(d.parseGrid(input))

	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day12{})
}
