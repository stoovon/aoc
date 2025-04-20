package solve2024

import (
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
	"aoc/utils/maths"
)

type Day20 struct {
}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 20}
}

func (d Day20) solve(input string) (int, int) {
	grid, start := map[image.Point]rune{}, image.Point{}
	for y, s := range strings.Fields(string(input)) {
		for x, r := range s {
			if r == 'S' {
				start = image.Point{X: x, Y: y}
			}
			grid[image.Point{X: x, Y: y}] = r
		}
	}

	queue, dist := []image.Point{start}, map[image.Point]int{start: 0}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, d := range grids.URDL() {
			n := p.Add(d)
			if _, ok := dist[n]; !ok && grid[n] != '#' {
				queue, dist[n] = append(queue, n), dist[p]+1
			}
		}
	}

	part1, part2 := 0, 0
	for p1 := range dist {
		for p2 := range dist {
			d := maths.Abs(p2.X-p1.X) + maths.Abs(p2.Y-p1.Y)
			if d <= 20 && dist[p2] >= dist[p1]+d+100 {
				if d <= 2 {
					part1++
				}
				part2++
			}
		}
	}

	return part1, part2
}

func (d Day20) Part1(input string) (string, error) {
	p1, _ := d.solve(input)

	return strconv.Itoa(p1), nil
}

func (d Day20) Part2(input string) (string, error) {
	_, p2 := d.solve(input)

	return strconv.Itoa(p2), nil
}

func init() {
	solve.Register(Day20{})
}
