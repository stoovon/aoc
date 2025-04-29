package solve2016

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day15 struct {
}

type Disc struct {
	Index     int
	Positions int
	Current   int
}

var discRe = regexp.MustCompile(`#(\d+).* (\d+) positions.* (\d+)[.]`)

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 15}
}

// Parse the input into a slice of Disc structs
func (d Day15) parseInput(input string) ([]Disc, error) {
	input = strings.TrimSpace(input)
	matches := discRe.FindAllStringSubmatch(input, -1)

	if matches == nil {
		return nil, fmt.Errorf("invalid input format")
	}

	var discs []Disc
	for _, match := range matches {
		index, _ := strconv.Atoi(match[1])
		positions, _ := strconv.Atoi(match[2])
		current, _ := strconv.Atoi(match[3])
		discs = append(discs, Disc{Index: index, Positions: positions, Current: current})
	}

	return discs, nil
}

// Check if the capsule falls through all slots at time t
func (d Day15) falls(t int, discs []Disc) bool {
	for _, disc := range discs {
		if (disc.Current+t+disc.Index)%disc.Positions != 0 {
			return false
		}
	}
	return true
}

// Find the first time t that satisfies the condition
func (d Day15) findFirstValidTime(discs []Disc) int {
	t := 0
	for {
		if d.falls(t, discs) {
			return t
		}
		t++
	}
}

func (d Day15) Part1(input string) (string, error) {
	discs, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	result := d.findFirstValidTime(discs)
	return strconv.Itoa(result), nil
}

func (d Day15) Part2(input string) (string, error) {
	discs, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	// Add the extra disc for Part 2
	discs = append(discs, Disc{Index: len(discs) + 1, Positions: 11, Current: 0})

	result := d.findFirstValidTime(discs)
	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day15{})
}
