package solve2016

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"aoc/solve"
)

type Day5 struct {
}

func (d Day5) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 5}
}

// findPassword finds the password by appending an incrementing number to the door ID
func (d Day5) findPassword(door string) string {
	password := ""
	for i := 0; ; i++ {
		hash := md5.Sum([]byte(fmt.Sprintf("%s%d", door, i)))
		hexHash := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hexHash, "00000") {
			password += string(hexHash[5])
			// Don't really need for debug, but it looks 1337 haxor 🆒
			// fmt.Printf("Progress: %d %s %s\n", i, hexHash, password)
			if len(password) == 8 {
				return password
			}
		}
	}
}

// findPasswordRedux finds the password using the 6th character as an index
func (d Day5) findPasswordRedux(door string) string {
	password := make([]rune, 8)
	for i := range password {
		password[i] = '_'
	}
	filled := 0

	for i := 0; ; i++ {
		hash := md5.Sum([]byte(fmt.Sprintf("%s%d", door, i)))
		hexHash := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hexHash, "00000") {
			index := int(hexHash[5] - '0')
			if index >= 0 && index < 8 && password[index] == '_' {
				password[index] = rune(hexHash[6])
				filled++
				// Don't really need for debug, but it looks 1337 haxor 🆒
				// fmt.Printf("Progress: %d %s %s\n", i, hexHash, string(password))
				if filled == 8 {
					return string(password)
				}
			}
		}
	}
}

func (d Day5) Part1(input string) (string, error) {
	return d.findPassword(strings.TrimSpace(input)), nil
}

func (d Day5) Part2(input string) (string, error) {
	return d.findPasswordRedux(strings.TrimSpace(input)), nil
}

func init() {
	solve.Register(Day5{})
}
