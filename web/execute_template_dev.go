// +build dev

package web

import (
	"html/template"
)

var (
	templates = template.Must(template.ParseGlob("../web/templates/*.go.html"))
)

func LookupTemplate(name string) *template.Template {
	return templates.Lookup(name + ".go.html")
}
