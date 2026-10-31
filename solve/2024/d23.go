package solve2024

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day23 struct{}

func (d Day23) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2024, Day: 23} }

type lanGraph struct {
	names    []string
	adjacent [][]uint64
}

func parseLAN(input string) (lanGraph, error) {
	var graph lanGraph
	edges := [][2]string{}
	names := map[string]bool{}
	for line, record := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {
		record = strings.TrimSuffix(record, "\r")
		pair := strings.Split(record, "-")
		if len(pair) != 2 {
			return graph, fmt.Errorf("line %d: expected computer-computer", line+1)
		}
		for _, name := range pair {
			if len(name) != 2 || name[0] < 'a' || name[0] > 'z' || name[1] < 'a' || name[1] > 'z' {
				return graph, fmt.Errorf("line %d: expected two lowercase letters per computer", line+1)
			}
		}
		if pair[0] == pair[1] {
			return graph, fmt.Errorf("line %d: self connection", line+1)
		}
		edges = append(edges, [2]string{pair[0], pair[1]})
		names[pair[0]] = true
		names[pair[1]] = true
	}
	for name := range names {
		graph.names = append(graph.names, name)
	}
	sort.Strings(graph.names)
	ids := map[string]int{}
	graph.adjacent = make([][]uint64, len(graph.names))
	for id, name := range graph.names {
		ids[name] = id
		graph.adjacent[id] = make([]uint64, (len(graph.names)+63)/64)
	}
	for _, edge := range edges {
		a, b := ids[edge[0]], ids[edge[1]]
		graph.adjacent[a][b/64] |= 1 << (b % 64)
		graph.adjacent[b][a/64] |= 1 << (a % 64)
	}
	return graph, nil
}

func (g lanGraph) connected(a, b int) bool { return g.adjacent[a][b/64]&(1<<(b%64)) != 0 }

func (g lanGraph) triangles() uint64 {
	var count uint64
	for a := range g.names {
		for b := a + 1; b < len(g.names); b++ {
			if !g.connected(a, b) {
				continue
			}
			for c := b + 1; c < len(g.names); c++ {
				if g.connected(a, c) && g.connected(b, c) && (g.names[a][0] == 't' || g.names[b][0] == 't' || g.names[c][0] == 't') {
					count++
				}
			}
		}
	}
	return count
}

func (g lanGraph) largestClique() []int {
	candidates := make([]int, len(g.names))
	for i := range candidates {
		candidates[i] = i
	}
	best, path := []int{}, []int{}
	var search func([]int)
	search = func(remaining []int) {
		if len(path) > len(best) {
			best = append([]int(nil), path...)
		}
		for i, vertex := range remaining {
			if len(path)+len(remaining)-i <= len(best) {
				return
			}
			next := make([]int, 0, len(remaining)-i-1)
			for _, candidate := range remaining[i+1:] {
				if g.connected(vertex, candidate) {
					next = append(next, candidate)
				}
			}
			path = append(path, vertex)
			search(next)
			path = path[:len(path)-1]
		}
	}
	search(candidates)
	return best
}

func (d Day23) Part1(input string) (string, error) {
	graph, err := parseLAN(input)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(graph.triangles(), 10), nil
}
func (d Day23) Part2(input string) (string, error) {
	graph, err := parseLAN(input)
	if err != nil {
		return "", err
	}
	clique := graph.largestClique()
	names := make([]string, len(clique))
	for i, id := range clique {
		names[i] = graph.names[id]
	}
	return strings.Join(names, ","), nil
}
func init() { solve.Register(Day23{}) }
