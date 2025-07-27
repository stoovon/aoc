package solve2020

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day17 struct{}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 17}
}

func (d Day17) Part1(input string) (string, error) {
	active := parseActive3D(strings.TrimSpace(input))
	return strconv.Itoa(simulateCycles3D(active, 6)), nil
}

func (d Day17) Part2(input string) (string, error) {
	active := parseActive4D(strings.TrimSpace(input))
	return strconv.Itoa(simulateCycles4D(active, 6)), nil
}

type coord3 struct{ x, y, z int }
type coord4 struct{ x, y, z, w int }

func parseActive3D(input string) map[coord3]struct{} {
	active := make(map[coord3]struct{})
	for y, line := range strings.Split(input, "\n") {
		for x, ch := range line {
			if ch == '#' {
				active[coord3{x, y, 0}] = struct{}{}
			}
		}
	}
	return active
}

func parseActive4D(input string) map[coord4]struct{} {
	active := make(map[coord4]struct{})
	for y, line := range strings.Split(input, "\n") {
		for x, ch := range line {
			if ch == '#' {
				active[coord4{x, y, 0, 0}] = struct{}{}
			}
		}
	}
	return active
}

func simulateCycles3D(active map[coord3]struct{}, cycles int) int {
	neighbors := [][3]int{}
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				if dx != 0 || dy != 0 || dz != 0 {
					neighbors = append(neighbors, [3]int{dx, dy, dz})
				}
			}
		}
	}
	for c := 0; c < cycles; c++ {
		neighborCount := make(map[coord3]int)
		for cube := range active {
			for _, d := range neighbors {
				n := coord3{cube.x + d[0], cube.y + d[1], cube.z + d[2]}
				neighborCount[n]++
			}
		}
		next := make(map[coord3]struct{})
		for cube, count := range neighborCount {
			_, isActive := active[cube]
			if (isActive && (count == 2 || count == 3)) || (!isActive && count == 3) {
				next[cube] = struct{}{}
			}
		}
		active = next
	}
	return len(active)
}

func simulateCycles4D(active map[coord4]struct{}, cycles int) int {
	neighbors := [][4]int{}
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				for dw := -1; dw <= 1; dw++ {
					if dx != 0 || dy != 0 || dz != 0 || dw != 0 {
						neighbors = append(neighbors, [4]int{dx, dy, dz, dw})
					}
				}
			}
		}
	}
	for range cycles {
		neighborCount := make(map[coord4]int)
		for cube := range active {
			for _, d := range neighbors {
				n := coord4{cube.x + d[0], cube.y + d[1], cube.z + d[2], cube.w + d[3]}
				neighborCount[n]++
			}
		}
		next := make(map[coord4]struct{})
		for cube, count := range neighborCount {
			_, isActive := active[cube]
			if (isActive && (count == 2 || count == 3)) || (!isActive && count == 3) {
				next[cube] = struct{}{}
			}
		}
		active = next
	}
	return len(active)
}

func init() {
	solve.Register(Day17{})
}
