package presentation

import (
	"bytes"
	"html/template"
)

func renderTemplate(path string, data interface{}) string {
	htmlStr, _ := views.ReadFile(path)
	t := template.Must(template.New(path).Parse(string(htmlStr)))

	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
