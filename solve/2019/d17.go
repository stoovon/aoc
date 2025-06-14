package solve2019

import (
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
	//fmt.Println("")
	//fmt.Println("scaffold map:")
	//fmt.Print(sb.String())

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

	// Based on following the scaffold map and manually finding the path:
	// R,12,L,8,L,4,L,4,L,8,R,6,L,6,R,12,L,8,L,4,L,4,L,8,R,6,L,6,L,8,L,4,R,12,L,6,L,4,R,12,L,8,L,4,L,4,L,8,L,4,R,12,L,6,L,4,R,12,L,8,L,4,L,4,L,8,L,4,R,12,L,6,L,4,L,8,R,6,L,6")

	// Giving each part unique names
	// A = R,12
	// B = L,8
	// C = L,4
	// D = R,6
	// E = L,6

	// A B C C B D E A B C C B D E B C A E C A B C C B C A E C A B C C B C A E C B D E
	// [A B C C] [B D E] [A B C C] [B D E] [B C A E C] [A B C C] [B C A E C] [A B C C] [B C A E C] [B D E]
	// Main function: A,B,A,B,C,A,C,A,C,B

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

	c.input = inputSeq
	output := c.Run()

	//for _, v := range output {
	//	if v > 255 {
	//		fmt.Println("Dust collected:", v)
	//	} else {
	//		fmt.Printf("%c", v)
	//		if v == 10 { // Newline
	//			time.Sleep(20 * time.Millisecond) // Adjust speed as needed
	//		}
	//	}
	//}

	return strconv.FormatInt(output[len(output)-1], 10), nil
}

func init() {
	solve.Register(Day17{})
}
