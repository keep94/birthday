package birthday_test

import (
	"testing"

	"github.com/keep94/birthday"
	asserts "github.com/stretchr/testify/assert"
)

func TestAsDays(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Year: 2020, Month: 10, Day: 10}
	assert.Equal(6, b.AsDays()%7)
	start := birthday.Birthday{Year: 2020, Month: 2, Day: 29}
	end := birthday.Birthday{Year: 2020, Month: 3, Day: 1}
	assert.Equal(1, end.AsDays()-start.AsDays())
	start = birthday.Birthday{Year: 2001, Month: 9, Day: 17}
	end = birthday.Birthday{Year: 2018, Month: 8, Day: 5}
	assert.Equal(6166, end.AsDays()-start.AsDays())
	b = birthday.Birthday{Year: 1, Month: 1, Day: 1}
	assert.Equal(1, b.AsDays())
	b = birthday.Birthday{Month: 5, Day: 31}
	assert.Panics(func() { b.AsDays() })
}

func TestIsValid(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Year: 2000, Month: 2, Day: 29}
	assert.True(b.IsValid())
	b = birthday.Birthday{Year: 2100, Month: 2, Day: 29}
	assert.False(b.IsValid())
	b = birthday.Birthday{Year: 2008, Month: 2, Day: 29}
	assert.True(b.IsValid())
	b = birthday.Birthday{Year: 2007, Month: 2, Day: 29}
	assert.False(b.IsValid())
	b = birthday.Birthday{Year: 2008, Month: 4, Day: 31}
	assert.False(b.IsValid())
	b = birthday.Birthday{Year: 2008, Month: 4, Day: 30}
	assert.True(b.IsValid())
	b = birthday.Birthday{Year: 2008, Month: 4, Day: 1}
	assert.True(b.IsValid())
	b = birthday.Birthday{Year: 2008, Month: 4}
	assert.False(b.IsValid())
	b = birthday.Birthday{Year: 2008}
	assert.False(b.IsValid())
	b = birthday.Birthday{Year: 2008, Month: 13, Day: 1}
	assert.False(b.IsValid())
	b = birthday.Birthday{Year: 2008, Month: 12, Day: 31}
	assert.True(b.IsValid())
	b = birthday.Birthday{Year: 2008, Month: 1, Day: 1}
	assert.True(b.IsValid())
	b = birthday.Birthday{Month: 2, Day: 29}
	assert.True(b.IsValid())
	b = birthday.Birthday{Month: 2, Day: 1}
	assert.True(b.IsValid())
	b = birthday.Birthday{Month: 2}
	assert.False(b.IsValid())
	b = birthday.Birthday{Month: 2, Day: 30}
	assert.False(b.IsValid())
	b = birthday.Birthday{Month: 6, Day: 30}
	assert.True(b.IsValid())
}

func TestNormalize(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Year: 2021, Month: 2, Day: 29}
	expected := birthday.Birthday{Year: 2021, Month: 3, Day: 1}
	assert.Equal(expected, b.Normalize())
	b = birthday.Birthday{Year: 2021, Month: 12, Day: 3}
	expected = birthday.Birthday{Year: 2021, Month: 12, Day: 3}
	assert.Equal(expected, b.Normalize())
	b = birthday.Birthday{Year: 2021, Month: 1, Day: 3}
	expected = birthday.Birthday{Year: 2021, Month: 1, Day: 3}
	assert.Equal(expected, b.Normalize())
	b = birthday.Birthday{Year: 2021, Month: 25, Day: 29}
	expected = birthday.Birthday{Year: 2023, Month: 1, Day: 29}
	assert.Equal(expected, b.Normalize())
	b = birthday.Birthday{Year: 2021, Month: -100, Day: -7}
	expected = birthday.Birthday{Year: 2012, Month: 7, Day: 24}
	assert.Equal(expected, b.Normalize())
	b = birthday.Birthday{Year: 2024, Month: 2, Day: 29}
	expected = birthday.Birthday{Year: 2024, Month: 2, Day: 29}
	assert.Equal(expected, b.Normalize())
	b = birthday.Birthday{Year: 2000, Month: 2, Day: 29}
	expected = birthday.Birthday{Year: 2000, Month: 2, Day: 29}
	assert.Equal(expected, b.Normalize())
	b = birthday.Birthday{Year: 1900, Month: 2, Day: 28}
	expected = birthday.Birthday{Year: 1900, Month: 2, Day: 28}
	assert.Equal(expected, b.Normalize())
	b = birthday.Birthday{Year: 1900, Month: 3, Day: 1}
	expected = birthday.Birthday{Year: 1900, Month: 3, Day: 1}
	assert.Equal(expected, b.Normalize())
}

