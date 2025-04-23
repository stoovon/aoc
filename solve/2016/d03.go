package solve2016

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day3 struct {
}

func (d Day3) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 3}
}

// Checks if the sides form a valid triangle
func (d Day3) isTriangular(sides []int) bool {
	if len(sides) != 3 {
		return false
	}
	x, y, z := sides[0], sides[1], sides[2]
	return x+y > z && x+z > y && y+z > x
}

// Parses the input into a slice of triangles
func (d Day3) parseInput(input string) [][]int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var triangles [][]int
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		sides := make([]int, 3)
		for i, field := range fields {
			sides[i], _ = strconv.Atoi(field)
		}
		triangles = append(triangles, sides)
	}
	return triangles
}

// Transposes the triangles for column-wise processing
func (d Day3) transposeTriangles(triangles [][]int) [][]int {
	var transposed [][]int
	for i := 0; i < len(triangles); i += 3 {
		if i+2 >= len(triangles) {
			break
		}
		for j := 0; j < 3; j++ {
			transposed = append(transposed, []int{
				triangles[i][j],
				triangles[i+1][j],
				triangles[i+2][j],
			})
		}
	}
	return transposed
}

func (d Day3) Part1(input string) (string, error) {
	triangles := d.parseInput(input)
	count := 0
	for _, triangle := range triangles {
		if d.isTriangular(triangle) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func (d Day3) Part2(input string) (string, error) {
	triangles := d.parseInput(input)
	transposed := d.transposeTriangles(triangles)
	count := 0
	for _, triangle := range transposed {
		if d.isTriangular(triangle) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func init() {
	solve.Register(Day3{})
}
