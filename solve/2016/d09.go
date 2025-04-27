package solve2016

import (
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day9 struct {
}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 9}
}

var markerRegex = regexp.MustCompile(`\((\d+)x(\d+)\)`)

// Decompresses the string in a single pass
func (d Day9) decompressSingleContent(s string) string {
	s = strings.ReplaceAll(s, " ", "") // Remove whitespace
	var result strings.Builder
	i := 0

	for i < len(s) {
		if loc := markerRegex.FindStringIndex(s[i:]); loc != nil && loc[0] == 0 {
			matches := markerRegex.FindStringSubmatch(s[i:])
			C, _ := strconv.Atoi(matches[1])
			R, _ := strconv.Atoi(matches[2])
			i += loc[1] // Move past the marker
			repeated := strings.Repeat(s[i:i+C], R)
			result.WriteString(repeated)
			i += C // Skip the repeated characters
		} else {
			result.WriteByte(s[i]) // Add regular character
			i++
		}
	}

	return result.String()
}

// Calculates the decompressed length recursively
func (d Day9) decompressRecursiveLengthOnly(s string) int {
	s = strings.ReplaceAll(s, " ", "") // Remove whitespace
	length := 0
	i := 0

	for i < len(s) {
		if loc := markerRegex.FindStringIndex(s[i:]); loc != nil && loc[0] == 0 {
			matches := markerRegex.FindStringSubmatch(s[i:])
			C, _ := strconv.Atoi(matches[1])
			R, _ := strconv.Atoi(matches[2])
			i += loc[1] // Move past the marker
			length += R * d.decompressRecursiveLengthOnly(s[i:i+C])
			i += C // Skip the repeated characters
		} else {
			length++ // Add regular character to length
			i++
		}
	}

	return length
}

func (d Day9) Part1(input string) (string, error) {
	decompressed := d.decompressSingleContent(strings.TrimSpace(input))
	return strconv.Itoa(len(decompressed)), nil
}

func (d Day9) Part2(input string) (string, error) {
	length := d.decompressRecursiveLengthOnly(strings.TrimSpace(input))
	return strconv.Itoa(length), nil
}

func init() {
	solve.Register(Day9{})
}
