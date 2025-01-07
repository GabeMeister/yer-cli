package presentation

import (
	helpers "GabeMeister/yer-cli/presentation/helpers"
	"GabeMeister/yer-cli/presentation/routes"
	presentation_views_pages "GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"
	"GabeMeister/yer-cli/utils"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

	routes.Init(e)

	/*
	 * PRESENTATION
	 */

	e.GET("/", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.Intro(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/:page/title", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		page := c.Param("page")

		titleSlideData := helpers.GetTitleSlideData(page, recap)
		component := presentation_views_pages.Title(presentation_views_pages.TitleParams{
			Title:       titleSlideData.Title,
			Description: titleSlideData.Description,
			NextBtnUrl:  titleSlideData.NextBtnUrl,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/shortest-commit-message-curr-year/title", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		nextBtnUrl := helpers.GetNextButtonLink("/shortest-commit-message-curr-year/title", recap)
		component := presentation_views_pages.Title(presentation_views_pages.TitleParams{
			Title:       "Shortest Commit Messages",
			Description: "The absolute shortest, low-effort commit messages engineers made this year.",
			NextBtnUrl:  nextBtnUrl,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)

	})

	e.GET("/new-engineer-count-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.NewEngineerCountCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/new-engineer-list-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.NewEngineerListCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/engineer-count-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCountCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/engineer-count-all-time", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCountAllTime(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/num-commits-prev-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.NumCommitsPrevYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/num-commits-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.NumCommitsCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/num-commits-all-time", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.NumCommitsAllTime(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-commits-over-time-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCommitsOverTimeCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-file-changes-over-time-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerFileChangesOverTimeCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/file-count-prev-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.FileCountPrevYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/file-count-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.FileCountCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/third-largest-file", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.ThirdLargestFile(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/second-largest-file", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.SecondLargestFile(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/largest-file", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.LargestFile(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/total-lines-of-code-prev-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.TotalLinesOfCodePrevYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/total-lines-of-code-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.TotalLinesOfCodeCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-single-day-commits-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.MostSingleDayCommitsByEngineerCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-single-day-commits-by-engineer-curr-year-commit-list", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.MostSingleDayCommitsByEngineerCurrYearCommitList(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-insertions-in-single-commit-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.MostInsertionsInSingleCommitCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-deletions-in-single-commit-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.MostDeletionsInSingleCommitCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/largest-commit-message-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.LargestCommitMessageCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/shortest-commit-message-curr-year/:index", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			panic(err)
		}

		component := presentation_views_pages.SmallestCommitMessagesCurrYear(recap, index)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-merges-in-one-day-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.MostMergesInOneDayCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-merges-in-one-day-commit-messages-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.MostMergesInOneDayCommitMessagesCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/avg-merges-per-day-to-master-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.AvgMergesToMasterPerDayCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/size-of-repo-by-week-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.SizeOfRepoByWeekCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/file-changes-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.FileChangesByEngineerCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/file-change-ratio-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.FileChangeRatioByEngineerCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-commit-counts-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCommitCountsCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-commit-counts-all-time", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCommitCountsAllTime(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/direct-pushes-on-master-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.DirectPushesOnMasterByEngineerCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/merges-to-master-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.MergesToMasterByEngineerCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/total-lines-of-code-in-repo-by-engineer", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.TotalLinesOfCodeInRepoByEngineer(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commits-by-month-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.CommitsByMonthCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commits-by-weekday-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.CommitsByWeekDayCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commits-by-hour-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.CommitsByHourCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commit-message-length-histogram-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.CommitMessageLengthHistogramCurrYear(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commonly-changed-files", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.CommonlyChangedFiles(recap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/end", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := presentation_views_pages.TheEnd()
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
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

		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	fmt.Println("\nDone! Browse to http://localhost:4000/")
	e.Logger.Fatal(e.Start(":4000"))
}
