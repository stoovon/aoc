package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day8 struct {
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 8}
}

func (d Day8) Part1(input string) (string, error) {
	charCount := 0
	symbolCount := 0

	for line := range strings.Lines(input) {
		charCount += len(line)
		subline := line[1 : len(line)-1] // Remove the surrounding quotes
		result := strings.Builder{}

		for i := 0; i < len(subline); i++ {
			ch := subline[i]
			if ch == '\\' {
				i++
				if subline[i] == 'x' {
					// Skip the next two characters (hex code)
					i += 2
					result.WriteByte('X') // Add a placeholder for the decoded character
				} else {
					result.WriteByte(subline[i]) // Add the escaped character
				}
			} else {
				result.WriteByte(ch) // Add the regular character
			}
		}

		symbolCount += result.Len()
	}

	return strconv.Itoa(charCount - symbolCount), nil
}

func (d Day8) Part2(input string) (string, error) {
	charCount := 0
	escapedCount := 0

	for line := range strings.Lines(input) {
		charCount += len(line)
		for _, ch := range line {
			if ch == '\\' || ch == '"' {
				escapedCount += 2
			} else {
				escapedCount++
			}
		}
		escapedCount += 2 // Add 2 for the surrounding quotes
	}

	return strconv.Itoa(escapedCount - charCount), nil
}

func init() {
	solve.Register(Day8{})
}
