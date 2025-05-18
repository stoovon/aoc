package solve2023

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day13 struct {
}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 13}
}

func (d Day13) parsePatterns(input string) [][2][]uint64 {
	blocks := strings.Split(strings.TrimSpace(input), "\n\n")
	patterns := make([][2][]uint64, len(blocks))
	for i, block := range blocks {
		lines := strings.Split(block, "\n")
		rows := make([]uint64, len(lines))
		cols := make([]uint64, len(lines[0]))
		for r, line := range lines {
			var row uint64
			for _, c := range line {
				row <<= 1
				if c == '#' {
					row |= 1
				}
			}
			rows[r] = row
		}
		for c := range cols {
			var col uint64
			for r := range lines {
				col <<= 1
				if lines[r][c] == '#' {
					col |= 1
				}
			}
			cols[c] = col
		}
		patterns[i] = [2][]uint64{rows, cols}
	}
	return patterns
}

func findReflection(values []uint64) (int, bool) {
	for i := 0; i < len(values)-1; i++ {
		if values[i] != values[i+1] {
			continue
		}
		reflection := true
		for d := 1; i-d >= 0 && i+d+1 < len(values); d++ {
			if values[i-d] != values[i+d+1] {
				reflection = false
				break
			}
		}
		if reflection {
			return i, true
		}
	}
	return 0, false
}

func hamming1(a, b uint64) bool {
	x := a ^ b
	return x != 0 && (x&(x-1)) == 0
}

func findSmudgedReflection(values []uint64) (int, bool) {
	for i := 0; i < len(values)-1; i++ {
		smudgeSeen := false
		if values[i] != values[i+1] {
			if !hamming1(values[i], values[i+1]) {
				continue
			}
			smudgeSeen = true
		}
		reflection := true
		for d := 1; i-d >= 0 && i+d+1 < len(values); d++ {
			a, b := values[i-d], values[i+d+1]
			if a != b {
				if !hamming1(a, b) || smudgeSeen {
					reflection = false
					break
				}
				smudgeSeen = true
			}
		}
		if reflection && smudgeSeen {
			return i, true
		}
	}
	return 0, false
}

func (d Day13) Part1(input string) (string, error) {
	patterns := d.parsePatterns(input)
	var sum int
	for _, pat := range patterns {
		if i, ok := findReflection(pat[0]); ok {
			sum += 100 * (i + 1)
		} else if i, ok := findReflection(pat[1]); ok {
			sum += i + 1
		} else {
			return "", errors.New("no reflection found")
		}
	}
	return strconv.Itoa(sum), nil
}

func (d Day13) Part2(input string) (string, error) {
	patterns := d.parsePatterns(input)
	var sum int
	for _, pat := range patterns {
		if i, ok := findSmudgedReflection(pat[0]); ok {
			sum += 100 * (i + 1)
		} else if i, ok := findSmudgedReflection(pat[1]); ok {
			sum += i + 1
		} else {
			return "", errors.New("no smudged reflection found")
		}
	}
	return strconv.Itoa(sum), nil
}

func init() {
	solve.Register(Day13{})
}
