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
	// Parse the Intcode program
	program := parseIntcode(input)

	// Set up input/output for manual play
	reader := bufio.NewReader(os.Stdin)
	var output strings.Builder

	for {
		// Run until input is needed or program halts
		out, needInput := program.RunUntilInputOrHalt()
		output.WriteString(out)

		fmt.Print(out)
		if !needInput {
			break
		}

		// Read user command
		fmt.Print("Command? ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if cmd == "cheat" {
			items := []string{"dark matter", "food ration", "fixed point", "astronaut ice cream", "polygon", "asterisk", "easter egg", "weather machine"}
			n := len(items)
			total := 1 << n // 256

			for perm := 0; perm < total; perm++ {
				// Drop all items
				for _, item := range items {
					program.ProvideASCIIInput("drop " + item + "\n")
				}

				// Take items according to current permutation
				for i := 0; i < n; i++ {
					if perm&(1<<i) != 0 {
						program.ProvideASCIIInput("take " + items[i] + "\n")
					}
				}

				// Try to pass the checkpoint
				program.ProvideASCIIInput("north\n")

				continue
			}
		}
		program.ProvideASCIIInput(cmd + "\n")
	}

	// The password is usually printed in the output
	return output.String(), nil
}

func (d Day25) Part2(input string) (string, error) {
	return "", errors.New("Not implemented")
}

func init() {
	solve.Register(Day25{})
}
