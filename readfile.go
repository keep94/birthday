package birthday

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ReadFile reads a file of birthdays and returns upcoming milestones.
// filename is the birthday file; daysAhead is how many days of milestones
// to show. Line in the birtday file are of format Name<tab>birthday. Blank
// lines and lines starting with '#' are ignored.
func ReadFile(
	filename string,
	currentDate time.Time,
	daysAhead int) ([]Milestone, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	remind := NewRemind(currentDate, daysAhead)
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
			return nil, fmt.Errorf("Line %d malformatted", lineNo)
		}
		name := strings.TrimSpace(parts[0])
		bdayStr := strings.TrimSpace(parts[1])
		bday, err := Parse(bdayStr)
		if err != nil {
			return nil, fmt.Errorf("Line %d contains invalid birthday", lineNo)
		}
		remind.Add(name, bday)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return remind.Reminders(), nil
}
