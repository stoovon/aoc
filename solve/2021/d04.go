package solve2021

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 4}
}

type board struct {
	nums  [5][5]int
	marks [5][5]bool
}

func (d Day4) parseInput(input string) ([]int, []board, error) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	if len(parts) < 2 {
		return nil, nil, errors.New("invalid input")
	}

	numStrs := strings.Split(parts[0], ",")
	draws := make([]int, len(numStrs))
	for i, s := range numStrs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, nil, err
		}
		draws[i] = n
	}

	var boards []board
	for _, bstr := range parts[1:] {
		var b board
		lines := strings.Split(strings.TrimSpace(bstr), "\n")
		for i, line := range lines {
			fields := strings.Fields(line)
			for j, f := range fields {
				n, err := strconv.Atoi(f)
				if err != nil {
					return nil, nil, err
				}
				b.nums[i][j] = n
			}
		}
		boards = append(boards, b)
	}
	return draws, boards, nil
}

func mark(b *board, n int) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if b.nums[i][j] == n {
				b.marks[i][j] = true
			}
		}
	}
}

func isWinner(b *board) bool {
	for i := 0; i < 5; i++ {
		row, col := true, true
		for j := 0; j < 5; j++ {
			row = row && b.marks[i][j]
			col = col && b.marks[j][i]
		}
		if row || col {
			return true
		}
	}
	return false
}

func unmarkedSum(b *board) int {
	sum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if !b.marks[i][j] {
				sum += b.nums[i][j]
			}
		}
	}
	return sum
}

func (d Day4) Part1(input string) (string, error) {
	draws, boards, err := d.parseInput(input)
	if err != nil {
		return "", err
	}
	for _, n := range draws {
		for i := range boards {
			mark(&boards[i], n)
			if isWinner(&boards[i]) {
				score := unmarkedSum(&boards[i]) * n
				return fmt.Sprintf("%d", score), nil
			}
		}
	}
	return "", errors.New("no winner")
}

func (d Day4) Part2(input string) (string, error) {
	draws, boards, err := d.parseInput(input)
	if err != nil {
		return "", err
	}
	won := make([]bool, len(boards))
	winners := 0
	var lastScore int
	for _, n := range draws {
		for i := range boards {
			if won[i] {
				continue
			}
			mark(&boards[i], n)
			if isWinner(&boards[i]) {
				won[i] = true
				winners++
				if winners == len(boards) {
					lastScore = unmarkedSum(&boards[i]) * n
					return fmt.Sprintf("%d", lastScore), nil
				}
			}
		}
	}
	return "", errors.New("no last winner")
}

func init() {
	solve.Register(Day4{})
}
