package solve2019

import (
	"fmt"

	"aoc/solve"
)

type Day13 struct {
}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 13}
}

func (d Day13) Part1(input string) (string, error) {
	code := parseIntcode(input)
	outputs := code.Run() // Run with no input, collect all outputs

	blockCount := 0
	for i := 0; i+2 < len(outputs); i += 3 {
		tileID := outputs[i+2]
		if tileID == 2 {
			blockCount++
		}
	}
	return fmt.Sprint(blockCount), nil
}

func (d Day13) Part2(input string) (string, error) {
	code := parseIntcode(input)
	code.mem[0] = 2 // Free play

	var (
		ballX, paddleX int64
		score          int64
	)

	var inputVal int64
	for !code.halted {
		outputs := make([]int64, 0, 3)
		for len(outputs) < 3 && !code.halted {
			// Provide input only if needed
			var in *int64
			if len(code.input) > 0 {
				inputVal = code.input[0]
				in = &inputVal
				code.input = nil
			}
			res := code.Step(in)
			if res.Halted {
				break
			}
			if res.NeedInput {
				// Set joystick input for next step
				if ballX < paddleX {
					inputVal = -1
				} else if ballX > paddleX {
					inputVal = 1
				} else {
					inputVal = 0
				}
				in = &inputVal
				res = code.Step(in)
				if res.Halted {
					break
				}
			}
			if res.Output != nil {
				outputs = append(outputs, *res.Output)
			}
		}
		if len(outputs) < 3 {
			break
		}
		x, y, v := outputs[0], outputs[1], outputs[2]
		if x == -1 && y == 0 {
			score = v
		} else {
			if v == 3 {
				paddleX = x
			} else if v == 4 {
				ballX = x
			}
		}
	}
	return fmt.Sprint(score), nil
}

func init() {
	solve.Register(Day13{})
}
