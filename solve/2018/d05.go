package solve2018

import (
	"fmt"
	"strings"

	"aoc/solve"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 5}
}

func (d Day5) Part1(input string) (string, error) {
	polymer := []rune(strings.TrimSpace(input))
	stack := make([]rune, 0, len(polymer))
	for _, unit := range polymer {
		if len(stack) > 0 {
			top := stack[len(stack)-1]
			if top != unit && (top^32) == unit {
				stack = stack[:len(stack)-1]
				continue
			}
		}
		stack = append(stack, unit)
	}
	return fmt.Sprintf("%d", len(stack)), nil
}

func (d Day5) Part2(input string) (string, error) {
	original := []rune(strings.TrimSpace(input))
	minLen := len(original)

	for unit := 'A'; unit <= 'Z'; unit++ {
		filtered := make([]rune, 0, len(original))
		for _, r := range original {
			if r != unit && r != unit+32 {
				filtered = append(filtered, r)
			}
		}
		stack := make([]rune, 0, len(filtered))
		for _, r := range filtered {
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				if top != r && (top^32) == r {
					stack = stack[:len(stack)-1]
					continue
				}
			}
			stack = append(stack, r)
		}
		if len(stack) < minLen {
			minLen = len(stack)
		}
	}
	return fmt.Sprintf("%d", minLen), nil
}

func init() {
	solve.Register(Day5{})
}
