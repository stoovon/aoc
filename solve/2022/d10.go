package solve2022

import (
	"aoc/solve"
	"aoc/utils/grids"
	"fmt"
	"strconv"
	"strings"
)

type Day10 struct{}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 10}
}

func (d Day10) Part1(input string) (string, error) {
	instructions := parseInstructions(input)
	cycles := []int{20, 60, 100, 140, 180, 220}
	signalStrengths := calculateSignalStrengths(instructions, cycles)
	sum := 0
	for _, strength := range signalStrengths {
		sum += strength
	}
	return fmt.Sprintf("%d", sum), nil
}

func (d Day10) Part2(input string) (string, error) {
	instructions := parseInstructions(input)
	crt := simulateCRT(instructions)
	result := grids.OCR(crt)
	return result, nil
}

func init() {
	solve.Register(Day10{})
}

func parseInstructions(input string) []string {
	return strings.Split(strings.TrimSpace(input), "\n")
}

func calculateSignalStrengths(instructions []string, cycles []int) []int {
	signalStrengths := []int{}
	cycle := 0
	x := 1
	cycleIndex := 0

	for _, instruction := range instructions {
		if cycleIndex >= len(cycles) {
			break
		}

		if instruction == "noop" {
			cycle++
			if cycle == cycles[cycleIndex] {
				signalStrengths = append(signalStrengths, cycle*x)
				cycleIndex++
				if cycleIndex >= len(cycles) {
					break
				}
			}
		} else if strings.HasPrefix(instruction, "addx") {
			value := parseAddxValue(instruction)
			cycle++
			if cycle == cycles[cycleIndex] {
				signalStrengths = append(signalStrengths, cycle*x)
				cycleIndex++
				if cycleIndex >= len(cycles) {
					break
				}
			}
			cycle++
			if cycle == cycles[cycleIndex] {
				signalStrengths = append(signalStrengths, cycle*x)
				cycleIndex++
				if cycleIndex >= len(cycles) {
					break
				}
			}
			x += value
		}
	}

	return signalStrengths
}

func parseAddxValue(instruction string) int {
	parts := strings.Split(instruction, " ")
	value, _ := strconv.Atoi(parts[1])
	return value
}

func simulateCRT(instructions []string) [][]int {
	const width, height = 40, 6
	crt := make([][]int, height)
	for i := range crt {
		crt[i] = make([]int, width)
	}

	cycle := 0
	x := 1
	row := 0

	for _, instruction := range instructions {
		if row >= height {
			break
		}

		if instruction == "noop" {
			drawPixel(crt, row, cycle%width, x)
			cycle++
			if cycle%width == 0 {
				row++
			}
		} else if strings.HasPrefix(instruction, "addx") {
			value := parseAddxValue(instruction)
			drawPixel(crt, row, cycle%width, x)
			cycle++
			if cycle%width == 0 {
				row++
			}
			drawPixel(crt, row, cycle%width, x)
			cycle++
			if cycle%width == 0 {
				row++
			}
			x += value
		}
	}

	return crt
}

func drawPixel(crt [][]int, row, col, x int) {
	if row >= len(crt) || col >= len(crt[row]) {
		return
	}
	if col >= x-1 && col <= x+1 {
		crt[row][col] = 1
	}
}
