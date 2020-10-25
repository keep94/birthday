package birthday

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/keep94/toolbox/date_util"
	"github.com/keep94/toolbox/str_util"
)

// SafeYMD works like YMD except it returns false if year, month, and day
// aren't valid.
func SafeYMD(year, month, day int) (t time.Time, ok bool) {
	result := date_util.YMD(year, month, day)
	y, m, d := result.Date()
	if y != year || int(m) != month || d != day {
		return
	}
	return result, true
}

// Today returns today's date at midnight in UTC.
func Today() time.Time {
	y, m, d := time.Now().Date()
	return date_util.YMD(y, int(m), d)
}

// ToString returns t as MM/dd/yyyy or as just MM/dd if t falls before
// 1 Jan 0001.
func ToString(t time.Time) string {
	if !HasYear(t) {
		return t.Format("01/02")
	}
	return t.Format("01/02/2006")
}

// ToStringWithWeekday works like ToString but adds weekday.
// ToStringWithWeekday panics if t falls before 1 Jan 0001.
func ToStringWithWeekDay(t time.Time) string {
	if !HasYear(t) {
		panic("no year")
	}
	return t.Format("Mon 01/02/2006")
}

// Parse converts s to a time in UTC. s must be of form MM/dd/yyyy or
// MM/dd.  If s is of form MM/dd, the year of returned time is 0.
// s must be a valid date as no normalizing is done.  Invalid dates like
// '08/32/2006' return an error.
func Parse(s string) (parsed time.Time, err error) {
	parts := strings.Split(s, "/")
	if len(parts) < 2 || len(parts) > 3 {
		return time.Time{}, errors.New("must be of form mm/dd or mm/dd/yyyy")
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return
	}
	day, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	var t time.Time
	var ok bool
	if len(parts) == 2 {
		t, ok = SafeYMD(0, month, day)
	} else {
		var year int
		year, err = strconv.Atoi(parts[2])
		if err != nil {
			return
		}
		t, ok = SafeYMD(year, month, day)
	}
	if !ok {
		return time.Time{}, fmt.Errorf("Invalid date: %s", s)
	}
	return t, nil
}

// AsDays returns the day number for t. 0 is 1 Jan 1970; 1 is 2 Jan 1970 etc.
func AsDays(t time.Time) int {
	unix := t.Unix()
	days := int(unix / 86400)
	seconds := int(unix % 86400)
	if seconds < 0 {
		days--
	}
	return days
}

// FromDays converts a day number to a time. Returned time is always in UTC.
func FromDays(days int) time.Time {
	return time.Unix(int64(days)*86400, 0).UTC()
}

// HasYear returns true if t has a year. That is t falls on or after
// 1 Jan 0001
func HasYear(t time.Time) bool {
	return t.Year() > 0
}

// DiffInYears returns the number of years between start and end rounded down.
func DiffInYears(end, start time.Time) int {
	eyear, emonth, eday := end.Date()
	syear, smonth, sday := start.Date()
	diff := eyear - syear
	if emonth < smonth || (emonth == smonth && eday < sday) {
		diff--
	}
	return diff
}

// Milestone represents a milestone day.
type Milestone struct {

	// The name of the person having the milestone
	Name string

	// The date of the milestone day
	Date time.Time

	// How many days in the future this milestone day is.
	DaysAway int

	// The age of the person on this mileestone day in years or days.
	// -1 if age of person is unknown.
	Age int

	// Set to true if age is in days
	AgeInDays bool
}

// Remind reminds of upcoming milestones for people.
// Caller adds people with the Add() method then the caller calls Reminders()
// to see all the people with upcoming milestones.
type Remind struct {
	currentDate time.Time
	daysAhead   int
	milestones  []Milestone
}

// NewRemind creates a new Remind instance. currentDate is the current date.
// daysAhead controls how many days in the future milestones can be.
func NewRemind(currentDate time.Time, daysAhead int) *Remind {
	return &Remind{
		currentDate: currentDate,
		daysAhead:   daysAhead}
}

// Add adds a person.
func (r *Remind) Add(name string, bday time.Time) {
	r.addYearMilestones(name, bday)
	if HasYear(bday) {
		r.addDayMilestones(name, bday)
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

func (r *Remind) addYearMilestones(name string, bday time.Time) {
	hasYear := HasYear(bday)
	nextAge := DiffInYears(r.currentDate.AddDate(0, 0, -1), bday) + 1
	if nextAge < 0 {
		nextAge = 0
	}
	nextMilestone := bday.AddDate(nextAge, 0, 0)
	daysAway := AsDays(nextMilestone) - AsDays(r.currentDate)
	for daysAway < r.daysAhead {
		age := -1
		if hasYear {
			age = nextAge
		}
		r.milestones = append(r.milestones, Milestone{
			Name:     name,
			Date:     nextMilestone,
			DaysAway: daysAway,
			Age:      age,
		})
		nextAge++
		nextMilestone = bday.AddDate(nextAge, 0, 0)
		daysAway = AsDays(nextMilestone) - AsDays(r.currentDate)
	}
}

func (r *Remind) addDayMilestones(name string, bday time.Time) {
	bAsDays := AsDays(bday)
	nextMilestoneAsDays := bAsDays
	currentDay := AsDays(r.currentDate)
	if nextMilestoneAsDays < currentDay {
		diff := currentDay - nextMilestoneAsDays
		nextMilestoneAsDays += ((diff + 999) / 1000) * 1000
	}
	for nextMilestoneAsDays-currentDay < r.daysAhead {
		r.milestones = append(r.milestones, Milestone{
			Name:      name,
			Date:      FromDays(nextMilestoneAsDays),
			DaysAway:  nextMilestoneAsDays - currentDay,
			Age:       nextMilestoneAsDays - bAsDays,
			AgeInDays: true,
		})
		nextMilestoneAsDays += 1000
	}
}

// Person represents a person
type Person struct {
	// Name of person
	Name string

	// Birthday of person
	Birthday time.Time

	// Age in years. 0 if age unknown
	AgeInYears int

	// Age in days. 0 if age unknown
	AgeInDays int
}

// Filter filters people by name
type Filter struct {
	currentDate time.Time
	query       string
	persons     []Person
}

// NewFilter returns a new Filter. query is a search string. Searches ignore
// case and extra whitespace.
func NewFilter(currentDate time.Time, query string) *Filter {
	return &Filter{currentDate: currentDate, query: str_util.Normalize(query)}
}

// Add adds a person. If name does not match the query of this instance,
// add does nothing.
func (f *Filter) Add(name string, birthday time.Time) {
	if strings.Contains(str_util.Normalize(name), f.query) {
		ageInYears := 0
		ageInDays := 0
		if HasYear(birthday) {
			ageInYears = DiffInYears(f.currentDate, birthday)
			ageInDays = AsDays(f.currentDate) - AsDays(birthday)
		}
		f.persons = append(f.persons, Person{
			Name:       name,
			Birthday:   birthday,
			AgeInYears: ageInYears,
			AgeInDays:  ageInDays,
		})
	}
}

// Persons returns the people that match the query string for this instance
// sorted by name.
func (f *Filter) Persons() []Person {
	result := make([]Person, len(f.persons))
	copy(result, f.persons)
	sort.SliceStable(
		result,
		func(i, j int) bool { return result[i].Name < result[j].Name },
	)
	return result
}
