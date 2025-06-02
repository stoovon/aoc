package acstrings

import (
	"strconv"
	"strings"
)

func MustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func Lines(input string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}
