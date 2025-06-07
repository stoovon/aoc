package solve2019

import (
	"fmt"
	"strconv"
	"strings"
)

type IntCode struct {
	mem     []int
	ip      int
	inputs  []int
	outputs []int
	halted  bool
}

func NewIntCode(prog []int, inputs []int) *IntCode {
	mem := append([]int(nil), prog...)
	return &IntCode{mem: mem, ip: 0, inputs: inputs}
}

func (ic *IntCode) Step() (int, bool) {
	get := func(mode, val int) int {
		if mode == 0 {
			return ic.mem[val]
		}
		return val
	}
	for {
		op := ic.mem[ic.ip] % 100
		mode1 := (ic.mem[ic.ip] / 100) % 10
		mode2 := (ic.mem[ic.ip] / 1000) % 10
		switch op {
		case 1:
			a := get(mode1, ic.mem[ic.ip+1])
			b := get(mode2, ic.mem[ic.ip+2])
			ic.mem[ic.mem[ic.ip+3]] = a + b
			ic.ip += 4
		case 2:
			a := get(mode1, ic.mem[ic.ip+1])
			b := get(mode2, ic.mem[ic.ip+2])
			ic.mem[ic.mem[ic.ip+3]] = a * b
			ic.ip += 4
		case 3:
			if len(ic.inputs) == 0 {
				return 0, false // needs input
			}
			ic.mem[ic.mem[ic.ip+1]] = ic.inputs[0]
			ic.inputs = ic.inputs[1:]
			ic.ip += 2
		case 4:
			out := get(mode1, ic.mem[ic.ip+1])
			ic.outputs = append(ic.outputs, out)
			ic.ip += 2
			return out, true // produced output
		case 5:
			if get(mode1, ic.mem[ic.ip+1]) != 0 {
				ic.ip = get(mode2, ic.mem[ic.ip+2])
			} else {
				ic.ip += 3
			}
		case 6:
			if get(mode1, ic.mem[ic.ip+1]) == 0 {
				ic.ip = get(mode2, ic.mem[ic.ip+2])
			} else {
				ic.ip += 3
			}
		case 7:
			if get(mode1, ic.mem[ic.ip+1]) < get(mode2, ic.mem[ic.ip+2]) {
				ic.mem[ic.mem[ic.ip+3]] = 1
			} else {
				ic.mem[ic.mem[ic.ip+3]] = 0
			}
			ic.ip += 4
		case 8:
			if get(mode1, ic.mem[ic.ip+1]) == get(mode2, ic.mem[ic.ip+2]) {
				ic.mem[ic.mem[ic.ip+3]] = 1
			} else {
				ic.mem[ic.mem[ic.ip+3]] = 0
			}
			ic.ip += 4
		case 99:
			ic.halted = true
			return 0, false
		default:
			ic.halted = true
			return 0, false
		}
	}
}

func (ic *IntCode) Run() []int {
	for !ic.halted {
		_, produced := ic.Step()
		if !produced && !ic.halted {
			continue
		}
	}
	return ic.outputs
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
