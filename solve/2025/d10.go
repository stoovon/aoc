package solve2025

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"aoc/solve"
)

type Day10 struct {
}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 10}
}

func (d Day10) Part1(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "0", nil
	}

	var sum int64
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		presses, err := minPressesForMachine(line)
		if err != nil {
			return "", fmt.Errorf("line %d: %w", i+1, err)
		}
		sum += int64(presses)
	}

	return strconv.FormatInt(sum, 10), nil
}

func (d Day10) Part2(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "0", nil
	}

	var sum int64
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		start := time.Now()
		presses, err := minPressesForMachineJolts(line)
		elapsed := time.Since(start)
		if elapsed > 500*time.Millisecond {
			fmt.Printf("[Day10 Part2] slow machine line %d (elapsed %s): %s\n", i+1, elapsed, line)
		}
		if err != nil {
			return "", fmt.Errorf("line %d: %w\n  input: %s", i+1, err, line)
		}
		sum += presses
	}

	return strconv.FormatInt(sum, 10), nil
}

func init() {
	solve.Register(Day10{})
}

// minPressesForMachineJolts computes the minimal number of button presses
// needed to reach the joltage targets for a single machine description line.
// Uses integer linear algebra: solve Ax = b, find free variables (at most 3),
// enumerate all combinations up to maxTarget, and minimize sum(x).
func minPressesForMachineJolts(line string) (int64, error) {
	pattern, buttonToggles, targets, err := parseMachineLine(line)
	if err != nil {
		return 0, err
	}
	if pattern == "" {
		return 0, nil
	}
	if len(targets) == 0 {
		return 0, fmt.Errorf("missing joltage targets")
	}

	k := len(targets)       // number of counters (rows)
	m := len(buttonToggles) // number of buttons (columns)

	maxTarget := int64(0)
	for _, t := range targets {
		if t < 0 {
			return 0, fmt.Errorf("negative joltage target %d", t)
		}
		if t > maxTarget {
			maxTarget = t
		}
	}

	if m == 0 {
		for i := 0; i < k; i++ {
			if targets[i] != 0 {
				return 0, fmt.Errorf("no buttons to satisfy non-zero joltage")
			}
		}
		return 0, nil
	}

	// Since buttons with the same signature are interchangeable, we work with signatures.
	sigMap := make(map[uint16]int)
	for _, idxs := range buttonToggles {
		if len(idxs) == 0 {
			continue
		}
		var sig uint16
		for _, i := range idxs {
			sig |= 1 << i
		}
		sigMap[sig]++
	}

	if len(sigMap) == 0 {
		for i := 0; i < k; i++ {
			if targets[i] != 0 {
				return 0, fmt.Errorf("no usable buttons")
			}
		}
		return 0, nil
	}

	sigs := make([]uint16, 0, len(sigMap))
	for s := range sigMap {
		sigs = append(sigs, s)
	}
	numSigs := len(sigs)

	// Use rational arithmetic for Gaussian elimination.
	// Represent each entry as (numerator, denominator) pair.
	type rat struct {
		num, den int64
	}

	gcd := func(a, b int64) int64 {
		if a < 0 {
			a = -a
		}
		if b < 0 {
			b = -b
		}
		for b != 0 {
			a, b = b, a%b
		}
		if a == 0 {
			return 1
		}
		return a
	}

	simplify := func(r rat) rat {
		if r.num == 0 {
			return rat{0, 1}
		}
		g := gcd(r.num, r.den)
		r.num /= g
		r.den /= g
		if r.den < 0 {
			r.num = -r.num
			r.den = -r.den
		}
		return r
	}

	mat := make([][]rat, k)
	for i := 0; i < k; i++ {
		mat[i] = make([]rat, numSigs+1)
		for j := 0; j <= numSigs; j++ {
			mat[i][j] = rat{0, 1}
		}
		mat[i][numSigs] = rat{targets[i], 1}
	}
	for j, sig := range sigs {
		for i := 0; i < k; i++ {
			if sig&(1<<i) != 0 {
				mat[i][j] = rat{1, 1}
			}
		}
	}

	// Gaussian elimination to RREF.
	pivotCol := make([]int, k)
	for i := range pivotCol {
		pivotCol[i] = -1
	}
	pivotRow := make([]int, numSigs)
	for j := range pivotRow {
		pivotRow[j] = -1
	}

	row := 0
	for col := 0; col < numSigs && row < k; col++ {
		pivot := -1
		for i := row; i < k; i++ {
			if mat[i][col].num != 0 {
				pivot = i
				break
			}
		}
		if pivot == -1 {
			continue
		}

		mat[row], mat[pivot] = mat[pivot], mat[row]
		pivotCol[row] = col
		pivotRow[col] = row

		pv := mat[row][col]
		for j := col; j <= numSigs; j++ {
			mat[row][j] = simplify(rat{mat[row][j].num * pv.den, mat[row][j].den * pv.num})
		}

		for i := 0; i < k; i++ {
			if i == row || mat[i][col].num == 0 {
				continue
			}
			factor := mat[i][col]
			for j := col; j <= numSigs; j++ {
				// a/b - c/d * e/f = a/b - ce/df = (adf - bce) / bdf
				a, b := mat[i][j].num, mat[i][j].den
				c, d := factor.num, factor.den
				e, f := mat[row][j].num, mat[row][j].den
				num := a*d*f - b*c*e
				den := b * d * f
				mat[i][j] = simplify(rat{num, den})
			}
		}
		row++
	}
	numPivotRows := row

	for i := 0; i < k; i++ {
		allZero := true
		for j := 0; j < numSigs; j++ {
			if mat[i][j].num != 0 {
				allZero = false
				break
			}
		}
		if allZero && mat[i][numSigs].num != 0 {
			return 0, fmt.Errorf("no solution for machine")
		}
	}

	freeCols := make([]int, 0)
	for col := 0; col < numSigs; col++ {
		if pivotRow[col] == -1 {
			freeCols = append(freeCols, col)
		}
	}
	numFree := len(freeCols)

	// Now mat is in RREF with pivot entries = 1.
	// For pivot row r with pivot column c:
	// x[c] + sum_{free f} mat[r][f] * x[f] = mat[r][numSigs]
	// So x[c] = mat[r][numSigs] - sum_{free f} mat[r][f] * x[f]

	best := int64(-1)

	var enumerate func(freeIdx int, freeVals []int64)
	enumerate = func(freeIdx int, freeVals []int64) {
		if freeIdx == numFree {
			x := make([]rat, numSigs)
			for j := 0; j < numSigs; j++ {
				x[j] = rat{0, 1}
			}
			for i, col := range freeCols {
				x[col] = rat{freeVals[i], 1}
			}

			valid := true
			for r := 0; r < numPivotRows; r++ {
				col := pivotCol[r]
				if col == -1 {
					continue
				}
				res := mat[r][numSigs]
				for _, f := range freeCols {
					a, b := res.num, res.den
					c, d := mat[r][f].num, mat[r][f].den
					e, ff := x[f].num, x[f].den
					num := a*d*ff - b*c*e
					den := b * d * ff
					res = simplify(rat{num, den})
				}
				if res.den != 1 {
					valid = false
					break
				}
				if res.num < 0 {
					valid = false
					break
				}
				x[col] = res
			}

			if valid {
				total := int64(0)
				for _, v := range x {
					if v.den != 1 || v.num < 0 {
						valid = false
						break
					}
					total += v.num
				}
				if valid && (best == -1 || total < best) {
					best = total
				}
			}
			return
		}

		for v := int64(0); v <= maxTarget; v++ {
			freeVals[freeIdx] = v
			enumerate(freeIdx+1, freeVals)
		}
	}

	freeVals := make([]int64, numFree)
	enumerate(0, freeVals)

	if best == -1 {
		return 0, fmt.Errorf("no non-negative solution found")
	}
	return best, nil
}

