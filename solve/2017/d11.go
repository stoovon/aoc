package solve2017

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day11 struct {
}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 11}
}

func hexDistance(moves []string) (finalDist, maxDist int) {
	x, y, z := 0, 0, 0
	maxDist = 0
	for _, move := range moves {
		switch move {
		case "n":
			y++
			z--
		case "ne":
			x++
			z--
		case "se":
			x++
			y--
		case "s":
			y--
			z++
		case "sw":
			x--
			z++
		case "nw":
			x--
			y++
		}
		dist := int(math.Max(math.Max(math.Abs(float64(x)), math.Abs(float64(y))), math.Abs(float64(z))))
		if dist > maxDist {
			maxDist = dist
		}
	}
	finalDist = int(math.Max(math.Max(math.Abs(float64(x)), math.Abs(float64(y))), math.Abs(float64(z))))
	return
}

func (d Day11) Part1(input string) (string, error) {
	moves := strings.Split(strings.TrimSpace(input), ",")
	if len(moves) == 0 {
		return "", errors.New("no moves")
	}
	finalDist, _ := hexDistance(moves)
	return strconv.Itoa(finalDist), nil
}

func (d Day11) Part2(input string) (string, error) {
	moves := strings.Split(strings.TrimSpace(input), ",")
	if len(moves) == 0 {
		return "", errors.New("no moves")
	}
	_, maxDist := hexDistance(moves)
	return strconv.Itoa(maxDist), nil
}

func init() {
	solve.Register(Day11{})
}
