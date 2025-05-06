package solve2016

import (
	"container/heap"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day22 struct {
}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 22}
}

var (
	digitRe = regexp.MustCompile(`\d+`)
)

type Node struct {
	X, Y, Size, Used, Avail, Pct int
}

type State struct {
	DataPos, EmptyPos [2]int
}

type PriorityQueueItem struct {
	State    State
	Priority int
	Index    int
}

type PriorityQueue []*PriorityQueueItem

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*PriorityQueueItem)
	item.Index = len(*pq)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

func parseInts(s string) []int {
	matches := digitRe.FindAllString(s, -1)
	ints := make([]int, len(matches))
	for i, m := range matches {
		ints[i], _ = strconv.Atoi(m)
	}
	return ints
}

func manhattanDistance(a, b [2]int) int {
	return maths.Abs(a[0]-b[0]) + maths.Abs(a[1]-b[1])
}

func manhattanNeighbours(pos [2]int) [][2]int {
	return [][2]int{
		{pos[0] - 1, pos[1]},
		{pos[0] + 1, pos[1]},
		{pos[0], pos[1] - 1},
		{pos[0], pos[1] + 1},
	}
}

func (d Day22) viable(a, b Node) bool {
	return a != b && a.Used > 0 && a.Used <= b.Avail
}

func (d Day22) astarSearch(start State, hFunc func(State) int, movesFunc func(State) []State) []State {
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &PriorityQueueItem{State: start, Priority: hFunc(start)})

	previous := make(map[State]*State)
	pathCost := make(map[State]int)
	pathCost[start] = 0

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*PriorityQueueItem).State
		if hFunc(current) == 0 {
			return d.reconstructPath(previous, current)
		}
		for _, next := range movesFunc(current) {
			newCost := pathCost[current] + 1
			if cost, ok := pathCost[next]; !ok || newCost < cost {
				pathCost[next] = newCost
				heap.Push(pq, &PriorityQueueItem{State: next, Priority: newCost + hFunc(next)})
				previous[next] = &current
			}
		}
	}
	return nil
}

func (d Day22) reconstructPath(previous map[State]*State, current State) []State {
	var path []State
	for currentPtr := &current; currentPtr != nil; currentPtr = previous[*currentPtr] {
		path = append([]State{*currentPtr}, path...)
	}
	return path
}

func (d Day22) solve(input string) (string, string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var nodes []Node
	for _, line := range lines {
		if strings.HasPrefix(line, "/dev") {
			ints := parseInts(line)
			nodes = append(nodes, Node{X: ints[0], Y: ints[1], Size: ints[2], Used: ints[3], Avail: ints[4], Pct: ints[5]})
		}
	}

	var empty Node
	maxX := 0
	grid := make(map[[2]int]Node)
	for _, node := range nodes {
		grid[[2]int{node.X, node.Y}] = node
		if node.Used == 0 {
			empty = node
		}
		if node.X > maxX {
			maxX = node.X
		}
	}

	start := State{DataPos: [2]int{maxX, 0}, EmptyPos: [2]int{empty.X, empty.Y}}
	hFunc := func(s State) int { return manhattanDistance(s.DataPos, [2]int{0, 0}) }
	movesFunc := func(s State) []State {
		var moves []State
		for _, pos := range manhattanNeighbours(s.EmptyPos) {
			if node, ok := grid[pos]; ok {
				empty := grid[s.EmptyPos]
				if node.Used <= empty.Size {
					newDataPos := s.DataPos
					if pos == s.DataPos {
						newDataPos = s.EmptyPos
					}
					moves = append(moves, State{DataPos: newDataPos, EmptyPos: pos})
				}
			}
		}
		return moves
	}

	path := d.astarSearch(start, hFunc, movesFunc)
	viablePairs := 0
	for _, a := range nodes {
		for _, b := range nodes {
			if d.viable(a, b) {
				viablePairs++
			}
		}
	}

	return strconv.Itoa(viablePairs), strconv.Itoa(len(path) - 1), nil
}

func (d Day22) Part1(input string) (string, error) {
	data, _, err := d.solve(input)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (d Day22) Part2(input string) (string, error) {
	_, data, err := d.solve(input)
	if err != nil {
		return "", err
	}
	return data, nil
}

func init() {
	solve.Register(Day22{})
}
