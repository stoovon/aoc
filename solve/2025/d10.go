package solve2025

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day10 struct{}

func (d Day10) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2025, Day: 10}
}

func (d Day10) Part1(input string) (string, error) {
	return sumMachinePresses(input, func(line string) (int64, error) {
		count, err := minPressesForMachine(line)
		return int64(count), err
	})
}

func (d Day10) Part2(input string) (string, error) {
	return sumMachinePresses(input, minPressesForMachineJolts)
}

func sumMachinePresses(input string, minimum func(string) (int64, error)) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "0", nil
	}
	var total int64
	for line, record := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {
		count, err := minimum(record)
		if err != nil {
			return "", fmt.Errorf("line %d: %w", line+1, err)
		}
		if total > math.MaxInt64-count {
			return "", fmt.Errorf("line %d: total press count exceeds int64", line+1)
		}
		total += count
	}
	return strconv.FormatInt(total, 10), nil
}

func init() {
	solve.Register(Day10{})
}

func minPressesForMachineJolts(line string) (int64, error) {
	_, buttonToggles, targets, err := parseMachineLine(line)
	if err != nil {
		return 0, err
	}
	if len(targets) == 0 {
		return 0, fmt.Errorf("missing joltage targets")
	}

	k := len(targets)
	buttons := distinctMachineButtons(buttonToggles, k)
	sigs := make([]string, 0, len(buttons))
	for _, indices := range buttons {
		sigs = append(sigs, machineButtonSignature(indices, k))
	}
	numSigs := len(sigs)

	mat := make([][]big.Rat, k)
	for i := 0; i < k; i++ {
		mat[i] = make([]big.Rat, numSigs+1)
		mat[i][numSigs].SetInt64(targets[i])
	}
	for j, sig := range sigs {
		for i := 0; i < k; i++ {
			if sig[i/8]&(1<<(i%8)) != 0 {
				mat[i][j].SetInt64(1)
			}
		}
	}

	pivotRow := make([]int, numSigs)
	for j := range pivotRow {
		pivotRow[j] = -1
	}

	row := 0
	for col := 0; col < numSigs && row < k; col++ {
		pivot := -1
		for i := row; i < k; i++ {
			if mat[i][col].Sign() != 0 {
				pivot = i
				break
			}
		}
		if pivot == -1 {
			continue
		}

		mat[row], mat[pivot] = mat[pivot], mat[row]
		pivotRow[col] = row

		var pv big.Rat
		pv.Set(&mat[row][col])
		for j := col; j <= numSigs; j++ {
			mat[row][j].Quo(&mat[row][j], &pv)
		}

		for i := 0; i < k; i++ {
			if i == row || mat[i][col].Sign() == 0 {
				continue
			}
			var factor, product big.Rat
			factor.Set(&mat[i][col])
			for j := col; j <= numSigs; j++ {
				product.Mul(&factor, &mat[row][j])
				mat[i][j].Sub(&mat[i][j], &product)
			}
		}
		row++
	}
	numPivotRows := row

	for i := 0; i < k; i++ {
		allZero := true
		for j := 0; j < numSigs; j++ {
			if mat[i][j].Sign() != 0 {
				allZero = false
				break
			}
		}
		if allZero && mat[i][numSigs].Sign() != 0 {
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
	freeBounds := make([]int64, numFree)
	for i, col := range freeCols {
		first := true
		for counter, target := range targets {
			if sigs[col][counter/8]&(1<<(counter%8)) != 0 && (first || target < freeBounds[i]) {
				freeBounds[i] = target
				first = false
			}
		}
	}

	var best *big.Int

	var enumerate func(freeIdx int, freeVals []int64, subtotal *big.Int)
	enumerate = func(freeIdx int, freeVals []int64, subtotal *big.Int) {
		if best != nil && subtotal.Cmp(best) >= 0 {
			return
		}
		if freeIdx == numFree {
			var total big.Int
			total.Set(subtotal)
			var result, product, coefficient big.Rat
			for r := 0; r < numPivotRows; r++ {
				result.Set(&mat[r][numSigs])
				for i, col := range freeCols {
					coefficient.SetInt64(freeVals[i])
					product.Mul(&mat[r][col], &coefficient)
					result.Sub(&result, &product)
				}
				if !result.IsInt() || result.Sign() < 0 {
					return
				}
				total.Add(&total, result.Num())
			}
			if best == nil || total.Cmp(best) < 0 {
				best = new(big.Int).Set(&total)
			}
			return
		}

		var nextTotal, value big.Int
		for v := int64(0); ; v++ {
			value.SetInt64(v)
			nextTotal.Add(subtotal, &value)
			if best != nil && nextTotal.Cmp(best) >= 0 {
				break
			}
			freeVals[freeIdx] = v
			enumerate(freeIdx+1, freeVals, &nextTotal)
			if v == freeBounds[freeIdx] {
				break
			}
		}
	}

	freeVals := make([]int64, numFree)
	enumerate(0, freeVals, new(big.Int))

	if best == nil {
		return 0, fmt.Errorf("no non-negative solution found")
	}
	if !best.IsInt64() {
		return 0, fmt.Errorf("minimum press count exceeds int64")
	}
	return best.Int64(), nil
}

func machineNumber(token string) (int64, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0, fmt.Errorf("empty integer")
	}
	for i := 0; i < len(token); i++ {
		if token[i] < '0' || token[i] > '9' {
			return 0, fmt.Errorf("expected nonnegative decimal integer")
		}
	}
	value, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("integer exceeds int64")
	}
	return value, nil
}

