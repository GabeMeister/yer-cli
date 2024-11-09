package presentation

import (
	presentation_views_pages "GabeMeister/yer-cli/presentation/views/pages"
	"GabeMeister/yer-cli/utils"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

//go:embed views/*
var views embed.FS

//go:embed static/*
var static embed.FS

func RunLocalServer() {
	godotenv.Load()

	isDevMode := os.Getenv("DEV_MODE") == "true"

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	recap, _ := getRecap()

	e.GET("/hello", func(c echo.Context) error {
		component := presentation_views_pages.Hello("Dog", "36")
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			component := presentation_views_pages.RepoNotFound()
			content := render(RenderParams{
				c:         c,
				component: component,
			})

			return c.HTML(http.StatusOK, content)
		}

		component := presentation_views_pages.Intro(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/prev-year-commits", func(c echo.Context) error {
		type PrevYearCommitsView struct {
			RepoName string
			Count    int
		}
		content := renderOld(TemplateParams{
			c:    c,
			path: "pages/prev-year-commits.html",
			data: PrevYearCommitsView{Count: recap.NumCommitsPrevYear, RepoName: recap.Name},
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-commits-over-time", func(c echo.Context) error {
		type EngineerCommitsOverTimeView struct {
			Commits string
		}
		commitsOverTimeJson, err := json.Marshal(recap.EngineerCommitsOverTimeCurrYear)
		if err != nil {
			panic(err)
		}

		content := renderOld(TemplateParams{
			c:    c,
			path: "pages/engineer-commits-over-time.html",
			data: EngineerCommitsOverTimeView{Commits: string(commitsOverTimeJson)},
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-file-changes-over-time", func(c echo.Context) error {
		type EngineerFileChangesOverTimeView struct {
			FileChangesOverTime string
		}
		fileChangesOverTimeJson, err := json.Marshal(recap.EngineerFileChangesOverTimeCurrYear)
		if err != nil {
			panic(err)
		}

		content := renderOld(TemplateParams{
			c:    c,
			path: "pages/engineer-file-changes-over-time.html",
			data: EngineerFileChangesOverTimeView{FileChangesOverTime: string(fileChangesOverTimeJson)},
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	/*
	 * RESOURCES
	 */

	e.GET("/env", func(c echo.Context) error {
		text := "Production"
		if isDevMode {
			text = "Development"
		}

		component := presentation_views_pages.Env(text)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/favicon.ico", func(c echo.Context) error {
		data, _ := static.ReadFile("static/images/favicon.ico")
		return c.Blob(200, "image/x-icon", data)
	})

	e.GET("/css/styles.css", func(c echo.Context) error {
		var data []byte
		var err error

		data, err = static.ReadFile("static/css/styles.css")
		if err != nil {
			log.Fatal(err)
		}

		// Directly read from the file on disk when developing, so we can get the
		// fast hot module reloading for style tweaks, instead of fully rebuilding
		// the whole go app every time
		if isDevMode {
			data, err = os.ReadFile("presentation/static/css/styles.css")
			if err != nil {
				log.Fatal(err)
			}
		}

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
