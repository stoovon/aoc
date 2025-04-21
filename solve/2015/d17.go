package solve2015

import (
	"strconv"
	"strings"

	"aoc/solve"
)

type Day17 struct {
}

func (d Day17) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 17}
}

func countCombinations(containerSizes []int, capacity int) int {
	combinations := 0
	n := len(containerSizes)

	for i := 0; i < (1 << n); i++ {
		stored := 0
		for j := 0; j < n; j++ {
			if i&(1<<j) != 0 {
				stored += containerSizes[j]
			}
		}
		if stored == capacity {
			combinations++
		}
	}
	return combinations
}

func countMinimalCombinations(containerSizes []int, targetSum int) int {
	n := len(containerSizes)
	for k := 1; k <= n; k++ {
		count := 0
		combinations := generateCombinations(containerSizes, k)
		for _, combination := range combinations {
			sum := 0
			for _, value := range combination {
				sum += value
			}
			if sum == targetSum {
				count++
			}
		}
		if count > 0 {
			return count
		}
	}
	return 0
}

func generateCombinations(arr []int, k int) [][]int {
	var result [][]int
	var helper func(start int, comb []int)
	helper = func(start int, comb []int) {
		if len(comb) == k {
			tmp := make([]int, k)
			copy(tmp, comb)
			result = append(result, tmp)
			return
		}
		for i := start; i < len(arr); i++ {
			helper(i+1, append(comb, arr[i]))
		}
	}
	helper(0, []int{})
	return result
}

func (d Day17) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	containerSizes := make([]int, len(lines))
	for i, line := range lines {
		containerSizes[i], _ = strconv.Atoi(line)
	}
	return strconv.Itoa(countCombinations(containerSizes, 150)), nil
}

func (d Day17) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	containerSizes := make([]int, len(lines))
	for i, line := range lines {
		containerSizes[i], _ = strconv.Atoi(line)
	}
	return strconv.Itoa(countMinimalCombinations(containerSizes, 150)), nil
}

func init() {
	solve.Register(Day17{})
}
