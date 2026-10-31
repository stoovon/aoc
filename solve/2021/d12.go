package solve2021

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day12 struct{}

func (d Day12) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2021, Day: 12} }

type caveGraph struct {
	edges             [][]int
	small             []int
	start, end, bytes int
}

type cavePathState struct {
	cave       int
	visited    string
	usedDouble bool
}

func parseCaves(input string) (caveGraph, error) {
	g := caveGraph{start: -1, end: -1}
	ids := map[string]int{}
	smallCount := 0
	intern := func(name string) int {
		if id, ok := ids[name]; ok {
			return id
		}
		id := len(g.edges)
		ids[name] = id
		g.edges = append(g.edges, nil)
		bit := -1
		if name == strings.ToLower(name) {
			bit = smallCount
			smallCount++
		}
		g.small = append(g.small, bit)
		if name == "start" {
			g.start = id
		}
		if name == "end" {
			g.end = id
		}
		return id
	}
	seen := map[[2]int]bool{}
	for line, record := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {
		record = strings.TrimSuffix(record, "\r")
		names := strings.Split(record, "-")
		if len(names) != 2 {
			return g, fmt.Errorf("line %d: expected cave-cave", line+1)
		}
		for _, name := range names {
			lower, upper := true, true
			for _, c := range name {
				lower = lower && c >= 'a' && c <= 'z'
				upper = upper && c >= 'A' && c <= 'Z'
			}
			if name == "" || (!lower && !upper) {
				return g, fmt.Errorf("line %d: invalid cave name", line+1)
			}
		}
		a, b := intern(names[0]), intern(names[1])
		for _, edge := range [][2]int{{a, b}, {b, a}} {
			if !seen[edge] {
				seen[edge] = true
				g.edges[edge[0]] = append(g.edges[edge[0]], edge[1])
			}
		}
	}
	if g.start < 0 || g.end < 0 {
		return g, fmt.Errorf("missing start or end cave")
	}
	g.bytes = (smallCount + 7) / 8
	return g, nil
}

func (g caveGraph) count(usedDouble bool) (uint64, error) {
	visited := make([]byte, g.bytes)
	visited[g.small[g.start]/8] |= 1 << (g.small[g.start] % 8)
	initial := cavePathState{g.start, string(visited), usedDouble}
	states := []cavePathState{initial}
	ids := map[cavePathState]int{initial: 0}
	edges, reverse := [][]int{nil}, [][]int{nil}
	for i := 0; i < len(states); i++ {
		cur := states[i]
		if cur.cave == g.end {
			continue
		}
		for _, cave := range g.edges[cur.cave] {
			if cave == g.start {
				continue
			}
			next := cavePathState{cave, cur.visited, cur.usedDouble}
			if bit := g.small[cave]; bit >= 0 {
				if cur.visited[bit/8]&(1<<(bit%8)) != 0 {
					if cur.usedDouble {
						continue
					}
					next.usedDouble = true
				} else {
					bits := []byte(cur.visited)
					bits[bit/8] |= 1 << (bit % 8)
					next.visited = string(bits)
				}
			}
			id, ok := ids[next]
			if !ok {
				id = len(states)
				ids[next] = id
				states = append(states, next)
				edges = append(edges, nil)
				reverse = append(reverse, nil)
			}
			edges[i] = append(edges[i], id)
			reverse[id] = append(reverse[id], i)
		}
	}
	productive := make([]bool, len(states))
	queue := []int{}
	for i, s := range states {
		if s.cave == g.end {
			productive[i] = true
			queue = append(queue, i)
		}
	}
	for head := 0; head < len(queue); head++ {
		for _, p := range reverse[queue[head]] {
			if !productive[p] {
				productive[p] = true
				queue = append(queue, p)
			}
		}
	}
	remaining := make([]int, len(states))
	counts := make([]uint64, len(states))
	ready := []int{}
	for i, s := range states {
		if !productive[i] {
			continue
		}
		for _, to := range edges[i] {
			if productive[to] {
				remaining[i]++
			}
		}
		if s.cave == g.end {
			counts[i] = 1
		}
		if remaining[i] == 0 {
			ready = append(ready, i)
		}
	}
	processed := 0
	for head := 0; head < len(ready); head++ {
		id := ready[head]
		processed++
		for _, p := range reverse[id] {
			if counts[p] > math.MaxUint64-counts[id] {
				return 0, fmt.Errorf("path count exceeds uint64")
			}
			counts[p] += counts[id]
			remaining[p]--
			if remaining[p] == 0 {
				ready = append(ready, p)
			}
		}
	}
	if processed != len(queue) {
		return 0, fmt.Errorf("infinitely many paths: repeatable state can reach end")
	}
	return counts[0], nil
}

func cavePaths(input string, usedDouble bool) (string, error) {
	g, err := parseCaves(input)
	if err != nil {
		return "", err
	}
	count, err := g.count(usedDouble)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(count, 10), nil
}
func (d Day12) Part1(input string) (string, error) { return cavePaths(input, true) }
func (d Day12) Part2(input string) (string, error) { return cavePaths(input, false) }
func init()                                        { solve.Register(Day12{}) }
