package search

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/keep94/birthday"
	"github.com/keep94/birthday/cmd/remind/common"
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
      <th>In Years</th>
      <th>In Days</th>
    </tr>
    {{with $top := .}}
    {{range .Results}}
    <tr>
      <td>{{.Name}}</td>
      <td>{{$top.BirthdayStr .}}</td>
      <td>{{$top.InYearsStr .}}</td>
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
	search := birthday.NewSearch(birthday.Today(), r.Form.Get("q"))
	err := birthday.ReadFile(h.File, search)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	http_util.WriteTemplate(w, kTemplate, &view{
		Values:  http_util.Values{r.Form},
		Results: search.Results(),
	})
}

type view struct {
	http_util.Values
	Results []birthday.Result
}

func (b *view) BirthdayStr(result birthday.Result) string {
	return birthday.ToString(result.Birthday)
}

func (v *view) InYearsStr(result birthday.Result) string {
	if birthday.HasYear(result.Birthday) {
		return strconv.Itoa(result.AgeInYears)
	}
	return "--"
}

func (v *view) InDaysStr(result birthday.Result) string {
	if birthday.HasYear(result.Birthday) {
		return strconv.Itoa(result.AgeInDays)
	}
	return "--"
}

func init() {
	kTemplate = common.NewTemplate("search", kTemplateSpec)
}
