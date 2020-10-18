package common

import (
	"html/template"
)

// NewTemplate returns a new template instance. name is the name
// of the template; templateStr is the template string.
func NewTemplate(name, templateStr string) *template.Template {
	return template.Must(template.New(name).Parse(templateStr))
}
