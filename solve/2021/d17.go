package solve2021

import (
	"aoc/solve"
	"errors"
	"fmt"
)

type Day17 struct{}

func parseTarget(input string) (xMin, xMax, yMin, yMax int, err error) {
	n, err := fmt.Sscanf(input, "target area: x=%d..%d, y=%d..%d", &xMin, &xMax, &yMin, &yMax)
	if n != 4 || err != nil {
		return 0, 0, 0, 0, errors.New("invalid input")
	}
	return
}

func hitsTarget(vx0, vy0, xMin, xMax, yMin, yMax int) (hit bool, maxYSeen int) {
	x, y := 0, 0
	vx, vy := vx0, vy0
	maxY := 0
	for x <= xMax && y >= yMin {
		x += vx
		y += vy
		if y > maxY {
			maxY = y
		}
		if vx > 0 {
			vx--
		} else if vx < 0 {
			vx++
		}
		vy--
		if x >= xMin && x <= xMax && y >= yMin && y <= yMax {
			return true, maxY
		}
	}
	return false, maxY
}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 17}
}

func (d Day17) Part1(input string) (string, error) {
	xMin, xMax, yMin, yMax, err := parseTarget(input)
	if err != nil {
		return "", err
	}
	highestY := 0
	for vx0 := 1; vx0 <= xMax; vx0++ {
		for vy0 := yMin; vy0 <= 1000; vy0++ {
			hit, maxY := hitsTarget(vx0, vy0, xMin, xMax, yMin, yMax)
			if hit && maxY > highestY {
				highestY = maxY
			}
		}
	}
	return fmt.Sprintf("%d", highestY), nil
}

func (d Day17) Part2(input string) (string, error) {
	xMin, xMax, yMin, yMax, err := parseTarget(input)
	if err != nil {
		return "", err
	}
	count := 0
	for vx0 := 1; vx0 <= xMax; vx0++ {
		for vy0 := yMin; vy0 <= 1000; vy0++ {
			hit, _ := hitsTarget(vx0, vy0, xMin, xMax, yMin, yMax)
			if hit {
				count++
			}
		}
	}
	return fmt.Sprintf("%d", count), nil
}

func init() {
	solve.Register(Day17{})
}
