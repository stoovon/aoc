package solve2021

import (
	"aoc/solve"
	"aoc/utils/maths"
	"fmt"
	"image"
	"strings"
)

type Day5 struct{}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 5}
}

func countOverlaps(input string, includeDiagonals bool) (int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	counts := make(map[image.Point]int)
	for _, l := range lines {
		var x1, y1, x2, y2 int
		_, err := fmt.Sscanf(l, "%d,%d -> %d,%d", &x1, &y1, &x2, &y2)
		if err != nil {
			return 0, err
		}
		switch {
		case x1 == x2:
			ystart, yend := y1, y2
			if ystart > yend {
				ystart, yend = yend, ystart
			}
			for y := ystart; y <= yend; y++ {
				counts[image.Point{x1, y}]++
			}
		case y1 == y2:
			xstart, xend := x1, x2
			if xstart > xend {
				xstart, xend = xend, xstart
			}
			for x := xstart; x <= xend; x++ {
				counts[image.Point{x, y1}]++
			}
		case includeDiagonals && maths.Abs(x2-x1) == maths.Abs(y2-y1):
			dx := 1
			if x2 < x1 {
				dx = -1
			}
			dy := 1
			if y2 < y1 {
				dy = -1
			}
			length := maths.Abs(x2 - x1)
			for i := 0; i <= length; i++ {
				counts[image.Point{x1 + i*dx, y1 + i*dy}]++
			}
		}
	}
	overlap := 0
	for _, c := range counts {
		if c >= 2 {
			overlap++
		}
	}
	return overlap, nil
}

func (d Day5) Part1(input string) (string, error) {
	overlap, err := countOverlaps(input, false)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", overlap), nil
}

func (d Day5) Part2(input string) (string, error) {
	overlap, err := countOverlaps(input, true)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", overlap), nil
}

func init() {
	solve.Register(Day5{})
}
