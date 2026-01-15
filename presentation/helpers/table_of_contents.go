package helpers

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"fmt"
	"slices"
)

var MULTI_REPO_TABLE_OF_CONTENTS = []string{
	"/",
	"/active-authors/title",
	"/active-authors",
	"/file-count-by-repo/title",
	"/file-count-by-repo",
	"/total-lines-of-code-by-repo/title",
	"/total-lines-of-code-by-repo",
	"/total-lines-of-code-per-week-by-repo/title",
	"/total-lines-of-code-per-week-by-repo",
	// "/commonly-changed-files/title",
	// "/commonly-changed-files",
	"/commits-made-by-repo/title",
	"/commits-made-by-repo",
	"/commits-made-by-author/title",
	"/commits-made-by-author",
	"/file-changes-made-by-author/title",
	"/file-changes-made-by-author",
	"/lines-of-code-owned-by-author-all-time/title",
	"/lines-of-code-owned-by-author-all-time",
	"/merge-commits-by-month/title",
	"/merge-commits-by-month",
	"/merge-commits-by-week-day/title",
	"/merge-commits-by-week-day",
	"/merge-commits-by-hour/title",
	"/merge-commits-by-hour",
	"/merges-to-master-by-repo/title",
	"/merges-to-master-by-repo",
	"/end",
}

var SINGLE_REPO_TABLE_OF_CONTENTS = []string{
	"/",
	"/new-author-count-curr-year/title",
	"/new-author-count-curr-year",
	"/new-author-list-curr-year",
	"/author-count-curr-year/title",
	"/author-count-curr-year",
	"/author-count-all-time/title",
	"/author-count-all-time",
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
	"/file-changes-by-author-curr-year/title",
	"/file-changes-by-author-curr-year",
	"/total-lines-of-code-in-repo-by-author/title",
	"/total-lines-of-code-in-repo-by-author",
	"/file-change-ratio-by-author-curr-year/title",
	"/file-change-ratio-by-author-curr-year",
	"/author-file-changes-over-time-curr-year/title",
	"/author-file-changes-over-time-curr-year",
	"/commonly-changed-files/title",
	"/commonly-changed-files",
	"/num-commits-prev-year/title",
	"/num-commits-prev-year",
	"/num-commits-curr-year/title",
	"/num-commits-curr-year",
	"/num-commits-all-time/title",
	"/num-commits-all-time",
	// "/author-commits-over-time-curr-year/title",
	// "/author-commits-over-time-curr-year",
	"/author-commit-counts-curr-year/title",
	"/author-commit-counts-curr-year",
	"/author-commit-counts-all-time/title",
	"/author-commit-counts-all-time",
	"/commits-by-month-curr-year/title",
	"/commits-by-month-curr-year",
	"/commits-by-weekday-curr-year/title",
	"/commits-by-weekday-curr-year",
	"/commits-by-hour-curr-year/title",
	"/commits-by-hour-curr-year",
	"/most-single-day-commits-by-author-curr-year/title",
	"/most-single-day-commits-by-author-curr-year",
	"/most-single-day-commits-by-author-curr-year-commit-list",
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
	"/direct-pushes-on-master-by-author-curr-year/title",
	"/direct-pushes-on-master-by-author-curr-year",
	"/merges-to-master-by-author-curr-year/title",
	"/merges-to-master-by-author-curr-year",
	"/most-merges-in-one-day-curr-year/title",
	"/most-merges-in-one-day-curr-year",
	"/avg-merges-per-day-to-master-curr-year/title",
	"/avg-merges-per-day-to-master-curr-year",
	"/end",
}

func GetMultiRepoTableOfContents(multiRepoRecap analyzer.MultiRepoRecap) []string {
	return MULTI_REPO_TABLE_OF_CONTENTS
}

