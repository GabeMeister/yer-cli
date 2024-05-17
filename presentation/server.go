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

type Greeting struct {
	Name string
}
type Repo struct {
	RepoName string
}

func RunLocalServer() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, renderTemplate("views/index.html", Greeting{Name: "Josh"}))
	})
	e.GET("/repo", func(c echo.Context) error {
		name := c.QueryParam("name")
		return c.HTML(http.StatusOK, renderTemplate("views/repo.html", Repo{RepoName: name}))
	})
	e.GET("/repo/:name", func(c echo.Context) error {
		name := c.Param("name")
		return c.HTML(http.StatusOK, renderTemplate("views/repo.html", Repo{RepoName: name}))
	})
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

	e.Logger.Fatal(e.Start(":4000"))
}
