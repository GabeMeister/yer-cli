package routes

import (
	"GabeMeister/yer-cli/presentation/views/components"
	t "GabeMeister/yer-cli/presentation/views/template"
	"GabeMeister/yer-cli/utils"
	"net/http"
	"strings"
	"time"

	"GabeMeister/yer-cli/presentation/views/pages"

	"github.com/labstack/echo/v4"
)

func addDebuggingRoutes(e *echo.Echo) {

	e.GET("/env", func(c echo.Context) error {
		text := "Production"
		if utils.IsDevMode() {
			text = "Development"
		}

		component := pages.Env(text)

		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/example", func(c echo.Context) error {
		time.Sleep(700 * time.Millisecond)

		return c.HTML(http.StatusOK, "<div>Success!</div>")
	})

	e.GET("/buttons", func(c echo.Context) error {
		component := pages.ButtonsPage()

		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/fade-in", func(c echo.Context) error {
		component := pages.FadeIn()

		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/fade-in-request", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<div class='bg-green-400'>Hey there!</div>")
	})

	e.GET("/hx-swap-oob", func(c echo.Context) error {
		component := pages.HxSwapOobExample(pages.HxSwapOobExampleProps{
			Animal: "",
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/animals-example", func(c echo.Context) error {
		formParams, _ := c.FormParams()

		component := components.HxSwapOobExampleContent(components.HxSwapOobExampleContentProps{
			Animal: strings.Join(formParams["animal"], ","),
			Food:   strings.Join(formParams["food"], ","),
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

}
