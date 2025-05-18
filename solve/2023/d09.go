package solve2023

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day9 struct {
}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 9}
}

func (d Day9) parseInput(input string) (int64, int64) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var partOne, partTwo int64

	for _, line := range lines {
		fields := strings.Fields(line)
		current := make([]int64, len(fields))
		for i, f := range fields {
			n, _ := strconv.ParseInt(f, 10, 64)
			current[i] = n
		}

		var starts, ends []int64
		for {
			starts = append(starts, current[0])
			ends = append(ends, current[len(current)-1])
			allZero := true
			next := make([]int64, len(current)-1)
			for i := 0; i < len(current)-1; i++ {
				next[i] = current[i+1] - current[i]
				if next[i] != 0 {
					allZero = false
				}
			}
			if allZero {
				break
			}
			current = next
		}

		// Part 1: sum of ends
		for _, v := range ends {
			partOne += v
		}
		// Part 2: fold starts in reverse
		acc := int64(0)
		for i := len(starts) - 1; i >= 0; i-- {
			acc = starts[i] - acc
		}
		partTwo += acc
	}
	return partOne, partTwo
}

func (d Day9) Part1(input string) (string, error) {
	partOne, _ := d.parseInput(input)
	return strconv.FormatInt(partOne, 10), nil
}

func (d Day9) Part2(input string) (string, error) {
	_, partTwo := d.parseInput(input)
	return strconv.FormatInt(partTwo, 10), nil
}

func init() {
	solve.Register(Day9{})
}