func parseMachineLine(line string) (pattern string, buttons [][]int, targets []int64, err error) {
	line = strings.TrimSpace(line)
	fail := func(message string) (string, [][]int, []int64, error) { return "", nil, nil, fmt.Errorf("%s", message) }
	if !strings.HasPrefix(line, "[") {
		return fail("expected light pattern")
	}
	end := strings.IndexByte(line, ']')
	if end < 2 {
		return fail("invalid light pattern")
	}
	pattern = line[1:end]
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != '.' && pattern[i] != '#' {
			return fail("light pattern must contain only . and #")
		}
	}
	rest := strings.TrimSpace(line[end+1:])
	for strings.HasPrefix(rest, "(") {
		end := strings.IndexByte(rest, ')')
		if end < 0 {
			return fail("unclosed button descriptor")
		}
		fields := strings.Split(rest[1:end], ",")
		indices := make([]int, 0, len(fields))
		seen := make(map[int]bool, len(fields))
		for _, field := range fields {
			value, parseErr := machineNumber(field)
			if parseErr != nil {
				return fail("invalid button index: " + parseErr.Error())
			}
			if value >= int64(len(pattern)) {
				return fail("button index out of range")
			}
			index := int(value)
			if seen[index] {
				return fail(fmt.Sprintf("repeated button index %d", index))
			}
			seen[index] = true
			indices = append(indices, index)
		}
		buttons = append(buttons, indices)
		rest = strings.TrimSpace(rest[end+1:])
	}
	if rest == "" {
		return pattern, buttons, nil, nil
	}
	if !strings.HasPrefix(rest, "{") || !strings.HasSuffix(rest, "}") {
		return fail("expected final joltage targets")
	}
	for _, field := range strings.Split(rest[1:len(rest)-1], ",") {
		value, parseErr := machineNumber(field)
		if parseErr != nil {
			return fail("invalid joltage target: " + parseErr.Error())
		}
		targets = append(targets, value)
	}
	if len(targets) != len(pattern) {
		return fail("light pattern and joltage target lengths differ")
	}
	return pattern, buttons, targets, nil
}

func machineButtonSignature(indices []int, width int) string {
	bits := make([]byte, (width+7)/8)
	for _, index := range indices {
		bits[index/8] |= 1 << (index % 8)
	}
	return string(bits)
}

func distinctMachineButtons(buttons [][]int, width int) [][]int {
	unique := make([][]int, 0, len(buttons))
	seen := make(map[string]bool, len(buttons))
	for _, indices := range buttons {
		signature := machineButtonSignature(indices, width)
		if !seen[signature] {
			seen[signature] = true
			unique = append(unique, indices)
		}
	}
	return unique
}

func minPressesForMachine(line string) (int, error) {
	pattern, buttonToggles, _, err := parseMachineLine(line)
	if err != nil {
		return 0, err
	}

	k := len(pattern)
	target := make([]uint8, k)
	for i := 0; i < k; i++ {
		if pattern[i] == '#' {
			target[i] = 1
		}
	}

	uniqueButtons := distinctMachineButtons(buttonToggles, k)
	m := len(uniqueButtons)

	mat := make([][]uint8, k)
	for i := 0; i < k; i++ {
		row := make([]uint8, m+1)
		mat[i] = row
	}

	for j, idxs := range uniqueButtons {
		for _, v := range idxs {
			mat[v][j] ^= 1
		}
	}
	for i := 0; i < k; i++ {
		mat[i][m] = target[i]
	}

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
	if row < len(freeCols) {
		return minimumLightPatternPresses(target, uniqueButtons)
	}
	best := m + 1
	x := make([]uint8, m)
	var enumerate func(int, int)
	enumerate = func(index, freePresses int) {
		if freePresses >= best {
			return
		}
		if index < len(freeCols) {
			col := freeCols[index]
			x[col] = 0
			enumerate(index+1, freePresses)
			x[col] = 1
			enumerate(index+1, freePresses+1)
			x[col] = 0
			return
		}
		total := freePresses
		for _, r := range where {
			if r < 0 {
				continue
			}
			value := mat[r][m]
			for _, free := range freeCols {
				value ^= mat[r][free] & x[free]
			}
			total += int(value)
		}
		if total < best {
			best = total
		}
	}
	enumerate(0, 0)
	return best, nil
}

func minimumLightPatternPresses(target []uint8, buttons [][]int) (int, error) {
	width := (len(target) + 7) / 8
	goal := make([]byte, width)
	for index, value := range target {
		if value != 0 {
			goal[index/8] |= 1 << (index % 8)
		}
	}
	goalKey := string(goal)
	masks := make([]string, 0, len(buttons))
	for _, indices := range buttons {
		masks = append(masks, machineButtonSignature(indices, len(target)))
	}
	initial := string(make([]byte, width))
	queue := []string{initial}
	seen := map[string]bool{initial: true}
	head := 0
	for presses := 0; head < len(queue); presses++ {
		end := len(queue)
		for ; head < end; head++ {
			current := queue[head]
			if current == goalKey {
				return presses, nil
			}
			for _, mask := range masks {
				bits := make([]byte, width)
				for i := 0; i < width; i++ {
					bits[i] = current[i] ^ mask[i]
				}
				next := string(bits)
				if !seen[next] {
					seen[next] = true
					queue = append(queue, next)
				}
			}
		}
	}
	return 0, fmt.Errorf("no solution for machine")
}
