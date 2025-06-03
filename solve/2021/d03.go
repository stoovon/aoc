package solve2021

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 3}
}

func (d Day3) Part1(input string) (string, error) {
	lines := strings.Fields(strings.TrimSpace(input))
	if len(lines) == 0 {
		return "", errors.New("no input")
	}
	bitLen := len(lines[0])
	counts := make([]int, bitLen)
	for _, line := range lines {
		for i, c := range line {
			if c == '1' {
				counts[i]++
			}
		}
	}
	var gamma, epsilon int
	half := len(lines) / 2
	for _, c := range counts {
		gamma <<= 1
		epsilon <<= 1
		if c > half {
			gamma |= 1
		} else {
			epsilon |= 1
		}
	}
	return fmt.Sprintf("%d", gamma*epsilon), nil
}

func filterByBitCriteria(lines []string, mostCommon bool) int {
	remaining := make([]string, len(lines))
	copy(remaining, lines)
	bitLen := len(lines[0])
	for i := 0; i < bitLen && len(remaining) > 1; i++ {
		ones, zeros := 0, 0
		for _, line := range remaining {
			if line[i] == '1' {
				ones++
			} else {
				zeros++
			}
		}
		var keep byte
		if mostCommon {
			if ones >= zeros {
				keep = '1'
			} else {
				keep = '0'
			}
		} else {
			if zeros <= ones {
				keep = '0'
			} else {
				keep = '1'
			}
		}
		filtered := remaining[:0]
		for _, line := range remaining {
			if line[i] == keep {
				filtered = append(filtered, line)
			}
		}
		remaining = filtered
	}
	val, _ := strconv.ParseInt(remaining[0], 2, 64)
	return int(val)
}

func (d Day3) Part2(input string) (string, error) {
	lines := strings.Fields(strings.TrimSpace(input))
	if len(lines) == 0 {
		return "", errors.New("no input")
	}
	oxygen := filterByBitCriteria(lines, true)
	co2 := filterByBitCriteria(lines, false)
	return fmt.Sprintf("%d", oxygen*co2), nil
}

func init() {
	solve.Register(Day3{})
}
