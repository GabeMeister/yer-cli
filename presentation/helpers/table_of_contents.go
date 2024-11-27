package presentation_helpers

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"slices"
)

var TABLE_OF_CONTENTS = []string{
	"/",
	"/new-engineer-count-curr-year/title",
	"/new-engineer-count-curr-year",
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

func GetTableOfContents() []string {
	return TABLE_OF_CONTENTS
}

func GetSingleYearRepoTableOfContents() []string {
	// Single year repos can't involve anything with the previous year. But a lot
	// of slides are still relevant, so we just filter the irrelevant pages
	return utils.Filter(TABLE_OF_CONTENTS, func(s string) bool {
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

		return !slices.Contains(pagesRequiringMultipleYears, s)
	})
}

func GetNextButtonLink(currUrl string, recap analyzer.Recap) string {
	currPageIdx := utils.FindIndex(TABLE_OF_CONTENTS, func(page string) bool {
		return page == currUrl
	})

	if recap.IsMultiYearRepo {
		tableOfContents := GetTableOfContents()
		return tableOfContents[getNextIdx(currPageIdx, len(tableOfContents))]
	} else {
		tableOfContents := GetSingleYearRepoTableOfContents()
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
