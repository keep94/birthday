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

const (
	kInvalidPeriod = "invalid period"
)

var yearly = Period{Years: 1}

var defaultPeriods = []Period{
	{Years: 1},
	{Months: 100},
	{Months: 6, Normalize: true},
	{Weeks: 100},
	{Days: 1000},
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
		t, ok = safeYMD(0, month, day)
	} else {
		var year int
		year, err = strconv.Atoi(parts[2])
		if err != nil {
			return
		}
		t, ok = safeYMD(year, month, day)
	}
	if !ok {
		return time.Time{}, fmt.Errorf("Invalid date: %s", s)
	}
	return t, nil
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
	return asDays(end) - asDays(start)
}

// Entry represents a single entry in the birthday database
type Entry struct {
	Name     string
	Birthday time.Time
}

// Period represents a period of time
type Period struct {
	Years  int
	Months int
	Weeks  int
	Days   int

	// If true, Multiply normalizes.
	Normalize bool
}

// Valid returns true if p represents a net positive period.
func (p Period) Valid() bool {
	return p.approxDays() > 0.0
}

// Less orders Periods. Less orders first by Days, then by Weeks,
// then by Months, and finally by Years.
func (p Period) Less(other Period) bool {
	if p.Days < other.Days {
		return true
	}
	if p.Days > other.Days {
		return false
	}
	if p.Weeks < other.Weeks {
		return true
	}
	if p.Weeks > other.Weeks {
		return false
	}
	if p.Months < other.Months {
		return true
	}
	if p.Months > other.Months {
		return false
	}
	return p.Years < other.Years
}

// Diff returns the number of this period between end and start rounded down.
// Diff panics if this period is not valid.
func (p Period) Diff(end, start time.Time) int {
	diff := float64(DiffInDays(end, start))
	approxDays := p.approxDays()
	if approxDays <= 0.0 {
		panic(kInvalidPeriod)
	}
	result := int(diff / approxDays)
	for !p.Add(start, result).After(end) {
		result++
	}
	for p.Add(start, result).After(end) {
		result--
	}
	return result
}

// Add adds count of this period to start and returns the result.
func (p Period) Add(start time.Time, count int) time.Time {
	return start.AddDate(
		count*p.Years, count*p.Months, count*(p.Weeks*7+p.Days))
}

func (p Period) String() string {
	var parts []string
	if p.Years != 0 {
		parts = append(parts, fmt.Sprintf("%d years", p.Years))
	}
	if p.Months != 0 {
		parts = append(parts, fmt.Sprintf("%d months", p.Months))
	}
	if p.Weeks != 0 {
		parts = append(parts, fmt.Sprintf("%d weeks", p.Weeks))
	}
	if p.Days != 0 {
		parts = append(parts, fmt.Sprintf("%d days", p.Days))
	}
	if len(parts) == 0 {
		return "0 days"
	}
	return strings.Join(parts, " ")
}

// Multiply returns p * count. If p.Normalize is true, the returned period is
// normalized. The Normalize field of returned Period is set to false.
func (p Period) Multiply(count int) Period {
	var result Period
	result.Years = p.Years * count
	result.Months = p.Months * count
	result.Weeks = p.Weeks * count
	result.Days = p.Days * count
	if p.Normalize {
		result.normalize()
	}
	return result
}

func (p Period) approxDays() float64 {
	years := float64(p.Years) + float64(p.Months)/12.0
	days := 7.0*float64(p.Weeks) + float64(p.Days)
	return years*365.2425 + days
}

func (p *Period) normalize() {
	p.normalizeMonths()
	p.normalizeDays()
}

func (p *Period) normalizeMonths() {
	monthsOver12 := floorDiv(p.Months, 12)
	p.Years += monthsOver12
	p.Months -= 12 * monthsOver12
}

func (p *Period) normalizeDays() {
	daysOver7 := floorDiv(p.Days, 7)
	p.Weeks += daysOver7
	p.Days -= 7 * daysOver7
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
	Age Period

	// If true, age is unknown
	AgeUnknown bool
}

// Less orders Milestones. Less orders first by DaysAway then by Name
// then by AgeUnknown and finally by Age.
func (m *Milestone) Less(other *Milestone) bool {
	if m.DaysAway < other.DaysAway {
		return true
	}
	if m.DaysAway > other.DaysAway {
		return false
	}
	if m.Name < other.Name {
		return true
	}
	if m.Name > other.Name {
		return false
	}
	if !m.AgeUnknown && other.AgeUnknown {
		return true
	}
	if m.AgeUnknown && !other.AgeUnknown {
		return false
	}
	return m.Age.Less(other.Age)
}

// AgeString returns the age as a string e.g "57 years"
func (m *Milestone) AgeString() string {
	if m.AgeUnknown {
		return "? years"
	}
	return m.Age.String()
}

