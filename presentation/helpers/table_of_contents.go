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

func GetTitleSlideData(page string, recap analyzer.Recap) TitleSlideData {
	nextBtnUrl := GetNextButtonLink(fmt.Sprintf("/%s/title", page), recap)
	data := TitleSlideData{
		Title:       "",
		Description: "",
		NextBtnUrl:  nextBtnUrl,
	}

	switch page {
	case "new-author-count-curr-year":
		data.Title = "New Author Count"
		data.Description = fmt.Sprintf("New authors in %d who committed to %s.", analyzer.CURR_YEAR, recap.Name)

	case "author-count-curr-year":
		data.Title = "Total Author Count"
		data.Description = fmt.Sprintf("Total number of authors who committed to %s in %d.", recap.Name, analyzer.CURR_YEAR)

	case "author-count-all-time":
		data.Title = "All Time Authors"
		data.Description = fmt.Sprintf("Total number of authors who committed to %s, since the beginning.", recap.Name)

	case "file-count-prev-year":
		data.Title = "Previous File Count"
		data.Description = fmt.Sprintf("Total number of files that existed in %s last year (%d).", recap.Name, analyzer.PREV_YEAR)

	case "file-count-curr-year":
		data.Title = "Current File Count"
		data.Description = fmt.Sprintf("Total number of files that exist in %s this year (%d).", recap.Name, analyzer.CURR_YEAR)

	case "third-largest-file":
		data.Title = "Third Largest File"
		data.Description = "The third largest file in the repo right now."

	case "second-largest-file":
		data.Title = "Second Largest File"
		data.Description = "The second largest file in the repo right now."

	case "largest-file":
		data.Title = "Largest File"
		data.Description = "The absolute largest file in the entire repo right now."

	case "total-lines-of-code-prev-year":
		data.Title = fmt.Sprintf("Total Lines of Code (%d)", analyzer.PREV_YEAR)
		data.Description = fmt.Sprintf("Total lines of code in the entire repo as of the end of last year (%d).", analyzer.PREV_YEAR)

	case "total-lines-of-code-curr-year":
		data.Title = fmt.Sprintf("Total Lines of Code (%d)", analyzer.CURR_YEAR)
		data.Description = fmt.Sprintf("Total lines of code in the entire repo as of this year (%d).", analyzer.CURR_YEAR)

	case "size-of-repo-by-week-curr-year":
		data.Title = "Weekly Repo Size"
		data.Description = fmt.Sprintf("Size of Repo by Week (%d)", analyzer.CURR_YEAR)

	case "total-lines-of-code-in-repo-by-author":
		data.Title = "Total Lines of Code"
		data.Description = fmt.Sprintf("The total number of lines of code in %s, categorized by author.", recap.Name)

	case "file-changes-by-author-curr-year":
		data.Title = fmt.Sprintf("Line Changes (%d)", analyzer.CURR_YEAR)
		data.Description = fmt.Sprintf("The total number of line changes made in %d by author.", analyzer.CURR_YEAR)

	case "file-change-ratio-by-author-curr-year":
		data.Title = "Line Change Ratios"
		data.Description = fmt.Sprintf("The ratio of line insertions to deletions by author. A higher number means an author adds in more code to the repo than removes it. (%d)", analyzer.CURR_YEAR)

	case "commonly-changed-files":
		data.Title = "Commonly Changed Files"
		data.Description = fmt.Sprintf("The files that seem to be changed the most frequently throughout %d.", analyzer.CURR_YEAR)

	case "num-commits-prev-year":
		data.Title = fmt.Sprintf("Number of Commits (%d)", analyzer.PREV_YEAR)
		data.Description = fmt.Sprintf("The total number of commits made by authors last year (%d).", analyzer.PREV_YEAR)

	case "num-commits-curr-year":
		data.Title = fmt.Sprintf("Number of Commits (%d)", analyzer.CURR_YEAR)
		data.Description = fmt.Sprintf("The total number of commits made by authors this year (%d).", analyzer.CURR_YEAR)

	case "num-commits-all-time":
		data.Title = "Number of Commits (All Time)"
		data.Description = "The total number of commits made by authors, since the very beginning."

	case "author-commits-over-time-curr-year":
		data.Title = fmt.Sprintf("Commits Over Time (%d)", analyzer.CURR_YEAR)
		data.Description = fmt.Sprintf("The number of commits made by each author, throughout the duration of %d.", analyzer.CURR_YEAR)

	case "author-file-changes-over-time-curr-year":
		data.Title = fmt.Sprintf("Line Changes Over Time (%d)", analyzer.CURR_YEAR)
		data.Description = fmt.Sprintf("The number of line changes made by author, throughout the duration of %d.", analyzer.CURR_YEAR)

	case "author-commit-counts-curr-year":
		data.Title = fmt.Sprintf("Commit Counts (%d)", analyzer.CURR_YEAR)
		data.Description = fmt.Sprintf("The number of commits by each author in %d.", analyzer.CURR_YEAR)

	case "author-commit-counts-all-time":
		data.Title = "Commit Counts (All Time)"
		data.Description = "The number of commits by each author, since the beginning."

	case "commits-by-weekday-curr-year":
		data.Title = "Commits by Weekday"
		data.Description = fmt.Sprintf("Number of commits made each week day, throughout %d.", analyzer.CURR_YEAR)

	case "commits-by-hour-curr-year":
		data.Title = "Commits by Hour"
		data.Description = fmt.Sprintf("Number of commits made each hour of the day, throughout %d.", analyzer.CURR_YEAR)

	case "commits-by-month-curr-year":
		data.Title = "Commits by Month"
		data.Description = fmt.Sprintf("Number of commits made each month of the year, throughout %d.", analyzer.CURR_YEAR)

	case "most-single-day-commits-by-author-curr-year":
		data.Title = "Most Single-Day Commits by Author"
		data.Description = fmt.Sprintf("The most commits made in one day by an author in %d.", analyzer.CURR_YEAR)

	case "most-insertions-in-single-commit-curr-year":
		data.Title = "Most Code Added in Single Commit"
		data.Description = fmt.Sprintf("The most massive single code change in %d.", analyzer.CURR_YEAR)

	case "most-deletions-in-single-commit-curr-year":
		data.Title = "Most Code Removed in Single Commit"
		data.Description = fmt.Sprintf("The most code nuked from the codebase in a single commit in %d.", analyzer.CURR_YEAR)

	case "largest-commit-message-curr-year":
		data.Title = "Largest Commit Message"
		data.Description = fmt.Sprintf("Largest commit message written by an author in %d.", analyzer.CURR_YEAR)

	case "shortest-commit-message-curr-year":
		data.Title = "Shortest Commit Message"
		data.Description = fmt.Sprintf("The shortest, low-effort commit messages written in %d.", analyzer.CURR_YEAR)

	case "commit-message-length-histogram-curr-year":
		data.Title = "Commit Message Lengths"
		data.Description = fmt.Sprintf("A histogram tracking the frequency of git commit message lengths in %d.", analyzer.CURR_YEAR)

	case "direct-pushes-on-master-by-author-curr-year":
		data.Title = "Direct Pushes on Master"
		data.Description = fmt.Sprintf("The number of direct pushes to master, by author, in %d.", analyzer.CURR_YEAR)

	case "merges-to-master-by-author-curr-year":
		data.Title = "Testing on Merge Requests"
		data.Description = fmt.Sprintf("The authors who tested and merged code the most in %d.", analyzer.CURR_YEAR)

	case "most-merges-in-one-day-curr-year":
		data.Title = "Most Merges in One Day"
		data.Description = fmt.Sprintf("The most amount of merges done in a single day in %d.", analyzer.CURR_YEAR)

	case "avg-merges-per-day-to-master-curr-year":
		data.Title = "Average Merges Per Day"
		data.Description = fmt.Sprintf("The average number of merges the team did per day in %d.", analyzer.CURR_YEAR)

	default:
		panic(fmt.Sprintf("Unrecognized page for title slide: %s", page))
	}

	return data
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

	default:
		panic(fmt.Sprintf("Unrecognized page for multi repo recap title slide: %s", page))
	}

	return data
}
