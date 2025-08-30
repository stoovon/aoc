package solve2021

import (
	"aoc/solve"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Day22 struct{}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 22}
}

func (d Day22) Part1(input string) (string, error) {
	// Define the initialization region bounds
	const minBound, maxBound = -50, 50

	steps, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	// Create a 3D grid for the initialization region
	grid := make(map[[3]int]bool)

	// Apply each step
	for _, step := range steps {
		for x := max(step.x1, minBound); x <= min(step.x2, maxBound); x++ {
			for y := max(step.y1, minBound); y <= min(step.y2, maxBound); y++ {
				for z := max(step.z1, minBound); z <= min(step.z2, maxBound); z++ {
					grid[[3]int{x, y, z}] = step.on
				}
			}
		}
	}

	// Count the number of cubes that are on
	onCount := 0
	for _, on := range grid {
		if on {
			onCount++
		}
	}

	return fmt.Sprintf("%d", onCount), nil
}

type rebootStep struct {
	on     bool
	x1, x2 int
	y1, y2 int
	z1, z2 int
}

func (d Day22) parseInput(input string) ([]rebootStep, error) {
	var steps []rebootStep
	lines := strings.Split(strings.TrimSpace(input), "\n")
	re := regexp.MustCompile(`(on|off) x=(-?\d+)..(-?\d+),y=(-?\d+)..(-?\d+),z=(-?\d+)..(-?\d+)`)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			return nil, errors.New("invalid input format")
		}
		on := matches[1] == "on"
		x1, _ := strconv.Atoi(matches[2])
		x2, _ := strconv.Atoi(matches[3])
		y1, _ := strconv.Atoi(matches[4])
		y2, _ := strconv.Atoi(matches[5])
		z1, _ := strconv.Atoi(matches[6])
		z2, _ := strconv.Atoi(matches[7])
		steps = append(steps, rebootStep{on, x1, x2, y1, y2, z1, z2})
	}
	return steps, nil
}

func (d Day22) Part2(input string) (string, error) {
	steps, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	// List of cuboids representing the "on" regions
	onCuboids := []rebootStep{}

	// Process each step
	for _, step := range steps {
		newCuboids := []rebootStep{}

		// Handle intersections with existing cuboids
		for _, existing := range onCuboids {
			if intersection, ok := intersectCuboids(existing, step); ok {
				// Toggle the state of the intersection based on the existing cuboid
				intersection.on = !existing.on
				newCuboids = append(newCuboids, intersection)
			}
		}

		// Add the new cuboid if it's an "on" operation
		if step.on {
			newCuboids = append(newCuboids, step)
		}

		onCuboids = append(onCuboids, newCuboids...)
	}

	// Calculate the total volume of "on" cuboids
	totalVolume := 0
	for _, cuboid := range onCuboids {
		volume := calculateVolume(cuboid)
		if cuboid.on {
			totalVolume += volume
		} else {
			totalVolume -= volume
		}
	}

	return fmt.Sprintf("%d", totalVolume), nil
}

func intersectCuboids(a, b rebootStep) (rebootStep, bool) {
	x1 := max(a.x1, b.x1)
	x2 := min(a.x2, b.x2)
	y1 := max(a.y1, b.y1)
	y2 := min(a.y2, b.y2)
	z1 := max(a.z1, b.z1)
	z2 := min(a.z2, b.z2)

	if x1 <= x2 && y1 <= y2 && z1 <= z2 {
		return rebootStep{on: false, x1: x1, x2: x2, y1: y1, y2: y2, z1: z1, z2: z2}, true
	}
	return rebootStep{}, false
}

func calculateVolume(cuboid rebootStep) int {
	return (cuboid.x2 - cuboid.x1 + 1) * (cuboid.y2 - cuboid.y1 + 1) * (cuboid.z2 - cuboid.z1 + 1)
}

func init() {
	solve.Register(Day22{})
}
