package solve2015

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestMoleculeLiteralRules(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
	}{
		{"Q => QR\nR => QR\n\nQR", "2"}, {"QQ => R\n\nQQQ", "2"}, {"A => AA\n\nAA", "1"}, {"Q => XYZ\n\nQ7R", "1"}, {"e => Q\ne => R\nQ => Qe\n\nQe", "3"},
	} {
		got, err := (Day19{}).Part1(tc.input)
		if err != nil || got != tc.want {
			t.Fatalf("one-step replacements: %s %v", got, err)
		}
	}
	p, _ := parseMolecules("QQ => R\n\nQQQ")
	results, _ := p.oneStep(defaultMoleculeLimits)
	if len(results) != 2 {
		t.Fatal("lost overlap")
	}
	for _, want := range []string{"RQ", "QR"} {
		if _, ok := results[want]; !ok {
			t.Fatal("wrong replacement")
		}
	}
}
func TestMoleculeShortestPaths(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"B => A\ne => B\ne => A\n\nA", "1"},
		{"e => C\nC => AB\nA => AB\n\nAB", "2"},
		{"e => Q\nQ => QR\nR => SS\n\nQRR", "3"},
		{"e => Q\ne => R\nQ => Qe\n\nQQ", "3"},
		{"e => A\nA => A\n\nA", "1"}, {"A => A\n\ne", "0"},
	} {
		got, err := (Day19{}).Part2(tc.input)
		if err != nil || got != tc.want {
			t.Fatalf("shortest path: got %s %v, want %s", got, err, tc.want)
		}
	}
}
func TestMoleculeCounterRequiresWitness(t *testing.T) {
	p, _ := parseMolecules("e => AA\n\nAA")
	c := deriveMoleculeCounter(p.rules)
	if c == nil || !c.verifies(p.rules) {
		t.Fatal("missing exact counter")
	}
	if steps, impossible, known := c.predict("AA"); !known || impossible || steps != 1 {
		t.Fatal("wrong prediction")
	}
	c.weights['A'].SetInt64(9)
	if c.verifies(p.rules) {
		t.Fatal("unchecked counter accepted")
	}
	for _, input := range []string{"e => AA\n\nA", "e => AA\n\nBB", "A => AA\n\nAA"} {
		_, err := (Day19{}).Part2(input)
		if !errors.Is(err, errMoleculeUnreachable) {
			t.Fatalf("counter mistaken for reachability: %v", err)
		}
	}
	p, _ = parseMolecules("e => A\ne => B\nB => A\n\nA")
	if deriveMoleculeCounter(p.rules) != nil {
		t.Fatal("variable length paths cannot have unit-increment proof")
	}
	p, _ = parseMolecules("e => C\nC => AB\nA => AB\n\nAB")
	reverse := []moleculeRule{p.rules[2], p.rules[1], p.rules[0]}
	if moleculeWitness(p.target, reverse, 2, defaultMoleculeLimits) {
		t.Fatal("stuck candidate must not supply a step count")
	}
	var large strings.Builder
	for i := 0; i < 257; i++ {
		fmt.Fprintf(&large, "e => A%d\n", i)
	}
	large.WriteString("\nA0")
	p, _ = parseMolecules(large.String())
	if deriveMoleculeCounter(p.rules) != nil {
		t.Fatal("proof rule budget ignored")
	}
	if got, err := p.minimumSteps(defaultMoleculeLimits); err != nil || got != 1 {
		t.Fatal("proof refusal must permit exact fallback")
	}
}
func TestMoleculeBudgetsAndCycles(t *testing.T) {
	p, _ := parseMolecules("e => A\nBB => B\nB => B\n\nB")
	if _, err := p.minimumSteps(moleculeLimits{8, 128, 1000}); !errors.Is(err, errMoleculeLimit) {
		t.Fatalf("growth not bounded: %v", err)
	}
	p, _ = parseMolecules("e => A\nA => A\n\nA")
	for _, limits := range []moleculeLimits{{0, 128, 1000}, {8, 0, 1000}, {8, 128, 0}} {
		if _, err := p.minimumSteps(limits); !errors.Is(err, errMoleculeLimit) {
			t.Fatal("budget exhaustion is not unreachability")
		}
	}
	zero, _ := parseMolecules("A => A\n\ne")
	if steps, err := zero.minimumSteps(moleculeLimits{}); err != nil || steps != 0 {
		t.Fatal("zero steps needs no search")
	}
	p, _ = parseMolecules("e => Q\n\nQQ")
	if _, err := p.oneStep(moleculeLimits{}); !errors.Is(err, errMoleculeLimit) {
		t.Fatal("one-step generation budget ignored")
	}
}
func TestMoleculeParser(t *testing.T) {
	for _, input := range []string{"", "\nA", "e => A", "e => A\n\n", "e => A\n\nA\nB", "e => A\n\nA\n\n", "e => A\n\n A", "e => A\n\nÅ", " => A\n\nA", "e => \n\nA", "e->A\n\nA", "e => A => B\n\nA", "e => A\ne => A\n\nA", "e => A\n\n\nA"} {
		if _, err := parseMolecules(input); err == nil {
			t.Fatalf("accepted malformed input %q", input)
		}
	}
	got, err := (Day19{}).Part2("e => A\r\nA => AB\r\n\r\nAB\r\n")
	if err != nil || got != "2" {
		t.Fatal("CRLF parse failed")
	}
	for _, part := range []func(string) (string, error){(Day19{}).Part1, (Day19{}).Part2} {
		if _, err := part("e => A\nbad\n\nA"); err == nil || !strings.HasPrefix(err.Error(), "line 2:") {
			t.Fatal("missing line context")
		}
	}
}

