package solve2025

import (
	"container/heap"
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day8 struct {
}

type vec3 struct {
	x, y, z int64
}

type edge struct {
	i, j int
	d2   int64
}

type dsuResult struct {
	componentSizes []int64
	lastMergeI     int
	lastMergeJ     int
}

func parsePoints(input string) ([]vec3, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	points := make([]vec3, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return nil, strconv.ErrSyntax
		}
		x, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			return nil, err
		}
		z, err := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
		if err != nil {
			return nil, err
		}
		points = append(points, vec3{x: x, y: y, z: z})
	}

	if len(points) < 2 {
		return nil, errors.New("at least two junction boxes are required")
	}
	return points, nil
}

func edgeLess(a, b edge) bool {
	if a.d2 != b.d2 {
		return a.d2 < b.d2
	}
	if a.i != b.i {
		return a.i < b.i
	}
	return a.j < b.j
}

type farthestEdges []edge

func (h farthestEdges) Len() int           { return len(h) }
func (h farthestEdges) Less(i, j int) bool { return edgeLess(h[j], h[i]) }
func (h farthestEdges) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *farthestEdges) Push(value any)    { *h = append(*h, value.(edge)) }
func (h *farthestEdges) Pop() any {
	last := len(*h) - 1
	value := (*h)[last]
	*h = (*h)[:last]
	return value
}

func squaredDistance(a, b vec3) (int64, error) {
	var total int64
	for _, pair := range [3][2]int64{{a.x, b.x}, {a.y, b.y}, {a.z, b.z}} {
		x, y := pair[0], pair[1]
		if (y > 0 && x < math.MinInt64+y) || (y < 0 && x > math.MaxInt64+y) {
			return 0, errors.New("coordinate difference overflows int64")
		}
		delta := x - y
		if delta == math.MinInt64 {
			return 0, errors.New("squared distance overflows int64")
		}
		if delta < 0 {
			delta = -delta
		}
		if delta != 0 && delta > math.MaxInt64/delta {
			return 0, errors.New("squared distance overflows int64")
		}
		square := delta * delta
		if total > math.MaxInt64-square {
			return 0, errors.New("squared distance overflows int64")
		}
		total += square
	}
	return total, nil
}

func multiplyCircuitAnswer(a, b int64) (int64, error) {
	if (a == math.MinInt64 && b == -1) || (b == math.MinInt64 && a == -1) {
		return 0, errors.New("answer product overflows int64")
	}
	product := a * b
	if a != 0 && product/a != b {
		return 0, errors.New("answer product overflows int64")
	}
	return product, nil
}

func buildEdges(points []vec3, maxEdges int) ([]edge, error) {
	n := len(points)
	capacity := 0
	if n > 1 {
		a, b := n, n-1
		if a%2 == 0 {
			a /= 2
		} else {
			b /= 2
		}
		if a > int(^uint(0)>>1)/b {
			return nil, errors.New("pair count overflows int")
		}
		capacity = a * b
	}
	if maxEdges > 0 {
		capacity = min(capacity, maxEdges)
	}
	edges := make(farthestEdges, 0, capacity)
	for i := range n {
		for j := i + 1; j < n; j++ {
			d2, err := squaredDistance(points[i], points[j])
			if err != nil {
				return nil, err
			}
			candidate := edge{i: i, j: j, d2: d2}
			if maxEdges <= 0 || len(edges) < maxEdges {
				edges = append(edges, candidate)
				if maxEdges > 0 && len(edges) == maxEdges {
					heap.Init(&edges)
				}
			} else if edgeLess(candidate, edges[0]) {
				edges[0] = candidate
				heap.Fix(&edges, 0)
			}
		}
	}
	return edges, nil
}

func processEdges(n int, edges []edge, maxConnections int, stopWhenSingle bool) dsuResult {
	sort.Slice(edges, func(a, b int) bool { return edgeLess(edges[a], edges[b]) })

	parent := make([]int, n)
	size := make([]int64, n)
	for i := range n {
		parent[i] = i
		size[i] = 1
	}

	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(a, b int) bool {
		ra := find(a)
		rb := find(b)
		if ra == rb {
			return false
		}
		if size[ra] < size[rb] {
			ra, rb = rb, ra
		}
		parent[rb] = ra
		size[ra] += size[rb]
		return true
	}

	components := n
	limit := len(edges)
	if maxConnections > 0 && maxConnections < limit {
		limit = maxConnections
	}

	res := dsuResult{lastMergeI: -1, lastMergeJ: -1}

	for k := 0; k < limit; k++ {
		e := edges[k]
		if union(e.i, e.j) {
			components--
			res.lastMergeI = e.i
			res.lastMergeJ = e.j
			if stopWhenSingle && components == 1 {
				break
			}
		}
	}

	componentSizes := make([]int64, 0, n)
	for i := range n {
		if parent[i] == i {
			componentSizes = append(componentSizes, size[i])
		}
	}
	res.componentSizes = componentSizes
	return res
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 8}
}

func (d Day8) Part1(input string) (string, error) {
	points, err := parsePoints(input)
	if err != nil {
		return "", err
	}

	n := len(points)

	edges, err := buildEdges(points, 1000)
	if err != nil {
		return "", err
	}
	res := processEdges(n, edges, 1000, false)

	componentSizes := res.componentSizes
	if len(componentSizes) < 3 {
		return "", errors.New("fewer than three circuits remain after the connections")
	}

	sort.Slice(componentSizes, func(a, b int) bool { return componentSizes[a] > componentSizes[b] })

	product, err := multiplyCircuitAnswer(componentSizes[0], componentSizes[1])
	if err != nil {
		return "", err
	}
	product, err = multiplyCircuitAnswer(product, componentSizes[2])
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(product, 10), nil
}

func (d Day8) Part2(input string) (string, error) {
	points, err := parsePoints(input)
	if err != nil {
		return "", err
	}

	n := len(points)

	edges, err := buildEdges(points, 0)
	if err != nil {
		return "", err
	}
	res := processEdges(n, edges, 0, true)
	if res.lastMergeI == -1 || res.lastMergeJ == -1 {
		return "0", nil
	}

	product, err := multiplyCircuitAnswer(points[res.lastMergeI].x, points[res.lastMergeJ].x)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(product, 10), nil
}

func init() {
	solve.Register(Day8{})
}
