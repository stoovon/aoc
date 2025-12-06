package solve2025

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 5}
}

func (d Day5) Part1(input string) (string, error) {
	parts := strings.Split(input, "\n\n")
	if len(parts) < 2 {
		return "0", nil
	}

	ranges := parseRanges(strings.TrimSpace(parts[0]))
	availableLines := strings.TrimSpace(parts[1])

	count := 0
	for _, line := range strings.Split(availableLines, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		id, _ := strconv.ParseInt(line, 10, 64)

		for _, r := range ranges {
			if id >= r[0] && id <= r[1] {
				count++
				break
			}
		}
	}

	return strconv.Itoa(count), nil
}

func (d Day5) Part2(input string) (string, error) {
	parts := strings.Split(input, "\n\n")
	if len(parts) < 1 {
		return "0", nil
	}

	ranges := parseRanges(strings.TrimSpace(parts[0]))

	merged := mergeRanges(ranges)

	var count int64
	for _, r := range merged {
		count += r[1] - r[0] + 1
	}

	return strconv.FormatInt(count, 10), nil
}

func parseRanges(rangeLines string) [][2]int64 {
	var ranges [][2]int64
	for _, line := range strings.Split(rangeLines, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		components := strings.Split(line, "-")
		if len(components) != 2 {
			continue
		}
		start, _ := strconv.ParseInt(strings.TrimSpace(components[0]), 10, 64)
		end, _ := strconv.ParseInt(strings.TrimSpace(components[1]), 10, 64)

		if start > end {
			start, end = end, start
		}
		ranges = append(ranges, [2]int64{start, end})
	}
	return ranges
}

func mergeRanges(ranges [][2]int64) [][2]int64 {
	if len(ranges) == 0 {
		return ranges
	}

	for i := range ranges {
		for j := i + 1; j < len(ranges); j++ {
			if ranges[j][0] < ranges[i][0] {
				ranges[i], ranges[j] = ranges[j], ranges[i]
			}
		}
	}

	var merged [][2]int64
	current := ranges[0]

	for i := 1; i < len(ranges); i++ {
		next := ranges[i]
		if next[0] <= current[1]+1 {
			if next[1] > current[1] {
				current[1] = next[1]
			}
		} else {
			merged = append(merged, current)
			current = next
		}
	}
	merged = append(merged, current)

	return merged
}

func init() {
	solve.Register(Day5{})
}
