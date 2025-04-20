package solve2024

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day18 struct {
}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 18}
}

func (d Day18) solve(input string, part1 bool) string {
	size := 70

	var bytes []image.Point
	for _, s := range strings.Fields(input) {
		var x, y int
		_, err := fmt.Sscanf(s, "%d,%d", &x, &y)
		if err != nil {
			return ""
		}
		bytes = append(bytes, image.Point{X: x, Y: y})
	}

	grid := map[image.Point]bool{}
	for y := range size + 1 {
		for x := range size + 1 {
			grid[image.Point{X: x, Y: y}] = true
		}
	}

loop:
	for b := range bytes {
		grid[bytes[b]] = false

		// Breadth-First Search (BFS) initialization
		queue, dist := []image.Point{{0, 0}}, map[image.Point]int{{0, 0}: 0}
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]

			if p == (image.Point{X: size, Y: size}) {
				if b == 1024 {
					if part1 {
						return strconv.Itoa(dist[p])
					}
				}
				continue loop
			}

			// BFS: Explore neighbors
			for _, d := range grids.URDL() {
				n := p.Add(d)
				if _, ok := dist[n]; !ok && grid[n] {
					queue, dist[n] = append(queue, n), dist[p]+1
				}
			}
		}

		if !part1 {
			return fmt.Sprintf("%d,%d", bytes[b].X, bytes[b].Y)
		}
		break
	}

	return "nope"
}

func (d Day18) Part1(input string) (string, error) {
	return d.solve(input, true), nil
}

func (d Day18) Part2(input string) (string, error) {
	return d.solve(input, false), nil
}

func init() {
	solve.Register(Day18{})
}
