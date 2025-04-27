package solve2016

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
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

// Matches a rendered letter to predefined patterns
func (d Day8) matchCharacter(letterRender string) string {
	renders := map[string]string{
		"A": `
 XX  
X  X 
X  X 
XXXX 
X  X 
X  X 
`,
		"C": `
 XX  
X  X 
X    
X    
X  X 
 XX  
`,
		"E": `
XXXX 
X    
XXX  
X    
X    
XXXX 
`,
		"G": `
 XX  
X  X 
X    
X XX 
X  X 
 XXX 
`,
		"H": `
X  X 
X  X 
XXXX 
X  X 
X  X 
X  X 
`,
		"J": `
  XX 
   X 
   X 
   X 
X  X 
 XX  
`,
		"K": `
X  X 
X X  
XX   
X X  
X X  
X  X 
`,
		"L": `
X    
X    
X    
X    
X    
XXXX 
`,
		"O": `
 XX  
X  X 
X  X 
X  X 
X  X 
 XX  
`,
		"P": `
XXX  
X  X 
X  X 
XXX  
X    
X    
`,
		"R": `
XXX  
X  X 
X  X 
XXX  
X X  
X  X 
`,
		"Y": `
X   X
X   X
 X X 
  X  
  X  
  X  
`,
		"Z": `
XXXX 
   X 
  X  
 X   
X    
XXXX 
`,
	}

	for k, v := range renders {
		if "\n"+letterRender == v {
			return k
		}
	}

	fmt.Printf("Letter not recognised: \n\n%s", letterRender)
	return "!"
}

// Decodes the screen into letters
func (d Day8) ocr(screen [][]int) string {
	result := ""
	for col := 0; col < len(screen[0]); col += 5 {
		letterRender := ""
		for row := 0; row < len(screen); row++ {
			line := ""
			for c := col; c < col+5 && c < len(screen[0]); c++ {
				if screen[row][c] == 1 {
					line += "X"
				} else {
					line += " "
				}
			}
			letterRender += line + "\n"
		}
		result += d.matchCharacter(letterRender)
	}
	return result
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
	return d.ocr(screen), nil
}

func init() {
	solve.Register(Day8{})
}
