package solve2016

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day21 struct {
}

func (d Day21) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 21}
}

var intRe = regexp.MustCompile(`\d+`)

func (d Day21) scramble(password string, instructions []string) string {
	pw := []rune(password)
	rot := func(n int) {
		n = (n%len(pw) + len(pw)) % len(pw)
		pw = append(pw[len(pw)-n:], pw[:len(pw)-n]...)
	}
	swap := func(a, b int) {
		pw[a], pw[b] = pw[b], pw[a]
	}
	for _, line := range instructions {
		words := strings.Fields(line)
		switch {
		case strings.HasPrefix(line, "swap position"):
			nums := d.parseInts(line)
			a, b := nums[0], nums[1]
			swap(a, b)
		case strings.HasPrefix(line, "swap letter"):
			swap(d.indexOf(pw, rune(words[2][0])), d.indexOf(pw, rune(words[5][0])))
		case strings.HasPrefix(line, "rotate right"):
			nums := d.parseInts(line)
			a := nums[0]
			rot(a)
		case strings.HasPrefix(line, "rotate left"):
			nums := d.parseInts(line)
			a := nums[0]
			rot(-a)
		case strings.HasPrefix(line, "reverse"):
			nums := d.parseInts(line)
			a, b := nums[0], nums[1]
			d.reverse(pw, a, b)
		case strings.HasPrefix(line, "move"):
			nums := d.parseInts(line)
			a, b := nums[0], nums[1]
			d.move(&pw, a, b)
		case strings.HasPrefix(line, "rotate based"):
			i := d.indexOf(pw, rune(words[6][0]))
			rot((i + 1 + map[bool]int{true: 1, false: 0}[i >= 4]) % len(pw))
		}
	}
	return string(pw)
}

func (d Day21) parseInts(s string) []int {
	matches := intRe.FindAllString(s, -1)
	nums := make([]int, len(matches))
	for i, m := range matches {
		nums[i], _ = strconv.Atoi(m)
	}
	return nums
}

func (d Day21) indexOf(pw []rune, r rune) int {
	for i, c := range pw {
		if c == r {
			return i
		}
	}
	return -1
}

func (d Day21) reverse(pw []rune, a, b int) {
	for i, j := a, b; i < j; i, j = i+1, j-1 {
		pw[i], pw[j] = pw[j], pw[i]
	}
}

func (d Day21) move(pw *[]rune, a, b int) {
	char := (*pw)[a]
	*pw = append((*pw)[:a], (*pw)[a+1:]...)
	*pw = append((*pw)[:b], append([]rune{char}, (*pw)[b:]...)...)
}

func (d Day21) permutations(s string) []string {
	var result []string
	var permute func([]rune, int)
	permute = func(arr []rune, n int) {
		if n == 1 {
			result = append(result, string(arr))
			return
		}
		for i := 0; i < n; i++ {
			permute(arr, n-1)
			if n%2 == 1 {
				arr[0], arr[n-1] = arr[n-1], arr[0]
			} else {
				arr[i], arr[n-1] = arr[n-1], arr[i]
			}
		}
	}
	permute([]rune(s), len(s))
	sort.Strings(result)
	return result
}

func (d Day21) Part1(input string) (string, error) {
	return d.scramble("abcdefgh", strings.Split(strings.TrimSpace(input), "\n")), nil
}

func (d Day21) Part2(input string) (string, error) {
	instructions := strings.Split(strings.TrimSpace(input), "\n")
	perms := d.permutations("fbgdceah")
	for _, p := range perms {
		if d.scramble(p, instructions) == "fbgdceah" {
			return p, nil
		}
	}
	return "", errors.New("no valid password found")
}

func init() {
	solve.Register(Day21{})
}
