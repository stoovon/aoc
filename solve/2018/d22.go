package solve2018

import (
	"aoc/solve"
	"container/heap"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Day22 struct{}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 22}
}

func (d Day22) parseInput(input string) (int, [2]int) {
	lines := strings.Split(input, "\n")
	depth := 0
	target := [2]int{}
	fmt.Sscanf(lines[0], "depth: %d", &depth)
	fmt.Sscanf(lines[1], "target: %d,%d", &target[0], &target[1])
	return depth, target
}

func calculateRiskLevel(depth int, target [2]int) int {
	width, height := target[0]+1, target[1]+1
	geologicIndex := make([][]int, height+1) // Adjusted to height+1 to avoid out-of-range errors
	erosionLevel := make([][]int, height+1)  // Adjusted to height+1 to avoid out-of-range errors
	for y := 0; y <= height; y++ {
		geologicIndex[y] = make([]int, width+1) // Adjusted to width+1 to avoid out-of-range errors
		erosionLevel[y] = make([]int, width+1)  // Adjusted to width+1 to avoid out-of-range errors
	}

	totalRisk := 0

	for y := 0; y <= target[1]; y++ {
		for x := 0; x <= target[0]; x++ {
			if x == 0 && y == 0 || (x == target[0] && y == target[1]) {
				geologicIndex[y][x] = 0
			} else if y == 0 {
				geologicIndex[y][x] = x * 16807
			} else if x == 0 {
				geologicIndex[y][x] = y * 48271
			} else {
				geologicIndex[y][x] = erosionLevel[y][x-1] * erosionLevel[y-1][x]
			}

			erosionLevel[y][x] = (geologicIndex[y][x] + depth) % 20183
			typeMod := erosionLevel[y][x] % 3
			totalRisk += typeMod
		}
	}

	return totalRisk
}

func (d Day22) Part1(input string) (string, error) {
	depth, target := d.parseInput(input)
	totalRisk := calculateRiskLevel(depth, target)
	return strconv.Itoa(totalRisk), nil
}

type State struct {
	x, y, tool, time int
}

type PriorityQueue []State

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].time < pq[j].time }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(State)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func (d Day22) Part2(input string) (string, error) {
	depth, target := d.parseInput(input)
	width, height := target[0]+50, target[1]+50 // Expand search area beyond target
	geologicIndex := make([][]int, height+1)
	erosionLevel := make([][]int, height+1)
	types := make([][]int, height+1)
	for y := 0; y <= height; y++ {
		geologicIndex[y] = make([]int, width+1)
		erosionLevel[y] = make([]int, width+1)
		types[y] = make([]int, width+1)
	}

	// Precompute region types
	for y := 0; y <= height; y++ {
		for x := 0; x <= width; x++ {
			if x == 0 && y == 0 || (x == target[0] && y == target[1]) {
				geologicIndex[y][x] = 0
			} else if y == 0 {
				geologicIndex[y][x] = x * 16807
			} else if x == 0 {
				geologicIndex[y][x] = y * 48271
			} else {
				geologicIndex[y][x] = erosionLevel[y][x-1] * erosionLevel[y-1][x]
			}

			erosionLevel[y][x] = (geologicIndex[y][x] + depth) % 20183
			types[y][x] = erosionLevel[y][x] % 3
		}
	}

	// Valid tools for each region type
	validTools := map[int][]int{
		0: {1, 2}, // Rocky: Torch, Climbing gear
		1: {0, 2}, // Wet: Climbing gear, Neither
		2: {0, 1}, // Narrow: Torch, Neither
	}

	// Dijkstra's algorithm
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, State{0, 0, 1, 0}) // Start at (0,0) with torch
	visited := make(map[[3]int]bool)

	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(State)
		if current.x == target[0] && current.y == target[1] && current.tool == 1 {
			return strconv.Itoa(current.time), nil
		}

		key := [3]int{current.x, current.y, current.tool}
		if visited[key] {
			continue
		}
		visited[key] = true

		// Move to adjacent regions
		for _, dir := range directions {
			nx, ny := current.x+dir[0], current.y+dir[1]
			if nx < 0 || ny < 0 || nx > width || ny > height {
				continue
			}
			for _, tool := range validTools[types[ny][nx]] {
				if tool == current.tool {
					heap.Push(pq, State{nx, ny, tool, current.time + 1})
				}
			}
		}

		// Switch tools
		for _, tool := range validTools[types[current.y][current.x]] {
			if tool != current.tool {
				heap.Push(pq, State{current.x, current.y, tool, current.time + 7})
			}
		}
	}

	return "", errors.New("No path found")
}

func init() {
	solve.Register(Day22{})
}
