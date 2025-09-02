package solve2022

import (
	"aoc/solve"
	"errors"
	"strconv"
	"strings"
)

type Day17 struct{}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 17}
}

func (d Day17) Part1(input string) (string, error) {
	return d.simulate(input, 2022)
}

func (d Day17) Part2(input string) (string, error) {
	return d.simulate(input, 1000000000000)
}

func (d Day17) simulate(input string, nPentominos int) (string, error) {
	jets := strings.TrimSpace(input)
	if len(jets) == 0 {
		return "", errors.New("empty input")
	}

	type Point struct {
		X, Y int
	}

	pentominos := [][]Point{
		{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
		{{1, 2}, {0, 1}, {1, 1}, {2, 1}, {1, 0}},
		{{2, 2}, {2, 1}, {0, 0}, {1, 0}, {2, 0}},
		{{0, 3}, {0, 2}, {0, 1}, {0, 0}},
		{{0, 1}, {1, 1}, {0, 0}, {1, 0}},
	}

	grid := map[Point]struct{}{}

	// Helper function to move pentominos
	movePentomino := func(pentomino []Point, delta Point) bool {
		npentomino := make([]Point, len(pentomino))
		for i, p := range pentomino {
			p.X += delta.X
			p.Y += delta.Y
			if _, ok := grid[p]; ok || p.X < 0 || p.X >= 7 || p.Y < 0 {
				return false
			}
			npentomino[i] = p
		}
		copy(pentomino, npentomino)
		return true
	}

	initializePentomino := func(basePentomino []Point, height int) []Point {
		newPentomino := make([]Point, len(basePentomino))
		for i, p := range basePentomino {
			newPentomino[i] = Point{p.X + 2, p.Y + height + 3}
		}
		return newPentomino
	}

	cache := map[[2]int][]int{}
	height, jet := 0, 0
	for i := 0; i <= nPentominos; i++ {
		if i == 2022 && nPentominos == 2022 {
			return strconv.Itoa(height), nil
		}

		k := [2]int{i % len(pentominos), jet}
		if c, ok := cache[k]; ok {
			if n, d := 1000000000000-i, i-c[0]; n%d == 0 {
				return strconv.Itoa(height + n/d*(height-c[1])), nil
			}
		}
		cache[k] = []int{i, height}

		pentomino := initializePentomino(pentominos[i%len(pentominos)], height)

		for {
			movePentomino(pentomino, Point{int(jets[jet]) - int('='), 0})
			jet = (jet + 1) % len(jets)

			if !movePentomino(pentomino, Point{0, -1}) {
				for _, p := range pentomino {
					grid[p] = struct{}{}
					if p.Y+1 > height {
						height = p.Y + 1
					}
				}
				break
			}
		}
	}

	return "", errors.New("simulation did not complete")
}

func init() {
	solve.Register(Day17{})
}
