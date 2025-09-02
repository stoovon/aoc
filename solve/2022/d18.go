package solve2022

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day18 struct{}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 18}
}

func (d Day18) Part1(input string) (string, error) {
	// Parse the input into a set of coordinates
	cubes := d.parseInput(input)

	totalSurfaceArea := 0

	// Directions for neighbors (x, y, z)
	directions := [][3]int{
		{1, 0, 0}, {-1, 0, 0},
		{0, 1, 0}, {0, -1, 0},
		{0, 0, 1}, {0, 0, -1},
	}

	// Calculate the surface area
	for cube := range cubes {
		exposedSides := 6
		for _, dir := range directions {
			neighbor := [3]int{cube[0] + dir[0], cube[1] + dir[1], cube[2] + dir[2]}
			if cubes[neighbor] {
				exposedSides--
			}
		}
		totalSurfaceArea += exposedSides
	}

	return strconv.Itoa(totalSurfaceArea), nil
}

func (d Day18) Part2(input string) (string, error) {
	// Parse the input into a set of coordinates
	cubes := d.parseInput(input)

	// Directions for neighbors (x, y, z)
	directions := [][3]int{
		{1, 0, 0}, {-1, 0, 0},
		{0, 1, 0}, {0, -1, 0},
		{0, 0, 1}, {0, 0, -1},
	}

	// Find the bounding box of the droplet
	min, max := [3]int{1<<31 - 1, 1<<31 - 1, 1<<31 - 1}, [3]int{-1 << 31, -1 << 31, -1 << 31}
	for cube := range cubes {
		for i := 0; i < 3; i++ {
			if cube[i] < min[i] {
				min[i] = cube[i]
			}
			if cube[i] > max[i] {
				max[i] = cube[i]
			}
		}
	}

	// Expand the bounding box slightly to ensure we start outside the droplet
	for i := 0; i < 3; i++ {
		min[i]--
		max[i]++
	}

	// Perform a flood fill to find all reachable air cubes
	reachable := make(map[[3]int]bool)
	queue := [][3]int{min}
	reachable[min] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, dir := range directions {
			neighbor := [3]int{current[0] + dir[0], current[1] + dir[1], current[2] + dir[2]}

			// Skip if out of bounds, already visited, or part of the droplet
			if neighbor[0] < min[0] || neighbor[1] < min[1] || neighbor[2] < min[2] ||
				neighbor[0] > max[0] || neighbor[1] > max[1] || neighbor[2] > max[2] ||
				reachable[neighbor] || cubes[neighbor] {
				continue
			}

			reachable[neighbor] = true
			queue = append(queue, neighbor)
		}
	}

	// Calculate the exterior surface area
	exteriorSurfaceArea := 0
	for cube := range cubes {
		for _, dir := range directions {
			neighbor := [3]int{cube[0] + dir[0], cube[1] + dir[1], cube[2] + dir[2]}
			if reachable[neighbor] {
				exteriorSurfaceArea++
			}
		}
	}

	return strconv.Itoa(exteriorSurfaceArea), nil
}

func init() {
	solve.Register(Day18{})
}

func (d Day18) parseInput(input string) map[[3]int]bool {
	cubes := make(map[[3]int]bool)
	lines := strings.SplitSeq(strings.TrimSpace(input), "\n")
	for line := range lines {
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		cubes[[3]int{x, y, z}] = true
	}
	return cubes
}
