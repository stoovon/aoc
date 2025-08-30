package solve2022

import (
	"aoc/solve"
	"fmt"
	"strconv"
	"strings"
)

type Day9 struct{}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 9}
}

func (d Day9) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	head := [2]int{0, 0}
	tail := [2]int{0, 0}
	visited := make(map[[2]int]bool)
	visited[tail] = true

	for _, motion := range lines {
		direction := motion[:1]
		steps, err := strconv.Atoi(motion[2:])
		if err != nil {
			return "", fmt.Errorf("invalid input: %v", err)
		}

		for range steps {
			switch direction {
			case "U":
				head[1]++
			case "D":
				head[1]--
			case "R":
				head[0]++
			case "L":
				head[0]--
			}

			diffX := head[0] - tail[0]
			diffY := head[1] - tail[1]
			if abs(diffX) > 1 || abs(diffY) > 1 {
				tail[0] += sign(diffX)
				tail[1] += sign(diffY)
				visited[tail] = true
			}
		}
	}

	return strconv.Itoa(len(visited)), nil
}

func (d Day9) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	const ropeLength = 10
	rope := make([][2]int, ropeLength)
	visited := make(map[[2]int]bool)
	visited[rope[ropeLength-1]] = true

	for _, motion := range lines {
		direction := motion[:1]
		steps, err := strconv.Atoi(motion[2:])
		if err != nil {
			return "", fmt.Errorf("invalid input: %v", err)
		}

		for i := 0; i < steps; i++ {
			switch direction {
			case "U":
				rope[0][1]++
			case "D":
				rope[0][1]--
			case "R":
				rope[0][0]++
			case "L":
				rope[0][0]--
			}

			for j := 1; j < ropeLength; j++ {
				diffX := rope[j-1][0] - rope[j][0]
				diffY := rope[j-1][1] - rope[j][1]
				if abs(diffX) > 1 || abs(diffY) > 1 {
					rope[j][0] += sign(diffX)
					rope[j][1] += sign(diffY)
				}
			}

			visited[rope[ropeLength-1]] = true
		}
	}

	return strconv.Itoa(len(visited)), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

func init() {
	solve.Register(Day9{})
}
