// Package birthday contains routines for tracking birthdays.
package birthday

import (
	"container/heap"
	"errors"
	"fmt"
	"iter"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/keep94/consume2"
	"github.com/keep94/itertools"
	"github.com/keep94/toolbox/date_util"
	"github.com/keep94/toolbox/str_util"
)

const (
	kInvalidPeriod = "invalid period"
)

var (
	// Currently yearly, 100 months, 100 weeks, 1000 days.
	DefaultPeriods = []Period{
		{Years: 1},
		{Months: 100},
		{Weeks: 100},
		{Days: 1000},
	}
)

var yearly = Period{Years: 1}

// Today returns today's date at midnight in UTC.
func Today(clock date_util.Clock) time.Time {
	y, m, d := clock.Now().Date()
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

// ToStringWithWeekDay works like ToString but adds weekday.
// ToStringWithWeekDay returns a string such as 'Mon 01/02/2006'.
// ToStringWithWeekDay panics if t falls before 1 Jan 0001.
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

// Entry represents a single entry in the birthday database
type Entry struct {
	Name     string
	Birthday time.Time
}

// EntriesSortedByName returns entries sorted by name while leaving the
// original entries slice unchanged.
func EntriesSortedByName(entries []*Entry) []*Entry {
	result := make([]*Entry, len(entries))
	copy(result, entries)
	sort.SliceStable(
		result,
		func(i, j int) bool { return result[i].Name < result[j].Name },
	)
	return result
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
	diff := float64(asDays(end) - asDays(start))
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
	monthsOver12 := p.Months / 12
	p.Years += monthsOver12
	p.Months -= 12 * monthsOver12
}

func (p *Period) normalizeDays() {
	daysOver7 := p.Days / 7
	p.Weeks += daysOver7
	p.Days -= 7 * daysOver7
}

// Milestone represents a milestone day.
type Milestone struct {

	// The person having the milestone
	EntryPtr *Entry

	// The date of the milestone day
	Date time.Time

	// The age of the person on this milestone day
	Age Period

	// If true, age is unknown
	AgeUnknown bool
}

// Less orders Milestones. Less orders first by Date then by Name
// then by AgeUnknown and finally by Age.
func (m *Milestone) Less(other *Milestone) bool {
	if m.Date.Before(other.Date) {
		return true
	}
	if m.Date.After(other.Date) {
		return false
	}
	if m.EntryPtr.Name < other.EntryPtr.Name {
		return true
	}
	if m.EntryPtr.Name > other.EntryPtr.Name {
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

// Query returns a function that returns true if the Entry instance passed
// to it matches query.
func Query(query string) func(entry Entry) bool {
	query = str_util.Normalize(query)
	if query == "" {
		return consume2.ComposeFilters[Entry]()
	}
	return func(entry Entry) bool {
		return strings.Contains(str_util.Normalize(entry.Name), query)
	}
}

// Remind returns all upcoming Milestones for the specified entries and
// periods starting at the date specified by current. Remind returns
// Milestone instances in chronological order.
func Remind(
	entries []*Entry,
	periods []Period,
	current time.Time) iter.Seq[Milestone] {
	checkPeriods(periods)
	base := createMilestoneBase(entries, periods, current)
	if len(base) == 0 {
		return itertools.Chain[Milestone]()
	}
	return func(yield func(Milestone) bool) {
		mh := createMilestoneHeap(base)
		milestone := mh[0].Milestone
		for {
			if !yield(milestone) {
				return
			}
			for !milestone.Less(&mh[0].Milestone) {
				mh[0].Advance()
				heap.Fix(&mh, 0)
			}
			milestone = mh[0].Milestone
		}
	}
}

// RemindPtrs works like Remind except that it returns Milestone pointers.
func RemindPtrs(
	entries []*Entry,
	periods []Period,
	current time.Time) iter.Seq[*Milestone] {
	return itertools.Map(
		func(m Milestone) *Milestone { return &m },
		Remind(entries, periods, current))
}

func createMilestoneBase(
	entries []*Entry,
	periods []Period,
	current time.Time) []milestoneGenerator {
	size := 0
	for i := range entries {
		hasYear := HasYear(entries[i].Birthday)
		for j := range periods {
			if hasYear || periods[j] == yearly {
				size++
			}
		}
	}
	if size == 0 {
		return nil
	}
	result := make([]milestoneGenerator, size)
	index := 0
	for i := range entries {
		hasYear := HasYear(entries[i].Birthday)
		for j := range periods {
			if hasYear || periods[j] == yearly {
				result[index].Init(entries[i], periods[j], current)
				index++
			}
		}
	}
	return result
}

func createMilestoneHeap(base []milestoneGenerator) milestoneHeap {
	allocatedSpace := slices.Clone(base)
	result := make(milestoneHeap, len(allocatedSpace))
	for i := range result {
		result[i] = &allocatedSpace[i]
	}
	heap.Init(&result)
	return result
}

type milestoneGenerator struct {
	Generator generator
	Milestone Milestone
}

func (gm *milestoneGenerator) Init(
	entry *Entry, period Period, current time.Time) {
	gm.Generator.Init(entry, period, current)
	gm.Advance()
}

func (gm *milestoneGenerator) Advance() {
	gm.Milestone = gm.Generator.Next()
}

type generator struct {
	entryPtr *Entry
	period   Period
	count    int
}

func (g *generator) Init(
	entry *Entry, period Period, current time.Time) {
	yesterday := current.AddDate(0, 0, -1)
	count := period.Diff(yesterday, entry.Birthday) + 1
	if count < 0 {
		count = 0
	}
	*g = generator{entryPtr: entry, period: period, count: count}
}

func (g *generator) Next() Milestone {
	hasYear := HasYear(g.entryPtr.Birthday)
	var age Period
	if hasYear {
		age = g.period.Multiply(g.count)
	}
	nextMilestone := g.period.Add(g.entryPtr.Birthday, g.count)
	result := Milestone{
		EntryPtr:   g.entryPtr,
		Date:       nextMilestone,
		Age:        age,
		AgeUnknown: !hasYear,
	}
	g.count++
	return result
}

type milestoneHeap []*milestoneGenerator

func (m milestoneHeap) Less(i, j int) bool {
	return m[i].Milestone.Less(&m[j].Milestone)
}

func (m milestoneHeap) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m milestoneHeap) Len() int {
	return len(m)
}

func (m *milestoneHeap) Push(x interface{}) {
	mg := x.(*milestoneGenerator)
	*m = append(*m, mg)
}

func (m *milestoneHeap) Pop() interface{} {
	old := *m
	n := len(old)
	x := old[n-1]
	*m = old[0 : n-1]
	return x
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

func checkPeriods(periods []Period) {
	for _, p := range periods {
		if !p.Valid() {
			panic(kInvalidPeriod)
		}
	}
}
