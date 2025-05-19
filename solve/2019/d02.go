package solve2019

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 2}
}

func (d Day2) parseIntcode(input string) ([]int, error) {
	parts := strings.Split(strings.TrimSpace(input), ",")
	program := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid intcode at position %d: %w", i, err)
		}
		program[i] = n
	}
	return program, nil
}

func runIntcode(prog []int) {
	ip := 0
	for {
		switch prog[ip] {
		case 1:
			a, b, c := prog[ip+1], prog[ip+2], prog[ip+3]
			prog[c] = prog[a] + prog[b]
			ip += 4
		case 2:
			a, b, c := prog[ip+1], prog[ip+2], prog[ip+3]
			prog[c] = prog[a] * prog[b]
			ip += 4
		case 99:
			return
		default:
			panic(fmt.Sprintf("unknown opcode %d at %d", prog[ip], ip))
		}
	}
}

func (d Day2) Part1(input string) (string, error) {
	prog, err := d.parseIntcode(input)
	if err != nil {
		return "", err
	}
	prog[1] = 12
	prog[2] = 2
	runIntcode(prog)
	return strconv.Itoa(prog[0]), nil
}

func (d Day2) Part2(input string) (string, error) {
	orig, err := d.parseIntcode(input)
	if err != nil {
		return "", err
	}
	const target = 19690720
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			prog := make([]int, len(orig))
			copy(prog, orig)
			prog[1] = noun
			prog[2] = verb
			runIntcode(prog)
			if prog[0] == target {
				return strconv.Itoa(100*noun + verb), nil
			}
		}
	}
	return "", fmt.Errorf("no noun/verb found")
}

func init() {
	solve.Register(Day2{})
}
