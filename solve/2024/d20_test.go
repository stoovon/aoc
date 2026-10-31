package solve2024

import (
	"strconv"
	"strings"
	"testing"
)

func raceSyntheticSnake(width, height int) string {
	rows := make([][]byte, height)
	for y := range rows {
		rows[y] = []byte(strings.Repeat("#", width))
	}
	for y := 1; y < height-1; y += 2 {
		for x := 1; x < width-1; x++ {
			rows[y][x] = '.'
		}
		if y+1 < height-1 {
			x := width - 2
			if ((y-1)/2)%2 != 0 {
				x = 1
			}
			rows[y+1][x] = '.'
		}
	}
	rows[1][1] = 'S'
	endX := width - 2
	if ((height-3)/2)%2 != 0 {
		endX = 1
	}
	rows[height-2][endX] = 'E'
	text := make([]string, height)
	for y, row := range rows {
		text[y] = string(row)
	}
	return strings.Join(text, "\n")
}

func raceAllPairDistanceReference(input string, radius int) int {
	rows := strings.Split(input, "\n")
	width := len(rows[0])
	size := width * len(rows)
	distances := make([]int, size)
	start := 0
	for y, row := range rows {
		for x := range row {
			distances[y*width+x] = -1
			if row[x] == 'S' {
				start = y*width + x
			}
		}
	}
	distances[start] = 0
	queue := []int{start}
	for head := 0; head < len(queue); head++ {
		id := queue[head]
		x, y := id%width, id/width
		for _, step := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nx, ny := x+step[0], y+step[1]
			if nx < 0 || nx >= width || ny < 0 || ny >= len(rows) || rows[ny][nx] == '#' {
				continue
			}
			next := ny*width + nx
			if distances[next] < 0 {
				distances[next] = distances[id] + 1
				queue = append(queue, next)
			}
		}
	}
	count := 0
	for a, da := range distances {
		if da < 0 {
			continue
		}
		for b, db := range distances {
			if db < 0 {
				continue
			}
			dx, dy := a%width-b%width, a/width-b/width
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			length := dx + dy
			if length > 0 && length <= radius && db-da-length >= 100 {
				count++
			}
		}
	}
	return count
}

func TestRaceOffsetsMatchAllPairs(t *testing.T) {
	for _, width := range []int{9, 17, 23} {
		for _, height := range []int{7, 15, 21} {
			input := raceSyntheticSnake(width, height)
			for part, solver := range []func(string) (string, error){(Day20{}).Part1, (Day20{}).Part2} {
				radius := 2
				if part == 1 {
					radius = 20
				}
				want := raceAllPairDistanceReference(input, radius)
				got, err := solver(input)
				if err != nil || got != strconv.Itoa(want) {
					t.Fatalf("%dx%d radius %d: got %q %v want %d", width, height, radius, got, err, want)
				}
			}
		}
	}
}

func raceVirtualEdgeReference(input string, radius, minSaving int) (uint64, bool) {
	rows := strings.Split(input, "\n")
	width := len(rows[0])
	size := width * len(rows)
	open := make([]bool, size)
	start, end := 0, 0
	edges := make([][]int, size)
	for y, row := range rows {
		for x := range row {
			id := y*width + x
			open[id] = row[x] != '#'
			if row[x] == 'S' {
				start = id
			}
			if row[x] == 'E' {
				end = id
			}
		}
	}
	for a := 0; a < size; a++ {
		if !open[a] {
			continue
		}
		for b := 0; b < size; b++ {
			if !open[b] {
				continue
			}
			dx, dy := a%width-b%width, a/width-b/width
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			if dx+dy == 1 {
				edges[a] = append(edges[a], b)
			}
		}
	}
	shortest := func(a, b, length int) int {
		distances := make([]int, size)
		for i := range distances {
			distances[i] = size + 1
		}
		distances[start] = 0
		for pass := 0; pass < size; pass++ {
			for id := 0; id < size; id++ {
				for _, next := range edges[id] {
					if distances[id]+1 < distances[next] {
						distances[next] = distances[id] + 1
					}
				}
				if id == a && distances[id]+length < distances[b] {
					distances[b] = distances[id] + length
				}
			}
		}
		return distances[end]
	}
	normal := shortest(-1, -1, 0)
	if normal > size {
		return 0, false
	}
	var count uint64
	for a := 0; a < size; a++ {
		if !open[a] {
			continue
		}
		for b := 0; b < size; b++ {
			if !open[b] || a == b {
				continue
			}
			dx, dy := a%width-b%width, a/width-b/width
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			length := dx + dy
			if length <= radius && normal-shortest(a, b, length) >= minSaving {
				count++
			}
		}
	}
	return count, true
}

