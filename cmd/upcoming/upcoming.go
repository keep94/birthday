package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/keep94/birthday"
	"github.com/keep94/consume2"
	"github.com/keep94/itertools"
	"github.com/keep94/toolbox/date_util"
)

const (
	kMaxRows = 100
)

var (
	fFile      string
	fDaysAhead int
)

var (
	kClock date_util.Clock = date_util.SystemClock{}
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
	today := birthday.Today(kClock)
	endTime := today.AddDate(0, 0, fDaysAhead)
	seq := birthday.RemindPtrs(entries, birthday.DefaultPeriods, today)
	seq = itertools.TakeWhile(
		func(m *birthday.Milestone) bool { return m.Date.Before(endTime) },
		seq)
	seq = itertools.Take(kMaxRows, seq)
	for milestonePtr := range seq {
		printMilestone(milestonePtr, today)
	}
}

func printMilestone(milestone *birthday.Milestone, today time.Time) {
	astricks := " "
	if milestone.Date.Equal(today) {
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
