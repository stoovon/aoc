package solve2023

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func hailCrossesInWindow(a, b hailExactStone, min, max *big.Int) (bool, error) {
	var scratch hailIntersectionScratch
	return scratch.crosses(a, b, min, max)
}

func hailCrossingReference(a, b hailExactStone, min, max *big.Int) (bool, error) {
	if a.velocity[0].Sign() == 0 && a.velocity[1].Sign() == 0 {
		return hailStationaryCrossing(a, b, min, max), nil
	}
	if b.velocity[0].Sign() == 0 && b.velocity[1].Sign() == 0 {
		return hailStationaryCrossing(b, a, min, max), nil
	}
	var delta [3]*big.Int
	for axis := range delta {
		delta[axis] = new(big.Int).Sub(b.position[axis], a.position[axis])
	}
	determinant := hailCrossXY(a.velocity, b.velocity)
	if determinant.Sign() == 0 {
		if hailCrossXY(delta, a.velocity).Sign() == 0 {
			return false, fmt.Errorf("collinear projected paths do not have a unique line intersection")
		}
		return false, nil
	}
	aTime, bTime := hailCrossXY(delta, b.velocity), hailCrossXY(delta, a.velocity)
	if determinant.Sign() < 0 {
		determinant.Neg(determinant)
		aTime.Neg(aTime)
		bTime.Neg(bTime)
	}
	if aTime.Sign() < 0 || bTime.Sign() < 0 {
		return false, nil
	}
	// Compare scaled integer coordinates instead of allocating fractions per pair.
	lower, upper := new(big.Int).Mul(min, determinant), new(big.Int).Mul(max, determinant)
	for axis := 0; axis < 2; axis++ {
		coordinate := new(big.Int).Add(new(big.Int).Mul(a.position[axis], determinant), new(big.Int).Mul(a.velocity[axis], aTime))
		if coordinate.Cmp(lower) < 0 || coordinate.Cmp(upper) > 0 {
			return false, nil
		}
	}
	return true, nil
}
func TestHailProjectedVerticalTimeAndBoundaries(t *testing.T) {
	for _, tc := range []struct {
		input    string
		min, max int64
		want     int
	}{
		{"0,0,9 @ 0,2,-1\n-3,4,7 @ 1,0,4", 0, 4, 1},
		{"0,0,9 @ 0,-2,-1\n-3,4,7 @ 1,0,4", 0, 4, 0},
		{"4,0,9 @ 0,2,0\n0,4,7 @ 1,0,0", 4, 4, 1},
		{"0,0,0 @ 1,1,0\n1,0,0 @ -1,1,0", 0, 1, 1},
		{"0,0,0 @ 1,1,0\n0,1,0 @ 1,1,0", 0, 4, 0},
		{"4,4,0 @ 0,0,1\n0,4,0 @ 1,0,1", 0, 4, 1},
		{"4,4,0 @ 0,0,1\n4,4,0 @ 0,0,-1", 4, 4, 1},
		{"4,4,0 @ 1,0,1\n4,4,0 @ 0,1,-1", 4, 4, 1},
	} {
		got, err := hailCountCrossings(tc.input, big.NewInt(tc.min), big.NewInt(tc.max))
		if err != nil || got != tc.want {
			t.Fatalf("got %d,%v want %d", got, err, tc.want)
		}
	}
	if _, err := hailCountCrossings("0,0,0 @ 1,1,0\n1,1,0 @ 1,1,0", big.NewInt(0), big.NewInt(4)); err == nil {
		t.Fatal("ambiguous line crossing")
	}
	if _, err := hailCountCrossings("0,0,0 @ 1,1,0", big.NewInt(4), big.NewInt(0)); err == nil {
		t.Fatal("reversed bounds")
	}
}
func TestHailProjectedLargeExactBoundary(t *testing.T) {
	boundary := new(big.Int).Exp(big.NewInt(10), big.NewInt(40), nil)
	previous := new(big.Int).Sub(boundary, big.NewInt(1))
	input := fmt.Sprintf("%s,%s,0 @ 0,1,0\n%s,%s,0 @ 3,1,0", boundary, boundary, previous, boundary)
	if got, err := hailCountCrossings(input, boundary, boundary); err != nil || got != 0 {
		t.Fatal("rounded fractional point onto boundary", got, err)
	}
}
func TestHailProjectedConstructedIntersections(t *testing.T) {
	checked := 0
	for seed := int64(0); seed < 512; seed++ {
		x, y := seed%17-8, (seed/17)%17-8
		av, bv := [3]int64{seed%7 - 3, (seed/7)%7 - 3, 0}, [3]int64{(seed/11)%7 - 3, (seed/13)%7 - 3, 0}
		if av[0]*bv[1] == av[1]*bv[0] {
			continue
		}
		at, bt := (seed/3)%7-3, (seed/5)%7-3
		var a, b hailExactStone
		for axis, value := range [3]int64{x - av[0]*at, y - av[1]*at, 0} {
			a.position[axis] = big.NewInt(value)
			a.velocity[axis] = big.NewInt(av[axis])
		}
		for axis, value := range [3]int64{x - bv[0]*bt, y - bv[1]*bt, 0} {
			b.position[axis] = big.NewInt(value)
			b.velocity[axis] = big.NewInt(bv[axis])
		}
		want := at >= 0 && bt >= 0 && x >= -5 && x <= 5 && y >= -5 && y <= 5
		for _, pair := range [][2]hailExactStone{{a, b}, {b, a}} {
			got, err := hailCrossesInWindow(pair[0], pair[1], big.NewInt(-5), big.NewInt(5))
			reference, referenceErr := hailCrossingReference(pair[0], pair[1], big.NewInt(-5), big.NewInt(5))
			if referenceErr != nil || got != reference {
				t.Fatal("scratch differs from reference", seed, referenceErr)
			}
			if err != nil || got != want {
				t.Fatal("constructed crossing mismatch", seed, err)
			}
		}
		checked++
	}
	if checked < 300 {
		t.Fatal("too few reference cases")
	}
}

func TestHailProjectedScratchReuseAndInputOwnership(t *testing.T) {
	input := "0,0,0 @ 1,2,0\n4,3,0 @ -2,1,0\n10000000000000000000000000000000000000000,0,0 @ -3,1,0\n0,5,0 @ 1,2,0\n0,0,0 @ 0,0,1"
	stones, err := hailParseExact(input)
	if err != nil {
		t.Fatal(err)
	}
	original, err := hailParseExact(input)
	if err != nil {
		t.Fatal(err)
	}
	var scratch hailIntersectionScratch
	for repeat := 0; repeat < 3; repeat++ {
		for _, a := range stones {
			for _, b := range stones {
				got, err := scratch.crosses(a, b, big.NewInt(-5), big.NewInt(5))
				want, referenceErr := hailCrossingReference(a, b, big.NewInt(-5), big.NewInt(5))
				if (err == nil) != (referenceErr == nil) || got != want {
					t.Fatal("reused scratch mismatch")
				}
			}
		}
	}
	if !reflect.DeepEqual(stones, original) {
		t.Fatal("scratch mutated parsed hailstones")
	}
}
