package solve2017

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day23 struct {
}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 23}
}

func runProgram(lines []string, registers map[string]int, onMul func()) map[string]int {
	getVal := func(x string) int {
		if v, err := strconv.Atoi(x); err == nil {
			return v
		}
		return registers[x]
	}

	for i := 0; i < len(lines); {
		parts := strings.Fields(lines[i])
		switch parts[0] {
		case "set":
			registers[parts[1]] = getVal(parts[2])
		case "sub":
			registers[parts[1]] -= getVal(parts[2])
		case "mul":
			registers[parts[1]] *= getVal(parts[2])
			if onMul != nil {
				onMul()
			}
		case "jnz":
			if getVal(parts[1]) != 0 {
				i += getVal(parts[2])
				continue
			}
		}
		i++
	}
	return registers
}

func (d Day23) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	mulCount := 0
	registers := make(map[string]int)
	runProgram(lines, registers, func() { mulCount++ })
	return strconv.Itoa(mulCount), nil
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func (d Day23) Part2(input string) (string, error) {
	b := 105700
	c := 122700
	step := 17
	h := 0
	for n := b; n <= c; n += step {
		if !isPrime(n) {
			h++
		}
	}
	return strconv.Itoa(h), nil
}

func init() {
	solve.Register(Day23{})
}
