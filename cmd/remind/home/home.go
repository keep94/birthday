package home

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/keep94/birthday"
	"github.com/keep94/birthday/cmd/remind/common"
	"github.com/keep94/consume2"
	"github.com/keep94/toolbox/date_util"
	"github.com/keep94/toolbox/http_util"
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
      <td {{if $top.Today .}}class="today"{{end}}>{{.EntryPtr.Name}}</td>
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
	Store          birthday.Store
	DaysAhead      int
	FirstN         consume2.Pipeline[birthday.Milestone, birthday.Milestone]
	DefaultPeriods []birthday.Period
	Clock          date_util.Clock
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var entries []*birthday.Entry
	err := h.Store.Read(
		consume2.Filter(
			consume2.AppendPtrsTo(&entries),
			birthday.Query(r.Form.Get("q"))))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	daysAhead := h.parseDays(r.Form.Get("days"))
	pipeline := consume2.PTakeWhile(func(m birthday.Milestone) bool {
		return m.DaysAway < daysAhead
	})
	pipeline = consume2.Join(pipeline, h.FirstN)
	var milestones []*birthday.Milestone
	birthday.Remind(
		entries,
		common.ParsePeriods(r.Form.Get("p"), h.DefaultPeriods),
		common.ParseDate(h.Clock, r.Form.Get("date")),
		pipeline.Run(consume2.AppendPtrsTo(&milestones)),
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
	Milestones []*birthday.Milestone
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
