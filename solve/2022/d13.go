package solve2022

import (
	"aoc/solve"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

type Day13 struct{}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 13}
}

func (d Day13) Part1(input string) (string, error) {
	// Parse the input into packets
	packets, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	// Sum indices of packets in the right order
	sum := 0
	for i, pair := range packets {
		if rightOrder(pair[0], pair[1]) {
			sum += i + 1 // Indices are 1-based
		}
	}

	return strconv.Itoa(sum), nil
}

func (d Day13) Part2(input string) (string, error) {
	// Parse the input into packets
	packets, err := d.parseInput(input)
	if err != nil {
		return "", err
	}

	// Flatten the packet pairs into a single list
	var allPackets []any
	for _, pair := range packets {
		allPackets = append(allPackets, pair...)
	}

	// Add the divider packets
	dividers := []any{[]any{float64(2)}, []any{float64(6)}}
	allPackets = append(allPackets, dividers...)

	// Sort the packets using the compare function
	sortedPackets := allPackets
	sort.Slice(sortedPackets, func(i, j int) bool {
		return compare(sortedPackets[i], sortedPackets[j]) < 0
	})

	// Find the indices of the divider packets
	decoderKey := 1
	for i, packet := range sortedPackets {
		for _, divider := range dividers {
			if compare(packet, divider) == 0 {
				decoderKey *= (i + 1) // Indices are 1-based
			}
		}
	}

	return strconv.Itoa(decoderKey), nil
}

func init() {
	solve.Register(Day13{})
}

func (d Day13) parseInput(input string) ([][]any, error) {
	var packets [][]any
	paragraphs := strings.Split(strings.TrimSpace(input), "\n\n")
	for _, par := range paragraphs {
		var pair []any
		lines := strings.Split(par, "\n")
		for _, line := range lines {
			var packet any
			if err := json.Unmarshal([]byte(line), &packet); err != nil {
				return nil, err
			}
			pair = append(pair, packet)
		}
		packets = append(packets, pair)
	}
	return packets, nil
}

func rightOrder(left, right any) bool {
	return compare(left, right) <= 0
}

func compare(left, right any) int {
	switch l := left.(type) {
	case float64: // JSON numbers are unmarshalled as float64
		switch r := right.(type) {
		case float64:
			return int(l - r)
		case []any:
			return compare([]any{l}, r)
		}
	case []any:
		switch r := right.(type) {
		case float64:
			return compare(l, []any{r})
		case []any:
			for i := 0; i < len(l) && i < len(r); i++ {
				if cmp := compare(l[i], r[i]); cmp != 0 {
					return cmp
				}
			}
			return len(l) - len(r)
		}
	}
	return 0
}
