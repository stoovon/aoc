package solve2022

import (
	"aoc/solve"
	"strconv"
	"strings"
)

type Day8 struct{}

func (d Day8) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 8}
}

func (d Day8) parseInput(input string) ([][]int, int, int) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rows := len(lines)
	cols := len(lines[0])
	grid := make([][]int, rows)
	for i := range lines {
		grid[i] = make([]int, cols)
		for j, char := range lines[i] {
			grid[i][j] = int(char - '0')
		}
	}
	return grid, rows, cols
}

func (d Day8) Part1(input string) (string, error) {
	grid, rows, cols := d.parseInput(input)

	// Count visible trees
	visible := 0

	// Check edge trees
	visible += 2*rows + 2*cols - 4 // All edge trees are visible

	// Check interior trees
	for i := 1; i < rows-1; i++ {
		for j := 1; j < cols-1; j++ {
			if isVisible(grid, i, j, rows, cols) {
				visible++
			}
		}
	}

	return strconv.Itoa(visible), nil
}

func isVisible(grid [][]int, i, j, rows, cols int) bool {
	height := grid[i][j]

	// Check up
	for x := i - 1; x >= 0; x-- {
		if grid[x][j] >= height {
			break
		}
		if x == 0 {
			return true
		}
	}

	// Check down
	for x := i + 1; x < rows; x++ {
		if grid[x][j] >= height {
			break
		}
		if x == rows-1 {
			return true
		}
	}

	// Check left
	for y := j - 1; y >= 0; y-- {
		if grid[i][y] >= height {
			break
		}
		if y == 0 {
			return true
		}
	}

	// Check right
	for y := j + 1; y < cols; y++ {
		if grid[i][y] >= height {
			break
		}
		if y == cols-1 {
			return true
		}
	}

	return false
}

func (d Day8) Part2(input string) (string, error) {
	grid, rows, cols := d.parseInput(input)

	// Calculate the maximum scenic score
	maxScenicScore := 0

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			scenicScore := calculateScenicScore(grid, i, j, rows, cols)
			if scenicScore > maxScenicScore {
				maxScenicScore = scenicScore
			}
		}
	}

	return strconv.Itoa(maxScenicScore), nil
}

func calculateScenicScore(grid [][]int, i, j, rows, cols int) int {
	height := grid[i][j]

	// Calculate viewing distances
	up := 0
	for x := i - 1; x >= 0; x-- {
		up++
		if grid[x][j] >= height {
			break
		}
	}
	down := 0
	for x := i + 1; x < rows; x++ {
		down++
		if grid[x][j] >= height {
			break
		}
	}

	left := 0
	for y := j - 1; y >= 0; y-- {
		left++
		if grid[i][y] >= height {
			break
		}
	}

	right := 0
	for y := j + 1; y < cols; y++ {
		right++
		if grid[i][y] >= height {
			break
		}
	}

	// Calculate scenic score
	return up * down * left * right
}

func init() {
	solve.Register(Day8{})
}
