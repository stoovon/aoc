package solve2025

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day3 struct{}

func (d Day3) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2025, Day: 3} }

// maxConcatK selects exactly k digits from `line` (preserving order) to form the
// largest possible k-digit number and returns it as an int64. Non-digits return 0.
func maxConcatK(line string, k int) int64 {
	n := len(line)
	if k <= 0 || n < k {
		return 0
	}

	for i := 0; i < n; i++ {
		if line[i] < '0' || line[i] > '9' {
			return 0
		}
	}

	var res int64
	start := 0
	remaining := k
	for remaining > 0 {
		end := n - remaining
		bestDigit := byte('0')
		bestIdx := start
		for i := start; i <= end; i++ {
			if line[i] > bestDigit {
				bestDigit = line[i]
				bestIdx = i
				if bestDigit == '9' {
					break
				}
			}
		}

		res = res*10 + int64(bestDigit-'0')
		start = bestIdx + 1
		remaining--
	}

	return res
}

func (d Day3) Part1(input string) (string, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return "0", nil
	}

	lines := strings.Split(s, "\n")
	var total int64
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		total += maxConcatK(line, 2)
	}

	return strconv.FormatInt(total, 10), nil
}

func (d Day3) Part2(input string) (string, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return "0", nil
	}

	lines := strings.Split(s, "\n")
	var total int64
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		total += maxConcatK(line, 12)
	}

	return strconv.FormatInt(total, 10), nil
}

func init() { solve.Register(Day3{}) }
