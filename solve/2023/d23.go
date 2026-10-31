package solve2023

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day23 struct {
}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 23}
}

type day23Grid struct {
	data   []byte
	width  int
	height int
	start  int
	end    int
}

var arrows = []byte{'^', '>', 'v', '<'}

func newDay23Grid(input string) (*day23Grid, error) {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	width := len(lines[0])
	if width == 0 {
		return nil, fmt.Errorf("expected a nonempty hiking map")
	}
	data := make([]byte, 0)
	for _, line := range lines {
		if len(line) != width {
			return nil, fmt.Errorf("expected rectangular hiking map")
		}
		for _, tile := range []byte(line) {
			if !strings.ContainsRune("#.^>v<", rune(tile)) {
				return nil, fmt.Errorf("invalid hiking tile")
			}
		}
		data = append(data, line...)
	}
	opening := func(row string) (int, error) {
		position := -1
		for i := range row {
			if row[i] != '#' {
				if position != -1 {
					return 0, fmt.Errorf("expected exactly one opening in each boundary row")
				}
				position = i
			}
		}
		if position == -1 {
			return 0, fmt.Errorf("expected exactly one opening in each boundary row")
		}
		return position, nil
	}
	start, err := opening(lines[0])
	if err != nil {
		return nil, err
	}
	end, err := opening(lines[len(lines)-1])
	if err != nil {
		return nil, err
	}
	return &day23Grid{data: data, width: width, height: len(lines), start: start, end: (len(lines)-1)*width + end}, nil
}

func (g *day23Grid) nextPos(p int, dir int) (int, bool) {
	switch dir {
	case 0:
		if p >= g.width {
			return p - g.width, true
		}
	case 1:
		if (p+1)%g.width != 0 {
			return p + 1, true
		}
	case 2:
		if p < len(g.data)-g.width {
			return p + g.width, true
		}
	case 3:
		if p%g.width != 0 {
			return p - 1, true
		}
	}
	return 0, false
}

func (g *day23Grid) canLeave(position, direction int, slopes bool) bool {
	if !slopes {
		return true
	}
	for arrow, tile := range arrows {
		if g.data[position] == tile {
			return arrow == direction
		}
	}
	return true
}

func (g *day23Grid) intersectionsAndDistances(slopes bool) [][][2]int {
	vertices := []int{g.start, g.end}
	for position, tile := range g.data {
		if tile == '#' {
			continue
		}
		degree := 0
		for direction := 0; direction < 4; direction++ {
			if neighbor, ok := g.nextPos(position, direction); ok && g.data[neighbor] != '#' {
				degree++
			}
		}
		if degree > 2 {
			vertices = append(vertices, position)
		}
	}
	sort.Ints(vertices)
	unique := vertices[:0]
	for _, position := range vertices {
		if len(unique) == 0 || unique[len(unique)-1] != position {
			unique = append(unique, position)
		}
	}
	vertices = unique
	vertexAt := make([]int, len(g.data))
	for i := range vertexAt {
		vertexAt[i] = -1
	}
	for id, position := range vertices {
		vertexAt[position] = id
	}
	edges := make([][][2]int, len(vertices))
	for source, position := range vertices {
		for direction := 0; direction < 4; direction++ {
			if !g.canLeave(position, direction, slopes) {
				continue
			}
			current, ok := g.nextPos(position, direction)
			if !ok || g.data[current] == '#' {
				continue
			}
			previous, distance := position, 1
			for {
				if target := vertexAt[current]; target >= 0 {
					if target != source {
						found := false
						for i, edge := range edges[source] {
							if edge[0] == target {
								edges[source][i][1] = max(edge[1], distance)
								found = true
								break
							}
						}
						if !found {
							edges[source] = append(edges[source], [2]int{target, distance})
						}
					}
					break
				}
				next := -1
				for d := 0; d < 4; d++ {
					if !g.canLeave(current, d, slopes) {
						continue
					}
					if neighbor, ok := g.nextPos(current, d); ok && neighbor != previous && g.data[neighbor] != '#' {
						next = neighbor
						break
					}
				}
				if next < 0 {
					break
				}
				previous, current = current, next
				distance++
			}
		}
	}
	return edges
}

func hikeLongestPath(position int, edges [][][2]int, visited []bool, end int) (int, bool) {
	if position == end {
		return 0, true
	}
	visited[position] = true
	best, exists := 0, false
	for _, edge := range edges[position] {
		if !visited[edge[0]] {
			if remaining, ok := hikeLongestPath(edge[0], edges, visited, end); ok {
				best = max(best, edge[1]+remaining)
				exists = true
			}
		}
	}
	visited[position] = false
	return best, exists
}

func hikeLongestRoute(edges [][][2]int, start, end int) (int, bool) {
	if start == end {
		return 0, true
	}
	visited := make([]bool, len(edges))
	predecessor, distance, incoming := 0, 0, 0
	for source, outgoing := range edges {
		for _, edge := range outgoing {
			if edge[0] == end {
				predecessor, distance = source, edge[1]
				incoming++
			}
		}
	}
	if incoming == 1 {
		visited[end] = true
		length, exists := hikeLongestPath(start, edges, visited, predecessor)
		if !exists {
			return 0, false
		}
		return length + distance, true
	}
	return hikeLongestPath(start, edges, visited, end)
}

func solveHike(input string, slopes bool) (string, error) {
	grid, err := newDay23Grid(input)
	if err != nil {
		return "", err
	}
	edges := grid.intersectionsAndDistances(slopes)
	distance, exists := hikeLongestRoute(edges, 0, len(edges)-1)
	if !exists {
		return "", fmt.Errorf("no route from entrance to exit")
	}
	return strconv.Itoa(distance), nil
}

func (d Day23) Part1(input string) (string, error) { return solveHike(input, true) }
func (d Day23) Part2(input string) (string, error) { return solveHike(input, false) }

func init() {
	solve.Register(Day23{})
}
