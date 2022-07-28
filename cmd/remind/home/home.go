package home

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/keep94/birthday"
	"github.com/keep94/birthday/cmd/remind/common"
	"github.com/keep94/consume2"
	"github.com/keep94/toolbox/http_util"
)

const (
	kMaxMilestones = 100
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
	var entries []birthday.Entry
	err := birthday.ReadFile(
		h.File,
		consume2.Filter(
			consume2.AppendTo(&entries),
			birthday.Query(r.Form.Get("q"))))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	daysAhead := h.parseDays(r.Form.Get("days"))
	var milestones []birthday.Milestone
	birthday.Remind(
		entries,
		common.ParsePeriods(r.Form.Get("p")),
		common.ParseDate(r.Form.Get("date")),
		consume2.TakeWhile(
			consume2.Slice(
				consume2.AppendTo(&milestones), 0, kMaxMilestones),
			func(m birthday.Milestone) bool {
				return m.DaysAway < daysAhead
			},
		),
	)
	http_util.WriteTemplate(w, kTemplate, &view{Milestones: milestones})
}

func (h *Handler) parseDays(daysStr string) int {
	result, err := strconv.Atoi(daysStr)
	if err != nil {
		return h.DaysAhead
	}
	return result
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
