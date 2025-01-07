package template

import (
	presentation_views_layouts "GabeMeister/yer-cli/presentation/views/layouts"
	presentation_views_pages "GabeMeister/yer-cli/presentation/views/pages"
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type RenderParams struct {
	C         echo.Context
	Component templ.Component
}

func Render(params RenderParams) string {
	component := params.Component
	c := params.C

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

func RenderRepoNotFound(c echo.Context) error {
	component := presentation_views_pages.RepoNotFound()
	content := Render(RenderParams{
		C:         c,
		Component: component,
	})

	return c.HTML(http.StatusOK, content)
}