// parseMachineLine parses a single machine description line into its
// pattern, button wiring, and joltage targets.
func parseMachineLine(line string) (pattern string, buttonToggles [][]int, targets []int64, err error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", nil, nil, nil
	}

	lb := strings.IndexByte(line, '[')
	rb := strings.IndexByte(line, ']')
	if lb == -1 || rb == -1 || rb <= lb+1 {
		return "", nil, nil, fmt.Errorf("invalid pattern segment")
	}
	pattern = line[lb+1 : rb]
	k := len(pattern)
	if k == 0 {
		return pattern, nil, nil, nil
	}

	curly := strings.IndexByte(line, '{')
	buttonsPart := line[rb+1:]
	if curly != -1 {
		buttonsPart = line[rb+1 : curly]
	}
	buttonsPart = strings.TrimSpace(buttonsPart)

	for {
		lb = strings.IndexByte(buttonsPart, '(')
		if lb == -1 {
			break
		}
		rb = strings.IndexByte(buttonsPart[lb+1:], ')')
		if rb == -1 {
			return "", nil, nil, fmt.Errorf("unclosed button descriptor")
		}
		rb += lb + 1
		inside := strings.TrimSpace(buttonsPart[lb+1 : rb])
		if inside != "" {
			parts := strings.Split(inside, ",")
			idxs := make([]int, 0, len(parts))
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				v, convErr := strconv.Atoi(p)
				if convErr != nil {
					return "", nil, nil, fmt.Errorf("invalid button index %q", p)
				}
				if v < 0 || v >= k {
					return "", nil, nil, fmt.Errorf("button index %d out of range", v)
				}
				idxs = append(idxs, v)
			}
			buttonToggles = append(buttonToggles, idxs)
		}
		buttonsPart = buttonsPart[rb+1:]
	}

	if curly != -1 {
		endCurly := strings.LastIndexByte(line, '}')
		if endCurly == -1 || endCurly <= curly+1 {
			return "", nil, nil, fmt.Errorf("invalid joltage segment")
		}
		inside := strings.TrimSpace(line[curly+1 : endCurly])
		if inside != "" {
			parts := strings.Split(inside, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				v, convErr := strconv.ParseInt(p, 10, 64)
				if convErr != nil {
					return "", nil, nil, fmt.Errorf("invalid joltage %q", p)
				}
				targets = append(targets, v)
			}
		}
	}

	if len(targets) > 0 && len(targets) != k {
		return "", nil, nil, fmt.Errorf("pattern length %d and targets length %d mismatch", k, len(targets))
	}

	return pattern, buttonToggles, targets, nil
}

