package solve2024

import (
	"fmt"
	"image"
	"math"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/acstrings"
	"aoc/utils/grids"
	"aoc/utils/maps"
)

type Day21 struct {
	minDistanceCache map[string]int
	pathsCache       map[string][]string
}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 21}
}

var (
	dirMap = grids.Directions("^v><")

	numericKeypadPerButton = grids.NewGridOptions().Parse(`
789
456
123
#0A
`).PointsFromTopLeft()

	dirKeypadPerButton = grids.NewGridOptions().Parse(`
#^A
<v>
`).PointsFromTopLeft()

	dirKeypadPerPosition     = maps.Inverted(dirKeypadPerButton)
	numericKeypadPerPosition = maps.Inverted(numericKeypadPerButton)
)

func (d Day21) cost(str string, depth int) (res int) {
	for i := 0; i < len(str)-1; i++ {
		currPairCost := d.pairCost(rune(str[i]), rune(str[i+1]), numericKeypadPerButton, numericKeypadPerPosition, depth)
		res += currPairCost
	}
	return
}

func getCachedCost(key string, cache map[string]int, computeCost func() int) int {
	if cost, exists := cache[key]; exists {
		return cost
	}
	cost := computeCost()
	cache[key] = cost
	return cost
}

func (d Day21) pairCost(a, b rune, charToIndex map[rune]image.Point, indexToChar map[image.Point]rune, depth int) int {
	code := 'd'
	if _, ok := charToIndex['0']; ok {
		code = 'n'
	}
	key := fmt.Sprintf("%c%c%c%d", a, b, code, depth)

	return getCachedCost(key, d.minDistanceCache, func() int {
		if depth == 0 {
			minCost := math.MaxInt
			for _, path := range d.allPaths(a, b, charToIndex, indexToChar) {
				minCost = min(minCost, len(path))
			}
			return minCost
		}

		paths := d.allPaths(a, b, charToIndex, indexToChar)
		minCost := math.MaxInt
		for _, path := range paths {
			path = "A" + path
			currCost := 0
			for i := 0; i < len(path)-1; i++ {
				currCost += d.pairCost(rune(path[i]), rune(path[i+1]), dirKeypadPerButton, dirKeypadPerPosition, depth-1)
			}
			minCost = min(minCost, currCost)
		}
		return minCost
	})
}

func (d Day21) allPaths(a, b rune, charToIndex map[rune]image.Point, indexToChar map[image.Point]rune) []string {
	key := fmt.Sprintf("%c %c", a, b)
	if paths, ok := d.pathsCache[key]; ok {
		return paths
	}
	var paths []string
	dfs(charToIndex[a], charToIndex[b], []rune{}, charToIndex, indexToChar, make(map[image.Point]bool), &paths)
	d.pathsCache[key] = paths
	return paths
}

func dfs(curr, end image.Point, path []rune, charToIndex map[rune]image.Point, indexToChar map[image.Point]rune, visited map[image.Point]bool, allPaths *[]string) {
	if curr == end {
		*allPaths = append(*allPaths, string(append(path, 'A')))
		return
	}

	visited[curr] = true
	defer func() { visited[curr] = false }()

	for char, dir := range dirMap {
		next := image.Point{X: curr.X + dir.X, Y: curr.Y + dir.Y}
		if _, ok := indexToChar[next]; ok && !visited[next] {
			dfs(next, end, append(path, char), charToIndex, indexToChar, visited, allPaths)
		}
	}
}

func (d Day21) solve(input string, depth int) (res int) {
	d.minDistanceCache = make(map[string]int)
	d.pathsCache = make(map[string][]string)

	for _, str := range strings.Fields(input) {
		res += d.cost("A"+str, depth) * acstrings.MustInt(str[:len(str)-1])
	}
	return
}

func (d Day21) Part1(input string) (string, error) {
	result := d.solve(input, 2)

	return strconv.Itoa(result), nil
}

func (d Day21) Part2(input string) (string, error) {
	result := d.solve(input, 25)

	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day21{})
}
