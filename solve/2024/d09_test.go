package solve2024

import (
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// The oracle follows the puzzle literally, restarting both searches after each move.
func compactBlocksReference(blocks []int) {
	for {
		gap := slices.Index(blocks, -1)
		last := len(blocks) - 1
		for last >= 0 && blocks[last] == -1 {
			last--
		}
		if gap < 0 || gap >= last {
			return
		}
		blocks[gap], blocks[last] = blocks[last], -1
	}
}

func TestDay9Part1EdgeCases(t *testing.T) {
	for _, test := range []struct{ input, want string }{
		{"12322", "18"},
		{"", "0"},
		{"00000", "0"},
		{"09", "0"},
		{"909", "117"},
	} {
		got, err := (Day9{}).Part1(test.input)
		if err != nil || got != test.want {
			t.Errorf("Part1(%q) = %q, %v; want %q", test.input, got, err, test.want)
		}
	}
}

func TestDay9BlockCompactionMatchesLiteralMoves(t *testing.T) {
	day := Day9{}
	cases := 0
	for length := 0; length <= 7; length++ {
		for encoded := 0; encoded < 1<<(2*length); encoded++ {
			input := make([]byte, length)
			value := encoded
			var expected []int
			for index := range input {
				size := value & 3
				value >>= 2
				input[index] = byte('0' + size)
				block := -1
				if index%2 == 0 {
					block = index / 2
				}
				for range size {
					expected = append(expected, block)
				}
			}
			files, spaces := day.parseDiskMap(string(input))
			disk := day.createInitialDiskState(files, spaces)
			compactBlocksReference(expected)
			disk.compactBlocks()
			if !slices.Equal(disk.Blocks, expected) {
				t.Fatalf("input %q: got %v, want %v", input, disk.Blocks, expected)
			}
			cases++
		}
	}
	if cases != 21845 {
		t.Fatalf("checked %d cases, want 21845", cases)
	}
}

func compactFilesReference(blocks []int) {
	maxID := -1
	for _, id := range blocks {
		maxID = max(maxID, id)
	}
	for id := maxID; id >= 0; id-- {
		start := slices.Index(blocks, id)
		if start < 0 {
			continue
		}
		end := start
		for end < len(blocks) && blocks[end] == id {
			end++
		}
		size := end - start
		for gap := 0; gap+size <= start; gap++ {
			if !slices.ContainsFunc(blocks[gap:gap+size], func(block int) bool { return block != -1 }) {
				for offset := range size {
					blocks[gap+offset], blocks[start+offset] = id, -1
				}
				break
			}
		}
	}
}

func checkDay9WholeFiles(t *testing.T, input string) {
	t.Helper()
	var expected []int
	for index := range len(input) {
		block := -1
		if index%2 == 0 {
			block = index / 2
		}
		for range int(input[index] - '0') {
			expected = append(expected, block)
		}
	}
	compactFilesReference(expected)
	day := Day9{}
	sizes, spaces := day.parseDiskMap(input)
	files, gaps := day.createFileSpans(sizes, spaces)
	compactDay9Files(files, gaps)
	actual := make([]int, len(expected))
	for index := range actual {
		actual[index] = -1
	}
	for id, file := range files {
		for offset := range file.Size {
			position := file.Start + offset
			if position < 0 || position >= len(actual) || actual[position] != -1 {
				t.Fatalf("input %q: overlapping or out-of-range file %+v", input, file)
			}
			actual[position] = id
		}
	}
	if !slices.Equal(actual, expected) {
		t.Fatalf("input %q: got %v, want %v", input, actual, expected)
	}
	checksum := 0
	for position, id := range expected {
		if id >= 0 {
			checksum += position * id
		}
	}
	got, err := day.Part2(input)
	if err != nil || got != strconv.Itoa(checksum) {
		t.Fatalf("Part2(%q) = %q, %v; want %d", input, got, err, checksum)
	}
}

func TestDay9WholeFileCompactionMatchesLiteralMoves(t *testing.T) {
	for length := 0; length <= 7; length++ {
		for encoded := 0; encoded < 1<<(2*length); encoded++ {
			input := make([]byte, length)
			value := encoded
			for index := range input {
				input[index] = byte('0' + (value & 3))
				value >>= 2
			}
			checkDay9WholeFiles(t, string(input))
		}
	}
}

func TestDay9WholeFilesWithAllDigitSizes(t *testing.T) {
	random := rand.New(rand.NewPCG(2024, 9))
	for range 1000 {
		input := make([]byte, random.IntN(50))
		for index := range input {
			input[index] = byte('0' + random.IntN(10))
		}
		checkDay9WholeFiles(t, string(input))
	}
}

func TestDay9WholeFileEdgeCases(t *testing.T) {
	for _, input := range []string{
		"", "0", "09", "909", "12322",
		"12023", "19099", "1010109", "111111111", "0909099",
	} {
		checkDay9WholeFiles(t, input)
	}

}