func GetSingleRepoTableOfContents(recap analyzer.Recap) []string {
	if recap.IncludesFileBlames {
		return SINGLE_REPO_TABLE_OF_CONTENTS
	} else {
		pagesRequiringFileBlames := []string{
			"/total-lines-of-code-in-repo-by-author/title",
			"/total-lines-of-code-in-repo-by-author",
			"/third-largest-file/title",
			"/third-largest-file",
			"/second-largest-file/title",
			"/second-largest-file",
			"/largest-file/title",
			"/largest-file",

			// TODO: re-write these just manually counting lines of code and not using
			// git blame
			"/total-lines-of-code-prev-year/title",
			"/total-lines-of-code-prev-year",
			"/total-lines-of-code-curr-year/title",
			"/total-lines-of-code-curr-year",

			"/author-file-changes-over-time-curr-year/title",
			"/author-file-changes-over-time-curr-year",
		}

		return utils.Filter(SINGLE_REPO_TABLE_OF_CONTENTS, func(s string) bool {
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
		"/new-author-count-curr-year/title",
		"/new-author-count-curr-year",
		"/author-commit-counts-curr-year/title",
		"/author-commit-counts-curr-year",
		"/author-count-curr-year/title",
		"/author-count-curr-year",
	}

	// Single year repos can't involve anything with the previous year. But a lot
	// of slides are still relevant, so we just filter the irrelevant pages
	return utils.Filter(GetSingleRepoTableOfContents(recap), func(s string) bool {
		return !slices.Contains(pagesRequiringMultipleYears, s)
	})
}

func GetMultiRepoNextButtonLink(currUrl string, multiRepoRecap analyzer.MultiRepoRecap) string {
	tableOfContents := GetMultiRepoTableOfContents(multiRepoRecap)
	currPageIdx := utils.FindIndex(tableOfContents, func(page string) bool {
		return page == currUrl
	})
	return tableOfContents[getNextIdx(currPageIdx, len(tableOfContents))]

}

func GetNextButtonLink(currUrl string, recap analyzer.Recap) string {
	tableOfContents := GetSingleRepoTableOfContents(recap)
	currPageIdx := utils.FindIndex(tableOfContents, func(page string) bool {
		return page == currUrl
	})
	return tableOfContents[getNextIdx(currPageIdx, len(tableOfContents))]
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

func GetMultiRepoTitleSlideData(page string, recap analyzer.MultiRepoRecap) TitleSlideData {
	nextBtnUrl := GetMultiRepoNextButtonLink(fmt.Sprintf("/%s/title", page), recap)
	data := TitleSlideData{
		Title:       "",
		Description: "",
		NextBtnUrl:  nextBtnUrl,
	}

	switch page {
	case "active-authors":
		data.Title = "Active Authors"
		data.Description = "The amount of unique authors that contributed to your repos."
	case "file-count-by-repo":
		data.Title = "File Count by Repo"
		data.Description = "The unique file count of your repos."
	case "total-lines-of-code-by-repo":
		data.Title = "Total Lines of Code by Repo"
		data.Description = "The total lines of code in each of your repos."
	case "total-lines-of-code-per-week-by-repo":
		data.Title = "Total Lines of Code per Week"
		data.Description = "The total lines of code per week throughout the past year, combined across all repos."
	case "commits-made-by-repo":
		data.Title = "Commits Made by Repo"
		data.Description = "The total amount of commits made over the past year, split out by repo."
	case "commits-made-by-author":
		data.Title = "Commits Made by Author"
		data.Description = "The total amount of commits across all repos made over the past year, split out by author."
	case "file-changes-made-by-author":
		data.Title = "File Changes by Author"
		data.Description = "The total amount of insertions/deletions made over the past year, split out by author."
	case "lines-of-code-owned-by-author-all-time":
		data.Title = "Lines of Code by Author"
		data.Description = "The total lines of code that each author owns, across all repos."
	case "merge-commits-by-month":
		data.Title = "Merges by Month"
		data.Description = "The number of times code was merged across all repos, broken out by month."
	case "merge-commits-by-week-day":
		data.Title = "Merges by Week Day"
		data.Description = "The number of times code was merged across all repos, broken out by week day."
	case "merge-commits-by-hour":
		data.Title = "Merges by Hour"
		data.Description = "The number of times code was merged across all repos, broken out by hour of day."
	case "merges-to-master-by-repo":
		data.Title = "Merges to Master by Repo"
		data.Description = "The total amount of merges into the master branch, broken out by repo."

	default:
		panic(fmt.Sprintf("Unrecognized page for multi repo recap title slide: %s", page))
	}

	return data
}
