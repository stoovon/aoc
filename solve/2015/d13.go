package solve2015

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day13 struct {
}

var effectRegex = regexp.MustCompile(`(\w+) .* (lose|gain) (\d+) .* (\w+)\.`)

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 13}
}

func happiness(combo []string, mapping map[string]map[string]int) int {
	total := 0
	n := len(combo)
	for i := 0; i < n; i++ {
		a, b := combo[i], combo[(i+1)%n]
		total += mapping[a][b] + mapping[b][a]
	}
	return total
}

func permutations(arr []string) [][]string {
	var helper func([]string, int)
	var res [][]string

	helper = func(arr []string, n int) {
		if n == 1 {
			tmp := make([]string, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
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

	helper(arr, len(arr))
	return res
}

func findMaxHappiness(mapping map[string]map[string]int) int {
	var people []string
	for person := range mapping {
		people = append(people, person)
	}

	maxHappiness := 0
	for _, combo := range permutations(people) {
		h := happiness(combo, mapping)
		if h > maxHappiness {
			maxHappiness = h
		}
	}
	return maxHappiness
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
	mapping := d.parse(input)

	answerA := findMaxHappiness(mapping)

	return fmt.Sprintf("%d", answerA), nil
}

func (d Day13) Part2(input string) (string, error) {
	mapping := d.parse(input)

	for person := range mapping {
		if mapping["me"] == nil {
			mapping["me"] = make(map[string]int)
		}
		mapping[person]["me"] = 0
		mapping["me"][person] = 0
	}

	answerB := findMaxHappiness(mapping)

	return strconv.Itoa(answerB), nil
}

func init() {
	solve.Register(Day13{})
}
