package solve2020

import (
	"aoc/solve"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Day8 struct{}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 8}
}

func (d Day8) Part1(input string) (string, error) {
	lines := strings.FieldsFunc(strings.ReplaceAll(input, "\r\n", "\n"), func(r rune) bool { return r == '\n' })
	acc := 0
	visited := make(map[int]bool)
	for i := 0; i < len(lines); {
		if visited[i] {
			return strconv.Itoa(acc), nil
		}
		visited[i] = true
		var op string
		var arg int
		n, _ := fmt.Sscanf(lines[i], "%s %d", &op, &arg)
		if n != 2 {
			return "", nil
		}
		switch op {
		case "acc":
			acc += arg
			i++
		case "jmp":
			i += arg
		case "nop":
			i++
		default:
			i++
		}
	}
	return strconv.Itoa(acc), nil
}

func (d Day8) Part2(input string) (string, error) {
	lines := strings.FieldsFunc(strings.ReplaceAll(input, "\r\n", "\n"), func(r rune) bool { return r == '\n' })
	for swapIdx := 0; swapIdx < len(lines); swapIdx++ {
		var op string
		var arg int
		n, _ := fmt.Sscanf(lines[swapIdx], "%s %d", &op, &arg)
		if n != 2 {
			continue
		}
		if op != "jmp" && op != "nop" {
			continue
		}
		// Copy instructions and swap one
		inst := make([]string, len(lines))
		copy(inst, lines)
		if op == "jmp" {
			inst[swapIdx] = "nop " + inst[swapIdx][4:]
		} else {
			inst[swapIdx] = "jmp " + inst[swapIdx][4:]
		}
		acc := 0
		visited := make(map[int]bool)
		i := 0
		for i < len(inst) {
			if visited[i] {
				break
			}
			visited[i] = true
			var op2 string
			var arg2 int
			n2, _ := fmt.Sscanf(inst[i], "%s %d", &op2, &arg2)
			if n2 != 2 {
				break
			}
			switch op2 {
			case "acc":
				acc += arg2
				i++
			case "jmp":
				i += arg2
			case "nop":
				i++
			default:
				i++
			}
		}
		if i == len(inst) {
			return strconv.Itoa(acc), nil
		}
	}
	return "", errors.New("no solution found")
}

func init() {
	solve.Register(Day8{})
}
