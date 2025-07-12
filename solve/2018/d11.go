package solve2018

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day11 struct{}

func (d Day11) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 11}
}

func powerLevel(x, y, serial int) int {
	rackID := x + 10
	power := rackID * y
	power += serial
	power *= rackID
	power = (power / 100) % 10
	return power - 5
}

func (d Day11) Part1(input string) (string, error) {
	serial, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	const size = 300
	var grid [size + 1][size + 1]int // 1-based
	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			grid[y][x] = powerLevel(x, y, serial)
		}
	}
	maxSum := -1 << 31
	var maxX, maxY int
	for y := 1; y <= size-2; y++ {
		for x := 1; x <= size-2; x++ {
			sum := 0
			for dy := 0; dy < 3; dy++ {
				for dx := 0; dx < 3; dx++ {
					sum += grid[y+dy][x+dx]
				}
			}
			if sum > maxSum {
				maxSum = sum
				maxX, maxY = x, y
			}
		}
	}
	return strconv.Itoa(maxX) + "," + strconv.Itoa(maxY), nil
}

func (d Day11) Part2(input string) (string, error) {
	serial, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	const size = 300
	var grid [size + 1][size + 1]int // 1-based
	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			grid[y][x] = powerLevel(x, y, serial)
		}
	}
	// Build summed-area table
	var sat [size + 1][size + 1]int
	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			sat[y][x] = grid[y][x] + sat[y-1][x] + sat[y][x-1] - sat[y-1][x-1]
		}
	}
	maxSum := -1 << 31
	var maxX, maxY, maxS int
	for s := 1; s <= size; s++ {
		for y := 1; y <= size-s+1; y++ {
			for x := 1; x <= size-s+1; x++ {
				y2, x2 := y+s-1, x+s-1
				sum := sat[y2][x2] - sat[y-1][x2] - sat[y2][x-1] + sat[y-1][x-1]
				if sum > maxSum {
					maxSum = sum
					maxX, maxY, maxS = x, y, s
				}
			}
		}
	}
	return strconv.Itoa(maxX) + "," + strconv.Itoa(maxY) + "," + strconv.Itoa(maxS), nil
}

func init() {
	solve.Register(Day11{})
}
