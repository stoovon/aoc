package solve2015

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2015, Day: 4}
}

func (d Day4) mine(input string, difficulty int, start int) int {
	// Trim any leading or trailing whitespace/newlines from the input
	input = strings.TrimSpace(input)

	prefix := fmt.Sprintf("%0*s", difficulty, "") // Generate the prefix of zeros
	i := start
	for {
		data := fmt.Sprintf("%s%d", input, i)
		hash := md5.Sum([]byte(data))
		hashStr := hex.EncodeToString(hash[:])
		if hashStr[:difficulty] == prefix {
			return i
		}
		i++
	}
}

func (d Day4) Part1(input string) (string, error) {
	result := d.mine(input, 5, 0)
	return strconv.Itoa(result), nil
}

func (d Day4) Part2(input string) (string, error) {
	a := d.mine(input, 5, 0)
	result := d.mine(input, 6, a-1)
	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day4{})
}
