package solve2016

import (
	"errors"
	"image"
	"math"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day1 struct {
}

func (d Day1) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 1}
}

var parseRe = regexp.MustCompile(`([RL])(\d+)`)

func distance(p image.Point) int {
	return int(math.Abs(float64(p.X)) + math.Abs(float64(p.Y)))
}

func (d Day1) parse(input string) ([][2]interface{}, error) {
	matches := parseRe.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return nil, errors.New("invalid input format")
	}

	var moves [][2]interface{}
	for _, match := range matches {
		turn := match[1]
		dist, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}
		moves = append(moves, [2]interface{}{turn, dist})
	}
	return moves, nil
}

func (d Day1) howFarFromOrigin(input string) (int, error) {
	moves, err := d.parse(input)
	if err != nil {
		return 0, err
	}

	loc := image.Point{}
	heading := image.Point{Y: 1} // Start facing North

	for _, move := range moves {
		turn := move[0].(string)
		dist := move[1].(int)

		// Update heading based on turn
		if turn == "R" {
			heading = image.Point{X: heading.Y, Y: -heading.X}
		} else if turn == "L" {
			heading = image.Point{X: -heading.Y, Y: heading.X}
		}

		// Move in the current heading
		loc.X += heading.X * dist
		loc.Y += heading.Y * dist
	}

	return distance(loc), nil
}

func (d Day1) distanceToFirstLocationVisitedTwice(input string) (int, error) {
	moves, err := d.parse(input)
	if err != nil {
		return 0, err
	}

	loc := image.Point{}
	heading := image.Point{Y: 1} // Start facing North
	visited := map[image.Point]bool{loc: true}

	for _, move := range moves {
		turn := move[0].(string)
		dist := move[1].(int)

		// Update heading based on turn
		if turn == "R" {
			heading = image.Point{X: heading.Y, Y: -heading.X}
		} else if turn == "L" {
			heading = image.Point{X: -heading.Y, Y: heading.X}
		}

		// Move step by step
		for i := 0; i < dist; i++ {
			loc.X += heading.X
			loc.Y += heading.Y

			if visited[loc] {
				return distance(loc), nil
			}
			visited[loc] = true
		}
	}

	return 0, errors.New("no location visited twice")
}

func (d Day1) Part1(input string) (string, error) {

	result, err := d.howFarFromOrigin(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(result), nil
}

func (d Day1) Part2(input string) (string, error) {
	result, err := d.distanceToFirstLocationVisitedTwice(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day1{})
}
