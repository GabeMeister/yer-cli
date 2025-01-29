package routes

import (
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"
	"os"

	"GabeMeister/yer-cli/presentation/views/pages"

	"github.com/labstack/echo/v4"
)

func addDebuggingRoutes(e *echo.Echo) {
	isDevMode := os.Getenv("DEV_MODE") == "true"

	e.GET("/env", func(c echo.Context) error {
		text := "Production"
		if isDevMode {
			text = "Development"
		}

		component := pages.Env(text)

		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/test", func(c echo.Context) error {
		component := pages.Test()

		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})
}
