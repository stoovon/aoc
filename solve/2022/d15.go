package solve2022

import (
	"aoc/solve"
	"aoc/utils/maths"
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Day15 struct{}

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 15}
}

type P struct{ X, Y int }

type Sensor struct {
	Pos    P
	Beacon P

	// Can be recomputed as needed, but let's store it for convenience
	Dist int
}

func (d Day15) parse(input string) []Sensor {
	var sensors []Sensor
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()
		var sx, sy, bx, by int
		fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sx, &sy, &bx, &by)

		sensorPos := P{X: sx, Y: sy}
		beaconPos := P{X: bx, Y: by}
		sensors = append(sensors, Sensor{
			Pos:    sensorPos,
			Beacon: beaconPos,
			Dist:   manhattanDistance(sensorPos, beaconPos),
		})
	}

	return sensors
}

func manhattanDistance(a, b P) int {
	return maths.Abs(a.X-b.X) + maths.Abs(a.Y-b.Y)
}

func findBoundaries(sensors []Sensor) (int, int) {
	minX, maxX := 1<<63-1, 0
	for _, sensor := range sensors {
		leftBoundary := sensor.Pos.X - sensor.Dist
		rightBoundary := sensor.Pos.X + sensor.Dist
		if leftBoundary < minX {
			minX = leftBoundary
		}
		if rightBoundary > maxX {
			maxX = rightBoundary
		}
	}
	return minX, maxX
}

func (d Day15) Part1(input string) (string, error) {
	sensors := d.parse(input)
	minX, maxX := findBoundaries(sensors)

	y := 2000000

	count := 0
Point:
	for x := minX; x < maxX; x++ {
		p := P{x, y}
		for _, s := range sensors {
			if manhattanDistance(p, s.Pos) <= s.Dist && p != s.Beacon {
				count++
				continue Point
			}
		}
	}
	return strconv.Itoa(count), nil
}

// calculateSkip determines how many X values can be skipped when a point is inside a sensor's range.
// It calculates the horizontal distance to the closest boundary of the sensor's diamond-shaped range.
func calculateSkip(p, s P, dist int) int {
	return dist - maths.Abs(s.Y-p.Y) + s.X - p.X
}

func (d Day15) Part2(input string) (string, error) {
	sensors := d.parse(input)

	for y := 0; y <= 4000000; y++ {
	Point:
		for x := 0; x <= 4000000; x++ {
			p := P{x, y}

			for _, s := range sensors {
				if manhattanDistance(p, s.Pos) < s.Dist {
					x += calculateSkip(p, s.Pos, s.Dist)
					continue Point
				}
			}

			// Sensor just out of range: found it!
			return strconv.Itoa(4000000*x + y), nil
		}
	}

	return "", errors.New("not found")
}

func init() {
	solve.Register(Day15{})
}
