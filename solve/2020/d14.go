package solve2020

import (
	"aoc/solve"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Day14 struct{}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 14}
}

func applyMask(value int64, mask string) int64 {
	result := value

	for i, bit := range mask {
		bitPosition := 35 - i // Convert from left-to-right to bit position

		switch bit {
		case '0':
			result &= ^(1 << bitPosition)
		case '1':
			result |= (1 << bitPosition)
		case 'X':
			continue
		}
	}

	return result
}

func (d Day14) parseInput(input string) ([]string, *regexp.Regexp, *regexp.Regexp) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	maskPattern := regexp.MustCompile(`^mask = ([X01]{36})$`)
	memPattern := regexp.MustCompile(`^mem\[(\d+)\] = (\d+)$`)

	return lines, maskPattern, memPattern
}

func sumMemory(memory map[int64]int64) string {
	var sum int64
	for _, value := range memory {
		sum += value
	}
	return strconv.FormatInt(sum, 10)
}

func (d Day14) Part1(input string) (string, error) {
	lines, maskPattern, memPattern := d.parseInput(input)

	memory := make(map[int64]int64)
	var currentMask string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if maskMatch := maskPattern.FindStringSubmatch(line); maskMatch != nil {
			currentMask = maskMatch[1]
		} else if memMatch := memPattern.FindStringSubmatch(line); memMatch != nil {
			address, err := strconv.ParseInt(memMatch[1], 10, 64)
			if err != nil {
				return "", err
			}

			value, err := strconv.ParseInt(memMatch[2], 10, 64)
			if err != nil {
				return "", err
			}

			maskedValue := applyMask(value, currentMask)
			memory[address] = maskedValue
		} else {
			return "", errors.New("invalid instruction: " + line)
		}
	}

	return sumMemory(memory), nil
}

func applyAddressMask(address int64, mask string) []int64 {
	// Apply the mask to create a template with floating bits
	var template []rune
	for i, bit := range mask {
		bitPosition := 35 - i // Convert from left-to-right to bit position

		switch bit {
		case '0':
			if (address & (1 << bitPosition)) != 0 {
				template = append(template, '1')
			} else {
				template = append(template, '0')
			}
		case '1':
			template = append(template, '1')
		case 'X':
			template = append(template, 'X')
		}
	}

	return expandFloatingBits(template)
}

func expandFloatingBits(template []rune) []int64 {
	var floatingPositions []int
	for i, bit := range template {
		if bit == 'X' {
			floatingPositions = append(floatingPositions, i)
		}
	}

	numCombinations := 1 << len(floatingPositions)
	addresses := make([]int64, numCombinations)

	for combo := 0; combo < numCombinations; combo++ {
		current := make([]rune, len(template))
		copy(current, template)

		for i, pos := range floatingPositions {
			if (combo & (1 << i)) != 0 {
				current[pos] = '1'
			} else {
				current[pos] = '0'
			}
		}

		var address int64
		for i, bit := range current {
			if bit == '1' {
				bitPosition := 35 - i
				address |= (1 << bitPosition)
			}
		}

		addresses[combo] = address
	}

	return addresses
}

func (d Day14) Part2(input string) (string, error) {
	lines, maskPattern, memPattern := d.parseInput(input)

	memory := make(map[int64]int64)
	var currentMask string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if maskMatch := maskPattern.FindStringSubmatch(line); maskMatch != nil {
			currentMask = maskMatch[1]
		} else if memMatch := memPattern.FindStringSubmatch(line); memMatch != nil {
			address, err := strconv.ParseInt(memMatch[1], 10, 64)
			if err != nil {
				return "", err
			}

			value, err := strconv.ParseInt(memMatch[2], 10, 64)
			if err != nil {
				return "", err
			}

			addresses := applyAddressMask(address, currentMask)
			for _, addr := range addresses {
				memory[addr] = value
			}
		} else {
			return "", errors.New("invalid instruction: " + line)
		}
	}

	return sumMemory(memory), nil
}

func init() {
	solve.Register(Day14{})
}
