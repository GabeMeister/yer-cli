package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"time"
)

func GetNewAuthorCommitsCurrYear() []GitCommit {
	pastCommits := getPastGitCommits()

	// userName -> throwaway
	authorsFromPast := make(map[string]int)
	for _, commit := range pastCommits {
		authorsFromPast[commit.Author] = 1
	}

	currYearCommits := getCurrYearGitCommits()

	// username -> bool on whether they have been processed or not
	newAuthors := make(map[string]bool)
	for _, commit := range currYearCommits {
		if _, ok := authorsFromPast[commit.Author]; !ok {
			newAuthors[commit.Author] = false
		}
	}

	newAuthorCommits := []GitCommit{}

	for _, commit := range currYearCommits {
		userName := commit.Author
		processed, ok := newAuthors[userName]

		if ok && !processed {
			newAuthorCommits = append(newAuthorCommits, commit)
			newAuthors[userName] = true
		}
	}

	return newAuthorCommits
}

func GetAuthorCommitCountCurrYear() map[string]int {
	commits := getGitCommits()
	authors := make(map[string]int)

	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			userName := commit.Author
			authors[userName] += 1
		}
	}

	return authors
}

func GetAuthorCommitCountPrevYear() map[string]int {
	commits := getGitCommits()
	authors := make(map[string]int)

	for _, commit := range commits {
		commitYear := utils.GetYearFromDateStr(commit.Date)

		if commitYear < CURR_YEAR {
			userName := commit.Author
			authors[userName] += 1
		}
	}

	return authors
}

func GetAllUsernames() []string {
	authors := GetAuthorCommitCountAllTime()

	usernames := []string{}
	for username := range authors {
		usernames = append(usernames, username)
	}

	return usernames
}

func GetAuthorCommitCountAllTime() map[string]int {
	commits := getGitCommits()
	authors := make(map[string]int)

	for _, commit := range commits {
		userName := commit.Author
		authors[userName] += 1
	}

	return authors
}

func GetAuthorCountCurrYear() int {
	authors := GetAuthorCommitCountCurrYear()

	return len(authors)
}

func GetAuthorCountAllTime() int {
	authors := GetAuthorCommitCountAllTime()

	return len(authors)
}

// Example:
//
// [
//
//	{ Date: '2023-01-03T08:00:00.000Z', Name: 'Steve Bremer', Value: 24 },
//	{ Date: '2023-01-03T08:00:00.000Z', Name: 'Gabe Jensen', Value: 340 },
//	...
//
// ]

func GetAuthorCommitsOverTimeCurrYear() []TotalCommitCount {
	dates := utils.GetDaysOfYear(CURR_YEAR)

	// Create map of all possible dates this year
	dateMap := make(map[string][]GitCommit)
	for _, d := range dates {
		dateMap[d] = nil
	}

	commitTracker := make(map[string]int)

	// Bucket commit counts for all enginers in past
	pastCommits := getPastGitCommits()
	for _, commit := range pastCommits {
		commitTracker[commit.Author] += 1
	}

	// Get current year commits, and bucket them under whatever date they fall on
	currCommits := getCurrYearGitCommits()
	for _, commit := range currCommits {
		commitDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", commit.Date)
		if err != nil {
			panic("Invalid dates found in commits: " + commit.Date)
		}

		commitDateStr := commitDate.Format("2006-01-02")
		dateMap[commitDateStr] = append(dateMap[commitDateStr], commit)
	}

	final := []TotalCommitCount{}

	for _, dateStr := range dates {
		commitsOnDay := dateMap[dateStr]

		for _, commit := range commitsOnDay {
			commitTracker[commit.Author] += 1
		}

		for userName, numCommits := range commitTracker {
			final = append(final, TotalCommitCount{
				Name:  userName,
				Date:  dateStr,
				Value: numCommits,
			})
		}
	}

	return final
}

