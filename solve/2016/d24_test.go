package solve2016

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func referenceDuctPermutation(paths [][]int, order []int, used []bool, returnHome bool) int {
	if len(order) == len(paths) {
		total := 0
		for i := 1; i < len(order); i++ {
			total += paths[order[i-1]][order[i]]
		}
		if returnHome {
			total += paths[order[len(order)-1]][0]
		}
		return total
	}
	best := math.MaxInt
	for next := 1; next < len(paths); next++ {
		if used[next] {
			continue
		}
		used[next] = true
		cost := referenceDuctPermutation(paths, append(order, next), used, returnHome)
		used[next] = false
		if cost < best {
			best = cost
		}
	}
	return best
}

func TestDuctMemoAgainstAllRouteOrders(t *testing.T) {
	seed := uint32(314159)
	for count := 1; count <= 7; count++ {
		for trial := 0; trial < 20; trial++ {
			paths := make([][]int, count)
			for i := range paths {
				paths[i] = make([]int, count)
			}
			for i := range paths {
				for j := i + 1; j < count; j++ {
					seed = seed*1664525 + 1013904223
					paths[i][j] = int(seed%37) + 1
					paths[j][i] = paths[i][j]
				}
			}
			for _, home := range []bool{false, true} {
				used := make([]bool, count)
				used[0] = true
				want := referenceDuctPermutation(paths, []int{0}, used, home)
				got, err := (Day24{}).shortestRoute(paths, home)
				if err != nil || got != want {
					t.Fatalf("%d points home%t: %d %v want%d", count, home, got, err, want)
				}
			}
		}
	}
}

func TestDuctRejectsUnreachableOrOverflowingRoutes(t *testing.T) {
	for _, paths := range [][][]int{
		nil, {{0}, {1, 0}}, {{0, math.MaxInt}, {math.MaxInt, 0}}, {{0, -1}, {-1, 0}},
		{{0, math.MaxInt - 1}, {math.MaxInt - 1, 0}},
	} {
		if _, err := (Day24{}).shortestRoute(paths, true); err == nil {
			t.Fatal("accepted invalid or overflowing route")
		}
	}
}

type referenceDuctState struct{ x, y, visited int }

func referenceDuctGridRoute(input string, home bool) int {
	rows := strings.Split(input, "\n")
	var start referenceDuctState
	all := 0
	for y, row := range rows {
		for x := range row {
			if row[x] >= '0' && row[x] <= '9' {
				all |= 1 << int(row[x]-'0')
				if row[x] == '0' {
					start = referenceDuctState{x, y, 1}
				}
			}
		}
	}
	type entry struct {
		state    referenceDuctState
		distance int
	}
	queue := []entry{{start, 0}}
	seen := map[referenceDuctState]bool{start: true}
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		state := current.state
		if state.visited == all && (!home || (state.x == start.x && state.y == start.y)) {
			return current.distance
		}
		for _, offset := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			next := referenceDuctState{state.x + offset[0], state.y + offset[1], state.visited}
			if next.y < 0 || next.y >= len(rows) || next.x < 0 || next.x >= len(rows[0]) {
				continue
			}
			cell := rows[next.y][next.x]
			if cell == '#' {
				continue
			}
			if cell >= '0' && cell <= '9' {
				next.visited |= 1 << int(cell-'0')
			}
			if seen[next] {
				continue
			}
			seen[next] = true
			queue = append(queue, entry{next, current.distance + 1})
		}
	}
	return -1
}

func TestDuctCompressedRoutesAgainstGridStateSearch(t *testing.T) {
	for zero := 0; zero < 6; zero++ {
		for two := 0; two < 6; two++ {
			for seven := 0; seven < 6; seven++ {
				if zero == two || zero == seven || two == seven {
					continue
				}
				for walls := 0; walls < 8; walls++ {
					var input strings.Builder
					bit := 0
					for cell := 0; cell < 6; cell++ {
						if cell == 3 {
							input.WriteByte('\n')
						}
						switch cell {
						case zero:
							input.WriteByte('0')
						case two:
							input.WriteByte('2')
						case seven:
							input.WriteByte('7')
						default:
							content := byte('.')
							if walls&(1<<bit) != 0 {
								content = '#'
							}
							input.WriteByte(content)
							bit++
						}
					}
					text := input.String()
					for part, fn := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
						want := referenceDuctGridRoute(text, part == 1)
						got, err := fn(text)
						if want < 0 {
							if err == nil {
								t.Fatalf("disconnected maze %q accepted", text)
							}
							continue
						}
						if err != nil || got != strconv.Itoa(want) {
							t.Fatalf("part%d maze %q: %q %v want%d", part+1, text, got, err, want)
						}
					}
				}
			}
		}
	}
}

func TestDuctCrossingOtherNumberedPoints(t *testing.T) {
	for part, tc := range []struct {
		fn   func(string) (string, error)
		want string
	}{{(Day24{}).Part1, "4"}, {(Day24{}).Part2, "8"}} {
		got, err := tc.fn("0.7.2")
		if err != nil || got != tc.want {
			t.Fatalf("part%d crossed numbered point: %q %v", part+1, got, err)
		}
	}
}

func TestDuctMazeValidation(t *testing.T) {
	for _, input := range []string{"", "\n", "...", "00", "0x", "0 ", "0é", "0\n..", "0\n\n", "0\r.", "0#2"} {
		for _, fn := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
			if _, err := fn(input); err == nil {
				t.Fatalf("accepted malformed or disconnected maze %q", input)
			}
		}
	}
	for _, input := range []string{"0", "0\n", "0\r\n", "0\n.", "0\r\n.\r\n"} {
		for _, fn := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
			if got, err := fn(input); err != nil || got != "0" {
				t.Fatalf("single numbered point %q: %q %v", input, got, err)
			}
		}
	}
}

func TestDuctRouteArithmeticAtIntLimit(t *testing.T) {
	paths := [][]int{{0, math.MaxInt - 1, math.MaxInt - 1}, {math.MaxInt - 1, 0, 1}, {math.MaxInt - 1, 1, 0}}
	if got, err := (Day24{}).shortestRoute(paths, false); err != nil || got != math.MaxInt {
		t.Fatalf("exact maximum route: %d %v", got, err)
	}
	if _, err := (Day24{}).shortestRoute(paths, true); err == nil {
		t.Fatal("accepted overflowing round trip")
	}
	paths = [][]int{{0, math.MaxInt - 1, 1}, {math.MaxInt - 1, 0, 3}, {1, 3, 0}}
	if got, err := (Day24{}).shortestRoute(paths, false); err != nil || got != 4 {
		t.Fatalf("finite alternative to overflow: %d %v", got, err)
	}
}

func TestDuctAllTenNumberedPoints(t *testing.T) {
	for part, tc := range []struct {
		fn   func(string) (string, error)
		want string
	}{{(Day24{}).Part1, "9"}, {(Day24{}).Part2, "18"}} {
		if got, err := tc.fn("0123456789"); err != nil || got != tc.want {
			t.Fatalf("part%d all digit labels: %q %v", part+1, got, err)
		}
	}
}
