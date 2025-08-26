package solve2017

import (
	"errors"
	"fmt"

	"aoc/solve"
)

type Day15 struct {
}

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 15}
}

const (
	factorA = 16807
	factorB = 48271
	mod     = 2147483647
)

func (d Day15) parseInput(input string) (int, int, error) {
	var startA, startB int
	_, err := fmt.Sscanf(input, "Generator A starts with %d\nGenerator B starts with %d", &startA, &startB)
	if err != nil {
		return 0, 0, errors.New("invalid input")
	}
	return startA, startB, nil
}

func nextGen(start, factor, multiple int) func() int {
	val := start
	return func() int {
		for {
			val = (val * factor) % mod
			if val%multiple == 0 {
				return val
			}
		}
	}
}

func judgeCount(startA, startB, multipleA, multipleB, pairs int) int {
	genA := nextGen(startA, factorA, multipleA)
	genB := nextGen(startB, factorB, multipleB)
	matches := 0
	for i := 0; i < pairs; i++ {
		a := genA()
		b := genB()
		if (a & 0xFFFF) == (b & 0xFFFF) {
			matches++
		}
	}
	return matches
}

func (d Day15) Part1(input string) (string, error) {
	startA, startB, err := d.parseInput(input)
	if err != nil {
		return "", err
	}
	matches := judgeCount(startA, startB, 1, 1, 40000000)
	return fmt.Sprintf("%d", matches), nil
}

func (d Day15) Part2(input string) (string, error) {
	startA, startB, err := d.parseInput(input)
	if err != nil {
		return "", err
	}
	matches := judgeCount(startA, startB, 4, 8, 5000000)
	return fmt.Sprintf("%d", matches), nil
}

func init() {
	solve.Register(Day15{})
}
