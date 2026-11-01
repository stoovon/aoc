package solve2016

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day24 struct{}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 24}
}

type ductMaze struct {
	width, height int
	cells         []byte
	points        []int
}

func (d Day24) parseInput(input string) (ductMaze, error) {
	var maze ductMaze
	text := strings.TrimSuffix(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	rows := strings.Split(text, "\n")
	maze.width, maze.height = len(rows[0]), len(rows)
	if maze.width == 0 {
		return maze, fmt.Errorf("empty maze")
	}
	if maze.height > math.MaxInt/maze.width {
		return maze, fmt.Errorf("maze size exceeds int")
	}
	var labels [10]int
	for i := range labels {
		labels[i] = -1
	}
	for y, row := range rows {
		if len(row) != maze.width {
			return maze, fmt.Errorf("row %d: inconsistent width", y+1)
		}
		for x := 0; x < len(row); x++ {
			cell := row[x]
			if cell == '.' || cell == '#' {
				continue
			}
			if cell < '0' || cell > '9' {
				return maze, fmt.Errorf("row %d: expected wall, floor or numbered point", y+1)
			}
			label := int(cell - '0')
			if labels[label] >= 0 {
				return maze, fmt.Errorf("duplicate numbered point")
			}
			labels[label] = y*maze.width + x
		}
	}
	if labels[0] < 0 {
		return maze, fmt.Errorf("missing starting point zero")
	}
	maze.cells = make([]byte, 0, maze.width*maze.height)
	for _, row := range rows {
		maze.cells = append(maze.cells, row...)
	}
	for _, cell := range labels {
		if cell >= 0 {
			maze.points = append(maze.points, cell)
		}
	}
	return maze, nil
}

func (maze ductMaze) distances() ([][]int, error) {
	paths := make([][]int, len(maze.points))
	distance := make([]int, len(maze.cells))
	queue := make([]int, 0, len(maze.cells))
	for index, start := range maze.points {
		for i := range distance {
			distance[i] = -1
		}
		queue = append(queue[:0], start)
		distance[start] = 0
		for head := 0; head < len(queue); head++ {
			cell := queue[head]
			x, y := cell%maze.width, cell/maze.width
			for _, step := range [...][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
				nx, ny := x+step[0], y+step[1]
				if nx < 0 || nx >= maze.width || ny < 0 || ny >= maze.height {
					continue
				}
				next := ny*maze.width + nx
				if maze.cells[next] == '#' || distance[next] >= 0 {
					continue
				}
				distance[next] = distance[cell] + 1
				queue = append(queue, next)
			}
		}
		paths[index] = make([]int, len(maze.points))
		for target, cell := range maze.points {
			if distance[cell] < 0 {
				return nil, fmt.Errorf("numbered points are disconnected")
			}
			paths[index][target] = distance[cell]
		}
	}
	return paths, nil
}

func (d Day24) shortestRoute(paths [][]int, returnToStart bool) (int, error) {
	count := len(paths)
	if count == 0 || count > 10 {
		return 0, fmt.Errorf("expected one to ten numbered points")
	}
	for _, row := range paths {
		if len(row) != count {
			return 0, fmt.Errorf("incomplete distance matrix")
		}
		for _, distance := range row {
			if distance < 0 {
				return 0, fmt.Errorf("numbered points are disconnected")
			}
		}
	}
	// Cache remaining costs, independent of the route prefix that reached a state.
	memo := make([]int, (1<<count)*count)
	for i := range memo {
		memo[i] = -1
	}
	var visit func(int, int) int
	visit = func(current, remaining int) int {
		if remaining == 0 {
			if returnToStart {
				return paths[current][0]
			}
			return 0
		}
		key := remaining*count + current
		if memo[key] != -1 {
			return memo[key]
		}
		best := -2 // All routes from this state may exceed int.
		for next := 1; next < count; next++ {
			bit := 1 << next
			if remaining&bit == 0 {
				continue
			}
			tail := visit(next, remaining&^bit)
			if tail < 0 || paths[current][next] > math.MaxInt-tail {
				continue
			}
			cost := paths[current][next] + tail
			if best < 0 || cost < best {
				best = cost
			}
		}
		memo[key] = best
		return best
	}
	result := visit(0, ((1<<count)-1)&^1)
	if result < 0 {
		return 0, fmt.Errorf("route cost exceeds supported int range")
	}
	return result, nil
}

func (d Day24) solvePart(input string, returnToStart bool) (string, error) {
	maze, err := d.parseInput(input)
	if err != nil {
		return "", err
	}
	paths, err := maze.distances()
	if err != nil {
		return "", err
	}

	result, err := d.shortestRoute(paths, returnToStart)
	if err != nil {
		return "", err
	}

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
