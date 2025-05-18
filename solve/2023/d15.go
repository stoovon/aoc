package solve2023

import (
	"errors"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day15 struct {
}

func (d Day15) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2023, Day: 15}
}

func day15HashFn(s string) int {
	hash := 0
	for _, c := range s {
		hash = (hash + int(c)) * 17 % 256
	}
	return hash
}

type lens struct {
	label string
	value int
}

type focusHashmap struct {
	data [][]lens
}

func newFocusHashmap() *focusHashmap {
	data := make([][]lens, 256)
	return &focusHashmap{data: data}
}

func (f *focusHashmap) insert(key string, value int) {
	bucket := day15HashFn(key)
	for i := range f.data[bucket] {
		if f.data[bucket][i].label == key {
			f.data[bucket][i].value = value
			return
		}
	}
	f.data[bucket] = append(f.data[bucket], lens{label: key, value: value})
}

func (f *focusHashmap) remove(key string) {
	bucket := day15HashFn(key)
	for i, l := range f.data[bucket] {
		if l.label == key {
			f.data[bucket] = append(f.data[bucket][:i], f.data[bucket][i+1:]...)
			return
		}
	}
}

func (f *focusHashmap) totalPower() int {
	total := 0
	for boxIdx, box := range f.data {
		for slotIdx, l := range box {
			total += (boxIdx + 1) * (slotIdx + 1) * l.value
		}
	}
	return total
}

func (d Day15) Part1(input string) (string, error) {
	parts := strings.Split(strings.TrimSpace(input), ",")
	sum := 0
	for _, s := range parts {
		sum += day15HashFn(s)
	}
	return strconv.Itoa(sum), nil
}

func (d Day15) Part2(input string) (string, error) {
	parts := strings.Split(strings.TrimSpace(input), ",")
	fh := newFocusHashmap()
	for _, s := range parts {
		if strings.Contains(s, "-") {
			key := strings.TrimSuffix(s, "-")
			fh.remove(key)
		} else {
			kv := strings.SplitN(s, "=", 2)
			if len(kv) != 2 {
				return "", errors.New("invalid input")
			}
			val, err := strconv.Atoi(kv[1])
			if err != nil {
				val = 0
			}
			fh.insert(kv[0], val)
		}
	}
	return strconv.Itoa(fh.totalPower()), nil
}

func init() {
	solve.Register(Day15{})
}
