package solve2018

import (
	"fmt"
	"strings"
	"testing"
)

func TestChronalPart1ReturnsFirstCandidate(t *testing.T) {
	program := "#ip 5\nseti 73 0 3\neqrr 3 0 1\naddr 1 5 5\nseti -1 0 5"
	got, err := (Day21{}).Part1(program)
	if err != nil || got != "73" {
		t.Fatalf("first candidate: %q %v", got, err)
	}
}

func TestChronalMachineOperations(t *testing.T) {
	for _, tc := range []struct {
		op         string
		a, b, want int
	}{
		{"addr", 0, 1, 10}, {"addi", 0, 4, 11}, {"mulr", 0, 1, 21}, {"muli", 1, 4, 12},
		{"banr", 0, 1, 3}, {"bani", 0, 5, 5}, {"borr", 0, 1, 7}, {"bori", 1, 8, 11},
		{"setr", 1, 99, 3}, {"seti", 13, 99, 13}, {"gtir", 8, 0, 1}, {"gtir", 7, 0, 0},
		{"gtri", 0, 6, 1}, {"gtri", 0, 7, 0}, {"gtrr", 0, 1, 1}, {"gtrr", 1, 0, 0},
		{"eqir", 7, 0, 1}, {"eqir", 8, 0, 0}, {"eqri", 1, 3, 1}, {"eqri", 1, 4, 0},
		{"eqrr", 0, 0, 1}, {"eqrr", 0, 1, 0},
	} {
		machine := chronalMachine{instructions: [][4]int{{opCode(tc.op), tc.a, tc.b, 2}}, ipRegister: 5, registers: [6]int{7, 3}}
		machine.step()
		if machine.registers[2] != tc.want || machine.ip != 1 {
			t.Fatalf("%s: result %d IP%d want %d", tc.op, machine.registers[2], machine.ip, tc.want)
		}
	}
}

func TestChronalBoundIPIsSeparateFromRegister(t *testing.T) {
	machine := chronalMachine{instructions: [][4]int{{opCode("seti"), 41, 0, 5}}, ipRegister: 5}
	machine.step()
	if machine.registers[5] != 41 || machine.ip != 42 {
		t.Fatalf("register %d IP%d", machine.registers[5], machine.ip)
	}
	if machine.running() {
		t.Fatal("out-of-range IP still running")
	}
}

func chronalSyntheticDivision(start, q, scratch, value, ip, divisor int) [][4]int {
	code := make([][4]int, start)
	return append(code,
		[4]int{opCode("seti"), 0, 91, q},
		[4]int{opCode("addi"), q, 1, scratch},
		[4]int{opCode("muli"), scratch, divisor, scratch},
		[4]int{opCode("gtrr"), scratch, value, scratch},
		[4]int{opCode("addr"), scratch, ip, ip},
		[4]int{opCode("addi"), ip, 1, ip},
		[4]int{opCode("seti"), start + 8, 37, ip},
		[4]int{opCode("addi"), q, 1, q},
		[4]int{opCode("seti"), start, 53, ip},
	)
}

func TestChronalDivisionMatchesCompleteMachineState(t *testing.T) {
	for _, roles := range [][4]int{{1, 2, 3, 5}, {4, 0, 2, 1}, {5, 4, 1, 2}} {
		q, scratch, value, ip := roles[0], roles[1], roles[2], roles[3]
		for _, start := range []int{0, 3} {
			for _, divisor := range []int{1, 2, 5, 11} {
				code := chronalSyntheticDivision(start, q, scratch, value, ip, divisor)
				shortcuts := findChronalDivisions(code, ip)
				if len(shortcuts) != 1 {
					t.Fatal("synthetic division not recognized")
				}
				for input := 0; input < 128; input++ {
					slow := chronalMachine{instructions: code, ipRegister: ip, ip: start, registers: [6]int{-5, 17, 33, -11, 22, 123}}
					slow.registers[value] = input
					fast := slow
					fast.divisions = shortcuts
					steps := 0
					for slow.ip != start+9 && steps < 2000 {
						slow.step()
						steps++
					}
					fast.step()
					if slow.ip != fast.ip || slow.registers != fast.registers {
						t.Fatalf("roles %v start%d divisor%d input%d: slow IP%d %v fast IP%d %v", roles, start, divisor, input, slow.ip, slow.registers, fast.ip, fast.registers)
					}
				}
			}
		}
	}
}

