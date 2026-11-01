package solve2018

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day21 struct{}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 21}
}

var chronalOpcodes = map[string]int{
	"addr": 0, "addi": 1, "mulr": 2, "muli": 3, "banr": 4, "bani": 5, "borr": 6, "bori": 7,
	"setr": 8, "seti": 9, "gtir": 10, "gtri": 11, "gtrr": 12, "eqir": 13, "eqri": 14, "eqrr": 15,
}

func opCode(name string) int { return chronalOpcodes[name] }

func chronalRegisterOperands(op int) (a, b bool) {
	switch op {
	case 0, 2, 4, 6, 12, 15:
		return true, true
	case 1, 3, 5, 7, 8, 11, 14:
		return true, false
	case 10, 13:
		return false, true
	default:
		return false, false
	}
}

func parseInstructions(input string) ([][4]int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	header := strings.Fields(lines[0])
	if len(header) != 2 || header[0] != "#ip" {
		return nil, 0, errors.New("expected instruction-pointer binding")
	}
	ip, err := strconv.Atoi(header[1])
	if err != nil || ip < 0 || ip >= 6 {
		return nil, 0, errors.New("invalid instruction-pointer register")
	}
	var program [][4]int
	for line, text := range lines[1:] {
		fields := strings.Fields(text)
		if len(fields) != 4 {
			return nil, 0, fmt.Errorf("line %d: expected opcode and three operands", line+2)
		}
		op, known := chronalOpcodes[fields[0]]
		if !known {
			return nil, 0, fmt.Errorf("line %d: unknown opcode", line+2)
		}
		inst := [4]int{op}
		for i := 1; i < 4; i++ {
			value, err := strconv.Atoi(fields[i])
			if err != nil {
				return nil, 0, fmt.Errorf("line %d: invalid int operand", line+2)
			}
			inst[i] = value
		}
		a, b := chronalRegisterOperands(op)
		if inst[3] < 0 || inst[3] >= 6 || (a && (inst[1] < 0 || inst[1] >= 6)) || (b && (inst[2] < 0 || inst[2] >= 6)) {
			return nil, 0, fmt.Errorf("line %d: register outside 0..5", line+2)
		}
		program = append(program, inst)
	}
	if len(program) == 0 {
		return nil, 0, errors.New("empty program")
	}
	return program, ip, nil
}

type chronalProbe struct{ ip, candidate, result int }

// Follow only unconditional jumps, stopping at computation, exit, or a jump cycle.
func chronalJumpEnd(program [][4]int, ipRegister, start int) (int, bool) {
	seen := make(map[int]bool)
	for start >= 0 && start < len(program) {
		if seen[start] {
			return 0, false
		}
		seen[start] = true
		inst := program[start]
		if inst[0] != opCode("seti") || inst[3] != ipRegister {
			return start, true
		}
		start = inst[1] + 1
	}
	return start, true
}

func findChronalProbe(program [][4]int, ip int) (chronalProbe, error) {
	var probe chronalProbe
	found := false
	if ip == 0 {
		return probe, errors.New("register zero must not bind the instruction pointer")
	}
	for index, inst := range program {
		if inst[0] != opCode("eqrr") || (inst[1] != 0 && inst[2] != 0) {
			continue
		}
		if found {
			return probe, errors.New("ambiguous register-zero comparison")
		}
		candidate := inst[1]
		if candidate == 0 {
			candidate = inst[2]
		}
		if candidate == 0 || candidate == ip || inst[3] == 0 || inst[3] == ip {
			return probe, errors.New("unsupported halting comparison registers")
		}
		probe = chronalProbe{index, candidate, inst[3]}
		found = true
	}
	if !found {
		return probe, errors.New("no register-zero halting comparison")
	}
	for index, inst := range program {
		a, b := chronalRegisterOperands(inst[0])
		if inst[3] == 0 || (index != probe.ip && ((a && inst[1] == 0) || (b && inst[2] == 0))) {
			return probe, errors.New("register zero is used outside the halting comparison")
		}
	}
	if probe.ip+2 >= len(program) {
		return probe, errors.New("missing halting comparison branch")
	}
	branch := program[probe.ip+1]
	if branch != [4]int{opCode("addr"), probe.result, ip, ip} && branch != [4]int{opCode("addr"), ip, probe.result, ip} {
		return probe, errors.New("unsupported halting comparison branch")
	}
	continuing := program[probe.ip+2]
	if continuing[0] != opCode("seti") || continuing[3] != ip {
		return probe, errors.New("inequality branch must jump into the program")
	}
	falseTarget, ok := chronalJumpEnd(program, ip, probe.ip+2)
	if !ok || falseTarget < 0 || falseTarget >= len(program) {
		return probe, errors.New("inequality branch does not return to computation")
	}
	trueTarget, ok := chronalJumpEnd(program, ip, probe.ip+3)
	if !ok || (trueTarget >= 0 && trueTarget < len(program)) {
		return probe, errors.New("equality branch is not a proven exit")
	}

	return probe, nil
}

