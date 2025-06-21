package solve2017

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day10 struct {
}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 10}
}

func knotHashRound(list []int, lengths []int, rounds int) []int {
	size := len(list)
	pos, skip := 0, 0
	for r := 0; r < rounds; r++ {
		for _, length := range lengths {
			for i := 0; i < length/2; i++ {
				a := (pos + i) % size
				b := (pos + length - 1 - i) % size
				list[a], list[b] = list[b], list[a]
			}
			pos = (pos + length + skip) % size
			skip++
		}
	}
	return list
}

// knotHash computes the full knot hash hex string for a given input.
func knotHash(input string) string {
	const size = 256
	list := make([]int, size)
	for i := 0; i < size; i++ {
		list[i] = i
	}
	// Convert input to ASCII codes and append the standard suffix
	suffix := []int{17, 31, 73, 47, 23}
	var lengths []int
	for _, c := range strings.TrimSpace(input) {
		lengths = append(lengths, int(c))
	}
	lengths = append(lengths, suffix...)
	sparse := knotHashRound(list, lengths, 64)
	// Dense hash
	dense := make([]int, 16)
	for i := 0; i < 16; i++ {
		x := sparse[i*16]
		for j := 1; j < 16; j++ {
			x ^= sparse[i*16+j]
		}
		dense[i] = x
	}
	// Format as hex string
	hash := ""
	for _, n := range dense {
		hash += fmt.Sprintf("%02x", n)
	}
	return hash
}

func (d Day10) Part1(input string) (string, error) {
	const size = 256
	list := make([]int, size)
	for i := 0; i < size; i++ {
		list[i] = i
	}
	var lengths []int
	for _, s := range strings.Split(strings.TrimSpace(input), ",") {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return "", errors.New("invalid input")
		}
		lengths = append(lengths, n)
	}
	list = knotHashRound(list, lengths, 1)
	return strconv.Itoa(list[0] * list[1]), nil
}

func (d Day10) Part2(input string) (string, error) {
	return knotHash(input), nil
}

func init() {
	solve.Register(Day10{})
}
