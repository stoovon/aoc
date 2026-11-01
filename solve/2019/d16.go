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
	signal, err := parseFFTSignal(input)
	if err != nil {
		return "", err
	}
	n := len(signal)
	phases := 100

	prefix := make([]int, n+1)
	next := make([]int, n)
	for phase := 0; phase < phases; phase++ {
		for i, digit := range signal {
			prefix[i+1] = prefix[i] + digit
		}
		for i := range signal {
			width := i + 1
			sum := 0
			sign := 1
			// Nonzero runs alternate signs, with one zero run between them.
			for start := i; start < n; start += 2 * width {
				end := min(start+width, n)
				sum += sign * (prefix[end] - prefix[start])
				sign = -sign
			}
			if sum < 0 {
				sum = -sum
			}
			next[i] = sum % 10
		}
		signal, next = next, signal
	}

	var builder strings.Builder
	for i := 0; i < 8 && i < len(signal); i++ {
		builder.WriteByte(byte('0' + signal[i]))
	}
	result := builder.String()
	return result, nil
}

func (d Day16) Part2(input string) (string, error) {
	baseSignal, err := parseFFTSignal(input)
	if err != nil {
		return "", err
	}
	if len(baseSignal) < 7 {
		return "", errors.New("signal must contain at least seven digits for the offset")
	}
	if len(baseSignal) > int(^uint(0)>>1)/10000 {
		return "", errors.New("repeated signal length overflows int")
	}
	offset := 0
	for i := 0; i < 7; i++ {
		offset = offset*10 + baseSignal[i]
	}
	totalLen := len(baseSignal) * 10000
	if offset > totalLen-8 {
		return "", errors.New("message offset leaves fewer than eight digits")
	}
	if offset < totalLen/2 {
		return "", errors.New("message offset must be in the second half of the signal")
	}

	// In the second half, each output is the sum of its remaining suffix.
	sigLen := totalLen - offset
	signal := make([]byte, sigLen)
	for i := 0; i < sigLen; i++ {
		signal[i] = byte(baseSignal[(offset+i)%len(baseSignal)])
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

func parseFFTSignal(input string) ([]int, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil, errors.New("signal must contain digits")
	}
	signal := make([]int, len(input))
	for i := range input {
		if input[i] < '0' || input[i] > '9' {
			return nil, errors.New("signal must contain only decimal digits")
		}
		signal[i] = int(input[i] - '0')
	}
	return signal, nil
}

func init() {
	solve.Register(Day16{})
}