type chronalMachine struct {
	instructions [][4]int
	ipRegister   int
	registers    [6]int
	ip           int
	divisions    map[int]chronalDivision
	probe        chronalProbe
}

func newChronalMachine(input string) (chronalMachine, error) {
	instructions, ipRegister, err := parseInstructions(input)
	if err != nil {
		return chronalMachine{}, err
	}
	probe, err := findChronalProbe(instructions, ipRegister)
	if err != nil {
		return chronalMachine{}, err
	}
	return chronalMachine{instructions: instructions, ipRegister: ipRegister, probe: probe, divisions: findChronalDivisions(instructions, ipRegister)}, nil
}
func (machine *chronalMachine) running() bool {
	return machine.ip >= 0 && machine.ip < len(machine.instructions)
}
func (machine *chronalMachine) step() {
	if division, ok := machine.divisions[machine.ip]; ok && division.apply(machine) {
		return
	}

	machine.registers[machine.ipRegister] = machine.ip
	inst := machine.instructions[machine.ip]
	registers := &machine.registers
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

	machine.ip = machine.registers[machine.ipRegister] + 1
}

type chronalDivision struct{ quotient, scratch, value, divisor int }

func findChronalDivisions(program [][4]int, ip int) map[int]chronalDivision {
	var result map[int]chronalDivision
	for start := 0; start+9 <= len(program); start++ {
		code := program[start : start+9]
		q, t, v := code[0][3], code[1][3], code[3][2]
		roles := map[int]bool{q: true, t: true, v: true, ip: true}
		if len(roles) != 4 {
			continue
		}
		valid := true
		for role := range roles {
			valid = valid && role >= 0 && role < 6
		}
		if !valid {
			continue
		}
		if code[0][0] != opCode("seti") || code[0][1] != 0 {
			continue
		}
		if code[1] != [4]int{opCode("addi"), q, 1, t} {
			continue
		}
		divisor := code[2][2]
		if divisor <= 0 || code[2] != [4]int{opCode("muli"), t, divisor, t} {
			continue
		}
		if code[3] != [4]int{opCode("gtrr"), t, v, t} {
			continue
		}
		if code[4] != [4]int{opCode("addr"), t, ip, ip} && code[4] != [4]int{opCode("addr"), ip, t, ip} {
			continue
		}
		if code[5] != [4]int{opCode("addi"), ip, 1, ip} {
			continue
		}
		if code[6][0] != opCode("seti") || code[6][1] != start+8 || code[6][3] != ip {
			continue
		}
		if code[7] != [4]int{opCode("addi"), q, 1, q} {
			continue
		}
		if code[8][0] != opCode("seti") || code[8][1] != start || code[8][3] != ip {
			continue
		}
		if result == nil {
			result = make(map[int]chronalDivision)
		}
		result[start] = chronalDivision{q, t, v, divisor}
	}
	return result
}

func (division chronalDivision) apply(machine *chronalMachine) bool {
	value := machine.registers[division.value]
	if value < 0 {
		return false
	}
	quotient := value / division.divisor
	maxInt := int(^uint(0) >> 1)
	// The ordinary loop tests (quotient+1)*divisor; it must not overflow.
	if quotient >= maxInt/division.divisor {
		return false
	}
	machine.registers[division.quotient] = quotient
	machine.registers[division.scratch] = 1
	machine.registers[machine.ipRegister] = machine.ip + 8
	machine.ip += 9
	return true
}

func (d Day21) Part1(input string) (string, error) {
	machine, err := newChronalMachine(input)
	if err != nil {
		return "", err
	}
	for machine.running() {
		machine.registers[machine.ipRegister] = machine.ip
		if machine.ip == machine.probe.ip {
			return strconv.Itoa(machine.registers[machine.probe.candidate]), nil
		}
		machine.step()
	}
	return "", errors.New("no halting comparison reached")
}

func (d Day21) Part2(input string) (string, error) {
	machine, err := newChronalMachine(input)
	if err != nil {
		return "", err
	}
	seenStates := make(map[[6]int]bool)
	seenCandidates := make(map[int]bool)
	var last int
	for machine.running() {
		machine.registers[machine.ipRegister] = machine.ip
		if machine.ip == machine.probe.ip {
			if seenStates[machine.registers] {
				return strconv.Itoa(last), nil
			}
			seenStates[machine.registers] = true
			candidate := machine.registers[machine.probe.candidate]
			if !seenCandidates[candidate] {
				seenCandidates[candidate] = true
				last = candidate
			}

			// Explore the continuing branch even when the candidate equals zero.
			machine.registers[machine.probe.result] = 0
			machine.ip = machine.registers[machine.ipRegister] + 1
			continue
		}
		machine.step()
	}
	return "", errors.New("no repeated comparison state reached")
}
func init() { solve.Register(Day21{}) }
