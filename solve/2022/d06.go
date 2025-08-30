package solve2022

import (
	"aoc/solve"
	"errors"
	"fmt"
)

type Day6 struct{}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 6}
}

func (d Day6) findMarker(input string, markerLength int) (string, error) {
	// Iterate through the input string
	for i := markerLength - 1; i < len(input); i++ {
		charSet := make(map[byte]bool)
		for j := 0; j < markerLength; j++ {
			charSet[input[i-j]] = true
		}
		if len(charSet) == markerLength {
			// Return the position (1-based index)
			return fmt.Sprintf("%d", i+1), nil
		}
	}
	return "", errors.New("no marker found")
}

func (d Day6) Part1(input string) (string, error) {
	return d.findMarker(input, 4)
}

func (d Day6) Part2(input string) (string, error) {
	return d.findMarker(input, 14)
}

func init() {
	solve.Register(Day6{})
}
