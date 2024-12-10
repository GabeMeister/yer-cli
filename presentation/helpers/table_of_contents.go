package presentation_helpers

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"fmt"
	"slices"
)

var TABLE_OF_CONTENTS = []string{
	"/",
	"/new-engineer-count-curr-year/title",
	"/new-engineer-count-curr-year",
	"/new-engineer-list-curr-year",
	"/engineer-count-curr-year/title",
	"/engineer-count-curr-year",
	"/engineer-count-all-time/title",
	"/engineer-count-all-time",
	"/file-count-prev-year/title",
	"/file-count-prev-year",
	"/file-count-curr-year/title",
	"/file-count-curr-year",
	"/third-largest-file/title",
	"/third-largest-file",
	"/second-largest-file/title",
	"/second-largest-file",
	"/largest-file/title",
	"/largest-file",
	"/total-lines-of-code-prev-year/title",
	"/total-lines-of-code-prev-year",
	"/total-lines-of-code-curr-year/title",
	"/total-lines-of-code-curr-year",
	"/size-of-repo-by-week-curr-year/title",
	"/size-of-repo-by-week-curr-year",
	"/file-changes-by-engineer-curr-year/title",
	"/file-changes-by-engineer-curr-year",
	"/total-lines-of-code-in-repo-by-engineer/title",
	"/total-lines-of-code-in-repo-by-engineer",
	"/file-change-ratio-by-engineer-curr-year/title",
	"/file-change-ratio-by-engineer-curr-year",
	"/engineer-file-changes-over-time-curr-year/title",
	"/engineer-file-changes-over-time-curr-year",
	"/commonly-changed-files/title",
	"/commonly-changed-files",
	"/num-commits-prev-year/title",
	"/num-commits-prev-year",
	"/num-commits-curr-year/title",
	"/num-commits-curr-year",
	"/num-commits-all-time/title",
	"/num-commits-all-time",
	"/engineer-commits-over-time-curr-year/title",
	"/engineer-commits-over-time-curr-year",
	"/engineer-commit-counts-curr-year/title",
	"/engineer-commit-counts-curr-year",
	"/engineer-commit-counts-all-time/title",
	"/engineer-commit-counts-all-time",
	"/commits-by-month-curr-year/title",
	"/commits-by-month-curr-year",
	"/commits-by-weekday-curr-year/title",
	"/commits-by-weekday-curr-year",
	"/commits-by-hour-curr-year/title",
	"/commits-by-hour-curr-year",
	"/most-single-day-commits-by-engineer-curr-year/title",
	"/most-single-day-commits-by-engineer-curr-year",
	"/most-single-day-commits-by-engineer-curr-year-commit-list",
	"/most-insertions-in-single-commit-curr-year/title",
	"/most-insertions-in-single-commit-curr-year",
	"/most-deletions-in-single-commit-curr-year/title",
	"/most-deletions-in-single-commit-curr-year",
	"/largest-commit-message-curr-year/title",
	"/largest-commit-message-curr-year",
	"/shortest-commit-message-curr-year/title",
	"/shortest-commit-message-curr-year/5",
	"/shortest-commit-message-curr-year/4",
	"/shortest-commit-message-curr-year/3",
	"/shortest-commit-message-curr-year/2",
	"/shortest-commit-message-curr-year/1",
	"/commit-message-length-histogram-curr-year/title",
	"/commit-message-length-histogram-curr-year",
	"/direct-pushes-on-master-by-engineer-curr-year/title",
	"/direct-pushes-on-master-by-engineer-curr-year",
	"/merges-to-master-by-engineer-curr-year/title",
	"/merges-to-master-by-engineer-curr-year",
	"/most-merges-in-one-day-curr-year/title",
	"/most-merges-in-one-day-curr-year",
	"/avg-merges-per-day-to-master-curr-year/title",
	"/avg-merges-per-day-to-master-curr-year",
	"/end",
}

