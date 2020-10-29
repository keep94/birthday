package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/keep94/birthday"
)

var (
	fFile      string
	fDaysAhead int
)

func main() {
	flag.Parse()
	if fFile == "" {
		fmt.Println("Need to specify at least -file flag.")
		flag.Usage()
		os.Exit(1)
	}
	reminder := birthday.NewReminder(birthday.Today(), fDaysAhead)
	err := birthday.ReadFile(fFile, reminder)
	if err != nil {
		log.Fatal(err)
	}
	for _, milestone := range reminder.Milestones() {
		ageStr := "? years"
		if milestone.Age != -1 {
			ageStr = fmt.Sprintf("%d %s", milestone.Age, milestone.Unit)
		}
		astricks := " "
		if milestone.DaysAway == 0 {
			astricks = "*"
		}
		fmt.Printf(
			"%s %14s %12s %s\n",
			astricks,
			birthday.ToStringWithWeekDay(milestone.Date),
			ageStr,
			milestone.Name)
	}
}

func init() {
	flag.StringVar(&fFile, "file", "", "Birthday file")
	flag.IntVar(&fDaysAhead, "days_ahead", 21, "Days ahead")
}
