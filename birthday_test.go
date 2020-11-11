package birthday_test

import (
	"testing"
	"time"

	"github.com/keep94/birthday"
	"github.com/keep94/toolbox/date_util"
	asserts "github.com/stretchr/testify/assert"
)

var years = birthday.Period{Count: 1, Unit: birthday.Years}
var hundredMonths = birthday.Period{Count: 100, Unit: birthday.Months}
var hundredWeeks = birthday.Period{Count: 100, Unit: birthday.Weeks}
var thousandDays = birthday.Period{Count: 1000, Unit: birthday.Days}

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
	assert.Equal(-6166, birthday.Days.Diff(start, end))
	assert.Equal(-881, birthday.Weeks.Diff(start, end))
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
	assert.Equal(-139, birthday.Months.Diff(end, start))
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
	assert.Equal(3, birthday.Years.Diff(end, date_util.YMD(1948, 2, 15)))
}

func TestAdd(t *testing.T) {
	assert := asserts.New(t)
	start := date_util.YMD(1952, 2, 29)
	assert.Equal(date_util.YMD(1952, 2, 24), birthday.Days.Add(start, -5))
	assert.Equal(date_util.YMD(1952, 3, 21), birthday.Weeks.Add(start, 3))
	assert.Equal(date_util.YMD(1952, 6, 29), birthday.Months.Add(start, 4))
	assert.Equal(date_util.YMD(1950, 3, 1), birthday.Years.Add(start, -2))
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
	milestone := birthday.Milestone{Age: -1, Unit: birthday.Years}
	assert.Equal("? years", milestone.AgeString())
	milestone = birthday.Milestone{Age: 47, Unit: birthday.Years}
	assert.Equal("47 years", milestone.AgeString())
	milestone = birthday.Milestone{Age: 5, Unit: birthday.Months}
	assert.Equal("5 months", milestone.AgeString())
	milestone = birthday.Milestone{Age: 3, Unit: birthday.Weeks}
	assert.Equal("3 weeks", milestone.AgeString())
	milestone = birthday.Milestone{Age: 0, Unit: birthday.Days}
	assert.Equal("0 days", milestone.AgeString())
	milestone = birthday.Milestone{Age: -2, Unit: birthday.Days}
	assert.Equal("? days", milestone.AgeString())
}

func TestMilestonesBirthdayNextYear(t *testing.T) {
	assert := asserts.New(t)
	b := date_util.YMD(0, 1, 26)
	currentDate := date_util.YMD(2020, 10, 15)
	milestones := getMilestones(currentDate, b, 300)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:     date_util.YMD(2021, 1, 26),
				DaysAway: 103,
				Age:      -1,
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
				Date:     date_util.YMD(2021, 9, 25),
				DaysAway: 364,
				Age:      -1,
			},
			{
				Date:     date_util.YMD(2022, 9, 25),
				DaysAway: 729,
				Age:      -1,
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
				Date:     date_util.YMD(2020, 9, 25),
				DaysAway: 0,
				Age:      -1,
			},
			{
				Date:     date_util.YMD(2021, 9, 25),
				DaysAway: 365,
				Age:      -1,
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
				Age:      30,
			},
			{
				Date:     date_util.YMD(2001, 11, 3),
				DaysAway: 42,
				Age:      11000,
				Unit:     birthday.Days,
			},
			{
				Date:     date_util.YMD(2002, 9, 22),
				DaysAway: 365,
				Age:      31,
			},
			{
				Date:     date_util.YMD(2003, 9, 22),
				DaysAway: 730,
				Age:      32,
			},
			{
				Date:     date_util.YMD(2004, 7, 30),
				DaysAway: 1042,
				Age:      12000,
				Unit:     birthday.Days,
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
				Age:      0,
			},
			{
				Date:     date_util.YMD(2024, 2, 4),
				DaysAway: 1211,
				Age:      0,
				Unit:     birthday.Days,
			},
			{
				Date:     date_util.YMD(2025, 2, 4),
				DaysAway: 1577,
				Age:      1,
			},
			{
				Date:     date_util.YMD(2026, 2, 4),
				DaysAway: 1942,
				Age:      2,
			},
			{
				Date:     date_util.YMD(2026, 10, 31),
				DaysAway: 2211,
				Age:      1000,
				Unit:     birthday.Days,
			},
		},
		milestones)
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
				Age:      0,
				Unit:     birthday.Days,
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
			Age:      18000,
			Unit:     birthday.Days,
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2017, 12, 28),
			DaysAway: 200,
			Age:      2600,
			Unit:     birthday.Weeks,
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2018, 3, 1),
			DaysAway: 263,
			Age:      50,
			Unit:     birthday.Years,
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2018, 3, 1),
			DaysAway: 263,
			Age:      600,
			Unit:     birthday.Months,
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2019, 3, 1),
			DaysAway: 628,
			Age:      51,
			Unit:     birthday.Years,
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2019, 11, 28),
			DaysAway: 900,
			Age:      2700,
			Unit:     birthday.Weeks,
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2020, 2, 29),
			DaysAway: 993,
			Age:      52,
			Unit:     birthday.Years,
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2020, 3, 7),
			DaysAway: 1000,
			Age:      19000,
			Unit:     birthday.Days,
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
			Age:      2600,
			Unit:     birthday.Weeks,
		},
		{
			Name:     "Mark",
			Date:     date_util.YMD(2019, 11, 28),
			DaysAway: 700,
			Age:      2700,
			Unit:     birthday.Weeks,
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
				Age:      0,
			},
			{
				Name:     "Mark",
				Date:     date_util.YMD(2023, 1, 20),
				DaysAway: 0,
				Age:      0,
				Unit:     birthday.Days,
			},
			{
				Name:     "Steve",
				Date:     date_util.YMD(2023, 3, 1),
				DaysAway: 40,
				Age:      -1,
			},
			{
				Name:     "Mark",
				Date:     date_util.YMD(2024, 1, 20),
				DaysAway: 365,
				Age:      1,
			},
			{
				Name:     "Steve",
				Date:     date_util.YMD(2024, 2, 29),
				DaysAway: 405,
				Age:      -1,
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
				Age:      71,
			},
			{
				Name:     "Matt",
				Date:     date_util.YMD(2023, 5, 7),
				DaysAway: 107,
				Age:      26000,
				Unit:     birthday.Days,
			},
			{
				Name:     "Matt",
				Date:     date_util.YMD(2024, 2, 29),
				DaysAway: 405,
				Age:      72,
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
