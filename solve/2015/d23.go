package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day23 struct {
}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 23}
}

func (d Day23) execute(instructions []string, registers map[string]int) map[string]int {
	pc := 0

	for pc >= 0 && pc < len(instructions) {
		line := instructions[pc]
		tokens := strings.Fields(line)

		switch tokens[0] {
		case "hlf":
			reg := tokens[1]
			registers[reg] /= 2
			pc++
		case "tpl":
			reg := tokens[1]
			registers[reg] *= 3
			pc++
		case "inc":
			reg := tokens[1]
			registers[reg]++
			pc++
		case "jmp":
			dst, _ := strconv.Atoi(tokens[1])
			pc += dst
		case "jie":
			reg := strings.TrimSuffix(tokens[1], ",")
			dst, _ := strconv.Atoi(tokens[2])
			if registers[reg]%2 == 0 {
				pc += dst
			} else {
				pc++
			}
		case "jio":
			reg := strings.TrimSuffix(tokens[1], ",")
			dst, _ := strconv.Atoi(tokens[2])
			if registers[reg] == 1 {
				pc += dst
			} else {
				pc++
			}
		default:
			panic("Unknown instruction: " + line)
		}
	}

	return registers
}

func (d Day23) Part1(input string) (string, error) {
	instructions := strings.Split(strings.TrimSpace(input), "\n")
	registers := d.execute(instructions, map[string]int{"a": 0, "b": 0})
	return strconv.Itoa(registers["b"]), nil
}

func (d Day23) Part2(input string) (string, error) {
	instructions := strings.Split(strings.TrimSpace(input), "\n")
	registers := d.execute(instructions, map[string]int{"a": 1, "b": 0})
	return strconv.Itoa(registers["b"]), nil
}

func init() {
	solve.Register(Day23{})
}
