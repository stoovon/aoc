package solve2016

import (
	"image"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day24 struct {
}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 24}
}

func (d Day24) parseInput(input string) (map[image.Point][]image.Point, []image.Point) {
	lines := strings.Split(input, "\n")
	nodes := make(map[image.Point][]image.Point)
	points := make([]image.Point, 0)

	var grid [][]string
	for _, line := range lines {
		if line == "" {
			continue
		}
		grid = append(grid, strings.Split(line, ""))
	}

	for y, row := range grid {
		for x, item := range row {
			if item == "#" {
				continue
			}

			point := image.Point{X: x, Y: y}
			nodes[point] = []image.Point{}

			if num, err := strconv.Atoi(item); err == nil {
				if len(points) <= num {
					points = append(points, make([]image.Point, num-len(points)+1)...)
				}
				points[num] = point
			}
		}
	}

	for point := range nodes {
		candidates := []image.Point{
			{X: point.X + 1, Y: point.Y},
			{X: point.X - 1, Y: point.Y},
			{X: point.X, Y: point.Y + 1},
			{X: point.X, Y: point.Y - 1},
		}

		for _, neighbor := range candidates {
			if _, exists := nodes[neighbor]; exists {
				nodes[point] = append(nodes[point], neighbor)
			}
		}
	}

	return nodes, points
}

func (d Day24) bfs(nodes map[image.Point][]image.Point, points []image.Point, start int) []int {
	queue := []struct {
		Point image.Point
		Time  int
	}{{Point: points[start], Time: 0}}
	visited := make(map[image.Point]bool)
	visited[points[start]] = true

	paths := make([]int, len(points))
	for i := range paths {
		paths[i] = math.MaxInt
	}
	paths[start] = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for i, point := range points {
			if current.Point == point {
				paths[i] = current.Time
			}
		}

		allVisited := true
		for _, path := range paths {
			if path == math.MaxInt {
				allVisited = false
				break
			}
		}
		if allVisited {
			return paths
		}

		for _, neighbor := range nodes[current.Point] {
			if !visited[neighbor] {
				queue = append(queue, struct {
					Point image.Point
					Time  int
				}{Point: neighbor, Time: current.Time + 1})
				visited[neighbor] = true
			}
		}
	}

	return paths
}

func (d Day24) solve(paths [][]int, remaining []int, start, total int, returnToStart bool) int {
	if len(remaining) == 0 {
		if returnToStart {
			return total + paths[start][0] // Add the cost to return to the starting point
		}
		return total
	}

	minTotal := math.MaxInt
	for i, key := range remaining {
		newRemaining := append(append([]int{}, remaining[:i]...), remaining[i+1:]...)
		cost := d.solve(paths, newRemaining, key, total+paths[start][key], returnToStart)
		if cost < minTotal {
			minTotal = cost
		}
	}

	return minTotal
}

func (d Day24) solvePart(input string, returnToStart bool) (string, error) {
	nodes, points := d.parseInput(input)

	paths := make([][]int, len(points))
	for i := range points {
		paths[i] = d.bfs(nodes, points, i)
	}

	remaining := make([]int, len(points)-1)
	for i := range remaining {
		remaining[i] = i + 1
	}

	result := d.solve(paths, remaining, 0, 0, returnToStart)

	return strconv.Itoa(result), nil
}

func (d Day24) Part1(input string) (string, error) {
	return d.solvePart(input, false)
}

func (d Day24) Part2(input string) (string, error) {
	return d.solvePart(input, true)
}

func init() {
	solve.Register(Day24{})
}
