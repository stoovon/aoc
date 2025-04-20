package httpclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type cache struct {
}

func (c *cache) GetInput(year, day int) (string, bool, error) {
	return c.getFileData(year, day, "input")
}

func (c *cache) PutInput(year, day int, input string) error {
	return c.putFileData(year, day, "input", input)
}

func (c *cache) GetSolution(year, day, part int) (string, bool) {
	data, found, err := c.getFileData(year, day, fmt.Sprintf("solution_%d", part))
	if err != nil {
		return "", false
	}

	return data, found
}

func (c *cache) PutSolution(year, day, part int, solution string) error {
	return c.putFileData(year, day, fmt.Sprintf("solution_%d", part), solution)
}

func (c *cache) GetAttempt(year, day, part int, attempt string) (bool, error) {
	filename := fmt.Sprintf("attempts_%d", part)
	existingData, _, err := c.getFileData(year, day, filename)
	if err != nil {
		return false, fmt.Errorf("error retrieving existing attempts: %v", err)
	}

	if existingData == "" {
		return false, nil
	}

	attempts := strings.Split(existingData, "\n")
	for _, a := range attempts {
		if a == attempt {
			return true, nil
		}
	}

	return false, nil
}

func (c *cache) PutAttempt(year, day, part int, attempt string) error {
	filename := fmt.Sprintf("attempts_%d", part)
	existingData, _, err := c.getFileData(year, day, filename)
	if err != nil {
		return fmt.Errorf("error retrieving existing attempts: %v", err)
	}

	updatedData := existingData
	if existingData != "" {
		updatedData += "\n"
	}
	updatedData += attempt

	return c.putFileData(year, day, filename, updatedData)
}

func (c *cache) getFileData(year, day int, filename string) (string, bool, error) {
	cacheDir := ".cache"
	yearDir := filepath.Join(cacheDir, fmt.Sprintf("%d", year))
	dayDir := filepath.Join(yearDir, fmt.Sprintf("%02d", day))
	dataFile := filepath.Join(dayDir, filename+".txt")

	if _, err := os.Stat(dataFile); err == nil {
		content, err := os.ReadFile(dataFile)
		if err != nil {
			return "", false, err
		}
		return string(content), true, nil
	}

	return "", false, nil
}

func (c *cache) putFileData(year, day int, filename, data string) error {
	cacheDir := ".cache"
	yearDir := filepath.Join(cacheDir, fmt.Sprintf("%d", year))
	dayDir := filepath.Join(yearDir, fmt.Sprintf("%02d", day))
	dataFile := filepath.Join(dayDir, filename+".txt")

	if err := os.MkdirAll(dayDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directories: %v\n", err)
	}

	err := os.WriteFile(dataFile, []byte(data), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing data to file: %v\n", err)
	}

	return nil
}
