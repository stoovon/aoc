package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/acstrings"
)

type Day6 struct {
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 6}
}

func (d Day6) solve(input string) (part1, part2 int) {
	A := make([][]bool, 1000)
	B := make([][]int, 1000)
	for i := 0; i < 1000; i++ {
		A[i] = make([]bool, 1000)
		B[i] = make([]int, 1000)
	}

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Fields(line)
		action := ""
		if parts[0] == "toggle" {
			action = "toggle"
			parts = parts[1:]
		} else {
			action = parts[1]
			parts = parts[2:]
		}

		// Parse coordinates
		p1 := strings.Split(parts[0], ",")
		p2 := strings.Split(parts[2], ",")
		x1, y1 := acstrings.MustInt(p1[0]), acstrings.MustInt(p1[1])
		x2, y2 := acstrings.MustInt(p2[0]), acstrings.MustInt(p2[1])

		// Apply the action
		for x := x1; x <= x2; x++ {
			for y := y1; y <= y2; y++ {
				switch action {
				case "toggle":
					A[x][y] = !A[x][y]
					B[x][y] += 2
				case "on":
					A[x][y] = true
					B[x][y]++
				case "off":
					A[x][y] = false
					B[x][y]--
					if B[x][y] < 0 {
						B[x][y] = 0
					}
				}
			}
		}
	}

	// Calculate results
	for x := 0; x < 1000; x++ {
		for y := 0; y < 1000; y++ {
			if A[x][y] {
				part1++
			}
			part2 += B[x][y]
		}
	}

	return part1, part2
}

func (d Day6) Part1(input string) (string, error) {
	part1, _ := d.solve(input)
	return strconv.Itoa(part1), nil
}

func (d Day6) Part2(input string) (string, error) {
	_, part2 := d.solve(input)
	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day6{})
}
