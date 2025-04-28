package solve2016

import (
	"container/heap"
	"image"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/queues"
)

type Day13 struct {
}

type pqItem struct {
	Value    image.Point
	Priority int
}

func (p *pqItem) Score() int {
	return p.Priority
}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 13}
}

func (d Day13) isOpen(p image.Point, favorite int) bool {
	if p.X < 0 || p.Y < 0 {
		return false
	}
	value := p.X*p.X + 3*p.X + 2*p.X*p.Y + p.Y + p.Y*p.Y + favorite
	return d.bitsOn(value)%2 == 0
}

func (d Day13) bitsOn(n int) int {
	count := 0
	for n > 0 {
		count += n & 1
		n >>= 1
	}
	return count
}

func (d Day13) manhattanNeighbours(p image.Point) []image.Point {
	return []image.Point{
		{p.X + 1, p.Y},
		{p.X - 1, p.Y},
		{p.X, p.Y + 1},
		{p.X, p.Y - 1},
	}
}

func (d Day13) openNeighbours(p image.Point, favorite int) []image.Point {
	var neighbors []image.Point
	for _, n := range d.manhattanNeighbours(p) {
		if d.isOpen(n, favorite) {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (d Day13) astarSearch(start, goal image.Point, favorite int) int {
	frontier := &queues.PriorityQueue[*pqItem]{}
	heap.Init(frontier)
	heap.Push(frontier, &pqItem{Value: start, Priority: d.manhattanDistance(start, goal)})

	costSoFar := map[image.Point]int{start: 0}

	for frontier.Len() > 0 {
		current := heap.Pop(frontier).(*pqItem).Value

		if current == goal {
			return costSoFar[current]
		}

		for _, next := range d.openNeighbours(current, favorite) {
			newCost := costSoFar[current] + 1
			if oldCost, ok := costSoFar[next]; !ok || newCost < oldCost {
				costSoFar[next] = newCost
				heap.Push(frontier, &pqItem{Value: next, Priority: newCost + d.manhattanDistance(next, goal)})
			}
		}
	}
	return -1
}

func (d Day13) countLocationsWithin(start image.Point, steps, favorite int) int {
	queue := []image.Point{start}
	distances := map[image.Point]int{start: 0}
	count := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if distances[current] < steps {
			for _, next := range d.openNeighbours(current, favorite) {
				if _, seen := distances[next]; !seen {
					distances[next] = distances[current] + 1
					queue = append(queue, next)
					count++
				}
			}
		}
	}
	return len(distances)
}

func (d Day13) manhattanDistance(p, q image.Point) int {
	return int(math.Abs(float64(p.X-q.X)) + math.Abs(float64(p.Y-q.Y)))
}

func (d Day13) Part1(input string) (string, error) {
	favorite, _ := strconv.Atoi(strings.TrimSpace(input))
	start := image.Point{X: 1, Y: 1}
	goal := image.Point{X: 31, Y: 39}
	steps := d.astarSearch(start, goal, favorite)
	return strconv.Itoa(steps), nil
}

func (d Day13) Part2(input string) (string, error) {
	favorite, _ := strconv.Atoi(strings.TrimSpace(input))
	start := image.Point{X: 1, Y: 1}
	locations := d.countLocationsWithin(start, 50, favorite)
	return strconv.Itoa(locations), nil
}

func init() {
	solve.Register(Day13{})
}
