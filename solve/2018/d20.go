package solve2018

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day20 struct{}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 20}
}

func parseInputToGraph(input string) map[[2]int][][2]int {
	input = strings.Trim(input, "^$")
	graph := make(map[[2]int][][2]int)
	stack := [][2]int{}
	current := [2]int{0, 0}

	directions := map[byte][2]int{
		'N': {0, -1},
		'S': {0, 1},
		'E': {1, 0},
		'W': {-1, 0},
	}

	for i := 0; i < len(input); i++ {
		switch input[i] {
		case 'N', 'S', 'E', 'W':
			delta := directions[input[i]]
			next := [2]int{current[0] + delta[0], current[1] + delta[1]}
			graph[current] = append(graph[current], next)
			graph[next] = append(graph[next], current)
			current = next
		case '(':
			stack = append(stack, current)
		case ')':
			current = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case '|':
			current = stack[len(stack)-1]
		}
	}

	return graph
}

func bfsDistances(graph map[[2]int][][2]int, start [2]int) map[[2]int]int {
	distances := make(map[[2]int]int)
	queue := [][2]int{start}
	distances[start] = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, neighbor := range graph[current] {
			if _, visited := distances[neighbor]; !visited {
				distances[neighbor] = distances[current] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return distances
}

func (d Day20) Part1(input string) (string, error) {
	graph := parseInputToGraph(input)
	start := [2]int{0, 0}
	distances := bfsDistances(graph, start)

	maxDistance := 0
	for _, distance := range distances {
		if distance > maxDistance {
			maxDistance = distance
		}
	}

	return strconv.Itoa(maxDistance), nil
}

func (d Day20) Part2(input string) (string, error) {
	graph := parseInputToGraph(input)
	start := [2]int{0, 0}
	distances := bfsDistances(graph, start)

	count := 0
	for _, distance := range distances {
		if distance >= 1000 {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

func init() {
	solve.Register(Day20{})
}
