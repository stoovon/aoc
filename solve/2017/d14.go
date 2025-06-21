package solve2017

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day14 struct {
}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 14}
}

func hexToBinary(hexStr string) string {
	bin := ""
	for _, c := range hexStr {
		n, _ := strconv.ParseUint(string(c), 16, 4)
		bin += fmt.Sprintf("%04b", n)
	}
	return bin
}

func (d Day14) Part1(input string) (string, error) {
	used := 0
	for i := 0; i < 128; i++ {
		rowInput := fmt.Sprintf("%s-%d", strings.TrimSpace(input), i)
		hash := knotHash(rowInput)
		bin := hexToBinary(hash)
		for _, b := range bin {
			if b == '1' {
				used++
			}
		}
	}
	return fmt.Sprintf("%d", used), nil
}

func (d Day14) Part2(input string) (string, error) {
	// Build the 128x128 grid
	grid := make([][]byte, 128)
	for i := 0; i < 128; i++ {
		rowInput := fmt.Sprintf("%s-%d", strings.TrimSpace(input), i)
		hash := knotHash(rowInput)
		bin := hexToBinary(hash)
		grid[i] = []byte(bin)
	}

	visited := make([][]bool, 128)
	for i := range visited {
		visited[i] = make([]bool, 128)
	}

	var dfs func(x, y int)
	dfs = func(x, y int) {
		if x < 0 || x >= 128 || y < 0 || y >= 128 {
			return
		}
		if grid[x][y] != '1' || visited[x][y] {
			return
		}
		visited[x][y] = true
		dfs(x+1, y)
		dfs(x-1, y)
		dfs(x, y+1)
		dfs(x, y-1)
	}

	regions := 0
	for i := 0; i < 128; i++ {
		for j := 0; j < 128; j++ {
			if grid[i][j] == '1' && !visited[i][j] {
				regions++
				dfs(i, j)
			}
		}
	}
	return fmt.Sprintf("%d", regions), nil
}

func init() {
	solve.Register(Day14{})
}
