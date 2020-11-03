package birthday

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Interface Consumer consumes entries from a birthday file.
type Consumer interface {
	Consume(e *Entry)
}

// ReadFile reads a birthday file letting consumer consume each entry.
func ReadFile(filename string, consumer Consumer) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return Read(file, consumer)
}

// Read reads a birthday file letting consumer consume each entry.
func Read(r io.Reader, consumer Consumer) error {
	scanner := bufio.NewScanner(r)
	lineNo := 0
	var entry Entry
	for scanner.Scan() {
		entry = Entry{}
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
		consumer.Consume(&entry)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
