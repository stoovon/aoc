package solve2023

import (
	"fmt"
	"math"
	"math/big"
)

func hailCrossXY(a, b [3]*big.Int) *big.Int {
	return new(big.Int).Sub(new(big.Int).Mul(a[0], b[1]), new(big.Int).Mul(a[1], b[0]))
}
func hailStationaryCrossing(point, moving hailExactStone, min, max *big.Int) bool {
	for axis := 0; axis < 2; axis++ {
		if point.position[axis].Cmp(min) < 0 || point.position[axis].Cmp(max) > 0 {
			return false
		}
	}
	var delta [3]*big.Int
	for axis := range delta {
		delta[axis] = new(big.Int).Sub(point.position[axis], moving.position[axis])
	}
	if hailCrossXY(delta, moving.velocity).Sign() != 0 {
		return false
	}
	for axis := 0; axis < 2; axis++ {
		if moving.velocity[axis].Sign() != 0 {
			return delta[axis].Sign() == 0 || delta[axis].Sign() == moving.velocity[axis].Sign()
		}
	}
	return delta[0].Sign() == 0 && delta[1].Sign() == 0
}

// Scratch values belong to one solve and never alias the parsed hailstones.
type hailIntersectionScratch struct {
	delta                             [3]big.Int
	determinant, aTime, bTime         big.Int
	lower, upper, coordinate, product big.Int
}

func (scratch *hailIntersectionScratch) cross(out *big.Int, a, b [3]*big.Int) *big.Int {
	out.Mul(a[0], b[1])
	scratch.product.Mul(a[1], b[0])
	return out.Sub(out, &scratch.product)
}

func (scratch *hailIntersectionScratch) crosses(a, b hailExactStone, min, max *big.Int) (bool, error) {
	if a.velocity[0].Sign() == 0 && a.velocity[1].Sign() == 0 {
		return hailStationaryCrossing(a, b, min, max), nil
	}
	if b.velocity[0].Sign() == 0 && b.velocity[1].Sign() == 0 {
		return hailStationaryCrossing(b, a, min, max), nil
	}
	var delta [3]*big.Int
	for axis := range delta {
		delta[axis] = &scratch.delta[axis]
		delta[axis].Sub(b.position[axis], a.position[axis])
	}
	determinant := scratch.cross(&scratch.determinant, a.velocity, b.velocity)
	if determinant.Sign() == 0 {
		if scratch.cross(&scratch.aTime, delta, a.velocity).Sign() == 0 {
			return false, fmt.Errorf("collinear projected paths do not have a unique line intersection")
		}
		return false, nil
	}
	aTime := scratch.cross(&scratch.aTime, delta, b.velocity)
	bTime := scratch.cross(&scratch.bTime, delta, a.velocity)
	if determinant.Sign() < 0 {
		determinant.Neg(determinant)
		aTime.Neg(aTime)
		bTime.Neg(bTime)
	}
	if aTime.Sign() < 0 || bTime.Sign() < 0 {
		return false, nil
	}
	// Compare scaled integer coordinates instead of allocating fractions per pair.
	lower := scratch.lower.Mul(min, determinant)
	upper := scratch.upper.Mul(max, determinant)
	for axis := 0; axis < 2; axis++ {
		coordinate := scratch.coordinate.Mul(a.position[axis], determinant)
		scratch.product.Mul(a.velocity[axis], aTime)
		coordinate.Add(coordinate, &scratch.product)
		if coordinate.Cmp(lower) < 0 || coordinate.Cmp(upper) > 0 {
			return false, nil
		}
	}
	return true, nil
}
func hailCountCrossings(input string, min, max *big.Int) (int, error) {
	if min.Cmp(max) > 0 {
		return 0, fmt.Errorf("intersection window minimum exceeds maximum")
	}
	stones, err := hailParseExact(input)
	if err != nil {
		return 0, err
	}
	count := 0
	var scratch hailIntersectionScratch
	for i, a := range stones {
		for _, b := range stones[i+1:] {
			crosses, err := scratch.crosses(a, b, min, max)
			if err != nil {
				return 0, err
			}
			if crosses {
				if count == math.MaxInt {
					return 0, fmt.Errorf("intersection count overflow")
				}
				count++
			}
		}
	}
	return count, nil
}
