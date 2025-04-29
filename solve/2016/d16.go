package solve2016

import (
	"strings"

	"aoc/solve"
)

type Day16 struct {
}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 16}
}

// Expand the seed until it has length N
func (d Day16) expand(seed string, n int) string {
	for len(seed) < n {
		seed = seed + "0" + d.flip(d.reverse(seed))
	}
	return seed[:n]
}

// Flip the bits in a string (0 -> 1, 1 -> 0)
func (d Day16) flip(s string) string {
	var flipped strings.Builder
	for _, c := range s {
		if c == '0' {
			flipped.WriteByte('1')
		} else {
			flipped.WriteByte('0')
		}
	}
	return flipped.String()
}

// Reverse a string
func (d Day16) reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Compute the checksum of a string
func (d Day16) checksum(s string) string {
	for len(s)%2 == 0 {
		var sb strings.Builder
		for i := 0; i < len(s); i += 2 {
			if s[i] == s[i+1] {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		s = sb.String()
	}
	return s
}

func (d Day16) Part1(input string) (string, error) {
	seed := strings.TrimSpace(input)
	result := d.checksum(d.expand(seed, 272))
	return result, nil
}

func (d Day16) Part2(input string) (string, error) {
	seed := strings.TrimSpace(input)
	result := d.checksum(d.expand(seed, 35651584))
	return result, nil
}

func init() {
	solve.Register(Day16{})
}
