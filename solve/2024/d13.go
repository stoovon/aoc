package solve2024

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day13 struct {
}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 13}
}

func calc(a, b, c image.Point) int {
	ap := (b.Y*c.X - b.X*c.Y) / (a.X*b.Y - a.Y*b.X)
	bp := (a.Y*c.X - a.X*c.Y) / (a.Y*b.X - a.X*b.Y)
	if a.Mul(ap).Add(b.Mul(bp)) == c {
		return ap*3 + bp
	}
	return 0
}

func (d Day13) Part1(input string) (string, error) {
	part1, part2 := 0, 0
	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n\n") {
		var a, b, c image.Point
		_, err := fmt.Sscanf(s, "Button A: X+%d, Y+%d\nButton B: X+%d, Y+%d\nPrize: X=%d, Y=%d",
			&a.X, &a.Y, &b.X, &b.Y, &c.X, &c.Y)
		if err != nil {
			return "", err
		}
		part1 += calc(a, b, c)
		part2 += calc(a, b, c.Add(image.Point{10000000000000, 10000000000000}))
	}

	return strconv.Itoa(part1), nil
}

func (d Day13) Part2(input string) (string, error) {
	part1, part2 := 0, 0
	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n\n") {
		var a, b, c image.Point
		_, err := fmt.Sscanf(s, "Button A: X+%d, Y+%d\nButton B: X+%d, Y+%d\nPrize: X=%d, Y=%d",
			&a.X, &a.Y, &b.X, &b.Y, &c.X, &c.Y)
		if err != nil {
			return "", err
		}
		part1 += calc(a, b, c)
		part2 += calc(a, b, c.Add(image.Point{10000000000000, 10000000000000}))
	}

	return strconv.Itoa(part2), nil
}

func init() {
	solve.Register(Day13{})
}
