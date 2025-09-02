package solve2022

import (
	"aoc/solve"
	"image"
	"strconv"
	"strings"
)

type BlizzardState struct {
	P image.Point
	T int
}

type Day24 struct{}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 24}
}

func (d Day24) Part1(input string) (string, error) {
	valley, blizzards := d.parseInput(input)

	// Define start and end points
	start, end := blizzards.Min.Sub(image.Point{0, 1}), blizzards.Max.Sub(image.Point{1, 0})
	result := d.bfs(valley, blizzards, start, end, 0)
	return strconv.Itoa(result), nil
}

func (d Day24) Part2(input string) (string, error) {
	valley, blizzards := d.parseInput(input)

	// There, and back again, and there once more
	start, end := blizzards.Min.Sub(image.Point{0, 1}), blizzards.Max.Sub(image.Point{1, 0})
	result := d.bfs(valley, blizzards, start, end, 0)
	result = d.bfs(valley, blizzards, end, start, result)
	result = d.bfs(valley, blizzards, start, end, result)
	return strconv.Itoa(result), nil
}

func (d Day24) parseInput(input string) (map[image.Point]rune, image.Rectangle) {
	vall := map[image.Point]rune{}
	for y, s := range strings.Fields(input) {
		for x, r := range s {
			vall[image.Point{x, y}] = r
		}
	}

	var bliz image.Rectangle
	for p := range vall {
		bliz = bliz.Union(image.Rectangle{p, p.Add(image.Point{1, 1})})
	}
	bliz.Min, bliz.Max = bliz.Min.Add(image.Point{1, 1}), bliz.Max.Sub(image.Point{1, 1})

	return vall, bliz
}

func (d Day24) bfs(vall map[image.Point]rune, bliz image.Rectangle, start, end image.Point, time int) int {
	delta := map[rune]image.Point{
		'#': {0, 0}, '^': {0, -1}, '>': {1, 0}, 'v': {0, 1}, '<': {-1, 0},
	}

	queue := []BlizzardState{{start, time}}
	seen := map[BlizzardState]struct{}{queue[0]: {}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

	loop:
		for _, d := range delta {
			next := BlizzardState{cur.P.Add(d), cur.T + 1}
			if next.P == end {
				return next.T
			}

			if _, ok := seen[next]; ok {
				continue
			}
			if r, ok := vall[next.P]; !ok || r == '#' {
				continue
			}

			if next.P.In(bliz) {
				for r, d := range delta {
					if vall[next.P.Sub(d.Mul(next.T)).Mod(bliz)] == r {
						continue loop
					}
				}
			}

			seen[next] = struct{}{}
			queue = append(queue, next)
		}
	}
	return -1
}

func init() {
	solve.Register(Day24{})
}
