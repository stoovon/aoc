package solve2018

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"aoc/solve"
	"aoc/utils/acstrings"
)

type Day4 struct {
}

func (d Day4) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2018, Day: 4}
}

var (
	guardRe = regexp.MustCompile(`Guard #(\d+) begins shift`)
	sleepRe = regexp.MustCompile(`(\d+)-(\d+)-(\d+) (\d+):(\d+)`)
)

type record struct {
	raw    string
	minute int
	event  string
}

func parseRecords(input string) []record {
	lines := acstrings.Lines(strings.TrimSpace(input))
	sort.Strings(lines)
	records := make([]record, len(lines))
	for i, line := range lines {
		parts := strings.SplitN(line, "] ", 2)
		timeStr := parts[0][1:]
		event := parts[1]
		m := sleepRe.FindStringSubmatch(timeStr)
		minute, _ := strconv.Atoi(m[5])
		records[i] = record{raw: line, minute: minute, event: event}
	}
	return records
}

func (d Day4) Part1(input string) (string, error) {
	records := parseRecords(input)
	guardSleep := map[int][60]int{}
	var guardID, sleepStart int

	for _, rec := range records {
		if m := guardRe.FindStringSubmatch(rec.event); m != nil {
			guardID, _ = strconv.Atoi(m[1])
		} else if rec.event == "falls asleep" {
			sleepStart = rec.minute
		} else if rec.event == "wakes up" {
			sleep := guardSleep[guardID]
			for m := sleepStart; m < rec.minute; m++ {
				sleep[m]++
			}
			guardSleep[guardID] = sleep
		}
	}

	maxSleep, sleepiestGuard := 0, 0
	for id, mins := range guardSleep {
		total := 0
		for _, v := range mins {
			total += v
		}
		if total > maxSleep {
			maxSleep = total
			sleepiestGuard = id
		}
	}

	sleepiestMinute, maxCount := 0, 0
	for m, v := range guardSleep[sleepiestGuard] {
		if v > maxCount {
			maxCount = v
			sleepiestMinute = m
		}
	}

	return fmt.Sprintf("%d", sleepiestGuard*sleepiestMinute), nil
}

func (d Day4) Part2(input string) (string, error) {
	records := parseRecords(input)
	guardSleep := map[int][60]int{}
	var guardID, sleepStart int

	for _, rec := range records {
		if m := guardRe.FindStringSubmatch(rec.event); m != nil {
			guardID, _ = strconv.Atoi(m[1])
		} else if rec.event == "falls asleep" {
			sleepStart = rec.minute
		} else if rec.event == "wakes up" {
			sleep := guardSleep[guardID]
			for m := sleepStart; m < rec.minute; m++ {
				sleep[m]++
			}
			guardSleep[guardID] = sleep
		}
	}

	maxCount, resultGuard, resultMinute := 0, 0, 0
	for id, mins := range guardSleep {
		for m, count := range mins {
			if count > maxCount {
				maxCount = count
				resultGuard = id
				resultMinute = m
			}
		}
	}

	return fmt.Sprintf("%d", resultGuard*resultMinute), nil
}

func init() {
	solve.Register(Day4{})
}
