package solve2018

import (
	"errors"
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day6 struct {
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 6}
}

func parsePoints(input string) []image.Point {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var pts []image.Point
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ",")
		x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		pts = append(pts, image.Pt(x, y))
	}
	return pts
}

func manhattan(a, b image.Point) int {
	return maths.Abs(a.X-b.X) + maths.Abs(a.Y-b.Y)
}

func (d Day6) Part1(input string) (string, error) {
	pts := parsePoints(input)
	if len(pts) == 0 {
		return "", errors.New("no input")
	}

	// Find bounding box
	minX, maxX := pts[0].X, pts[0].X
	minY, maxY := pts[0].Y, pts[0].Y
	for _, p := range pts {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	area := make([]int, len(pts))
	infinite := make([]bool, len(pts))

	// For each grid cell, find closest point (if not tied)
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			minDist := 1 << 30
			closest := -1
			tied := false
			for i, p := range pts {
				dist := manhattan(image.Point{X: x, Y: y}, p)
				if dist < minDist {
					minDist = dist
					closest = i
					tied = false
				} else if dist == minDist {
					tied = true
				}
			}
			if !tied {
				area[closest]++
				// If on border, mark as infinite
				if x == minX || x == maxX || y == minY || y == maxY {
					infinite[closest] = true
				}
			}
		}
	}

	// Find the largest finite area
	maxArea := 0
	for i, a := range area {
		if !infinite[i] && a > maxArea {
			maxArea = a
		}
	}
	return strconv.Itoa(maxArea), nil
}

func (d Day6) Part2(input string) (string, error) {
	pts := parsePoints(input)
	if len(pts) == 0 {
		return "", errors.New("no input")
	}

	// Find bounding box
	minX, maxX := pts[0].X, pts[0].X
	minY, maxY := pts[0].Y, pts[0].Y
	for _, p := range pts {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	regionSize := 0
	const limit = 10000

	// Check each point in the bounding box
	for x := minX - limit/len(pts); x <= maxX+limit/len(pts); x++ {
		for y := minY - limit/len(pts); y <= maxY+limit/len(pts); y++ {
			sum := 0
			for _, p := range pts {
				sum += manhattan(image.Point{X: x, Y: y}, p)
				if sum >= limit {
					break
				}
			}
			if sum < limit {
				regionSize++
			}
		}
	}

	return strconv.Itoa(regionSize), nil
}

func init() {
	solve.Register(Day6{})
}

func init() {
	solve.Register(Day6{})
}
