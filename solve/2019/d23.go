package solve2019

import (
	"strconv"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day23 struct {
}

func (d Day23) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 23}
}

type packet struct{ grids.WidePoint }

func launchNIC(prog *Intcode, in chan int64, out chan int64) {
	go prog.ChannelRunner(in, out)
}

func initNICs(input string, nics int) ([]chan int64, []chan int64) {
	ins := make([]chan int64, nics)
	outs := make([]chan int64, nics)
	for i := 0; i < nics; i++ {
		ins[i] = make(chan int64, 1)
		outs[i] = make(chan int64, 1)
	}
	for i := 0; i < nics; i++ {
		prog := parseIntcode(input)
		launchNIC(prog, ins[i], outs[i])
	}
	for i := 0; i < nics; i++ {
		ins[i] <- int64(i)
	}
	return ins, outs
}

func routePackets(ins, outs []chan int64, natMode bool) (int64, error) {
	nics := len(ins)
	queues := make([][]int64, nics)
	outbufs := make([][]int64, nics)
	var natX, natY int64
	var lastNatY int64
	natSet := false
	lastNatYSet := false

	idleRounds := 0

	for {
		activity := false

		// Drain all outputs for all NICs, non-blocking, buffer until 3 values
		for i := 0; i < nics; i++ {
			select {
			case v, ok := <-outs[i]:
				if !ok {
					continue
				}
				outbufs[i] = append(outbufs[i], v)
				activity = true
				for len(outbufs[i]) >= 3 {
					dest, x, y := outbufs[i][0], outbufs[i][1], outbufs[i][2]
					outbufs[i] = outbufs[i][3:]
					if dest == 255 {
						if !natMode {
							return y, nil // Part 1: return immediately
						}
						natX, natY = x, y
						natSet = true
					} else if dest >= 0 && dest < int64(nics) {
						queues[dest] = append(queues[dest], x, y)
					}
				}
			default:
				// no output ready
			}
		}

		// Feed input to all NICs (buffered channels, so this won't block)
		idleCount := 0
		for i := 0; i < nics; i++ {
			if len(queues[i]) > 0 {
				select {
				case ins[i] <- queues[i][0]:
					queues[i] = queues[i][1:]
					activity = true
				default:
				}
			} else {
				select {
				case ins[i] <- -1:
					idleCount++
				default:
				}
			}
		}

		// Check if all queues are empty and all NICs are idle
		allIdle := idleCount == nics
		queuesEmpty := true
		for i := 0; i < nics; i++ {
			if len(queues[i]) > 0 {
				queuesEmpty = false
				break
			}
		}

		// Only increment idleRounds if truly idle
		if natMode && allIdle && queuesEmpty && !activity && natSet {
			idleRounds++
		} else {
			idleRounds = 0
		}

		// For part 2, check for idle and NAT injection
		if natMode && idleRounds > 1 {
			if lastNatYSet && lastNatY == natY {
				return natY, nil
			}
			queues[0] = append(queues[0], natX, natY)
			lastNatY = natY
			lastNatYSet = true
			idleRounds = 0
		}
	}
}

func (d Day23) Part1(input string) (string, error) {
	const nics = 50
	ins, outs := initNICs(input, nics)
	y, err := routePackets(ins, outs, false)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(y, 10), nil
}

func (d Day23) Part2(input string) (string, error) {
	const nics = 50
	ins, outs := initNICs(input, nics)
	y, err := routePackets(ins, outs, true)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(y, 10), nil
}

func init() {
	solve.Register(Day23{})
}
