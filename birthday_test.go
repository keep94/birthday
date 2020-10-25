package birthday_test

import (
	"testing"
	"time"

	"github.com/keep94/birthday"
	"github.com/keep94/toolbox/date_util"
	asserts "github.com/stretchr/testify/assert"
)

func TestAsDays(t *testing.T) {
	assert := asserts.New(t)
	start := date_util.YMD(2020, 2, 29)
	end := date_util.YMD(2020, 3, 1)
	assert.Equal(1, birthday.AsDays(end)-birthday.AsDays(start))
	start = date_util.YMD(2001, 9, 17)
	end = date_util.YMD(2018, 8, 5)
	assert.Equal(6166, birthday.AsDays(end)-birthday.AsDays(start))
	b := date_util.YMD(1970, 1, 1)
	assert.Equal(0, birthday.AsDays(b))
	b = date_util.YMD(1930, 1, 1)
	assert.Equal(-14610, birthday.AsDays(b))
}

func TestNormalize(t *testing.T) {
	assert := asserts.New(t)
	b := date_util.YMD(2021, 2, 29)
	expected := date_util.YMD(2021, 3, 1)
	assert.Equal(expected, birthday.FromDays(birthday.AsDays(b)))
	b = date_util.YMD(1900, 2, 28)
	expected = date_util.YMD(1900, 2, 28)
	assert.Equal(expected, birthday.FromDays(birthday.AsDays(b)))
}

func TestFromDays(t *testing.T) {
	assert := asserts.New(t)
	assert.Equal(date_util.YMD(1950, 1, 1), birthday.FromDays(-7305))
}

func TestDiffInYears(t *testing.T) {
	assert := asserts.New(t)
	end := date_util.YMD(1951, 2, 15)
	assert.Equal(0, birthday.DiffInYears(end, date_util.YMD(1951, 2, 15)))
	assert.Equal(-1, birthday.DiffInYears(end, date_util.YMD(1951, 2, 16)))
	assert.Equal(-1, birthday.DiffInYears(end, date_util.YMD(1951, 3, 1)))
	assert.Equal(0, birthday.DiffInYears(end, date_util.YMD(1951, 2, 14)))
	assert.Equal(0, birthday.DiffInYears(end, date_util.YMD(1951, 1, 31)))
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
				Date:      date_util.YMD(2001, 11, 3),
				DaysAway:  42,
				Age:       11000,
				AgeInDays: true,
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
				Date:      date_util.YMD(2004, 7, 30),
				DaysAway:  1042,
				Age:       12000,
				AgeInDays: true,
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
				Date:      date_util.YMD(2024, 2, 4),
				DaysAway:  1211,
				Age:       0,
				AgeInDays: true,
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
				Date:      date_util.YMD(2026, 10, 31),
				DaysAway:  2211,
				Age:       1000,
				AgeInDays: true,
			},
		},
		milestones)
}

func TestRemind(t *testing.T) {
	assert := asserts.New(t)
	currentDate := date_util.YMD(2023, 1, 20)
	r := birthday.NewRemind(currentDate, 500)
	r.Add("Mark", date_util.YMD(2023, 1, 20))
	r.Add("Steve", date_util.YMD(0, 2, 29))
	milestones := r.Reminders()
	assert.Equal(
		[]birthday.Milestone{
			{
				Name:     "Mark",
				Date:     date_util.YMD(2023, 1, 20),
				DaysAway: 0,
				Age:      0,
			},
			{
				Name:      "Mark",
				Date:      date_util.YMD(2023, 1, 20),
				DaysAway:  0,
				Age:       0,
				AgeInDays: true,
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
	r := birthday.NewRemind(currentDate, 406)
	r.Add("Matt", date_util.YMD(1952, 2, 29))
	milestones := r.Reminders()
	assert.Equal(
		[]birthday.Milestone{
			{
				Name:     "Matt",
				Date:     date_util.YMD(2023, 3, 1),
				DaysAway: 40,
				Age:      71,
			},
			{
				Name:      "Matt",
				Date:      date_util.YMD(2023, 5, 7),
				DaysAway:  107,
				Age:       26000,
				AgeInDays: true,
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
	today := date_util.YMD(2020, 10, 15)
	filter := birthday.NewFilter(today, "")
	filter.Add("Bob", date_util.YMD(0, 10, 15))
	filter.Add("Billy", date_util.YMD(1968, 11, 1))
	assert.Equal(
		[]birthday.Person{
			{
				Name:       "Billy",
				Birthday:   date_util.YMD(1968, 11, 1),
				AgeInYears: 51,
				AgeInDays:  18976,
			},
			{
				Name:     "Bob",
				Birthday: date_util.YMD(0, 10, 15),
			},
		},
		filter.Persons())
}

func TestFilterSome(t *testing.T) {
	assert := asserts.New(t)
	today := date_util.YMD(2020, 10, 15)
	filter := birthday.NewFilter(today, "jOHN  dOe")
	filter.Add("John Doe", date_util.YMD(2019, 10, 15))
	filter.Add("Billy", date_util.YMD(1968, 11, 1))
	assert.Equal(
		[]birthday.Person{
			{
				Name:       "John Doe",
				Birthday:   date_util.YMD(2019, 10, 15),
				AgeInYears: 1,
				AgeInDays:  366,
			},
		},
		filter.Persons())
}

func getMilestones(
	currentDate, b time.Time, daysAhead int) []birthday.Milestone {
	r := birthday.NewRemind(currentDate, daysAhead)
	r.Add("", b)
	return r.Reminders()
}
