package solve2024

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/maths"
)

type Day2 struct {
}

func (d Day2) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 2}
}

func (d Day2) parseInput(input string) [][]int {
	scanner := bufio.NewScanner(strings.NewReader(input))
	parsedInput := make([][]int, 0)
	for scanner.Scan() {
		parsedInput = append(parsedInput, d.parseIntegers(scanner))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return parsedInput
}

func (d Day2) parseIntegers(in *bufio.Scanner) []int {
	numbs := Day2{}.parseStrings(in)
	arr := make([]int, 0)
	for _, n := range numbs {
		val, _ := strconv.Atoi(n)
		arr = append(arr, val)
	}
	return arr
}

func (d Day2) parseStrings(in *bufio.Scanner) []string {
	line := d.parseString(in)
	numbs := strings.Split(line, " ")
	return numbs
}

func (d Day2) parseString(in *bufio.Scanner) string {
	nStr := in.Text()
	nStr = strings.ReplaceAll(nStr, "\r", "")
	nStr = strings.ReplaceAll(nStr, "\n", "")
	nStr = strings.TrimSpace(nStr)
	nStr = strings.Trim(nStr, "\t \n")
	return nStr
}

func (d Day2) decreasingRun(input []int) bool {
	for i := 1; i < len(input); i++ {
		if input[i] > input[i-1] {
			return false
		}
	}
	return true
}

func (d Day2) increasingRun(input []int) bool {
	for i := 1; i < len(input); i++ {
		if input[i] < input[i-1] {
			return false
		}
	}
	return true
}

func (d Day2) maxDiff(input []int, max int) bool {
	for i := 1; i < len(input); i++ {
		if maths.Abs(input[i]-input[i-1]) > max {
			return false
		}
	}
	return true
}

func (d Day2) minDiff(input []int, min int) bool {
	for i := 1; i < len(input); i++ {
		if maths.Abs(input[i]-input[i-1]) < min {
			return false
		}
	}
	return true
}

func (d Day2) validate(data []int) bool {
	decr := d.decreasingRun(data)
	incr := d.increasingRun(data)
	maxd := d.maxDiff(data, 3)
	mind := d.minDiff(data, 1)
	return (decr || incr) && maxd && mind
}

func (d Day2) Part1(input string) (string, error) {
	count := 0

	inputData := d.parseInput(input)

	for _, data := range inputData {
		if d.validate(data) {
			count++
		}
	}

	return fmt.Sprintf("%d", count), nil
}

func (d Day2) validateDamped(data []int) bool {
	for i := 0; i < len(data); i++ {
		dampedData := make([]int, len(data)-1)
		copy(dampedData, data[:i])
		copy(dampedData[i:], data[i+1:])
		if d.validate(dampedData) {
			return true
		}
	}
	return false
}

func (d Day2) Part2(input string) (string, error) {
	count := 0

	inputData := d.parseInput(input)

	for _, data := range inputData {
		if d.validate(data) || d.validateDamped(data) {
			count++
		}
	}

	return fmt.Sprintf("%d", count), nil
}

func init() {
	solve.Register(Day2{})
}
