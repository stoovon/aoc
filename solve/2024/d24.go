package solve2024

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maps"
)

type Day24 struct {
}

func (d Day24) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 24}
}

type Circuit struct {
	values map[string]int
	gates  map[string]GateInfo
	cache  map[string]int
}

func newCircuit() *Circuit {
	return &Circuit{
		values: make(map[string]int),
		gates:  make(map[string]GateInfo),
		cache:  make(map[string]int),
	}
}

func (c *Circuit) evaluate(op string) int {
	if val, ok := c.cache[op]; ok {
		return val
	}
	if val, ok := c.values[op]; ok {
		return val
	}
	gate, ok := c.gates[op]
	if !ok {
		panic(fmt.Sprintf("gate not found for operation: %s", op))
	}
	if len(gate.inputs) < 2 {
		panic(fmt.Sprintf("invalid inputs for gate: %v", gate))
	}
	op1 := c.evaluate(gate.inputs[0])
	op2 := c.evaluate(gate.inputs[1])
	var result int
	switch gate.operation {
	case "AND":
		result = op1 & op2
	case "OR":
		result = op1 | op2
	case "XOR":
		result = op1 ^ op2
	default:
		panic(fmt.Sprintf("unknown operation: %s", gate.operation))
	}
	c.cache[op] = result
	return result
}

func parseInput24(input string) *Circuit {
	circuit := newCircuit()
	gates := make(map[string]GateInfo)
	sections := strings.Split(input, "\n\n")

	// Parse values
	for _, line := range strings.Split(sections[0], "\n") {
		parts := strings.Split(line, " ")
		op, val := parts[0], parts[1]
		circuit.values[op[:len(op)-1]], _ = strconv.Atoi(val)
	}

	// Parse gates
	for _, line := range strings.Split(sections[1], "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		op1, op, op2, outOp := parts[0], parts[1], parts[2], parts[4]
		gates[outOp] = GateInfo{
			operation: op,
			inputs:    []string{op1, op2},
			output:    outOp,
		}
	}

	circuit.gates = gates
	return circuit
}

func (d Day24) Part1(input string) (string, error) {
	circuit := parseInput24(input)

	digits := make(map[int]int)
	for op := range circuit.gates {
		if strings.HasPrefix(op, "z") {
			index, _ := strconv.Atoi(op[1:])
			digits[index] = circuit.evaluate(op)
		}
	}

	x := maps.MaxKey(digits)
	num := 0
	for i := x; i >= 0; i-- {
		num = 2*num + digits[i]
	}

	return strconv.Itoa(num), nil
}

func find(wire1, wire2, op string, gateStrings []string) string {
	for _, gate := range gateStrings {
		if strings.Contains(gate, wire1) && strings.Contains(gate, wire2) && strings.Contains(gate, op) {
			parts := strings.Split(gate, " ")
			return parts[len(parts)-1]
		}
	}
	return ""
}

type GateInfo struct {
	operation string   // The operation performed by the gate (e.g., "AND", "OR", "XOR").
	inputs    []string // The input wires for the gate.
	output    string   // The output wire of the gate.
}

func parseGatesToStrings(gates map[string]GateInfo) []string {
	var gateStrings []string
	for wireName, gate := range gates {
		gateStr := fmt.Sprintf("%s %s %s -> %s",
			gate.inputs[0], // First input
			gate.operation, // Operation
			gate.inputs[1], // Second input
			wireName)
		gateStrings = append(gateStrings, gateStr)
	}
	return gateStrings
}

func findHalfAdder(n string, gateStrings []string) (string, string) {
	m1 := find(fmt.Sprintf("x%s", n), fmt.Sprintf("y%s", n), "XOR", gateStrings)
	n1 := find(fmt.Sprintf("x%s", n), fmt.Sprintf("y%s", n), "AND", gateStrings)
	return m1, n1
}

func processFullAdder(carry string, m1, n1 *string, gateStrings []string, swapped *[]string) (string, string) {
	r1 := find(carry, *m1, "AND", gateStrings)
	if r1 == "" {
		*m1, *n1 = *n1, *m1
		*swapped = append(*swapped, *m1, *n1)
		r1 = find(carry, *m1, "AND", gateStrings)
	}

	z1 := find(carry, *m1, "XOR", gateStrings)
	return r1, z1
}

func checkAndSwap(wire1, wire2 string, swapped *[]string) (string, string) {
	if strings.HasPrefix(wire1, "z") {
		wire1, wire2 = wire2, wire1
		*swapped = append(*swapped, wire1, wire2)
	}
	return wire1, wire2
}

func (d Day24) Part2(input string) (string, error) {
	circuit := parseInput24(input)

	var swapped []string
	var carry string

	gateStrings := parseGatesToStrings(circuit.gates)

	// Check each bit position for adder structure
	for i := 0; i < 45; i++ {
		n := fmt.Sprintf("%02d", i)
		var m1, n1, r1, z1, c1 string

		m1, n1 = findHalfAdder(n, gateStrings)

		if carry != "" {
			r1, z1 = processFullAdder(carry, &m1, &n1, gateStrings, &swapped)

			// Check for misplaced z wires
			m1, z1 = checkAndSwap(m1, z1, &swapped)
			n1, z1 = checkAndSwap(n1, z1, &swapped)
			r1, z1 = checkAndSwap(r1, z1, &swapped)

			c1 = find(r1, n1, "OR", gateStrings)
		}

		if strings.HasPrefix(c1, "z") && c1 != "z45" {
			c1, z1 = z1, c1
			swapped = append(swapped, c1, z1)
		}

		if carry == "" {
			carry = n1
		} else {
			carry = c1
		}
	}

	sort.Strings(swapped)
	return strings.Join(swapped, ","), nil
}

func init() {
	solve.Register(Day24{})
}
