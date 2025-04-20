package solve2024

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day22 struct {
}

func (d Day22) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 22}
}

func parseInput(input string) (nums []int) {
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		temp, _ := strconv.Atoi(line)
		nums = append(nums, temp)
	}
	return
}

func getNthSecretNumber(num int, n int) int {
	for ; n > 0; n-- {
		num = ((num << 6) ^ num) & (16777216 - 1)
		num = ((num >> 5) ^ num) & (16777216 - 1)
		num = ((num << 11) ^ num) & (16777216 - 1)
	}
	return num
}

func (d Day22) Part1(input string) (string, error) {
	nums := parseInput(input)
	res := 0

	for _, num := range nums {
		temp := getNthSecretNumber(num, 2000)
		res += temp
	}

	return strconv.Itoa(res), nil
}

type pair struct {
	delta   int
	bananas int
}

type window struct {
	a, b, c, d int
}

func updateCache(num int, n int, cache map[window]int) (delBan []pair) {
	visited := make(map[window]bool)
	delBan = make([]pair, 0, n-1)
	currWindow := make([]int, 0, 8)
	prev := num % 10

	for ; n > 1; n-- {
		num = ((num << 6) ^ num) & (16777216 - 1)
		num = ((num >> 5) ^ num) & (16777216 - 1)
		num = ((num << 11) ^ num) & (16777216 - 1)

		delta, bananas := num%10-prev, num%10
		delBan = append(delBan, pair{delta: delta, bananas: bananas})
		currWindow = append(currWindow, delta)

		if len(currWindow) == 4 {
			key := window{a: currWindow[0], b: currWindow[1], c: currWindow[2], d: currWindow[3]}
			if !visited[key] {
				cache[key] += num % 10
				visited[key] = true
			}
			currWindow = currWindow[1:]
		}
		prev = num % 10
	}
	return
}

func (d Day22) Part2(input string) (string, error) {
	nums := parseInput(input)
	cache := make(map[window]int)
	res := 0

	for _, num := range nums {
		updateCache(num, 2000, cache)
	}

	for _, value := range cache {
		res = max(res, value)
	}

	return strconv.Itoa(res), nil
}

func init() {
	solve.Register(Day22{})
}
