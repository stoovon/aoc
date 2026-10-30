package solve2015

import (
	"aoc/solve"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

type Day19 struct{}

func (Day19) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2015, Day: 19} }

type moleculeRule struct{ from, to string }
type moleculeProblem struct {
	rules  []moleculeRule
	target string
}
type moleculeLimits struct{ states, bytes, checks int }

var defaultMoleculeLimits = moleculeLimits{65536, 32 * 1024 * 1024, 10_000_000}
var errMoleculeLimit = errors.New("molecule search resource limit exceeded")
var errMoleculeUnreachable = errors.New("target cannot be derived from e")

func validMolecule(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range []byte(s) {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}
func parseMolecules(input string) (moleculeProblem, error) {
	var p moleculeProblem
	separated := false
	seen := map[moleculeRule]bool{}
	// Remove one terminal newline, preserving unexpected blank target records.
	lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		invalid := func(message string) error { return fmt.Errorf("line %d: %s", i+1, message) }
		if !separated {
			if line == "" {
				if len(p.rules) == 0 {
					return p, invalid("expected replacement rules")
				}
				separated = true
				continue
			}
			from, to, ok := strings.Cut(line, " => ")
			if !ok || !validMolecule(from) || !validMolecule(to) {
				return p, invalid("expected nonempty ASCII alphanumeric source => replacement")
			}
			rule := moleculeRule{from, to}
			if seen[rule] {
				return p, invalid("duplicate replacement rule")
			}
			seen[rule] = true
			p.rules = append(p.rules, rule)
		} else {
			if p.target != "" || !validMolecule(line) {
				return p, invalid("expected exactly one ASCII alphanumeric target")
			}
			p.target = line
		}
	}
	if p.target == "" {
		return p, errors.New("expected blank separator and one target molecule")
	}
	return p, nil
}
func moleculeChecks(text, pattern string) int {
	if len(pattern) > len(text) {
		return 0
	}
	return len(text) - len(pattern) + 1
}
func spendMoleculeChecks(remaining *int, text, pattern string) error {
	n := moleculeChecks(text, pattern)
	if n > *remaining {
		return errMoleculeLimit
	}
	*remaining -= n
	return nil
}
func replaceMolecule(text string, offset int, from, to string) string {
	return text[:offset] + to + text[offset+len(from):]
}
func (p moleculeProblem) oneStep(limits moleculeLimits) (map[string]struct{}, error) {
	result := map[string]struct{}{}
	bytes, checks := 0, limits.checks
	for _, rule := range p.rules {
		if err := spendMoleculeChecks(&checks, p.target, rule.from); err != nil {
			return nil, err
		}
		for start := 0; start <= len(p.target)-len(rule.from); {
			found := strings.Index(p.target[start:], rule.from)
			if found < 0 {
				break
			}
			offset := start + found
			start = offset + 1
			length := len(p.target) - len(rule.from) + len(rule.to)
			if length > limits.bytes {
				return nil, errMoleculeLimit
			}
			next := replaceMolecule(p.target, offset, rule.from, rule.to)
			if _, ok := result[next]; ok {
				continue
			}
			if len(result) >= limits.states || length > limits.bytes-bytes {
				return nil, errMoleculeLimit
			}
			bytes += length
			result[next] = struct{}{}
		}
	}
	return result, nil
}

// A verified character potential increases by one for every rule, in any context.
// It fixes path length, but only an actual path establishes reachability.
type moleculeCounter struct{ weights [128]big.Rat }

func (c *moleculeCounter) total(text string) *big.Rat {
	total := new(big.Rat)
	for _, b := range []byte(text) {
		total.Add(total, &c.weights[b])
	}
	return total
}
func (c *moleculeCounter) verifies(rules []moleculeRule) bool {
	one := big.NewRat(1, 1)
	for _, r := range rules {
		if new(big.Rat).Sub(c.total(r.to), c.total(r.from)).Cmp(one) != 0 {
			return false
		}
	}
	return true
}
func moleculeProofSmall(r *big.Rat) bool {
	return r.Num().BitLen() <= 4096 && r.Denom().BitLen() <= 4096
}
func deriveMoleculeCounter(rules []moleculeRule) *moleculeCounter {
	if len(rules) > 256 {
		return nil
	}
	present := [128]bool{}
	present['e'] = true
	for _, r := range rules {
		for _, text := range []string{r.from, r.to} {
			for _, b := range []byte(text) {
				present[b] = true
			}
		}
	}
	var alphabet []byte
	columns := [128]int{}
	for b, used := range present {
		if used {
			columns[b] = len(alphabet)
			alphabet = append(alphabet, byte(b))
		}
	}
	n := len(alphabet)
	if n > 64 {
		return nil
	}
	rows := make([][]big.Rat, len(rules))
	for i, r := range rules {
		rows[i] = make([]big.Rat, n+1)
		counts := make([]int64, n)
		for _, b := range []byte(r.from) {
			counts[columns[b]]--
		}
		for _, b := range []byte(r.to) {
			counts[columns[b]]++
		}
		for col, value := range counts {
			rows[i][col].SetInt64(value)
		}
		rows[i][n].SetInt64(1)
	}
	var pivots []int
	for col := 0; col < n; col++ {
		rank := len(pivots)
		pivot := rank
		for pivot < len(rows) && rows[pivot][col].Sign() == 0 {
			pivot++
		}
		if pivot == len(rows) {
			continue
		}
		rows[rank], rows[pivot] = rows[pivot], rows[rank]
		divisor := new(big.Rat).Set(&rows[rank][col])
		for k := col; k <= n; k++ {
			rows[rank][k].Quo(&rows[rank][k], divisor)
			if !moleculeProofSmall(&rows[rank][k]) {
				return nil
			}
		}
		for i := range rows {
			if i == rank || rows[i][col].Sign() == 0 {
				continue
			}
			factor := new(big.Rat).Set(&rows[i][col])
			for k := col; k <= n; k++ {
				term := new(big.Rat).Mul(factor, &rows[rank][k])
				rows[i][k].Sub(&rows[i][k], term)
				if !moleculeProofSmall(&rows[i][k]) {
					return nil
				}
			}
		}
		pivots = append(pivots, col)
	}
	for _, row := range rows {
		zero := true
		for col := 0; col < n; col++ {
			if row[col].Sign() != 0 {
				zero = false
				break
			}
		}
		if zero && row[n].Sign() != 0 {
			return nil
		}
	}
	counter := new(moleculeCounter)
	for i, col := range pivots {
		counter.weights[alphabet[col]].Set(&rows[i][n])
	}
	if !counter.verifies(rules) {
		return nil
	}
	return counter
}
func (c *moleculeCounter) predict(target string) (steps int, impossible bool, known bool) {
	difference := new(big.Rat).Sub(c.total(target), c.total("e"))
	if difference.Sign() < 0 || !difference.IsInt() {
		return 0, true, true
	}
	if !difference.Num().IsInt64() {
		return 0, false, false
	}
	value := difference.Num().Int64()
	if int64(int(value)) != value {
		return 0, false, false
	}
	return int(value), false, true
}

