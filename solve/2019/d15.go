package solve2019

import (
	"errors"
	"fmt"
	"image"

	"aoc/solve"
)

type Day15 struct {
}

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 15}
}

var dirs = []struct {
	d   image.Point
	cmd int64
}{
	{image.Pt(0, -1), 1}, // North
	{image.Pt(0, 1), 2},  // South
	{image.Pt(-1, 0), 3}, // West
	{image.Pt(1, 0), 4},  // East
}

func (d Day15) Part1(input string) (string, error) {
	type state struct {
		p     image.Point
		steps int
		code  *Intcode
	}
	start := image.Pt(0, 0)
	visited := map[image.Point]bool{start: true}
	queue := []state{{start, 0, parseIntcode(input)}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, dir := range dirs {
			np := cur.p.Add(dir.d)
			if visited[np] {
				continue
			}
			code := cur.code.Clone()
			outs, _ := code.Step([]int64{dir.cmd}, 1)
			status := outs[0]
			if status == 0 {
				visited[np] = true
				continue
			}
			if status == 2 {
				return fmt.Sprint(cur.steps + 1), nil
			}
			visited[np] = true
			queue = append(queue, state{np, cur.steps + 1, code})
		}
	}
	return "", errors.New("Oxygen system not found")
}

// Explore the maze, return open tiles and oxygen system location
func exploreMaze(start *Intcode) (map[image.Point]bool, image.Point, error) {
	type state struct {
		p    image.Point
		code *Intcode
	}
	open := map[image.Point]bool{}
	var oxygen image.Point
	foundOxygen := false
	queue := []state{{image.Pt(0, 0), start.Clone()}}
	open[image.Pt(0, 0)] = true

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, dir := range dirs {
			np := cur.p.Add(dir.d)
			if open[np] { // already visited
				continue
			}
			code := cur.code.Clone()
			outs, _ := code.Step([]int64{dir.cmd}, 1)
			status := outs[0]
			if status == 0 {
				continue // wall
			}
			open[np] = true
			if status == 2 {
				oxygen = np
				foundOxygen = true
			}
			queue = append(queue, state{np, code})
		}
	}
	if !foundOxygen {
		return nil, image.Point{}, errors.New("oxygen system not found")
	}
	return open, oxygen, nil
}

// BFS from oxygen system to fill all open tiles, return minutes needed
func fillOxygen(open map[image.Point]bool, oxygen image.Point) int {
	filled := map[image.Point]bool{oxygen: true}
	queue := []image.Point{oxygen}
	minutes := 0

	for {
		next := []image.Point{}
		for _, p := range queue {
			for _, dir := range dirs {
				np := p.Add(dir.d)
				if open[np] && !filled[np] {
					filled[np] = true
					next = append(next, np)
				}
			}
		}
		if len(next) == 0 {
			break
		}
		queue = next
		minutes++
	}
	return minutes
}

func (d Day15) Part2(input string) (string, error) {
	code := parseIntcode(input)
	open, oxygen, err := exploreMaze(code)
	if err != nil {
		return "", err
	}
	minutes := fillOxygen(open, oxygen)
	return fmt.Sprint(minutes), nil
}

func init() {
	solve.Register(Day15{})
}
