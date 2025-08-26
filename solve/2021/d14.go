package solve2021

import (
	"aoc/solve"
	"fmt"
	"strings"
)

type Day14 struct{}

// parseInput parses the template and rules from the input string
func (d Day14) parseInput(input string) (string, map[string]string) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	template := lines[0]
	rules := make(map[string]string)
	for _, line := range lines[2:] {
		if len(line) < 6 {
			continue
		}
		parts := strings.Split(line, " -> ")
		if len(parts) == 2 {
			rules[parts[0]] = parts[1]
		}
	}
	return template, rules
}

// simulatePolymerNaive runs the naive simulation for small step counts
func simulatePolymerNaive(template string, rules map[string]string, steps int) string {
	polymer := template
	for step := 0; step < steps; step++ {
		var sb strings.Builder
		for i := 0; i < len(polymer)-1; i++ {
			pair := polymer[i : i+2]
			sb.WriteByte(polymer[i])
			if insert, ok := rules[pair]; ok {
				sb.WriteString(insert)
			}
		}
		sb.WriteByte(polymer[len(polymer)-1])
		polymer = sb.String()
	}
	return polymer
}

// simulatePolymerPairs runs the efficient simulation for large step counts
func simulatePolymerPairs(template string, rules map[string]string, steps int) map[byte]int {
	pairCounts := make(map[string]int)
	elemCounts := make(map[byte]int)
	// Initialize pair counts
	for i := 0; i < len(template)-1; i++ {
		pair := template[i : i+2]
		pairCounts[pair]++
	}
	// Initialize element counts
	for i := 0; i < len(template); i++ {
		elemCounts[template[i]]++
	}
	for step := 0; step < steps; step++ {
		newPairCounts := make(map[string]int)
		for pair, count := range pairCounts {
			if insert, ok := rules[pair]; ok {
				// The pair AB -> C produces AC and CB
				newPairCounts[string(pair[0])+insert] += count
				newPairCounts[insert+string(pair[1])] += count
				elemCounts[insert[0]] += count
			} else {
				newPairCounts[pair] += count
			}
		}
		pairCounts = newPairCounts
	}
	return elemCounts
}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 14}
}

func (d Day14) Part1(input string) (string, error) {
	template, rules := d.parseInput(input)
	polymer := simulatePolymerNaive(template, rules, 10)
	freq := make(map[byte]int)
	for i := 0; i < len(polymer); i++ {
		freq[polymer[i]]++
	}
	min, max := -1, -1
	for _, v := range freq {
		if min == -1 || v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return fmt.Sprintf("%d", max-min), nil
}

func (d Day14) Part2(input string) (string, error) {
	template, rules := d.parseInput(input)
	freq := simulatePolymerPairs(template, rules, 40)
	min, max := -1, -1
	for _, v := range freq {
		if min == -1 || v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return fmt.Sprintf("%d", max-min), nil
}

func init() {
	solve.Register(Day14{})
}
