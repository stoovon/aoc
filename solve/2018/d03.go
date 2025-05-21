package solve2018

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 3}
}

var (
	// Regex to parse: #123 @ 3,2: 5x4
	fabricSpecRE = regexp.MustCompile(`#\d+ @ (\d+),(\d+): (\d+)x(\d+)`)
)

func (d Day3) Solve(input string) (part1 string, part2 string, err error) {
	const size = 1000
	fabric := [size][size]int{}
	lines := strings.Split(strings.TrimSpace(input), "\n")

	type claim struct {
		id, left, top, width, height int
	}
	claims := make([]claim, 0, len(lines))
	overlap := make(map[int]bool)

	// Parse and fill fabric
	for _, line := range lines {
		m := fabricSpecRE.FindStringSubmatch(line)
		if m == nil {
			return "", "", fmt.Errorf("invalid claim: %s", line)
		}
		id, _ := strconv.Atoi(strings.Fields(line)[0][1:])
		left, _ := strconv.Atoi(m[1])
		top, _ := strconv.Atoi(m[2])
		width, _ := strconv.Atoi(m[3])
		height, _ := strconv.Atoi(m[4])
		claims = append(claims, claim{id, left, top, width, height})
		overlap[id] = false
		for i := left; i < left+width; i++ {
			for j := top; j < top+height; j++ {
				if fabric[i][j] == 0 {
					fabric[i][j] = id
				} else {
					overlap[fabric[i][j]] = true
					overlap[id] = true
					fabric[i][j] = -1 // Mark as overlapping
				}
			}
		}
	}

	// Part 1: count overlaps
	count := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if fabric[i][j] == -1 {
				count++
			}
		}
	}
	part1 = strconv.Itoa(count)

	// Part 2: find non-overlapping claim
	for _, c := range claims {
		if !overlap[c.id] {
			part2 = strconv.Itoa(c.id)
			return
		}
	}
	err = fmt.Errorf("no non-overlapping claim found")
	return
}

func (d Day3) Part1(input string) (string, error) {
	part1, _, err := d.Solve(input)
	if err != nil {
		return "", err
	}
	return part1, nil
}

func (d Day3) Part2(input string) (string, error) {
	_, part2, err := d.Solve(input)
	if err != nil {
		return "", err
	}
	return part2, nil
}

func init() {
	solve.Register(Day3{})
}
