package solve2025

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day11 struct {
}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 11}
}

func parseGraph(input string) map[string][]string {
	graph := make(map[string][]string)
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		node := strings.TrimSpace(parts[0])
		children := strings.Fields(strings.TrimSpace(parts[1]))
		graph[node] = children
	}
	return graph
}

// countPaths returns the number of distinct paths from start to end in the DAG
// that visit ALL nodes in the required slice (if any).
// With required=nil, counts all paths unconditionally.
func countPaths(graph map[string][]string, start, end string, required []string) int64 {
	reqIndex := make(map[string]int)
	for i, r := range required {
		reqIndex[r] = i
	}
	fullMask := (1 << len(required)) - 1

	// memo[node][mask] = number of paths from node to end with current visited mask
	memo := make(map[string]map[int]int64)

	var dfs func(node string, mask int) int64
	dfs = func(node string, mask int) int64 {
		if idx, ok := reqIndex[node]; ok {
			mask |= (1 << idx)
		}

		if node == end {
			if mask == fullMask {
				return 1
			}
			return 0
		}

		if memo[node] == nil {
			memo[node] = make(map[int]int64)
		}
		if v, ok := memo[node][mask]; ok {
			return v
		}

		var total int64
		for _, child := range graph[node] {
			total += dfs(child, mask)
		}
		memo[node][mask] = total
		return total
	}

	return dfs(start, 0)
}

func (d Day11) Part1(input string) (string, error) {
	graph := parseGraph(input)
	count := countPaths(graph, "you", "out", nil)
	return strconv.FormatInt(count, 10), nil
}

func (d Day11) Part2(input string) (string, error) {
	graph := parseGraph(input)
	count := countPaths(graph, "svr", "out", []string{"dac", "fft"})
	return strconv.FormatInt(count, 10), nil
}

func init() {
	solve.Register(Day11{})
}
