package solve2022

import (
	"aoc/solve"
	"errors"
	"strconv"
	"strings"
)

type Day20 struct{}

func (d Day20) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 20}
}

func (d Day20) Part1(input string) (string, error) {
	// Parse the input into a list of integers
	numbers := []int{}
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		num, err := strconv.Atoi(line)
		if err != nil {
			return "", err
		}
		numbers = append(numbers, num)
	}

	// Create a list of indices to track the original order
	indices := make([]int, len(numbers))
	for i := range indices {
		indices[i] = i
	}

	// Perform the mixing process
	for i := 0; i < len(numbers); i++ {
		// Find the current index of the number to move
		currentIndex := -1
		for j, idx := range indices {
			if idx == i {
				currentIndex = j
				break
			}
		}

		// Remove the number from the current position
		value := numbers[currentIndex]
		indices = append(indices[:currentIndex], indices[currentIndex+1:]...)
		numbers = append(numbers[:currentIndex], numbers[currentIndex+1:]...)

		// Calculate the new position
		newIndex := (currentIndex + value) % len(numbers)
		if newIndex < 0 {
			newIndex += len(numbers)
		}

		// Insert the number at the new position
		indices = append(indices[:newIndex], append([]int{i}, indices[newIndex:]...)...)
		numbers = append(numbers[:newIndex], append([]int{value}, numbers[newIndex:]...)...)
	}

	// Find the grove coordinates
	zeroIndex := -1
	for i, num := range numbers {
		if num == 0 {
			zeroIndex = i
			break
		}
	}

	if zeroIndex == -1 {
		return "", errors.New("0 not found in the list")
	}

	sum := 0
	for _, offset := range []int{1000, 2000, 3000} {
		index := (zeroIndex + offset) % len(numbers)
		sum += numbers[index]
	}

	return strconv.Itoa(sum), nil
}

func (d Day20) Part2(input string) (string, error) {
	const decryptionKey = 811589153
	const mixRounds = 10

	// Parse the input into a list of integers
	numbers := []int{}
	for line := range strings.SplitSeq(strings.TrimSpace(input), "\n") {
		num, err := strconv.Atoi(line)
		if err != nil {
			return "", err
		}
		numbers = append(numbers, num*decryptionKey)
	}

	// Create a list of indices to track the original order
	indices := make([]int, len(numbers))
	for i := range indices {
		indices[i] = i
	}

	// Perform the mixing process for the specified number of rounds
	for range mixRounds {
		for i := range numbers {
			// Find the current index of the number to move
			currentIndex := -1
			for j, idx := range indices {
				if idx == i {
					currentIndex = j
					break
				}
			}

			// Remove the number from the current position
			value := numbers[currentIndex]
			indices = append(indices[:currentIndex], indices[currentIndex+1:]...)
			numbers = append(numbers[:currentIndex], numbers[currentIndex+1:]...)

			// Calculate the new position
			newIndex := (currentIndex + value) % len(numbers)
			if newIndex < 0 {
				newIndex += len(numbers)
			}

			// Insert the number at the new position
			indices = append(indices[:newIndex], append([]int{i}, indices[newIndex:]...)...)
			numbers = append(numbers[:newIndex], append([]int{value}, numbers[newIndex:]...)...)
		}
	}

	// Find the grove coordinates
	zeroIndex := -1
	for i, num := range numbers {
		if num == 0 {
			zeroIndex = i
			break
		}
	}

	if zeroIndex == -1 {
		return "", errors.New("0 not found in the list")
	}

	sum := 0
	for _, offset := range []int{1000, 2000, 3000} {
		index := (zeroIndex + offset) % len(numbers)
		sum += numbers[index]
	}

	return strconv.Itoa(sum), nil
}

func init() {
	solve.Register(Day20{})
}