// All generated forward rules preserve or increase length, so the target-length
// bound makes this graph complete. It shares no reverse-search or counter code.
func moleculeForwardReference(rules []moleculeRule, target string) (int, bool) {
	type entry struct {
		text  string
		steps int
	}
	queue := []entry{{"e", 0}}
	seen := map[string]bool{"e": true}
	for head := 0; head < len(queue); head++ {
		s := queue[head]
		if s.text == target {
			return s.steps, true
		}
		for _, r := range rules {
			for pos := 0; pos+len(r.from) <= len(s.text); pos++ {
				if s.text[pos:pos+len(r.from)] != r.from {
					continue
				}
				next := s.text[:pos] + r.to + s.text[pos+len(r.from):]
				if len(next) <= len(target) && !seen[next] {
					seen[next] = true
					queue = append(queue, entry{next, s.steps + 1})
				}
			}
		}
	}
	return 0, false
}
func TestMoleculeAgainstCompleteForwardGraphs(t *testing.T) {
	pool := []moleculeRule{{"e", "A"}, {"e", "B"}, {"A", "B"}, {"B", "A"}, {"A", "AA"}, {"A", "AB"}, {"B", "BA"}, {"B", "BB"}}
	targets := []string{"A", "B", "AA", "AB", "BA", "BB", "AAA", "AAB", "ABA", "ABB", "BAA", "BAB", "BBA", "BBB"}
	cases := 0
	for mask := 1; mask < 1<<len(pool); mask++ {
		var rules []moleculeRule
		for i, r := range pool {
			if mask&(1<<i) != 0 {
				rules = append(rules, r)
			}
		}
		for _, target := range targets {
			want, reachable := moleculeForwardReference(rules, target)
			got, err := (moleculeProblem{rules, target}).minimumSteps(moleculeLimits{128, 4096, 10000})
			if reachable {
				if err != nil || got != want {
					t.Fatalf("rules %d target %s: %d %v vs %d", mask, target, got, err, want)
				}
			} else if !errors.Is(err, errMoleculeUnreachable) {
				t.Fatalf("rules %d target %s: %v", mask, target, err)
			}
			cases++
		}
	}
	if cases != 3570 {
		t.Fatal("incomplete oracle corpus")
	}
}
