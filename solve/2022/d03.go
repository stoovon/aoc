package solve2022

import (
	"errors"
	"fmt"
	"strings"

	"aoc/solve"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 3}
}

func (d Day3) Part1(input string) (string, error) {
	lines := strings.Fields(strings.TrimSpace(input))
	sum := 0
	for _, line := range lines {
		n := len(line) / 2
		left := make(map[rune]bool)
		for _, c := range line[:n] {
			left[c] = true
		}
		var dup rune
		for _, c := range line[n:] {
			if left[c] {
				dup = c
				break
			}
		}
		var priority int
		if dup >= 'a' && dup <= 'z' {
			priority = int(dup-'a') + 1
		} else if dup >= 'A' && dup <= 'Z' {
			priority = int(dup-'A') + 27
		}
		sum += priority
	}
	return fmt.Sprintf("%d", sum), nil
}

func (d Day3) Part2(input string) (string, error) {
	lines := strings.Fields(strings.TrimSpace(input))
	if len(lines)%3 != 0 {
		return "", errors.New("input does not contain a multiple of 3 lines")
	}
	sum := 0
	for i := 0; i < len(lines); i += 3 {
		// Build sets for each rucksack
		set1 := make(map[rune]bool)
		set2 := make(map[rune]bool)
		for _, c := range lines[i] {
			set1[c] = true
		}
		for _, c := range lines[i+1] {
			set2[c] = true
		}
		// Find the common item in all three
		var badge rune
		for _, c := range lines[i+2] {
			if set1[c] && set2[c] {
				badge = c
				break
			}
		}
		var priority int
		if badge >= 'a' && badge <= 'z' {
			priority = int(badge-'a') + 1
		} else if badge >= 'A' && badge <= 'Z' {
			priority = int(badge-'A') + 27
		}
		sum += priority
	}
	return fmt.Sprintf("%d", sum), nil
}

func init() {
	solve.Register(Day3{})
}