func TestChronalDivisionRejectsNearMatches(t *testing.T) {
	code := chronalSyntheticDivision(0, 1, 2, 3, 5, 7)
	for index := range code {
		changed := append([][4]int(nil), code...)
		changed[index][0] = (changed[index][0] + 1) % 16
		if len(findChronalDivisions(changed, 5)) != 0 {
			t.Fatalf("accepted changed opcode at%d", index)
		}
	}
	for _, changed := range [][][4]int{
		chronalSyntheticDivision(0, 1, 2, 3, 5, 0),
		chronalSyntheticDivision(0, 1, 2, 3, 5, -3),
		chronalSyntheticDivision(0, 1, 1, 3, 5, 7),
		chronalSyntheticDivision(0, 1, 2, 3, 3, 7),
		chronalSyntheticDivision(0, 1, 2, 7, 5, 7),
		code[:8],
	} {
		if len(findChronalDivisions(changed, 5)) != 0 {
			t.Fatal("accepted unsafe or incomplete loop")
		}
	}
	changed := append([][4]int(nil), code...)
	changed[6][1]++
	if len(findChronalDivisions(changed, 5)) != 0 {
		t.Fatal("accepted wrong exit jump")
	}
	changed = append([][4]int(nil), code...)
	changed[8][1]++
	if len(findChronalDivisions(changed, 5)) != 0 {
		t.Fatal("accepted wrong loop jump")
	}
}

func TestChronalDivisionGuardsArithmeticDomain(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	division := chronalDivision{1, 2, 3, 5}
	for _, value := range []int{-1, maxInt} {
		machine := chronalMachine{ipRegister: 4, registers: [6]int{11, 12, 13, value, 15, 16}}
		before := machine
		if division.apply(&machine) {
			t.Fatalf("unsafe shortcut for %d", value)
		}
		if machine.ip != before.ip || machine.registers != before.registers {
			t.Fatal("rejected shortcut changed machine")
		}
	}
	machine := chronalMachine{ipRegister: 4, registers: [6]int{11, 12, 13, maxInt - 5, 15, 16}}
	if !division.apply(&machine) {
		t.Fatal("safe large quotient rejected")
	}
}

func chronalSyntheticCandidates(values [4]int) string {
	lines := []string{"#ip 5", fmt.Sprintf("seti %d 0 3", values[0])}
	for counter := 1; counter < 4; counter++ {
		lines = append(lines, fmt.Sprintf("eqri 1 %d 2", counter), fmt.Sprintf("muli 2 %d 2", values[counter]-values[0]), "addr 3 2 3")
	}
	lines = append(lines, "seti 27 0 5")
	for len(lines) < 29 {
		lines = append(lines, "seti 0 0 4")
	}
	lines = append(lines, "eqrr 3 0 2", "addr 2 5 5", "seti 31 0 5", "seti 99 0 5", "addi 1 1 1", "bani 1 3 1", "seti -1 0 5")
	return strings.Join(lines, "\n")
}

func TestChronalCycleNeedsFullStateAndLatestNewCandidate(t *testing.T) {
	for _, tc := range []struct {
		values [4]int
		want   string
	}{
		{[4]int{7, 11, 7, 19}, "19"},
		{[4]int{7, 11, 19, 7}, "19"},
		{[4]int{0, 11, 0, 19}, "19"},
		{[4]int{0, 0, 0, 0}, "0"},
	} {
		program := chronalSyntheticCandidates(tc.values)
		got, err := (Day21{}).Part2(program)
		if err != nil || got != tc.want {
			t.Fatalf("candidates %v: %q %v want %q", tc.values, got, err, tc.want)
		}
		got, err = (Day21{}).Part1(program)
		if err != nil || got != fmt.Sprint(tc.values[0]) {
			t.Fatalf("first candidate %v: %q %v", tc.values, got, err)
		}
	}
}

