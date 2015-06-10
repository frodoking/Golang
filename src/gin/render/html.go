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

func (this *HTMLProduction) Instance(name string, data interface{}) Render {
	return &HTML{
		Template: this.Template,
		Name:     name,
		Data:     data,
	}
}

func (this *HTMLDebug) Instance(name string, data interface{}) Render {
	return &HTML{
		Template: this.loadTemplate(),
		Name:     name,
		Data:     data,
	}
}

func (this *HTMLDebug) loadTemplate() *template.Template {
	if len(this.Files) > 0 {
		return template.Must(template.ParseFiles(this.Files...))
	}

	if len(this.Glob) > 0 {
		return template.Must(template.ParseGlob(this.Glob))
	}

	panic("the HTML debug render was created without files or glob pattern")
}

func (this *HTML) Write(w http.ResponseWriter) error {
	w.Header()["Content-Type"] = htmlContentType
	if len(this.Name) == 0 {
		return this.Template.Execute(w, this.Data)
	} else {
		return this.Template.ExecuteTemplate(w, this.Name, this.Data)
	}
}
