package solve2024

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day19 struct{}

func (d Day19) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2024, Day: 19} }

type towelPatterns [256][]string

type towelInput struct {
	patterns towelPatterns
	designs  []string
}

func validTowelText(text string) bool {
	if text == "" {
		return false
	}
	for i := 0; i < len(text); i++ {
		if text[i] < 'a' || text[i] > 'z' {
			return false
		}
	}
	return true
}

func parseTowels(input string) (towelInput, error) {
	var result towelInput
	text := strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n"))
	stock, designs, ok := strings.Cut(text, "\n\n")
	if !ok || strings.Contains(stock, "\n") {
		return result, fmt.Errorf("expected pattern list and design section")
	}
	seen := make(map[string]bool)
	for _, item := range strings.Split(stock, ",") {
		pattern := strings.TrimSpace(item)
		if !validTowelText(pattern) {
			return result, fmt.Errorf("expected nonempty lowercase ASCII pattern")
		}
		if seen[pattern] {
			continue
		}
		seen[pattern] = true
		result.patterns[pattern[0]] = append(result.patterns[pattern[0]], pattern)
	}
	for line, item := range strings.Split(designs, "\n") {
		design := strings.TrimSpace(item)
		if !validTowelText(design) {
			return result, fmt.Errorf("design %d: expected nonempty lowercase ASCII text", line+1)
		}
		result.designs = append(result.designs, design)
	}
	return result, nil
}

func (patterns *towelPatterns) canMake(design string, cache map[string]bool) bool {
	if design == "" {
		return true
	}
	if possible, known := cache[design]; known {
		return possible
	}
	for _, pattern := range patterns[design[0]] {
		if strings.HasPrefix(design, pattern) && patterns.canMake(design[len(pattern):], cache) {
			cache[design] = true
			return true
		}
	}
	cache[design] = false
	return false
}

// Cached counts are immutable; add them into a distinct result value.
func (patterns *towelPatterns) countWays(design string, cache map[string]*big.Int) *big.Int {
	if count, known := cache[design]; known {
		return count
	}
	count := new(big.Int)
	if design == "" {
		count.SetInt64(1)
	} else {
		for _, pattern := range patterns[design[0]] {
			if strings.HasPrefix(design, pattern) {
				count.Add(count, patterns.countWays(design[len(pattern):], cache))
			}
		}
	}
	cache[design] = count
	return count
}

func (d Day19) Part1(input string) (string, error) {
	towels, err := parseTowels(input)
	if err != nil {
		return "", err
	}
	cache := make(map[string]bool)
	possible := 0
	for _, design := range towels.designs {
		if towels.patterns.canMake(design, cache) {
			possible++
		}
	}
	return strconv.Itoa(possible), nil
}
func (d Day19) Part2(input string) (string, error) {
	towels, err := parseTowels(input)
	if err != nil {
		return "", err
	}
	cache := make(map[string]*big.Int)
	var total big.Int
	for _, design := range towels.designs {
		total.Add(&total, towels.patterns.countWays(design, cache))
	}
	return total.String(), nil
}
func init() { solve.Register(Day19{}) }
