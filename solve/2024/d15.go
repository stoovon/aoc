package solve2024

import (
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/grids"
)

type Day15 struct {
}

func (day Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 15}
}

func (day Day15) convertMoves(moves string) [][]int {
	dirs := grids.DURL()

	var D = []rune{'>', '<', 'v', '^'}

	moves = strings.ReplaceAll(moves, "\n", "")
	result := make([][]int, len(moves))
	for i, m := range moves {
		for j, d := range D {
			if m == d {
				result[i] = []int{dirs[j].X, dirs[j].Y}
				break
			}
		}
	}
	return result
}

func (day Day15) addPair(x, y []int) []int {
	if len(x) == 2 {
		return []int{x[0] + y[0], x[1] + y[1]}
	}
	result := make([]int, len(x))
	for i := range x {
		result[i] = x[i] + y[i]
	}
	return result
}

func (day Day15) push(box [2]int, d []int, walls map[[2]int]struct{}, boxes map[[2]int]struct{}) bool {
	nxt := day.addPair(box[:], d)
	nxtBox := [2]int{nxt[0], nxt[1]}

	if _, exists := walls[nxtBox]; exists {
		return false
	}
	if _, exists := boxes[nxtBox]; exists {
		if !day.push(nxtBox, d, walls, boxes) {
			return false
		}
	}
	delete(boxes, box)
	boxes[nxtBox] = struct{}{}
	return true
}

func (day Day15) doubleWidePush(box [2]int, d []int, walls map[[2]int]struct{}, boxes map[[2]int]struct{}) bool {
	_, boxExists := boxes[box]
	if !boxExists {
		panic("box should be in boxes")
	}

	nxt := day.addPair(box[:], d)
	nxtPos := [2]int{nxt[0], nxt[1]}
	rightNxt := day.turnRight(nxtPos)

	if _, exists := walls[nxtPos]; exists || func() bool { _, exists := walls[rightNxt]; return exists }() {
		return false
	}

	if d[0] != 0 {
		// we are moving up/down
		if _, exists := boxes[nxtPos]; exists {
			if !day.doubleWidePush(nxtPos, d, walls, boxes) {
				return false
			}
		}
		leftNxt := day.turnLeft(nxtPos)
		if _, exists := boxes[leftNxt]; exists {
			if !day.doubleWidePush(leftNxt, d, walls, boxes) {
				return false
			}
		}
		if _, exists := boxes[rightNxt]; exists {
			if !day.doubleWidePush(rightNxt, d, walls, boxes) {
				return false
			}
		}
	}

	if d[1] == 1 {
		// we are pushing right
		if _, exists := boxes[rightNxt]; exists {
			if !day.doubleWidePush(rightNxt, d, walls, boxes) {
				return false
			}
		}
	}

	if d[1] == -1 {
		// we are pushing left
		leftNxt := day.turnLeft(nxtPos)
		if _, exists := boxes[leftNxt]; exists {
			if !day.doubleWidePush(leftNxt, d, walls, boxes) {
				return false
			}
		}
	}

	delete(boxes, box)
	boxes[nxtPos] = struct{}{}
	return true
}

func (day Day15) turnLeft(pos [2]int) [2]int {
	return [2]int{pos[0], pos[1] - 1}
}

func (day Day15) turnRight(pos [2]int) [2]int {
	return [2]int{pos[0], pos[1] + 1}
}

func (day Day15) Part1(input string) (string, error) {
	sections := strings.Split(strings.TrimSpace(input), "\n\n")
	ll := strings.Split(sections[0], "\n")
	moves := day.convertMoves(sections[1])

	walls := make(map[[2]int]struct{})
	boxes := make(map[[2]int]struct{})
	var robot [2]int

	for i, l := range ll {
		for j, ch := range l {
			switch ch {
			case '#':
				walls[[2]int{i, j}] = struct{}{}
			case 'O':
				boxes[[2]int{i, j}] = struct{}{}
			case '@':
				robot = [2]int{i, j}
			}
		}
	}

	for _, move := range moves {
		nxt := day.addPair(robot[:], move)
		nxtPos := [2]int{nxt[0], nxt[1]}

		if _, exists := walls[nxtPos]; exists {
			continue
		}
		if _, exists := boxes[nxtPos]; exists {
			if !day.push(nxtPos, move, walls, boxes) {
				continue
			}
		}
		if _, exists := boxes[nxtPos]; exists {
			panic("nxt should not be in boxes")
		}
		robot = nxtPos
	}

	c := 0
	for box := range boxes {
		c += 100*box[0] + box[1]
	}
	return strconv.Itoa(c), nil
}

func (day Day15) Part2(input string) (string, error) {
	sections := strings.Split(strings.TrimSpace(input), "\n\n")
	ll := strings.Split(sections[0], "\n")
	moves := day.convertMoves(sections[1])

	walls := make(map[[2]int]struct{})
	boxes := make(map[[2]int]struct{})
	var robot [2]int

	for i, l := range ll {
		for j, ch := range l {
			j *= 2
			switch ch {
			case '#':
				walls[[2]int{i, j}] = struct{}{}
				walls[[2]int{i, j + 1}] = struct{}{}
			case 'O':
				boxes[[2]int{i, j}] = struct{}{}
			case '@':
				robot = [2]int{i, j}
			}
		}
	}

	for _, move := range moves {
		for box := range boxes {
			if _, exists := boxes[day.turnRight(box)]; exists {
				panic("right(box) should not be in boxes")
			}
			if _, exists := walls[day.turnRight(box)]; exists {
				panic("right(box) should not be in walls")
			}
		}

		nxt := day.addPair(robot[:], move)
		nxtPos := [2]int{nxt[0], nxt[1]}

		if _, exists := walls[nxtPos]; exists {
			continue
		}

		if _, exists := boxes[nxtPos]; exists {
			copyBoxes := make(map[[2]int]struct{})
			for k, v := range boxes {
				copyBoxes[k] = v
			}
			if !day.doubleWidePush(nxtPos, move, walls, boxes) {
				boxes = copyBoxes
				continue
			}
		} else if _, exists := boxes[day.turnLeft(nxtPos)]; exists {
			copyBoxes := make(map[[2]int]struct{})
			for k, v := range boxes {
				copyBoxes[k] = v
			}
			if !day.doubleWidePush(day.turnLeft(nxtPos), move, walls, boxes) {
				boxes = copyBoxes
				continue
			}
		}

		if _, exists := boxes[nxtPos]; exists {
			panic("nxt should not be in boxes")
		}
		if _, exists := boxes[day.turnLeft(nxtPos)]; exists {
			panic("left(nxt) should not be in boxes")
		}

		robot = nxtPos
	}

	c := 0
	for box := range boxes {
		c += 100*box[0] + box[1]
	}
	return strconv.Itoa(c), nil
}

func init() {
	solve.Register(Day15{})
}
