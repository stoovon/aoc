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

func bfsExplore(
	startCode *Intcode,
	process func(p image.Point, steps int, status int64, code *Intcode) (stop bool),
) (found bool, stopPoint image.Point, stopSteps int, stopCode *Intcode) {
	type state struct {
		p     image.Point
		steps int
		code  *Intcode
	}
	start := image.Pt(0, 0)
	visited := map[image.Point]bool{start: true}
	queue := []state{{start, 0, startCode.Clone()}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, dir := range dirs {
			np := cur.p.Add(dir.d)
			if visited[np] {
				continue
			}
			code := cur.code.Clone()
			inputVal := dir.cmd
			var status int64
			for {
				res := code.Step(&inputVal)
				inputVal = 0
				if res.Halted {
					break
				}
				if res.NeedInput {
					break
				}
				if res.Output != nil {
					status = *res.Output
					break
				}
			}
			visited[np] = true
			if stop := process(np, cur.steps+1, status, code); stop {
				return true, np, cur.steps + 1, code
			}
			if status != 0 {
				queue = append(queue, state{np, cur.steps + 1, code})
			}
		}
	}
	return false, image.Point{}, 0, nil
}

func (d Day15) Part1(input string) (string, error) {
	found, _, steps, _ := bfsExplore(parseIntcode(input), func(p image.Point, steps int, status int64, code *Intcode) bool {
		return status == 2 // Stop when oxygen system found
	})
	if !found {
		return "", errors.New("Oxygen system not found")
	}
	return fmt.Sprint(steps), nil
}

// Explore the maze, return open tiles and oxygen system location
func exploreMaze(start *Intcode) (map[image.Point]bool, image.Point, error) {
	open := map[image.Point]bool{image.Pt(0, 0): true}
	var oxygen image.Point
	found, _, _, _ := bfsExplore(start, func(p image.Point, steps int, status int64, code *Intcode) bool {
		if status != 0 {
			open[p] = true
		}
		if status == 2 {
			oxygen = p
			return false // Continue to fill open map
		}
		return false
	})
	if !found && oxygen == (image.Point{}) {
		return nil, image.Point{}, errors.New("Oxygen system not found")
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
