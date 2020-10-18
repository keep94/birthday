package home

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/keep94/birthday"
	"github.com/keep94/birthday/cmd/remind/common"
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
      <td {{if $top.Today .}}class="today"{{end}}>{{.Date.StringWithWeekDay}}</td>
      <td {{if $top.Today .}}class="today"{{end}}>{{.Name}}</td>
      <td {{if $top.Today .}}class="today"{{end}}>{{$top.AgeStr .}}</td>
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
	milestones, err := birthday.ReadFile(h.File, birthday.Now(), h.DaysAhead)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	http_util.WriteTemplate(w, kTemplate, &view{Milestones: milestones})
}

type view struct {
	Milestones []birthday.Milestone
}

func (v *view) AgeStr(milestone birthday.Milestone) string {
	if milestone.Age < 0 {
		return "? Years"
	}
	if milestone.AgeInDays {
		return fmt.Sprintf("%d Days", milestone.Age)
	}
	return fmt.Sprintf("%d Years", milestone.Age)
}

func (v *view) Today(milestone birthday.Milestone) bool {
	return milestone.DaysAway == 0
}

func init() {
	kTemplate = common.NewTemplate("home", kTemplateSpec)
}
