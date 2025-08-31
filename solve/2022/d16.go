package solve2022

import (
	"aoc/solve"
	"fmt"
	"math"
	"strings"
)

type Day16 struct{}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 16}
}

func (d Day16) Part1(input string) (string, error) {
	// Parse input
	valves, tunnels := d.parseInput(input)

	// Compute shortest paths using Floyd-Warshall
	costs := floydWarshall(valves, tunnels)

	// Filter relevant valves
	relevantValves := filterRelevantValves(valves)

	// Perform DFS to maximize pressure
	maxPressure := dfs("AA", 30, 0, relevantValves, costs, map[string]bool{})

	return fmt.Sprintf("%d", maxPressure), nil
}

func (d Day16) Part2(input string) (string, error) {
	// Parse input
	valves, tunnels := d.parseInput(input)

	// Compute shortest paths using Floyd-Warshall
	costs := floydWarshall(valves, tunnels)

	// Filter relevant valves
	relevantValves := filterRelevantValves(valves)

	// First worker (you)
	maxPressure1, path1 := dfsWithPath("AA", 26, 0, relevantValves, costs, map[string]bool{}, []string{})

	// Mark valves in path1 as opened
	for _, valve := range path1 {
		relevantValves[valve] = 0
	}

	// Second worker (elephant)
	maxPressure2, _ := dfsWithPath("AA", 26, 0, relevantValves, costs, map[string]bool{}, []string{})

	// Combine results
	totalPressure := maxPressure1 + maxPressure2

	return fmt.Sprintf("%d", totalPressure), nil
}

func init() {
	solve.Register(Day16{})
}

func (d Day16) parseInput(input string) (map[string]int, map[string][]string) {
	valves := make(map[string]int)
	tunnels := make(map[string][]string)
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		// Extract the valve and flow rate
		parts := strings.Split(line, ";")
		valveInfo := strings.Fields(parts[0])
		valve := valveInfo[1]
		flow := 0
		fmt.Sscanf(valveInfo[4], "rate=%d", &flow)
		valves[valve] = flow

		// Extract the connections
		connections := []string{}
		if len(parts) > 1 {
			connectionsPart := strings.TrimPrefix(parts[1], " tunnels lead to valves ")
			connectionsPart = strings.TrimPrefix(connectionsPart, " tunnel leads to valve ")
			connections = strings.Split(connectionsPart, ", ")
		}
		tunnels[valve] = connections
	}
	return valves, tunnels
}

func floydWarshall(valves map[string]int, tunnels map[string][]string) map[string]map[string]int {
	costs := make(map[string]map[string]int)
	for v := range valves {
		costs[v] = make(map[string]int)
		for u := range valves {
			if v == u {
				costs[v][u] = 0
			} else {
				costs[v][u] = math.MaxInt32
			}
		}
	}
	for v, neighbors := range tunnels {
		for _, u := range neighbors {
			costs[v][u] = 1
		}
	}
	for k := range valves {
		for i := range valves {
			for j := range valves {
				if costs[i][j] > costs[i][k]+costs[k][j] {
					costs[i][j] = costs[i][k] + costs[k][j]
				}
			}
		}
	}
	return costs
}

func filterRelevantValves(valves map[string]int) map[string]int {
	relevant := make(map[string]int)
	for valve, flow := range valves {
		if flow > 0 || valve == "AA" {
			relevant[valve] = flow
		}
	}
	return relevant
}

func dfs(current string, timeLeft, pressure int, valves map[string]int, costs map[string]map[string]int, opened map[string]bool) int {
	if timeLeft <= 0 {
		return pressure
	}
	maxPressure := pressure
	for valve, flow := range valves {
		if !opened[valve] && flow > 0 {
			timeToOpen := costs[current][valve] + 1
			if timeLeft >= timeToOpen {
				opened[valve] = true
				totalPressure := dfs(valve, timeLeft-timeToOpen, pressure+(timeLeft-timeToOpen)*flow, valves, costs, opened)
				if totalPressure > maxPressure {
					maxPressure = totalPressure
				}
				opened[valve] = false
			}
		}
	}
	return maxPressure
}

func dfsWithPath(current string, timeLeft, pressure int, valves map[string]int, costs map[string]map[string]int, opened map[string]bool, path []string) (int, []string) {
	if timeLeft <= 0 {
		return pressure, path
	}
	maxPressure := pressure
	bestPath := append([]string{}, path...)
	for valve, flow := range valves {
		if !opened[valve] && flow > 0 {
			timeToOpen := costs[current][valve] + 1
			if timeLeft >= timeToOpen {
				opened[valve] = true
				newPath := append(path, valve)
				totalPressure, resultPath := dfsWithPath(valve, timeLeft-timeToOpen, pressure+(timeLeft-timeToOpen)*flow, valves, costs, opened, newPath)
				if totalPressure > maxPressure {
					maxPressure = totalPressure
					bestPath = resultPath
				}
				opened[valve] = false
			}
		}
	}
	return maxPressure, bestPath
}
