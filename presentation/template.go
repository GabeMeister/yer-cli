package presentation

import (
	presentation_views_layouts "GabeMeister/yer-cli/presentation/views/layouts"
	"bytes"
	"context"
	"html/template"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

var templateGlobal *template.Template
var templateErr error

func getTemplate() *template.Template {
	if templateGlobal == nil {
		templateGlobal, templateErr = template.ParseFS(views, "views/layouts/*", "views/pages/*", "views/partials/*")
		if templateErr != nil {
			panic(templateErr)
		}
	}

	return templateGlobal
}

type TemplateParams struct {
	c      echo.Context
	path   string
	layout string
	data   any
}

func renderOld(params TemplateParams) string {
	c := params.c

	if c == nil {
		panic("Need to pass echo.Context to render() function")
	}

	path := params.path
	if path == "" {
		panic("Need to pass path to render() function")
	}

	buf := new(bytes.Buffer)
	templ := getTemplate()
	data := params.data

	err := templ.ExecuteTemplate(buf, path, data)
	if err != nil {
		panic(err)
	}

	htmlContent := buf.String()

	htmxRequestHeader := c.Request().Header["Hx-Request"]
	isHtmxRequest := len(htmxRequestHeader) > 0 && htmxRequestHeader[0] == "true"

	layout := params.layout
	if layout == "" {
		// Always default to standard if not specified
		layout = "layouts/standard.html"
	}

	if isHtmxRequest {
		// If it's specifically an htmx request, we already have the boilerplate of
		// the page, so we only care about the content. No need for rendering the layout
		layout = ""
	}

	if layout != "" {
		buf = new(bytes.Buffer)
		err = templ.ExecuteTemplate(buf, layout, template.HTML(htmlContent))
		if err != nil {
			panic(err)
		}
	}

	return buf.String()
}

type RenderParams struct {
	c         echo.Context
	component templ.Component
}

func render(params RenderParams) string {
	component := params.component
	c := params.c

	htmxRequestHeader := c.Request().Header["Hx-Request"]
	isHtmxRequest := len(htmxRequestHeader) > 0 && htmxRequestHeader[0] == "true"
	buf := templ.GetBuffer()

	if isHtmxRequest {
		// If it's an Htmx request, then that means the headers/styling has already
		// loaded, so no need to add that into the response again
		err := component.Render(context.Background(), buf)
		if err != nil {
			panic(err)
		}

		return buf.String()
	} else {
		ctx := templ.WithChildren(context.Background(), component)
		standardLayout := presentation_views_layouts.StandardLayout()
		err := standardLayout.Render(ctx, buf)
		if err != nil {
			panic(err)
		}

		return buf.String()
	}
}
