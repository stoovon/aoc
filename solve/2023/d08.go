package solve2023

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day8 struct {
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 8}
}

type day8Node struct {
	left, right string
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a * (b / gcd(a, b))
}

func (d Day8) parseInput(input string) (string, map[string]day8Node) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	instructions := lines[0]
	nodes := make(map[string]day8Node)
	for _, line := range lines[2:] {
		name := line[0:3]
		left := line[7:10]
		right := line[12:15]
		nodes[name] = day8Node{left, right}
	}
	return instructions, nodes
}

func (d Day8) Part1(input string) (string, error) {
	instructions, nodes := d.parseInput(input)
	curr := "AAA"
	steps := 0
	for {
		for _, dir := range instructions {
			if curr == "ZZZ" {
				return strconv.Itoa(steps), nil
			}
			if dir == 'L' {
				curr = nodes[curr].left
			} else {
				curr = nodes[curr].right
			}
			steps++
		}
	}
}

func (d Day8) Part2(input string) (string, error) {
	instructions, nodes := d.parseInput(input)
	var starts []string
	for k := range nodes {
		if strings.HasSuffix(k, "A") {
			starts = append(starts, k)
		}
	}
	cycles := make([]int64, len(starts))
	for i, start := range starts {
		curr := start
		steps := 0
		for {
			for _, dir := range instructions {
				if strings.HasSuffix(curr, "Z") {
					cycles[i] = int64(steps)
					goto next
				}
				if dir == 'L' {
					curr = nodes[curr].left
				} else {
					curr = nodes[curr].right
				}
				steps++
			}
		}
	next:
	}
	res := cycles[0]
	for _, c := range cycles[1:] {
		res = lcm(res, c)
	}
	return strconv.FormatInt(res, 10), nil
}

func init() {
	solve.Register(Day8{})
}
