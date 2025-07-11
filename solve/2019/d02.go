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

func (d Day2) Part1(input string) (string, error) {
	prog, err := parseIntcodeV1(input)
	if err != nil {
		return "", err
	}
	prog[1] = 12
	prog[2] = 2
	runIntcodeV1(prog)
	return strconv.Itoa(prog[0]), nil
}

func (d Day2) Part2(input string) (string, error) {
	orig, err := parseIntcodeV1(input)
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
			runIntcodeV1(prog)
			if prog[0] == target {
				return strconv.Itoa(100*noun + verb), nil
			}
		}
	}
	return "", fmt.Errorf("no noun/verb found")
}

func parseIntcodeV1(input string) ([]int, error) {
	parts := strings.Split(strings.TrimSpace(input), ",")
	res := make([]int, len(parts))
	for i, s := range parts {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

// (no parameter modes, no input/output)
func runIntcodeV1(mem []int) {
	ip := 0
	for {
		switch mem[ip] {
		case 1:
			a, b, c := mem[ip+1], mem[ip+2], mem[ip+3]
			mem[c] = mem[a] + mem[b]
			ip += 4
		case 2:
			a, b, c := mem[ip+1], mem[ip+2], mem[ip+3]
			mem[c] = mem[a] * mem[b]
			ip += 4
		case 99:
			return
		default:
			panic(fmt.Sprintf("unknown opcode %d at %d", mem[ip], ip))
		}
	}
}

func init() {
	solve.Register(Day2{})
}
