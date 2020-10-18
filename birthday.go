package birthday

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var marToFebOffsets = []int{
	0, 31, 61, 92, 122, 153, 184, 214, 245, 275, 306, 337}

var monthLengths = []int{
	31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

var daysOfWeek = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// Birthday represents a birthday or date
type Birthday struct {

	// Year may be 0 if unknown
	Year int

	// 1=January, 2=February etc.
	Month int

	Day int
}

// String returns this instance as 'mm/dd' or 'mm/dd/yyyy'
func (b Birthday) String() string {
	if !b.YearSet() {
		return fmt.Sprintf("%02d/%02d", b.Month, b.Day)
	}
	return fmt.Sprintf("%02d/%02d/%04d", b.Month, b.Day, b.Year)
}

// StringWithWeekDay returns this instance as a string preceded by
// the day of the week e.g "Thu 10/15/2020" StringWithWeekDay panics
// if it can't determine the day of the week such as with a birthday
// that doesn't have a year.
func (b Birthday) StringWithWeekDay() string {
	dayOfWeek := b.AsDays() % 7
	return fmt.Sprintf("%s %s", daysOfWeek[dayOfWeek], b.String())
}

// Now returns today.
func Now() Birthday {
	t := time.Now()
	return Birthday{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}
}

// Parse converts a string of form 'mm/dd' or 'mm/dd/yyyy' into a Birthday
// Returns an error if string cannot be parsed.
func Parse(s string) (birthday Birthday, err error) {
	parts := strings.Split(s, "/")
	if len(parts) < 2 || len(parts) > 3 {
		return Birthday{}, errors.New("must be of form mm/dd or mm/dd/yyyy")
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return
	}
	day, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	var b Birthday
	if len(parts) == 2 {
		b = Birthday{Month: month, Day: day}
	} else {
		var year int
		year, err = strconv.Atoi(parts[2])
		if err != nil {
			return
		}
		b = Birthday{Year: year, Month: month, Day: day}
	}
	birthday = b
	return
}

// YearSet returns true if Year field is set.
func (b Birthday) YearSet() bool {
	return b.Year > 0
}

// IsValid returns true if this birthday is valid. That is the Month and Day
// fields are within range.
func (b Birthday) IsValid() bool {
	if b.Month < 1 || b.Month > 12 {
		return false
	}
	daysInMonth := monthLengths[b.Month-1]
	if b.Month == 2 && (b.Year <= 0 || isLeapYear(b.Year)) {
		daysInMonth += 1
	}
	return b.Day >= 1 && b.Day <= daysInMonth
}

// Normalize returns an equivalent valid Birthday for this Birthday. For
// example 08/32/2020 -> 09/01/2020. Normalize panics if this Birthday
// represents a day before 1 Jan 0001.
func (b Birthday) Normalize() Birthday {
	return FromDays(b.AsDays())
}

// AsDays returns this birthday as days since 31 Dec 0000.
// AsDays panics if this birthday falls before 1 Jan 0001.
// AsDays() % 7 gives the day of the week. 0=Sunday, 1=Monday, etc.
func (b Birthday) AsDays() int {
	years := b.Month / 12
	months := b.Month % 12
	ymNormalized := Birthday{Year: b.Year + years, Month: months, Day: b.Day}
	if ymNormalized.Month <= 0 {
		ymNormalized.Month += 12
		ymNormalized.Year--
	}
	return ymNormalized.asDays()
}

func (b Birthday) asDays() int {
	if !b.YearSet() {
		panic("normalized year is 0 or less.")
	}
	year := b.Year
	monthIndex := b.Month - 3
	if monthIndex < 0 {
		year--
		monthIndex += 12
	}
	leapDays := (year / 4) - (year / 100) + (year / 400)
	result := 365*year + leapDays
	result += marToFebOffsets[monthIndex] + (b.Day - 1)

	// Make it so that 1 Jan 0001 is day 1. That day happens to be Monday
	result -= 305
	if result <= 0 {
		panic("birthday falls before 1 Jan 0001")
	}
	return result
}

// FromDays is the inverse of AsDays. FromDays always returns a normalized
// birthday. FromDays panics if days <= 0.
func FromDays(days int) Birthday {
	if days <= 0 {
		panic("Days must be at least 1")
	}
	days += 305
	year400 := days / 146097
	days -= year400 * 146097
	year100 := days / 36524
	if year100 == 4 {
		year100 = 3
	}
	days -= year100 * 36524
	year4 := days / 1461
	days -= year4 * 1461
	year1 := days / 365
	if year1 == 4 {
		year1 = 3
	}
	days -= year1 * 365
	year := year400*400 + year100*100 + year4*4 + year1
	monthIndex := 11
	for marToFebOffsets[monthIndex] > days {
		monthIndex--
	}
	day := days - marToFebOffsets[monthIndex]
	if monthIndex < 10 {
		return Birthday{Year: year, Month: monthIndex + 3, Day: day + 1}
	}
	return Birthday{Year: year + 1, Month: monthIndex - 9, Day: day + 1}
}

// Milestone represents a milestone day.
type Milestone struct {

	// The name of the person having the milestone
	Name string

	// The date of the milestone day
	Date Birthday

	// How many days in the future this milestone day is.
	DaysAway int

	// The age of the person on this mileestone day in years or days.
	Age int

	// Set to true if age is in days
	AgeInDays bool
}

// Remind reminds of upcoming milestones for people.
// Caller adds people with the Add() method then the caller calls Reminders()
// to see all the people with upcoming milestones.
type Remind struct {
	currentDate Birthday
	currentDay  int
	daysAhead   int
	milestones  []Milestone
}

// NewRemind creates a new Remind instance. currentDate is the current date.
// daysAhead controls how many days in the future milestones can be.
func NewRemind(currentDate Birthday, daysAhead int) *Remind {
	return &Remind{
		currentDate: currentDate.Normalize(),
		currentDay:  currentDate.AsDays(),
		daysAhead:   daysAhead}
}

// Add adds a person. Add panics if b is invalid.
func (r *Remind) Add(name string, b Birthday) {
	if !b.IsValid() {
		panic("b must be valid")
	}
	r.addYearMilestones(name, b)
	if b.YearSet() {
		r.addDayMilestones(name, b)
	}
}

// Reminders returns upcoming reminders for people added so far. Milestones
// happening soonest come first followed by milestones happining later.
func (r *Remind) Reminders() []Milestone {
	result := make([]Milestone, len(r.milestones))
	copy(result, r.milestones)
	sort.SliceStable(
		result,
		func(i, j int) bool { return result[i].DaysAway < result[j].DaysAway },
	)
	return result
}

func (r *Remind) addYearMilestones(name string, b Birthday) {
	nextMilestone := b
	if !nextMilestone.YearSet() || nextMilestone.AsDays() < r.currentDay {
		nextMilestone.Year = r.currentDate.Year
	}
	if nextMilestone.AsDays() < r.currentDay {
		nextMilestone.Year++
	}
	daysAway := nextMilestone.AsDays() - r.currentDay
	for daysAway < r.daysAhead {
		age := -1
		if b.YearSet() {
			age = nextMilestone.Year - b.Year
		}
		r.milestones = append(r.milestones, Milestone{
			Name:     name,
			Date:     nextMilestone.Normalize(),
			DaysAway: daysAway,
			Age:      age,
		})
		nextMilestone.Year++
		daysAway = nextMilestone.AsDays() - r.currentDay
	}
}

func (r *Remind) addDayMilestones(name string, b Birthday) {
	bAsDays := b.AsDays()
	nextMilestoneAsDays := bAsDays
	if nextMilestoneAsDays < r.currentDay {
		diff := r.currentDay - nextMilestoneAsDays
		nextMilestoneAsDays += ((diff + 999) / 1000) * 1000
	}
	for nextMilestoneAsDays-r.currentDay < r.daysAhead {
		r.milestones = append(r.milestones, Milestone{
			Name:      name,
			Date:      FromDays(nextMilestoneAsDays),
			DaysAway:  nextMilestoneAsDays - r.currentDay,
			Age:       nextMilestoneAsDays - bAsDays,
			AgeInDays: true,
		})
		nextMilestoneAsDays += 1000
	}
}

func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	if year%4 == 0 {
		return true
	}
	return false
}
