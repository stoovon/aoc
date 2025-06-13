package solve2019

import (
	"errors"
	"strings"

	"aoc/solve"
)

type Day16 struct {
}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 16}
}

func (d Day16) Part1(input string) (string, error) {
	basePattern := []int{0, 1, 0, -1}

	signal := make([]int, 0, len(input))
	for _, ch := range input {
		if ch >= '0' && ch <= '9' {
			signal = append(signal, int(ch-'0'))
		}
	}
	n := len(signal)
	phases := 100

	for phase := 0; phase < phases; phase++ {
		next := make([]int, n)
		for i := 0; i < n; i++ {
			sum := 0
			for j := 0; j < n; j++ {
				// Pattern repeats every (i+1) times, skip first value
				patternIdx := ((j + 1) / (i + 1)) % 4
				sum += signal[j] * basePattern[patternIdx]
			}
			if sum < 0 {
				sum = -sum
			}
			next[i] = sum % 10
		}
		signal = next
	}

	var builder strings.Builder
	for i := 0; i < 8 && i < len(signal); i++ {
		builder.WriteByte(byte('0' + signal[i]))
	}
	result := builder.String()
	return result, nil
}

func (d Day16) Part2(input string) (string, error) {
	baseSignal := make([]int, 0, len(input))
	for _, ch := range input {
		if ch >= '0' && ch <= '9' {
			baseSignal = append(baseSignal, int(ch-'0'))
		}
	}
	offset := 0
	for i := 0; i < 7; i++ {
		offset = offset*10 + baseSignal[i]
	}
	totalLen := len(baseSignal) * 10000
	if offset < totalLen/2 {
		return "", errors.New("Offset too small for efficient solution")
	}

	// Only need the part from offset to end
	sigLen := totalLen - offset
	signal := make([]int, sigLen)
	for i := 0; i < sigLen; i++ {
		signal[i] = baseSignal[(offset+i)%len(baseSignal)]
	}
	for phase := 0; phase < 100; phase++ {
		for i := sigLen - 2; i >= 0; i-- {
			signal[i] = (signal[i] + signal[i+1]) % 10
		}
	}
	var builder strings.Builder
	for i := 0; i < 8; i++ {
		builder.WriteByte(byte('0' + signal[i]))
	}
	return builder.String(), nil
}

func init() {
	solve.Register(Day16{})
}
