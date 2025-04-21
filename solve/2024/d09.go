package solve2024

import (
	"strconv"

	"aoc/solve"
)

type Day9 struct {
}

func (d Day9) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2024, Day: 9}
}

type file struct {
	ID   int
	Size int
}

type day9Disk struct {
	Blocks []int // -1 represents free space, non-negative numbers represent file IDs
}

// parseDiskMap coverts the input string into a sequence of files and free spaces
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
	fileID := 0

	for i := 0; i < len(files); i++ {
		// file blocks
		for j := 0; j < files[i]; j++ {
			blocks = append(blocks, fileID)
		}
		fileID++

		// free space
		if i < len(spaces) {
			for j := 0; j < spaces[i]; j++ {
				blocks = append(blocks, -1)
			}
		}
	}

	return day9Disk{Blocks: blocks}
}

// findFirstFreeSpace finds the leftmost free space in the disk
func (d *day9Disk) findFirstFreeSpace() int {
	for i, block := range d.Blocks {
		if block == -1 {
			return i
		}
	}
	return -1
}

// findLastFileBlock finds the position of the rightmost block of any file
func (d *day9Disk) findLastFileBlock() int {
	for i := len(d.Blocks) - 1; i >= 0; i-- {
		if d.Blocks[i] != -1 {
			return i
		}
	}
	return -1
}

func (d *day9Disk) moveOneBlock(fromIndex, toIndex int) {
	fileID := d.Blocks[fromIndex]
	d.Blocks[fromIndex] = -1
	d.Blocks[toIndex] = fileID
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

// findFileSize finds the size of the file containing pos
func (d *day9Disk) findFileSize(pos int) int {
	if pos < 0 || pos >= len(d.Blocks) || d.Blocks[pos] == -1 {
		return 0
	}

	fileID := d.Blocks[pos]
	start := pos
	for start >= 0 && d.Blocks[start] == fileID {
		start--
	}
	start++

	end := pos
	for end < len(d.Blocks) && d.Blocks[end] == fileID {
		end++
	}

	return end - start
}

// findFileStart finds the starting position of the file containing pos
func (d *day9Disk) findFileStart(pos int) int {
	if pos < 0 || pos >= len(d.Blocks) || d.Blocks[pos] == -1 {
		return -1
	}

	fileID := d.Blocks[pos]
	start := pos
	for start >= 0 && d.Blocks[start] == fileID {
		start--
	}
	return start + 1
}

// findFreeSpaceSize finds size of contiguous free space starting at pos
func (d *day9Disk) findFreeSpaceSize(pos int) int {
	size := 0
	for i := pos; i < len(d.Blocks) && d.Blocks[i] == -1; i++ {
		size++
	}
	return size
}

func (d *day9Disk) moveWholeFile(fromIndex, toIndex, size int) {
	fileID := d.Blocks[fromIndex]

	// Clear old location
	for i := 0; i < size; i++ {
		d.Blocks[fromIndex+i] = -1
	}

	// Place at new location
	for i := 0; i < size; i++ {
		d.Blocks[toIndex+i] = fileID
	}
}

// compactDisk performs the disk compaction process
func (d Day9) compactDisk(diskMap string, part int) int {
	// Parse the input
	files, spaces := d.parseDiskMap(diskMap)
	disk := d.createInitialDiskState(files, spaces)

	if part == 1 {
		// Part 1: Original block-by-block movement
		for {
			freeSpace := disk.findFirstFreeSpace()
			if freeSpace == -1 {
				break
			}

			lastBlock := disk.findLastFileBlock()
			if lastBlock == -1 || lastBlock <= freeSpace {
				break
			}

			disk.moveOneBlock(lastBlock, freeSpace)
		}
	} else {
		// Part 2: Move whole files in decreasing file ID order
		maxFileID := -1
		for _, block := range disk.Blocks {
			if block > maxFileID {
				maxFileID = block
			}
		}

		for fileID := maxFileID; fileID >= 0; fileID-- {
			var filePos int = -1
			for i := len(disk.Blocks) - 1; i >= 0; i-- {
				if disk.Blocks[i] == fileID {
					filePos = disk.findFileStart(i)
					break
				}
			}

			if filePos == -1 {
				continue
			}

			fileSize := disk.findFileSize(filePos)

			var bestFreeSpace int = -1
			for i := 0; i < filePos; i++ {
				if disk.Blocks[i] == -1 {
					freeSize := disk.findFreeSpaceSize(i)
					if freeSize >= fileSize {
						bestFreeSpace = i
						break
					}
					i += freeSize - 1 // Skip to end of this free space
				}
			}

			// Move the file if we found suitable space
			if bestFreeSpace != -1 && bestFreeSpace < filePos {
				disk.moveWholeFile(filePos, bestFreeSpace, fileSize)
			}
		}
	}

	return disk.calculateChecksum()
}

func (d Day9) Part1(input string) (string, error) {
	return strconv.Itoa(d.compactDisk(input, 1)), nil
}

func (d Day9) Part2(input string) (string, error) {
	return strconv.Itoa(d.compactDisk(input, 2)), nil
}

func init() {
	solve.Register(Day9{})
}
