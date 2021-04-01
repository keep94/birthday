package birthday_test

import (
	"testing"
	"time"

	"github.com/keep94/birthday"
	"github.com/keep94/toolbox/date_util"
	asserts "github.com/stretchr/testify/assert"
)

var years = birthday.Period{Years: 1}
var hundredMonths = birthday.Period{Months: 100}
var hundredWeeks = birthday.Period{Weeks: 100}
var thousandDays = birthday.Period{Days: 1000}
var sixMonths = birthday.Period{Months: 6, Normalize: true}

func TestDiffInDaysAndWeeks(t *testing.T) {
	assert := asserts.New(t)
	start := date_util.YMD(2020, 2, 29)
	end := date_util.YMD(2020, 3, 1)
	assert.Equal(1, birthday.DiffInDays(end, start))
	assert.Equal(0, birthday.DiffInWeeks(end, start))
	start = date_util.YMD(2001, 9, 17)
	end = date_util.YMD(2018, 8, 5)
	assert.Equal(6166, birthday.DiffInDays(end, start))
	assert.Equal(880, birthday.DiffInWeeks(end, start))
	assert.Equal(-6166, birthday.DiffInDays(start, end))
	assert.Equal(-881, birthday.DiffInWeeks(start, end))
}

func TestDiffInMonths(t *testing.T) {
	assert := asserts.New(t)
	start := date_util.YMD(2019, 12, 31)
	end := date_util.YMD(2021, 3, 3)
	assert.Equal(14, birthday.DiffInMonths(end, start))
	end = date_util.YMD(2021, 3, 2)
	assert.Equal(13, birthday.DiffInMonths(end, start))
	end = date_util.YMD(2019, 12, 31)
	assert.Equal(0, birthday.DiffInMonths(end, start))
	end = date_util.YMD(2019, 12, 30)
	assert.Equal(-1, birthday.DiffInMonths(end, start))
	end = date_util.YMD(1983, 5, 26)
	start = date_util.YMD(1971, 11, 26)
	assert.Equal(138, birthday.DiffInMonths(end, start))
	end = date_util.YMD(1971, 11, 26)
	start = date_util.YMD(1983, 5, 26)
	assert.Equal(-138, birthday.DiffInMonths(end, start))
	end = date_util.YMD(1971, 11, 25)
	start = date_util.YMD(1983, 5, 26)
	assert.Equal(-139, birthday.DiffInMonths(end, start))
}

func TestDiffInYears(t *testing.T) {
	assert := asserts.New(t)
	end := date_util.YMD(1951, 2, 15)
	assert.Equal(0, birthday.DiffInYears(end, date_util.YMD(1951, 2, 15)))
	assert.Equal(-1, birthday.DiffInYears(end, date_util.YMD(1951, 2, 16)))
	assert.Equal(-1, birthday.DiffInYears(end, date_util.YMD(1951, 3, 1)))
	assert.Equal(-1, birthday.DiffInYears(end, date_util.YMD(1952, 2, 15)))
	assert.Equal(-2, birthday.DiffInYears(end, date_util.YMD(1952, 2, 16)))
	assert.Equal(0, birthday.DiffInYears(end, date_util.YMD(1951, 2, 14)))
	assert.Equal(0, birthday.DiffInYears(end, date_util.YMD(1951, 1, 31)))
	assert.Equal(0, birthday.DiffInYears(end, date_util.YMD(1950, 2, 16)))
	assert.Equal(1, birthday.DiffInYears(end, date_util.YMD(1950, 2, 15)))
	assert.Equal(3, birthday.DiffInYears(end, date_util.YMD(1948, 2, 15)))
}

func TestToString(t *testing.T) {
	assert := asserts.New(t)
	b := date_util.YMD(1992, 7, 4)
	assert.Equal("07/04/1992", birthday.ToString(b))
	b = date_util.YMD(953, 11, 30)
	assert.Equal("11/30/0953", birthday.ToString(b))
	b = date_util.YMD(0, 2, 29)
	assert.Equal("02/29", birthday.ToString(b))
	b = date_util.YMD(0, 12, 31)
	assert.Equal("12/31", birthday.ToString(b))
	b = date_util.YMD(0, 1, 1)
	assert.Equal("01/01", birthday.ToString(b))
}

func TestStringWithWeekday(t *testing.T) {
	assert := asserts.New(t)
	b := date_util.YMD(2020, 10, 15)
	assert.Equal("Thu 10/15/2020", birthday.ToStringWithWeekDay(b))
	b = date_util.YMD(0, 4, 2)
	assert.Panics(func() { birthday.ToStringWithWeekDay(b) })
}

