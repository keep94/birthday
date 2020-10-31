// Package birthday contains routines for tracking birthdays.
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

// HasYear returns true if t has a year. That is t falls on or after
// 1 Jan 0001
func HasYear(t time.Time) bool {
	return t.Year() > 0
}

// DiffInYears returns the number of years between start and end rounded down.
func DiffInYears(end, start time.Time) int {
	return floorDiv(DiffInMonths(end, start), 12)
}

// DiffInMonths returns the number of months between start and end rounded
// down.
func DiffInMonths(end, start time.Time) int {
	syear, smonth, sday := start.Date()
	end = end.AddDate(0, 0, 1-sday)
	eyear, emonth, _ := end.Date()
	return (eyear-syear)*12 + int(emonth) - int(smonth)
}

// DiffInWeeks returns the number of weeks between start and end rounded down
func DiffInWeeks(end, start time.Time) int {
	return floorDiv(DiffInDays(end, start), 7)
}

// DiffInDays returns the number of days between start and end rounded down
func DiffInDays(end, start time.Time) int {
	return AsDays(end) - AsDays(start)
}

// Entry represents a single entry in the birthday database
type Entry struct {
	Name     string
	Birthday time.Time
}

// Unit is unit of time.
type Unit int

const (
	Years Unit = iota
	Months
	Weeks
	Days
)

func (u Unit) String() string {
	switch u {
	case Years:
		return "years"
	case Months:
		return "months"
	case Weeks:
		return "weeks"
	case Days:
		return "days"
	default:
		return "unknown"
	}
}

// Diff returns the number of this unit between start and end rounded down.
func (u Unit) Diff(end, start time.Time) int {
	switch u {
	case Years:
		return DiffInYears(end, start)
	case Months:
		return DiffInMonths(end, start)
	case Weeks:
		return DiffInWeeks(end, start)
	case Days:
		return DiffInDays(end, start)
	default:
		panic("unknown unit")
	}
}

// Add returns start plus x of this unit.
func (u Unit) Add(start time.Time, x int) time.Time {
	switch u {
	case Years:
		return start.AddDate(x, 0, 0)
	case Months:
		return start.AddDate(0, x, 0)
	case Weeks:
		return start.AddDate(0, 0, 7*x)
	case Days:
		return start.AddDate(0, 0, x)
	default:
		panic("unknown unit")
	}
}

// Milestone represents a milestone day.
type Milestone struct {

	// The name of the person having the milestone
	Name string

	// The date of the milestone day
	Date time.Time

	// How many days in the future this milestone day is.
	DaysAway int

	// The age of the person on this mileestone day.
	// -1 if age of person is unknown.
	Age int

	// The age units of the person. e.g Years, Weeks, Days, etc.
	Unit Unit
}

// Types contains what milestone types a Reminder instance will give
type Types struct {
	Years         bool
	HundredMonths bool
	HundredWeeks  bool
	ThousandDays  bool
}

// Flip returns the opposite of this instance.
func (t Types) Flip() Types {
	return Types{
		Years:         !t.Years,
		HundredMonths: !t.HundredMonths,
		HundredWeeks:  !t.HundredWeeks,
		ThousandDays:  !t.ThousandDays,
	}
}

// Reminder reminds of upcoming milestones for people.
// Caller adds people with the Consume() method then the caller calls
// Milestones() to see all the people with upcoming milestones.
type Reminder struct {
	currentDate time.Time
	daysAhead   int
	types       Types
	milestones  []Milestone
}

// NewReminder creates a new Reminder instance. currentDate is the current
// date. daysAhead controls how many days in the future milestones can be.
func NewReminder(currentDate time.Time, daysAhead int) *Reminder {
	result := &Reminder{
		currentDate: currentDate,
		daysAhead:   daysAhead,
		types:       Types{}.Flip()}
	return result
}

// SetTypes sets the milestone types this instance will give. The default
// is all types.
func (r *Reminder) SetTypes(types Types) {
	r.types = types
}

// Consume consumes an entry.
func (r *Reminder) Consume(e *Entry) {
	if r.types.Years {
		r.addUnitMilestones(e, 1, Years)
	}
	if !HasYear(e.Birthday) {
		return
	}
	if r.types.HundredMonths {
		r.addUnitMilestones(e, 100, Months)
	}
	if r.types.HundredWeeks {
		r.addUnitMilestones(e, 100, Weeks)
	}
	if r.types.ThousandDays {
		r.addUnitMilestones(e, 1000, Days)
	}
}

// Milestones returns upcoming milestones for people consumed so far.
// Milestones happening soonest come first followed by milestones happining
// later.
func (r *Reminder) Milestones() []Milestone {
	result := make([]Milestone, len(r.milestones))
	copy(result, r.milestones)
	sort.SliceStable(
		result,
		func(i, j int) bool { return result[i].DaysAway < result[j].DaysAway },
	)
	return result
}

func (r *Reminder) addUnitMilestones(e *Entry, num int, unit Unit) {
	hasYear := HasYear(e.Birthday)
	yesterday := r.currentDate.AddDate(0, 0, -1)
	nextAge := floorDiv(unit.Diff(yesterday, e.Birthday), num)*num + num
	if nextAge < 0 {
		nextAge = 0
	}
	nextMilestone := unit.Add(e.Birthday, nextAge)
	daysAway := DiffInDays(nextMilestone, r.currentDate)
	for daysAway < r.daysAhead {
		age := -1
		if hasYear {
			age = nextAge
		}
		r.milestones = append(r.milestones, Milestone{
			Name:     e.Name,
			Date:     nextMilestone,
			DaysAway: daysAway,
			Age:      age,
			Unit:     unit,
		})
		nextAge += num
		nextMilestone = unit.Add(e.Birthday, nextAge)
		daysAway = DiffInDays(nextMilestone, r.currentDate)
	}
}

// Search searches for people by name
type Search struct {
	query   string
	entries []Entry
}

// NewSearch returns a new search. query is a search string. Searches ignore
// case and extra whitespace.
func NewSearch(query string) *Search {
	return &Search{query: str_util.Normalize(query)}
}

// Consume consumes an entry.
func (s *Search) Consume(e *Entry) {
	if strings.Contains(str_util.Normalize(e.Name), s.query) {
		s.entries = append(s.entries, *e)
	}
}

// Results returns the results that match the query string for this instance
// sorted by name.
func (s *Search) Results() []Entry {
	result := make([]Entry, len(s.entries))
	copy(result, s.entries)
	sort.SliceStable(
		result,
		func(i, j int) bool { return result[i].Name < result[j].Name },
	)
	return result
}

func floorDiv(x, positiveY int) int {
	if positiveY <= 0 {
		panic("positiveY must be positive")
	}
	result := x / positiveY
	if x%positiveY < 0 {
		result--
	}
	return result
}
