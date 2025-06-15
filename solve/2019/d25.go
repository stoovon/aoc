package solve2019

import (
	"errors"
	"regexp"
	"strings"

	"aoc/solve"
)

type Day25 struct {
}

func (d Day25) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2019, Day: 25}
}

func (d Day25) Part1(input string) (string, error) {
	program := parseIntcode(input)

	var output strings.Builder

	for {
		out, needInput := program.RunUntilInputOrHalt()
		output.WriteString(out)

		if !needInput {
			break
		}

		// Initial commands to get us to the checkpoint with the safe items
		program.ProvideASCIIInput("west\n")
		program.ProvideASCIIInput("west\n")
		program.ProvideASCIIInput("west\n")
		program.ProvideASCIIInput("west\n")
		program.ProvideASCIIInput("take dark matter\n")
		program.ProvideASCIIInput("east\n")
		program.ProvideASCIIInput("south\n")
		program.ProvideASCIIInput("west\n")
		program.ProvideASCIIInput("take food ration\n")
		program.ProvideASCIIInput("east\n")
		program.ProvideASCIIInput("north\n")
		program.ProvideASCIIInput("east\n")
		program.ProvideASCIIInput("east\n")
		program.ProvideASCIIInput("south\n")
		program.ProvideASCIIInput("south\n")
		program.ProvideASCIIInput("south\n")
		program.ProvideASCIIInput("take asterisk\n")
		program.ProvideASCIIInput("north\n")
		program.ProvideASCIIInput("north\n")
		program.ProvideASCIIInput("north\n")
		program.ProvideASCIIInput("west\n")
		program.ProvideASCIIInput("south\n")
		program.ProvideASCIIInput("take astronaut ice cream\n")
		program.ProvideASCIIInput("south\n")
		program.ProvideASCIIInput("take polygon\n")
		program.ProvideASCIIInput("east\n")
		program.ProvideASCIIInput("take easter egg\n")
		program.ProvideASCIIInput("east\n")
		program.ProvideASCIIInput("take weather machine\n")
		program.ProvideASCIIInput("north\n")

		items := []string{"dark matter", "food ration", "fixed point", "astronaut ice cream", "polygon", "asterisk", "easter egg", "weather machine"}
		n := len(items)
		total := 1 << n

		for perm := 0; perm < total; perm++ {
			for _, item := range items {
				program.ProvideASCIIInput("drop " + item + "\n")
			}

			for i := 0; i < n; i++ {
				if perm&(1<<i) != 0 {
					program.ProvideASCIIInput("take " + items[i] + "\n")
				}
			}

			program.ProvideASCIIInput("north\n")

			continue
		}
	}

	re := regexp.MustCompile(`typing (\d+)`)
	matches := re.FindStringSubmatch(output.String())
	if len(matches) > 1 {
		return matches[1], nil
	}
	return output.String(), nil
}

func (d Day25) Part2(input string) (string, error) {
	return "", errors.New("Not implemented")
}

func init() {
	solve.Register(Day25{})
}
