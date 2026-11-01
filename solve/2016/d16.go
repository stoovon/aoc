package solve2016

import (
	"errors"
	"strings"

	"aoc/solve"
)

type Day16 struct {
}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 16}
}

func (d Day16) expand(seed string, n int) []byte {
	data := make([]byte, 0, n)
	data = append(data, seed[:min(len(seed), n)]...)
	for len(data) < n {
		originalLen := len(data)
		data = append(data, '0')
		for i := originalLen - 1; i >= 0 && len(data) < n; i-- {
			data = append(data, '0'+'1'-data[i])
		}
	}
	return data
}

func (d Day16) checksum(data []byte) string {
	for len(data) > 0 && len(data)%2 == 0 {
		for i := 0; i < len(data)/2; i++ {
			if data[2*i] == data[2*i+1] {
				data[i] = '1'
			} else {
				data[i] = '0'
			}
		}
		data = data[:len(data)/2]
	}
	return string(data)
}

func (d Day16) Part1(input string) (string, error) {
	return d.solve(input, 272)
}

func (d Day16) Part2(input string) (string, error) {
	return d.solve(input, 35651584)
}

func (d Day16) solve(input string, size int) (string, error) {
	seed := strings.TrimSpace(input)
	if len(seed) == 0 {
		return "", errors.New("seed must contain binary digits")
	}
	for i := range seed {
		if seed[i] != '0' && seed[i] != '1' {
			return "", errors.New("seed must contain only binary digits")
		}
	}
	return d.checksum(d.expand(seed, size)), nil
}

func init() {
	solve.Register(Day16{})
}
