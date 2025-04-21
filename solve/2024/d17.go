package solve2024

import (
	"slices"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day17 struct {
}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 17}
}

const (
	InsAdv int = iota
	InsBxl
	InsBst
	InsJnz
	InsBxc
	InsOut
	InsBdv
	InsCdv
)

type d17Program struct {
	data []int
	regA int
	regB int
	regC int
}

func (prog *d17Program) Execute() []int {
	insPtr := 0
	output := make([]int, 0, len(prog.data))

	for insPtr < len(prog.data) {
		instruction := prog.data[insPtr]
		literal := prog.data[insPtr+1]
		combo := comboOperand(literal, prog.regA, prog.regB, prog.regC)

		switch instruction {
		case InsAdv:
			prog.regA >>= combo
		case InsBxl:
			prog.regB ^= literal
		case InsBst:
			prog.regB = combo % 8
		case InsJnz:
			if prog.regA != 0 {
				insPtr = literal - 2
			}
		case InsBxc:
			prog.regB = prog.regB ^ prog.regC
		case InsOut:
			output = append(output, combo%8)
		case InsBdv:
			prog.regB = prog.regA >> combo
		case InsCdv:
			prog.regC = prog.regA >> combo
		}
		insPtr += 2
	}

	return output
}

func comboOperand(operand int, regA int, regB int, regC int) int {
	switch {
	case operand >= 0 && operand <= 3:
		return operand
	case operand == 4:
		return regA
	case operand == 5:
		return regB
	case operand == 6:
		return regC
	default:
		panic("Unexpected operand")
	}
}

func (d Day17) mustAtoi(input string) int {
	result, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	return result
}

func (d Day17) intArrayToString(input []int) string {
	result := make([]string, len(input))
	for i, v := range input {
		result[i] = strconv.Itoa(v)
	}
	return strings.Join(result, ",")
}

func (d Day17) intsInString(input string) []int {
	splits := strings.Split(input, ",")
	result := make([]int, len(splits))
	for i, v := range splits {
		result[i] = d.mustAtoi(v)
	}

	return result
}

func (d Day17) Part1(input string) (string, error) {
	program := d.parseRegisters(input)

	output := program.Execute()

	return d.intArrayToString(output), nil
}

func (d Day17) Part2(input string) (string, error) {
	program := d.parseRegisters(input)

	candidateA := 0
	// Reduce search space by starting with a small A value.
	// Shift 3 bits left once we find it (can still change later if necessary); this is because
	// the algo we're reverse engineering seems to use the last 3 bits of A to determine the output value,
	// LSB corresponding to first such value.
	// Increment to find e.g. first and second values. Then loop until we find all values.

	// This approach cuts the number of iterations considerably.
	for pos := len(program.data) - 1; pos >= 0; pos-- {
		candidateA <<= 3
		for {
			program.regA = candidateA
			output := program.Execute()
			if slices.Equal(output, program.data[pos:]) {
				break
			}
			candidateA++
		}
	}

	return strconv.Itoa(candidateA), nil
}

func (d Day17) parseRegisters(input string) d17Program {
	a := 0
	b := 0
	c := 0
	program := make([]int, 0)

	for _, line := range strings.Split(input, "\n") {
		if strings.HasPrefix(line, "Register A:") {
			a = d.mustAtoi(strings.TrimSpace(line[12:]))
		} else if strings.HasPrefix(line, "Register B:") {
			b = d.mustAtoi(strings.TrimSpace(line[12:]))
		} else if strings.HasPrefix(line, "Register C:") {
			c = d.mustAtoi(strings.TrimSpace(line[12:]))
		} else if strings.HasPrefix(line, "Program:") {
			program = d.intsInString(strings.TrimSpace(line[9:]))
		}
	}

	return d17Program{
		data: program,
		regA: a,
		regB: b,
		regC: c,
	}
}

func init() {
	solve.Register(Day17{})
}
