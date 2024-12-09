package presentation

import (
	"GabeMeister/yer-cli/analyzer"
	presentation_helpers "GabeMeister/yer-cli/presentation/helpers"
	presentation_views_pages "GabeMeister/yer-cli/presentation/views/pages"
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

	e.GET("/:page/title", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		page := c.Param("page")

		pageToTitleMap := map[string]string{
			"new-engineer-count-curr-year":                  fmt.Sprintf("New Engineer Count (%d)", analyzer.CURR_YEAR),
			"engineer-count-curr-year":                      fmt.Sprintf("Number of Engineers who committed to %s (%d)", recap.Name, analyzer.CURR_YEAR),
			"engineer-count-all-time":                       fmt.Sprintf("Number of Engineers who committed to %s (all time)", recap.Name),
			"file-count-prev-year":                          fmt.Sprintf("File Count (%d)", analyzer.PREV_YEAR),
			"file-count-curr-year":                          fmt.Sprintf("File Count (%d)", analyzer.CURR_YEAR),
			"third-largest-file":                            "Third Largest File",
			"second-largest-file":                           "Second Largest File",
			"largest-file":                                  "Largest File",
			"total-lines-of-code-prev-year":                 fmt.Sprintf("Total Lines of Code (%d)", analyzer.PREV_YEAR),
			"total-lines-of-code-curr-year":                 fmt.Sprintf("Total Lines of Code (%d)", analyzer.CURR_YEAR),
			"size-of-repo-by-week-curr-year":                fmt.Sprintf("Size of Repo by Week (%d)", analyzer.CURR_YEAR),
			"total-lines-of-code-in-repo-by-engineer":       "Total Lines of Code by Engineer",
			"file-changes-by-engineer-curr-year":            fmt.Sprintf("File Changes by Engineer (%d)", analyzer.CURR_YEAR),
			"file-change-ratio-by-engineer-curr-year":       fmt.Sprintf("Insertions/Deletions Ratio by Engineer (%d)", analyzer.CURR_YEAR),
			"commonly-changed-files":                        "Most Commonly Changed Files",
			"num-commits-prev-year":                         fmt.Sprintf("Number of Commits (%d)", analyzer.PREV_YEAR),
			"num-commits-curr-year":                         fmt.Sprintf("Number of Commits (%d)", analyzer.CURR_YEAR),
			"num-commits-all-time":                          "Number of Commits (all time)",
			"engineer-commits-over-time-curr-year":          fmt.Sprintf("Engineer Commits Over Time (%d)", analyzer.CURR_YEAR),
			"engineer-file-changes-over-time-curr-year":     fmt.Sprintf("Engineer File Changes Over Time (%d)", analyzer.CURR_YEAR),
			"engineer-commit-counts-curr-year":              fmt.Sprintf("Engineer Commit Counts (%d)", analyzer.CURR_YEAR),
			"engineer-commit-counts-all-time":               "Engineer Commit Counts (all time)",
			"commits-by-month-curr-year":                    fmt.Sprintf("Commits by Month (%d)", analyzer.CURR_YEAR),
			"commits-by-weekday-curr-year":                  fmt.Sprintf("Commits by Weekday (%d)", analyzer.CURR_YEAR),
			"commits-by-hour-curr-year":                     fmt.Sprintf("Commits by Hour (%d)", analyzer.CURR_YEAR),
			"most-single-day-commits-by-engineer-curr-year": fmt.Sprintf("Most Single-Day Commits by Engineer (%d)", analyzer.CURR_YEAR),
			"most-insertions-in-single-commit-curr-year":    fmt.Sprintf("Most Insertions in a Single Commit (%d)", analyzer.CURR_YEAR),
			"most-deletions-in-single-commit-curr-year":     fmt.Sprintf("Most Deletions in a Single Commit (%d)", analyzer.CURR_YEAR),
			"largest-commit-message-curr-year":              fmt.Sprintf("Largest Commit Message (%d)", analyzer.CURR_YEAR),
			"shortest-commit-message-curr-year":             fmt.Sprintf("Shortest Commit Message (%d)", analyzer.CURR_YEAR),
			"commit-message-length-histogram-curr-year":     fmt.Sprintf("Commit Message Length Frequencies (%d)", analyzer.CURR_YEAR),
			"direct-pushes-on-master-by-engineer-curr-year": fmt.Sprintf("Direct Pushes on Master by Engineer (%d)", analyzer.CURR_YEAR),
			"merges-to-master-by-engineer-curr-year":        fmt.Sprintf("Merges to Master by Engineer (%d)", analyzer.CURR_YEAR),
			"most-merges-in-one-day-curr-year":              fmt.Sprintf("Most Merges in One Day (%d)", analyzer.CURR_YEAR),
			"avg-merges-per-day-to-master-curr-year":        fmt.Sprintf("Average Merges per Day to Master (%d)", analyzer.CURR_YEAR),
		}

		nextBtnUrl := presentation_helpers.GetNextButtonLink(fmt.Sprintf("/%s/title", page), recap)
		component := presentation_views_pages.Title(
			pageToTitleMap[page],
			nextBtnUrl,
		)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/shortest-commit-message-curr-year/title", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		nextBtnUrl := presentation_helpers.GetNextButtonLink("/shortest-commit-message-curr-year/title", recap)
		component := presentation_views_pages.Title(
			"Shortest Commit Messages",
			nextBtnUrl,
		)
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

		component := presentation_views_pages.NewEngineerCountCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/new-engineer-list-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.NewEngineerListCurrYear(recap)
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

		component := presentation_views_pages.EngineerCountCurrYear(recap)
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

		component := presentation_views_pages.EngineerCountAllTime(recap)
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

		component := presentation_views_pages.NumCommitsPrevYear(recap)
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

		component := presentation_views_pages.NumCommitsCurrYear(recap)
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

		component := presentation_views_pages.NumCommitsAllTime(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-commits-over-time-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCommitsOverTimeCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-file-changes-over-time-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerFileChangesOverTimeCurrYear(recap)
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

		component := presentation_views_pages.ThirdLargestFile(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/second-largest-file", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.SecondLargestFile(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/largest-file", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.LargestFile(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/total-lines-of-code-prev-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.TotalLinesOfCodePrevYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/total-lines-of-code-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.TotalLinesOfCodeCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-single-day-commits-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.MostSingleDayCommitsByEngineerCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-single-day-commits-by-engineer-curr-year-commit-list", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.MostSingleDayCommitsByEngineerCurrYearCommitList(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-insertions-in-single-commit-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.MostInsertionsInSingleCommitCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-deletions-in-single-commit-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.MostDeletionsInSingleCommitCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/largest-commit-message-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.LargestCommitMessageCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/shortest-commit-message-curr-year/:index", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			panic(err)
		}

		component := presentation_views_pages.SmallestCommitMessagesCurrYear(recap, index)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-merges-in-one-day-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.MostMergesInOneDayCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/most-merges-in-one-day-commit-messages-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.MostMergesInOneDayCommitMessagesCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/avg-merges-per-day-to-master-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.AvgMergesToMasterPerDayCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/size-of-repo-by-week-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.SizeOfRepoByWeekCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/file-changes-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.FileChangesByEngineerCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/file-change-ratio-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.FileChangeRatioByEngineerCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-commit-counts-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCommitCountsCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/engineer-commit-counts-all-time", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.EngineerCommitCountsAllTime(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/direct-pushes-on-master-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.DirectPushesOnMasterByEngineerCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/merges-to-master-by-engineer-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.MergesToMasterByEngineerCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/total-lines-of-code-in-repo-by-engineer", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.TotalLinesOfCodeInRepoByEngineer(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commits-by-month-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.CommitsByMonthCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commits-by-weekday-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.CommitsByWeekDayCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commits-by-hour-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.CommitsByHourCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commit-message-length-histogram-curr-year", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.CommitMessageLengthHistogramCurrYear(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/commonly-changed-files", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.CommonlyChangedFiles(recap)
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(
			http.StatusOK,
			content,
		)
	})

	e.GET("/end", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return renderRepoNotFound(c)
		}

		component := presentation_views_pages.TheEnd()
		content := render(RenderParams{
			c:         c,
			component: component,
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
		content := render(RenderParams{
			c:         c,
			component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	fmt.Println("\nDone! Browse to http://localhost:4000/")
	e.Logger.Fatal(e.Start(":4000"))
}
