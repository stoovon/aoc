package solve2019

import (
	"strconv"

	"aoc/solve"
)

type Day7 struct {
}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 7}
}

func permute(a []int, f func([]int)) {
	var generate func(int)
	generate = func(n int) {
		if n == 1 {
			tmp := make([]int, len(a))
			copy(tmp, a)
			f(tmp)
			return
		}
		for i := 0; i < n; i++ {
			generate(n - 1)
			if n%2 == 1 {
				a[0], a[n-1] = a[n-1], a[0]
			} else {
				a[i], a[n-1] = a[n-1], a[i]
			}
		}
	}
	generate(len(a))
}

func (d Day7) Part1(input string) (string, error) {
	prog := parseIntcode(input)
	phases := []int{0, 1, 2, 3, 4}
	maxSignal := int64(0)
	permute(phases, func(p []int) {
		signal := int64(0)
		for i := 0; i < 5; i++ {
			ic := prog.Clone()
			out := ic.Run(int64(p[i]), signal)
			if len(out) == 0 {
				return
			}
			signal = out[len(out)-1]
		}
		if signal > maxSignal {
			maxSignal = signal
		}
	})
	return strconv.FormatInt(maxSignal, 10), nil
}

func (d Day7) Part2(input string) (string, error) {
	prog := parseIntcode(input)
	phases := []int{5, 6, 7, 8, 9}
	maxSignal := int64(0)
	permute(phases, func(p []int) {
		amps := make([]*Intcode, 5)
		for i := 0; i < 5; i++ {
			amps[i] = prog.Clone()
			amps[i].input = append(amps[i].input, int64(p[i]))
		}
		amps[0].input = append(amps[0].input, 0)
		lastOutput := int64(0)
		ampIdx := 0
		for {
			ic := amps[ampIdx%5]
			if ic.halted {
				break
			}
			out := ic.RunUntilOutput()
			if len(out) > 0 {
				next := amps[(ampIdx+1)%5]
				next.input = append(next.input, out[0])
				if ampIdx%5 == 4 {
					lastOutput = out[0]
				}
			}
			ampIdx++
		}
		if lastOutput > maxSignal {
			maxSignal = lastOutput
		}
	})
	return strconv.FormatInt(maxSignal, 10), nil
}

func init() {
	solve.Register(Day7{})
}
