package solve2020

import (
	"aoc/solve"
	"errors"
	"strconv"
	"strings"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 5}
}

func (d Day5) Part1(input string) (string, error) {
	maxID := 0
	for _, line := range strings.Fields(input) {
		row := 0
		for i := 0; i < 7; i++ {
			row <<= 1
			if line[i] == 'B' {
				row |= 1
			}
		}
		col := 0
		for i := 7; i < 10; i++ {
			col <<= 1
			if line[i] == 'R' {
				col |= 1
			}
		}
		seatID := row*8 + col
		if seatID > maxID {
			maxID = seatID
		}
	}
	return strconv.Itoa(maxID), nil
}

func (d Day5) Part2(input string) (string, error) {
	seats := make(map[int]bool)
	minID, maxID := 1<<31-1, 0
	for _, line := range strings.Fields(input) {
		row := 0
		for i := 0; i < 7; i++ {
			row <<= 1
			if line[i] == 'B' {
				row |= 1
			}
		}
		col := 0
		for i := 7; i < 10; i++ {
			col <<= 1
			if line[i] == 'R' {
				col |= 1
			}
		}
		seatID := row*8 + col
		seats[seatID] = true
		if seatID < minID {
			minID = seatID
		}
		if seatID > maxID {
			maxID = seatID
		}
	}
	for id := minID + 1; id < maxID; id++ {
		if !seats[id] && seats[id-1] && seats[id+1] {
			return strconv.Itoa(id), nil
		}
	}
	return "", errors.New("seat not found")
}

func init() {
	solve.Register(Day5{})
}
