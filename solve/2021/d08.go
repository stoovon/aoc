package solve2021

import (
	"aoc/solve"
	"errors"
	"sort"
	"strconv"
	"strings"
)

type Day8 struct{}

func sortString(s string) string {
	chars := strings.Split(s, "")
	sort.Strings(chars)
	return strings.Join(chars, "")
}

func deduceDigitMap(patterns []string) map[string]int {
	byLen := make(map[int][]string)
	for _, p := range patterns {
		byLen[len(p)] = append(byLen[len(p)], p)
	}
	digit := make(map[int]string)

	// Unique lengths
	digit[1] = byLen[2][0]
	digit[4] = byLen[4][0]
	digit[7] = byLen[3][0]
	digit[8] = byLen[7][0]

	// Helper: contains all chars of b in a
	contains := func(a, b string) bool {
		for _, c := range b {
			if !strings.ContainsRune(a, c) {
				return false
			}
		}
		return true
	}

	// Deduce 0, 6, 9 (len 6)
	for _, p := range byLen[6] {
		if contains(p, digit[4]) {
			digit[9] = p
		} else if contains(p, digit[1]) {
			digit[0] = p
		} else {
			digit[6] = p
		}
	}

	// Deduce 2, 3, 5 (len 5)
	for _, p := range byLen[5] {
		if contains(p, digit[1]) {
			digit[3] = p
		} else if contains(digit[6], p) {
			digit[5] = p
		} else {
			digit[2] = p
		}
	}

	// Build map from sorted pattern to digit
	result := make(map[string]int)
	for d, pat := range digit {
		result[sortString(pat)] = d
	}

	return result
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 8}
}

func (d Day8) parse(input string) ([][2][]string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, errors.New("empty input")
	}
	lines := strings.Split(input, "\n")
	entries := make([][2][]string, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}
		patterns := strings.Fields(parts[0])
		outputs := strings.Fields(parts[1])
		entries = append(entries, [2][]string{patterns, outputs})
	}
	return entries, nil
}

func (d Day8) Part1(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("empty input")
	}
	lines := strings.Split(input, "\n")
	uniqueLens := map[int]bool{2: true, 3: true, 4: true, 7: true}
	count := 0
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}
		outputs := strings.Fields(parts[1])
		for _, out := range outputs {
			if uniqueLens[len(out)] {
				count++
			}
		}
	}
	return strconv.Itoa(count), nil
}

func (d Day8) Part2(input string) (string, error) {
	entries, err := d.parse(input)
	if err != nil {
		return "", err
	}
	sum := 0
	for _, entry := range entries {
		patterns, outputs := entry[0], entry[1]
		digitMap := deduceDigitMap(patterns)
		value := 0
		for _, out := range outputs {
			sorted := sortString(out)
			digit := digitMap[sorted]
			value = value*10 + digit
		}
		sum += value
	}
	return strconv.Itoa(sum), nil
}

func init() {
	solve.Register(Day8{})
}