func TestFromDays(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Year: 1991, Month: 1, Day: 15}
	assert.Equal(b, birthday.FromDays(b.AsDays()))
	b = birthday.Birthday{Year: 1977, Month: 7, Day: 14}
	assert.Equal(b, birthday.FromDays(b.AsDays()))
	b = birthday.Birthday{Year: 1, Month: 1, Day: 1}
	assert.Equal(b, birthday.FromDays(1))
	assert.Panics(func() { birthday.FromDays(0) })

	today := birthday.Birthday{Year: 2020, Month: 10, Day: 11}
	future := birthday.Birthday{Year: 2020, Month: 12, Day: 31}
	assert.Equal(future, birthday.FromDays(today.AsDays()+81))
	future = birthday.Birthday{Year: 2021, Month: 1, Day: 1}
	assert.Equal(future, birthday.FromDays(today.AsDays()+82))
}

func TestMilestonesPanic(t *testing.T) {
	assert := asserts.New(t)
	currentDate := birthday.Birthday{Month: 7, Day: 12}
	b := birthday.Birthday{Month: 8, Day: 4}
	assert.Panics(func() { getMilestones(currentDate, b, 200) })
	currentDate = birthday.Birthday{Year: 2020, Month: 1, Day: 12}
	b = birthday.Birthday{Month: 8, Day: 32}
	assert.Panics(func() { getMilestones(currentDate, b, 200) })
}

func TestMilestonesBirthdayNextYear(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Month: 1, Day: 26}
	currentDate := birthday.Birthday{Year: 2020, Month: 10, Day: 15}
	milestones := getMilestones(currentDate, b, 300)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:     birthday.Birthday{Year: 2021, Month: 1, Day: 26},
				DaysAway: 103,
				Age:      -1,
			},
		},
		milestones)
}

func TestMilestonesNoYear(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Month: 9, Day: 25}
	currentDate := birthday.Birthday{Year: 2020, Month: 9, Day: 26}
	milestones := getMilestones(currentDate, b, 730)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:     birthday.Birthday{Year: 2021, Month: 9, Day: 25},
				DaysAway: 364,
				Age:      -1,
			},
			{
				Date:     birthday.Birthday{Year: 2022, Month: 9, Day: 25},
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

	currentDate = birthday.Birthday{Year: 2020, Month: 9, Day: 25}
	milestones = getMilestones(currentDate, b, 366)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:     birthday.Birthday{Year: 2020, Month: 9, Day: 25},
				DaysAway: 0,
				Age:      -1,
			},
			{
				Date:     birthday.Birthday{Year: 2021, Month: 9, Day: 25},
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
	b := birthday.Birthday{Year: 1971, Month: 9, Day: 22}
	currentDate := birthday.Birthday{Year: 2001, Month: 9, Day: 22}
	milestones := getMilestones(currentDate, b, 1043)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:     birthday.Birthday{Year: 2001, Month: 9, Day: 22},
				DaysAway: 0,
				Age:      30,
			},
			{
				Date:      birthday.Birthday{Year: 2001, Month: 11, Day: 3},
				DaysAway:  42,
				Age:       11000,
				AgeInDays: true,
			},
			{
				Date:     birthday.Birthday{Year: 2002, Month: 9, Day: 22},
				DaysAway: 365,
				Age:      31,
			},
			{
				Date:     birthday.Birthday{Year: 2003, Month: 9, Day: 22},
				DaysAway: 730,
				Age:      32,
			},
			{
				Date:      birthday.Birthday{Year: 2004, Month: 7, Day: 30},
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

	currentDate = birthday.Birthday{Year: 2001, Month: 9, Day: 23}
	milestones = getMilestones(currentDate, b, 1043)
	assert.Len(milestones, 4)
	currentDate = birthday.Birthday{Year: 2001, Month: 11, Day: 3}
	milestones = getMilestones(currentDate, b, 1043)
	assert.Len(milestones, 4)
	currentDate = birthday.Birthday{Year: 2001, Month: 11, Day: 4}
	milestones = getMilestones(currentDate, b, 1043)
	assert.Len(milestones, 3)
}

