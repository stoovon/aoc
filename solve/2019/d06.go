package solve2019

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day6 struct {
}

func (d Day6) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 6}
}

func (d Day6) Part1(input string) (string, error) {
	orbitMap := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Split(line, ")")
		if len(parts) != 2 {
			continue
		}
		orbitMap[parts[1]] = parts[0]
	}

	countOrbits := func(obj string) int {
		count := 0
		for {
			parent, ok := orbitMap[obj]
			if !ok {
				break
			}
			count++
			obj = parent
		}
		return count
	}

	total := 0
	for obj := range orbitMap {
		total += countOrbits(obj)
	}
	return strconv.Itoa(total), nil
}

func (d Day6) Part2(input string) (string, error) {
	orbitMap := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Split(line, ")")
		if len(parts) != 2 {
			continue
		}
		orbitMap[parts[1]] = parts[0]
	}

	// Build path from YOU to COM
	pathToCOM := func(start string) map[string]int {
		path := make(map[string]int)
		steps := 0
		for obj := orbitMap[start]; obj != ""; obj = orbitMap[obj] {
			path[obj] = steps
			steps++
		}
		return path
	}

	youPath := pathToCOM("YOU")
	sanPath := pathToCOM("SAN")

	// Find the first common ancestor
	minTransfers := -1
	for obj, youSteps := range youPath {
		if sanSteps, ok := sanPath[obj]; ok {
			transfers := youSteps + sanSteps
			if minTransfers == -1 || transfers < minTransfers {
				minTransfers = transfers
			}
		}
	}

	if minTransfers == -1 {
		return "", errors.New("no path found for YOU -> SAN")
	}
	return strconv.Itoa(minTransfers), nil
}

func init() {
	solve.Register(Day6{})
}
