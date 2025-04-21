package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day16 struct {
}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 16}
}

type candidate map[string]int

func (d Day16) parseInput(data string) map[int]candidate {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	candidates := make(map[int]candidate)

	for _, line := range lines {
		parts := strings.Fields(line)
		id, _ := strconv.Atoi(strings.TrimSuffix(parts[1], ":"))
		candidates[id] = make(candidate)

		for i := 2; i < len(parts); i += 2 {
			prop := strings.TrimSuffix(parts[i], ":")
			value, _ := strconv.Atoi(strings.TrimSuffix(parts[i+1], ","))
			candidates[id][prop] = value
		}
	}

	return candidates
}

func filterCandidates(candidates map[int]candidate, evidence map[string]int, predicate func(candidate, map[string]int) bool) int {
	for id, candidate := range candidates {
		if predicate(candidate, evidence) {
			return id
		}
	}
	return 0
}

func simplePredicate(currentCandidate candidate, evidence map[string]int) bool {
	for prop, value := range evidence {
		if setProp, ok := currentCandidate[prop]; ok {
			if setProp != value {
				return false
			}
		}
	}
	return true
}

func (d Day16) Part1(input string) (string, error) {
	candidates := d.parseInput(input)
	evidence := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}
	return strconv.Itoa(filterCandidates(candidates, evidence, simplePredicate)), nil
}

func rangePredicate(currentCandidate candidate, evidence map[string]int) bool {
	for prop, value := range evidence {
		setProp, ok := currentCandidate[prop]

		if !ok {
			continue
		}

		switch prop {
		case "cats", "trees":
			if setProp <= value {
				return false
			}
		case "pomeranians", "goldfish":
			if setProp >= value {
				return false
			}
		default:
			if setProp != value {
				return false
			}
		}
	}
	return true
}

func (d Day16) Part2(input string) (string, error) {
	candidates := d.parseInput(input)
	evidence := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}
	return strconv.Itoa(filterCandidates(candidates, evidence, rangePredicate)), nil
}

func init() {
	solve.Register(Day16{})
}
