package presentation

import (
	"bytes"
	"html/template"
)

var templateGlobal *template.Template
var templateErr error

func getTemplate() *template.Template {
	if templateGlobal == nil {
		templateGlobal, templateErr = template.ParseFS(views, "views/layout/*", "views/pages/*", "views/partials/*")
		if templateErr != nil {
			panic(templateErr)
		}
	}

	return templateGlobal
}

type TemplateParams struct {
	path   string
	layout string
	data   any
}

func render(params TemplateParams) string {
	templ := getTemplate()
	buf := new(bytes.Buffer)

	err := templ.ExecuteTemplate(buf, params.path, params.data)
	if err != nil {
		panic(err)
	}

	content := buf.String()

	if params.layout != "" {
		buf = new(bytes.Buffer)
		err = templ.ExecuteTemplate(buf, params.layout, template.HTML(content))
		if err != nil {
			panic(err)
		}
	}

	return buf.String()
}
