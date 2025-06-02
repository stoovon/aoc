package solve2019

import (
	"errors"
	"fmt"
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/acstrings"
	"aoc/utils/maths"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 3}
}

func parseWire(line string) []string {
	return strings.Split(strings.TrimSpace(line), ",")
}

func traceWire(moves []string) (visitedPoints map[image.Point]struct{}, visitSteps map[image.Point]int) {
	x, y, steps := 0, 0, 0
	visited := make(map[image.Point]struct{})
	stepsMap := make(map[image.Point]int)
	for _, move := range moves {
		dir := move[0]
		dist, _ := strconv.Atoi(move[1:])
		for i := 0; i < dist; i++ {
			switch dir {
			case 'U':
				y++
			case 'D':
				y--
			case 'L':
				x--
			case 'R':
				x++
			}
			p := image.Point{X: x, Y: y}
			visited[p] = struct{}{}
			if _, exists := stepsMap[p]; !exists {
				steps++
				stepsMap[p] = steps
			} else {
				steps++
			}
		}
	}
	return visited, stepsMap
}

func (d Day3) Part1(input string) (string, error) {
	lines := acstrings.Lines(input)
	if len(lines) < 2 {
		return "", errors.New("expected two wire paths")
	}
	wire1 := parseWire(lines[0])
	wire2 := parseWire(lines[1])

	points1, _ := traceWire(wire1)
	points2, _ := traceWire(wire2)

	minDist := -1
	for p := range points1 {
		if p == (image.Point{}) {
			continue
		}
		if _, ok := points2[p]; ok {
			dist := maths.Abs(p.X) + maths.Abs(p.Y)
			if minDist == -1 || dist < minDist {
				minDist = dist
			}
		}
	}
	if minDist == -1 {
		return "", errors.New("no intersection found")
	}
	return fmt.Sprintf("%d", minDist), nil
}

func (d Day3) Part2(input string) (string, error) {
	lines := acstrings.Lines(input)
	if len(lines) < 2 {
		return "", errors.New("expected two wire paths")
	}
	wire1 := parseWire(lines[0])
	wire2 := parseWire(lines[1])

	_, steps1 := traceWire(wire1)
	_, steps2 := traceWire(wire2)

	minSteps := -1
	for p, s1 := range steps1 {
		if s2, ok := steps2[p]; ok {
			total := s1 + s2
			if minSteps == -1 || total < minSteps {
				minSteps = total
			}
		}
	}
	if minSteps == -1 {
		return "", errors.New("no intersection found")
	}
	return fmt.Sprintf("%d", minSteps), nil
}

func init() {
	solve.Register(Day3{})
}
