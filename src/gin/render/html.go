package render

import (
	"html/template"
	"net/http"
)

type (
	HTMLRender interface {
		Instance(string, interface{}) Render
	}

	HTMLProduction struct {
		Template *template.Template
	}

	HTMLDebug struct {
		Files []string
		Glob  string
	}

	HTML struct {
		Template *template.Template
		Name     string
		Data     interface{}
	}
)

var htmlContentType = []string{"text/html; chartset=utf-8"}

func (p HTMLProduction) Instance(name string, data interface{}) Render {
	return HTML{
		Template: p.Template,
		Name:     name,
		Data:     data,
	}
}

func (d HTMLDebug) Instance(name string, data interface{}) Render {
	return HTML{
		Template: d.loadTemplate(),
		Name:     name,
		Data:     data,
	}
}

func (d HTMLDebug) loadTemplate() *template.Template {
	if len(d.Files) > 0 {
		return template.Must(template.ParseFiles(d.Files...))
	}

	if len(d.Glob) > 0 {
		return template.Must(template.ParseGlob(d.Glob))
	}

	panic("the HTML debug render was created without files or glob pattern")
}

func (h HTML) Write(w http.ResponseWriter) error {
	w.Header()["Content-Type"] = htmlContentType
	if len(h.Name) == 0 {
		return h.Template.Execute(w, h.Data)
	} else {
		return h.Template.ExecuteTemplate(w, h.Name, h.Data)
	}
}
