package presentation

import (
	presentation_views_pages "GabeMeister/yer-cli/presentation/views/pages"
	"GabeMeister/yer-cli/utils"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

//go:embed static/*
var static embed.FS

func RunLocalServer() {
	godotenv.Load()

	isDevMode := os.Getenv("DEV_MODE") == "true"

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	recap, _ := getRecap()

	/*
	 * PRESENTATION
	 */

	e.GET("/", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.Intro(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/new-engineer-count-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.NewEngineerCountCurrYear(recap.NewEngineerCountCurrYear)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/engineer-count-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCountCurrYear(recap.EngineerCountAllTime)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/engineer-count-all-time", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCountAllTime(recap.EngineerCountAllTime, recap.Name)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/num-commits-prev-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.NumCommitsPrevYear(recap.NumCommitsPrevYear)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/num-commits-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.NumCommitsCurrYear(recap.NumCommitsCurrYear)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/num-commits-all-time", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.NumCommitsAllTime(recap.NumCommitsAllTime)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-commits-over-time", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCommitsOverTimeCurrYear(recap.EngineerCommitsOverTimeCurrYear)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-file-changes-over-time", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerFileChangesOverTimeCurrYear(recap.EngineerFileChangesOverTimeCurrYear)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/file-count-prev-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.FileCountPrevYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/file-count-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.FileCountCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/third-largest-file", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.FileCountCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	/*
	 * RESOURCES
	 */

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

	/*
	 * DEBUGGING
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

	fmt.Println("\nDone! Browse to http://localhost:4000/")
	e.Logger.Fatal(e.Start(":4000"))
}
