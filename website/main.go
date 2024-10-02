package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"title":   "Welcome to my Golang Echo Web Server",
			"message": "Hello, World!",
		})
	})

	e.GET("/install", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "https://fly.storage.tigris.dev/year-end-recap-storage/install")
	})

	e.Start(":8080")
}
