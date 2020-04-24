package roadmap

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

// NewIO creates a IO instance
func NewIO() IO {
	return IO{}
}

// IO represents a persistence layer using the file system (or standard i/o)
type IO struct {
}

// Get reads a Roadmap from the file system (or standard i/o)
func (frw IO) Read(input string) ([]string, error) {
	var (
		file = os.Stdin
		err  error
	)

	if input != "" {
		file, err = os.Open(input)
		if err != nil {
			return nil, fmt.Errorf("can't open file (%s): %w", input, err)
		}
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file (%s): %w", input, err)
	}

	return lines, nil
}

// Write writes a roadmap to the file system (or standard i/o)
func (frw IO) Write(output string, content string) error {
	if output == "" {
		_, err := fmt.Print(content)

		return err
	}

	d1 := []byte(content)
	err := ioutil.WriteFile(output, d1, 0644)

	return err
}
