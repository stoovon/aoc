package solve2021

import (
	"aoc/solve"
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SnailfishNumber interface {
	IsRegular() bool
	Magnitude() int
	String() string
}

type Regular struct {
	Value int
}

func (r *Regular) IsRegular() bool { return true }
func (r *Regular) Magnitude() int  { return r.Value }
func (r *Regular) String() string  { return fmt.Sprintf("%d", r.Value) }

type Pair struct {
	Left, Right SnailfishNumber
}

func (p *Pair) IsRegular() bool { return false }
func (p *Pair) Magnitude() int {
	return 3*p.Left.Magnitude() + 2*p.Right.Magnitude()
}
func (p *Pair) String() string {
	return fmt.Sprintf("[%s,%s]", p.Left.String(), p.Right.String())
}

// Parse a snailfish number from string
func ParseSnailfish(s string) SnailfishNumber {
	s = strings.TrimSpace(s)
	if s[0] != '[' {
		v, _ := strconv.Atoi(s)
		return &Regular{Value: v}
	}
	// Find comma separating top-level pair
	depth := 0
	for i := 1; i < len(s)-1; i++ {
		switch s[i] {
		case '[':
			depth++
		case ']':
			depth--
		case ',':
			if depth == 0 {
				left := ParseSnailfish(s[1:i])
				right := ParseSnailfish(s[i+1 : len(s)-1])
				return &Pair{Left: left, Right: right}
			}
		}
	}
	panic("invalid snailfish number: " + s)
}

// Add two snailfish numbers and reduce
func AddSnailfish(a, b SnailfishNumber) SnailfishNumber {
	sum := &Pair{Left: a, Right: b}
	return Reduce(sum)
}

// Reduce a snailfish number
func Reduce(n SnailfishNumber) SnailfishNumber {
	for {
		exploded, _, _, changed := Explode(n, 0)
		if changed {
			n = exploded
			continue
		}
		splitted, changed := Split(n)
		if changed {
			n = splitted
			continue
		}
		break
	}
	return n
}

// Explode returns (new number, left add, right add, changed)
func Explode(n SnailfishNumber, depth int) (SnailfishNumber, int, int, bool) {
	if n.IsRegular() {
		return n, 0, 0, false
	}
	p := n.(*Pair)
	if depth == 4 {
		// Explode this pair
		l := p.Left.(*Regular).Value
		r := p.Right.(*Regular).Value
		return &Regular{Value: 0}, l, r, true
	}
	// Explode left
	left, lAdd, rAdd, changed := Explode(p.Left, depth+1)
	if changed {
		// Add rAdd to first regular on right
		newRight := AddLeftmost(p.Right, rAdd)
		return &Pair{Left: left, Right: newRight}, lAdd, 0, true
	}
	// Explode right
	right, lAdd, rAdd, changed := Explode(p.Right, depth+1)
	if changed {
		// Add lAdd to first regular on left
		newLeft := AddRightmost(p.Left, lAdd)
		return &Pair{Left: newLeft, Right: right}, 0, rAdd, true
	}
	return n, 0, 0, false
}

// Add value to leftmost regular number
func AddLeftmost(n SnailfishNumber, v int) SnailfishNumber {
	if v == 0 {
		return n
	}
	if n.IsRegular() {
		return &Regular{Value: n.(*Regular).Value + v}
	}
	p := n.(*Pair)
	return &Pair{Left: AddLeftmost(p.Left, v), Right: p.Right}
}

// Add value to rightmost regular number
func AddRightmost(n SnailfishNumber, v int) SnailfishNumber {
	if v == 0 {
		return n
	}
	if n.IsRegular() {
		return &Regular{Value: n.(*Regular).Value + v}
	}
	p := n.(*Pair)
	return &Pair{Left: p.Left, Right: AddRightmost(p.Right, v)}
}

// Split returns (new number, changed)
func Split(n SnailfishNumber) (SnailfishNumber, bool) {
	if n.IsRegular() {
		v := n.(*Regular).Value
		if v >= 10 {
			left := v / 2
			right := (v + 1) / 2
			return &Pair{Left: &Regular{Value: left}, Right: &Regular{Value: right}}, true
		}
		return n, false
	}
	p := n.(*Pair)
	left, changed := Split(p.Left)
	if changed {
		return &Pair{Left: left, Right: p.Right}, true
	}
	right, changed := Split(p.Right)
	if changed {
		return &Pair{Left: p.Left, Right: right}, true
	}
	return n, false
}

type Day18 struct{}

func (d Day18) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2021, Day: 18}
}

func (d Day18) Part1(input string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var nums []SnailfishNumber
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		nums = append(nums, ParseSnailfish(line))
	}
	if len(nums) == 0 {
		return "", errors.New("No input numbers found.")
	}
	sum := nums[0]
	for i := 1; i < len(nums); i++ {
		sum = AddSnailfish(sum, nums[i])
	}
	return strconv.Itoa(sum.Magnitude()), nil
}

func (d Day18) Part2(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var nums []SnailfishNumber
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		nums = append(nums, ParseSnailfish(line))
	}
	maxMag := 0
	for i := range nums {
		for j := range nums {
			if i == j {
				continue
			}
			// Deep copy both numbers
			a := ParseSnailfish(nums[i].String())
			b := ParseSnailfish(nums[j].String())
			sum := AddSnailfish(a, b)
			mag := sum.Magnitude()
			if mag > maxMag {
				maxMag = mag
			}
		}
	}
	return strconv.Itoa(maxMag), nil
}

func init() {
	solve.Register(Day18{})
}
