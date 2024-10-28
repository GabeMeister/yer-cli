package presentation

import (
	"GabeMeister/yer-cli/utils"
	"embed"
	"encoding/json"
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

	recap, _ := getRecap()

	e.GET("/", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return RepoNotFoundPage(c)
		}

		type NextButton struct {
			Href string
		}

		dateStr, err := utils.FormatISODate(recap.DateAnalyzed)
		if err != nil {
			panic(err)
		}

		data := struct {
			Title        string
			DateAnalyzed string
			NextButton   NextButton
		}{
			Title:        recap.Name,
			DateAnalyzed: dateStr,
			NextButton: NextButton{
				Href: "/prev-year-commits",
			},
		}

		content := render(TemplateParams{
			c:    c,
			path: "pages/intro.html",
			data: data,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/prev-year-commits", func(c echo.Context) error {
		type PrevYearCommitsView struct {
			RepoName string
			Count    int
		}
		content := render(TemplateParams{
			c:    c,
			path: "pages/prev-year-commits.html",
			data: PrevYearCommitsView{Count: recap.NumCommitsPrevYear, RepoName: recap.Name},
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	// e.GET("/presentation/curr-year-commits", func(c echo.Context) error {
	// 	type CurrYearCommitsView struct {
	// 		RepoName string
	// 		Count    int
	// 	}
	// 	return c.HTML(
	// 		http.StatusOK,
	// 		renderTemplate(
	// 			"views/curr-year-commits.html",
	// 			CurrYearCommitsView{Count: recap.NumCommitsCurrYear, RepoName: recap.Name}))
	// })

	// e.GET("/presentation/all-time-commits", func(c echo.Context) error {
	// 	type AllTimeCommitsView struct {
	// 		RepoName string
	// 		Count    int
	// 	}
	// 	return c.HTML(
	// 		http.StatusOK,
	// 		renderTemplate(
	// 			"views/all-time-commits.html",
	// 			AllTimeCommitsView{Count: recap.NumCommitsAllTime, RepoName: recap.Name}))
	// })

	e.GET("/engineer-commits-over-time", func(c echo.Context) error {
		type EngineerCommitsOverTimeView struct {
			Commits string
		}
		commitsOverTimeJson, err := json.Marshal(recap.EngineerCommitsOverTimeCurrYear)
		if err != nil {
			panic(err)
		}

		content := render(TemplateParams{
			c:    c,
			path: "pages/engineer-commits-over-time.html",
			data: EngineerCommitsOverTimeView{Commits: string(commitsOverTimeJson)},
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	/*
	 * RESOURCES
	 */

	e.GET("/example", func(c echo.Context) error {
		content := render(TemplateParams{
			c:    c,
			path: "pages/example.html",
		})

		return c.HTML(http.StatusOK, content)
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
		// TODO: figure out how to return any kind of image
		return c.Blob(200, "image/jpeg", data)
	})

	e.GET("/scripts/:name", func(c echo.Context) error {
		data, _ := static.ReadFile(fmt.Sprintf("static/scripts/%s", c.Param("name")))
		return c.Blob(200, "text/javascript", data)
	})

	fmt.Println("\nDone! Browse to http://localhost:4000/")
	e.Logger.Fatal(e.Start(":4000"))
}
