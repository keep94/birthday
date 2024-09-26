package common

import (
	"html/template"
	"strings"
	"time"

	"github.com/keep94/birthday"
	"github.com/keep94/toolbox/date_util"
)

// NewTemplate returns a new template instance. name is the name
// of the template; templateStr is the template string.
func NewTemplate(name, templateStr string) *template.Template {
	return template.Must(template.New(name).Parse(templateStr))
}

// ParseDate parses dateStr to a time in UTC.
// dateStr can be of form mm/dd or mm/dd/yyyy. If dateStr is of form
// mm/dd, then the current year is used as the year. If there is an error
// parsing dateStr, then the ParseDate() returns the current date.
func ParseDate(clock date_util.Clock, dateStr string) time.Time {
	today := birthday.Today(clock)
	result, err := birthday.Parse(dateStr)
	if err != nil {
		return today
	}
	return fixMissingYear(today, result)
}

func fixMissingYear(today, date time.Time) time.Time {
	if birthday.HasYear(date) {
		return date
	}
	return date_util.YMD(today.Year(), int(date.Month()), date.Day())
}

// ParsePeriods parses a periodStr of form 'ymwdh' into a slice of periods.
// y stands for year; m stands for 100 months; w stands for 100 weeks;
// d stands for 1000 days; h stands for half-year. If periodStr is empty,
// ParsePeriods returns a copy of defaultPeriods. periodStr can be any
// subset of 'ymwdh'
func ParsePeriods(
	periodStr string, defaultPeriods []birthday.Period) []birthday.Period {
	var result []birthday.Period
	if periodStr == "" {
		result = make([]birthday.Period, len(defaultPeriods))
		copy(result, defaultPeriods)
		return result
	}
	if strings.Contains(periodStr, "y") {
		result = append(result, birthday.Period{Years: 1})
	}
	if strings.Contains(periodStr, "m") {
		result = append(result, birthday.Period{Months: 100})
	}
	if strings.Contains(periodStr, "w") {
		result = append(result, birthday.Period{Weeks: 100})
	}
	if strings.Contains(periodStr, "d") {
		result = append(result, birthday.Period{Days: 1000})
	}
	if strings.Contains(periodStr, "h") {
		result = append(result, birthday.Period{Months: 6, Normalize: true})
	}
	return result
}
