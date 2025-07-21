package solve2018

import (
	"aoc/solve"
	"fmt"
	"strconv"
	"strings"
)

type Day25 struct{}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 25}
}

func (d Day25) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	points := make([][4]int, 0, len(lines))
	for _, line := range lines {
		fields := strings.Split(line, ",")
		var pt [4]int
		for i := 0; i < 4; i++ {
			n, _ := strconv.Atoi(strings.TrimSpace(fields[i]))
			pt[i] = n
		}
		points = append(points, pt)
	}

	constellations := [][][4]int{}
	for _, pt := range points {
		joined := []int{}
		for i, c := range constellations {
			for _, other := range c {
				if manhattan4d(pt, other) <= 3 {
					joined = append(joined, i)
					break
				}
			}
		}
		if len(joined) == 0 {
			constellations = append(constellations, [][4]int{pt})
		} else {
			// Merge all joined constellations
			merged := [][4]int{pt}
			// Add all points from joined constellations
			for i := len(joined) - 1; i >= 0; i-- {
				idx := joined[i]
				merged = append(merged, constellations[idx]...)
				constellations = append(constellations[:idx], constellations[idx+1:]...)
			}
			constellations = append(constellations, merged)
		}
	}
	return fmt.Sprintf("%d", len(constellations)), nil
}

func manhattan4d(a, b [4]int) int {
	sum := 0
	for i := 0; i < 4; i++ {
		d := a[i] - b[i]
		if d < 0 {
			d = -d
		}
		sum += d
	}
	return sum
}

func (d Day25) Part2(input string) (string, error) {
	return "", fmt.Errorf("No Part 2 for Day 25!")
}

func init() {
	solve.Register(Day25{})
}
