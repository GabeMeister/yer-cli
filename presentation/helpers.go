package presentation

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RepoNotFoundPage(c echo.Context) error {
	return c.HTML(http.StatusOK, renderOld(TemplateParams{
		c:    c,
		path: "pages/repo-not-found.html",
	}))
}
