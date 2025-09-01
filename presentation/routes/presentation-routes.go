package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"context"
	"fmt"
	"net/http"
	"time"

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

	e.GET("/shutdown", func(c echo.Context) error {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10)
			defer cancel()

			yellow := "\033[1;33m"
			reset := "\033[0m"

			fmt.Printf("%s┌──────────────────────────────────────┐%s\n", yellow, reset)
			fmt.Printf("%s│ Great! Now run the following command │%s\n", yellow, reset)
			fmt.Printf("%s│ to analyze your stats:               │%s\n", yellow, reset)
			fmt.Printf("%s│                                      │%s\n", yellow, reset)
			fmt.Printf("%s│ ./year-end-recap -a                  │%s\n", yellow, reset)
			fmt.Printf("%s└──────────────────────────────────────┘%s\n", yellow, reset)
			fmt.Println()

			time.Sleep(100 * time.Millisecond)
			e.Shutdown(ctx)
		}()

		// TODO: give nice message
		return c.String(http.StatusOK, "Go run ./year-end-recap -a")
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

	e.GET("/toc", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.TableOfContents()
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

	e.GET("/file-count-by-repo", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.FileCountByRepo(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/total-lines-of-code-by-repo", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.TotalLinesOfCodeByRepo(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/total-lines-of-code-per-week-by-repo", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.TotalLinesOfCodePerWeekByRepo(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/commits-made-by-repo", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.CommitsMadeByRepo(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/commits-made-by-author", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.CommitsMadeByAuthor(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/file-changes-made-by-author", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.FileChangesMadeByAuthor(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/lines-of-code-owned-by-author-all-time", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.LinesOfCodeOwnedByAuthorAllTime(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/aggregate-commits-by-month", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.AggregateCommitsByMonth(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/aggregate-commits-by-week-day", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.AggregateCommitsByWeekDay(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/aggregate-commits-by-hour", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.AggregateCommitsByHour(multiRepoRecap)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/merges-to-master-by-repo", func(c echo.Context) error {
		if !analyzer.HasRecapBeenRan() {
			return t.RenderRepoNotFound(c)
		}

		component := pages.MergesToMasterByRepo(multiRepoRecap)
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
