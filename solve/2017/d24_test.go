package solve2017

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

type referenceBridgePiece struct{ a, b int }
type referenceBridgeResult struct{ strongest, length, longestStrength int }

func referenceBridge(pieces []referenceBridgePiece, open int) referenceBridgeResult {
	var best referenceBridgeResult
	for index, piece := range pieces {
		if piece.a != open && piece.b != open {
			continue
		}
		next := piece.a
		if next == open {
			next = piece.b
		}
		remaining := make([]referenceBridgePiece, 0, len(pieces)-1)
		remaining = append(remaining, pieces[:index]...)
		remaining = append(remaining, pieces[index+1:]...)
		tail := referenceBridge(remaining, next)
		weight := piece.a + piece.b
		if candidate := weight + tail.strongest; candidate > best.strongest {
			best.strongest = candidate
		}
		length, strength := tail.length+1, tail.longestStrength+weight
		if length > best.length || (length == best.length && strength > best.longestStrength) {
			best.length = length
			best.longestStrength = strength
		}
	}
	return best
}

func TestBridgeIndexedCandidatesAgainstRemovingPieces(t *testing.T) {
	options := []referenceBridgePiece{{0, 0}, {0, 1}, {1, 0}, {1, 1}, {1, 2}, {2, 2}}
	for size := 0; size <= 5; size++ {
		count := 1
		for i := 0; i < size; i++ {
			count *= len(options)
		}
		for encoding := 0; encoding < count; encoding++ {
			rest := encoding
			pieces := make([]referenceBridgePiece, size)
			comps := make([]component, size)
			for i := range pieces {
				pieces[i] = options[rest%len(options)]
				rest /= len(options)
				comps[i] = component{a: pieces[i].a, b: pieces[i].b}
			}
			want := referenceBridge(pieces, 0)
			index := bridgeCandidates(comps)
			got := maxBridgeStrength(comps, index, 0)
			length, strength := maxLongestBridge(comps, index, 0)
			if got != want.strongest || length != want.length || strength != want.longestStrength {
				t.Fatalf("pieces%v: %d/%d/%d want%v", pieces, got, length, strength, want)
			}
			for _, piece := range comps {
				if piece.used {
					t.Fatal("search left a component marked used")
				}
			}
		}
	}
}

func TestBridgeSelfLoopsAndPhysicalDuplicates(t *testing.T) {
	comps := []component{{a: 0, b: 4}, {a: 4, b: 4}, {a: 4, b: 4}}
	if got := len(bridgeCandidates(comps)[4]); got != 3 {
		t.Fatalf("port4 has%d candidates want3", got)
	}
	for _, fn := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
		if got, err := fn("0/4\n4/4\n4/4"); err != nil || got != "20" {
			t.Fatalf("physical duplicates: %q %v", got, err)
		}
	}
}

func TestBridgeMoreThanOneMachineWordOfPieces(t *testing.T) {
	var input strings.Builder
	for port := 0; port < 70; port++ {
		if port > 0 {
			input.WriteByte('\n')
		}
		input.WriteString(strconv.Itoa(port))
		input.WriteByte('/')
		input.WriteString(strconv.Itoa(port + 1))
	}
	for _, fn := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
		if got, err := fn(input.String()); err != nil || got != "4900" {
			t.Fatalf("70-piece chain: %q %v", got, err)
		}
	}
}

func TestBridgeScoringPoliciesStaySeparate(t *testing.T) {
	for _, tc := range []struct{ input, p1, p2 string }{
		{"0/50\n0/1\n1/2\n2/3", "50", "9"},
		{"0/2\n2/2\n0/3\n3/3", "9", "9"},
		{"3/4\n4/4", "0", "0"},
		{"2/0\n3/2", "7", "7"},
		{"0/0\n0/0", "0", "0"},
	} {
		for part, item := range []struct {
			fn   func(string) (string, error)
			want string
		}{{(Day24{}).Part1, tc.p1}, {(Day24{}).Part2, tc.p2}} {
			if got, err := item.fn(tc.input); err != nil || got != item.want {
				t.Fatalf("part%d input%q: %q %v want%q", part+1, tc.input, got, err, item.want)
			}
		}
	}
}

func TestBridgeCheckedParsingAndStrengthBound(t *testing.T) {
	max := strconv.Itoa(math.MaxInt)
	for _, input := range []string{"", " ", "0", "0/", "/1", "0/1/2", "0/x", "-1/2", "+1/2", "0/1\n\n1/2", max + "/1", max + "/0\n0/1", max + "0/0"} {
		for _, fn := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
			if _, err := fn(input); err == nil {
				t.Fatalf("accepted malformed or out-of-range input%q", input)
			}
		}
	}
	for _, input := range []string{"0/" + max, " 0 / 4 \r\n4 / 4\r\n"} {
		want := "12"
		if strings.HasPrefix(input, "0/") {
			want = max
		}
		for _, fn := range []func(string) (string, error){(Day24{}).Part1, (Day24{}).Part2} {
			if got, err := fn(input); err != nil || got != want {
				t.Fatalf("valid bounded input: %q %v want%q", got, err, want)
			}
		}
	}
}
