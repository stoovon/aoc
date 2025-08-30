package solve2021

import (
	"aoc/solve"
	"fmt"
	"strconv"
	"strings"
	"slices"
)

type Day24 struct{}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 24}
}

func fillSlice(value, size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = value
	}
	return slice
}

func (d Day24) solvePart(input string, fillValue int) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	cmds := [][]string{}
	for _, line := range lines {
		cmds = append(cmds, strings.Fields(line))
	}

	result := Solve(fillSlice(fillValue, 14), cmds)
	return result, nil
}

func (d Day24) Part1(input string) (string, error) {
	return d.solvePart(input, 9)
}

func (d Day24) Part2(input string) (string, error) {
	return d.solvePart(input, 1)
}

func init() {
	solve.Register(Day24{})
}

type ALU struct {
	cmds [][]string
	reg  map[string]int
}

func NewALU(cmds [][]string) *ALU {
	return &ALU{
		cmds: cmds,
		reg:  map[string]int{"w": 0, "x": 0, "y": 0, "z": 0},
	}
}

func (alu *ALU) resolve(val string) int {
	if v, ok := alu.reg[val]; ok {
		return v
	}
	res, _ := strconv.Atoi(val)
	return res
}

func (alu *ALU) Run(input []int) map[string]int {
	inputs := make(chan int, len(input))
	for _, v := range input {
		inputs <- v
	}
	close(inputs)

	for _, cmd := range alu.cmds {
		switch cmd[0] {
		case "inp":
			alu.reg[cmd[1]] = <-inputs
		case "add":
			alu.reg[cmd[1]] += alu.resolve(cmd[2])
		case "mul":
			alu.reg[cmd[1]] *= alu.resolve(cmd[2])
		case "div":
			alu.reg[cmd[1]] /= alu.resolve(cmd[2])
		case "mod":
			alu.reg[cmd[1]] %= alu.resolve(cmd[2])
		case "eql":
			if alu.resolve(cmd[1]) == alu.resolve(cmd[2]) {
				alu.reg[cmd[1]] = 1
			} else {
				alu.reg[cmd[1]] = 0
			}
		default:
			panic(fmt.Sprintf("Unknown operator: %s", cmd[0]))
		}
	}

	state := make(map[string]int)
	for k, v := range alu.reg {
		state[k] = v
	}
	alu.reg = map[string]int{"w": 0, "x": 0, "y": 0, "z": 0}
	return state
}

func CorrectInput(input []int, cmds [][]string) []int {
	inp := slices.Clone(input)
	subLen := 18
	lineNos := []int{4, 5, 15}
	stack := []struct {
		idx int
		add int
	}{}

	for i := 0; i < 14; i++ {
		divChkAdd := []string{
			cmds[i*subLen+lineNos[0]][2],
			cmds[i*subLen+lineNos[1]][2],
			cmds[i*subLen+lineNos[2]][2],
		}
		div, _ := strconv.Atoi(divChkAdd[0])
		chk, _ := strconv.Atoi(divChkAdd[1])
		add, _ := strconv.Atoi(divChkAdd[2])

		if div == 1 {
			stack = append(stack, struct {
				idx int
				add int
			}{i, add})
		} else if div == 26 {
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			inp[i] = inp[last.idx] + last.add + chk
			if inp[i] > 9 {
				inp[last.idx] -= inp[i] - 9
				inp[i] = 9
			}
			if inp[i] < 1 {
				inp[last.idx] += 1 - inp[i]
				inp[i] = 1
			}
		}
	}

	return inp
}

func Check(input []int, cmds [][]string) bool {
	alu := NewALU(cmds)
	reg := alu.Run(input)
	return reg["z"] == 0
}

func Solve(wishInput []int, cmds [][]string) string {
	input := CorrectInput(wishInput, cmds)
	solution := ""
	for _, v := range input {
		solution += strconv.Itoa(v)
	}

	if !Check(input, cmds) {
		panic(fmt.Sprintf("Solution doesn't pass the test: %s", solution))
	}

	return solution
}
