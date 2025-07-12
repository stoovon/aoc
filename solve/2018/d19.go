package solve2018

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day19 struct{}

func (d Day19) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 19}
}

type Instruction struct {
	opcode  string
	a, b, c int
}

type CPU struct {
	registers [6]int
	ip        int
	ipReg     int
}

func (cpu *CPU) execute(inst Instruction) {
	cpu.registers[cpu.ipReg] = cpu.ip

	switch inst.opcode {
	case "addr":
		cpu.registers[inst.c] = cpu.registers[inst.a] + cpu.registers[inst.b]
	case "addi":
		cpu.registers[inst.c] = cpu.registers[inst.a] + inst.b

	case "mulr":
		cpu.registers[inst.c] = cpu.registers[inst.a] * cpu.registers[inst.b]
	case "muli":
		cpu.registers[inst.c] = cpu.registers[inst.a] * inst.b

	case "banr":
		cpu.registers[inst.c] = cpu.registers[inst.a] & cpu.registers[inst.b]
	case "bani":
		cpu.registers[inst.c] = cpu.registers[inst.a] & inst.b

	case "borr":
		cpu.registers[inst.c] = cpu.registers[inst.a] | cpu.registers[inst.b]
	case "bori":
		cpu.registers[inst.c] = cpu.registers[inst.a] | inst.b

	case "setr":
		cpu.registers[inst.c] = cpu.registers[inst.a]
	case "seti":
		cpu.registers[inst.c] = inst.a

	case "gtir":
		if inst.a > cpu.registers[inst.b] {
			cpu.registers[inst.c] = 1
		} else {
			cpu.registers[inst.c] = 0
		}
	case "gtri":
		if cpu.registers[inst.a] > inst.b {
			cpu.registers[inst.c] = 1
		} else {
			cpu.registers[inst.c] = 0
		}
	case "gtrr":
		if cpu.registers[inst.a] > cpu.registers[inst.b] {
			cpu.registers[inst.c] = 1
		} else {
			cpu.registers[inst.c] = 0
		}

	case "eqir":
		if inst.a == cpu.registers[inst.b] {
			cpu.registers[inst.c] = 1
		} else {
			cpu.registers[inst.c] = 0
		}
	case "eqri":
		if cpu.registers[inst.a] == inst.b {
			cpu.registers[inst.c] = 1
		} else {
			cpu.registers[inst.c] = 0
		}
	case "eqrr":
		if cpu.registers[inst.a] == cpu.registers[inst.b] {
			cpu.registers[inst.c] = 1
		} else {
			cpu.registers[inst.c] = 0
		}
	}

	cpu.ip = cpu.registers[cpu.ipReg]

	cpu.ip++
}

func parseProgram(input string) (int, []Instruction) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	ipReg := 0
	if strings.HasPrefix(lines[0], "#ip ") {
		ipReg, _ = strconv.Atoi(strings.TrimPrefix(lines[0], "#ip "))
		lines = lines[1:]
	}

	var instructions []Instruction
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 4 {
			a, _ := strconv.Atoi(parts[1])
			b, _ := strconv.Atoi(parts[2])
			c, _ := strconv.Atoi(parts[3])
			instructions = append(instructions, Instruction{
				opcode: parts[0],
				a:      a, b: b, c: c,
			})
		}
	}

	return ipReg, instructions
}

func runProgram(input string, initialReg0 int) (string, error) {
	ipReg, instructions := parseProgram(input)

	cpu := &CPU{
		registers: [6]int{initialReg0, 0, 0, 0, 0, 0},
		ip:        0,
		ipReg:     ipReg,
	}

	// For Part 2, we need to optimize - the program is calculating sum of divisors
	if initialReg0 == 1 {
		// Let the program run until it calculates the target number
		cycles := 0
		for cpu.ip >= 0 && cpu.ip < len(instructions) && cycles < 1000 {
			cpu.execute(instructions[cpu.ip])
			cycles++
		}

		// At this point, one of the registers should contain the target number
		// The program is calculating the sum of all divisors of that number
		target := 0
		for _, reg := range cpu.registers {
			if reg > target {
				target = reg
			}
		}

		if target > 1 {
			sum := 0
			for i := 1; i <= target; i++ {
				if target%i == 0 {
					sum += i
				}
			}
			return strconv.Itoa(sum), nil
		}
	}

	for cpu.ip >= 0 && cpu.ip < len(instructions) {
		cpu.execute(instructions[cpu.ip])
	}

	return strconv.Itoa(cpu.registers[0]), nil
}

func (d Day19) Part1(input string) (string, error) {
	return runProgram(input, 0)
}

func (d Day19) Part2(input string) (string, error) {
	return runProgram(input, 1)
}

func init() {
	solve.Register(Day19{})
}