func TestRaceCompleteRouteAgainstVirtualEdges(t *testing.T) {
	for code := 0; code < 256; code++ {
		width, height := 3, 3
		if code >= 128 {
			width, height = 4, 3
		}
		cells := []byte(strings.Repeat(".", width*height))
		seed := uint32(code)
		bit := 0
		for id := range cells {
			if id == 0 || id == width-1 {
				continue
			}
			wall := false
			if code < 128 {
				wall = code&(1<<bit) != 0
				bit++
			} else {
				seed = seed*1664525 + 1013904223
				wall = seed>>30 == 0
			}
			if wall {
				cells[id] = '#'
			}
		}
		cells[0] = 'S'
		cells[width-1] = 'E'
		rows := make([]string, height)
		for y := range rows {
			rows[y] = string(cells[y*width : (y+1)*width])
		}
		input := strings.Join(rows, "\n")
		track, err := parseRaceTrack(input)
		if err != nil {
			t.Fatal(err)
		}
		for _, radius := range []int{1, 2, 4} {
			for _, saving := range []int{1, 2, 3} {
				want, reachable := raceVirtualEdgeReference(input, radius, saving)
				got, err := track.countCheats(radius, saving)
				if !reachable {
					if err == nil {
						t.Fatalf("case %d accepted unreachable finish", code)
					}
				} else if err != nil || got != want {
					t.Fatalf("case %d radius %d saving %d: got %d %v want %d", code, radius, saving, got, err, want)
				}
			}
		}
	}
}

func TestRaceDeadEndDoesNotCountAsProgress(t *testing.T) {
	rows := strings.Split(strings.ReplaceAll(raceSyntheticSnake(73, 9), "E", "."), "\n")
	row := []byte(rows[1])
	row[3] = 'E'
	rows[1] = string(row)
	input := strings.Join(rows, "\n")
	for _, solver := range []func(string) (string, error){(Day20{}).Part1, (Day20{}).Part2} {
		if got, err := solver(input); err != nil || got != "0" {
			t.Fatalf("two-step normal race: %q %v", got, err)
		}
	}
}

func TestRaceParserAndBoundaries(t *testing.T) {
	for _, input := range []string{"", "...", "S..", "..E", "SSE", "SEE", "SXE", "SéE", " S.E", "S.E\n..", "S.E\n\n..."} {
		for _, solver := range []func(string) (string, error){(Day20{}).Part1, (Day20{}).Part2} {
			if _, err := solver(input); err == nil {
				t.Fatalf("accepted %q", input)
			}
		}
	}
	for _, input := range []string{"S.E", "S\n.\nE", "S.E\r\n...\r\n"} {
		for _, solver := range []func(string) (string, error){(Day20{}).Part1, (Day20{}).Part2} {
			if got, err := solver(input); err != nil || got != "0" {
				t.Fatalf("open boundary: %q %v", got, err)
			}
		}
	}
	for _, solver := range []func(string) (string, error){(Day20{}).Part1, (Day20{}).Part2} {
		if _, err := solver("S#E"); err == nil {
			t.Fatal("expected unreachable finish error")
		}
	}
}
