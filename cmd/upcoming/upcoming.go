package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

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
	milestones, err := birthday.ReadFile(fFile, birthday.Today(), fDaysAhead)
	if err != nil {
		log.Fatal(err)
	}
	for _, milestone := range milestones {
		ageStr := "?????"
		if milestone.Age != -1 {
			ageStr = strconv.Itoa(milestone.Age)
		}
		astricks := " "
		if milestone.DaysAway == 0 {
			astricks = "*"
		}
		fmt.Printf(
			"%s %14s %5s %s\n",
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
