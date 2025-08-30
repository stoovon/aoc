package solve2021

import (
	"aoc/solve"
	"aoc/utils/maths"
	"strconv"
	"strings"
)

var (
	RoomA = [][2]int{{3, 2}, {3, 3}, {3, 4}, {3, 5}}
	RoomB = [][2]int{{5, 2}, {5, 3}, {5, 4}, {5, 5}}
	RoomC = [][2]int{{7, 2}, {7, 3}, {7, 4}, {7, 5}}
	RoomD = [][2]int{{9, 2}, {9, 3}, {9, 4}, {9, 5}}

	Rooms   = [][][2]int{RoomA, RoomB, RoomC, RoomD}
	Hallway = [][2]int{{1, 1}, {2, 1}, {4, 1}, {6, 1}, {8, 1}, {10, 1}, {11, 1}}
)

var memo map[Burrow]int

const (
	Amber  byte = 'A'
	Bronze byte = 'B'
	Copper byte = 'C'
	Desert byte = 'D'
)

var AmphipodConsumption = map[byte]int{
	Amber:  1,
	Bronze: 10,
	Copper: 100,
	Desert: 1000,
}

var AmphipodRoom = map[byte][][2]int{
	Amber:  RoomA,
	Bronze: RoomB,
	Copper: RoomC,
	Desert: RoomD,
}

func loc(xy [2]int) int {
	return xy[1]*14 + xy[0]
}

type Step struct {
	burrow Burrow
	energy int
}

type Burrow string

// Move amphipod from its room to the hallway
func (burrow Burrow) moveToHallway(a [2]int) []Step {
	results := []Step{}
	reachable := [][2]int{}

	var room [][2]int
	if r, ok := AmphipodRoom[burrow[loc(a)]]; !ok {
		panic("bad room")
	} else {
		room = r
	}
	roomX := room[0][0]

	// Check if amphipod can move to the hallway
	movable := false
	slots := room
	if burrow.folded() {
		slots = room[:2]
	}
	if roomX == a[0] {
		for _, r := range slots {
			if !burrow.empty(r) && burrow[loc(r)] != burrow[loc(a)] {
				movable = true
				break
			}
		}
	} else {
		movable = true
	}

	if !movable {
		return results
	}

	// Find reachable hallway positions
	for _, h := range Hallway {
		if h[0] < a[0] {
			if burrow.empty(h) {
				reachable = append(reachable, h)
			} else {
				reachable = reachable[:0]
			}
		} else {
			if burrow.empty(h) {
				reachable = append(reachable, h)
			} else {
				break
			}
		}
	}

	// Generate steps to move to reachable hallway positions
	for _, n := range reachable {
		next := Burrow(burrow[:loc(n)] + Burrow(burrow[loc(a)]) + burrow[loc(n)+1:loc(a)] + "." + burrow[loc(a)+1:])
		energy := (maths.Abs(a[0]-n[0]) + maths.Abs(a[1]-n[1])) * AmphipodConsumption[burrow[loc(a)]]
		results = append(results, Step{next, energy})
	}

	return results
}

// Check if the burrow is in a folded state
func (burrow Burrow) folded() bool {
	return len(burrow)/14 <= 5
}

