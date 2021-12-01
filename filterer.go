package birthday

import (
	"github.com/keep94/consume"
)

type entryFilterer func(ptr *Entry) bool

func EntryFilterer(f func(ptr *Entry) bool) consume.Filterer {
	return entryFilterer(f)
}

func (e entryFilterer) Filter(ptr interface{}) bool {
	return e(ptr.(*Entry))
}

type milestoneFilterer func(ptr *Milestone) bool

func MilestoneFilterer(f func(ptr *Milestone) bool) consume.Filterer {
	return milestoneFilterer(f)
}

func (m milestoneFilterer) Filter(ptr interface{}) bool {
	return m(ptr.(*Milestone))
}
