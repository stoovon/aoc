package solve2022

import (
	"aoc/solve"
	"fmt"
	"strconv"
	"strings"
)

type Day23 struct{}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 23}
}

func parseInput(input string) ([][]rune, [][]int) {
	grid := [][]rune{}
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		grid = append(grid, []rune(line))
	}

	elfPositions := [][]int{}
	for rowIdx, row := range grid {
		for colIdx, cell := range row {
			if cell == '#' {
				elfPositions = append(elfPositions, []int{rowIdx, colIdx})
			}
		}
	}

	return grid, elfPositions
}

func simulateElves(grid [][]rune, elfPositions [][]int, stopAtRound int) (int, int) {
	directions := map[string][][]int{
		"N": {{-1, 0}, {-1, 1}, {-1, -1}},
		"S": {{1, 0}, {1, 1}, {1, -1}},
		"W": {{0, -1}, {-1, -1}, {1, -1}},
		"E": {{0, 1}, {-1, 1}, {1, 1}},
	}
	allDirections := [][]int{}
	for _, dirSet := range directions {
		allDirections = append(allDirections, dirSet...)
	}
	directionQueue := []string{"N", "S", "W", "E"}

	addCoordinates := func(coord1, coord2 []int) []int {
		return []int{coord1[0] + coord2[0], coord1[1] + coord2[1]}
	}

	round := 0
	for {
		round++
		elfPositionSet := map[[2]int]bool{}
		for _, position := range elfPositions {
			elfPositionSet[[2]int{position[0], position[1]}] = true
		}

		proposedMoves := map[[2]int][]int{}
		for elfIdx, position := range elfPositions {
			canMove := false
			for _, direction := range allDirections {
				if elfPositionSet[[2]int{position[0] + direction[0], position[1] + direction[1]}] {
					canMove = true
					break
				}
			}
			if !canMove {
				continue
			}

			for _, dir := range directionQueue {
				canMoveInDirection := true
				for _, direction := range directions[dir] {
					if elfPositionSet[[2]int{position[0] + direction[0], position[1] + direction[1]}] {
						canMoveInDirection = false
						break
					}
				}
				if canMoveInDirection {
					newPosition := addCoordinates(position, directions[dir][0])
					proposedMoves[[2]int{newPosition[0], newPosition[1]}] = append(proposedMoves[[2]int{newPosition[0], newPosition[1]}], elfIdx)
					break
				}
			}
		}

		directionQueue = append(directionQueue[1:], directionQueue[0])

		anyElfMoved := false
		for newPosition, elfIndices := range proposedMoves {
			if len(elfIndices) == 1 {
				elfPositions[elfIndices[0]] = []int{newPosition[0], newPosition[1]}
				anyElfMoved = true
			}
		}

		if stopAtRound > 0 && round == stopAtRound {
			minRow, maxRow := elfPositions[0][0], elfPositions[0][0]
			minCol, maxCol := elfPositions[0][1], elfPositions[0][1]
			for _, position := range elfPositions {
				if position[0] < minRow {
					minRow = position[0]
				}
				if position[0] > maxRow {
					maxRow = position[0]
				}
				if position[1] < minCol {
					minCol = position[1]
				}
				if position[1] > maxCol {
					maxCol = position[1]
				}
			}
			emptySpaces := (maxRow-minRow+1)*(maxCol-minCol+1) - len(elfPositions)
			return emptySpaces, round
		}

		if !anyElfMoved {
			return -1, round
		}
	}
}

func (d Day23) Part1(input string) (string, error) {
	_, elfPositions := parseInput(input)
	emptySpaces, _ := simulateElves(nil, elfPositions, 10)
	return fmt.Sprintf("%d", emptySpaces), nil
}

func (d Day23) Part2(input string) (string, error) {
	_, elfPositions := parseInput(input)
	_, rounds := simulateElves(nil, elfPositions, -1)
	return strconv.Itoa(rounds), nil
}

func init() {
	solve.Register(Day23{})
}
