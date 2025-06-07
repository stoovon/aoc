package solve2019

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 5}
}

func parseIntCode(input string) ([]int, error) {
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

func runIntCode(mem []int, inputVal int) (int, error) {
	mem = append([]int(nil), mem...) // copy to avoid mutating input
	ip := 0
	var output int

	get := func(mode, val int) int {
		if mode == 0 {
			return mem[val]
		}
		return val
	}

	for {
		op := mem[ip] % 100
		mode1 := (mem[ip] / 100) % 10
		mode2 := (mem[ip] / 1000) % 10

		switch op {
		case 1: // add
			a := get(mode1, mem[ip+1])
			b := get(mode2, mem[ip+2])
			mem[mem[ip+3]] = a + b
			ip += 4
		case 2: // multiply
			a := get(mode1, mem[ip+1])
			b := get(mode2, mem[ip+2])
			mem[mem[ip+3]] = a * b
			ip += 4
		case 3: // input
			mem[mem[ip+1]] = inputVal
			ip += 2
		case 4: // output
			output = get(mode1, mem[ip+1])
			ip += 2
		case 5: // jump-if-true
			if get(mode1, mem[ip+1]) != 0 {
				ip = get(mode2, mem[ip+2])
			} else {
				ip += 3
			}
		case 6: // jump-if-false
			if get(mode1, mem[ip+1]) == 0 {
				ip = get(mode2, mem[ip+2])
			} else {
				ip += 3
			}
		case 7: // less than
			if get(mode1, mem[ip+1]) < get(mode2, mem[ip+2]) {
				mem[mem[ip+3]] = 1
			} else {
				mem[mem[ip+3]] = 0
			}
			ip += 4
		case 8: // equals
			if get(mode1, mem[ip+1]) == get(mode2, mem[ip+2]) {
				mem[mem[ip+3]] = 1
			} else {
				mem[mem[ip+3]] = 0
			}
			ip += 4
		case 99:
			return output, nil
		default:
			return 0, fmt.Errorf("unknown opcode %d at %d", op, ip)
		}
	}
}

func (d Day5) Part1(input string) (string, error) {
	prog, err := parseIntCode(input)
	if err != nil {
		return "", err
	}
	output, err := runIntCode(prog, 1)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", output), nil
}

func (d Day5) Part2(input string) (string, error) {
	prog, err := parseIntCode(input)
	if err != nil {
		return "", err
	}
	output, err := runIntCode(prog, 5)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", output), nil
}

func init() {
	solve.Register(Day5{})
}
