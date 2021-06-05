package common

import (
	"html/template"
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
func ParseDate(dateStr string) time.Time {
	result, err := birthday.Parse(dateStr)
	if err != nil {
		return birthday.Today()
	}
	return fixMissingYear(result)
}

func fixMissingYear(date time.Time) time.Time {
	if birthday.HasYear(date) {
		return date
	}
	today := birthday.Today()
	return date_util.YMD(today.Year(), int(date.Month()), date.Day())
}