func TestParse(t *testing.T) {
	assert := asserts.New(t)
	b, err := birthday.Parse("12/31")
	assert.NoError(err)
	assert.Equal(date_util.YMD(0, 12, 31), b)
	b, err = birthday.Parse("2/29")
	assert.NoError(err)
	assert.Equal(date_util.YMD(0, 2, 29), b)
	b, err = birthday.Parse("1/1")
	assert.NoError(err)
	assert.Equal(date_util.YMD(0, 1, 1), b)
	b, err = birthday.Parse("3/28/2017")
	assert.NoError(err)
	assert.Equal(date_util.YMD(2017, 3, 28), b)
	_, err = birthday.Parse("wrong")
	assert.Error(err)
	_, err = birthday.Parse("4/2/3/1")
	assert.Error(err)
	_, err = birthday.Parse("wrong/2")
	assert.Error(err)
	_, err = birthday.Parse("2/wrong")
	assert.Error(err)
	_, err = birthday.Parse("5/31/wrong")
	assert.Error(err)
	_, err = birthday.Parse("4/31")
	assert.Error(err)
	_, err = birthday.Parse("4/31/2017")
	assert.Error(err)
}

func TestMilestoneAgeString(t *testing.T) {
	assert := asserts.New(t)
	milestone := birthday.Milestone{AgeUnknown: true}
	assert.Equal("? years", milestone.AgeString())
	milestone = birthday.Milestone{Age: birthday.Period{Years: 47}}
	assert.Equal("47 years", milestone.AgeString())
	milestone = birthday.Milestone{Age: birthday.Period{Months: 5}}
	assert.Equal("5 months", milestone.AgeString())
	milestone = birthday.Milestone{Age: birthday.Period{Weeks: 3}}
	assert.Equal("3 weeks", milestone.AgeString())
	milestone = birthday.Milestone{Age: birthday.Period{Days: 1}}
	assert.Equal("1 days", milestone.AgeString())
	milestone = birthday.Milestone{}
	assert.Equal("0 days", milestone.AgeString())
	milestone = birthday.Milestone{Age: birthday.Period{Years: 12, Months: 6}}
	assert.Equal("12 years 6 months", milestone.AgeString())
}

func TestMilestoneLess(t *testing.T) {
	assert := asserts.New(t)
	lhs := birthday.Milestone{}
	rhs := birthday.Milestone{AgeUnknown: true}
	assert.True(lhs.Less(&rhs))
	assert.False(rhs.Less(&lhs))
}

func TestPeriodLess(t *testing.T) {
	assert := asserts.New(t)
	var lhs birthday.Period
	rhs := birthday.Period{Days: 17}
	assert.True(lhs.Less(rhs))
	assert.False(rhs.Less(lhs))
	rhs = birthday.Period{Weeks: 17}
	assert.True(lhs.Less(rhs))
	assert.False(rhs.Less(lhs))
}

func TestDiffPanics(t *testing.T) {
	assert := asserts.New(t)
	var p birthday.Period
	assert.Panics(func() {
		p.Diff(date_util.YMD(2020, 10, 15), date_util.YMD(2020, 10, 14))
	})
}

func TestPeriodMultiply(t *testing.T) {
	assert := asserts.New(t)
	p := birthday.Period{Days: 5, Normalize: true}
	assert.Equal(birthday.Period{Weeks: 4, Days: 2}, p.Multiply(6))
	p = birthday.Period{Days: 5}
	assert.Equal(birthday.Period{Days: 30}, p.Multiply(6))
}

func TestMilestonesBirthdayNextYear(t *testing.T) {
	assert := asserts.New(t)
	b := date_util.YMD(0, 1, 26)
	currentDate := date_util.YMD(2020, 10, 15)
	milestones := getMilestones(currentDate, b, 300)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:       date_util.YMD(2021, 1, 26),
				DaysAway:   103,
				AgeUnknown: true,
			},
		},
		milestones)
}

