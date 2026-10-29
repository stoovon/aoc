package solve2024

import (
	"strconv"

	"aoc/solve"
)

type Day9 struct{}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 9}
}

type day9Span struct {
	Start int
	Size  int
}

type day9Disk struct {
	Blocks []int // -1 represents free space, non-negative numbers represent file IDs
}

func (d Day9) parseDiskMap(input string) ([]int, []int) {
	files := make([]int, 0)
	spaces := make([]int, 0)

	for i := 0; i < len(input); i++ {
		size, _ := strconv.Atoi(string(input[i]))
		if i%2 == 0 {
			files = append(files, size)
		} else {
			spaces = append(spaces, size)
		}
	}

	return files, spaces
}

func (d Day9) createInitialDiskState(files, spaces []int) day9Disk {
	var blocks []int
	for fileID, size := range files {
		for range size {
			blocks = append(blocks, fileID)
		}
		if fileID < len(spaces) {
			for range spaces[fileID] {
				blocks = append(blocks, -1)
			}
		}
	}

	return day9Disk{Blocks: blocks}
}

func (d *day9Disk) compactBlocks() {
	left, right := 0, len(d.Blocks)-1
	for left < right {
		for left < right && d.Blocks[left] != -1 {
			left++
		}
		for left < right && d.Blocks[right] == -1 {
			right--
		}
		if left < right {
			d.Blocks[left], d.Blocks[right] = d.Blocks[right], d.Blocks[left]
			left++
			right--
		}
	}
}

func (d *day9Disk) calculateChecksum() int {
	cs := 0
	for pos, fileID := range d.Blocks {
		if fileID != -1 {
			cs += pos * fileID
		}
	}
	return cs
}

func (d Day9) createFileSpans(sizes, spaces []int) ([]day9Span, []day9Span) {
	files := make([]day9Span, len(sizes))
	gaps := make([]day9Span, 0, len(spaces))
	position := 0
	for id, size := range sizes {
		files[id] = day9Span{Start: position, Size: size}
		position += size
		if id < len(spaces) && spaces[id] > 0 {
			// Zero-length files do not separate adjacent gaps.
			if last := len(gaps) - 1; last >= 0 && gaps[last].Start+gaps[last].Size == position {
				gaps[last].Size += spaces[id]
			} else {
				gaps = append(gaps, day9Span{Start: position, Size: spaces[id]})
			}
			position += spaces[id]
		}
	}
	return files, gaps
}

func compactDay9Files(files, gaps []day9Span) {
	// Freed locations are right of every unprocessed file. Original gaps only
	// shrink, so each size's cursor can permanently skip undersized gaps.
	var nextGap [10]int
	for id := len(files) - 1; id >= 0; id-- {
		file := &files[id]
		if file.Size == 0 {
			continue
		}
		gapIndex := nextGap[file.Size]
		for gapIndex < len(gaps) && gaps[gapIndex].Size < file.Size {
			gapIndex++
		}
		nextGap[file.Size] = gapIndex
		if gapIndex == len(gaps) || gaps[gapIndex].Start >= file.Start {
			continue
		}
		gap := &gaps[gapIndex]
		file.Start = gap.Start
		gap.Start += file.Size
		gap.Size -= file.Size
	}
}

func (d Day9) Part1(input string) (string, error) {
	files, spaces := d.parseDiskMap(input)
	disk := d.createInitialDiskState(files, spaces)
	disk.compactBlocks()
	return strconv.Itoa(disk.calculateChecksum()), nil
}

func (d Day9) Part2(input string) (string, error) {
	sizes, spaces := d.parseDiskMap(input)
	files, gaps := d.createFileSpans(sizes, spaces)
	compactDay9Files(files, gaps)
	checksum := 0
	for id, file := range files {
		checksum += id * (file.Start*file.Size + file.Size*(file.Size-1)/2)
	}
	return strconv.Itoa(checksum), nil
}

func init() {
	solve.Register(Day9{})
}
