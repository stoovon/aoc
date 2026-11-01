package solve2022

import (
	"fmt"
	"math"
	"math/big"
	"math/rand/v2"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day20 struct{}

func (d Day20) Coords() solve.SolutionCoords { return solve.SolutionCoords{Year: 2022, Day: 20} }

type mixNode struct {
	value               int64
	priority            uint64
	size                int
	left, right, parent *mixNode
}

func mixSize(node *mixNode) int {
	if node == nil {
		return 0
	}
	return node.size
}
func mixUpdate(node *mixNode) {
	node.size = 1 + mixSize(node.left) + mixSize(node.right)
	if node.left != nil {
		node.left.parent = node
	}
	if node.right != nil {
		node.right.parent = node
	}
	node.parent = nil
}

func mixMerge(left, right *mixNode) *mixNode {
	if left == nil {
		if right != nil {
			right.parent = nil
		}
		return right
	}
	if right == nil {
		left.parent = nil
		return left
	}
	if left.priority >= right.priority {
		left.right = mixMerge(left.right, right)
		mixUpdate(left)
		return left
	}
	right.left = mixMerge(left, right.left)
	mixUpdate(right)
	return right
}

func mixSplit(root *mixNode, count int) (*mixNode, *mixNode) {
	if root == nil {
		return nil, nil
	}
	leftSize := mixSize(root.left)
	if count <= leftSize {
		left, rest := mixSplit(root.left, count)
		root.left = rest
		mixUpdate(root)
		return left, root
	}
	rest, right := mixSplit(root.right, count-leftSize-1)
	root.right = rest
	mixUpdate(root)
	return root, right
}

func mixRank(node *mixNode) int {
	rank := mixSize(node.left)
	for node.parent != nil {
		parent := node.parent
		if node == parent.right {
			rank += mixSize(parent.left) + 1
		}
		node = parent
	}
	return rank
}

func mixAt(root *mixNode, index int) *mixNode {
	for root != nil {
		left := mixSize(root.left)
		if index < left {
			root = root.left
		} else if index == left {
			return root
		} else {
			index -= left + 1
			root = root.right
		}
	}
	return nil
}

func mixDestination(position int, value int64, span int) int {
	position %= span
	shift := int(value % int64(span))
	if shift < 0 {
		shift += span
	}
	if position >= span-shift {
		return position - (span - shift)
	}
	return position + shift
}

func mixTreap(values []int64, rounds int) (*mixNode, []mixNode) {
	nodes := make([]mixNode, len(values))
	random := rand.New(rand.NewPCG(0x6a09e667f3bcc909, 0xbb67ae8584caa73b))
	var root *mixNode
	for i, value := range values {
		nodes[i] = mixNode{value: value, priority: random.Uint64(), size: 1}
		root = mixMerge(root, &nodes[i])
	}
	if len(nodes) < 2 {
		return root, nodes
	}
	for round := 0; round < rounds; round++ {
		for i := range nodes {
			node := &nodes[i]
			position := mixRank(node)
			before, rest := mixSplit(root, position)
			removed, after := mixSplit(rest, 1)
			root = mixMerge(before, after)
			destination := mixDestination(position, node.value, len(nodes)-1)
			before, after = mixSplit(root, destination)
			root = mixMerge(mixMerge(before, removed), after)
		}
	}
	return root, nodes
}

func parseMixing(input string, key int64) ([]int64, int, error) {
	values := []int64{}
	zero := -1
	for line, record := range strings.Split(strings.TrimSuffix(input, "\n"), "\n") {
		value, err := strconv.ParseInt(strings.TrimSpace(record), 10, 64)
		if err != nil {
			return nil, 0, fmt.Errorf("line %d: invalid signed decimal integer", line+1)
		}
		if value > math.MaxInt64/key || value < math.MinInt64/key {
			return nil, 0, fmt.Errorf("line %d: scaled value exceeds int64", line+1)
		}
		if value == 0 {
			if zero >= 0 {
				return nil, 0, fmt.Errorf("line %d: repeated zero anchor", line+1)
			}
			zero = len(values)
		}
		values = append(values, value*key)
	}
	if zero < 0 {
		return nil, 0, fmt.Errorf("missing zero anchor")
	}
	return values, zero, nil
}

func groveCoordinates(input string, key int64, rounds int) (string, error) {
	values, zero, err := parseMixing(input, key)
	if err != nil {
		return "", err
	}
	root, nodes := mixTreap(values, rounds)
	position := mixRank(&nodes[zero])
	var sum, term big.Int
	for _, offset := range []int64{1000, 2000, 3000} {
		id := mixDestination(position, offset, len(nodes))
		term.SetInt64(mixAt(root, id).value)
		sum.Add(&sum, &term)
	}
	if !sum.IsInt64() {
		return "", fmt.Errorf("grove coordinate sum exceeds int64")
	}
	return strconv.FormatInt(sum.Int64(), 10), nil
}
func (d Day20) Part1(input string) (string, error) { return groveCoordinates(input, 1, 1) }
func (d Day20) Part2(input string) (string, error) { return groveCoordinates(input, 811589153, 10) }
func init()                                        { solve.Register(Day20{}) }
