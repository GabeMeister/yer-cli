package routes

import (
	"GabeMeister/yer-cli/presentation/views/components"
	"GabeMeister/yer-cli/presentation/views/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RenderErrorMessage(c echo.Context, err error) error {
	content := template.Render(template.RenderParams{
		C: c,
		Component: components.ErrorMessage(components.ErrorMessageProps{
			Msg: err.Error(),
		}),
	})
	return c.HTML(http.StatusOK, content)
}

func RenderMessage(c echo.Context, msg string) error {
	content := template.Render(template.RenderParams{
		C: c,
		Component: components.ErrorMessage(components.ErrorMessageProps{
			Msg: msg,
		}),
	})
	return c.HTML(http.StatusOK, content)
}
