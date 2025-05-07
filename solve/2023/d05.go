package solve2023

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 5}
}

func (d Day5) translate(mapping [][]int, pairs [][2]int) [][2]int {
	var result [][2]int
	for _, pair := range pairs {
		start, end := pair[0], pair[1]
		for _, m := range mapping {
			a1, a2, d := m[0], m[1], m[2]
			result = append(result, [2]int{start, maths.Min(a1, end)})
			result = append(result, [2]int{max(a1, start) + d, maths.Min(a2, end) + d})
			start = maths.Max(start, maths.Min(a2, end))
		}
		result = append(result, [2]int{start, end})
	}
	return result
}

func (d Day5) solve(mappings [][][]int, seed [][2]int) int {
	for _, mapping := range mappings {
		newSeed := [][2]int{}
		for _, pair := range d.translate(mapping, seed) {
			if pair[0] < pair[1] {
				newSeed = append(newSeed, pair)
			}
		}
		seed = newSeed
	}
	minValue := math.MaxInt
	for _, pair := range seed {
		if pair[0] < minValue {
			minValue = pair[0]
		}
	}
	return minValue
}

func (d Day5) parse(input []string) [][][]int {
	var mappings [][][]int
	i := 3
	for i < len(input) {
		var current [][]int
		for i < len(input) && input[i] != "" {
			parts := strings.Fields(input[i])
			s2, _ := strconv.Atoi(parts[0])
			s1, _ := strconv.Atoi(parts[1])
			length, _ := strconv.Atoi(parts[2])
			current = append(current, []int{s1, s1 + length, s2 - s1})
			i++
		}
		sort.Slice(current, func(i, j int) bool {
			return current[i][0] < current[j][0]
		})
		mappings = append(mappings, current)
		i += 2
	}
	return mappings
}

func (d Day5) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	seeds := parseSeeds(lines[0])
	mappings := d.parse(lines)
	seed := make([][2]int, len(seeds))
	for i, x := range seeds {
		seed[i] = [2]int{x, x + 1}
	}
	answer := d.solve(mappings, seed)
	return strconv.Itoa(answer), nil
}

func (d Day5) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	seeds := parseSeeds(lines[0])
	mappings := d.parse(lines)

	// Ensure seeds have an even number of elements
	if len(seeds)%2 != 0 {
		return "", errors.New("invalid input: seeds must have an even number of elements")
	}

	seed := make([][2]int, len(seeds)/2)
	for i := 0; i < len(seeds); i += 2 {
		seed[i/2] = [2]int{seeds[i], seeds[i] + seeds[i+1]}
	}
	answer := d.solve(mappings, seed)
	return strconv.Itoa(answer), nil
}

func parseSeeds(line string) []int {
	parts := strings.Split(strings.Split(line, ":")[1], " ")
	seeds := make([]int, len(parts)-1)
	for i, part := range parts {
		if part == "" {
			continue
		}
		seeds[i-1], _ = strconv.Atoi(part)
	}
	return seeds
}

func init() {
	solve.Register(Day5{})
}