func GetTableOfContents(recap analyzer.Recap) []string {
	if recap.IncludesFileBlames {
		return TABLE_OF_CONTENTS
	} else {
		pagesRequiringFileBlames := []string{
			"/total-lines-of-code-in-repo-by-engineer/title",
			"/total-lines-of-code-in-repo-by-engineer",
			"/third-largest-file/title",
			"/third-largest-file",
			"/second-largest-file/title",
			"/second-largest-file",
			"/largest-file/title",
			"/largest-file",

			// TODO: re-write these next two slides so that they don't require file blames
			"/file-count-prev-year/title",
			"/file-count-prev-year",
			"/file-count-curr-year/title",
			"/file-count-curr-year",

			// TODO: re-write these just manually counting lines of code and not using
			// git blame
			"/total-lines-of-code-prev-year/title",
			"/total-lines-of-code-prev-year",
			"/total-lines-of-code-curr-year/title",
			"/total-lines-of-code-curr-year",
		}

		return utils.Filter(TABLE_OF_CONTENTS, func(s string) bool {
			return !slices.Contains(pagesRequiringFileBlames, s)
		})
	}
}

func GetSingleYearRepoTableOfContents(recap analyzer.Recap) []string {
	pagesRequiringMultipleYears := []string{
		"/file-count-prev-year/title",
		"/file-count-prev-year",
		"/total-lines-of-code-prev-year/title",
		"/total-lines-of-code-prev-year",
		"/num-commits-prev-year/title",
		"/num-commits-prev-year",
		"/new-engineer-count-curr-year/title",
		"/new-engineer-count-curr-year",
		"/engineer-commit-counts-curr-year/title",
		"/engineer-commit-counts-curr-year",
		"/engineer-count-curr-year/title",
		"/engineer-count-curr-year",
	}

	// Single year repos can't involve anything with the previous year. But a lot
	// of slides are still relevant, so we just filter the irrelevant pages
	return utils.Filter(GetTableOfContents(recap), func(s string) bool {
		return !slices.Contains(pagesRequiringMultipleYears, s)
	})
}

func GetNextButtonLink(currUrl string, recap analyzer.Recap) string {
	if recap.IsMultiYearRepo {
		tableOfContents := GetTableOfContents(recap)
		currPageIdx := utils.FindIndex(tableOfContents, func(page string) bool {
			return page == currUrl
		})
		final := tableOfContents[getNextIdx(currPageIdx, len(tableOfContents))]

		return final
	} else {
		tableOfContents := GetSingleYearRepoTableOfContents(recap)
		currPageIdx := utils.FindIndex(tableOfContents, func(page string) bool {
			return page == currUrl
		})
		return tableOfContents[getNextIdx(currPageIdx, len(tableOfContents))]
	}
}

func getNextIdx(idx int, length int) int {
	if idx+1 >= length {
		return length
	} else {
		return idx + 1
	}
}

type TitleSlideData struct {
	Title       string
	Description string
	NextBtnUrl  string
}

