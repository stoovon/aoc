package solve2015

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day24 struct {
}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 24}
}

func balance(weights []int, numGroups int) int {
	groupSize := sum(weights) / numGroups
	minQE := math.MaxInt64

	for i := 1; i <= len(weights); i++ {
		combs := combinations(weights, i)
		for _, comb := range combs {
			if sum(comb) == groupSize {
				qe := quantumEntanglement(comb)
				if qe < minQE {
					minQE = qe
				}
			}
		}
		if minQE != math.MaxInt64 {
			return minQE
		}
	}
	return minQE
}

func combinations(weights []int, n int) [][]int {
	if n == 0 {
		return [][]int{{}}
	}
	if len(weights) < n {
		return nil
	}

	withFirst := combinations(weights[1:], n-1)
	for i := range withFirst {
		withFirst[i] = append([]int{weights[0]}, withFirst[i]...)
	}

	withoutFirst := combinations(weights[1:], n)
	return append(withFirst, withoutFirst...)
}

func quantumEntanglement(group []int) int {
	qe := 1
	for _, weight := range group {
		qe *= weight
	}
	return qe
}

func sum(weights []int) int {
	total := 0
	for _, weight := range weights {
		total += weight
	}
	return total
}

func (d Day24) parseInput(input string) []int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	weights := make([]int, len(lines))
	for i, line := range lines {
		weights[i], _ = strconv.Atoi(line)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(weights))) // Sort in descending order for efficiency
	return weights
}

func (d Day24) Part1(input string) (string, error) {
	weights := d.parseInput(input)
	result := balance(weights, 3)
	return strconv.Itoa(result), nil
}

func (d Day24) Part2(input string) (string, error) {
	weights := d.parseInput(input)
	result := balance(weights, 4)
	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day24{})
}
