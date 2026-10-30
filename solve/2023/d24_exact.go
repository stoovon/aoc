package solve2023

import (
	"fmt"
	"math/big"
	"strings"
)

type hailExactStone struct{ position, velocity [3]*big.Int }

func hailParseIntegerVector(text string) ([3]*big.Int, error) {
	var result [3]*big.Int
	fields := strings.Split(text, ",")
	if len(fields) != 3 {
		return result, fmt.Errorf("expected three coordinates")
	}
	for i, field := range fields {
		field = strings.TrimSpace(field)
		var digits string
		if len(field) > 1 && (field[0] == '-' || field[0] == '+') {
			digits = field[1:]
		} else {
			digits = field
		}
		if digits == "" {
			return result, fmt.Errorf("expected a signed decimal integer coordinate")
		}
		for _, c := range digits {
			if c < '0' || c > '9' {
				return result, fmt.Errorf("expected a signed decimal integer coordinate")
			}
		}
		value, ok := new(big.Int).SetString(field, 10)
		if !ok {
			return result, fmt.Errorf("invalid integer coordinate")
		}
		result[i] = value
	}
	return result, nil
}
func hailParseExact(input string) ([]hailExactStone, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("puzzle input is empty")
	}
	input = strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	stones := make([]hailExactStone, 0, len(lines))
	for i, line := range lines {
		parse := func() (hailExactStone, error) {
			fields := strings.Split(line, "@")
			if len(fields) != 2 {
				return hailExactStone{}, fmt.Errorf("expected position @ velocity")
			}
			p, err := hailParseIntegerVector(fields[0])
			if err != nil {
				return hailExactStone{}, err
			}
			v, err := hailParseIntegerVector(fields[1])
			if err != nil {
				return hailExactStone{}, err
			}
			return hailExactStone{p, v}, nil
		}
		stone, err := parse()
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		stones = append(stones, stone)
	}
	return stones, nil
}
func hailAddPair(system *hailLinearSystem, a, b hailExactStone) error {
	// Subtracting two collision equations cancels the common rock-only product.
	for _, axes := range [][2]int{{0, 1}, {0, 2}, {1, 2}} {
		i, j := axes[0], axes[1]
		var coefficients [7]*big.Int
		for k := range coefficients {
			coefficients[k] = new(big.Int)
		}
		coefficients[i].Sub(a.velocity[j], b.velocity[j])
		coefficients[j].Sub(b.velocity[i], a.velocity[i])
		coefficients[3+i].Sub(b.position[j], a.position[j])
		coefficients[3+j].Sub(a.position[i], b.position[i])
		rhs := coefficients[6]
		rhs.Mul(a.position[i], a.velocity[j])
		rhs.Sub(rhs, new(big.Int).Mul(a.position[j], a.velocity[i]))
		rhs.Sub(rhs, new(big.Int).Mul(b.position[i], b.velocity[j]))
		rhs.Add(rhs, new(big.Int).Mul(b.position[j], b.velocity[i]))
		var row hailRow
		for k, c := range coefficients {
			row[k] = new(big.Rat).SetInt(c)
		}
		if err := system.addRow(row); err != nil {
			return err
		}
	}
	return nil
}
func hailCollisionTime(rock [6]*big.Rat, stone hailExactStone) (*big.Rat, error) {
	var time *big.Rat
	for axis := 0; axis < 3; axis++ {
		separation := new(big.Rat).Sub(new(big.Rat).SetInt(stone.position[axis]), rock[axis])
		relative := new(big.Rat).Sub(rock[axis+3], new(big.Rat).SetInt(stone.velocity[axis]))
		if relative.Sign() == 0 {
			if separation.Sign() != 0 {
				return nil, fmt.Errorf("rock does not collide with hailstone")
			}
		} else {
			axisTime := new(big.Rat).Quo(separation, relative)
			if axisTime.Sign() < 0 || (time != nil && time.Cmp(axisTime) != 0) {
				return nil, fmt.Errorf("rock must collide at one nonnegative time in all three coordinates")
			}
			time = axisTime
		}
	}
	if time == nil {
		time = new(big.Rat)
	}
	return time, nil
}
func hailSolveExact(stones []hailExactStone) ([6]*big.Rat, error) {
	var system hailLinearSystem
	for _, stone := range stones[1:] {
		if err := hailAddPair(&system, stones[0], stone); err != nil {
			return [6]*big.Rat{}, err
		}
		complete := true
		for _, pivot := range system.pivots {
			if pivot == nil {
				complete = false
				break
			}
		}
		if complete {
			break
		}
	}
	rock, err := system.solution()
	if err != nil {
		return rock, err
	}
	for i, stone := range stones {
		if _, err := hailCollisionTime(rock, stone); err != nil {
			return rock, fmt.Errorf("line %d: %w", i+1, err)
		}
	}
	return rock, nil
}
func hailExactPart2(input string) (string, error) {
	stones, err := hailParseExact(input)
	if err != nil {
		return "", err
	}
	rock, err := hailSolveExact(stones)
	if err != nil {
		return "", err
	}
	for _, value := range rock {
		if !value.IsInt() {
			return "", fmt.Errorf("rock position and velocity must be integral")
		}
	}
	sum := new(big.Int).Add(rock[0].Num(), rock[1].Num())
	sum.Add(sum, rock[2].Num())
	return sum.String(), nil
}