func chronalSimpleProbe(padding, ip, candidate, result int, reverse bool) string {
	lines := []string{fmt.Sprintf("#ip %d", ip), fmt.Sprintf("seti 73 91 %d", candidate)}
	for i := 0; i < padding; i++ {
		lines = append(lines, fmt.Sprintf("seti 0 37 %d", result))
	}
	a, b := candidate, 0
	if reverse {
		a, b = b, a
	}
	lines = append(lines, fmt.Sprintf("eqrr %d %d %d", a, b, result), fmt.Sprintf("addr %d %d %d", result, ip, ip), fmt.Sprintf("seti -1 53 %d", ip))
	return strings.Join(lines, "\n")
}

func TestChronalFindsRelocatedRenamedHaltingProbe(t *testing.T) {
	for _, padding := range []int{0, 4, 9} {
		for _, roles := range [][3]int{{5, 3, 2}, {1, 4, 3}, {4, 2, 5}} {
			for _, reverse := range []bool{false, true} {
				program := chronalSimpleProbe(padding, roles[0], roles[1], roles[2], reverse)
				for part, fn := range []func(string) (string, error){(Day21{}).Part1, (Day21{}).Part2} {
					got, err := fn(program)
					if err != nil || got != "73" {
						t.Fatalf("padding%d roles%v reverse%t part%d: %q %v", padding, roles, reverse, part+1, got, err)
					}
				}
			}
		}
	}
	// Equality may leave through an explicit unconditional jump as well.
	program := chronalSimpleProbe(0, 5, 3, 2, false) + "\nseti 87 0 5"
	if got, err := (Day21{}).Part1(program); err != nil || got != "73" {
		t.Fatalf("explicit exit: %q %v", got, err)
	}
}

func TestChronalRejectsUnprovenHaltingComparisons(t *testing.T) {
	base := chronalSimpleProbe(0, 5, 3, 2, false)
	for _, program := range []string{
		strings.Replace(base, "eqrr 3 0 2", "eqrr 3 1 2", 1),
		strings.Replace(base, "eqrr 3 0 2", "eqrr 0 0 2", 1),
		strings.Replace(base, "eqrr 3 0 2", "eqrr 5 0 2", 1),
		strings.Replace(base, "eqrr 3 0 2", "eqrr 3 0 0", 1),
		strings.Replace(base, "eqrr 3 0 2", "eqrr 3 0 5", 1),
		strings.Replace(base, "addr 2 5 5", "addi 5 1 5", 1),
		strings.Replace(base, "seti -1 53 5", "seti 87 53 5", 1),
		base + "\nseti -1 0 5",
		base + "\naddi 3 1 3",
		base + "\neqrr 3 0 2",
		base + "\nseti 1 0 0",
		base + "\naddi 0 1 1",
		"#ip 0\nseti 1 0 1",
	} {
		for _, fn := range []func(string) (string, error){(Day21{}).Part1, (Day21{}).Part2} {
			if _, err := fn(program); err == nil {
				t.Fatalf("accepted unsupported halt layout %q", program)
			}
		}
	}
}

func TestChronalParserRejectsMalformedInstructions(t *testing.T) {
	for _, program := range []string{
		"", "#ip", "#ip 1 extra", "#ip x", "#ip -1", "#ip 6", "#ip 5", "seti 1 0 1",
		"#ip 5\nunknown 1 0 1", "#ip 5\nseti 1 0", "#ip 5\nseti 1 0 1 extra",
		"#ip 5\nseti x 0 1", "#ip 5\nseti 1 0 6", "#ip 5\nseti 1 0 -1",
		"#ip 5\naddr 6 0 1", "#ip 5\ngtir 7 -1 1", "#ip 5\nseti 1 0 1\n\nseti 2 0 1",
		"#ip 5\nseti 99999999999999999999999999 0 1",
	} {
		if _, _, err := parseInstructions(program); err == nil {
			t.Fatalf("accepted malformed program %q", program)
		}
	}
	valid := chronalSimpleProbe(0, 5, 3, 2, false)
	if _, err := newChronalMachine(strings.ReplaceAll(valid, "\n", "\r\n") + "\r\n"); err != nil {
		t.Fatalf("CRLF: %v", err)
	}
}
