package solve2019

import (
	"fmt"

	"aoc/solve"
)

type Day19 struct {
}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 19}
}

func (d Day19) Part1(input string) (string, error) {
	count := 0
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			out := int64(-1)
			p := parseIntcode(input) // create a fresh instance
			p.runCore([]int64{int64(x), int64(y)}, func(v int64) bool {
				out = v
				return false
			})
			if out == 1 {
				count++
			}
		}
	}
	return fmt.Sprintf("%d", count), nil
}

func (d Day19) Part2(input string) (string, error) {
	// Helper to check if (x, y) is in the beam
	inBeam := func(x, y int) bool {
		out := int64(-1)
		p := parseIntcode(input)
		p.runCore([]int64{int64(x), int64(y)}, func(v int64) bool {
			out = v
			return false
		})
		return out == 1
	}

	y := 100
	x := 0
	for {
		// Move x right until beam starts at this y
		for !inBeam(x, y) {
			x++
		}
		// Check if top-right of 100x100 square is also in beam
		if inBeam(x+99, y-99) {
			// Found: top-left is (x, y-99)
			result := x*10000 + (y - 99)
			return fmt.Sprintf("%d", result), nil
		}
		y++
	}
}

func init() {
	solve.Register(Day19{})
}