func GetAuthorFileChangesOverTimeCurrYear() []TotalFileChangeCount {
	if !HasPrevYearFileBlames() || !HasCurrYearFileBlames() {
		return []TotalFileChangeCount{}
	}

	// author => line change count
	fileChangeTracker := make(map[string]int)

	// Bucket file changes for all enginers in past
	prevFileBlames := GetPrevYearFileBlames()
	for _, fileBlame := range prevFileBlames {
		for author, lineCount := range fileBlame.GitBlame {
			fileChangeTracker[author] += lineCount
		}
	}

	dates := utils.GetDaysOfYear(CURR_YEAR)

	// Create map of all possible dates this year
	dateMap := make(map[string][]GitCommit)
	for _, d := range dates {
		dateMap[d] = nil
	}

	// Get current year commits, and bucket them under whatever date they fall on
	commits := getCurrYearGitCommits()
	for _, commit := range commits {
		commitDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", commit.Date)
		if err != nil {
			panic("Invalid dates found in commits: " + commit.Date)
		}

		commitDateStr := commitDate.Format("2006-01-02")
		dateMap[commitDateStr] = append(dateMap[commitDateStr], commit)
	}

	final := []TotalFileChangeCount{}

	for _, dateStr := range dates {
		commitsOnDay := dateMap[dateStr]

		for _, commit := range commitsOnDay {
			for _, fileChange := range commit.FileChanges {
				fileChangeTracker[commit.Author] += fileChange.Insertions
				fileChangeTracker[commit.Author] += fileChange.Deletions
			}
		}

		for userName, numFileChanges := range fileChangeTracker {
			final = append(final, TotalFileChangeCount{
				Name:  userName,
				Date:  dateStr,
				Value: numFileChanges,
			})

		}

	}

	return final
}

type CommitList []string
type AuthorCommitList map[string]CommitList
type DayCommitListByAuthor map[string]AuthorCommitList

func GetMostCommitsByAuthorCurrYear() MostSingleDayCommitsByAuthor {
	// Go from:

	// [
	//   {
	//     "commit": "eaeea5e7b50864bc2695f7bfa73b7106974f2165",
	//     "author": "Steve Bremer",
	//     "email": "steve@redballoon.work",
	//     "message": "Initial commit",
	//     "date": "Mon Aug 29 22:36:29 2022 +0000",
	//     "file_changes": [
	//       {
	//         "insertions": 92,
	//         "deletions": 0,
	//         "file_path": "README.md"
	//       }
	//     ]
	//   },
	//   .
	//   .
	//   .
	// ]

	// To:

	// {
	//   "2022-08-29": {
	//     "Steve": [
	//       "Initial commit",
	//     ],
	//     "Ezra": [
	//       "this is broken",
	//       "i fix",
	//     ],
	//   },
	//   .
	//   .
	//   .
	// }

	commits := getCurrYearGitCommits()
	commitListByDay := DayCommitListByAuthor{}

	for _, commit := range commits {
		dateStr := utils.GetMachineReadableDateStr(commit.Date)
		author := commit.Author
		msg := commit.Message

		if commitListByDay[dateStr] == nil {
			commitListByDay[dateStr] = make(AuthorCommitList)
		}

		if commitListByDay[dateStr][author] == nil {
			commitListByDay[dateStr][author] = CommitList{}
		}

		commitListByDay[dateStr][author] = append(commitListByDay[dateStr][author], msg)
	}

	mostCommitsAuthor := MostSingleDayCommitsByAuthor{
		Username: "",
		Date:     "2024-01-01",
		Count:    0,
		Commits:  []string{},
	}

	for _, date := range utils.GetDaysOfYear(CURR_YEAR) {
		currDay := commitListByDay[date]

		for author, commits := range currDay {
			if len(commits) > mostCommitsAuthor.Count {
				mostCommitsAuthor = MostSingleDayCommitsByAuthor{
					Username: author,
					Date:     date,
					Count:    len(commits),
					Commits:  commits,
				}
			}
		}
	}

	return mostCommitsAuthor
}

func GetTotalLinesOfCodeInRepoByAuthor() map[string]int {
	if !HasCurrYearFileBlames() {
		return make(map[string]int)
	}

	authorLineCountMap := make(map[string]int)

	fileBlames := GetCurrYearFileBlames()
	for _, fileBlame := range fileBlames {
		for author, lineCount := range fileBlame.GitBlame {
			authorLineCountMap[author] += lineCount
		}
	}

	return authorLineCountMap
}