// Move amphipod from the hallway to its room
func (burrow Burrow) moveToRoom(a [2]int) []Step {
	results := []Step{}
	reachable := true

	var room [][2]int
	if r, ok := AmphipodRoom[burrow[loc(a)]]; !ok {
		panic("bad room")
	} else {
		room = r
	}
	roomX := room[0][0]

	// Check if target room is reachable from the hallway
	if a[0] < roomX {
		for _, h := range Hallway {
			if h[0] <= a[0] {
				continue
			}
			if h[0] > roomX {
				break
			}
			if !burrow.empty(h) {
				reachable = false
				break
			}
		}
	} else {
		for i := len(Hallway) - 1; i >= 0; i-- {
			h := Hallway[i]
			if h[0] >= a[0] {
				continue
			}
			if h[0] < roomX {
				break
			}
			if !burrow.empty(h) {
				reachable = false
				break
			}
		}
	}

	// Check if room is available
	slots := room
	if burrow.folded() {
		slots = room[:2]
	}
	for _, r := range slots {
		if !burrow.empty(r) && burrow[loc(r)] != burrow[loc(a)] {
			reachable = false
			break
		}
	}

	if reachable {
		for i := len(slots) - 1; i >= 0; i-- {
			r := room[i]
			if burrow.empty(r) {
				next := Burrow(burrow[:loc(a)] + "." + burrow[loc(a)+1:loc(r)] + Burrow(burrow[loc(a)]) + burrow[loc(r)+1:])
				e := (maths.Abs(a[0]-r[0]) + maths.Abs(a[1]-r[1])) * AmphipodConsumption[burrow[loc(a)]]
				results = append(results, Step{next, e})
				break
			}
		}
	}

	return results
}

// Check if a specific position in the burrow is empty
func (burrow Burrow) empty(pos [2]int) bool {
	return burrow[loc(pos)] == '.'
}

// Generate all possible next steps for the current burrow state
func (burrow Burrow) move() []Step {
	nexts := []Step{}

	// Move amphipods from rooms to the hallway
	for _, room := range Rooms {
		end := 4 // Default depth for unfolded burrow
		if burrow.folded() {
			end = 2 // Adjust depth for folded burrow
		}
	inner:
		for _, n := range room[:end] {
			if !burrow.empty(n) {
				nexts = append(nexts, burrow.moveToHallway(n)...) // Move top amphipod in the room
				break inner
			}
		}
	}

	// Move amphipods from the hallway to their target rooms
	for _, h := range Hallway {
		if !burrow.empty(h) {
			nexts = append(nexts, burrow.moveToRoom(h)...) // Move amphipod to its room
		}
	}

	return nexts
}

func organize(input []string, target Burrow) int {
	start := Burrow(strings.Join(input, "\n"))

	memo = map[Burrow]int{
		start: 0,
	}

	queue := []Burrow{start}
	for len(queue) > 0 {
		burrow := queue[0]
		queue = queue[1:]

		for _, step := range burrow.move() {
			if _, ok := memo[step.burrow]; !ok {
				memo[step.burrow] = memo[burrow] + step.energy
				queue = append(queue, step.burrow)
			} else {
				if memo[step.burrow] > memo[burrow]+step.energy {
					memo[step.burrow] = memo[burrow] + step.energy
					queue = append(queue, step.burrow)
				}
			}
		}
	}

	return memo[target]
}

var targetBurrow = `#############
#...........#
###A#B#C#D###
  #A#B#C#D#  
  #########  `

var targetUnfoldedBurrow = `#############
#...........#
###A#B#C#D###
  #A#B#C#D#  
  #A#B#C#D#  
  #A#B#C#D#  
  #########  `

type Day23 struct{}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 23}
}

func (d Day23) Part1(input string) (string, error) {
	lines := strings.Split(input, "\n")

	validLines := []string{}
	for _, line := range lines {
		if line != "" {
			if len(line) < 13 {
				line += strings.Repeat(" ", 13-len(line))
			}

			validLines = append(validLines, line)
		}
	}
	lines = validLines
	result := organize(lines, Burrow(targetBurrow))

	return strconv.Itoa(result), nil
}

func (d Day23) Part2(input string) (string, error) {
	lines := strings.Split(input, "\n")

	validLines := []string{}
	for _, line := range lines {
		if line != "" {
			if len(line) < 13 {
				line += strings.Repeat(" ", 13-len(line))
			}

			validLines = append(validLines, line)
		}
	}
	lines = validLines

	lines = []string{lines[0], lines[1], lines[2], "  #D#C#B#A#  ", "  #D#B#A#C#  ", lines[3], lines[4]}
	result := organize(lines, Burrow(targetUnfoldedBurrow))

	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day23{})
}
