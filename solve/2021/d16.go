package solve2021

import (
	"aoc/solve"
	"fmt"
	"strconv"
	"strings"
)

type Day16 struct{}

func (d Day16) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 16}
}

func (d Day16) Part1(input string) (string, error) {
	bin := hexToBin(strings.TrimSpace(input))
	_, versionSum, _ := parsePacketFull(bin, 0)
	return strconv.Itoa(versionSum), nil
}

func hexToBin(hex string) string {
	var sb strings.Builder
	for _, c := range hex {
		n, _ := strconv.ParseUint(string(c), 16, 4)
		sb.WriteString(fmt.Sprintf("%04b", n))
	}
	return sb.String()
}

// Returns: next index after packet, sum of version numbers, value of packet
func parsePacketFull(bin string, i int) (int, int, int) {
	if i+6 > len(bin) {
		return len(bin), 0, 0
	}
	version, _ := strconv.ParseInt(bin[i:i+3], 2, 64)
	typeID, _ := strconv.ParseInt(bin[i+3:i+6], 2, 64)
	i += 6
	versionSum := int(version)
	if typeID == 4 {
		// Literal value
		val := 0
		for {
			if i+5 > len(bin) {
				break
			}
			group := bin[i : i+5]
			i += 5
			n, _ := strconv.ParseInt(group[1:], 2, 64)
			val = (val << 4) | int(n)
			if group[0] == '0' {
				break
			}
		}
		return i, versionSum, val
	}
	var values []int
	if i >= len(bin) {
		return i, versionSum, 0
	}
	lengthTypeID := bin[i]
	i++
	if lengthTypeID == '0' {
		// Next 15 bits: total length in bits of sub-packets
		if i+15 > len(bin) {
			return i, versionSum, 0
		}
		totalLen, _ := strconv.ParseInt(bin[i:i+15], 2, 64)
		i += 15
		end := i + int(totalLen)
		for i < end {
			ni, vs, val := parsePacketFull(bin, i)
			versionSum += vs
			values = append(values, val)
			i = ni
		}
	} else {
		// Next 11 bits: number of sub-packets
		if i+11 > len(bin) {
			return i, versionSum, 0
		}
		numPackets, _ := strconv.ParseInt(bin[i:i+11], 2, 64)
		i += 11
		for j := 0; j < int(numPackets); j++ {
			ni, vs, val := parsePacketFull(bin, i)
			versionSum += vs
			values = append(values, val)
			i = ni
		}
	}
	var result int
	switch typeID {
	case 0:
		result = 0
		for _, v := range values {
			result += v
		}
	case 1:
		result = 1
		for _, v := range values {
			result *= v
		}
	case 2:
		result = values[0]
		for _, v := range values {
			if v < result {
				result = v
			}
		}
	case 3:
		result = values[0]
		for _, v := range values {
			if v > result {
				result = v
			}
		}
	case 5:
		if len(values) == 2 && values[0] > values[1] {
			result = 1
		} else {
			result = 0
		}
	case 6:
		if len(values) == 2 && values[0] < values[1] {
			result = 1
		} else {
			result = 0
		}
	case 7:
		if len(values) == 2 && values[0] == values[1] {
			result = 1
		} else {
			result = 0
		}
	}
	return i, versionSum, result
}

func (d Day16) Part2(input string) (string, error) {
	bin := hexToBin(strings.TrimSpace(input))
	_, _, value := parsePacketFull(bin, 0)
	return strconv.Itoa(value), nil
}

func init() {
	solve.Register(Day16{})
}