func TestMilestonesYearAfter(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Year: 2024, Month: 2, Day: 4}
	currentDate := birthday.Birthday{Year: 2020, Month: 10, Day: 11}
	milestones := getMilestones(currentDate, b, 2212)
	assert.Equal(
		[]birthday.Milestone{
			{
				Date:     birthday.Birthday{Year: 2024, Month: 2, Day: 4},
				DaysAway: 1211,
				Age:      0,
			},
			{
				Date:      birthday.Birthday{Year: 2024, Month: 2, Day: 4},
				DaysAway:  1211,
				Age:       0,
				AgeInDays: true,
			},
			{
				Date:     birthday.Birthday{Year: 2025, Month: 2, Day: 4},
				DaysAway: 1577,
				Age:      1,
			},
			{
				Date:     birthday.Birthday{Year: 2026, Month: 2, Day: 4},
				DaysAway: 1942,
				Age:      2,
			},
			{
				Date:      birthday.Birthday{Year: 2026, Month: 10, Day: 31},
				DaysAway:  2211,
				Age:       1000,
				AgeInDays: true,
			},
		},
		milestones)
}

func TestRemindPanic(t *testing.T) {
	assert := asserts.New(t)
	currentDate := birthday.Birthday{Month: 5, Day: 30}
	assert.Panics(func() { birthday.NewRemind(currentDate, 21) })
}

func TestRemind(t *testing.T) {
	assert := asserts.New(t)
	currentDate := birthday.Birthday{Year: 2023, Month: 1, Day: 20}
	r := birthday.NewRemind(currentDate, 500)
	r.Add("Mark", birthday.Birthday{Year: 2023, Month: 1, Day: 20})
	r.Add("Steve", birthday.Birthday{Month: 2, Day: 29})
	milestones := r.Reminders()
	assert.Equal(
		[]birthday.Milestone{
			{
				Name:     "Mark",
				Date:     birthday.Birthday{Year: 2023, Month: 1, Day: 20},
				DaysAway: 0,
				Age:      0,
			},
			{
				Name:      "Mark",
				Date:      birthday.Birthday{Year: 2023, Month: 1, Day: 20},
				DaysAway:  0,
				Age:       0,
				AgeInDays: true,
			},
			{
				Name:     "Steve",
				Date:     birthday.Birthday{Year: 2023, Month: 3, Day: 1},
				DaysAway: 40,
				Age:      -1,
			},
			{
				Name:     "Mark",
				Date:     birthday.Birthday{Year: 2024, Month: 1, Day: 20},
				DaysAway: 365,
				Age:      1,
			},
			{
				Name:     "Steve",
				Date:     birthday.Birthday{Year: 2024, Month: 2, Day: 29},
				DaysAway: 405,
				Age:      -1,
			},
		},
		milestones)
}

func TestRemindAgain(t *testing.T) {
	assert := asserts.New(t)
	currentDate := birthday.Birthday{Year: 2023, Month: 1, Day: 20}
	r := birthday.NewRemind(currentDate, 406)
	r.Add("Matt", birthday.Birthday{Year: 1972, Month: 2, Day: 29})
	milestones := r.Reminders()
	assert.Equal(
		[]birthday.Milestone{
			{
				Name:     "Matt",
				Date:     birthday.Birthday{Year: 2023, Month: 3, Day: 1},
				DaysAway: 40,
				Age:      51,
			},
			{
				Name:     "Matt",
				Date:     birthday.Birthday{Year: 2024, Month: 2, Day: 29},
				DaysAway: 405,
				Age:      52,
			},
		},
		milestones)
}

func TestString(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Month: 7, Day: 4, Year: 1992}
	assert.Equal("07/04/1992", b.String())
	b = birthday.Birthday{Month: 11, Day: 30, Year: 953}
	assert.Equal("11/30/0953", b.String())
	b = birthday.Birthday{Month: 2, Day: 28}
	assert.Equal("02/28", b.String())
	b = birthday.Birthday{Month: 11, Day: 4}
	assert.Equal("11/04", b.String())
}

func TestStringWithWeekday(t *testing.T) {
	assert := asserts.New(t)
	b := birthday.Birthday{Month: 10, Day: 15, Year: 2020}
	assert.Equal("Thu 10/15/2020", b.StringWithWeekDay())
}

func TestParse(t *testing.T) {
	assert := asserts.New(t)
	b, err := birthday.Parse("11/04")
	assert.NoError(err)
	assert.Equal(birthday.Birthday{Month: 11, Day: 4}, b)
	b, err = birthday.Parse("3/28/2017")
	assert.NoError(err)
	assert.Equal(birthday.Birthday{Month: 3, Day: 28, Year: 2017}, b)
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
}

func getMilestones(
	currentDate, b birthday.Birthday, daysAhead int) []birthday.Milestone {
	r := birthday.NewRemind(currentDate, daysAhead)
	r.Add("", b)
	return r.Reminders()
}
