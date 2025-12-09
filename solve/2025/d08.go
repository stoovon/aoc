package solve2025

import (
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

	return points, nil
}

func buildEdges(points []vec3) []edge {
	n := len(points)
	edges := make([]edge, 0, n*(n-1)/2)
	for i := range n {
		for j := i + 1; j < n; j++ {
			dx := points[i].x - points[j].x
			dy := points[i].y - points[j].y
			dz := points[i].z - points[j].z
			d2 := dx*dx + dy*dy + dz*dz
			edges = append(edges, edge{i: i, j: j, d2: d2})
		}
	}
	return edges
}

func processEdges(n int, edges []edge, maxConnections int, stopWhenSingle bool) dsuResult {
	sort.Slice(edges, func(a, b int) bool {
		if edges[a].d2 == edges[b].d2 {
			if edges[a].i == edges[b].i {
				return edges[a].j < edges[b].j
			}
			return edges[a].i < edges[b].i
		}
		return edges[a].d2 < edges[b].d2
	})

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
	if n == 0 {
		return "0", nil
	}

	edges := buildEdges(points)
	res := processEdges(n, edges, 1000, false)

	componentSizes := res.componentSizes
	if len(componentSizes) < 3 {
		for len(componentSizes) < 3 {
			componentSizes = append(componentSizes, 1)
		}
	}

	sort.Slice(componentSizes, func(a, b int) bool { return componentSizes[a] > componentSizes[b] })

	product := componentSizes[0] * componentSizes[1] * componentSizes[2]
	return strconv.FormatInt(product, 10), nil
}

func (d Day8) Part2(input string) (string, error) {
	points, err := parsePoints(input)
	if err != nil {
		return "", err
	}

	n := len(points)
	if n == 0 {
		return "0", nil
	}

	edges := buildEdges(points)
	res := processEdges(n, edges, 0, true)
	if res.lastMergeI == -1 || res.lastMergeJ == -1 {
		return "0", nil
	}

	product := points[res.lastMergeI].x * points[res.lastMergeJ].x
	return strconv.FormatInt(product, 10), nil
}

func init() {
	solve.Register(Day8{})
}
