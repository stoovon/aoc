package solve2015

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day9 struct{}

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

type routeRange struct {
	shortest, longest int
	reachable         bool
}

func routeExtremes(graph map[string]map[string]int) (int, int, error) {
	nodes := slices.Sorted(maps.Keys(graph))
	n := len(nodes)
	if n == 0 {
		return 0, 0, fmt.Errorf("route graph has no cities")
	}
	// Subset DP remains exponential; bound its table before shifting or allocating.
	const maxEntries = 1 << 20
	if n >= 20 || (1<<n) > maxEntries/n {
		return 0, 0, fmt.Errorf("route graph exceeds the subset table budget")
	}
	subsets := 1 << n
	distances := make([][]int, n)
	for from, name := range nodes {
		distances[from] = make([]int, n)
		for to, other := range nodes {
			distance, connected := graph[name][other]
			if !connected {
				distance = -1
			}
			distances[from][to] = distance
		}
	}
	table := make([]routeRange, subsets*n)
	for city := range n {
		table[(1<<city)*n+city].reachable = true
	}
	for visited := 1; visited < subsets; visited++ {
		for last := range n {
			current := table[visited*n+last]
			if !current.reachable {
				continue
			}
			for next, distance := range distances[last] {
				if visited&(1<<next) != 0 || distance < 0 {
					continue
				}
				entry := &table[(visited|1<<next)*n+next]
				shortest := current.shortest + distance
				longest := current.longest + distance
				if !entry.reachable {
					*entry = routeRange{shortest: shortest, longest: longest, reachable: true}
				} else {
					entry.shortest = min(entry.shortest, shortest)
					entry.longest = max(entry.longest, longest)
				}
			}
		}
	}
	var answer routeRange
	for _, entry := range table[(subsets-1)*n:] {
		if !entry.reachable {
			continue
		}
		if !answer.reachable {
			answer = entry
		} else {
			answer.shortest = min(answer.shortest, entry.shortest)
			answer.longest = max(answer.longest, entry.longest)
		}
	}
	if !answer.reachable {
		return 0, 0, fmt.Errorf("no route visits every city exactly once")
	}
	return answer.shortest, answer.longest, nil
}

func (d Day9) Part1(input string) (string, error) {
	shortest, _, err := routeExtremes(d.parseInput(input))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(shortest), nil
}

func (d Day9) Part2(input string) (string, error) {
	_, longest, err := routeExtremes(d.parseInput(input))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(longest), nil
}

func init() {
	solve.Register(Day9{})
}
