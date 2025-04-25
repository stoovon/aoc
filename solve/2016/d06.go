package solve2016

import (
	"sort"
	"strings"

	"aoc/solve"
)

type Day6 struct {
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 6}
}

// Transpose converts rows into columns
func (d Day6) transpose(lines []string) [][]rune {
	if len(lines) == 0 {
		return nil
	}
	columns := make([][]rune, len(lines[0]))
	for _, line := range lines {
		for i, char := range line {
			columns[i] = append(columns[i], char)
		}
	}
	return columns
}

// Count characters and return sorted by frequency and alphabetically
func (d Day6) countCharacters(column []rune) []struct {
	char  rune
	count int
} {
	counts := make(map[rune]int)
	for _, char := range column {
		counts[char]++
	}

	var result []struct {
		char  rune
		count int
	}
	for char, count := range counts {
		result = append(result, struct {
			char  rune
			count int
		}{char, count})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].count == result[j].count {
			return result[i].char < result[j].char
		}
		return result[i].count > result[j].count
	})

	return result
}

func (d Day6) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	columns := d.transpose(lines)

	var result strings.Builder
	for _, column := range columns {
		counts := d.countCharacters(column)
		if len(counts) > 0 {
			result.WriteRune(counts[0].char) // Most common character
		}
	}

	return result.String(), nil
}

func (d Day6) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	columns := d.transpose(lines)

	var result strings.Builder
	for _, column := range columns {
		counts := d.countCharacters(column)
		if len(counts) > 0 {
			result.WriteRune(counts[len(counts)-1].char) // Least common character
		}
	}

	return result.String(), nil
}

func init() {
	solve.Register(Day6{})
}
