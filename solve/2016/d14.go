package solve2016

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day14 struct {
	hashCache map[int]string
}

func (d Day14) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2016, Day: 14}
}

// Compute the MD5 hash of a string
func (d Day14) hashval(index int, salt string, stretch int) string {
	if hash, exists := d.hashCache[index]; exists {
		return hash
	}

	h := md5.Sum([]byte(salt + strconv.Itoa(index)))
	hash := hex.EncodeToString(h[:])

	for i := 0; i < stretch; i++ {
		h = md5.Sum([]byte(hash))
		hash = hex.EncodeToString(h[:])
	}

	d.hashCache[index] = hash
	return hash
}

func (d Day14) hasTripleDigits(input string) byte {
	for i := 0; i < len(input)-2; i++ {
		if input[i] == input[i+1] && input[i+1] == input[i+2] {
			return input[i]
		}
	}
	return 0
}

func (d Day14) hasQuintuple(hash string, char byte) bool {
	target := strings.Repeat(string(char), 5)
	return strings.Contains(hash, target)
}

func (d Day14) isKey(index int, salt string, stretch int) bool {
	hash := d.hashval(index, salt, stretch)

	tripleMatch := d.hasTripleDigits(hash)
	if tripleMatch != 0 {
		for i := 1; i <= 1000; i++ {
			if d.hasQuintuple(d.hashval(index+i, salt, stretch), tripleMatch) {
				return true
			}
		}
	}

	return false
}

func (d Day14) nthKey(n int, salt string, stretch int) int {
	d.hashCache = make(map[int]string)
	count := 0
	for i := 0; ; i++ {
		if d.isKey(i, salt, stretch) {
			count++
			if count == n {
				return i
			}
		}
	}
}

func (d Day14) Part1(input string) (string, error) {
	result := d.nthKey(64, strings.TrimSpace(input), 0)
	return strconv.Itoa(result), nil
}

func (d Day14) Part2(input string) (string, error) {
	result := d.nthKey(64, strings.TrimSpace(input), 2016)
	return strconv.Itoa(result), nil
}

func init() {
	solve.Register(Day14{})
}
