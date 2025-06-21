package solve2017

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day8 struct {
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 8}
}

func (d Day8) processInstructions(input string) (finalMax int, maxEver int, err error) {
	registers := make(map[string]int)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	maxEver = 0
	first := true

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 7 {
			return 0, 0, errors.New("invalid instruction format")
		}
		reg := parts[0]
		op := parts[1]
		val, _ := strconv.Atoi(parts[2])
		condReg := parts[4]
		condOp := parts[5]
		condVal, _ := strconv.Atoi(parts[6])

		cond := false
		switch condOp {
		case ">":
			cond = registers[condReg] > condVal
		case "<":
			cond = registers[condReg] < condVal
		case ">=":
			cond = registers[condReg] >= condVal
		case "<=":
			cond = registers[condReg] <= condVal
		case "==":
			cond = registers[condReg] == condVal
		case "!=":
			cond = registers[condReg] != condVal
		}

		if cond {
			if op == "inc" {
				registers[reg] += val
			} else if op == "dec" {
				registers[reg] -= val
			}
			if first || registers[reg] > maxEver {
				maxEver = registers[reg]
				first = false
			}
		}
	}

	finalMax = 0
	first = true
	for _, v := range registers {
		if first || v > finalMax {
			finalMax = v
			first = false
		}
	}
	return finalMax, maxEver, nil
}

func (d Day8) Part1(input string) (string, error) {
	finalMax, _, err := d.processInstructions(input)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(finalMax), nil
}

func (d Day8) Part2(input string) (string, error) {
	_, maxEver, err := d.processInstructions(input)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(maxEver), nil
}

func init() {
	solve.Register(Day8{})
}

func init() {
	solve.Register(Day8{})
}
