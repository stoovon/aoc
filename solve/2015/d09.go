package solve2015

import (
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day9 struct {
}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 9}
}

var costRegex = regexp.MustCompile(`(.+) to (.+) = (\d+)`)

func (d Day9) parseInput(input string) map[string]map[string]int {
	graph := make(map[string]map[string]int)

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		matches := costRegex.FindStringSubmatch(line)
		if len(matches) != 4 {
			continue
		}

		node1, node2 := matches[1], matches[2]
		distance, _ := strconv.Atoi(matches[3])

		if graph[node1] == nil {
			graph[node1] = make(map[string]int)
		}
		if graph[node2] == nil {
			graph[node2] = make(map[string]int)
		}

		graph[node1][node2] = distance
		graph[node2][node1] = distance
	}

	return graph
}

func generatePermutations(nodes []string) [][]string {
	var helper func([]string, int)
	var result [][]string

	helper = func(arr []string, n int) {
		if n == 1 {
			perm := make([]string, len(arr))
			copy(perm, arr)
			result = append(result, perm)
			return
		}

		for i := 0; i < n; i++ {
			helper(arr, n-1)
			if n%2 == 1 {
				arr[0], arr[n-1] = arr[n-1], arr[0]
			} else {
				arr[i], arr[n-1] = arr[n-1], arr[i]
			}
		}
	}

	helper(nodes, len(nodes))
	return result
}

func calculateDistances(graph map[string]map[string]int, permutations [][]string) []int {
	var distances []int

	for _, path := range permutations {
		totalDistance := 0
		for i := 0; i < len(path)-1; i++ {
			totalDistance += graph[path[i]][path[i+1]]
		}
		distances = append(distances, totalDistance)
	}

	return distances
}

func (d Day9) getAllDistances(input string) []int {
	graph := d.parseInput(input)

	nodes := make([]string, 0, len(graph))
	for node := range graph {
		nodes = append(nodes, node)
	}

	permutations := generatePermutations(nodes)
	return calculateDistances(graph, permutations)
}

func (d Day9) Part1(input string) (string, error) {
	return strconv.Itoa(maths.Min(d.getAllDistances(input)...)), nil
}

func (d Day9) Part2(input string) (string, error) {
	return strconv.Itoa(maths.Max(d.getAllDistances(input)...)), nil
}

func init() {
	solve.Register(Day9{})
}
