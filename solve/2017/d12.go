package solve2017

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day12 struct {
}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 12}
}

func parseGraph(input string) map[int][]int {
	graph := make(map[int][]int)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "<->")
		from, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		neighbors := strings.Split(parts[1], ",")
		for _, n := range neighbors {
			to, _ := strconv.Atoi(strings.TrimSpace(n))
			graph[from] = append(graph[from], to)
		}
	}
	return graph
}

func groupSize(graph map[int][]int, start int, visited map[int]bool) int {
	queue := []int{start}
	count := 0
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if visited[node] {
			continue
		}
		visited[node] = true
		count++
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}
	return count
}

func (d Day12) Part1(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("empty input")
	}
	graph := parseGraph(input)
	visited := make(map[int]bool)
	size := groupSize(graph, 0, visited)
	return strconv.Itoa(size), nil
}

func (d Day12) Part2(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("empty input")
	}
	graph := parseGraph(input)
	visited := make(map[int]bool)
	groups := 0
	for node := range graph {
		if !visited[node] {
			groupSize(graph, node, visited)
			groups++
		}
	}
	return strconv.Itoa(groups), nil
}

func init() {
	solve.Register(Day12{})
}
