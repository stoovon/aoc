package solve2016

import (
	"strings"
	"testing"
)

func TestElevatorKeyPreservesPairRelationships(t *testing.T) {
	first := Building{Elevator: 1, Floors: [4]Floor{1, 1 << 8, 2 | (2 << 8), 0}}
	renamed := Building{Elevator: 1, Floors: [4]Floor{2, 2 << 8, 1 | (1 << 8), 0}}
	if first.Norm() != renamed.Norm() {
		t.Fatal("consistent element renaming changed key")
	}
	crossed := Building{Elevator: 1, Floors: [4]Floor{1, 2 << 8, 2 | (1 << 8), 0}}
	if first.Norm() == crossed.Norm() {
		t.Fatal("different pair relationships collapsed")
	}
	moved := first
	moved.Elevator = 2
	if first.Norm() == moved.Norm() {
		t.Fatal("elevator position missing from key")
	}
	empty := Building{}
	pairOnBottom := Building{Floors: [4]Floor{1 | (1 << 8), 0, 0, 0}}
	if empty.Norm() == pairOnBottom.Norm() {
		t.Fatal("pair count missing from key")
	}
}

type labelledElevatorState struct {
	floors   [4]Floor
	elevator int
}

func TestElevatorSymmetryAgainstLabelledSearch(t *testing.T) {
	// A reverse BFS over labelled states is an oracle with no symmetry reduction.
	goal := Building{Floors: [4]Floor{0, 0, 0, 3 | (3 << 8)}}
	distances := make(map[labelledElevatorState]int)
	queue := []Building{}
	for elevator := 0; elevator < 4; elevator++ {
		state := goal
		state.Elevator = elevator
		distances[labelledElevatorState{state.Floors, elevator}] = 0
		queue = append(queue, state)
	}
	for head := 0; head < len(queue); head++ {
		state := queue[head]
		distance := distances[labelledElevatorState{state.Floors, state.Elevator}]
		for _, next := range elevatorReferenceMoves(state) {
			key := labelledElevatorState{next.Floors, next.Elevator}
			if _, seen := distances[key]; seen {
				continue
			}
			distances[key] = distance + 1
			queue = append(queue, next)
		}
	}
	day := Day11{}
	for key, want := range distances {
		start := Building{Floors: key.floors, Elevator: key.elevator}
		got, err := day.bfs(start, goal)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("state %v elevator %d: %d want %d", key.floors, key.elevator, got, want)
		}
	}
	if len(distances) < 100 {
		t.Fatalf("oracle covered only %d states", len(distances))
	}
}

func TestElevatorLayersAndNoRoute(t *testing.T) {
	day := Day11{}
	goal := Building{Floors: [4]Floor{0, 0, 0, 1 | (1 << 8)}}
	if got, err := day.bfs(goal, goal); err != nil || got != 0 {
		t.Fatalf("already finished: %d %v", got, err)
	}
	start := Building{Floors: [4]Floor{1 | (1 << 8), 0, 0, 0}}
	if got, err := day.bfs(start, goal); err != nil || got != 3 {
		t.Fatalf("three elevator rides: %d %v", got, err)
	}
	stranded := Building{Floors: [4]Floor{0, 1 | (1 << 8), 0, 0}}
	if _, err := day.bfs(stranded, goal); err == nil {
		t.Fatal("unreachable goal reported as zero moves")
	}
}

func elevatorReferenceSafe(building Building) bool {
	for _, floor := range building.Floors {
		hasGenerator := false
		for element := 0; element < 8; element++ {
			hasGenerator = hasGenerator || floor&(1<<(element+8)) != 0
		}
		if !hasGenerator {
			continue
		}
		for element := 0; element < 8; element++ {
			if floor&(1<<element) != 0 && floor&(1<<(element+8)) == 0 {
				return false
			}
		}
	}
	return true
}

func elevatorReferenceMoves(state Building) []Building {
	var result []Building
	present := uint16(state.Floors[state.Elevator])
	// Enumerate every subset independently, retaining exactly one or two items.
	for subset := present; subset != 0; subset = (subset - 1) & present {
		count := 0
		for bit := 0; bit < 16; bit++ {
			if subset&(1<<bit) != 0 {
				count++
			}
		}
		if count != 1 && count != 2 {
			continue
		}
		for target := 0; target < 4; target++ {
			if target != state.Elevator-1 && target != state.Elevator+1 {
				continue
			}
			next := state
			next.Floors[state.Elevator] &^= Floor(subset)
			next.Floors[target] |= Floor(subset)
			next.Elevator = target
			if elevatorReferenceSafe(next) {
				result = append(result, next)
			}
		}
	}
	return result
}

