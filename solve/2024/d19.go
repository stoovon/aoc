package solve2024

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day19 struct {
}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 19}
}

func (d Day19) permutations(cache map[string]int, split []string, design string) (n int) {
	if n, ok := cache[design]; ok {
		return n
	}
	defer func() { cache[design] = n }()

	if design == "" {
		return 1
	}
	for _, s := range strings.Split(split[0], ", ") {
		if strings.HasPrefix(design, s) {
			n += d.permutations(cache, split, design[len(s):])
		}
	}
	return n
}

func (d Day19) Part1(input string) (string, error) {
	split := strings.Split(strings.TrimSpace(input), "\n\n")

	cache := map[string]int{}

	part1 := 0
	for _, s := range strings.Fields(split[1]) {
		if w := d.permutations(cache, split, s); w >= 1 {
			part1++
		}
	}

	return strconv.Itoa(part1), nil
}

func (d Day19) Part2(input string) (string, error) {
	split := strings.Split(strings.TrimSpace(input), "\n\n")

	cache := map[string]int{}

	part2 := 0
	for _, s := range strings.Fields(split[1]) {
		if w := d.permutations(cache, split, s); w >= 1 {
			part2 += w
		}
	}

	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day19{})
}
