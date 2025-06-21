package solve2017

import (
	"errors"
	"strconv"
	"strings"
	"sync/atomic"

	"aoc/solve"
)

type Day18 struct {
}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 18}
}

func getVal(reg map[string]int, x string) int {
	if v, err := strconv.Atoi(x); err == nil {
		return v
	}
	return reg[x]
}

type progResult struct {
	sendCount int
}

type instrCallbacks struct {
	snd func(val int)
	rcv func(x string) (recovered bool)
}

var errNoRecovery = errors.New("no recovery")

func execInstructions(
	lines []string,
	reg map[string]int,
	cb instrCallbacks,
) (lastSound int, ip int, err error) {
	for ip = 0; ip >= 0 && ip < len(lines); {
		parts := strings.Fields(lines[ip])
		switch parts[0] {
		case "snd":
			val := getVal(reg, parts[1])
			cb.snd(val)
			lastSound = val
			ip++
		case "set":
			reg[parts[1]] = getVal(reg, parts[2])
			ip++
		case "add":
			reg[parts[1]] += getVal(reg, parts[2])
			ip++
		case "mul":
			reg[parts[1]] *= getVal(reg, parts[2])
			ip++
		case "mod":
			reg[parts[1]] %= getVal(reg, parts[2])
			ip++
		case "rcv":
			if cb.rcv(parts[1]) {
				return lastSound, ip, nil
			}
			ip++
		case "jgz":
			if getVal(reg, parts[1]) > 0 {
				ip += getVal(reg, parts[2])
			} else {
				ip++
			}
		}
	}
	return lastSound, ip, errNoRecovery
}

func (d Day18) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	reg := map[string]int{}

	cb := instrCallbacks{
		snd: func(val int) {},
		rcv: func(x string) bool { return getVal(reg, x) != 0 },
	}

	sound, _, err := execInstructions(lines, reg, cb)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(sound), nil
}

func (d Day18) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	ch0 := make(chan int, 1000)
	ch1 := make(chan int, 1000)
	done0 := make(chan progResult, 1)
	done1 := make(chan progResult, 1)

	var waiting0, waiting1 int32

	runProg := func(id int, in <-chan int, out chan<- int, done chan<- progResult, waiting *int32) {
		reg := map[string]int{"p": id}
		sendCount := 0
		cb := instrCallbacks{
			snd: func(val int) {
				out <- val
				sendCount++
			},
			rcv: func(x string) bool {
				atomic.StoreInt32(waiting, 1)
				val, ok := <-in
				if !ok {
					done <- progResult{sendCount}
					return true // signal to execInstructions to exit
				}
				atomic.StoreInt32(waiting, 0)
				reg[x] = val
				return false // continue execution
			},
		}
		_, _, err := execInstructions(lines, reg, cb)
		if err != nil && !errors.Is(err, errNoRecovery) {
			panic("unexpected error: " + err.Error())
		}
		done <- progResult{sendCount}
	}

	go runProg(0, ch0, ch1, done0, &waiting0)
	go runProg(1, ch1, ch0, done1, &waiting1)

	for {
		// Deadlock detection: both waiting and both channels empty
		if atomic.LoadInt32(&waiting0) == 1 && atomic.LoadInt32(&waiting1) == 1 && len(ch0) == 0 && len(ch1) == 0 {
			close(ch0)
			close(ch1)
			return strconv.Itoa((<-done1).sendCount), nil
		}
		select {
		case res := <-done1:
			return strconv.Itoa(res.sendCount), nil
		case <-done0:
			// ignore, only care about program 1's send count
		default:
			// allow the loop to check deadlock condition
		}
	}
}

func init() {
	solve.Register(Day18{})
}
