package solve2023

import (
	"fmt"
	"math/big"
)

const hailUnknowns = 6

type hailRow [hailUnknowns + 1]*big.Rat

type hailLinearSystem struct {
	pivots [hailUnknowns]*hailRow
}

func (s *hailLinearSystem) addRow(input hailRow) error {
	var row hailRow
	for i, value := range input {
		row[i] = new(big.Rat).Set(value)
	}
	for column := 0; column < hailUnknowns; column++ {
		if row[column].Sign() == 0 {
			continue
		}
		if pivot := s.pivots[column]; pivot != nil {
			factor := new(big.Rat).Set(row[column])
			for j := column; j <= hailUnknowns; j++ {
				product := new(big.Rat).Mul(factor, (*pivot)[j])
				row[j].Sub(row[j], product)
			}
		} else {
			divisor := new(big.Rat).Set(row[column])
			for j := column; j <= hailUnknowns; j++ {
				row[j].Quo(row[j], divisor)
			}
			s.pivots[column] = &row
			return nil
		}
	}
	if row[hailUnknowns].Sign() != 0 {
		return fmt.Errorf("hailstone equations are inconsistent")
	}
	return nil
}

func (s *hailLinearSystem) solution() ([hailUnknowns]*big.Rat, error) {
	var result [hailUnknowns]*big.Rat
	for column := hailUnknowns - 1; column >= 0; column-- {
		row := s.pivots[column]
		if row == nil {
			return result, fmt.Errorf("hailstone equations do not determine a unique rock")
		}
		value := new(big.Rat).Set((*row)[hailUnknowns])
		for j := column + 1; j < hailUnknowns; j++ {
			value.Sub(value, new(big.Rat).Mul((*row)[j], result[j]))
		}
		result[column] = value
	}
	return result, nil
}
