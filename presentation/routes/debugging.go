package routes

import (
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"
	"os"

	presentation_views_pages "GabeMeister/yer-cli/presentation/views/pages"

	"github.com/labstack/echo/v4"
)

func AddDebuggingRoutes(e *echo.Echo) {
	isDevMode := os.Getenv("DEV_MODE") == "true"

	e.GET("/env", func(c echo.Context) error {
		text := "Production"
		if isDevMode {
			text = "Development"
		}

		component := presentation_views_pages.Env(text)

		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})
}
