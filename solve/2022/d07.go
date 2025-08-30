package solve2022

import (
	"aoc/solve"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Day7 struct{}

func (d Day7) Coords() solve.SolutionCoords {
	return solve.SolutionCoords{Year: 2022, Day: 7}
}

func (d Day7) Part1(input string) (string, error) {
	return d.solve(input, 100000)
}

func (d Day7) Part2(input string) (string, error) {
	return d.solve(input, 30000000)
}

func (d Day7) solve(input string, limit int) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	root := newDirectory()
	currentPath := []string{}
	state := "normal"

	for _, line := range lines {
		if strings.HasPrefix(line, "$") {
			state = "normal"
		}
		if strings.HasPrefix(line, "$ cd ..") {
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
		} else if strings.HasPrefix(line, "$ cd ") {
			currentPath = append(currentPath, line[5:])
		} else if strings.HasPrefix(line, "$ ls") {
			state = "ls"
		} else if state == "ls" {
			if !strings.HasPrefix(line, "dir") {
				parts := strings.Fields(line)
				size, _ := strconv.Atoi(parts[0])
				name := parts[1]
				addData(root, currentPath, name, size)
			}
		}
	}

	var dirs []struct {
		Path []string
		Size int
	}
	totalSize := getSize(root, []string{}, &dirs)

	// Calculate score for Part 1
	score := 0
	for _, dir := range dirs {
		if dir.Size <= limit {
			score += dir.Size
		}
	}

	if limit == 100000 {
		return fmt.Sprintf("%d", score), nil
	}

	// Part 2
	freeSpace := 70000000 - totalSize
	spaceToBeFreed := 30000000 - freeSpace

	var possibleDirs []struct {
		Path []string
		Size int
	}
	for _, dir := range dirs {
		if dir.Size >= spaceToBeFreed {
			possibleDirs = append(possibleDirs, dir)
		}
	}

	sort.Slice(possibleDirs, func(i, j int) bool {
		return possibleDirs[i].Size < possibleDirs[j].Size
	})

	if len(possibleDirs) > 0 {
		return strconv.Itoa(possibleDirs[0].Size), nil
	}

	return "", nil
}

type Directory struct {
	Files map[string]int
	Dirs  map[string]*Directory
}

func newDirectory() *Directory {
	return &Directory{
		Files: make(map[string]int),
		Dirs:  make(map[string]*Directory),
	}
}

func addData(dir *Directory, path []string, name string, size int) {
	if len(path) == 0 {
		dir.Files[name] = size
		return
	}
	subDir, exists := dir.Dirs[path[0]]
	if !exists {
		subDir = newDirectory()
		dir.Dirs[path[0]] = subDir
	}
	addData(subDir, path[1:], name, size)
}

func getSize(dir *Directory, path []string, dirs *[]struct {
	Path []string
	Size int
}) int {
	size := 0
	for _, fileSize := range dir.Files {
		size += fileSize
	}
	for subDirName, subDir := range dir.Dirs {
		subPath := append(path, subDirName)
		size += getSize(subDir, subPath, dirs)
	}
	*dirs = append(*dirs, struct {
		Path []string
		Size int
	}{Path: path, Size: size})
	return size
}

func init() {
	solve.Register(Day7{})
}
