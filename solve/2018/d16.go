package solve2018

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day16 struct{}

func addr(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a] + result[b]
	return result
}

func addi(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a] + b
	return result
}

func mulr(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a] * result[b]
	return result
}

func muli(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a] * b
	return result
}

func banr(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a] & result[b]
	return result
}

func bani(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a] & b
	return result
}

func borr(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a] | result[b]
	return result
}

func bori(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a] | b
	return result
}

func setr(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = result[a]
	return result
}

func seti(registers [4]int, a, b, c int) [4]int {
	result := registers
	result[c] = a
	return result
}

func gtir(registers [4]int, a, b, c int) [4]int {
	result := registers
	if a > result[b] {
		result[c] = 1
	} else {
		result[c] = 0
	}
	return result
}

func gtri(registers [4]int, a, b, c int) [4]int {
	result := registers
	if result[a] > b {
		result[c] = 1
	} else {
		result[c] = 0
	}
	return result
}

func gtrr(registers [4]int, a, b, c int) [4]int {
	result := registers
	if result[a] > result[b] {
		result[c] = 1
	} else {
		result[c] = 0
	}
	return result
}

func eqir(registers [4]int, a, b, c int) [4]int {
	result := registers
	if a == result[b] {
		result[c] = 1
	} else {
		result[c] = 0
	}
	return result
}

func eqri(registers [4]int, a, b, c int) [4]int {
	result := registers
	if result[a] == b {
		result[c] = 1
	} else {
		result[c] = 0
	}
	return result
}

func eqrr(registers [4]int, a, b, c int) [4]int {
	result := registers
	if result[a] == result[b] {
		result[c] = 1
	} else {
		result[c] = 0
	}
	return result
}

var operations = []func([4]int, int, int, int) [4]int{
	addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr,
}

func possibleOperations(instruction []int, before, after [4]int) int {
	count := 0
	for _, operation := range operations {
		if operation(before, instruction[1], instruction[2], instruction[3]) == after {
			count++
		}
	}
	return count
}

func (d Day16) Part1(input string) (string, error) {
	lines := strings.Split(input, "\n")
	count := 0
	for i := 0; i < len(lines); i += 4 {
		if lines[i] == "" {
			break
		}
		before := parseRegisters(lines[i][8:])
		instruction := parseInstruction(lines[i+1])
		after := parseRegisters(lines[i+2][8:])
		if possibleOperations(instruction, before, after) >= 3 {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func parseRegisters(s string) [4]int {
	var registers [4]int
	values := strings.Split(strings.Trim(s, "[]"), ", ")
	for i, v := range values {
		registers[i], _ = strconv.Atoi(v)
	}
	return registers
}

func parseInstruction(s string) []int {
	values := strings.Split(s, " ")
	instruction := make([]int, len(values))
	for i, v := range values {
		instruction[i], _ = strconv.Atoi(v)
	}
	return instruction
}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 16}
}

func (d Day16) Part2(input string) (string, error) {
	lines := strings.Split(input, "\n")
	// Parse samples to deduce opcode mappings
	experiments := []struct {
		instruction []int
		before      [4]int
		after       [4]int
	}{}
	i := 0
	for i < len(lines) && lines[i] != "" {
		before := parseRegisters(lines[i][8:])
		instruction := parseInstruction(lines[i+1])
		after := parseRegisters(lines[i+2][8:])
		experiments = append(experiments, struct {
			instruction []int
			before      [4]int
			after       [4]int
		}{instruction, before, after})
		i += 4
	}

	// Deduce opcode mappings
	opcodeMappings := make(map[int]func([4]int, int, int, int) [4]int)
	remainingOperations := make(map[int]map[int]bool)
	for _, experiment := range experiments {
		opcode := experiment.instruction[0]
		if _, exists := remainingOperations[opcode]; !exists {
			remainingOperations[opcode] = make(map[int]bool)
			for idx := range operations {
				remainingOperations[opcode][idx] = true
			}
		}
		for idx, op := range operations {
			if remainingOperations[opcode][idx] && op(experiment.before, experiment.instruction[1], experiment.instruction[2], experiment.instruction[3]) != experiment.after {
				delete(remainingOperations[opcode], idx)
			}
		}
	}

	resolved := true
	for resolved {
		resolved = false
		for opcode, ops := range remainingOperations {
			if len(ops) == 1 {
				for idx := range ops {
					opcodeMappings[opcode] = operations[idx]
					delete(remainingOperations, opcode)
					for _, otherOps := range remainingOperations {
						delete(otherOps, idx)
					}
					resolved = true
				}
			}
		}
	}

	// Execute the test program
	registers := [4]int{}
	for i++; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		instruction := parseInstruction(lines[i])
		registers = opcodeMappings[instruction[0]](registers, instruction[1], instruction[2], instruction[3])
	}

	return strconv.Itoa(registers[0]), nil
}

func init() {
	solve.Register(Day16{})
}
