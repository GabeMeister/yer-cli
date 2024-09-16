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

func RunLocalServer() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	recap := getRecap()

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, renderTemplate("views/index.html", []any{}))
	})

	e.GET("/presentation/intro", func(c echo.Context) error {
		fmt.Println()
		type IntroView struct {
			RepoName   string
			ShowLayout bool
		}
		return c.HTML(
			http.StatusOK,
			renderTemplate(
				"views/repo.html",
				IntroView{RepoName: recap.Name, ShowLayout: c.Request().Header.Get("HX-Request") != "true"}))
	})

	e.GET("/presentation/prev-year-commits", func(c echo.Context) error {
		type PrevYearCommitsView struct {
			RepoName string
			Count    int
		}
		return c.HTML(
			http.StatusOK,
			renderTemplate(
				"views/prev-year-commits.html",
				PrevYearCommitsView{Count: recap.NumCommitsPrevYear, RepoName: recap.Name}))
	})

	e.GET("/presentation/curr-year-commits", func(c echo.Context) error {
		type CurrYearCommitsView struct {
			RepoName string
			Count    int
		}
		return c.HTML(
			http.StatusOK,
			renderTemplate(
				"views/curr-year-commits.html",
				CurrYearCommitsView{Count: recap.NumCommitsCurrYear, RepoName: recap.Name}))
	})

	e.GET("/presentation/all-time-commits", func(c echo.Context) error {
		type AllTimeCommitsView struct {
			RepoName string
			Count    int
		}
		return c.HTML(
			http.StatusOK,
			renderTemplate(
				"views/all-time-commits.html",
				AllTimeCommitsView{Count: recap.NumCommitsAllTime, RepoName: recap.Name}))
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

	e.GET("/scripts/:name", func(c echo.Context) error {
		data, _ := static.ReadFile(fmt.Sprintf("static/scripts/%s", c.Param("name")))
		return c.Blob(200, "text/javascript", data)
	})

	fmt.Println("\nDone! Browse to http://localhost:4000/")
	e.Logger.Fatal(e.Start(":4000"))
}
