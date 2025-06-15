package solve2019

import (
	"container/list"
	"errors"
	"image"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day20 struct {
}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 20}
}

type mazeData struct {
	grid    map[image.Point]rune
	portals map[string][]image.Point
	adj     map[image.Point][]image.Point
	width   int
	height  int
	outer   map[image.Point]bool
	inner   map[image.Point]bool
}

func parseMaze(input string) mazeData {
	lines := strings.Split(input, "\n")
	grid := make(map[image.Point]rune)
	portals := make(map[string][]image.Point)
	width, height := len(lines[0]), len(lines)
	for y, line := range lines {
		for x, ch := range line {
			grid[image.Point{X: x, Y: y}] = ch
		}
	}
	// Find portals
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			p := image.Point{X: x, Y: y}
			ch := grid[p]
			if ch >= 'A' && ch <= 'Z' {
				for _, dxy := range []image.Point{{1, 0}, {0, 1}} {
					np := image.Point{X: x + dxy.X, Y: y + dxy.Y}
					if grid[np] >= 'A' && grid[np] <= 'Z' {
						label := string([]rune{ch, grid[np]})
						var op image.Point
						if grid[image.Point{X: x - dxy.X, Y: y - dxy.Y}] == '.' {
							op = image.Point{X: x - dxy.X, Y: y - dxy.Y}
						} else if grid[image.Point{X: x + 2*dxy.X, Y: y + 2*dxy.Y}] == '.' {
							op = image.Point{X: x + 2*dxy.X, Y: y + 2*dxy.Y}
						} else {
							continue
						}
						portals[label] = append(portals[label], op)
					}
				}
			}
		}
	}
	// Build adjacency
	adj := make(map[image.Point][]image.Point)
	dirs := []image.Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for p, ch := range grid {
		if ch != '.' {
			continue
		}
		for _, dxy := range dirs {
			np := image.Point{X: p.X + dxy.X, Y: p.Y + dxy.Y}
			if grid[np] == '.' {
				adj[p] = append(adj[p], np)
			}
		}
	}
	// Add portal connections
	for _, pts := range portals {
		if len(pts) == 2 {
			adj[pts[0]] = append(adj[pts[0]], pts[1])
			adj[pts[1]] = append(adj[pts[1]], pts[0])
		}
	}
	// Find bounds of open area
	minX, maxX, minY, maxY := width, 0, height, 0
	for p, ch := range grid {
		if ch == '.' {
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
	}
	// Classify portals as inner/outer using open area bounds
	outer, inner := make(map[image.Point]bool), make(map[image.Point]bool)
	for _, pts := range portals {
		for _, p := range pts {
			if p.X == minX || p.X == maxX || p.Y == minY || p.Y == maxY {
				outer[p] = true
			} else {
				inner[p] = true
			}
		}
	}
	return mazeData{grid, portals, adj, width, height, outer, inner}
}

func bfsFlat(m mazeData, start, end image.Point) int {
	visited := map[image.Point]bool{start: true}
	queue := list.New()
	queue.PushBack(struct {
		p     image.Point
		steps int
	}{start, 0})
	for queue.Len() > 0 {
		curr := queue.Remove(queue.Front()).(struct {
			p     image.Point
			steps int
		})
		if curr.p == end {
			return curr.steps
		}
		for _, np := range m.adj[curr.p] {
			if !visited[np] {
				visited[np] = true
				queue.PushBack(struct {
					p     image.Point
					steps int
				}{np, curr.steps + 1})
			}
		}
	}
	return -1
}

func bfsRecursive(m mazeData, start, end image.Point) int {
	type state struct {
		p     image.Point
		level int
	}
	visited := map[state]bool{{start, 0}: true}
	queue := list.New()
	queue.PushBack(struct {
		s     state
		steps int
	}{state{start, 0}, 0})
	for queue.Len() > 0 {
		curr := queue.Remove(queue.Front()).(struct {
			s     state
			steps int
		})
		if curr.s.p == end && curr.s.level == 0 {
			return curr.steps
		}
		for _, np := range m.adj[curr.s.p] {
			nextLevel := curr.s.level
			skip := false
			for label, pts := range m.portals {
				if label == "AA" || label == "ZZ" {
					continue
				}
				if len(pts) == 2 && (pts[0] == curr.s.p || pts[1] == curr.s.p) && (pts[0] == np || pts[1] == np) {
					if m.outer[curr.s.p] {
						if curr.s.level == 0 {
							skip = true
							break
						}
						nextLevel--
					} else if m.inner[curr.s.p] {
						nextLevel++
					}
				}
			}
			if skip {
				continue
			}
			// Don't enter AA/ZZ except at level 0
			if (np == m.portals["AA"][0] || np == m.portals["ZZ"][0]) && nextLevel != 0 {
				continue
			}
			st := state{np, nextLevel}
			if !visited[st] {
				visited[st] = true
				queue.PushBack(struct {
					s     state
					steps int
				}{st, curr.steps + 1})
			}
		}
	}
	return -1
}

func (d Day20) Part1(input string) (string, error) {
	m := parseMaze(input)
	start, end := m.portals["AA"][0], m.portals["ZZ"][0]
	steps := bfsFlat(m, start, end)
	if steps < 0 {
		return "", errors.New("No path found")
	}
	return strconv.Itoa(steps), nil
}

func (d Day20) Part2(input string) (string, error) {
	m := parseMaze(input)
	start, end := m.portals["AA"][0], m.portals["ZZ"][0]
	steps := bfsRecursive(m, start, end)
	if steps < 0 {
		return "", errors.New("No path found")
	}
	return strconv.Itoa(steps), nil
}

func init() {
	solve.Register(Day20{})
}
