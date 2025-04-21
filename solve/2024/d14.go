package solve2024

import (
	"errors"
	"fmt"
	"image"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day14 struct {
}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 14}
}

type robot struct {
	Position, Velocity image.Point
}

func parseRobots(input string, area image.Rectangle) ([]robot, map[image.Point]int) {
	var robots []robot
	quadrants := map[image.Point]int{}

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		var r robot
		_, err := fmt.Sscanf(line, "p=%d,%d v=%d,%d",
			&r.Position.X, &r.Position.Y, &r.Velocity.X, &r.Velocity.Y)
		if err != nil {
			continue
		}
		robots = append(robots, r)
		r.Position = r.Position.Add(r.Velocity.Mul(100)).Mod(area)
		quadrant := image.Point{
			X: maths.Sign(r.Position.X - area.Dx()/2),
			Y: maths.Sign(r.Position.Y - area.Dy()/2),
		}
		quadrants[quadrant]++
	}

	return robots, quadrants
}

func moveRobots(robots []robot, area image.Rectangle) map[image.Point]struct{} {
	seen := map[image.Point]struct{}{}
	for i := range robots {
		robots[i].Position = robots[i].Position.Add(robots[i].Velocity).Mod(area)
		seen[robots[i].Position] = struct{}{}
	}
	return seen
}

func (d Day14) Part1(input string) (string, error) {
	area := image.Rectangle{Min: image.Point{}, Max: image.Point{X: 101, Y: 103}}
	_, quadrants := parseRobots(input, area)

	result := quadrants[image.Point{X: -1, Y: -1}] *
		quadrants[image.Point{X: 1, Y: -1}] *
		quadrants[image.Point{X: 1, Y: 1}] *
		quadrants[image.Point{X: -1, Y: 1}]

	return strconv.Itoa(result), nil
}

func (d Day14) Part2(input string) (string, error) {
	area := image.Rectangle{Min: image.Point{}, Max: image.Point{X: 101, Y: 103}}

	robots, _ := parseRobots(input, area)

	for steps := 1; steps <= 10000; steps++ {
		seen := moveRobots(robots, area)
		if len(seen) == len(robots) {
			return strconv.Itoa(steps), nil
		}
	}

	return "", errors.New("no solution found after 10000 steps")
}

func init() {
	solve.Register(Day14{})
}
