package solve2023

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day25 struct {
}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 25}
}

func (d Day25) Part1(input string) (string, error) {
	// Build the graph
	G := map[string]map[string]struct{}{}
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		parts := strings.Fields(strings.ReplaceAll(line, ":", ""))
		u := parts[0]
		if _, ok := G[u]; !ok {
			G[u] = map[string]struct{}{}
		}
		for _, v := range parts[1:] {
			G[u][v] = struct{}{}
			if _, ok := G[v]; !ok {
				G[v] = map[string]struct{}{}
			}
			G[v][u] = struct{}{}
		}
	}

	// Initialize S as the set of all nodes
	S := map[string]struct{}{}
	for k := range G {
		S[k] = struct{}{}
	}

	// Helper to count external connections
	count := func(v string) int {
		c := 0
		for n := range G[v] {
			if _, ok := S[n]; !ok {
				c++
			}
		}
		return c
	}

	// Remove nodes until only 3 edges cross the cut
	for {
		sum := 0
		for v := range S {
			sum += count(v)
		}
		if sum == 3 {
			break
		}
		// Find node with max external connections
		var maxV string
		maxC := -1
		for v := range S {
			c := count(v)
			if c > maxC {
				maxC = c
				maxV = v
			}
		}
		delete(S, maxV)
	}

	// Compute the product of the sizes of the two partitions
	other := map[string]struct{}{}
	for k := range G {
		if _, ok := S[k]; !ok {
			other[k] = struct{}{}
		}
	}
	return strconv.Itoa(len(S) * len(other)), nil
}

func (d Day25) Part2(_ string) (string, error) {
	return "", nil
}

func init() {
	solve.Register(Day25{})
}
