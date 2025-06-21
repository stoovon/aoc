package solve2017

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day25 struct {
}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 25}
}

type turingAction struct {
	write int
	move  int // -1 for left, +1 for right
	next  string
}

type turingState struct {
	actions [2]turingAction // actions[0] for value 0, actions[1] for value 1
}

func parseTuring(input string) (start string, steps int, states map[string]turingState, err error) {
	lines := strings.Split(input, "\n")
	reStart := regexp.MustCompile(`Begin in state ([A-Z])\.`)
	reSteps := regexp.MustCompile(`Perform a diagnostic checksum after (\d+) steps\.`)
	reState := regexp.MustCompile(`In state ([A-Z]):`)
	reIf := regexp.MustCompile(`If the current value is ([01]):`)
	reWrite := regexp.MustCompile(`- Write the value ([01])\.`)
	reMove := regexp.MustCompile(`- Move one slot to the (right|left)\.`)
	reNext := regexp.MustCompile(`- Continue with state ([A-Z])\.`)

	states = make(map[string]turingState)
	var curState string
	var curAction int
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if m := reStart.FindStringSubmatch(line); m != nil {
			start = m[1]
		} else if m := reSteps.FindStringSubmatch(line); m != nil {
			steps, _ = strconv.Atoi(m[1])
		} else if m := reState.FindStringSubmatch(line); m != nil {
			curState = m[1]
			states[curState] = turingState{}
		} else if m := reIf.FindStringSubmatch(line); m != nil {
			curAction, _ = strconv.Atoi(m[1])
		} else if m := reWrite.FindStringSubmatch(line); m != nil {
			state := states[curState]
			state.actions[curAction].write, _ = strconv.Atoi(m[1])
			states[curState] = state
		} else if m := reMove.FindStringSubmatch(line); m != nil {
			state := states[curState]
			if m[1] == "right" {
				state.actions[curAction].move = 1
			} else {
				state.actions[curAction].move = -1
			}
			states[curState] = state
		} else if m := reNext.FindStringSubmatch(line); m != nil {
			state := states[curState]
			state.actions[curAction].next = m[1]
			states[curState] = state
		}
	}
	if start == "" || steps == 0 || len(states) == 0 {
		return "", 0, nil, errors.New("parse error")
	}
	return start, steps, states, nil
}

func (d Day25) Part1(input string) (string, error) {
	start, steps, states, err := parseTuring(input)
	if err != nil {
		return "", err
	}
	tape := make(map[int]int)
	pos := 0
	state := start
	for i := 0; i < steps; i++ {
		val := tape[pos]
		act := states[state].actions[val]
		tape[pos] = act.write
		pos += act.move
		state = act.next
	}
	ones := 0
	for _, v := range tape {
		if v == 1 {
			ones++
		}
	}
	return strconv.Itoa(ones), nil
}

func (d Day25) Part2(_ string) (string, error) {
	return "", errors.New("not implemented")
}

func init() {
	solve.Register(Day25{})
}