func TestElevatorMovesAgainstAllSubsets(t *testing.T) {
	for placement := 0; placement < 256; placement++ {
		var building Building
		rest := placement
		for _, component := range []Component{0, 8, 1, 9} {
			building.Floors[rest%4].Add(component)
			rest /= 4
		}
		if !elevatorReferenceSafe(building) {
			continue
		}
		for elevator := 0; elevator < 4; elevator++ {
			building.Elevator = elevator
			want := make(map[Building]bool)
			for _, next := range elevatorReferenceMoves(building) {
				want[next] = true
			}
			got := building.AllCombinations()
			if len(got) != len(want) {
				t.Fatalf("state %v elevator %d: %d moves want %d", building.Floors, elevator, len(got), len(want))
			}
			for _, next := range got {
				if !want[next] {
					t.Fatalf("unexpected or duplicated move: %v", next)
				}
				delete(want, next)
			}
		}
	}
}

func TestElevatorRejectsUnsafeStart(t *testing.T) {
	unsafe := Building{Floors: [4]Floor{1 | (2 << 8), 2 | (1 << 8), 0, 0}}
	if _, err := (Day11{}).bfs(unsafe, Building{}); err == nil {
		t.Fatal("accepted unsafe initial building")
	}
	for _, floor := range []int{-1, 4} {
		if _, err := (Day11{}).bfs(Building{Elevator: floor}, Building{}); err == nil {
			t.Fatalf("accepted elevator floor %d", floor)
		}
	}
}

func elevatorSyntheticDescription(contents [4]string) string {
	var lines []string
	for floor, ordinal := range []string{"first", "second", "third", "fourth"} {
		body := contents[floor]
		if body == "" {
			body = "nothing relevant"
		}
		lines = append(lines, "The "+ordinal+" floor contains "+body+".")
	}
	return strings.Join(lines, "\n")
}

func TestElevatorParsesArbitraryElementNames(t *testing.T) {
	day := Day11{}
	input := elevatorSyntheticDescription([4]string{"an amber generator and an amber-compatible microchip", "", "", ""})
	for _, text := range []string{input, input + "\n", strings.ReplaceAll(input, "\n", "\r\n") + "\r\n"} {
		if got, err := day.Part1(text); err != nil || got != "3" {
			t.Fatalf("arbitrary name: %q %v", got, err)
		}
	}
	// Names formerly reserved for part 2 must still get two fresh pair slots.
	input = elevatorSyntheticDescription([4]string{"", "", "", "a dilithium generator and a dilithium-compatible microchip"})
	if got, err := day.Part2(input); err != nil || got != "15" {
		t.Fatalf("fresh additions: %q %v", got, err)
	}
}

func TestElevatorParserRejectsBadDescriptions(t *testing.T) {
	valid := elevatorSyntheticDescription([4]string{"an amber generator and an amber-compatible microchip", "", "", ""})
	for _, input := range []string{
		"", strings.Join(strings.Split(valid, "\n")[:3], "\n"), valid + "\nextra", strings.Replace(valid, "first", "second", 1),
		strings.Replace(valid, "an amber generator", "an amber gizmo", 1),
		strings.Replace(valid, "an amber generator", "an amber generator and an amber generator", 1),
		strings.Replace(valid, "an amber-compatible microchip", "an brass-compatible microchip", 1),
		strings.Replace(valid, "an amber generator", "an 123 generator", 1),
		strings.Replace(valid, "an amber generator", "amber generator", 1),
		elevatorSyntheticDescription([4]string{"an amber generator", "an amber generator and an amber-compatible microchip", "", ""}),
	} {
		if _, err := (Day11{}).parseInput(input); err == nil {
			t.Fatalf("accepted malformed description %q", input)
		}
	}
	unsafe := elevatorSyntheticDescription([4]string{"an amber generator and a brass-compatible microchip", "a brass generator and an amber-compatible microchip", "", ""})
	if _, err := (Day11{}).Part1(unsafe); err == nil {
		t.Fatal("accepted unsafe parsed building")
	}
	exposed := elevatorSyntheticDescription([4]string{"an amber-compatible microchip", "an amber generator", "", ""})
	if _, err := (Day11{}).Part2(exposed); err == nil {
		t.Fatal("accepted building made unsafe by new generators")
	}
}

func TestElevatorParserCapacity(t *testing.T) {
	names := []string{"amber", "brass", "cinder", "dusk", "ember", "frost", "glint", "haze", "indigo"}
	for _, count := range []int{7, 8, 9} {
		var items []string
		for _, name := range names[:count] {
			items = append(items, "a "+name+" generator", "a "+name+"-compatible microchip")
		}
		input := elevatorSyntheticDescription([4]string{"", "", "", strings.Join(items, ", ")})
		_, err := (Day11{}).parseInput(input)
		if (err != nil) != (count > 8) {
			t.Fatalf("%d pairs parse: %v", count, err)
		}
		if count <= 8 {
			if got, err := (Day11{}).Part1(input); err != nil || got != "0" {
				t.Fatalf("completed %d pairs: %q %v", count, got, err)
			}
		}
		if _, err := (Day11{}).Part2(input); err == nil {
			t.Fatalf("accepted %d+2 pairs", count)
		}
	}
}
