package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"net/http"

	helpers "GabeMeister/yer-cli/presentation/helpers"
	"GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"

	"github.com/labstack/echo/v4"
)

func addPresentationRoutes(e *echo.Echo) {
	multiRepoRecap, _ := analyzer.GetMultiRepoRecapFromTmpDir()

	e.GET("/", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.Intro(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/:page/title", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		page := c.Param("page")

		titleSlideData := helpers.GetMultiRepoTitleSlideData(page, multiRepoRecap)
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

	e.GET("/active-authors", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.ActiveAuthors(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	// e.GET("/shortest-commit-message-curr-year/title", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	nextBtnUrl := helpers.GetNextButtonLink("/shortest-commit-message-curr-year/title", recap)
	// 	component := pages.Title(pages.TitleParams{
	// 		Title:       "Shortest Commit Messages",
	// 		Description: "The absolute shortest, low-effort commit messages authors made this year.",
	// 		NextBtnUrl:  nextBtnUrl,
	// 	})
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)

	// })

	// e.GET("/new-author-count-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.NewAuthorCountCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)
	// })

	// e.GET("/new-author-list-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.NewAuthorListCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)
	// })

	// e.GET("/author-count-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.AuthorCountCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)
	// })

	// e.GET("/author-count-all-time", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.AuthorCountAllTime(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)
	// })

	// e.GET("/num-commits-prev-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.NumCommitsPrevYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/num-commits-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.NumCommitsCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/num-commits-all-time", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.NumCommitsAllTime(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/author-file-changes-over-time-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.AuthorFileChangesOverTimeCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/file-count-prev-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.FileCountPrevYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/file-count-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.FileCountCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/third-largest-file", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.ThirdLargestFile(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/second-largest-file", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.SecondLargestFile(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/largest-file", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.LargestFile(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/total-lines-of-code-prev-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.TotalLinesOfCodePrevYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/total-lines-of-code-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.TotalLinesOfCodeCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/most-single-day-commits-by-author-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.MostSingleDayCommitsByAuthorCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/most-single-day-commits-by-author-curr-year-commit-list", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.MostSingleDayCommitsByAuthorCurrYearCommitList(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/most-insertions-in-single-commit-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.MostInsertionsInSingleCommitCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/most-deletions-in-single-commit-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.MostDeletionsInSingleCommitCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/largest-commit-message-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.LargestCommitMessageCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/shortest-commit-message-curr-year/:index", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	index, err := strconv.Atoi(c.Param("index"))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	component := pages.SmallestCommitMessagesCurrYear(recap, index)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/most-merges-in-one-day-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.MostMergesInOneDayCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/most-merges-in-one-day-commit-messages-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.MostMergesInOneDayCommitMessagesCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/avg-merges-per-day-to-master-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.AvgMergesToMasterPerDayCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/size-of-repo-by-week-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.SizeOfRepoByWeekCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/file-changes-by-author-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.FileChangesByAuthorCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/file-change-ratio-by-author-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.FileChangeRatioByAuthorCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/author-commit-counts-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.AuthorCommitCountsCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/author-commit-counts-all-time", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.AuthorCommitCountsAllTime(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/direct-pushes-on-master-by-author-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.DirectPushesOnMasterByAuthorCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/merges-to-master-by-author-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.MergesToMasterByAuthorCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/total-lines-of-code-in-repo-by-author", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.TotalLinesOfCodeInRepoByAuthor(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/commits-by-month-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.CommitsByMonthCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/commits-by-weekday-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.CommitsByWeekDayCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/commits-by-hour-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.CommitsByHourCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/commit-message-length-histogram-curr-year", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.CommitMessageLengthHistogramCurrYear(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	// e.GET("/commonly-changed-files", func(c echo.Context) error {
	// 	if !analyzer.HasRecapBeenRan() {
	// 		return t.RenderRepoNotFound(c)
	// 	}

	// 	component := pages.CommonlyChangedFiles(recap)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(
	// 		http.StatusOK,
	// 		content,
	// 	)
	// })

	e.GET("/end", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
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
