package birthday

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/keep94/consume2"
)

// Interface Store abstracts away reading the birthday file for testability.
type Store interface {
	Read(consumer consume2.Consumer[Entry]) error
}

// SystemStore reads birthday file from given file path.
type SystemStore string

// Read reads the birthday file at path s.
// consumer consumes the Entry instances read.
func (s SystemStore) Read(consumer consume2.Consumer[Entry]) error {
	return ReadFile(string(s), consumer)
}

// ReadFile reads a birthday file.
// consumer consumes the Entry instances read.
func ReadFile(filename string, consumer consume2.Consumer[Entry]) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return Read(file, consumer)
}

// Read reads a birthday file. consumer consumes the Entry instances read.
func Read(r io.Reader, consumer consume2.Consumer[Entry]) error {
	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() && consumer.CanConsume() {
		var entry Entry
		lineNo++
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			return fmt.Errorf("Line %d malformatted", lineNo)
		}
		entry.Name = strings.TrimSpace(parts[0])
		bdayStr := strings.TrimSpace(parts[1])
		var err error
		entry.Birthday, err = Parse(bdayStr)
		if err != nil {
			return fmt.Errorf("Line %d contains invalid birthday", lineNo)
		}
		consumer.Consume(entry)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
