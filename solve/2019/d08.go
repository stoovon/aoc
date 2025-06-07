package solve2019

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day8 struct {
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 8}
}

func (d Day8) Part1(input string) (string, error) {
	const width, height = 25, 6
	layerSize := width * height
	data := strings.TrimSpace(input)
	if len(data)%layerSize != 0 {
		return "", errors.New("input length is not a multiple of layer size")
	}
	minZeros := layerSize + 1
	result := 0
	for i := 0; i < len(data); i += layerSize {
		layer := data[i : i+layerSize]
		zeros, ones, twos := 0, 0, 0
		for _, c := range layer {
			switch c {
			case '0':
				zeros++
			case '1':
				ones++
			case '2':
				twos++
			}
		}
		if zeros < minZeros {
			minZeros = zeros
			result = ones * twos
		}
	}
	return strconv.Itoa(result), nil
}

func (d Day8) Part2(input string) (string, error) {
	const width, height = 25, 6
	layerSize := width * height
	data := strings.TrimSpace(input)
	if len(data)%layerSize != 0 {
		return "", errors.New("input length is not a multiple of layer size")
	}
	layers := len(data) / layerSize
	final := make([]byte, layerSize)
	for i := 0; i < layerSize; i++ {
		final[i] = '2'
		for l := 0; l < layers; l++ {
			c := data[l*layerSize+i]
			if final[i] == '2' && c != '2' {
				final[i] = c
			}
		}
	}
	// Build grid for OCR
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			if final[y*width+x] == '1' {
				grid[y][x] = 1
			} else {
				grid[y][x] = 0
			}
		}
	}
	return grids.OCR(grid), nil
}

func init() {
	solve.Register(Day8{})
}
