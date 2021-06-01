package home

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/keep94/birthday"
	"github.com/keep94/birthday/cmd/remind/common"
	"github.com/keep94/toolbox/date_util"
	"github.com/keep94/toolbox/http_util"
)

const (
	kMaxDaysAhead = 365
)

var (
	kTemplateSpec = `
<html>
<head>
  <title>Birthday Reminders</title>
  <style>
  h1 {
    font-size: 40px;
  }
  th {
    font-size: 30px;
  }
  td {
    font-size: 30px;
  }
  td.today {
    font-style: italic;
  }
  </style>
</head>
<body>
  <h1>Birthday Reminders</h1>
  <table border=1>
    <tr>
      <th>Date</th>
      <th>Name</th>
      <th>Age</th>
    </tr>
    {{with $top := .}}
    {{range .Milestones}}
    <tr>
      <td {{if $top.Today .}}class="today"{{end}}>{{$top.DateStr .}}</td>
      <td {{if $top.Today .}}class="today"{{end}}>{{.Name}}</td>
      <td {{if $top.Today .}}class="today"{{end}}>{{.AgeString}}</td>
    </tr>
    {{end}}
    {{end}}
  </table>
</body>
</html>`
)

var (
	kTemplate *template.Template
)

type Handler struct {
	File      string
	DaysAhead int
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	reminder := birthday.NewReminder(
		getDate(r.Form.Get("date")), h.getDays(r.Form.Get("days")))
	err := birthday.ReadFile(h.File, reminder)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	milestones := reminder.Milestones()
	http_util.WriteTemplate(w, kTemplate, &view{Milestones: milestones})
}

func (h *Handler) getDays(daysStr string) int {
	result, err := strconv.Atoi(daysStr)
	if err != nil {
		return h.DaysAhead
	}
	if result > kMaxDaysAhead {
		result = kMaxDaysAhead
	}
	return result
}

func fixMissingYear(date time.Time) time.Time {
	if birthday.HasYear(date) {
		return date
	}
	today := birthday.Today()
	return date_util.YMD(today.Year(), int(date.Month()), date.Day())
}

func getDate(dateStr string) time.Time {
	result, err := birthday.Parse(dateStr)
	if err != nil {
		return birthday.Today()
	}
	return fixMissingYear(result)
}

type view struct {
	Milestones []birthday.Milestone
}

func (b *view) DateStr(milestone *birthday.Milestone) string {
	return birthday.ToStringWithWeekDay(milestone.Date)
}

func (v *view) Today(milestone *birthday.Milestone) bool {
	return milestone.DaysAway == 0
}

func init() {
	kTemplate = common.NewTemplate("home", kTemplateSpec)
}
