package solve2018

import (
	"aoc/solve"
	"aoc/utils/grids"
	"image"
	"regexp"
	"strconv"
	"strings"
)

var (
	lineRe = regexp.MustCompile(`position=<\s*(-?\d+),\s*(-?\d+)> velocity=<\s*(-?\d+),\s*(-?\d+)>`)
)

type Day10 struct{}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 10}
}

type point struct {
	image.Point
	vx, vy int
}

func parsePoints(input string) []point {
	var points []point
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		m := lineRe.FindStringSubmatch(line)
		if len(m) != 5 {
			continue
		}
		x, _ := strconv.Atoi(m[1])
		y, _ := strconv.Atoi(m[2])
		vx, _ := strconv.Atoi(m[3])
		vy, _ := strconv.Atoi(m[4])
		points = append(points, point{image.Point{x, y}, vx, vy})
	}
	return points
}

// Simulate points and return (bestPoints, bestStep)
func simulatePoints(points []point, maxSteps int) ([]point, int) {
	minArea := int(^uint(0) >> 1)
	bestStep := 0
	bestPoints := make([]point, len(points))
	pts := make([]point, len(points))
	copy(pts, points)
	for step := 0; step < maxSteps; step++ {
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
		area := (maxX - minX) * (maxY - minY)
		if area < minArea {
			minArea = area
			bestStep = step
			copy(bestPoints, pts)
		}
		// Move points
		for i := range pts {
			pts[i].X += pts[i].vx
			pts[i].Y += pts[i].vy
		}
	}
	return bestPoints, bestStep
}

func renderPoints(points []point) [][]int {
	minX, maxX := points[0].X, points[0].X
	minY, maxY := points[0].Y, points[0].Y
	for _, p := range points {
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
	w, h := maxX-minX+1, maxY-minY+1
	grid := make([][]int, h)
	for i := range grid {
		grid[i] = make([]int, w)
	}
	for _, p := range points {
		grid[p.Y-minY][p.X-minX] = 1
	}
	return grid
}

func (d Day10) Part1(input string) (string, error) {
	points := parsePoints(input)
	bestPoints, _ := simulatePoints(points, 20000)
	return grids.OCRWithWidth(renderPoints(bestPoints), 8), nil
}

func (d Day10) Part2(input string) (string, error) {
	points := parsePoints(input)
	_, bestStep := simulatePoints(points, 20000)
	return strconv.Itoa(bestStep), nil
}

func init() {
	solve.Register(Day10{})
}
