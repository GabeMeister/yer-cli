package presentation

import (
	presentation_views_layouts "GabeMeister/yer-cli/presentation/views/layouts"
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

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
