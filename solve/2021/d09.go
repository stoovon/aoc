package solve2021

import (
	"aoc/solve"
	"errors"
	"sort"
	"strconv"
	"strings"
)

type Day9 struct{}

func parseHeightmap(input string) [][]int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	heightmap := make([][]int, len(lines))
	for i, line := range lines {
		row := make([]int, len(line))
		for j, c := range line {
			row[j] = int(c - '0')
		}
		heightmap[i] = row
	}
	return heightmap
}

func findLowPoints(heightmap [][]int) [][2]int {
	rows, cols := len(heightmap), len(heightmap[0])
	var lows [][2]int
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			h := heightmap[i][j]
			low := true
			for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < rows && nj >= 0 && nj < cols {
					if heightmap[ni][nj] <= h {
						low = false
						break
					}
				}
			}
			if low {
				lows = append(lows, [2]int{i, j})
			}
		}
	}
	return lows
}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 9}
}

func (d Day9) Part1(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("empty input")
	}
	heightmap := parseHeightmap(input)
	lows := findLowPoints(heightmap)
	riskSum := 0
	for _, pt := range lows {
		i, j := pt[0], pt[1]
		riskSum += heightmap[i][j] + 1
	}
	return strconv.Itoa(riskSum), nil
}

func (d Day9) Part2(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("empty input")
	}
	heightmap := parseHeightmap(input)
	lows := findLowPoints(heightmap)
	rows, cols := len(heightmap), len(heightmap[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}
	var basinSizes []int
	for _, pt := range lows {
		stack := [][2]int{pt}
		size := 0
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			i, j := cur[0], cur[1]
			if visited[i][j] {
				continue
			}
			visited[i][j] = true
			if heightmap[i][j] == 9 {
				continue
			}
			size++
			for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < rows && nj >= 0 && nj < cols {
					if !visited[ni][nj] && heightmap[ni][nj] != 9 {
						stack = append(stack, [2]int{ni, nj})
					}
				}
			}
		}
		basinSizes = append(basinSizes, size)
	}
	if len(basinSizes) < 3 {
		return "", errors.New("not enough basins")
	}
	sort.Slice(basinSizes, func(i, j int) bool { return basinSizes[i] > basinSizes[j] })
	prod := basinSizes[0] * basinSizes[1] * basinSizes[2]
	return strconv.Itoa(prod), nil
}

func init() {
	solve.Register(Day9{})
}
