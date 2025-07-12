package solve2018

import (
	"aoc/solve"
	"fmt"
	"strconv"
)

type Day9 struct{}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 9}
}

// Node represents a marble in the circular linked list
type Node struct {
	value int
	next  *Node
	prev  *Node
}

// CircularList represents the marble circle
type CircularList struct {
	current *Node
	size    int
}

// NewCircularList creates a new circular list with the initial marble 0
func NewCircularList() *CircularList {
	node := &Node{value: 0}
	node.next = node
	node.prev = node
	return &CircularList{current: node, size: 1}
}

// Insert adds a marble between 1 and 2 positions clockwise from current
func (cl *CircularList) Insert(value int) {
	// Move 1 position clockwise
	pos := cl.current.next

	// Create new node
	newNode := &Node{value: value}

	// Insert between pos and pos.next
	newNode.next = pos.next
	newNode.prev = pos
	pos.next.prev = newNode
	pos.next = newNode

	// New marble becomes current
	cl.current = newNode
	cl.size++
}

// Remove removes the marble 7 positions counter-clockwise and returns its value
func (cl *CircularList) Remove() int {
	// Move 7 positions counter-clockwise
	toRemove := cl.current
	for i := 0; i < 7; i++ {
		toRemove = toRemove.prev
	}

	// Remove the node
	toRemove.prev.next = toRemove.next
	toRemove.next.prev = toRemove.prev

	// Set current to the node clockwise of removed
	cl.current = toRemove.next
	cl.size--

	return toRemove.value
}

func (d Day9) Part1(input string) (string, error) {
	var players, lastMarble int
	_, err := fmt.Sscanf(input, "%d players; last marble is worth %d points", &players, &lastMarble)
	if err != nil {
		return "", err
	}

	score := playMarbleGame(players, lastMarble)
	return strconv.Itoa(score), nil
}

func (d Day9) Part2(input string) (string, error) {
	var players, lastMarble int
	_, err := fmt.Sscanf(input, "%d players; last marble is worth %d points", &players, &lastMarble)
	if err != nil {
		return "", err
	}

	// For part 2, multiply the last marble by 100
	score := playMarbleGame(players, lastMarble*100)
	return strconv.Itoa(score), nil
}

func playMarbleGame(players, lastMarble int) int {
	scores := make([]int, players)
	circle := NewCircularList()

	for marble := 1; marble <= lastMarble; marble++ {
		currentPlayer := (marble - 1) % players

		if marble%23 == 0 {
			// Special case: marble is multiple of 23
			scores[currentPlayer] += marble
			scores[currentPlayer] += circle.Remove()
		} else {
			// Normal case: insert marble
			circle.Insert(marble)
		}
	}

	// Find the highest score
	maxScore := 0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	return maxScore
}

func init() {
	solve.Register(Day9{})
}
