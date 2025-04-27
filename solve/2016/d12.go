package solve2016

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day12 struct {
}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 12}
}

type Instruction struct {
	Op string
	X  string
	Y  string
}

func (d Day12) parseInput(input string) []Instruction {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var instructions []Instruction
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 {
			instructions = append(instructions, Instruction{Op: parts[0], X: parts[1]})
		} else if len(parts) == 3 {
			instructions = append(instructions, Instruction{Op: parts[0], X: parts[1], Y: parts[2]})
		}
	}
	return instructions
}

func (d Day12) interpret(code []Instruction, regs map[string]int) map[string]int {
	val := func(x string) int {
		if v, err := strconv.Atoi(x); err == nil {
			return v
		}
		return regs[x]
	}

	pc := 0
	for pc < len(code) {
		inst := code[pc]
		switch inst.Op {
		case "cpy":
			regs[inst.Y] = val(inst.X)
		case "inc":
			regs[inst.X]++
		case "dec":
			regs[inst.X]--
		case "jnz":
			if val(inst.X) != 0 {
				pc += val(inst.Y) - 1
			}
		}
		pc++
	}
	return regs
}

func (d Day12) Part1(input string) (string, error) {
	code := d.parseInput(input)
	regs := d.interpret(code, map[string]int{"a": 0, "b": 0, "c": 0, "d": 0})
	return strconv.Itoa(regs["a"]), nil
}

func (d Day12) Part2(input string) (string, error) {
	code := d.parseInput(input)
	regs := d.interpret(code, map[string]int{"a": 0, "b": 0, "c": 1, "d": 0})
	return strconv.Itoa(regs["a"]), nil
}

func init() {
	solve.Register(Day12{})
}
