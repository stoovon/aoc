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

	for !code.halted {
		outs, _ := code.Step(nil, 3)
		if len(outs) < 3 {
			break
		}
		x, y, v := outs[0], outs[1], outs[2]
		if x == -1 && y == 0 {
			score = v
		} else {
			if v == 3 {
				paddleX = x
			} else if v == 4 {
				ballX = x
			}
		}
		// Set joystick input for next step (overwrite, don't append)
		var input int64
		if ballX < paddleX {
			input = -1
		} else if ballX > paddleX {
			input = 1
		} else {
			input = 0
		}
		code.input = []int64{input}
	}
	return fmt.Sprint(score), nil
}

func init() {
	solve.Register(Day13{})
}
