package solve2021

import (
	"aoc/solve"
	"container/heap"
	"errors"
	"image"
	"strconv"
	"strings"
)

type Day15 struct{}

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 15}
}

func (d Day15) Part1(input string) (string, error) {
	grid := importGrid(input)
	height := len(grid)
	width := len(grid[0])

	dist := make(map[image.Point]int)
	start := image.Point{X: 0, Y: 0}
	end := image.Point{X: width - 1, Y: height - 1}
	dist[start] = 0
	pq := &PQ{}
	heap.Init(pq)
	heap.Push(pq, Item{Pt: start, Score: 0})

	dirs := []image.Point{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}}
	for pq.Len() > 0 {
		curr := heap.Pop(pq).(Item)
		if curr.Pt == end {
			return strconv.Itoa(curr.Score), nil
		}
		for _, d := range dirs {
			nx, ny := curr.Pt.X+d.X, curr.Pt.Y+d.Y
			if nx < 0 || ny < 0 || nx >= width || ny >= height {
				continue
			}
			next := image.Point{X: nx, Y: ny}
			newScore := curr.Score + grid[ny][nx]
			if old, ok := dist[next]; !ok || newScore < old {
				dist[next] = newScore
				heap.Push(pq, Item{Pt: next, Score: newScore})
			}
		}
	}
	return "", errors.New("no path found")
}

type PQ []Item
type Item struct {
	Pt    image.Point
	Score int
}

func (pq *PQ) Len() int           { return len(*pq) }
func (pq *PQ) Less(i, j int) bool { return (*pq)[i].Score < (*pq)[j].Score }
func (pq *PQ) Swap(i, j int)      { (*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i] }
func (pq *PQ) Push(x any)         { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
func importGrid(input string) [][]int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]int, len(lines))
	for y, line := range lines {
		grid[y] = make([]int, len(line))
		for x, ch := range line {
			grid[y][x] = int(ch - '0')
		}
	}
	return grid
}

func (d Day15) Part2(input string) (string, error) {
	base := importGrid(input)
	baseH := len(base)
	baseW := len(base[0])
	size := 5
	height := baseH * size
	width := baseW * size

	// Expand grid
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			inc := (x / baseW) + (y / baseH)
			val := base[y%baseH][x%baseW] + inc
			if val > 9 {
				val = (val-1)%9 + 1
			}
			grid[y][x] = val
		}
	}

	dist := make(map[image.Point]int)
	start := image.Point{X: 0, Y: 0}
	end := image.Point{X: width - 1, Y: height - 1}
	dist[start] = 0
	pq := &PQ{}
	heap.Init(pq)
	heap.Push(pq, Item{Pt: start, Score: 0})

	dirs := []image.Point{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}}
	for pq.Len() > 0 {
		curr := heap.Pop(pq).(Item)
		if curr.Pt == end {
			return strconv.Itoa(curr.Score), nil
		}
		for _, d := range dirs {
			nx, ny := curr.Pt.X+d.X, curr.Pt.Y+d.Y
			if nx < 0 || ny < 0 || nx >= width || ny >= height {
				continue
			}
			next := image.Point{X: nx, Y: ny}
			newScore := curr.Score + grid[ny][nx]
			if old, ok := dist[next]; !ok || newScore < old {
				dist[next] = newScore
				heap.Push(pq, Item{Pt: next, Score: newScore})
			}
		}
	}
	return "", errors.New("no path found")
}

func init() {
	solve.Register(Day15{})
}
