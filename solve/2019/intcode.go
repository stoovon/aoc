package solve2019

import (
	"strconv"
	"strings"
)

type Intcode struct {
	mem          map[int64]int64 // Use map for sparse memory
	ip           int64
	relativeBase int64
	input        []int64
	output       []int64
	halted       bool
}

func (c *Intcode) getParam(mode, offset int) int64 {
	addr := c.ip + int64(offset)
	switch mode {
	case 0: // position
		return c.mem[c.mem[addr]]
	case 1: // immediate
		return c.mem[addr]
	case 2: // relative
		return c.mem[c.mem[addr]+c.relativeBase]
	}
	panic("invalid mode")
}

func (c *Intcode) getAddr(mode, offset int) int64 {
	addr := c.ip + int64(offset)
	switch mode {
	case 0:
		return c.mem[addr]
	case 2:
		return c.mem[addr] + c.relativeBase
	}
	panic("invalid write mode")
}

func (c *Intcode) Clone() *Intcode {
	mem := make(map[int64]int64, len(c.mem))
	for k, v := range c.mem {
		mem[k] = v
	}
	return &Intcode{
		mem:          mem,
		ip:           0,
		relativeBase: 0,
		input:        nil,
		output:       nil,
		halted:       false,
	}
}

// stepOnce executes a single instruction and returns (output, needInput, halted)
func (c *Intcode) stepOnce(input *int64) (output *int64, needInput, halted bool) {
	if c.halted {
		return nil, false, true
	}
	instr := c.mem[c.ip]
	op := instr % 100
	m1 := int((instr / 100) % 10)
	m2 := int((instr / 1000) % 10)
	m3 := int((instr / 10000) % 10)
	switch op {
	case 1:
		c.mem[c.getAddr(m3, 3)] = c.getParam(m1, 1) + c.getParam(m2, 2)
		c.ip += 4
	case 2:
		c.mem[c.getAddr(m3, 3)] = c.getParam(m1, 1) * c.getParam(m2, 2)
		c.ip += 4
	case 3:
		if input == nil {
			return nil, true, false
		}
		c.mem[c.getAddr(m1, 1)] = *input
		c.ip += 2
	case 4:
		val := c.getParam(m1, 1)
		c.ip += 2
		return &val, false, false
	case 5:
		if c.getParam(m1, 1) != 0 {
			c.ip = c.getParam(m2, 2)
		} else {
			c.ip += 3
		}
	case 6:
		if c.getParam(m1, 1) == 0 {
			c.ip = c.getParam(m2, 2)
		} else {
			c.ip += 3
		}
	case 7:
		if c.getParam(m1, 1) < c.getParam(m2, 2) {
			c.mem[c.getAddr(m3, 3)] = 1
		} else {
			c.mem[c.getAddr(m3, 3)] = 0
		}
		c.ip += 4
	case 8:
		if c.getParam(m1, 1) == c.getParam(m2, 2) {
			c.mem[c.getAddr(m3, 3)] = 1
		} else {
			c.mem[c.getAddr(m3, 3)] = 0
		}
		c.ip += 4
	case 9:
		c.relativeBase += c.getParam(m1, 1)
		c.ip += 2
	case 99:
		c.halted = true
		return nil, false, true
	default:
		panic("unknown opcode")
	}
	return nil, false, false
}

func (c *Intcode) runCore(inputFn func() int64, outputFn func(int64) bool) {
	for !c.halted {
		var input *int64
		// Only provide input if the next instruction is opcode 3
		instr := c.mem[c.ip]
		op := instr % 100
		if op == 3 {
			val := inputFn()
			input = &val
		}
		output, needInput, halted := c.stepOnce(input)
		if halted {
			return
		}
		if needInput {
			// Should not happen, as we always provide input for opcode 3
			panic("runCore: input needed but not provided")
		}
		if output != nil {
			if stop := outputFn(*output); stop {
				return
			}
		}
	}
}

func (c *Intcode) Run(input ...int64) []int64 {
	c.output = nil
	idx := 0
	c.runCore(
		func() int64 {
			if idx < len(input) {
				v := input[idx]
				idx++
				return v
			}
			panic("no input available")
		},
		func(val int64) bool {
			c.output = append(c.output, val)
			return false
		},
	)
	return c.output
}

func (c *Intcode) ProvideASCIIInput(s string) {
	if len(s) == 0 || s[len(s)-1] != '\n' {
		s += "\n"
	}
	for _, ch := range s {
		c.input = append(c.input, int64(ch))
	}
}

func (c *Intcode) RunUntilInputOrHalt() (string, bool) {
	var output strings.Builder
	for !c.halted {
		var input *int64
		// Check if next instruction is input and input buffer is empty
		instr := c.mem[c.ip]
		op := instr % 100
		if op == 3 {
			if len(c.input) == 0 {
				// Needs input, return current output
				return output.String(), true
			}
			val := c.input[0]
			c.input = c.input[1:]
			input = &val
		}
		out, needInput, halted := c.stepOnce(input)
		if halted {
			return output.String(), false
		}
		if needInput {
			return output.String(), true
		}
		if out != nil {
			output.WriteByte(byte(*out))
		}
	}
	return output.String(), false
}

func (c *Intcode) RunUntilOutput() []int64 {
	c.output = nil
	c.runCore(
		func() int64 {
			if len(c.input) == 0 {
				panic("no input available")
			}
			v := c.input[0]
			c.input = c.input[1:]
			return v
		},
		func(val int64) bool {
			c.output = append(c.output, val)
			return true
		},
	)
	return c.output
}

func (c *Intcode) ChannelRunner(in <-chan int64, out chan<- int64) {
	c.runCore(
		func() int64 { return <-in },
		func(val int64) bool {
			out <- val
			return false
		},
	)
	close(out)
}

type StepResult struct {
	Output    *int64 // Output value, if any
	NeedInput bool   // True if waiting for input
	Halted    bool   // True if program halted
}

// Step executes a single instruction, optionally using input.
func (c *Intcode) Step(input *int64) StepResult {
	output, needInput, halted := c.stepOnce(input)
	return StepResult{Output: output, NeedInput: needInput, Halted: halted}
}

func parseIntcode(input string) *Intcode {
	parts := strings.Split(strings.TrimSpace(input), ",")
	mem := make(map[int64]int64)
	for i, s := range parts {
		v, _ := strconv.ParseInt(s, 10, 64)
		mem[int64(i)] = v
	}
	return &Intcode{mem: mem}
}
