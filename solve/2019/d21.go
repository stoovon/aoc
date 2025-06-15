package solve2019

import (
	"errors"
	"strconv"

	"aoc/solve"
)

type Day21 struct {
}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 21}
}

func runSpringDroid(input string, script []string, command string) (string, error) {
	// Prepare Springscript input
	var inputASCII []int64
	for _, line := range script {
		for _, ch := range line {
			inputASCII = append(inputASCII, int64(ch))
		}
		inputASCII = append(inputASCII, 10)
	}
	for _, ch := range command {
		inputASCII = append(inputASCII, int64(ch))
	}
	inputASCII = append(inputASCII, 10)

	prog := parseIntcode(input)
	inputIdx := 0
	outputs := make([]int64, 0, 1000)

	for {
		res := prog.Step(nil)
		if res.Halted {
			break
		}
		if res.NeedInput {
			if inputIdx >= len(inputASCII) {
				return "", errors.New("no input available")
			}
			val := inputASCII[inputIdx]
			inputIdx++
			res = prog.Step(&val)
			if res.Halted {
				break
			}
			if res.Output != nil {
				outputs = append(outputs, *res.Output)
			}
			continue
		}
		if res.Output != nil {
			outputs = append(outputs, *res.Output)
		}
	}

	for _, v := range outputs {
		if v > 255 {
			return strconv.FormatInt(v, 10), nil
		}
		// fmt.Print(string(rune(v)))
	}
	return "", errors.New("No damage output found")
}

func (d Day21) Part1(input string) (string, error) {
	script := []string{
		"NOT A J",
		"NOT B T",
		"OR T J",
		"NOT C T",
		"OR T J",
		"AND D J",
	}
	return runSpringDroid(input, script, "WALK")
}

func (d Day21) Part2(input string) (string, error) {
	script := []string{
		"NOT A J",
		"NOT B T",
		"OR T J",
		"NOT C T",
		"OR T J",
		"AND D J",
		"NOT H T",
		"NOT T T",
		"OR E T",
		"AND T J",
	}
	return runSpringDroid(input, script, "RUN")
}

func init() {
	solve.Register(Day21{})
}