func TestMilestonesNoYear(t *testing.T) {
	assert := asserts.New(t)
	b := date_util.YMD(0, 9, 25)
	currentDate := date_util.YMD(2020, 9, 26)
	milestones := getMilestones(currentDate, b, 730)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:       date_util.YMD(2021, 9, 25),
				DaysAway:   364,
				AgeUnknown: true,
			},
			{
				Date:       date_util.YMD(2022, 9, 25),
				DaysAway:   729,
				AgeUnknown: true,
			},
		},
		milestones)
	milestones = getMilestones(currentDate, b, 729)
	assert.Len(milestones, 1)
	milestones = getMilestones(currentDate, b, 365)
	assert.Len(milestones, 1)
	milestones = getMilestones(currentDate, b, 364)
	assert.Empty(milestones)

	currentDate = date_util.YMD(2020, 9, 25)
	milestones = getMilestones(currentDate, b, 366)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:       date_util.YMD(2020, 9, 25),
				DaysAway:   0,
				AgeUnknown: true,
			},
			{
				Date:       date_util.YMD(2021, 9, 25),
				DaysAway:   365,
				AgeUnknown: true,
			},
		},
		milestones)
	milestones = getMilestones(currentDate, b, 365)
	assert.Len(milestones, 1)
	milestones = getMilestones(currentDate, b, 1)
	assert.Len(milestones, 1)
	milestones = getMilestones(currentDate, b, 0)
	assert.Empty(milestones)
	milestones = getMilestones(currentDate, b, -1000000)
	assert.Empty(milestones)
}

func TestMilestonesYearBefore(t *testing.T) {
	assert := asserts.New(t)
	b := date_util.YMD(1971, 9, 22)
	currentDate := date_util.YMD(2001, 9, 22)
	milestones := getMilestones(currentDate, b, 1043)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:     date_util.YMD(2001, 9, 22),
				DaysAway: 0,
				Age:      birthday.Period{Years: 30},
			},
			{
				Date:     date_util.YMD(2001, 11, 3),
				DaysAway: 42,
				Age:      birthday.Period{Days: 11000},
			},
			{
				Date:     date_util.YMD(2002, 9, 22),
				DaysAway: 365,
				Age:      birthday.Period{Years: 31},
			},
			{
				Date:     date_util.YMD(2003, 9, 22),
				DaysAway: 730,
				Age:      birthday.Period{Years: 32},
			},
			{
				Date:     date_util.YMD(2004, 7, 30),
				DaysAway: 1042,
				Age:      birthday.Period{Days: 12000},
			},
		},
		milestones)
	milestones = getMilestones(currentDate, b, 1042)
	assert.Len(milestones, 4)
	milestones = getMilestones(currentDate, b, 731)
	assert.Len(milestones, 4)
	milestones = getMilestones(currentDate, b, 730)
	assert.Len(milestones, 3)
	milestones = getMilestones(currentDate, b, 366)
	assert.Len(milestones, 3)
	milestones = getMilestones(currentDate, b, 365)
	assert.Len(milestones, 2)
	milestones = getMilestones(currentDate, b, 43)
	assert.Len(milestones, 2)
	milestones = getMilestones(currentDate, b, 42)
	assert.Len(milestones, 1)
	milestones = getMilestones(currentDate, b, 1)
	assert.Len(milestones, 1)
	milestones = getMilestones(currentDate, b, 0)
	assert.Empty(milestones)
	milestones = getMilestones(currentDate, b, -1000000)
	assert.Empty(milestones)

	currentDate = date_util.YMD(2001, 9, 23)
	milestones = getMilestones(currentDate, b, 1043)
	assert.Len(milestones, 4)
	currentDate = date_util.YMD(2001, 11, 3)
	milestones = getMilestones(currentDate, b, 1043)
	assert.Len(milestones, 4)
	currentDate = date_util.YMD(2001, 11, 4)
	milestones = getMilestones(currentDate, b, 1043)
	assert.Len(milestones, 3)
}

func TestMilestonesYearAfter(t *testing.T) {
	assert := asserts.New(t)
	b := date_util.YMD(2024, 2, 4)
	currentDate := date_util.YMD(2020, 10, 11)
	milestones := getMilestones(currentDate, b, 2212)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:     date_util.YMD(2024, 2, 4),
				DaysAway: 1211,
			},
			{
				Date:     date_util.YMD(2025, 2, 4),
				DaysAway: 1577,
				Age:      birthday.Period{Years: 1},
			},
			{
				Date:     date_util.YMD(2026, 2, 4),
				DaysAway: 1942,
				Age:      birthday.Period{Years: 2},
			},
			{
				Date:     date_util.YMD(2026, 10, 31),
				DaysAway: 2211,
				Age:      birthday.Period{Days: 1000},
			},
		},
		milestones)
}

