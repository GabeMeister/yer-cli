package presentation

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed views/*
var views embed.FS

//go:embed static/*
var static embed.FS

type Commits struct {
	Count int
}
type Repo struct {
	RepoName string
}

func RunLocalServer() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	summary := getSummary()

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, renderTemplate("views/index.html", []any{}))
	})
	e.GET("/presentation/intro", func(c echo.Context) error {
		name := c.QueryParam("name")
		return c.HTML(http.StatusOK, renderTemplate("views/repo.html", Repo{RepoName: name}))
	})
	e.GET("/presentation/commits", func(c echo.Context) error {
		return c.HTML(http.StatusOK, renderTemplate("views/commits.html", Commits{Count: summary.NumCommitsAllTime}))
	})
	// e.GET("/repo/:name", func(c echo.Context) error {
	// 	name := c.Param("name")
	// 	return c.HTML(http.StatusOK, renderTemplate("views/repo.html", Repo{RepoName: name}))
	// })
	e.GET("/favicon.ico", func(c echo.Context) error {
		data, _ := static.ReadFile("static/images/favicon.ico")
		return c.Blob(200, "image/x-icon", data)
	})
	e.GET("/css/styles.css", func(c echo.Context) error {
		data, _ := static.ReadFile("static/css/styles.css")
		return c.Blob(200, "text/css; charset=utf-8", data)
	})
	e.GET("/images/:name", func(c echo.Context) error {
		data, _ := static.ReadFile(fmt.Sprintf("static/images/%s", c.Param("name")))
		return c.Blob(200, "image/jpeg", data)
	})

	fmt.Println("\nDone! Browse to http://localhost:4000/")
	e.Logger.Fatal(e.Start(":4000"))
}
