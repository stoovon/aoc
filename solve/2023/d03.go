package solve2023

import (
	"image"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day3 struct {
}

var digitRegex = regexp.MustCompile(`\d+`)

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 3}
}

func (d Day3) solve(input string, isPart2 bool) (int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	// Build the symbols grid
	symbols := make(map[image.Point]rune)
	for y, line := range lines {
		for x, char := range line {
			if char != '.' && (char < '0' || char > '9') {
				symbols[image.Point{X: x, Y: y}] = char
			}
		}
	}

	// Check adjacency and build the gear grid
	gears := make(map[image.Point][]int)
	partNumbersSum := 0

	for y, line := range lines {
		for _, match := range digitRegex.FindAllStringIndex(line, -1) {
			start, end := match[0], match[1]
			num := line[start:end]
			number, _ := strconv.Atoi(num)

			for pos, char := range symbols {
				symbolX, symbolY := pos.X, pos.Y
				if start-1 <= symbolX && symbolX <= end && y-1 <= symbolY && symbolY <= y+1 {
					partNumbersSum += number
					if char == '*' {
						gears[pos] = append(gears[pos], number)
					}
					break
				}
			}
		}
	}

	if isPart2 {
		totalProduct := 0
		for _, partNums := range gears {
			if len(partNums) == 2 {
				totalProduct += partNums[0] * partNums[1]
			}
		}
		return totalProduct, nil
	}

	return partNumbersSum, nil
}

func (d Day3) Part1(input string) (string, error) {
	result, err := d.solve(input, false)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(result), nil
}

func (d Day3) Part2(input string) (string, error) {
	result, err := d.solve(input, true)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day3{})
}
