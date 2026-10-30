package solve2023

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func hailTestScene(position, velocity [3]int64, fractional bool) string {
	directions := [][3]int64{{2, 4, -6}, {-4, 6, 2}, {8, -2, 4}, {6, 8, -4}, {-6, -4, 8}, {4, -8, -2}}
	var input strings.Builder
	for i, direction := range directions {
		numerator, denominator := int64(2*i+1), int64(1)
		if fractional {
			denominator = 2
		}
		var p, u [3]int64
		for axis := 0; axis < 3; axis++ {
			p[axis] = position[axis] + direction[axis]*numerator/denominator
			u[axis] = velocity[axis] - direction[axis]
		}
		fmt.Fprintf(&input, "%d, %d, %d @ %d, %d, %d\n", p[0], p[1], p[2], u[0], u[1], u[2])
	}
	return input.String()
}
func TestHailExactGeneratedRocks(t *testing.T) {
	for seed := int64(0); seed < 64; seed++ {
		position, velocity := [3]int64{17 + seed, -23 + 2*seed, 41 - seed}, [3]int64{701 + seed, -907 - seed, 1201 + seed}
		for _, fractional := range []bool{false, true} {
			input := hailTestScene(position, velocity, fractional)
			stones, err := hailParseExact(input)
			if err != nil {
				t.Fatal(err)
			}
			rock, err := hailSolveExact(stones)
			if err != nil {
				t.Fatal(err)
			}
			for axis := 0; axis < 3; axis++ {
				if rock[axis].Cmp(big.NewRat(position[axis], 1)) != 0 || rock[axis+3].Cmp(big.NewRat(velocity[axis], 1)) != 0 {
					t.Fatal("generated rock mismatch")
				}
			}
			got, err := (Day24{}).Part2(input)
			want := fmt.Sprint(position[0] + position[1] + position[2])
			if err != nil || got != want {
				t.Fatal("part2 mismatch", err)
			}
			time, err := hailCollisionTime(rock, stones[0])
			if err != nil || time.IsInt() == fractional {
				t.Fatal("wrong time integrality", err)
			}
		}
	}
}
func TestHailExactDependentAndLaterStones(t *testing.T) {
	input := hailTestScene([3]int64{17, -23, 41}, [3]int64{701, -907, 1201}, true)
	first, _, _ := strings.Cut(input, "\n")
	redundant := first + "\n" + first + "\n" + input
	if got, err := (Day24{}).Part2(redundant); err != nil || got != "35" {
		t.Fatal("dependent prefix", err)
	}
	stones, err := hailParseExact(input)
	if err != nil {
		t.Fatal(err)
	}
	stones[len(stones)-1].position[0].Add(stones[len(stones)-1].position[0], big.NewInt(1))
	if _, err := hailSolveExact(stones); err == nil {
		t.Fatal("ignored later bad stone")
	}
	if _, err := (Day24{}).Part2(first); err == nil || !strings.Contains(err.Error(), "unique rock") {
		t.Fatal("missing rank error")
	}
}
func TestHailExactCollisionTimes(t *testing.T) {
	var rock [6]*big.Rat
	for i, v := range [6]int64{17, -23, 41, 7, -9, 12} {
		rock[i] = big.NewRat(v, 1)
	}
	var stone hailExactStone
	for i, v := range [3]int64{17, -23, 41} {
		stone.position[i] = big.NewInt(v)
	}
	for i, v := range [3]int64{7, -9, 12} {
		stone.velocity[i] = big.NewInt(v)
	}
	time, err := hailCollisionTime(rock, stone)
	if err != nil || time.Sign() != 0 {
		t.Fatal("shared path", err)
	}
	stone.position[0].Add(stone.position[0], big.NewInt(1))
	if _, err := hailCollisionTime(rock, stone); err == nil {
		t.Fatal("parallel displaced path")
	}
	stone.velocity[0].Sub(stone.velocity[0], big.NewInt(2))
	time, err = hailCollisionTime(rock, stone)
	if err != nil || time.Cmp(big.NewRat(1, 2)) != 0 {
		t.Fatal("fractional time", err)
	}
	stone.position[0].Sub(stone.position[0], big.NewInt(2))
	if _, err := hailCollisionTime(rock, stone); err == nil {
		t.Fatal("past collision accepted")
	}
}
func TestHailVerifiesOriginalCollisionEquations(t *testing.T) {
	var input strings.Builder
	for _, triple := range [][3]int64{{1, 2, 3}, {2, 3, 5}, {3, 5, 7}, {5, 7, 11}, {7, 11, 13}, {11, 13, 17}} {
		a, b, c := triple[0], triple[1], triple[2]
		p := [3]int64{1 + a*b, b, 0}
		u := [3]int64{c*p[0] + a, c*p[1] + 1, 0}
		fmt.Fprintf(&input, "%d,%d,0 @ %d,%d,0\n", p[0], p[1], u[0], u[1])
	}
	stones, err := hailParseExact(input.String())
	if err != nil {
		t.Fatal(err)
	}
	var system hailLinearSystem
	for _, stone := range stones[1:] {
		if err := hailAddPair(&system, stones[0], stone); err != nil {
			t.Fatal(err)
		}
	}
	solution, err := system.solution()
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range solution {
		if value.Sign() != 0 {
			t.Fatal("expected zero linear solution")
		}
	}
	if _, err := (Day24{}).Part2(input.String()); err == nil {
		t.Fatal("accepted solution to differences that fails original collisions")
	}
}
func TestHailRejectsFractionalRocksAndPastCollisions(t *testing.T) {
	p2, v2 := [3]int64{35, -45, 83}, [3]int64{1403, -1813, 2403}
	var input strings.Builder
	for i, u := range [][3]int64{{-3, 2, 4}, {5, -7, 2}, {8, 9, -5}, {-11, 13, 17}, {19, -23, 29}, {-31, 37, 41}} {
		time := int64(2*i + 1)
		var p [3]int64
		for axis := 0; axis < 3; axis++ {
			p[axis] = (p2[axis]+v2[axis]*time)/2 - u[axis]*time
		}
		fmt.Fprintf(&input, "%d,%d,%d @ %d,%d,%d\n", p[0], p[1], p[2], u[0], u[1], u[2])
	}
	if _, err := (Day24{}).Part2(input.String()); err == nil || !strings.Contains(err.Error(), "must be integral") {
		t.Fatal("fractional rock was not rejected", err)
	}
	past := hailTestScene([3]int64{17, -23, 41}, [3]int64{701, -907, 1201}, true) + fmt.Sprintf("%d,%d,%d @ 5,7,11\n", 17-701+5, -23+907+7, 41-1201+11)
	if _, err := (Day24{}).Part2(past); err == nil || !strings.Contains(err.Error(), "nonnegative time") {
		t.Fatal("past collision accepted", err)
	}
}
func TestHailOrderAndLargeIntegerAnswer(t *testing.T) {
	input := hailTestScene([3]int64{17, -23, 41}, [3]int64{701, -907, 1201}, true)
	lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	for rotation := 0; rotation < len(lines); rotation++ {
		ordered := append(append([]string{}, lines[rotation:]...), lines[:rotation]...)
		for repeat := 0; repeat < 2; repeat++ {
			got, err := (Day24{}).Part2(strings.Join(ordered, "\n"))
			if err != nil || got != "35" {
				t.Fatal("input order changed solution", err)
			}
			for i, j := 0, len(ordered)-1; i < j; i, j = i+1, j-1 {
				ordered[i], ordered[j] = ordered[j], ordered[i]
			}
		}
	}
	shift := new(big.Int).Exp(big.NewInt(10), big.NewInt(80), nil)
	stones, err := hailParseExact(input)
	if err != nil {
		t.Fatal(err)
	}
	var shifted strings.Builder
	for _, stone := range stones {
		fmt.Fprintf(&shifted, "%s,%s,%s @ %s,%s,%s\n", new(big.Int).Add(stone.position[0], shift), new(big.Int).Add(stone.position[1], shift), new(big.Int).Add(stone.position[2], shift), stone.velocity[0], stone.velocity[1], stone.velocity[2])
	}
	want := new(big.Int).Add(new(big.Int).Mul(shift, big.NewInt(3)), big.NewInt(35)).String()
	got, err := (Day24{}).Part2(shifted.String())
	if err != nil || got != want {
		t.Fatal("large answer failed", err)
	}
}
func TestHailVariedGeneratedGeometries(t *testing.T) {
	for seed := uint64(0); seed < 32; seed++ {
		position, velocity := [3]int64{17 + int64(seed), -23 + int64(seed), 41 - int64(seed)}, [3]int64{701, -907, 1201}
		state := seed + 1
		var input strings.Builder
		for i := int64(0); i < 8; i++ {
			var relative, p, u [3]int64
			for axis := 0; axis < 3; axis++ {
				state = state*1664525 + 1013904223
				relative[axis] = 2 * (int64(state%19) - 9)
				p[axis] = position[axis] + relative[axis]*(2*i+1)/2
				u[axis] = velocity[axis] - relative[axis]
			}
			fmt.Fprintf(&input, "%d,%d,%d @ %d,%d,%d\n", p[0], p[1], p[2], u[0], u[1], u[2])
		}
		stones, err := hailParseExact(input.String())
		if err != nil {
			t.Fatal(err)
		}
		rock, err := hailSolveExact(stones)
		if err != nil {
			t.Fatal(err)
		}
		for axis := 0; axis < 3; axis++ {
			if rock[axis].Cmp(big.NewRat(position[axis], 1)) != 0 || rock[axis+3].Cmp(big.NewRat(velocity[axis], 1)) != 0 {
				t.Fatal("varied geometry mismatch", seed)
			}
		}
	}
}
func TestHailCommonParserBothPartsAndCRLF(t *testing.T) {
	for _, input := range []string{"", " \n", "1,2 @ 3,4,5", "1,2,3,4 @ 5,6,7", "1,2,3 @ 4,5,6 @ 7,8,9", "NaN,2,3 @ 4,5,6", "1.5,2,3 @ 4,5,6", "1,2,3 @ 4,5,6\n\n", "--1,2,3 @ 4,5,6", "-+1,2,3 @ 4,5,6"} {
		for _, solve := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
			if _, err := solve(input); err == nil {
				t.Fatal("accepted invalid record")
			}
		}
	}
	for _, solve := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
		if _, err := solve("1,2,3 @ 4,5,6\nbad"); err == nil || !strings.HasPrefix(err.Error(), "line 2:") {
			t.Fatal("missing line context")
		}
	}
	input := strings.ReplaceAll(hailTestScene([3]int64{17, -23, 41}, [3]int64{701, -907, 1201}, true), "\n", "\r\n")
	if got, err := (Day24{}).Part2(input); err != nil || got != "35" {
		t.Fatal("CRLF failed", err)
	}
	stones, err := hailParseExact(" +1, -2, 0003 @ -4, +5, 0\r\n")
	if err != nil || stones[0].position[2].Cmp(big.NewInt(3)) != 0 {
		t.Fatal("signed integers", err)
	}
}