// minPressesForMachine computes the minimal number of button presses needed
// to reach the target pattern for a single machine description line.
func minPressesForMachine(line string) (int, error) {
	pattern, buttonToggles, _, err := parseMachineLine(line)
	if err != nil {
		return 0, err
	}
	if pattern == "" {
		return 0, nil
	}

	k := len(pattern)
	target := make([]uint8, k)
	for i := 0; i < k; i++ {
		if pattern[i] == '#' {
			target[i] = 1
		}
	}

	m := len(buttonToggles)
	if m == 0 {
		for i := 0; i < k; i++ {
			if target[i] != 0 {
				return 0, fmt.Errorf("no buttons to satisfy non-zero pattern")
			}
		}
		return 0, nil
	}

	// Build augmented matrix mat[k][m+1] over GF(2): last column is target.
	mat := make([][]uint8, k)
	for i := 0; i < k; i++ {
		row := make([]uint8, m+1)
		mat[i] = row
	}

	for j, idxs := range buttonToggles {
		for _, v := range idxs {
			mat[v][j] ^= 1
		}
	}
	for i := 0; i < k; i++ {
		mat[i][m] = target[i]
	}

	// Gaussian elimination to reduced row-echelon form.
	where := make([]int, m)
	for i := range where {
		where[i] = -1
	}

	row := 0
	for col := 0; col < m && row < k; col++ {
		pivot := -1
		for i := row; i < k; i++ {
			if mat[i][col] != 0 {
				pivot = i
				break
			}
		}
		if pivot == -1 {
			continue
		}
		mat[row], mat[pivot] = mat[pivot], mat[row]
		where[col] = row

		for i := row + 1; i < k; i++ {
			if mat[i][col] != 0 {
				for j := col; j <= m; j++ {
					mat[i][j] ^= mat[row][j]
				}
			}
		}
		row++
	}

	for i := 0; i < k; i++ {
		allZero := true
		for j := 0; j < m; j++ {
			if mat[i][j] != 0 {
				allZero = false
				break
			}
		}
		if allZero && mat[i][m] != 0 {
			return 0, fmt.Errorf("no solution for machine")
		}
	}

	// Back substitution to get RREF.
	for col := m - 1; col >= 0; col-- {
		r := where[col]
		if r == -1 {
			continue
		}
		for i := 0; i < r; i++ {
			if mat[i][col] != 0 {
				for j := col; j <= m; j++ {
					mat[i][j] ^= mat[r][j]
				}
			}
		}
	}

	freeCols := make([]int, 0)
	for col := 0; col < m; col++ {
		if where[col] == -1 {
			freeCols = append(freeCols, col)
		}
	}
	d := len(freeCols)

	// Enumerate all assignments of free variables; for each, derive pivot vars.
	best := -1
	if d >= 0 && d <= 30 {
		limit := 1 << d
		for mask := 0; mask < limit; mask++ {
			x := make([]uint8, m)
			for idx, col := range freeCols {
				if (mask>>idx)&1 == 1 {
					x[col] = 1
				}
			}
			for col := 0; col < m; col++ {
				r := where[col]
				if r == -1 {
					continue
				}
				v := mat[r][m]
				for _, c := range freeCols {
					if mat[r][c] != 0 && x[c] != 0 {
						v ^= 1
					}
				}
				x[col] = v
			}

			w := 0
			for j := 0; j < m; j++ {
				if x[j] != 0 {
					w++
				}
			}
			if best == -1 || w < best {
				best = w
			}
		}
	} else {
		return 0, fmt.Errorf("too many free variables: %d", d)
	}

	return best, nil
}
