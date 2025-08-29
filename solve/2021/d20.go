package solve2021

import (
	"aoc/solve"
	"fmt"
	"strconv"
	"strings"
)

type Day20 struct{}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 20}
}

func (d Day20) parseInputMap(input string) ([]string, map[[2]int]string) {
	parts := strings.Split(input, "\n\n")
	algorithm := strings.Split(parts[0], "")

	image := map[[2]int]string{}
	for r, line := range strings.Split(parts[1], "\n") {
		for c, char := range strings.Split(line, "") {
			image[[2]int{r, c}] = char
		}
	}

	return algorithm, image
}

func (d Day20) getAlgIndex(image map[[2]int]string, r, c int, defaultChar string) int {
	binary := ""
	for _, d := range [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 0}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	} {
		coord := [2]int{r + d[0], c + d[1]}
		if val, exists := image[coord]; exists {
			if val == "#" {
				binary += "1"
			} else {
				binary += "0"
			}
		} else {
			if defaultChar == "#" {
				binary += "1"
			} else {
				binary += "0"
			}
		}
	}
	index, _ := strconv.ParseInt(binary, 2, 64)
	return int(index)
}

func (d Day20) enhanceImageMap(image map[[2]int]string, algorithm []string, infiniteChar string) map[[2]int]string {
	newImage := map[[2]int]string{}

	var minRow, maxRow, minCol, maxCol int
	for coord := range image {
		if coord[0] < minRow {
			minRow = coord[0]
		}
		if coord[0] > maxRow {
			maxRow = coord[0]
		}
		if coord[1] < minCol {
			minCol = coord[1]
		}
		if coord[1] > maxCol {
			maxCol = coord[1]
		}
	}

	for r := minRow - 1; r <= maxRow+1; r++ {
		for c := minCol - 1; c <= maxCol+1; c++ {
			index := d.getAlgIndex(image, r, c, infiniteChar)
			newImage[[2]int{r, c}] = algorithm[index]
		}
	}

	return newImage
}

func (d Day20) Part1(input string) (string, error) {
	algorithm, image := d.parseInputMap(input)
	infiniteChar := "."

	for step := 0; step < 2; step++ {
		image = d.enhanceImageMap(image, algorithm, infiniteChar)
		if algorithm[0] == "#" {
			if infiniteChar == "." {
				infiniteChar = "#"
			} else {
				infiniteChar = "."
			}
		}
	}

	count := 0
	for _, val := range image {
		if val == "#" {
			count++
		}
	}

	return fmt.Sprintf("%d", count), nil
}

func (d Day20) Part2(input string) (string, error) {
	algorithm, image := d.parseInputMap(input)
	infiniteChar := "."

	for step := 0; step < 50; step++ {
		image = d.enhanceImageMap(image, algorithm, infiniteChar)
		if algorithm[0] == "#" {
			if infiniteChar == "." {
				infiniteChar = "#"
			} else {
				infiniteChar = "."
			}
		}
	}

	count := 0
	for _, val := range image {
		if val == "#" {
			count++
		}
	}

	return fmt.Sprintf("%d", count), nil
}

func init() {
	solve.Register(Day20{})
}
