package solve2018

import (
	"aoc/solve"
	"bufio"
	"image"
	"regexp"
	"strconv"
	"strings"
)

type Day17 struct{}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 17}
}

func (d Day17) Part1(input string) (string, error) {
	grid, minX, maxX, minY, maxY := d.parseInput(input)
	simulateWater(grid, minX, maxX, minY, maxY)
	return strconv.Itoa(countWater(grid, minY, maxY)), nil
}

func (d Day17) Part2(input string) (string, error) {
	grid, minX, maxX, minY, maxY := d.parseInput(input)
	simulateWater(grid, minX, maxX, minY, maxY)

	// Count only settled water (~) tiles
	count := 0
	for pt, val := range grid {
		if pt.Y >= minY && pt.Y <= maxY && val == '~' {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

func (d Day17) parseInput(input string) (map[image.Point]rune, int, int, int, int) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	minX, maxX := 2000, 0
	minY, maxY := 2000, 0
	grid := make(map[image.Point]rune)

	re := regexp.MustCompile(`([xy])=(\d+), [xy]=(\d+)..(\d+)`)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if matches[1] == "x" {
			x, _ := strconv.Atoi(matches[2])
			y1, _ := strconv.Atoi(matches[3])
			y2, _ := strconv.Atoi(matches[4])
			for y := y1; y <= y2; y++ {
				grid[image.Point{x, y}] = '#'
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		} else {
			y, _ := strconv.Atoi(matches[2])
			x1, _ := strconv.Atoi(matches[3])
			x2, _ := strconv.Atoi(matches[4])
			for x := x1; x <= x2; x++ {
				grid[image.Point{x, y}] = '#'
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	return grid, minX, maxX, minY, maxY
}

func simulateWater(grid map[image.Point]rune, minX, maxX, minY, maxY int) {
	settled := make(map[image.Point]bool)
	flowing := make(map[image.Point]bool)

	var fill func(pt image.Point, direction image.Point) bool
	fill = func(pt image.Point, direction image.Point) bool {
		flowing[pt] = true

		below := image.Point{pt.X, pt.Y + 1}
		if grid[below] == 0 && !flowing[below] && below.Y <= maxY {
			fill(below, image.Point{0, 1})
		}

		if grid[below] == 0 && !settled[below] {
			return false
		}

		left := image.Point{pt.X - 1, pt.Y}
		right := image.Point{pt.X + 1, pt.Y}

		leftFilled := grid[left] == '#' || (!flowing[left] && fill(left, image.Point{-1, 0}))
		rightFilled := grid[right] == '#' || (!flowing[right] && fill(right, image.Point{1, 0}))

		if direction == (image.Point{0, 1}) && leftFilled && rightFilled {
			settled[pt] = true

			l := left
			for {
				if _, ok := flowing[l]; ok {
					settled[l] = true
					l = image.Point{l.X - 1, l.Y}
				} else {
					break
				}
			}

			r := right
			for {
				if _, ok := flowing[r]; ok {
					settled[r] = true
					r = image.Point{r.X + 1, r.Y}
				} else {
					break
				}
			}
		}

		return direction == (image.Point{-1, 0}) && (leftFilled || grid[left] == '#') ||
			direction == (image.Point{1, 0}) && (rightFilled || grid[right] == '#')
	}

	fill(image.Point{500, 0}, image.Point{0, 1})

	// Update the grid with settled and flowing water
	for pt := range flowing {
		if _, ok := settled[pt]; !ok {
			grid[pt] = '|'
		}
	}
	for pt := range settled {
		grid[pt] = '~'
	}
}

func countWater(grid map[image.Point]rune, minY, maxY int) int {
	count := 0
	for pt, val := range grid {
		if pt.Y >= minY && pt.Y <= maxY && (val == '|' || val == '~') {
			count++
		}
	}

	return count
}

func init() {
	solve.Register(Day17{})
}