func TestRemindPanic(t *testing.T) {
	assert := asserts.New(t)
	currentDate := date_util.YMD(2023, 1, 20)
	r := birthday.NewReminder(currentDate, 500)
	assert.Panics(func() { r.SetPeriods(birthday.Period{}) })
}

func TestRemindNoYears(t *testing.T) {
	assert := asserts.New(t)
	currentDate := date_util.YMD(2023, 1, 20)
	r := birthday.NewReminder(currentDate, 500)
	r.SetPeriods(thousandDays)
	e := birthday.Entry{
		Name:     "Mark",
		Birthday: date_util.YMD(2023, 1, 20),
	}
	r.Consume(&e)
	e = birthday.Entry{
		Name:     "Steve",
		Birthday: date_util.YMD(0, 2, 29),
	}
	r.Consume(&e)
	milestones := r.Milestones()
	assert.Equal(
		[]birthday.Milestone{
			{
				Name:     "Mark",
				Date:     date_util.YMD(2023, 1, 20),
				DaysAway: 0,
			},
		},
		milestones)
}

func TestRemindWithEverything(t *testing.T) {
	assert := asserts.New(t)
	currentDate := date_util.YMD(2017, 6, 11)
	r := birthday.NewReminder(currentDate, 1001)
	e := birthday.Entry{
		Name:     "Mark",
		Birthday: date_util.YMD(1968, 2, 29),
	}
	r.Consume(&e)
	milestones := r.Milestones()
	assert.Equal([]birthday.Milestone{
		{
			Name:     "Mark",
			Date:     date_util.YMD(2017, 6, 11),
			DaysAway: 0,
			Age:      birthday.Period{Days: 18000},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2017, 8, 29),
			DaysAway: 79,
			Age:      birthday.Period{Years: 49, Months: 6},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2017, 12, 28),
			DaysAway: 200,
			Age:      birthday.Period{Weeks: 2600},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2018, 3, 1),
			DaysAway: 263,
			Age:      birthday.Period{Years: 50},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2018, 3, 1),
			DaysAway: 263,
			Age:      birthday.Period{Months: 600},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2018, 8, 29),
			DaysAway: 444,
			Age:      birthday.Period{Years: 50, Months: 6},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2019, 3, 1),
			DaysAway: 628,
			Age:      birthday.Period{Years: 51},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2019, 8, 29),
			DaysAway: 809,
			Age:      birthday.Period{Years: 51, Months: 6},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2019, 11, 28),
			DaysAway: 900,
			Age:      birthday.Period{Weeks: 2700},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2020, 2, 29),
			DaysAway: 993,
			Age:      birthday.Period{Years: 52},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2020, 3, 7),
			DaysAway: 1000,
			Age:      birthday.Period{Days: 19000},
		},
	}, milestones)
}

func TestRemindWithWeeks(t *testing.T) {
	assert := asserts.New(t)
	currentDate := date_util.YMD(2017, 12, 28)
	r := birthday.NewReminder(currentDate, 701)
	r.SetPeriods(hundredWeeks)
	e := birthday.Entry{
		Name:     "Mark",
		Birthday: date_util.YMD(1968, 2, 29),
	}
	r.Consume(&e)
	milestones := r.Milestones()
	assert.Equal([]birthday.Milestone{
		{
			Name:     "Mark",
			Date:     date_util.YMD(2017, 12, 28),
			DaysAway: 0,
			Age:      birthday.Period{Weeks: 2600},
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2019, 11, 28),
			DaysAway: 700,
			Age:      birthday.Period{Weeks: 2700},
		},
	}, milestones)
}

func TestRemind(t *testing.T) {
	assert := asserts.New(t)
	currentDate := date_util.YMD(2023, 1, 20)
	r := birthday.NewReminder(currentDate, 500)
	r.SetPeriods(years, thousandDays)
	e := birthday.Entry{
		Name:     "Mark",
		Birthday: date_util.YMD(2023, 1, 20),
	}
	r.Consume(&e)
	e = birthday.Entry{
		Name:     "Steve",
		Birthday: date_util.YMD(0, 2, 29),
	}
	r.Consume(&e)
	milestones := r.Milestones()
	assert.Equal(
		[]birthday.Milestone{
			{
				Name:     "Mark",
				Date:     date_util.YMD(2023, 1, 20),
				DaysAway: 0,
			},
			{
				Name:       "Steve",
				Date:       date_util.YMD(2023, 3, 1),
				DaysAway:   40,
				AgeUnknown: true,
			},
			{
				Name:     "Mark",
				Date:     date_util.YMD(2024, 1, 20),
				DaysAway: 365,
				Age:      birthday.Period{Years: 1},
			},
			{
				Name:       "Steve",
				Date:       date_util.YMD(2024, 2, 29),
				DaysAway:   405,
				AgeUnknown: true,
			},
		},
		milestones)
}

