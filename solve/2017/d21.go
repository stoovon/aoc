package solve2017

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day21 struct {
}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 21}
}

func parseRules(input string) map[string][]string {
	rules := map[string][]string{}
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Split(line, " => ")
		in, out := parts[0], parts[1]
		inGrid := strings.Split(in, "/")
		for _, variant := range allVariants(inGrid) {
			rules[gridToString(variant)] = strings.Split(out, "/")
		}
	}
	return rules
}

func allVariants(grid []string) [][]string {
	var variants [][]string
	g := grid
	for i := 0; i < 4; i++ {
		g = rotate(g)
		variants = append(variants, g)
		variants = append(variants, flip(g))
	}
	return variants
}

func rotate(grid []string) []string {
	n := len(grid)
	res := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			row[j] = grid[n-j-1][i]
		}
		res[i] = string(row)
	}
	return res
}

func flip(grid []string) []string {
	n := len(grid)
	res := make([]string, n)
	for i := 0; i < n; i++ {
		row := []byte(grid[i])
		for j := 0; j < n/2; j++ {
			row[j], row[n-j-1] = row[n-j-1], row[j]
		}
		res[i] = string(row)
	}
	return res
}

func gridToString(grid []string) string {
	return strings.Join(grid, "/")
}

func enhance(grid []string, rules map[string][]string) []string {
	n := len(grid)
	var size, outSize int
	if n%2 == 0 {
		size, outSize = 2, 3
	} else {
		size, outSize = 3, 4
	}
	blocks := n / size
	newGrid := make([]string, blocks*outSize)
	for i := range newGrid {
		newGrid[i] = strings.Repeat(" ", blocks*outSize)
	}
	for by := 0; by < blocks; by++ {
		for bx := 0; bx < blocks; bx++ {
			block := make([]string, size)
			for y := 0; y < size; y++ {
				block[y] = grid[by*size+y][bx*size : bx*size+size]
			}
			out := rules[gridToString(block)]
			for y := 0; y < outSize; y++ {
				row := []byte(newGrid[by*outSize+y])
				copy(row[bx*outSize:bx*outSize+outSize], out[y])
				newGrid[by*outSize+y] = string(row)
			}
		}
	}
	return newGrid
}

func (d Day21) solve(input string, iterations int) (string, error) {
	rules := parseRules(input)
	grid := []string{".#.", "..#", "###"}
	for i := 0; i < iterations; i++ {
		grid = enhance(grid, rules)
	}
	count := 0
	for _, row := range grid {
		count += strings.Count(row, "#")
	}
	return strconv.Itoa(count), nil
}

func (d Day21) Part1(input string) (string, error) {
	return d.solve(input, 5)
}

func (d Day21) Part2(input string) (string, error) {
	return d.solve(input, 18)
}

func init() {
	solve.Register(Day21{})
}
