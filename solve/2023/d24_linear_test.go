package solve2023

import (
	"math/big"
	"strings"
	"testing"
)

func hailTestRow(values [7]int64) hailRow {
	var row hailRow
	for i, value := range values {
		row[i] = big.NewRat(value, 1)
	}
	return row
}
func TestHailLinearExactFractions(t *testing.T) {
	expected := [6]*big.Rat{big.NewRat(1, 2), big.NewRat(-2, 3), big.NewRat(5, 7), big.NewRat(11, 1), big.NewRat(-13, 1), big.NewRat(17, 1)}
	var system hailLinearSystem
	for i := hailUnknowns - 1; i >= 0; i-- {
		row := hailTestRow([7]int64{})
		row[i].SetInt64(2 * int64(i+1))
		row[6].Mul(row[i], expected[i])
		if i+1 < hailUnknowns {
			row[i+1].SetInt64(1)
			row[6].Add(row[6], expected[i+1])
		}
		for repeat := 0; repeat < 2; repeat++ {
			if err := system.addRow(row); err != nil {
				t.Fatal(err)
			}
		}
		// The system must own its coefficients independently of the caller.
		for _, value := range row {
			value.SetInt64(999)
		}
	}
	if err := system.addRow(hailTestRow([7]int64{})); err != nil {
		t.Fatal(err)
	}
	got, err := system.solution()
	if err != nil {
		t.Fatal(err)
	}
	for i := range expected {
		if got[i].Cmp(expected[i]) != 0 {
			t.Fatal("fraction mismatch", i)
		}
	}
}
func TestHailLinearCoupledEquations(t *testing.T) {
	var system hailLinearSystem
	for _, values := range [][7]int64{{1, 1, 0, 0, 0, 0, 5}, {1, -1, 0, 0, 0, 0, 1}, {0, 0, 1, 0, 0, 0, 0}, {0, 0, 0, 1, 0, 0, 0}, {0, 0, 0, 0, 1, 0, 0}, {0, 0, 0, 0, 0, 1, 0}} {
		if err := system.addRow(hailTestRow(values)); err != nil {
			t.Fatal(err)
		}
	}
	got, err := system.solution()
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range [6]int64{3, 2, 0, 0, 0, 0} {
		if got[i].Cmp(big.NewRat(want, 1)) != 0 {
			t.Fatal("elimination mismatch")
		}
	}
}
func TestHailLinearInsufficientAndContradictory(t *testing.T) {
	var system hailLinearSystem
	if _, err := system.solution(); err == nil || !strings.Contains(err.Error(), "unique rock") {
		t.Fatal("missing rank error")
	}
	row := hailTestRow([7]int64{1, 0, 0, 0, 0, 0, 0})
	row[6].SetFrac64(1, 3)
	if err := system.addRow(row); err != nil {
		t.Fatal(err)
	}
	if _, err := system.solution(); err == nil {
		t.Fatal("accepted insufficient equations")
	}
	row[6].Add(row[6], big.NewRat(1, 1000000000))
	if err := system.addRow(row); err == nil || !strings.Contains(err.Error(), "inconsistent") {
		t.Fatal("missed exact contradiction")
	}
}
