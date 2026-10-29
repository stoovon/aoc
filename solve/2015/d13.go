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

type Day13 struct{}

var effectRegex = regexp.MustCompile(`(\w+) .* (lose|gain) (\d+) .* (\w+)\.`)

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 13}
}

type seatingScore struct {
	happiness int
	reachable bool
}

func findMaxHappiness(mapping map[string]map[string]int, includeSelf bool) (int, error) {
	people := slices.Sorted(maps.Keys(mapping))
	n := len(people)
	if n == 0 {
		return 0, fmt.Errorf("seating plan has no guests")
	}
	if includeSelf {
		n++
	}
	if n == 1 {
		return 2 * mapping[people[0]][people[0]], nil
	}
	others := n - 1
	const maxEntries = 1 << 20
	if others >= 20 || (1<<others) > maxEntries/others {
		return 0, fmt.Errorf("seating plan exceeds the subset table budget")
	}
	weights := make([][]int, n)
	for i := range n {
		weights[i] = make([]int, n)
	}
	for a, person := range people {
		for b, neighbour := range people {
			weights[a][b] = mapping[person][neighbour] + mapping[neighbour][person]
		}
	}
	// Fix person 0 to remove rotations; subset bits describe the remaining guests.
	subsets := 1 << others
	table := make([]seatingScore, subsets*others)
	for last := range others {
		table[(1<<last)*others+last] = seatingScore{happiness: weights[0][last+1], reachable: true}
	}
	for seated := 1; seated < subsets; seated++ {
		for last := range others {
			current := table[seated*others+last]
			if !current.reachable {
				continue
			}
			for next := range others {
				if seated&(1<<next) != 0 {
					continue
				}
				entry := &table[(seated|1<<next)*others+next]
				happiness := current.happiness + weights[last+1][next+1]
				if !entry.reachable || happiness > entry.happiness {
					*entry = seatingScore{happiness: happiness, reachable: true}
				}
			}
		}
	}
	var answer seatingScore
	for last, entry := range table[(subsets-1)*others:] {
		happiness := entry.happiness + weights[last+1][0]
		if entry.reachable && (!answer.reachable || happiness > answer.happiness) {
			answer = seatingScore{happiness: happiness, reachable: true}
		}
	}
	return answer.happiness, nil
}

func (d Day13) parse(input string) map[string]map[string]int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	mapping := make(map[string]map[string]int)

	for _, line := range lines {
		matches := effectRegex.FindStringSubmatch(line)
		if len(matches) == 5 {
			a, change, n, b := matches[1], matches[2], matches[3], matches[4]
			value, _ := strconv.Atoi(n)
			if change == "lose" {
				value = -value
			}
			if mapping[a] == nil {
				mapping[a] = make(map[string]int)
			}
			mapping[a][b] = value
		}
	}

	return mapping
}

func (d Day13) Part1(input string) (string, error) {
	happiness, err := findMaxHappiness(d.parse(input), false)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(happiness), nil
}

func (d Day13) Part2(input string) (string, error) {
	happiness, err := findMaxHappiness(d.parse(input), true)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(happiness), nil
}

func init() {
	solve.Register(Day13{})
}
