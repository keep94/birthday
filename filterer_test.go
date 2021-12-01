package birthday_test

import (
	"testing"

	"github.com/keep94/birthday"
	"github.com/stretchr/testify/assert"
)

func TestEntryFilterer(t *testing.T) {
	assert := assert.New(t)
	var entry birthday.Entry
	filtererTrue := birthday.EntryFilterer(
		func(ptr *birthday.Entry) bool {
			return true
		})
	filtererFalse := birthday.EntryFilterer(
		func(ptr *birthday.Entry) bool {
			return false
		})
	assert.True(filtererTrue.Filter(&entry))
	assert.False(filtererFalse.Filter(&entry))
}

func TestMilestoneFilterer(t *testing.T) {
	assert := assert.New(t)
	var milestone birthday.Milestone
	filtererTrue := birthday.MilestoneFilterer(
		func(ptr *birthday.Milestone) bool {
			return true
		})
	filtererFalse := birthday.MilestoneFilterer(
		func(ptr *birthday.Milestone) bool {
			return false
		})
	assert.True(filtererTrue.Filter(&milestone))
	assert.False(filtererFalse.Filter(&milestone))
}
