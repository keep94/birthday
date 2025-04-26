package home

import (
	"fmt"
	"html/template"
	"iter"
	"net/http"
	"strconv"
	"time"

	"github.com/keep94/birthday"
	"github.com/keep94/birthday/cmd/remind/common"
	"github.com/keep94/consume2"
	"github.com/keep94/itertools"
	"github.com/keep94/toolbox/date_util"
	"github.com/keep94/toolbox/http_util"
)

var (
	kTemplateSpec = `
<html>
<head>
  <title>Birthdays</title>
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
  {{with .BuildId}}
      <h1>Birthdays {{.}}</h1>
  {{else}}
      <h1>Birthdays</h1>
  {{end}}
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
	MaxRows        int
	BuildId        string
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
	today := common.ParseDate(h.Clock, r.Form.Get("date"))
	endDate := today.AddDate(0, 0, daysAhead)
	seq := birthday.RemindPtrs(
		entries,
		common.ParsePeriods(r.Form.Get("p"), h.DefaultPeriods),
		today)
	seq = itertools.TakeWhile(
		func(m *birthday.Milestone) bool { return m.Date.Before(endDate) },
		seq)
	seq = itertools.Take(h.MaxRows, seq)
	http_util.WriteTemplate(
		w, kTemplate, &view{Milestones: seq, BuildId: h.BuildId, today: today})
}

func (h *Handler) parseDays(daysStr string) int {
	result, err := strconv.Atoi(daysStr)
	if err != nil {
		return h.DaysAhead
	}
	return result
}

type view struct {
	Milestones iter.Seq[*birthday.Milestone]
	BuildId    string
	today      time.Time
}

func (b *view) DateStr(milestone *birthday.Milestone) string {
	return birthday.ToStringWithWeekDay(milestone.Date)
}

func (v *view) Today(milestone *birthday.Milestone) bool {
	return milestone.Date.Equal(v.today)
}

func init() {
	kTemplate = common.NewTemplate("home", kTemplateSpec)
}
