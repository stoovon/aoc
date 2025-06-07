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
	prog, err := parseIntCode(input)
	if err != nil {
		return "", err
	}
	phases := []int{0, 1, 2, 3, 4}
	maxSignal := 0
	permute(phases, func(p []int) {
		signal := 0
		for i := 0; i < 5; i++ {
			ic := NewIntCode(prog, []int{p[i], signal})
			outputs := ic.Run()
			if len(outputs) == 0 {
				return
			}
			signal = outputs[len(outputs)-1]
		}
		if signal > maxSignal {
			maxSignal = signal
		}
	})
	return strconv.Itoa(maxSignal), nil
}

func (d Day7) Part2(input string) (string, error) {
	prog, err := parseIntCode(input)
	if err != nil {
		return "", err
	}
	phases := []int{5, 6, 7, 8, 9}
	maxSignal := 0

	permute(phases, func(p []int) {
		// Initialize amplifiers
		amps := make([]*IntCode, 5)
		for i := 0; i < 5; i++ {
			amps[i] = NewIntCode(prog, []int{p[i]})
		}
		amps[0].inputs = append(amps[0].inputs, 0) // initial input signal

		lastOutput := 0
		for {
			allHalted := true
			for i := 0; i < 5; i++ {
				out, produced := amps[i].Step()
				if amps[i].halted {
					continue
				}
				allHalted = false
				if produced {
					// Pass output to next amp's input
					amps[(i+1)%5].inputs = append(amps[(i+1)%5].inputs, out)
					if i == 4 {
						lastOutput = out
					}
				}
			}
			if allHalted {
				break
			}
		}
		if lastOutput > maxSignal {
			maxSignal = lastOutput
		}
	})
	return strconv.Itoa(maxSignal), nil
}

func init() {
	solve.Register(Day7{})
}
