package solve2018

import (
	"aoc/solve"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Day21 struct{}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 21}
}

func parseInstructions(input string) ([][4]int, int) {
	lines := strings.Split(input, "\n")
	ipRegister := 0
	fmt.Sscanf(lines[0], "#ip %d", &ipRegister)
	instructions := [][4]int{}

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		var op string
		var a, b, c int
		fmt.Sscanf(line, "%s %d %d %d", &op, &a, &b, &c)
		instructions = append(instructions, [4]int{opCode(op), a, b, c})
	}

	return instructions, ipRegister
}

func opCode(op string) int {
	switch op {
	case "addr":
		return 0
	case "addi":
		return 1
	case "mulr":
		return 2
	case "muli":
		return 3
	case "banr":
		return 4
	case "bani":
		return 5
	case "borr":
		return 6
	case "bori":
		return 7
	case "setr":
		return 8
	case "seti":
		return 9
	case "gtir":
		return 10
	case "gtri":
		return 11
	case "gtrr":
		return 12
	case "eqir":
		return 13
	case "eqri":
		return 14
	case "eqrr":
		return 15
	default:
		panic("unknown operation")
	}
}

func (d Day21) Part1(input string) (string, error) {
	instructions, ipRegister := parseInstructions(input)
	registers := [6]int{}
	seen := map[int]bool{}
	var firstHalt int
	foundFirst := false

	for {
		if registers[ipRegister] < 0 || registers[ipRegister] >= len(instructions) {
			break
		}
		inst := instructions[registers[ipRegister]]

		// Focus on the critical instruction (eqrr)
		if registers[ipRegister] == 28 { // Assuming instruction 28 is the critical one
			targetValue := registers[3] // Assuming R5 is the target register
			if !foundFirst {
				firstHalt = targetValue
				foundFirst = true
			}
			if _, exists := seen[targetValue]; exists {
				return strconv.Itoa(firstHalt), nil
			}
			seen[targetValue] = true
		}

		// Execute the instruction
		switch inst[0] {
		case 0: // addr
			registers[inst[3]] = registers[inst[1]] + registers[inst[2]]
		case 1: // addi
			registers[inst[3]] = registers[inst[1]] + inst[2]
		case 2: // mulr
			registers[inst[3]] = registers[inst[1]] * registers[inst[2]]
		case 3: // muli
			registers[inst[3]] = registers[inst[1]] * inst[2]
		case 4: // banr
			registers[inst[3]] = registers[inst[1]] & registers[inst[2]]
		case 5: // bani
			registers[inst[3]] = registers[inst[1]] & inst[2]
		case 6: // borr
			registers[inst[3]] = registers[inst[1]] | registers[inst[2]]
		case 7: // bori
			registers[inst[3]] = registers[inst[1]] | inst[2]
		case 8: // setr
			registers[inst[3]] = registers[inst[1]]
		case 9: // seti
			registers[inst[3]] = inst[1]
		case 10: // gtir
			if inst[1] > registers[inst[2]] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 11: // gtri
			if registers[inst[1]] > inst[2] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 12: // gtrr
			if registers[inst[1]] > registers[inst[2]] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 13: // eqir
			if inst[1] == registers[inst[2]] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 14: // eqri
			if registers[inst[1]] == inst[2] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 15: // eqrr
			if registers[inst[1]] == registers[inst[2]] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		}

		registers[ipRegister]++
	}

	return "", errors.New("No halting condition met")
}

func (d Day21) Part2(input string) (string, error) {
	instructions, ipRegister := parseInstructions(input)
	registers := [6]int{}
	seen := map[int]bool{}
	var lastHalt int

	for {
		if registers[ipRegister] < 0 || registers[ipRegister] >= len(instructions) {
			break
		}
		inst := instructions[registers[ipRegister]]

		// Focus on the critical instruction (eqrr)
		if registers[ipRegister] == 28 { // Assuming instruction 28 is the critical one
			targetValue := registers[3] // Assuming R3 is the target register
			if _, exists := seen[targetValue]; exists {
				return strconv.Itoa(lastHalt), nil
			}
			seen[targetValue] = true
			lastHalt = targetValue
		}

		// Execute the instruction
		switch inst[0] {
		case 0: // addr
			registers[inst[3]] = registers[inst[1]] + registers[inst[2]]
		case 1: // addi
			registers[inst[3]] = registers[inst[1]] + inst[2]
		case 2: // mulr
			registers[inst[3]] = registers[inst[1]] * registers[inst[2]]
		case 3: // muli
			registers[inst[3]] = registers[inst[1]] * inst[2]
		case 4: // banr
			registers[inst[3]] = registers[inst[1]] & registers[inst[2]]
		case 5: // bani
			registers[inst[3]] = registers[inst[1]] & inst[2]
		case 6: // borr
			registers[inst[3]] = registers[inst[1]] | registers[inst[2]]
		case 7: // bori
			registers[inst[3]] = registers[inst[1]] | inst[2]
		case 8: // setr
			registers[inst[3]] = registers[inst[1]]
		case 9: // seti
			registers[inst[3]] = inst[1]
		case 10: // gtir
			if inst[1] > registers[inst[2]] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 11: // gtri
			if registers[inst[1]] > inst[2] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 12: // gtrr
			if registers[inst[1]] > registers[inst[2]] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 13: // eqir
			if inst[1] == registers[inst[2]] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 14: // eqri
			if registers[inst[1]] == inst[2] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		case 15: // eqrr
			if registers[inst[1]] == registers[inst[2]] {
				registers[inst[3]] = 1
			} else {
				registers[inst[3]] = 0
			}
		}

		registers[ipRegister]++
	}

	return "", errors.New("No halting condition met")
}

func init() {
	solve.Register(Day21{})
}
