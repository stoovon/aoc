package solve2020

import (
	"fmt"
	"image"
	"math"
	"regexp"
	"strings"

	"aoc/solve"
)

type Day20 struct {
}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 20}
}

type tile []string

func (t tile) Col(i int) (col string) {
	for _, s := range t {
		col += string(s[i])
	}
	return
}

func (t tile) Edges() []string {
	return []string{t[0], t.Col(len(t[0]) - 1), t[len(t)-1], t.Col(0)}
}

func (t tile) Rotations() []tile {
	rotations := make([]tile, 0, 8)

	current := t
	for i := 0; i < 4; i++ {
		rotations = append(rotations, current)
		current = current.rotate90()
	}

	for _, r := range rotations[:4] {
		rotations = append(rotations, r.flip())
	}

	return rotations
}

func (t tile) rotate90() tile {
	size := len(t)
	rotated := make(tile, size)
	for i := 0; i < size; i++ {
		var row string
		for j := size - 1; j >= 0; j-- {
			row += string(t[j][i])
		}
		rotated[i] = row
	}
	return rotated
}

func (t tile) flip() tile {
	flipped := make(tile, len(t))
	for i, row := range t {
		flipped[i] = reverse(row)
	}
	return flipped
}

func reverse(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		builder.WriteByte(s[i])
	}
	return builder.String()
}

func parseTiles(input string) (map[int]tile, map[string]int, error) {
	tiles := map[int]tile{}
	counts := map[string]int{}
	sections := strings.Split(strings.TrimSpace(input), "\n\n")

	for _, section := range sections {
		var id int
		_, err := fmt.Sscanf(section, "Tile %d:", &id)
		if err != nil {
			return nil, nil, err
		}

		t := tile(strings.Split(section, "\n")[1:])
		tiles[id] = t
		for _, edge := range t.Edges() {
			counts[edge]++
			counts[reverse(edge)]++
		}
	}

	return tiles, counts, nil
}

func findCorners(tiles map[int]tile, counts map[string]int) (int, error) {
	part1 := 1
	for id, t := range tiles {
		uniqueEdges := 0
		for _, edge := range t.Edges() {
			if counts[edge] == 1 {
				uniqueEdges++
			}
		}
		if uniqueEdges == 2 {
			part1 *= id
		}
	}
	return part1, nil
}

func checkMonster(r tile, x, y int, monster []string) bool {
	for i, s := range monster {
		if match, _ := regexp.MatchString(s, r[y+i][x:x+len(s)]); !match {
			return false
		}
	}
	return true
}

func countMonsters(img tile, monster []string) int {
	count := 0
	for _, r := range img.Rotations() {
		for y := 0; y < len(r)-len(monster); y++ {
			for x := 0; x < len(r[0])-len(monster[0]); x++ {
				if checkMonster(r, x, y, monster) {
					count++
				}
			}
		}
	}
	return count
}

func findAndPlaceTile(tiles map[int]tile, counts map[string]int, order map[image.Point]tile, x, y, tileSize int) (tile, error) {
	for id, t := range tiles {
		for _, r := range t.Rotations() {
			if (y == 0 && counts[r[0]] == 1 || y != 0 && r[0] == order[image.Point{X: x, Y: y - 1}][tileSize-1]) &&
				(x == 0 && counts[r.Col(0)] == 1 || x != 0 && r.Col(0) == order[image.Point{X: x - 1, Y: y}].Col(tileSize-1)) {
				delete(tiles, id)
				return r, nil
			}
		}
	}
	return nil, fmt.Errorf("no matching tile found for position (%d, %d)", x, y)
}

func (d Day20) solve(input string) (int, int, error) {
	tiles, counts, err := parseTiles(input)
	if err != nil {
		return 0, 0, err
	}

	part1, err := findCorners(tiles, counts)
	if err != nil {
		return 0, 0, err
	}

	imageSize, tileSize := int(math.Sqrt(float64(len(tiles)))), 10
	img := make(tile, imageSize*(tileSize-2))
	order := map[image.Point]tile{}

	for y := 0; y < imageSize; y++ {
		for x := 0; x < imageSize; x++ {
			r, err := findAndPlaceTile(tiles, counts, order, x, y, tileSize)
			if err != nil {
				return 0, 0, err
			}
			order[image.Point{X: x, Y: y}] = r
			for i := 0; i < tileSize-2; i++ {
				img[(tileSize-2)*y+i] += r[i+1][1 : tileSize-1]
			}
		}
	}

	monster := []string{"..................#.", "#....##....##....###", ".#..#..#..#..#..#..."}
	count := countMonsters(img, monster)
	part2 := strings.Count(strings.Join(img, ""), "#") - count*strings.Count(strings.Join(monster, ""), "#")

	return part1, part2, nil
}

func (d Day20) Part1(input string) (string, error) {
	part1, _, err := d.solve(input)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(part1), nil
}

func (d Day20) Part2(input string) (string, error) {
	_, part2, err := d.solve(input)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(part2), nil
}

func init() {
	solve.Register(Day20{})
}
