package solve2025

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day7 struct {
}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 7}
}

// simulate runs the beam propagation and returns (splitCount, finalBeamSum)
func simulate(input string) (int64, int64) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	var startCol int
	for c, ch := range lines[0] {
		if ch == 'S' {
			startCol = c
			break
		}
	}

	splitCount := int64(0)

	curr := make(map[int]int)
	curr[startCol] = 1

	for row := 1; row < len(lines); row++ {
		nxt := make(map[int]int)

		for col, count := range curr {
			cell := rune(lines[row][col])

			if cell == '^' {
				splitCount++
				nxt[col-1] += count
				nxt[col+1] += count
			} else {
				nxt[col] += count
			}
		}

		curr = nxt
	}

	sum := int64(0)
	for _, count := range curr {
		sum += int64(count)
	}

	return splitCount, sum
}

func (d Day7) Part1(input string) (string, error) {
	splitCount, _ := simulate(input)
	return strconv.FormatInt(splitCount, 10), nil
}

func (d Day7) Part2(input string) (string, error) {
	_, beamSum := simulate(input)
	return strconv.FormatInt(beamSum, 10), nil
}

func init() {
	solve.Register(Day7{})
}
