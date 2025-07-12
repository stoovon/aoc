package solve2018

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day7 struct {
}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 7}
}

var depRe = regexp.MustCompile(`Step (\w) must be finished before step (\w) can begin.`)

func parseDeps(input string) (allSteps map[byte]bool, prereqs map[byte]map[byte]bool) {
	deps := make(map[byte][]byte)
	allSteps = make(map[byte]bool)
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		m := depRe.FindStringSubmatch(line)
		if len(m) != 3 {
			continue
		}
		before, after := m[1][0], m[2][0]
		deps[after] = append(deps[after], before)
		allSteps[before] = true
		allSteps[after] = true
	}
	prereqs = make(map[byte]map[byte]bool)
	for step := range allSteps {
		prereqs[step] = make(map[byte]bool)
	}
	for step, pre := range deps {
		for _, p := range pre {
			prereqs[step][p] = true
		}
	}
	return
}

func (d Day7) Part1(input string) (string, error) {
	allSteps, prereqs := parseDeps(input)
	var order []byte
	done := make(map[byte]bool)
	for len(order) < len(allSteps) {
		var available []byte
		for step := range allSteps {
			if done[step] {
				continue
			}
			ready := true
			for p := range prereqs[step] {
				if !done[p] {
					ready = false
					break
				}
			}
			if ready {
				available = append(available, step)
			}
		}
		sort.Slice(available, func(i, j int) bool { return available[i] < available[j] })
		if len(available) == 0 {
			return "", errors.New("no available steps (cycle?)")
		}
		next := available[0]
		order = append(order, next)
		done[next] = true
	}
	return string(order), nil
}

func (d Day7) Part2(input string) (string, error) {
	const numWorkers = 5
	const baseTime = 60

	allSteps, prereqs := parseDeps(input)

	type worker struct {
		step byte
		rem  int
		busy bool
	}
	workers := make([]worker, numWorkers)
	done := make(map[byte]bool)
	inProgress := make(map[byte]bool)
	time := 0
	totalSteps := len(allSteps)

	for completed := 0; completed < totalSteps; {
		// Assign available steps to idle workers
		var available []byte
		for step := range allSteps {
			if done[step] || inProgress[step] {
				continue
			}
			ready := true
			for p := range prereqs[step] {
				if !done[p] {
					ready = false
					break
				}
			}
			if ready {
				available = append(available, step)
			}
		}
		sort.Slice(available, func(i, j int) bool { return available[i] < available[j] })

		// Assign steps to idle workers
		ai := 0
		for i := range workers {
			if !workers[i].busy && ai < len(available) {
				step := available[ai]
				workers[i].step = step
				workers[i].rem = baseTime + int(step-'A'+1)
				workers[i].busy = true
				inProgress[step] = true
				ai++
			}
		}

		// Find next completion time
		minRem := 0
		for i := range workers {
			if workers[i].busy {
				if minRem == 0 || workers[i].rem < minRem {
					minRem = workers[i].rem
				}
			}
		}
		if minRem == 0 {
			break // done
		}
		time += minRem

		// Advance time and finish steps
		for i := range workers {
			if workers[i].busy {
				workers[i].rem -= minRem
				if workers[i].rem == 0 {
					done[workers[i].step] = true
					delete(inProgress, workers[i].step)
					workers[i].busy = false
					completed++
				}
			}
		}
	}

	return strconv.Itoa(time), nil
}

func init() {
	solve.Register(Day7{})
}
