package solve2020

import (
	"fmt"
	"strconv"
	"strings"

	"aoc/solve"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2020, Day: 4}
}

var requiredFields = []string{
	"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid",
}

func hasRequiredFields(passport string) bool {
	for _, field := range requiredFields {
		if !strings.Contains(passport, field+":") {
			return false
		}
	}
	return true
}

func (d Day4) Part1(input string) (string, error) {
	passports := strings.Split(strings.TrimSpace(input), "\n\n")
	count := 0
	for _, passport := range passports {
		passport = strings.ReplaceAll(passport, "\n", " ")
		if hasRequiredFields(passport) {
			count++
		}
	}
	return fmt.Sprintf("%d", count), nil
}

func isValidField(key, value string) bool {
	switch key {
	case "byr":
		year, err := strconv.Atoi(value)
		return err == nil && year >= 1920 && year <= 2002
	case "iyr":
		year, err := strconv.Atoi(value)
		return err == nil && year >= 2010 && year <= 2020
	case "eyr":
		year, err := strconv.Atoi(value)
		return err == nil && year >= 2020 && year <= 2030
	case "hgt":
		if strings.HasSuffix(value, "cm") {
			num, err := strconv.Atoi(strings.TrimSuffix(value, "cm"))
			return err == nil && num >= 150 && num <= 193
		} else if strings.HasSuffix(value, "in") {
			num, err := strconv.Atoi(strings.TrimSuffix(value, "in"))
			return err == nil && num >= 59 && num <= 76
		}
		return false
	case "hcl":
		return len(value) == 7 && value[0] == '#' && strings.Count(value[1:], "#") == 0 &&
			strings.IndexFunc(value[1:], func(r rune) bool {
				return !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f')
			}) == -1
	case "ecl":
		switch value {
		case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
			return true
		}
		return false
	case "pid":
		if len(value) != 9 {
			return false
		}
		_, err := strconv.Atoi(value)
		return err == nil
	case "cid":
		return true
	}
	return false
}

func parsePassport(passport string) map[string]string {
	fields := strings.Fields(passport)
	m := make(map[string]string)
	for _, f := range fields {
		parts := strings.SplitN(f, ":", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func (d Day4) Part2(input string) (string, error) {
	passports := strings.Split(strings.TrimSpace(input), "\n\n")
	count := 0
	for _, passport := range passports {
		passport = strings.ReplaceAll(passport, "\n", " ")
		m := parsePassport(passport)
		valid := true
		for _, field := range requiredFields {
			val, ok := m[field]
			if !ok || !isValidField(field, val) {
				valid = false
				break
			}
		}
		if valid {
			count++
		}
	}
	return fmt.Sprintf("%d", count), nil
}

func init() {
	solve.Register(Day4{})
}