// Reminder reminds of upcoming milestones for people.
// Caller adds people with the Consume() method then the caller calls
// Milestones() to see all the people with upcoming milestones.
// Reminder implements the consume.Consumer interface and consumes Entry
// instances.
type Reminder struct {
	currentDate time.Time
	daysAhead   int
	periods     []Period
	milestones  []Milestone
}

// NewReminder creates a new Reminder instance. currentDate is the current
// date. daysAhead controls how many days in the future milestones can be.
// By default the new instance reminds of yearly birthdays, each 100 months,
// each 100 weeks, each 1000 days, and each 6 months. Caller can change this
// by calling SetPeriods.
func NewReminder(currentDate time.Time, daysAhead int) *Reminder {
	result := &Reminder{
		currentDate: currentDate,
		daysAhead:   daysAhead,
		periods:     defaultPeriods}
	return result
}

// SetPeriods sets the periods for which this reminder will remind overriding
// previously set periods. SetPeriods panics if any of the periods passed to
// it are invalid.
func (r *Reminder) SetPeriods(periods ...Period) {
	for _, p := range periods {
		if !p.Valid() {
			panic(kInvalidPeriod)
		}
	}
	r.periods = make([]Period, len(periods))
	copy(r.periods, periods)
}

// CanConsume always returns true
func (r *Reminder) CanConsume() bool {
	return true
}

// Consume consumes an Entry instance. ptr points to the Entry instance
// being consumed.
func (r *Reminder) Consume(ptr interface{}) {
	e := ptr.(*Entry)
	hasYear := HasYear(e.Birthday)
	for _, p := range r.periods {
		if hasYear || p == yearly {
			r.addPeriodMilestones(e, p)
		}
	}
}

// Milestones returns upcoming milestones for people consumed so far.
// Milestones happening soonest come first followed by milestones happening
// later.
func (r *Reminder) Milestones() []Milestone {
	result := make([]Milestone, len(r.milestones))
	copy(result, r.milestones)
	sort.Slice(
		result,
		func(i, j int) bool { return result[i].Less(&result[j]) })
	removeDuplicateMilestones(&result)
	return result
}

func removeDuplicateMilestones(milestones *[]Milestone) {
	if len(*milestones) == 0 {
		return
	}
	idx := 1
	for i := 1; i < len(*milestones); i++ {
		if (*milestones)[idx-1].Less(&(*milestones)[i]) {
			(*milestones)[idx] = (*milestones)[i]
			idx++
		}
	}
	*milestones = (*milestones)[:idx]
}

func (r *Reminder) addPeriodMilestones(e *Entry, period Period) {
	hasYear := HasYear(e.Birthday)
	yesterday := r.currentDate.AddDate(0, 0, -1)
	count := period.Diff(yesterday, e.Birthday) + 1
	if count < 0 {
		count = 0
	}
	nextMilestone := period.Add(e.Birthday, count)
	daysAway := DiffInDays(nextMilestone, r.currentDate)
	for daysAway < r.daysAhead {
		var age Period
		if hasYear {
			age = period.Multiply(count)
		}
		r.milestones = append(r.milestones, Milestone{
			Name:       e.Name,
			Date:       nextMilestone,
			DaysAway:   daysAway,
			Age:        age,
			AgeUnknown: !hasYear,
		})
		count++
		nextMilestone = period.Add(e.Birthday, count)
		daysAway = DiffInDays(nextMilestone, r.currentDate)
	}
}

// EntryConsumer implements the consume.Consumer interface and consumes
// Entry instances.
type EntryConsumer struct {
	entries []Entry
}

// CanConsume always returns true
func (e *EntryConsumer) CanConsume() bool {
	return true
}

// Consume consumes an Entry instance. ptr points to the Entry instance
// being consumed.
func (e *EntryConsumer) Consume(ptr interface{}) {
	p := ptr.(*Entry)
	e.entries = append(e.entries, *p)
}

// Entries returns the consumed Entry instances sorted by Name in ascending
// order.
func (s *EntryConsumer) Entries() []Entry {
	result := make([]Entry, len(s.entries))
	copy(result, s.entries)
	sort.SliceStable(
		result,
		func(i, j int) bool { return result[i].Name < result[j].Name },
	)
	return result
}

// Query returns a function that returns true if the Entry instance passed
// to it matches query.
func Query(query string) func(entry *Entry) bool {
	query = str_util.Normalize(query)
	return func(entry *Entry) bool {
		return strings.Contains(str_util.Normalize(entry.Name), query)
	}
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

func safeYMD(year, month, day int) (t time.Time, ok bool) {
	result := date_util.YMD(year, month, day)
	y, m, d := result.Date()
	if y != year || int(m) != month || d != day {
		return
	}
	return result, true
}

func asDays(t time.Time) int {
	unix := t.Unix()
	days := int(unix / 86400)
	seconds := int(unix % 86400)
	if seconds < 0 {
		days--
	}
	return days
}