func GetTitleSlideData(page string, recap analyzer.Recap) TitleSlideData {
	nextBtnUrl := GetNextButtonLink(fmt.Sprintf("/%s/title", page), recap)
	data := TitleSlideData{
		Title:       "",
		Description: "",
		NextBtnUrl:  nextBtnUrl,
	}

	switch page {
	case "new-engineer-count-curr-year":
		data.Title = "New Engineer Count"
		data.Description = fmt.Sprintf("New engineers in %d who committed to %s.", analyzer.CURR_YEAR, recap.Name)
	case "engineer-count-curr-year":
		data.Title = "Total Engineer Count"
		data.Description = fmt.Sprintf("Total number of engineers who committed to %s in %d.", recap.Name, analyzer.CURR_YEAR)
	case "engineer-count-all-time":
		data.Title = "All Time Engineers"
		data.Description = fmt.Sprintf("Total number of engineers who committed to %s, since the beginning.", recap.Name)
	case "file-count-prev-year":
		data.Title = "Previous File Count"
		data.Description = fmt.Sprintf("Total number of files that existed in %s last year. (%d)", recap.Name, analyzer.PREV_YEAR)
	case "file-count-curr-year":
		data.Title = "Current File Count"
		data.Description = fmt.Sprintf("Total number of files that exist in %s this year. (%d)", recap.Name, analyzer.CURR_YEAR)
	default:
		panic(fmt.Sprintf("Unrecognized page for title slide: %s", page))
	}

	// "third-largest-file":                            "Third Largest File",
	// "second-largest-file":                           "Second Largest File",
	// "largest-file":                                  "Largest File",
	// "total-lines-of-code-prev-year":                 fmt.Sprintf("Total Lines of Code (%d)", analyzer.PREV_YEAR),
	// "total-lines-of-code-curr-year":                 fmt.Sprintf("Total Lines of Code (%d)", analyzer.CURR_YEAR),
	// "size-of-repo-by-week-curr-year":                fmt.Sprintf("Size of Repo by Week (%d)", analyzer.CURR_YEAR),
	// "total-lines-of-code-in-repo-by-engineer":       "Total Lines of Code by Engineer",
	// "file-changes-by-engineer-curr-year":            fmt.Sprintf("File Changes by Engineer (%d)", analyzer.CURR_YEAR),
	// "file-change-ratio-by-engineer-curr-year":       fmt.Sprintf("Insertions/Deletions Ratio by Engineer (%d)", analyzer.CURR_YEAR),
	// "commonly-changed-files":                        "Most Commonly Changed Files",
	// "num-commits-prev-year":                         fmt.Sprintf("Number of Commits (%d)", analyzer.PREV_YEAR),
	// "num-commits-curr-year":                         fmt.Sprintf("Number of Commits (%d)", analyzer.CURR_YEAR),
	// "num-commits-all-time":                          "Number of Commits (all time)",
	// "engineer-commits-over-time-curr-year":          fmt.Sprintf("Engineer Commits Over Time (%d)", analyzer.CURR_YEAR),
	// "engineer-file-changes-over-time-curr-year":     fmt.Sprintf("Engineer File Changes Over Time (%d)", analyzer.CURR_YEAR),
	// "engineer-commit-counts-curr-year":              fmt.Sprintf("Engineer Commit Counts (%d)", analyzer.CURR_YEAR),
	// "engineer-commit-counts-all-time":               "Engineer Commit Counts (all time)",
	// "commits-by-month-curr-year":                    fmt.Sprintf("Commits by Month (%d)", analyzer.CURR_YEAR),
	// "commits-by-weekday-curr-year":                  fmt.Sprintf("Commits by Weekday (%d)", analyzer.CURR_YEAR),
	// "commits-by-hour-curr-year":                     fmt.Sprintf("Commits by Hour (%d)", analyzer.CURR_YEAR),
	// "most-single-day-commits-by-engineer-curr-year": fmt.Sprintf("Most Single-Day Commits by Engineer (%d)", analyzer.CURR_YEAR),
	// "most-insertions-in-single-commit-curr-year":    fmt.Sprintf("Most Insertions in a Single Commit (%d)", analyzer.CURR_YEAR),
	// "most-deletions-in-single-commit-curr-year":     fmt.Sprintf("Most Deletions in a Single Commit (%d)", analyzer.CURR_YEAR),
	// "largest-commit-message-curr-year":              fmt.Sprintf("Largest Commit Message (%d)", analyzer.CURR_YEAR),
	// "shortest-commit-message-curr-year":             fmt.Sprintf("Shortest Commit Message (%d)", analyzer.CURR_YEAR),
	// "commit-message-length-histogram-curr-year":     fmt.Sprintf("Commit Message Length Frequencies (%d)", analyzer.CURR_YEAR),
	// "direct-pushes-on-master-by-engineer-curr-year": fmt.Sprintf("Direct Pushes on Master by Engineer (%d)", analyzer.CURR_YEAR),
	// "merges-to-master-by-engineer-curr-year":        fmt.Sprintf("Merges to Master by Engineer (%d)", analyzer.CURR_YEAR),
	// "most-merges-in-one-day-curr-year":              fmt.Sprintf("Most Merges in One Day (%d)", analyzer.CURR_YEAR),
	// "avg-merges-per-day-to-master-curr-year":        fmt.Sprintf("Average Merges per Day to Master (%d)", analyzer.CURR_YEAR),

	return data
}
