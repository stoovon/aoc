package solve2024

import (
	"image"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day8 struct {
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 8}
}

type Day8Parse struct {
	bounds map[image.Point]bool
	freq   map[rune][]image.Point
}

type Day8Result struct {
	part1 int
	part2 int
}

func (d Day8) parse(input string) Day8Parse {
	bounds := make(map[image.Point]bool)
	freq := make(map[rune][]image.Point)

	for y, line := range strings.Split(strings.TrimSpace(input), "\n") {
		for x, char := range line {
			point := image.Point{X: x, Y: y}
			bounds[point] = true
			if char != '.' {
				freq[char] = append(freq[char], point)
			}
		}
	}

	return Day8Parse{bounds: bounds, freq: freq}
}

func (d Day8) solve(parse Day8Parse) Day8Result {
	part1 := make(map[image.Point]struct{})
	part2 := make(map[image.Point]struct{})

	for _, antennas := range parse.freq {
		for _, a1 := range antennas {
			for _, a2 := range antennas {
				if a2 == a1 {
					continue
				}
				direction := a2.Sub(a1)
				next := a2.Add(direction)

				// Part 1: Check if the next point is within bounds
				if parse.bounds[next] {
					part1[next] = struct{}{}
				}

				// Part 2: Continue moving in the same direction while within bounds
				for next = a2; parse.bounds[next]; next = next.Add(direction) {
					part2[next] = struct{}{}
				}
			}
		}
	}

	return Day8Result{part1: len(part1), part2: len(part2)}
}

func (d Day8) Part1(input string) (string, error) {
	return strconv.Itoa(d.solve(d.parse(input)).part1), nil
}

func (d Day8) Part2(input string) (string, error) {
	return strconv.Itoa(d.solve(d.parse(input)).part2), nil
}

func init() {
	solve.Register(Day8{})
}
