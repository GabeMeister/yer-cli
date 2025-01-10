package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"net/http"
	"strconv"

	helpers "GabeMeister/yer-cli/presentation/helpers"
	"GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"

	"github.com/labstack/echo/v4"
)

func addPresentationRoutes(e *echo.Echo) {
	recap, _ := analyzer.GetRecap()

	e.GET("/", func(c echo.Context) error {
		if !utils.HasRepoBeenAnalyzed() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.Intro(recap)
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
		component := pages.Title(pages.TitleParams{
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
		component := pages.Title(pages.TitleParams{
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

		component := pages.NewEngineerCountCurrYear(recap)
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

		component := pages.NewEngineerListCurrYear(recap)
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

		component := pages.EngineerCountCurrYear(recap)
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

		component := pages.EngineerCountAllTime(recap)
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

		component := pages.NumCommitsPrevYear(recap)
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

		component := pages.NumCommitsCurrYear(recap)
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

		component := pages.NumCommitsAllTime(recap)
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

		component := pages.EngineerCommitsOverTimeCurrYear(recap)
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

		component := pages.EngineerFileChangesOverTimeCurrYear(recap)
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

		component := pages.FileCountPrevYear(recap)
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

		component := pages.FileCountCurrYear(recap)
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

		component := pages.ThirdLargestFile(recap)
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

		component := pages.SecondLargestFile(recap)
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

		component := pages.LargestFile(recap)
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

		component := pages.TotalLinesOfCodePrevYear(recap)
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

		component := pages.TotalLinesOfCodeCurrYear(recap)
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

		component := pages.MostSingleDayCommitsByEngineerCurrYear(recap)
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

		component := pages.MostSingleDayCommitsByEngineerCurrYearCommitList(recap)
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

		component := pages.MostInsertionsInSingleCommitCurrYear(recap)
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

		component := pages.MostDeletionsInSingleCommitCurrYear(recap)
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

		component := pages.LargestCommitMessageCurrYear(recap)
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

		component := pages.SmallestCommitMessagesCurrYear(recap, index)
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

		component := pages.MostMergesInOneDayCurrYear(recap)
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

		component := pages.MostMergesInOneDayCommitMessagesCurrYear(recap)
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

		component := pages.AvgMergesToMasterPerDayCurrYear(recap)
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

		component := pages.SizeOfRepoByWeekCurrYear(recap)
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

		component := pages.FileChangesByEngineerCurrYear(recap)
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

		component := pages.FileChangeRatioByEngineerCurrYear(recap)
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

		component := pages.EngineerCommitCountsCurrYear(recap)
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

		component := pages.EngineerCommitCountsAllTime(recap)
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

		component := pages.DirectPushesOnMasterByEngineerCurrYear(recap)
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

		component := pages.MergesToMasterByEngineerCurrYear(recap)
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

		component := pages.TotalLinesOfCodeInRepoByEngineer(recap)
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

		component := pages.CommitsByMonthCurrYear(recap)
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

		component := pages.CommitsByWeekDayCurrYear(recap)
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

		component := pages.CommitsByHourCurrYear(recap)
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

		component := pages.CommitMessageLengthHistogramCurrYear(recap)
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

		component := pages.CommonlyChangedFiles(recap)
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

		component := pages.TheEnd()
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})
}
