package solve2019

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

	reader := bufio.NewReader(os.Stdin)
	var output strings.Builder

	for {
		out, needInput := program.RunUntilInputOrHalt()
		output.WriteString(out)

		fmt.Print(out)
		if !needInput {
			break
		}

		fmt.Print("Command? ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if cmd == "cheat" {
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
		program.ProvideASCIIInput(cmd + "\n")
	}

	return output.String(), nil
}

func (d Day25) Part2(input string) (string, error) {
	return "", errors.New("Not implemented")
}

func init() {
	solve.Register(Day25{})
}
