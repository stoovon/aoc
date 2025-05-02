package solve2016

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day23 struct {
}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 23}
}

func (d Day23) isNum(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (d Day23) getVal(regs map[string]int, x string) int {
	if d.isNum(x) {
		val, _ := strconv.Atoi(x)
		return val
	}
	return regs[x]
}

func (d Day23) toggle(line string) string {
	parts := strings.Fields(line)
	switch parts[0] {
	case "inc":
		parts[0] = "dec"
	case "dec", "tgl":
		parts[0] = "inc"
	case "jnz":
		parts[0] = "cpy"
	case "cpy":
		parts[0] = "jnz"
	}
	return strings.Join(parts, " ")
}

func (d Day23) run(lines []string, part2 bool) int {
	pc := 0
	regs := map[string]int{"a": 7, "b": 0, "c": 0, "d": 0}
	if part2 {
		regs["a"] = 12
	}

	for pc < len(lines) {
		// Peephole optimization: Detect multiplication loop
		if pc+4 < len(lines) &&
			strings.HasPrefix(lines[pc], "cpy") &&
			strings.HasPrefix(lines[pc+1], "inc") &&
			strings.HasPrefix(lines[pc+2], "dec") &&
			strings.HasPrefix(lines[pc+3], "jnz") &&
			strings.HasPrefix(lines[pc+4], "dec") &&
			strings.HasPrefix(lines[pc+5], "jnz") {

			// Parse the registers/values involved
			cpyParts := strings.Fields(lines[pc])
			incParts := strings.Fields(lines[pc+1])
			dec1Parts := strings.Fields(lines[pc+2])
			jnz1Parts := strings.Fields(lines[pc+3])
			dec2Parts := strings.Fields(lines[pc+4])
			jnz2Parts := strings.Fields(lines[pc+5])

			// Ensure the pattern matches the multiplication loop
			if cpyParts[0] == "cpy" && incParts[0] == "inc" && dec1Parts[0] == "dec" &&
				jnz1Parts[0] == "jnz" && dec2Parts[0] == "dec" && jnz2Parts[0] == "jnz" &&
				jnz1Parts[2] == "-2" && jnz2Parts[2] == "-5" {

				// Perform the multiplication
				src := d.getVal(regs, cpyParts[1])
				target := incParts[1]
				loop1 := dec1Parts[1]
				loop2 := dec2Parts[1]

				regs[target] += src * regs[loop2]
				regs[loop1] = 0
				regs[loop2] = 0

				// Skip the optimized instructions
				pc += 6
				continue
			}
		}

		// Regular instruction execution
		line := lines[pc]
		parts := strings.Fields(line)

		switch parts[0] {
		case "cpy":
			val := d.getVal(regs, parts[1])
			if _, ok := regs[parts[2]]; ok {
				regs[parts[2]] = val
			}
			pc++
		case "inc":
			if _, ok := regs[parts[1]]; ok {
				regs[parts[1]]++
			}
			pc++
		case "dec":
			if _, ok := regs[parts[1]]; ok {
				regs[parts[1]]--
			}
			pc++
		case "jnz":
			val := d.getVal(regs, parts[1])
			offset := d.getVal(regs, parts[2])
			if val != 0 {
				pc += offset
			} else {
				pc++
			}
		case "tgl":
			offset := d.getVal(regs, parts[1])
			target := pc + offset
			if target >= 0 && target < len(lines) {
				lines[target] = d.toggle(lines[target])
			}
			pc++
		default:
			panic("Unknown instruction: " + parts[0])
		}
	}

	return regs["a"]
}

func (d Day23) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	return strconv.Itoa(d.run(lines, false)), nil
}

func (d Day23) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	return strconv.Itoa(d.run(lines, true)), nil
}

func init() {
	solve.Register(Day23{})
}
