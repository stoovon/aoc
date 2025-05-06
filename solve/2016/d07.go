package solve2016

import (
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day7 struct {
}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 7}
}

var segmentRe = regexp.MustCompile(`\[|\]`)

// Checks if a string contains an ABBA pattern
func (d Day7) abba(text string) bool {
	for i := 0; i < len(text)-3; i++ {
		a, b, c, d := text[i], text[i+1], text[i+2], text[i+3]
		if a == d && b == c && a != b {
			return true
		}
	}
	return false
}

// Splits a line into segments inside and outside square brackets
func (d Day7) segment(line string) []string {
	return segmentRe.Split(line, -1)
}

// Extracts segments outside square brackets
func (d Day7) outsides(segments []string) []string {
	var result []string
	for i := 0; i < len(segments); i += 2 {
		result = append(result, segments[i])
	}
	return result
}

// Extracts segments inside square brackets
func (d Day7) insides(segments []string) []string {
	var result []string
	for i := 1; i < len(segments); i += 2 {
		result = append(result, segments[i])
	}
	return result
}

// Checks if the segments support TLS
func (d Day7) tls(segments []string) bool {
	outside := strings.Join(d.outsides(segments), ",")
	inside := strings.Join(d.insides(segments), ",")
	return d.abba(outside) && !d.abba(inside)
}

// Checks if the segments support SSL
func (d Day7) ssl(segments []string) bool {
	outside := strings.Join(d.outsides(segments), ",")
	inside := strings.Join(d.insides(segments), ",")
	for _, a := range "abcdefghijklmnopqrstuvwxyz" {
		for _, b := range "abcdefghijklmnopqrstuvwxyz" {
			if a != b {
				aba := string([]rune{a, b, a})
				bab := string([]rune{b, a, b})
				if strings.Contains(outside, aba) && strings.Contains(inside, bab) {
					return true
				}
			}
		}
	}
	return false
}

func (d Day7) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	count := 0
	for _, line := range lines {
		if d.tls(d.segment(line)) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func (d Day7) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	count := 0
	for _, line := range lines {
		if d.ssl(d.segment(line)) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func init() {
	solve.Register(Day7{})
}
