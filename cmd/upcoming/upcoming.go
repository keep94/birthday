package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/keep94/birthday"
	"github.com/keep94/consume2"
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
	var entries []*birthday.Entry
	err := birthday.ReadFile(fFile, consume2.AppendPtrsTo(&entries))
	if err != nil {
		log.Fatal(err)
	}
	pipeline := consume2.PTakeWhile(func(m birthday.Milestone) bool {
		return m.DaysAway < fDaysAhead
	})
	birthday.Remind(
		entries,
		birthday.DefaultPeriods,
		birthday.Today(),
		pipeline.Call(printMilestone),
	)
}

func printMilestone(milestone birthday.Milestone) {
	astricks := " "
	if milestone.DaysAway == 0 {
		astricks = "*"
	}
	fmt.Printf(
		"%s %14s %20s %s\n",
		astricks,
		birthday.ToStringWithWeekDay(milestone.Date),
		milestone.AgeString(),
		milestone.EntryPtr.Name)
}

func init() {
	flag.StringVar(&fFile, "file", "", "Birthday file")
	flag.IntVar(&fDaysAhead, "days_ahead", 21, "Days ahead")
}
