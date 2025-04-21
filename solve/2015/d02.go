package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/acstrings"
	"aoc/utils/maths"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 2}
}

func boxArea(l, w, h int) int {
	return 2*l*w + 2*w*h + 2*h*l
}

func smallestSide(l, w, h int) int {
	return maths.Min(l*w, w*h, h*l)
}

func shortestPerimeter(l, w, h int) int {
	return 2 * maths.Min(l+w, w+h, h+l)
}

func boxVolume(l, w, h int) int {
	return l * w * h
}

func (d Day2) solve(input string) (area, length int) {
	area = 0
	length = 0
	for _, box := range strings.Fields(input) {
		sides := strings.Split(box, "x")
		l, w, h := acstrings.MustInt(sides[0]), acstrings.MustInt(sides[1]), acstrings.MustInt(sides[2])
		area += boxArea(l, w, h) + smallestSide(l, w, h)
		length += shortestPerimeter(l, w, h) + boxVolume(l, w, h)
	}
	
	return area, length
}

func (d Day2) Part1(input string) (string, error) {
	area, _ := d.solve(input)
	return strconv.Itoa(area), nil
}

func (d Day2) Part2(input string) (string, error) {
	_, length := d.solve(input)
	return strconv.Itoa(length), nil
}

func init() {
	solve.Register(Day2{})
}
