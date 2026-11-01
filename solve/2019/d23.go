package solve2019

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day23 struct{}

func (d Day23) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2019, Day: 23} }

type networkNIC struct {
	vm          *Intcode
	input       []int64
	output      [3]int64
	outputCount int
}

// A turn ends on output, or on an empty input request after processing -1.
func (nic *networkNIC) turn() (output *int64, waiting, activity bool, err error) {
	emptyRead := false
	for {
		result := nic.vm.Step(nil)
		if result.Halted {
			return nil, false, activity, fmt.Errorf("network computer halted")
		}
		if result.Output != nil {
			return result.Output, false, true, nil
		}
		if !result.NeedInput {
			continue
		}
		value := int64(-1)
		if len(nic.input) > 0 {
			value = nic.input[0]
			nic.input = nic.input[1:]
			activity = true
			emptyRead = false
		} else {
			if emptyRead {
				return nil, true, activity, nil
			}
			emptyRead = true
		}
		nic.vm.Step(&value)
	}
}

func initNICs(input string, count int) []networkNIC {
	nics := make([]networkNIC, count)
	for i := range nics {
		nics[i] = networkNIC{vm: parseIntcode(input), input: []int64{int64(i)}}
	}
	return nics
}

func networkIdle(nics []networkNIC, allWaiting, activity bool) bool {
	if !allWaiting || activity {
		return false
	}
	for i := range nics {
		if len(nics[i].input) != 0 || nics[i].outputCount != 0 {
			return false
		}
	}
	return true
}

type networkPacket struct{ x, y int64 }

type packetNetwork struct {
	nics       []networkNIC
	next       int
	allWaiting bool
	activity   bool
}

func newPacketNetwork(nics []networkNIC) *packetNetwork {
	return &packetNetwork{nics: nics, allWaiting: true}
}

// Retain the sweep position when reporting a packet, so neither part skips work.
func (network *packetNetwork) nextEvent() (packet *networkPacket, idle bool, err error) {
	for {
		if network.next == len(network.nics) {
			idle := networkIdle(network.nics, network.allWaiting, network.activity)
			network.next, network.allWaiting, network.activity = 0, true, false
			if idle {
				return nil, true, nil
			}
		}
		i := network.next
		network.next++
		nic := &network.nics[i]
		output, waiting, active, err := nic.turn()
		if err != nil {
			return nil, false, fmt.Errorf("computer %d: %w", i, err)
		}
		network.allWaiting = network.allWaiting && waiting
		network.activity = network.activity || active
		if output == nil {
			continue
		}
		nic.output[nic.outputCount] = *output
		nic.outputCount++
		if nic.outputCount != 3 {
			continue
		}
		dest, x, y := nic.output[0], nic.output[1], nic.output[2]
		nic.outputCount = 0
		if dest == 255 {
			return &networkPacket{x, y}, false, nil
		}
		if dest < 0 || dest >= int64(len(network.nics)) {
			return nil, false, fmt.Errorf("computer %d: invalid packet destination", i)
		}
		network.nics[dest].input = append(network.nics[dest].input, x, y)
	}
}

func firstNetworkY(nics []networkNIC) (int64, error) {
	network := newPacketNetwork(nics)
	for {
		packet, _, err := network.nextEvent()
		if err != nil {
			return 0, err
		}
		if packet != nil {
			return packet.y, nil
		}
	}
}

func repeatedNATY(nics []networkNIC) (int64, error) {
	network := newPacketNetwork(nics)
	var nat *networkPacket
	var lastY int64
	lastYSet := false
	for {
		packet, idle, err := network.nextEvent()
		if err != nil {
			return 0, err
		}
		if packet != nil {
			nat = packet
		}
		if !idle || nat == nil {
			continue
		}
		network.nics[0].input = append(network.nics[0].input, nat.x, nat.y)
		if lastYSet && lastY == nat.y {
			return nat.y, nil
		}
		lastY, lastYSet = nat.y, true
	}
}

func parseNetwork(input string) ([]networkNIC, error) {
	for index, token := range strings.Split(strings.TrimSpace(input), ",") {
		if _, err := strconv.ParseInt(token, 10, 64); err != nil {
			return nil, fmt.Errorf("word %d: expected signed int64", index+1)
		}
	}
	return initNICs(input, 50), nil
}
func (d Day23) Part1(input string) (string, error) {
	nics, err := parseNetwork(input)
	if err != nil {
		return "", err
	}
	y, err := firstNetworkY(nics)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(y, 10), nil
}
func (d Day23) Part2(input string) (string, error) {
	nics, err := parseNetwork(input)
	if err != nil {
		return "", err
	}
	y, err := repeatedNATY(nics)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(y, 10), nil
}
func init() { solve.Register(Day23{}) }
