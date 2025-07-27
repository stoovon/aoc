package solve2020

import (
	"aoc/solve"
	"errors"
	"math"
	"strconv"
	"strings"
)

type Day13 struct{}

func (d Day13) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 13}
}

func (d Day13) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) != 2 {
		return "", errors.New("invalid input format")
	}

	// Parse the earliest departure time
	earliestTime, err := strconv.Atoi(lines[0])
	if err != nil {
		return "", err
	}

	// Parse bus IDs (ignore 'x' entries)
	busIDsStr := strings.Split(lines[1], ",")
	var busIDs []int

	for _, busStr := range busIDsStr {
		if busStr != "x" {
			busID, err := strconv.Atoi(busStr)
			if err != nil {
				return "", err
			}
			busIDs = append(busIDs, busID)
		}
	}

	// Find the earliest bus we can take
	minWaitTime := math.MaxInt32
	bestBusID := 0

	for _, busID := range busIDs {
		// Calculate when this bus will next depart after earliestTime
		// If earliestTime is exactly divisible by busID, we can take it immediately
		waitTime := busID - (earliestTime % busID)
		if waitTime == busID {
			waitTime = 0
		}

		if waitTime < minWaitTime {
			minWaitTime = waitTime
			bestBusID = busID
		}
	}

	result := bestBusID * minWaitTime
	return strconv.Itoa(result), nil
}

// Extended Euclidean Algorithm to find modular inverse
func extendedGCD(a, b int) (int, int, int) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, x1, y1 := extendedGCD(b%a, a)
	x := y1 - (b/a)*x1
	y := x1
	return gcd, x, y
}

func modInverse(a, m int) int {
	gcd, x, _ := extendedGCD(a, m)
	if gcd != 1 {
		return -1 // No modular inverse exists
	}
	return (x%m + m) % m
}

// Chinese Remainder Theorem solver
func solveCRT(remainders, moduli []int) int {
	if len(remainders) != len(moduli) {
		return -1
	}

	// Calculate the product of all moduli
	totalProduct := 1
	for _, mod := range moduli {
		totalProduct *= mod
	}

	result := 0
	for i := 0; i < len(remainders); i++ {
		// Calculate Ni = totalProduct / moduli[i]
		ni := totalProduct / moduli[i]

		// Calculate Mi = modular inverse of Ni modulo moduli[i]
		mi := modInverse(ni, moduli[i])
		if mi == -1 {
			return -1 // No solution exists
		}

		// Add to result: remainders[i] * Ni * Mi
		result += remainders[i] * ni * mi
	}

	return ((result % totalProduct) + totalProduct) % totalProduct
}

func (d Day13) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) != 2 {
		return "", errors.New("invalid input format")
	}

	busIDsStr := strings.Split(lines[1], ",")

	var remainders []int
	var moduli []int

	for i, busStr := range busIDsStr {
		if busStr != "x" {
			busID, err := strconv.Atoi(busStr)
			if err != nil {
				return "", err
			}

			// We want: t + i ≡ 0 (mod busID)
			// Which means: t ≡ -i (mod busID)
			// Or equivalently: t ≡ (busID - i) (mod busID)
			remainder := (busID - i%busID) % busID

			remainders = append(remainders, remainder)
			moduli = append(moduli, busID)
		}
	}

	result := solveCRT(remainders, moduli)
	if result == -1 {
		return "", errors.New("no solution found")
	}

	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day13{})
}
