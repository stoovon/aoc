package solve2024

import (
	"container/heap"
	"fmt"
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
	"aoc/utils/queues"
)

type Day16 struct {
}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 16}
}

type Maze struct {
	grid  [][]string
	start image.Point
	end   image.Point
}

type QueueItem struct {
	pos   image.Point
	dir   int
	score int
	path  []image.Point
}

func (q QueueItem) Score() int {
	return q.score
}

var (
	directions = grids.URDL()
)

const (
	Wall     = "#"
	Start    = "S"
	End      = "E"
	TurnCost = 1000
	MoveCost = 1
	StartDir = 1 // start facing right
)

func parseMaze(input string) Maze {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]string, len(lines))
	var start, end image.Point

	for y, line := range lines {
		grid[y] = strings.Split(line, "")
		for x, ch := range grid[y] {
			switch ch {
			case Start:
				start = image.Point{X: x, Y: y}
			case End:
				end = image.Point{X: x, Y: y}
			}
		}
	}

	return Maze{grid, start, end}
}

func findLowestScore(m Maze) int {
	pq := &queues.PriorityQueue[QueueItem]{}
	heap.Init(pq)
	heap.Push(pq, QueueItem{m.start, StartDir, 0, nil})
	visited := make(map[string]bool)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(QueueItem)

		if m.isEnd(current.pos) {
			return current.score
		}

		key := key(current.pos, current.dir)
		if visited[key] {
			continue
		}
		visited[key] = true

		// try moving forward in current direction
		nextPos := current.pos.Add(directions[current.dir])
		if m.isValid(nextPos) {
			heap.Push(pq, QueueItem{
				nextPos,
				current.dir,
				current.score + MoveCost,
				nil,
			})
		}

		// try both possible 90 degree turns
		heap.Push(pq, QueueItem{current.pos, (current.dir + 1) % 4, current.score + TurnCost, nil})
		heap.Push(pq, QueueItem{current.pos, (current.dir + 3) % 4, current.score + TurnCost, nil})
	}

	return -1
}

// similar to findLowestScore but keeps track of all paths that achieve
// the target score. visited map stores scores now instead of just bools
// since we want all paths with exactly the target score
func findAllOptimalPaths(m Maze, targetScore int) [][]image.Point {
	queue := []QueueItem{{m.start, StartDir, 0, []image.Point{m.start}}}
	visited := make(map[string]int)
	var paths [][]image.Point

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.score > targetScore {
			continue
		}

		key := key(current.pos, current.dir)
		if score, exists := visited[key]; exists && score < current.score {
			continue
		}
		visited[key] = current.score

		if m.isEnd(current.pos) && current.score == targetScore {
			paths = append(paths, current.path)
			continue
		}

		// same movement logic as findLowestScore but we need to
		// track and copy paths now
		nextPos := current.pos.Add(directions[current.dir])
		if m.isValid(nextPos) {
			newPath := make([]image.Point, len(current.path))
			copy(newPath, current.path)
			queue = append(queue, QueueItem{
				nextPos,
				current.dir,
				current.score + MoveCost,
				append(newPath, nextPos),
			})
		}

		// handle turns, keeping the same path since we haven't moved
		for _, newDir := range []int{(current.dir + 1) % 4, (current.dir + 3) % 4} {
			queue = append(queue, QueueItem{
				current.pos,
				newDir,
				current.score + TurnCost,
				current.path,
			})
		}
	}

	return paths
}

func countUniqueTiles(paths [][]image.Point) int {
	unique := make(map[string]bool)
	for _, path := range paths {
		for _, p := range path {
			unique[key(p, 0)] = true
		}
	}
	return len(unique)
}

func key(p image.Point, dir int) string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, dir)
}

func (m Maze) isValid(p image.Point) bool {
	return p.Y >= 0 && p.Y < len(m.grid) &&
		p.X >= 0 && p.X < len(m.grid[0]) &&
		m.grid[p.Y][p.X] != Wall
}

func (m Maze) isEnd(p image.Point) bool {
	return p == m.end
}

func (d Day16) Part1(input string) (string, error) {
	maze := parseMaze(input)
	return strconv.Itoa(findLowestScore(maze)), nil
}

func (d Day16) Part2(input string) (string, error) {
	maze := parseMaze(input)

	lowestScore := findLowestScore(maze)
	paths := findAllOptimalPaths(maze, lowestScore)
	return strconv.Itoa(countUniqueTiles(paths)), nil
}

func init() {
	solve.Register(Day16{})
}
