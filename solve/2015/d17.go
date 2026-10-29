package solve2015

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day17 struct{}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 17}
}

type containerCounts struct {
	totalWays, minContainers, minWays int
}

func countContainerArrangements(sizes []int, target int) (containerCounts, error) {
	if target < 0 || target >= 1<<20 {
		return containerCounts{}, fmt.Errorf("target exceeds the volume table budget")
	}
	table := make([]containerCounts, target+1)
	table[0] = containerCounts{totalWays: 1, minWays: 1}
	maxCount := int(^uint(0) >> 1)
	for _, size := range sizes {
		if size < 0 {
			return containerCounts{}, fmt.Errorf("container size must be nonnegative")
		}
		// Descending volumes prevent reuse; equal-sized containers remain distinct.
		for volume := target; volume >= size; volume-- {
			source := table[volume-size]
			if source.totalWays == 0 {
				continue
			}
			entry := &table[volume]
			if source.totalWays > maxCount-entry.totalWays {
				return containerCounts{}, fmt.Errorf("container combination count overflows int")
			}
			entry.totalWays += source.totalWays
			used := source.minContainers + 1
			if entry.minWays == 0 || used < entry.minContainers {
				entry.minContainers, entry.minWays = used, source.minWays
			} else if used == entry.minContainers {
				entry.minWays += source.minWays
			}
		}
	}
	return table[target], nil
}

func (d Day17) arrangements(input string) (containerCounts, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return containerCounts{}, fmt.Errorf("container list is empty")
	}
	var sizes []int
	for index, line := range strings.Split(input, "\n") {
		size, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return containerCounts{}, fmt.Errorf("line %d: invalid container size: %w", index+1, err)
		}
		if size < 0 {
			return containerCounts{}, fmt.Errorf("line %d: container size must be nonnegative", index+1)
		}
		sizes = append(sizes, size)
	}
	return countContainerArrangements(sizes, 150)
}

func (d Day17) Part1(input string) (string, error) {
	counts, err := d.arrangements(input)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(counts.totalWays), nil
}

func (d Day17) Part2(input string) (string, error) {
	counts, err := d.arrangements(input)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(counts.minWays), nil
}

func init() {
	solve.Register(Day17{})
}
