package birthday

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ReadPersonsFromFile reads a file of birthdays and returns people whose
// names match query. queries ignore case and extra whitespace.
func ReadPersonsFromFile(
	filename string,
	currentDate time.Time,
	query string) ([]Person, error) {
	filter := NewFilter(currentDate, query)
	err := readFile(filename, filter)
	if err != nil {
		return nil, err
	}
	return filter.Persons(), nil
}

// ReadFile reads a file of birthdays and returns upcoming milestones.
// filename is the birthday file; daysAhead is how many days of milestones
// to show. Line in the birtday file are of format Name<tab>birthday. Blank
// lines and lines starting with '#' are ignored.
func ReadFile(
	filename string,
	currentDate time.Time,
	daysAhead int) ([]Milestone, error) {
	remind := NewRemind(currentDate, daysAhead)
	err := readFile(filename, remind)
	if err != nil {
		return nil, err
	}
	return remind.Reminders(), nil
}

type consumer interface {
	Add(name string, bday time.Time)
}

func readFile(filename string, c consumer) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Line %d malformatted", lineNo)
		}
		name := strings.TrimSpace(parts[0])
		bdayStr := strings.TrimSpace(parts[1])
		bday, err := Parse(bdayStr)
		if err != nil {
			return fmt.Errorf("Line %d contains invalid birthday", lineNo)
		}
		c.Add(name, bday)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
