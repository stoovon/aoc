package solve2016

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day25 struct {
}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 25}
}

func (d Day25) parseLine(line string) []interface{} {
	parts := strings.Fields(line)
	parsed := make([]interface{}, len(parts))
	for i, part := range parts {
		if val, err := strconv.Atoi(part); err == nil {
			parsed[i] = val
		} else {
			parsed[i] = part
		}
	}
	return parsed
}

func (d Day25) toggle(code [][]interface{}, i int) {
	if i < 0 || i >= len(code) {
		return
	}
	inst := code[i]
	switch inst[0] {
	case "inc":
		inst[0] = "dec"
	case "dec", "tgl":
		inst[0] = "inc"
	case "jnz":
		inst[0] = "cpy"
	case "cpy":
		inst[0] = "jnz"
	}
}

func (d Day25) interpret(code [][]interface{}, regs map[string]int, steps int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		val := func(x interface{}) int {
			switch v := x.(type) {
			case int:
				return v
			case string:
				return regs[v]
			default:
				return 0
			}
		}

		pc := 0
		for step := 0; step < steps; step++ {
			if pc < 0 || pc >= len(code) {
				return
			}
			inst := code[pc]
			op := inst[0].(string)
			x := inst[1]
			var y interface{}
			if len(inst) > 2 {
				y = inst[2]
			}
			pc++
			switch op {
			case "cpy":
				if reg, ok := y.(string); ok {
					regs[reg] = val(x)
				}
			case "inc":
				regs[x.(string)]++
			case "dec":
				regs[x.(string)]--
			case "jnz":
				if val(x) != 0 {
					pc += val(y) - 1
				}
			case "tgl":
				d.toggle(code, pc-1+val(x))
			case "out":
				out <- val(x)
			}
		}
	}()
	return out
}

func (d Day25) repeats(a int, code [][]interface{}, steps, minSignals int) bool {
	regs := map[string]int{"a": a, "b": 0, "c": 0, "d": 0}
	signals := d.interpret(code, regs, steps)
	expected := 0
	count := 0
	for signal := range signals {
		if signal != expected {
			return false
		}
		expected = 1 - expected
		count++
		if count >= minSignals {
			return true
		}
	}
	return false
}

func (d Day25) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	code := make([][]interface{}, len(lines))
	for i, line := range lines {
		code[i] = d.parseLine(line)
	}

	const steps = 1000000
	const minSignals = 100
	for a := 1; ; a++ {
		if d.repeats(a, code, steps, minSignals) {
			return strconv.Itoa(a), nil
		}
	}
}

func (d Day25) Part2(input string) (string, error) {
	return "", errors.New("Not implemented")
}

func init() {
	solve.Register(Day25{})
}
