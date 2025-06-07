package solve2016

import (
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day8 struct {
}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 8}
}

var digitRegex = regexp.MustCompile(`\d+`)

// Initializes a 6x50 screen
func (d Day8) Screen() [][]int {
	screen := make([][]int, 6)
	for i := range screen {
		screen[i] = make([]int, 50)
	}
	return screen
}

// Rotates a slice by n positions
func (d Day8) rotate(items []int, n int) []int {
	n = n % len(items)
	return append(items[len(items)-n:], items[:len(items)-n]...)
}

// Interprets a command and updates the screen
func (d Day8) interpret(cmd string, screen [][]int) {
	matches := digitRegex.FindAllString(cmd, -1)
	if len(matches) != 2 {
		return
	}
	A, _ := strconv.Atoi(matches[0])
	B, _ := strconv.Atoi(matches[1])

	if strings.HasPrefix(cmd, "rect") {
		for i := 0; i < B; i++ {
			for j := 0; j < A; j++ {
				screen[i][j] = 1
			}
		}
	} else if strings.HasPrefix(cmd, "rotate row") {
		screen[A] = d.rotate(screen[A], B)
	} else if strings.HasPrefix(cmd, "rotate col") {
		column := make([]int, len(screen))
		for i := range screen {
			column[i] = screen[i][A]
		}
		column = d.rotate(column, B)
		for i := range screen {
			screen[i][A] = column[i]
		}
	}
}

// Runs all commands and returns the final screen
func (d Day8) run(commands []string, screen [][]int) [][]int {
	for _, cmd := range commands {
		d.interpret(cmd, screen)
	}
	return screen
}

func (d Day8) Part1(input string) (string, error) {
	screen := d.run(strings.Split(strings.TrimSpace(input), "\n"), d.Screen())
	count := 0
	for _, row := range screen {
		for _, pixel := range row {
			if pixel == 1 {
				count++
			}
		}
	}
	return strconv.Itoa(count), nil
}

func (d Day8) Part2(input string) (string, error) {
	screen := d.run(strings.Split(strings.TrimSpace(input), "\n"), d.Screen())
	//fmt.Println()
	//for _, row := range screen {
	//	for _, pixel := range row {
	//		if pixel == 1 {
	//			fmt.Print("@")
	//		} else {
	//			fmt.Print(" ")
	//		}
	//	}
	//	fmt.Println()
	//}
	return grids.OCR(screen), nil
}

func init() {
	solve.Register(Day8{})
}
