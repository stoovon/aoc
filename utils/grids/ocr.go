package grids

import "fmt"

// Matches a rendered letter to predefined patterns
func matchCharacter(letterRender string) string {
	renders := map[string]string{
		"A": `
 XX  
X  X 
X  X 
XXXX 
X  X 
X  X 
`,
		"B": `
XXX  
X  X 
XXX  
X  X 
X  X 
XXX  
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
		"F": `
XXXX 
X    
XXX  
X    
X    
X    
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
		"U": `
X  X
X  X
X  X
X  X
X  X
 XX 
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

// OCR decodes the screen into letters
func OCR(screen [][]int) string {
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
		result += matchCharacter(letterRender)
	}
	return result
}