func moleculeWitness(target string, reverse []moleculeRule, required int, limits moleculeLimits) bool {
	if len(target) > limits.bytes {
		return false
	}
	checks, steps := limits.checks, 0
	for target != "e" && steps < required {
		changed := false
		for _, r := range reverse {
			if spendMoleculeChecks(&checks, target, r.to) != nil {
				return false
			}
			offset := strings.Index(target, r.to)
			if offset < 0 {
				continue
			}
			if len(target)-len(r.to)+len(r.from) > limits.bytes {
				return false
			}
			target = replaceMolecule(target, offset, r.to, r.from)
			steps++
			changed = true
			if target == "e" {
				return steps == required
			}
			if steps == required {
				return false
			}
		}
		if !changed {
			return false
		}
	}
	return target == "e" && steps == required
}
func exactMoleculeSearch(target string, reverse []moleculeRule, required *int, limits moleculeLimits) (int, error) {
	if limits.states < 1 || len(target) > limits.bytes {
		return 0, errMoleculeLimit
	}
	type state struct {
		text  string
		steps int
	}
	queue := []state{{target, 0}}
	seen := map[string]struct{}{target: {}}
	bytes, checks := len(target), limits.checks
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		queue[head] = state{}
		if required != nil && current.steps >= *required {
			continue
		}
		for _, r := range reverse {
			if err := spendMoleculeChecks(&checks, current.text, r.to); err != nil {
				return 0, err
			}
			for start := 0; start <= len(current.text)-len(r.to); {
				found := strings.Index(current.text[start:], r.to)
				if found < 0 {
					break
				}
				offset := start + found
				start = offset + 1
				length := len(current.text) - len(r.to) + len(r.from)
				if length > limits.bytes {
					return 0, errMoleculeLimit
				}
				next := replaceMolecule(current.text, offset, r.to, r.from)
				steps := current.steps + 1
				if next == "e" {
					return steps, nil
				}
				if _, ok := seen[next]; ok {
					continue
				}
				if len(seen) >= limits.states || length > limits.bytes-bytes {
					return 0, errMoleculeLimit
				}
				bytes += length
				seen[next] = struct{}{}
				queue = append(queue, state{next, steps})
			}
		}
	}
	return 0, errMoleculeUnreachable
}
func (p moleculeProblem) minimumSteps(limits moleculeLimits) (int, error) {
	if p.target == "e" {
		return 0, nil
	}
	hasStart := false
	for _, r := range p.rules {
		if r.from == "e" {
			hasStart = true
		}
	}
	if !hasStart {
		return 0, errMoleculeUnreachable
	}
	reverse := append([]moleculeRule(nil), p.rules...)
	sort.Slice(reverse, func(i, j int) bool {
		a, b := reverse[i], reverse[j]
		if len(a.to) != len(b.to) {
			return len(a.to) > len(b.to)
		}
		if a.to != b.to {
			return a.to < b.to
		}
		return a.from < b.from
	})
	var required *int
	if c := deriveMoleculeCounter(p.rules); c != nil {
		steps, impossible, known := c.predict(p.target)
		if impossible {
			return 0, errMoleculeUnreachable
		}
		if known {
			required = &steps
			if moleculeWitness(p.target, reverse, steps, limits) {
				return steps, nil
			}
		}
	}
	return exactMoleculeSearch(p.target, reverse, required, limits)
}
func (Day19) Part1(input string) (string, error) {
	p, err := parseMolecules(input)
	if err != nil {
		return "", err
	}
	distinct, err := p.oneStep(defaultMoleculeLimits)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(len(distinct)), nil
}
func (Day19) Part2(input string) (string, error) {
	p, err := parseMolecules(input)
	if err != nil {
		return "", err
	}
	steps, err := p.minimumSteps(defaultMoleculeLimits)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(steps), nil
}
func init() { solve.Register(Day19{}) }
