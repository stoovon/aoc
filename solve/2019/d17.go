package solve2019

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day17 struct {
}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 17}
}

func (d Day17) Part1(input string) (string, error) {
	c := parseIntcode(input)
	output := c.Run()
	var sb strings.Builder
	for _, v := range output {
		sb.WriteByte(byte(v))
	}

	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	height := len(lines)
	width := len(lines[0])

	sum := 0
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			if lines[y][x] == '#' &&
				lines[y-1][x] == '#' &&
				lines[y+1][x] == '#' &&
				lines[y][x-1] == '#' &&
				lines[y][x+1] == '#' {
				sum += x * y
			}
		}
	}
	return strconv.Itoa(sum), nil
}

func (d Day17) Part2(input string) (string, error) {
	c := parseIntcode(input)
	c.mem[0] = 2 // Wake up robot

	main := "A,B,A,B,C,A,C,A,C,B\n"
	funcA := "R,12,L,8,L,4,L,4\n"     // ABCC
	funcB := "L,8,R,6,L,6\n"          // BDE
	funcC := "L,8,L,4,R,12,L,6,L,4\n" // BCAEC
	video := "n\n"

	asciiInput := func(s string) []int64 {
		out := make([]int64, len(s))
		for i, ch := range s {
			out[i] = int64(ch)
		}
		return out
	}
	inputSeq := append(asciiInput(main),
		append(asciiInput(funcA),
			append(asciiInput(funcB),
				append(asciiInput(funcC), asciiInput(video)...)...)...)...)

	inputIdx := 0
	outputs := make([]int64, 0, 1000)
	for {
		res := c.Step(nil)
		if res.Halted {
			break
		}
		if res.NeedInput {
			if inputIdx >= len(inputSeq) {
				return "", fmt.Errorf("no input available")
			}
			val := inputSeq[inputIdx]
			inputIdx++
			res = c.Step(&val)
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
	if len(outputs) == 0 {
		return "", fmt.Errorf("no output produced")
	}

	return strconv.FormatInt(outputs[len(outputs)-1], 10), nil
}

func init() {
	solve.Register(Day17{})
}
