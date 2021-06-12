package search

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/keep94/birthday"
	"github.com/keep94/birthday/cmd/remind/common"
	"github.com/keep94/consume"
	"github.com/keep94/toolbox/http_util"
)

var (
	kYears  = birthday.Period{Years: 1}
	kMonths = birthday.Period{Months: 1}
	kWeeks  = birthday.Period{Weeks: 1}
	kDays   = birthday.Period{Days: 1}
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
  form {
    font-size: 30px;
  }
  input {
    font-size: 30px;
  }
  </style>
</head>
<body>
  <h1>Birthdays</h1>
  <form>
     Name: <input type="text" name="q" value="{{.Get "q"}}">
    <input type="submit" value="Search">
  </form>
  <hr>
  <table border=1>
    <tr>
      <th>Name</th>
      <th>Birthday</th>
      <th>Years</th>
      <th>Months</th>
      <th>Weeks</th>
      <th>Days</th>
    </tr>
    {{with $top := .}}
    {{range .Results}}
    <tr>
      <td>{{.Name}}</td>
      <td>{{$top.BirthdayStr .}}</td>
      <td>{{$top.InYearsStr .}}</td>
      <td>{{$top.InMonthsStr .}}</td>
      <td>{{$top.InWeeksStr .}}</td>
      <td>{{$top.InDaysStr .}}</td>
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
	File string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var entries []birthday.Entry
	err := birthday.ReadFile(
		h.File,
		consume.MapFilter(
			consume.AppendTo(&entries), birthday.Query(r.Form.Get("q"))))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	http_util.WriteTemplate(w, kTemplate, &view{
		Values:      http_util.Values{r.Form},
		Results:     birthday.EntriesSortedByName(entries),
		CurrentDate: common.ParseDate(r.Form.Get("date")),
	})
}

type view struct {
	http_util.Values
	Results     []birthday.Entry
	CurrentDate time.Time
}

func (b *view) BirthdayStr(entry *birthday.Entry) string {
	return birthday.ToString(entry.Birthday)
}

func (v *view) InYearsStr(entry *birthday.Entry) string {
	return v.inUnitStr(entry, kYears)
}

func (v *view) InMonthsStr(entry *birthday.Entry) string {
	return v.inUnitStr(entry, kMonths)
}

func (v *view) InWeeksStr(entry *birthday.Entry) string {
	return v.inUnitStr(entry, kWeeks)
}

func (v *view) InDaysStr(entry *birthday.Entry) string {
	return v.inUnitStr(entry, kDays)
}

func (v *view) inUnitStr(
	entry *birthday.Entry,
	period birthday.Period) string {
	if birthday.HasYear(entry.Birthday) {
		return strconv.Itoa(period.Diff(v.CurrentDate, entry.Birthday))
	}
	return "--"
}

func init() {
	kTemplate = common.NewTemplate("search", kTemplateSpec)
}
