package solve2024

import (
	"image"
	"strconv"
	"strings"
	"testing"
)

func referenceGuardPatrol(input string, extra image.Point) (int, bool) {
	rows := strings.Split(input, "\n")
	start := image.Point{}
	for y, row := range rows {
		for x := range row {
			if row[x] == '^' {
				start = image.Pt(x, y)
			}
		}
	}
	position, direction := start, 0
	directions := [4]image.Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}}
	states := make(map[[3]int]bool)
	cells := make(map[image.Point]bool)
	for position.Y >= 0 && position.Y < len(rows) && position.X >= 0 && position.X < len(rows[0]) {
		state := [3]int{position.X, position.Y, direction}
		if states[state] {
			return len(cells), true
		}
		states[state] = true
		cells[position] = true
		next := position.Add(directions[direction])
		wall := next.Y >= 0 && next.Y < len(rows) && next.X >= 0 && next.X < len(rows[0]) && rows[next.Y][next.X] == '#'
		if wall || next == extra {
			direction = (direction + 1) % 4
		} else {
			position = next
		}
	}
	return len(cells), false
}

func referenceGuardAnswers(input string) (int, int) {
	visited, _ := referenceGuardPatrol(input, image.Pt(-1, -1))
	loops := 0
	for y, row := range strings.Split(input, "\n") {
		for x := range row {
			if row[x] != '.' {
				continue
			}
			if _, loop := referenceGuardPatrol(input, image.Pt(x, y)); loop {
				loops++
			}
		}
	}
	return visited, loops
}

func TestGuardCandidateFilterAgainstEveryEmptyPlacement(t *testing.T) {
	for _, size := range [][2]int{{2, 3}, {3, 3}} {
		width, height := size[0], size[1]
		for start := 0; start < width*height; start++ {
			for walls := 0; walls < 1<<(width*height-1); walls++ {
				var board strings.Builder
				bit := 0
				for cell := 0; cell < width*height; cell++ {
					if cell > 0 && cell%width == 0 {
						board.WriteByte('\n')
					}
					if cell == start {
						board.WriteByte('^')
						continue
					}
					char := byte('.')
					if walls&(1<<bit) != 0 {
						char = '#'
					}
					board.WriteByte(char)
					bit++
				}
				input := board.String()
				want1, want2 := referenceGuardAnswers(input)
				for part, tc := range []struct {
					fn   func(string) (string, error)
					want int
				}{{(Day6{}).Part1, want1}, {(Day6{}).Part2, want2}} {
					got, err := tc.fn(input)
					if err != nil || got != strconv.Itoa(tc.want) {
						t.Fatalf("part%d board %q: %q %v want%d", part+1, input, got, err, tc.want)
					}
				}
			}
		}
	}
}

func TestGuardAlreadyLoopsAndOffRoutePlacements(t *testing.T) {
	input := ".#.\n#^#\n.#."
	if got, err := (Day6{}).Part1(input); err != nil || got != "1" {
		t.Fatalf("turning in place: %q %v", got, err)
	}
	if got, err := (Day6{}).Part2(input); err != nil || got != "4" {
		t.Fatalf("unvisited corners preserve loop: %q %v", got, err)
	}
}

func TestGuardTrialStampsDoNotNeedHistory(t *testing.T) {
	const input = ".....\n..#..\n.^...\n....."
	grid, err := (Day6{}).parseGrid(input)
	if err != nil {
		t.Fatal(err)
	}
	marks := make([]int, len(grid.ahead))
	for i := range marks {
		marks[i] = 70000
	}
	trial := 70001
	for cell, content := range grid.cells {
		if content != '.' {
			continue
		}
		_, want := referenceGuardPatrol(input, image.Pt(cell%5, cell/5))
		if got := grid.loops(cell, marks, trial); got != want {
			t.Fatalf("trial%d cell%d: %t want%t", trial, cell, got, want)
		}
		trial++
	}
}

func TestGuardGridValidation(t *testing.T) {
	for _, input := range []string{"", "\n", "..", "^^", "^x", "^ ", "^é", "^\n..", "^\n\n", "^\r."} {
		for _, fn := range []func(string) (string, error){(Day6{}).Part1, (Day6{}).Part2} {
			if _, err := fn(input); err == nil {
				t.Fatalf("accepted malformed grid %q", input)
			}
		}
	}
	for _, input := range []string{"^", "^\n", "^\r\n", ".\n^", ".\r\n^\r\n"} {
		want1, want2 := referenceGuardAnswers(strings.TrimSuffix(strings.ReplaceAll(input, "\r\n", "\n"), "\n"))
		for part, tc := range []struct {
			fn   func(string) (string, error)
			want int
		}{{(Day6{}).Part1, want1}, {(Day6{}).Part2, want2}} {
			got, err := tc.fn(input)
			if err != nil || got != strconv.Itoa(tc.want) {
				t.Fatalf("part%d thin grid %q: %q %v want%d", part+1, input, got, err, tc.want)
			}
		}
	}
}
