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

func (c *Intcode) RunUntilOutput() []int64 {
	c.output = nil
	for !c.halted {
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
			if len(c.input) == 0 {
				return c.output
			}
			c.mem[c.getAddr(m1, 1)] = c.input[0]
			c.input = c.input[1:]
			c.ip += 2
		case 4:
			val := c.getParam(m1, 1)
			c.output = append(c.output, val)
			c.ip += 2
			return c.output
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
			return c.output
		default:
			panic("unknown opcode")
		}
	}
	return c.output
}

func (c *Intcode) Run(input ...int64) []int64 {
	c.input = append(c.input, input...)
	for {
		instr := c.mem[c.ip]
		op := instr % 100
		m1 := int((instr / 100) % 10)
		m2 := int((instr / 1000) % 10)
		m3 := int((instr / 10000) % 10)
		switch op {
		case 1: // add
			c.mem[c.getAddr(m3, 3)] = c.getParam(m1, 1) + c.getParam(m2, 2)
			c.ip += 4
		case 2: // mul
			c.mem[c.getAddr(m3, 3)] = c.getParam(m1, 1) * c.getParam(m2, 2)
			c.ip += 4
		case 3: // input
			c.mem[c.getAddr(m1, 1)] = c.input[0]
			c.input = c.input[1:]
			c.ip += 2
		case 4: // output
			c.output = append(c.output, c.getParam(m1, 1))
			c.ip += 2
		case 5: // jump-if-true
			if c.getParam(m1, 1) != 0 {
				c.ip = c.getParam(m2, 2)
			} else {
				c.ip += 3
			}
		case 6: // jump-if-false
			if c.getParam(m1, 1) == 0 {
				c.ip = c.getParam(m2, 2)
			} else {
				c.ip += 3
			}
		case 7: // less than
			if c.getParam(m1, 1) < c.getParam(m2, 2) {
				c.mem[c.getAddr(m3, 3)] = 1
			} else {
				c.mem[c.getAddr(m3, 3)] = 0
			}
			c.ip += 4
		case 8: // equals
			if c.getParam(m1, 1) == c.getParam(m2, 2) {
				c.mem[c.getAddr(m3, 3)] = 1
			} else {
				c.mem[c.getAddr(m3, 3)] = 0
			}
			c.ip += 4
		case 9: // adjust relative base
			c.relativeBase += c.getParam(m1, 1)
			c.ip += 2
		case 99:
			c.halted = true
			return c.output
		default:
			panic("unknown opcode")
		}
	}
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
