package solve2022

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func mixingSliceReference(values []int64, rounds int) []int {
	order := make([]int, len(values))
	for i := range order {
		order[i] = i
	}
	if len(values) < 2 {
		return order
	}
	for round := 0; round < rounds; round++ {
		for id, value := range values {
			at := 0
			for order[at] != id {
				at++
			}
			order = append(order[:at], order[at+1:]...)
			destination := new(big.Int).Add(big.NewInt(int64(at)), big.NewInt(value))
			destination.Mod(destination, big.NewInt(int64(len(order))))
			to := int(destination.Int64())
			order = append(order, 0)
			copy(order[to+1:], order[to:])
			order[to] = id
		}
	}
	return order
}

func checkMixingTree(t *testing.T, root *mixNode, nodes []mixNode) []int {
	t.Helper()
	ids := map[*mixNode]int{}
	for i := range nodes {
		ids[&nodes[i]] = i
	}
	seen := map[*mixNode]bool{}
	order := []int{}
	var walk func(*mixNode, *mixNode) int
	walk = func(node, parent *mixNode) int {
		if node == nil {
			return 0
		}
		if seen[node] {
			t.Fatal("cycle or repeated node")
		}
		seen[node] = true
		if node.parent != parent {
			t.Fatal("incorrect parent pointer")
		}
		if parent != nil && node.priority > parent.priority {
			t.Fatal("heap property violated")
		}
		left := walk(node.left, node)
		id, exists := ids[node]
		if !exists {
			t.Fatal("node identity lost")
		}
		order = append(order, id)
		right := walk(node.right, node)
		if node.size != left+right+1 {
			t.Fatal("incorrect subtree size")
		}
		return left + right + 1
	}
	walk(root, nil)
	for rank, id := range order {
		if mixRank(&nodes[id]) != rank || mixAt(root, rank) != &nodes[id] {
			t.Fatal("rank/select disagreement")
		}
	}
	return order
}

func TestMixingTreapAgainstExactSliceReference(t *testing.T) {
	combinations := 1
	for length := 1; length <= 6; length++ {
		combinations *= 5
		for code := 0; code < combinations; code++ {
			values := make([]int64, length)
			rest := code
			for i := range values {
				values[i] = int64(rest%5) - 2
				rest /= 5
			}
			for _, rounds := range []int{0, 1, 2, 10} {
				want := mixingSliceReference(values, rounds)
				root, nodes := mixTreap(values, rounds)
				got := checkMixingTree(t, root, nodes)
				if len(got) != len(want) {
					t.Fatal("lost nodes")
				}
				for i := range want {
					if got[i] != want[i] {
						t.Fatalf("values %v rounds %d: got %v want %v", values, rounds, got, want)
					}
				}
			}
		}
	}
}

func TestMixingSplitMergeAndExtremeMoves(t *testing.T) {
	values := []int64{math.MinInt64, math.MaxInt64, 0, 7, -4, 7, 11, -9, 2, -2, 8}
	for _, rounds := range []int{1, 3, 10} {
		root, nodes := mixTreap(values, rounds)
		want := mixingSliceReference(values, rounds)
		for cut := 0; cut <= len(values); cut++ {
			left, right := mixSplit(root, cut)
			if mixSize(left) != cut || mixSize(right) != len(values)-cut {
				t.Fatal("incorrect split sizes")
			}
			checkMixingTree(t, left, nodes)
			checkMixingTree(t, right, nodes)
			root = mixMerge(left, right)
			got := checkMixingTree(t, root, nodes)
			for i := range want {
				if got[i] != want[i] {
					t.Fatal("split/merge changed order")
				}
			}
		}
	}
	for span := 1; span <= 17; span++ {
		for position := 0; position <= span; position++ {
			for _, value := range []int64{math.MinInt64, math.MaxInt64, -1, 0, 1} {
				want := new(big.Int).Add(big.NewInt(int64(position)), big.NewInt(value))
				want.Mod(want, big.NewInt(int64(span)))
				if got := mixDestination(position, value, span); got != int(want.Int64()) {
					t.Fatalf("destination %d %d %d: %d want %s", position, value, span, got, want)
				}
			}
		}
	}
}

func TestMixingParserAndArithmetic(t *testing.T) {
	for _, solver := range []func(string) (string, error){(Day20{}).Part1, (Day20{}).Part2} {
		for _, input := range []string{"", "1", "0\n0", "0\n\n1", "0\nnope", "0\n1x", "9223372036854775808\n0"} {
			if _, err := solver(input); err == nil {
				t.Fatalf("accepted %q", input)
			}
		}
		for _, input := range []string{"0", " 0 \r\n", "-0\n+2\n-2\n"} {
			if _, err := solver(input); err != nil {
				t.Fatalf("rejected valid signed input: %v", err)
			}
		}
		if got, err := solver("0"); err != nil || got != "0" {
			t.Fatalf("singleton zero: %q %v", got, err)
		}
	}
	key := int64(811589153)
	for _, value := range []int64{math.MaxInt64 / key, math.MinInt64 / key} {
		if _, _, err := parseMixing(fmt.Sprintf("0\n%d", value), key); err != nil {
			t.Fatal("rejected safe scaling boundary")
		}
	}
	for _, value := range []int64{math.MaxInt64/key + 1, math.MinInt64/key - 1} {
		if _, err := (Day20{}).Part2(fmt.Sprintf("0\n%d", value)); err == nil {
			t.Fatal("accepted scaling overflow")
		}
	}
	if _, err := (Day20{}).Part1("9223372036854775807\n9223372036854775807\n0"); err == nil {
		t.Fatal("accepted coordinate-sum overflow")
	}
	input := strings.Join([]string{"0", "1", "2", "3", strconv.FormatInt(math.MinInt64, 10), strconv.FormatInt(math.MaxInt64, 10), strconv.FormatInt(math.MaxInt64, 10)}, "\n")
	got, err := groveCoordinates(input, 1, 0)
	if err != nil || got != strconv.FormatInt(math.MaxInt64-1, 10) {
		t.Fatalf("exact cancelling sum: %q %v", got, err)
	}
}
