package solve2019

import (
	"container/heap"
	"errors"
	"fmt"
	"image"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day18 struct {
}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 18}
}

type keyEdge struct {
	to       int
	steps    int
	required uint32
}

type pqItem struct {
	robots [4]int
	keys   uint32
	steps  int
	index  int
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].steps < pq[j].steps }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i]; pq[j].index = i; pq[i].index = j }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*pqItem)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// Precompute all key-to-key distances and required keys (doors) to reach them
func buildKeyGraph(lines []string, starts []image.Point) ([][]keyEdge, map[byte]int, uint32) {
	h, w := len(lines), len(lines[0])
	keyPos := make(map[byte]image.Point)
	keyIdx := make(map[byte]int)
	var keyList []byte

	// Find all keys and assign indices
	for y, line := range lines {
		for x := range line {
			c := line[x]
			if c >= 'a' && c <= 'z' {
				keyPos[c] = image.Pt(x, y)
			}
		}
	}
	// Add robot starts as '@', '@'+1, ...
	for i, p := range starts {
		key := byte('@' + i)
		keyPos[key] = p
	}
	// Assign indices: robots first, then keys
	for i := 0; i < len(starts); i++ {
		key := byte('@' + i)
		keyIdx[key] = len(keyList)
		keyList = append(keyList, key)
	}
	for c := byte('a'); c <= 'z'; c++ {
		if _, ok := keyPos[c]; ok {
			keyIdx[c] = len(keyList)
			keyList = append(keyList, c)
		}
	}
	// Build edges: for each key/robot, BFS to all other keys
	graph := make([][]keyEdge, len(keyList))
	for from, c := range keyList {
		visited := make(map[image.Point]bool)
		type bfsNode struct {
			pos      image.Point
			steps    int
			required uint32
		}
		q := []bfsNode{{keyPos[c], 0, 0}}
		visited[keyPos[c]] = true
		for len(q) > 0 {
			cur := q[0]
			q = q[1:]
			cell := lines[cur.pos.Y][cur.pos.X]
			// If it's a key (not the start), add edge
			if idx, ok := keyIdx[cell]; ok && idx != from && cell >= 'a' && cell <= 'z' {
				graph[from] = append(graph[from], keyEdge{
					to:       idx,
					steps:    cur.steps,
					required: cur.required,
				})
			}
			// Explore neighbors
			for _, dxy := range []image.Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
				np := cur.pos.Add(dxy)
				if np.X < 0 || np.Y < 0 || np.X >= w || np.Y >= h {
					continue
				}
				if visited[np] {
					continue
				}
				cell := lines[np.Y][np.X]
				if cell == '#' {
					continue
				}
				nreq := cur.required
				if cell >= 'A' && cell <= 'Z' {
					nreq |= 1 << (cell - 'A')
				}
				visited[np] = true
				q = append(q, bfsNode{np, cur.steps + 1, nreq})
			}
		}
	}
	// All keys bitmask
	allKeys := uint32(0)
	for c := byte('a'); c <= 'z'; c++ {
		if _, ok := keyIdx[c]; ok {
			allKeys |= 1 << (c - 'a')
		}
	}
	return graph, keyIdx, allKeys
}

func shortestPathAllKeysFast(lines []string, starts []image.Point) (int, error) {
	graph, _, allKeys := buildKeyGraph(lines, starts)
	numRobots := len(starts)

	// Initial state: robots at their start indices, no keys
	var robots [4]int
	for i := 0; i < numRobots; i++ {
		robots[i] = i
	}
	type stateKey struct {
		robots [4]int
		keys   uint32
	}
	visited := make(map[stateKey]int)
	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{robots: robots, keys: 0, steps: 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*pqItem)
		if cur.keys == allKeys {
			return cur.steps, nil
		}
		sk := stateKey{robots: cur.robots, keys: cur.keys}
		if prev, ok := visited[sk]; ok && prev <= cur.steps {
			continue
		}
		visited[sk] = cur.steps
		// For each robot, try to collect a new key
		for i := 0; i < numRobots; i++ {
			from := cur.robots[i]
			for _, e := range graph[from] {
				// Only consider keys not yet collected
				if e.to < numRobots {
					continue // don't go to other robots' starts
				}
				keyBit := uint32(1) << (e.to - numRobots)
				if cur.keys&keyBit != 0 {
					continue
				}
				// Check if required keys are collected
				if (cur.keys & e.required) != e.required {
					continue
				}
				// Move robot i to this key
				var nextRobots [4]int
				copy(nextRobots[:], cur.robots[:])
				nextRobots[i] = e.to
				nextKeys := cur.keys | keyBit
				nextSteps := cur.steps + e.steps
				nextSk := stateKey{robots: nextRobots, keys: nextKeys}
				if prev, ok := visited[nextSk]; ok && prev <= nextSteps {
					continue
				}
				heap.Push(pq, &pqItem{robots: nextRobots, keys: nextKeys, steps: nextSteps})
			}
		}
	}
	return 0, errors.New("No solution")
}

func (d Day18) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var start image.Point
	for y, line := range lines {
		for x := range line {
			if line[x] == '@' {
				start = image.Pt(x, y)
			}
		}
	}
	steps, err := shortestPathAllKeysFast(lines, []image.Point{start})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", steps), nil
}

func (d Day18) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var sx, sy int
	for y, line := range lines {
		for x := range line {
			if line[x] == '@' {
				sx, sy = x, y
			}
		}
	}
	// Modify the map for four robots
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			c := '#'
			if maths.Abs(dx) == 1 && maths.Abs(dy) == 1 {
				c = '@'
			}
			row := []byte(lines[sy+dy])
			row[sx+dx] = byte(c)
			lines[sy+dy] = string(row)
		}
	}
	var starts []image.Point
	for y, line := range lines {
		for x := range line {
			if line[x] == '@' {
				starts = append(starts, image.Pt(x, y))
			}
		}
	}
	steps, err := shortestPathAllKeysFast(lines, starts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", steps), nil
}

func init() {
	solve.Register(Day18{})
}
