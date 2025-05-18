package solve2023

import (
	"container/heap"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day17 struct{}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 17}
}

type direction int

const (
	north direction = iota
	east
	south
	west
)

var directions = []direction{north, east, south, west}

func (d direction) reflect() direction {
	switch d {
	case north:
		return south
	case east:
		return west
	case south:
		return north
	case west:
		return east
	}
	panic("bad direction")
}

type day17Node struct {
	pos      int
	dir      *direction
	distance int
	cost     int
	index    int // for heap
}

type nodeHeap []*day17Node

func (h nodeHeap) Len() int            { return len(h) }
func (h nodeHeap) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h nodeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *nodeHeap) Push(x interface{}) { *h = append(*h, x.(*day17Node)) }
func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func parseGrid(input string) ([]int, int, int) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rows := len(lines)
	cols := len(lines[0])
	vals := make([]int, 0, rows*cols)
	for _, line := range lines {
		for _, c := range line {
			vals = append(vals, int(c-'0'))
		}
	}
	return vals, rows, cols
}

func (d Day17) solve(input string, min, max int) int {
	tiles, rows, cols := parseGrid(input)
	type histEntry struct {
		visited bool
		cost    int
	}
	history := make([]histEntry, len(tiles)*4*max)
	for i := range history {
		history[i].cost = 1 << 30
	}
	h := &nodeHeap{}
	heap.Init(h)
	heap.Push(h, &day17Node{pos: 0, dir: nil, distance: 0, cost: 0})

	for h.Len() > 0 {
		n := heap.Pop(h).(*day17Node)
		var dirIdx int
		if n.dir != nil {
			dirIdx = int(*n.dir)
			history[n.pos*4*max+dirIdx*max+n.distance].visited = true
		} else {
			for d := 0; d < 4; d++ {
				history[n.pos*4*max+d*max+n.distance].visited = true
			}
		}
		for _, extendDir := range directions {
			sameDir, oppDir := true, false
			if n.dir != nil {
				sameDir = *n.dir == extendDir
				oppDir = n.dir.reflect() == extendDir
			}
			if (n.distance < min && !sameDir) ||
				(n.distance > max-1 && sameDir) ||
				oppDir ||
				(extendDir == north && n.pos < cols) ||
				(extendDir == east && n.pos%cols == cols-1) ||
				(extendDir == south && n.pos/cols == rows-1) ||
				(extendDir == west && n.pos%cols == 0) {
				continue
			}
			var next int
			switch extendDir {
			case north:
				next = n.pos - cols
			case east:
				next = n.pos + 1
			case south:
				next = n.pos + cols
			case west:
				next = n.pos - 1
			}
			ndist := 1
			if sameDir {
				ndist += n.distance
			}
			nkey := next*4*max + int(extendDir)*max + ndist
			ncost := n.cost + tiles[next]
			entry := &history[nkey]
			if entry.visited || entry.cost <= ncost {
				continue
			}
			entry.cost = ncost
			dirCopy := extendDir
			heap.Push(h, &day17Node{pos: next, dir: &dirCopy, distance: ndist, cost: ncost})
		}
	}
	// Find min cost at the last tile
	minCost := 1 << 30
	for i := (len(tiles) - 1) * 4 * max; i < (len(tiles))*4*max; i++ {
		if history[i].cost < minCost {
			minCost = history[i].cost
		}
	}
	return minCost
}

func (d Day17) Part1(input string) (string, error) {
	return strconv.Itoa(d.solve(input, 0, 3)), nil
}

func (d Day17) Part2(input string) (string, error) {
	return strconv.Itoa(d.solve(input, 4, 10)), nil
}

func init() {
	solve.Register(Day17{})
}
