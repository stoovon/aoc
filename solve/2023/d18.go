package solve2023

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day18 struct {
}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 18}
}

type dirStep struct {
	dir  byte
	dist int64
}

var (
	directionInstructionRE = regexp.MustCompile(`(?m)^([RLDU]) (\d+)`)
	hexInstructionRE       = regexp.MustCompile(`\(#([0-9a-fA-F]{5})([0-3])\)$`)
)

func parseSmallDirections(input string) ([]dirStep, error) {
	var steps []dirStep
	for _, m := range directionInstructionRE.FindAllStringSubmatch(input, -1) {
		dist, err := strconv.ParseInt(m[2], 10, 64)
		if err != nil {
			return nil, err
		}
		steps = append(steps, dirStep{dir: m[1][0], dist: dist})
	}
	return steps, nil
}

func parseLargeDirections(input string) ([]dirStep, error) {
	dirs := []byte{'R', 'D', 'L', 'U'}
	var steps []dirStep
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		m := hexInstructionRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		dist, err := strconv.ParseInt(m[1], 16, 64)
		if err != nil {
			return nil, err
		}
		didx, _ := strconv.Atoi(m[2])
		steps = append(steps, dirStep{dir: dirs[didx], dist: dist})
	}
	return steps, nil
}

func getArea(steps []dirStep) int64 {
	var perimeter, area int64
	y, x := int64(0), int64(0)
	for _, s := range steps {
		switch s.dir {
		case 'U':
			perimeter += s.dist
			area -= x * s.dist
			y -= s.dist
		case 'R':
			perimeter += s.dist
			x += s.dist
		case 'D':
			perimeter += s.dist
			area += x * s.dist
			y += s.dist
		case 'L':
			perimeter += s.dist
			x -= s.dist
		default:
			panic(fmt.Sprintf("Unknown dir: %c", s.dir))
		}
	}
	return area + perimeter/2 + 1
}

func (d Day18) Part1(input string) (string, error) {
	steps, err := parseSmallDirections(input)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(getArea(steps), 10), nil
}

func (d Day18) Part2(input string) (string, error) {
	steps, err := parseLargeDirections(input)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(getArea(steps), 10), nil
}

func init() {
	solve.Register(Day18{})
}