func TestRemindHalfYear(t *testing.T) {
	assert := asserts.New(t)
	currentDate := date_util.YMD(2021, 3, 20)
	r := birthday.NewReminder(currentDate, 300)
	r.SetPeriods(years, sixMonths)
	e := birthday.Entry{
		Name:     "Mark",
		Birthday: date_util.YMD(1985, 3, 27),
	}
	r.Consume(&e)
	e = birthday.Entry{
		Name:     "Steve",
		Birthday: date_util.YMD(1984, 3, 27),
	}
	r.Consume(&e)
	milestones := r.Milestones()
	assert.Equal(
		[]birthday.Milestone{
			{
				Name:     "Mark",
				Date:     date_util.YMD(2021, 3, 27),
				DaysAway: 7,
				Age:      birthday.Period{Years: 36},
			},
			{
				Name:     "Steve",
				Date:     date_util.YMD(2021, 3, 27),
				DaysAway: 7,
				Age:      birthday.Period{Years: 37},
			},
			{
				Name:     "Mark",
				Date:     date_util.YMD(2021, 9, 27),
				DaysAway: 191,
				Age:      birthday.Period{Years: 36, Months: 6},
			},
			{
				Name:     "Steve",
				Date:     date_util.YMD(2021, 9, 27),
				DaysAway: 191,
				Age:      birthday.Period{Years: 37, Months: 6},
			},
		},
		milestones)
}

func TestRemindAgain(t *testing.T) {
	assert := asserts.New(t)
	currentDate := date_util.YMD(2023, 1, 20)
	r := birthday.NewReminder(currentDate, 406)
	r.SetPeriods(years, thousandDays)
	e := birthday.Entry{Name: "Matt", Birthday: date_util.YMD(1952, 2, 29)}
	r.Consume(&e)
	milestones := r.Milestones()
	assert.Equal(
		[]birthday.Milestone{
			{
				Name:     "Matt",
				Date:     date_util.YMD(2023, 3, 1),
				DaysAway: 40,
				Age:      birthday.Period{Years: 71},
			},
			{
				Name:     "Matt",
				Date:     date_util.YMD(2023, 5, 7),
				DaysAway: 107,
				Age:      birthday.Period{Days: 26000},
			},
			{
				Name:     "Matt",
				Date:     date_util.YMD(2024, 2, 29),
				DaysAway: 405,
				Age:      birthday.Period{Years: 72},
			},
		},
		milestones)
}

func TestFilterNone(t *testing.T) {
	assert := asserts.New(t)
	search := birthday.NewSearch("")
	e := birthday.Entry{
		Name:     "Bob",
		Birthday: date_util.YMD(0, 10, 15),
	}
	search.Consume(&e)
	e = birthday.Entry{
		Name:     "Billy",
		Birthday: date_util.YMD(1968, 11, 1),
	}
	search.Consume(&e)
	assert.Equal(
		[]birthday.Entry{
			{
				Name:     "Billy",
				Birthday: date_util.YMD(1968, 11, 1),
			},
			{
				Name:     "Bob",
				Birthday: date_util.YMD(0, 10, 15),
			},
		},
		search.Results())
}

func TestFilterSome(t *testing.T) {
	assert := asserts.New(t)
	search := birthday.NewSearch("jOHN  dOe")
	e := birthday.Entry{
		Name:     "John Doe",
		Birthday: date_util.YMD(2019, 10, 15),
	}
	search.Consume(&e)
	e = birthday.Entry{
		Name:     "Billy",
		Birthday: date_util.YMD(1968, 11, 1),
	}
	search.Consume(&e)
	assert.Equal(
		[]birthday.Entry{
			{
				Name:     "John Doe",
				Birthday: date_util.YMD(2019, 10, 15),
			},
		},
		search.Results())
}

func getMilestones(
	currentDate, b time.Time, daysAhead int) []birthday.Milestone {
	r := birthday.NewReminder(currentDate, daysAhead)
	r.SetPeriods(years, thousandDays)
	e := birthday.Entry{Birthday: b}
	r.Consume(&e)
	return r.Milestones()
}
