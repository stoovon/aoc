package solve2017

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day7 struct {
}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2017, Day: 7}
}

func (d Day7) Part1(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	all := make(map[string]struct{})
	children := make(map[string]struct{})

	for _, line := range lines {
		parts := strings.Fields(line)
		name := parts[0]
		all[name] = struct{}{}
		if len(parts) > 2 && parts[2] == "->" {
			for _, child := range parts[3:] {
				child = strings.TrimRight(child, ",")
				children[child] = struct{}{}
			}
		}
	}

	for name := range all {
		if _, ok := children[name]; !ok {
			return name, nil // This is the bottom program
		}
	}
	return "", errors.New("no bottom program found")
}

type node struct {
	name     string
	weight   int
	children []*node
	parent   *node
}

func (d Day7) Part2(input string) (string, error) {
	nodes := make(map[string]*node)
	childLinks := make(map[string][]string)

	// Parse input and build nodes
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		parts := strings.Fields(line)
		name := parts[0]
		weight, _ := strconv.Atoi(strings.Trim(parts[1], "()"))
		n, ok := nodes[name]
		if !ok {
			n = &node{name: name}
			nodes[name] = n
		}
		n.weight = weight
		if len(parts) > 2 && parts[2] == "->" {
			for _, child := range parts[3:] {
				child = strings.TrimRight(child, ",")
				childLinks[name] = append(childLinks[name], child)
			}
		}
	}

	// Link children and parents
	for parent, kids := range childLinks {
		for _, child := range kids {
			cn, ok := nodes[child]
			if !ok {
				cn = &node{name: child}
				nodes[child] = cn
			}
			nodes[parent].children = append(nodes[parent].children, cn)
			cn.parent = nodes[parent]
		}
	}

	// Find root
	var root *node
	for _, n := range nodes {
		if n.parent == nil {
			root = n
			break
		}
	}

	// Recursively find the unbalanced node
	_, answer, found := findUnbalanced(root)
	if !found {
		return "", errors.New("no unbalanced node found")
	}
	return strconv.Itoa(answer), nil
}

// Returns (totalWeight, correctedWeight, found)
func findUnbalanced(n *node) (int, int, bool) {
	if len(n.children) == 0 {
		return n.weight, 0, false
	}

	weights := make(map[int][]*node)
	childWeights := make([]int, len(n.children))
	for i, c := range n.children {
		w, corrected, found := findUnbalanced(c)
		if found {
			return 0, corrected, true
		}
		childWeights[i] = w
		weights[w] = append(weights[w], c)
	}

	if len(weights) == 1 {
		sum := n.weight
		for _, w := range childWeights {
			sum += w
		}
		return sum, 0, false
	}

	// Find the outlier
	var wrongWeight, correctWeight int
	for w, nodes := range weights {
		if len(nodes) == 1 {
			wrongWeight = w
		} else {
			correctWeight = w
		}
	}
	diff := correctWeight - wrongWeight
	badNode := weights[wrongWeight][0]
	return 0, badNode.weight + diff, true
}

func init() {
	solve.Register(Day7{})
}
