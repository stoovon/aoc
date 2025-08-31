package solve2022

import (
	"aoc/solve"
	"container/list"
	"errors"
	"strconv"
	"strings"
)

type Day12 struct{}

func (d Day12) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 12}
}

func (d Day12) Part1(input string) (string, error) {
	grid, start, end, err := d.parseHeightmap(input)
	if err != nil {
		return "", err
	}

	steps := bfs(grid, start, end)
	if steps == -1 {
		return "", errors.New("no path found")
	}

	return strconv.Itoa(steps), nil
}

func (d Day12) parseHeightmap(input string) ([][]rune, [2]int, [2]int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	grid := make([][]rune, len(lines))
	var start, end [2]int

	for i, line := range lines {
		grid[i] = []rune(line)
		for j, char := range grid[i] {
			switch char {
			case 'S':
				start = [2]int{i, j}
				grid[i][j] = 'a' // Treat 'S' as elevation 'a'
			case 'E':
				end = [2]int{i, j}
				grid[i][j] = 'z' // Treat 'E' as elevation 'z'
			}
		}
	}

	return grid, start, end, nil
}

func bfs(grid [][]rune, start, end [2]int) int {
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	rows, cols := len(grid), len(grid[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	queue := list.New()
	queue.PushBack([3]int{start[0], start[1], 0}) // {row, col, steps}
	visited[start[0]][start[1]] = true

	for queue.Len() > 0 {
		elem := queue.Remove(queue.Front()).([3]int)
		r, c, steps := elem[0], elem[1], elem[2]

		if [2]int{r, c} == end {
			return steps
		}

		for _, dir := range directions {
			nr, nc := r+dir[0], c+dir[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] {
				if grid[nr][nc] <= grid[r][c]+1 { // Elevation constraint
					visited[nr][nc] = true
					queue.PushBack([3]int{nr, nc, steps + 1})
				}
			}
		}
	}

	return -1 // No path found
}

func (d Day12) Part2(input string) (string, error) {
	grid, _, end, err := d.parseHeightmap(input)
	if err != nil {
		return "", err
	}

	rows, cols := len(grid), len(grid[0])
	minSteps := -1

	for r := range rows {
		for c := range cols {
			if grid[r][c] == 'a' {
				steps := bfs(grid, [2]int{r, c}, end)
				if steps != -1 && (minSteps == -1 || steps < minSteps) {
					minSteps = steps
				}
			}
		}
	}

	if minSteps == -1 {
		return "", errors.New("no path found")
	}

	return strconv.Itoa(minSteps), nil
}

func init() {
	solve.Register(Day12{})
}
