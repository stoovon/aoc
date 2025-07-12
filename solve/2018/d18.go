package solve2018

import (
	"aoc/solve"
	"bufio"
	"strconv"
	"strings"
)

type Day18 struct{}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 18}
}

func (d Day18) Part1(input string) (string, error) {
	grid := d.parseInput(input)
	for i := 0; i < 10; i++ {
		grid = d.simulate(grid)
	}
	wooded, lumberyards := d.countAcres(grid)
	return strconv.Itoa(wooded * lumberyards), nil
}

func (d Day18) Part2(input string) (string, error) {
	grid := d.parseInput(input)
	seen := make(map[string]int)
	minute := 0

	serialize := func(grid [][]rune) string {
		var sb strings.Builder
		for _, row := range grid {
			sb.WriteString(string(row))
		}
		return sb.String()
	}

	for minute < 1000000000 {
		key := serialize(grid)
		if prev, ok := seen[key]; ok {
			cycleLen := minute - prev
			remaining := (1000000000 - prev) % cycleLen
			for i := 0; i < remaining; i++ {
				grid = d.simulate(grid)
			}
			wooded, lumberyards := d.countAcres(grid)
			return strconv.Itoa(wooded * lumberyards), nil
		}
		seen[key] = minute
		grid = d.simulate(grid)
		minute++
	}

	wooded, lumberyards := d.countAcres(grid)
	return strconv.Itoa(wooded * lumberyards), nil
}

func (d Day18) parseInput(input string) [][]rune {
	scanner := bufio.NewScanner(strings.NewReader(input))
	grid := [][]rune{}
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []rune(line))
	}
	return grid
}

func (d Day18) simulate(grid [][]rune) [][]rune {
	height := len(grid)
	width := len(grid[0])
	newGrid := make([][]rune, height)
	for i := range newGrid {
		newGrid[i] = make([]rune, width)
	}

	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			trees, lumberyards := 0, 0
			for _, d := range directions {
				dy, dx := y+d[0], x+d[1]
				if dy >= 0 && dy < height && dx >= 0 && dx < width {
					if grid[dy][dx] == '|' {
						trees++
					} else if grid[dy][dx] == '#' {
						lumberyards++
					}
				}
			}

			switch grid[y][x] {
			case '.':
				if trees >= 3 {
					newGrid[y][x] = '|'
				} else {
					newGrid[y][x] = '.'
				}
			case '|':
				if lumberyards >= 3 {
					newGrid[y][x] = '#'
				} else {
					newGrid[y][x] = '|'
				}
			case '#':
				if lumberyards >= 1 && trees >= 1 {
					newGrid[y][x] = '#'
				} else {
					newGrid[y][x] = '.'
				}
			}
		}
	}

	return newGrid
}

func (d Day18) countAcres(grid [][]rune) (int, int) {
	wooded, lumberyards := 0, 0
	for _, row := range grid {
		for _, acre := range row {
			if acre == '|' {
				wooded++
			} else if acre == '#' {
				lumberyards++
			}
		}
	}
	return wooded, lumberyards
}

func init() {
	solve.Register(Day18{})
}
