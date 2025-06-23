package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

/*
 * COMMIT GETTERS
 */

func (r *RepoConfig) getGitCommits() []GitCommit {
	bytes, err := os.ReadFile(r.GetCommitsFile())
	if err != nil {
		panic(err)
	}

	var commits []GitCommit
	jsonErr := json.Unmarshal(bytes, &commits)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return commits
}

func (r *RepoConfig) getMergeGitCommits() []GitCommit {
	bytes, err := os.ReadFile(r.GetMergeCommitsFile())
	if err != nil {
		panic(err)
	}

	var commits []GitCommit
	jsonErr := json.Unmarshal(bytes, &commits)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return commits
}

func (r *RepoConfig) getDirectPushOnMasterCommits() []GitCommit {
	bytes, err := os.ReadFile(r.GetDirectPushesFile())
	if err != nil {
		panic(err)
	}

	var commits []GitCommit
	jsonErr := json.Unmarshal(bytes, &commits)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return commits
}

func (r *RepoConfig) getPastGitCommits() []GitCommit {
	// We return any commit that has come before the current year
	commits := r.getGitCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrBeforeYear(commit.Date, CURR_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func (r *RepoConfig) getPrevYearGitCommits() []GitCommit {
	commits := r.getGitCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, PREV_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func (r *RepoConfig) getCurrYearGitCommits() []GitCommit {
	commits := r.getGitCommits()
	final := []GitCommit{}

	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func (r *RepoConfig) getCurrYearMergeGitCommits() []GitCommit {
	commits := r.getMergeGitCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func (r *RepoConfig) getCurrYearDirectPushOnMasterCommits() []GitCommit {
	commits := r.getDirectPushOnMasterCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

/*
 * COMMIT STATS
 */

func (r *RepoConfig) GetIsMultiYearRepo() bool {
	commits := r.getGitCommits()
	firstCommit := commits[0]

	return utils.GetYearFromDateStr(firstCommit.Date) < CURR_YEAR
}

func (r *RepoConfig) GetNumCommitsAllTime() int {
	commits := r.getGitCommits()

	return len(commits)
}

func (r *RepoConfig) GetNumCommitsPrevYear() int {
	commits := r.getGitCommits()
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() == PREV_YEAR
	})

	return len(prevYearCommits)
}

func (r *RepoConfig) GetNumCommitsCurrYear() int {
	commits := r.getGitCommits()
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() == CURR_YEAR
	})

	return len(prevYearCommits)
}

func (r *RepoConfig) GetNumCommitsInPast() int {
	commits := r.getGitCommits()
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() < PREV_YEAR
	})

	return len(prevYearCommits)
}

func (r *RepoConfig) GetCommitsByMonthCurrYear() []CommitMonth {
	commits := r.getCurrYearGitCommits()

	monthMap := make(map[string]int)

	for _, commit := range commits {
		currDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", commit.Date)
		if err != nil {
			panic(err)
		}
		month := currDate.Month().String()
		monthMap[month] += 1
	}

	commitMonths := []CommitMonth{}

	for _, month := range MONTHS {
		commitMonths = append(commitMonths, CommitMonth{
			Month:   month,
			Commits: monthMap[month],
		})
	}

	return commitMonths
}

func (r *RepoConfig) GetCommitsByWeekDayCurrYear() []CommitWeekDay {
	commits := r.getCurrYearGitCommits()

	weekDayMap := make(map[string]int)

	for _, commit := range commits {
		currDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", commit.Date)
		if err != nil {
			panic(err)
		}
		weekDay := currDate.Weekday().String()
		weekDayMap[weekDay] += 1
	}

	commitWeekDays := []CommitWeekDay{}

	for _, weekDay := range WEEK_DAYS {
		commitWeekDays = append(commitWeekDays, CommitWeekDay{
			Day:     weekDay,
			Commits: weekDayMap[weekDay],
		})
	}

	return commitWeekDays
}

func (r *RepoConfig) GetCommitsByHourCurrYear() []CommitHour {
	commits := r.getCurrYearGitCommits()

	hourMap := make(map[int]int)

	for _, commit := range commits {
		currDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", commit.Date)
		if err != nil {
			panic(err)
		}
		hour := currDate.Hour()
		hourMap[hour] += 1
	}

	commitHours := []CommitHour{}

	for idx, hour := range HOURS {
		commitHours = append(commitHours, CommitHour{
			Hour:    hour,
			Commits: hourMap[idx],
		})
	}

	return commitHours
}

func (r *RepoConfig) GetMostInsertionsInCommitCurrYear() GitCommit {
	commits := r.getCurrYearGitCommits()

	mostInsertionsCommit := commits[0]
	mostInsertionsAmt := 0

	for _, commit := range commits {
		totalChanges := 0
		for _, fileChange := range commit.FileChanges {
			totalChanges += fileChange.Insertions
		}

		if totalChanges > mostInsertionsAmt {
			mostInsertionsCommit = commit
			mostInsertionsAmt = totalChanges
		}
	}

	return mostInsertionsCommit
}

func (r *RepoConfig) GetMostDeletionsInCommitCurrYear() GitCommit {
	commits := r.getCurrYearGitCommits()

	mostDeletionsCommit := commits[0]
	largestDeletionsAmt := 0

	for _, commit := range commits {
		totalChanges := 0
		for _, fileChange := range commit.FileChanges {
			totalChanges += fileChange.Deletions
		}

		if totalChanges > largestDeletionsAmt {
			mostDeletionsCommit = commit
			largestDeletionsAmt = totalChanges
		}
	}

	return mostDeletionsCommit
}

func (r *RepoConfig) GetLargestCommitMessageCurrYear() GitCommit {
	commits := r.getCurrYearGitCommits()

	largestLengthCommit := commits[0]

	for _, commit := range commits {
		if len(largestLengthCommit.Message) < len(commit.Message) {
			largestLengthCommit = commit
		}
	}

	return largestLengthCommit
}

func (r *RepoConfig) GetSmallestCommitMessagesCurrYear() []GitCommit {
	commits := r.getCurrYearGitCommits()

	sort.Slice(commits, func(i int, j int) bool {
		return len(commits[i].Message) < len(commits[j].Message)
	})

	return commits[0:5]
}

func (r *RepoConfig) GetCommitMessageHistogramCurrYear() []CommitMessageLengthFrequency {
	commits := r.getCurrYearGitCommits()

	lengthFrequencyMap := make(map[int]int)

	for _, commit := range commits {
		// "Convert" the `|||` back to newlines
		msg := strings.ReplaceAll(commit.Message, "|||", "\n")
		length := len(msg)
		lengthFrequencyMap[length] += 1
	}

	maxLength := 0
	for length := range lengthFrequencyMap {
		if length > maxLength {
			maxLength = length
		}
	}

	// Index of the array represents the length, and each value is the frequency
	messageLengthFrequencies := []CommitMessageLengthFrequency{}

	for length, frequency := range lengthFrequencyMap {
		messageLengthFrequencies = append(messageLengthFrequencies, CommitMessageLengthFrequency{length, frequency})
	}

	// Sort by length
	sort.Slice(messageLengthFrequencies, func(i, j int) bool {
		return messageLengthFrequencies[i][0] < messageLengthFrequencies[j][0]
	})

	return messageLengthFrequencies
}

func (r *RepoConfig) GetDirectPushesOnMasterByAuthorCurrYear() map[string]int {
	commits := r.getCurrYearDirectPushOnMasterCommits()

	authorToCommitMap := make(map[string]int)

	for _, commit := range commits {
		authorToCommitMap[commit.Author] += 1
	}

	return authorToCommitMap
}

func (r *RepoConfig) GetMergesToMasterByAuthorCurrYear() map[string]int {
	commits := r.getCurrYearMergeGitCommits()

	authorToCommitMap := make(map[string]int)

	for _, commit := range commits {
		authorToCommitMap[commit.Author] += 1
	}

	return authorToCommitMap
}

func (r *RepoConfig) GetMostMergesInOneDayCurrYear() MostMergesInOneDay {
	commits := r.getCurrYearMergeGitCommits()

	dayCommitMap := make(map[string][]GitCommit)

	for _, commit := range commits {
		day := utils.GetSimpleDateStr(commit.Date)
		dayCommitMap[day] = append(dayCommitMap[day], commit)
	}

	mostMergesInOneDay := MostMergesInOneDay{
		Count: 0,
	}

	for day, commits := range dayCommitMap {
		if mostMergesInOneDay.Count < len(commits) {
			mostMergesInOneDay = MostMergesInOneDay{
				Count:   len(commits),
				Date:    day,
				Commits: commits,
			}
		}
	}

	return mostMergesInOneDay
}

func (r *RepoConfig) GetAvgMergesToMasterPerDayCurrYear() float64 {
	commits := r.getCurrYearMergeGitCommits()
	if len(commits) == 0 {
		return 0.0
	}

	numWorkDays := GetNumWorkDaysInCurrYear()
	final := float64(len(commits)) / float64(numWorkDays)

	if math.IsNaN(final) {
		panic("GetAvgMergesToMasterPerDayCurrYear() is NaN!")
	}

	return final
}

func (r *RepoConfig) GetFileChangesByAuthorCurrYear() map[string]int {
	commits := r.getCurrYearGitCommits()

	authorInsertionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorInsertionsMap[commit.Author] += change.Insertions
			authorInsertionsMap[commit.Author] += change.Deletions
		}
	}

	return authorInsertionsMap
}

func (r *RepoConfig) GetCodeInsertionsByAuthorCurrYear() map[string]int {
	commits := r.getCurrYearGitCommits()

	authorDeletionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorDeletionsMap[commit.Author] += change.Insertions
		}
	}

	return authorDeletionsMap
}

func (r *RepoConfig) GetCodeDeletionsByAuthorCurrYear() map[string]int {
	commits := r.getCurrYearGitCommits()

	authorDeletionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorDeletionsMap[commit.Author] += change.Deletions
		}
	}

	return authorDeletionsMap
}

/*
 * TEAM STATS
 */

func (r *RepoConfig) GetNewAuthorCommitsCurrYear() []GitCommit {
	pastCommits := r.getPastGitCommits()

	// userName -> throwaway
	authorsFromPast := make(map[string]int)
	for _, commit := range pastCommits {
		authorsFromPast[commit.Author] = 1
	}

	currYearCommits := r.getCurrYearGitCommits()

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

func (r *RepoConfig) GetAuthorCommitCountCurrYear() map[string]int {
	commits := r.getGitCommits()
	authors := make(map[string]int)

	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			userName := commit.Author
			authors[userName] += 1
		}
	}

	return authors
}

func (r *RepoConfig) GetAuthorCommitCountPrevYear() map[string]int {
	commits := r.getGitCommits()
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

func (r *RepoConfig) GetAllUsernames() []string {
	authors := r.GetAuthorCommitCountAllTime()

	usernames := []string{}
	for username := range authors {
		usernames = append(usernames, username)
	}

	return usernames
}

func (r *RepoConfig) GetAuthorCommitCountAllTime() map[string]int {
	commits := r.getGitCommits()
	authors := make(map[string]int)

	for _, commit := range commits {
		userName := commit.Author
		authors[userName] += 1
	}

	return authors
}

func (r *RepoConfig) GetAuthorCountCurrYear() int {
	authors := r.GetAuthorCommitCountCurrYear()

	return len(authors)
}

func (r *RepoConfig) GetAuthorCountAllTime() int {
	authors := r.GetAuthorCommitCountAllTime()

	return len(authors)
}

func (r *RepoConfig) GetAuthorCommitsOverTimeCurrYear() []TotalCommitCount {
	// Example:
	//
	// [
	//
	//	{ Date: '2023-01-03T08:00:00.000Z', Name: 'Steve Bremer', Value: 24 },
	//	{ Date: '2023-01-03T08:00:00.000Z', Name: 'Gabe Jensen', Value: 340 },
	//	...
	//
	// ]
	dates := utils.GetDaysOfYear(CURR_YEAR)

	// Create map of all possible dates this year
	dateMap := make(map[string][]GitCommit)
	for _, d := range dates {
		dateMap[d] = nil
	}

	commitTracker := make(map[string]int)

	// Bucket commit counts for all enginers in past
	pastCommits := r.getPastGitCommits()
	for _, commit := range pastCommits {
		commitTracker[commit.Author] += 1
	}

	// Get current year commits, and bucket them under whatever date they fall on
	currCommits := r.getCurrYearGitCommits()
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

func (r *RepoConfig) GetAllAuthorsList() []string {
	allAuthorsMap := make(map[string]bool)

	if r.hasPrevYearCommits() {
		for _, commit := range r.getPrevYearGitCommits() {
			allAuthorsMap[commit.Author] = true
		}
	}

	for _, commit := range r.getCurrYearGitCommits() {
		allAuthorsMap[commit.Author] = true
	}

	return utils.MapKeysToSlice(allAuthorsMap)
}

// Used mainly as initial starting data for author file changes over time
func (r *RepoConfig) GetAuthorTotalFileChangesPrevYear() map[string]int {
	if !r.HasPrevYearFileBlames() || !r.HasCurrYearFileBlames() {
		return make(map[string]int)
	}

	// author => line change count
	fileChangeTracker := make(map[string]int)

	// Bucket file changes for all enginers in past
	prevFileBlames := r.GetPrevYearFileBlames()
	for _, fileBlame := range prevFileBlames {
		for author, lineCount := range fileBlame.GitBlame {
			fileChangeTracker[author] += lineCount
		}
	}

	return fileChangeTracker
}

func (r *RepoConfig) GetAuthorFileChangesOverTimeCurrYear() TotalFileChangeCount {
	if !r.HasPrevYearFileBlames() || !r.HasCurrYearFileBlames() {
		return make(TotalFileChangeCount)
	}

	// author => line change count
	fileChangeTracker := r.GetAuthorTotalFileChangesPrevYear()

	// e.g. [ '2025-01-01', '2025-01-02', ... ]
	dates := utils.GetDaysOfYear(CURR_YEAR)

	// Create map of all possible dates this year
	// date -> slice of GitCommit
	// 	e.g. 2025-01-02 -> [{ Commit: "29jf382gh892j38", ... }, { ... }, ...]
	dateMap := make(map[string][]GitCommit)
	for _, d := range dates {
		dateMap[d] = nil
	}

	// Get current year commits, and bucket them under whatever date they fall on
	commits := r.getCurrYearGitCommits()
	for _, commit := range commits {
		commitDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", commit.Date)
		if err != nil {
			panic("Invalid dates found in commits: " + commit.Date)
		}

		commitDateStr := commitDate.Format("2006-01-02")
		dateMap[commitDateStr] = append(dateMap[commitDateStr], commit)
	}

	// {
	//   "2024-01-01|Kenny": 29838,
	//   "2024-01-02|Isaac": 29838,
	//	 ...
	// }
	final := make(TotalFileChangeCount)

	for _, dateStr := range dates {
		commitsOnDay := dateMap[dateStr]

		for _, commit := range commitsOnDay {
			for _, fileChange := range commit.FileChanges {
				fileChangeTracker[commit.Author] += fileChange.Insertions
				fileChangeTracker[commit.Author] += fileChange.Deletions
			}
		}

		uniqAuthorMap := make(map[string]bool)
		for _, commit := range commitsOnDay {
			uniqAuthorMap[commit.Author] = true
		}
		uniqueAuthors := []string{}
		for author := range uniqAuthorMap {
			uniqueAuthors = append(uniqueAuthors, author)
		}

		// Add entries for ONLY the authors that actually committed on this day
		for _, author := range uniqueAuthors {
			key := fmt.Sprintf("%s|%s", dateStr, author)
			final[key] = fileChangeTracker[author]
		}
	}

	return final
}

func (r *RepoConfig) GetMostCommitsByAuthorCurrYear() MostSingleDayCommitsByAuthor {
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

	commits := r.getCurrYearGitCommits()
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

func (r *RepoConfig) GetTotalLinesOfCodeInRepoByAuthor() map[string]int {
	if !r.HasCurrYearFileBlames() {
		return make(map[string]int)
	}

	authorLineCountMap := make(map[string]int)

	fileBlames := r.GetCurrYearFileBlames()
	for _, fileBlame := range fileBlames {
		for author, lineCount := range fileBlame.GitBlame {
			authorLineCountMap[author] += lineCount
		}
	}

	return authorLineCountMap
}

/*
 * FILE STATS
 */

func (r *RepoConfig) GetPrevYearFileList() []string {
	bytes, err := os.ReadFile(r.GetPrevYearFileListFile())
	if err != nil {
		panic(err)
	}

	var files []string
	jsonErr := json.Unmarshal(bytes, &files)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return files

}

func (r *RepoConfig) GetCurrYearFileList() []string {
	bytes, err := os.ReadFile(r.GetCurrYearFileListFile())
	if err != nil {
		panic(err)
	}

	var files []string
	jsonErr := json.Unmarshal(bytes, &files)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return files

}

func (r *RepoConfig) HasPrevYearFileBlames() bool {
	_, err := os.ReadFile(r.GetPrevYearFileBlamesFile())

	return err == nil
}

func (r *RepoConfig) HasCurrYearFileBlames() bool {
	_, err := os.ReadFile(r.GetCurrYearFileBlamesFile())

	return err == nil
}

func (r *RepoConfig) GetPrevYearFileBlames() []FileBlame {
	bytes, err := os.ReadFile(r.GetPrevYearFileBlamesFile())
	if err != nil {
		panic(err)
	}

	var fileBlames []FileBlame
	jsonErr := json.Unmarshal(bytes, &fileBlames)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return fileBlames
}

func (r *RepoConfig) GetCurrYearFileBlames() []FileBlame {
	bytes, err := os.ReadFile(r.GetCurrYearFileBlamesFile())
	if err != nil {
		panic(err)
	}

	var fileBlames []FileBlame
	jsonErr := json.Unmarshal(bytes, &fileBlames)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return fileBlames
}

func (r *RepoConfig) GetFileChangeRatio(insertionsByAuthor map[string]int, deletionsByAuthor map[string]int) map[string]float64 {
	ratios := make(map[string]float64)

	for author, insertionCount := range insertionsByAuthor {
		insertions := float64(insertionCount)
		deletions := float64(deletionsByAuthor[author])

		if deletions == 0 {
			ratios[author] = 1
		} else {
			ratios[author] = insertions / deletions
		}
	}

	for _, ratio := range ratios {
		if math.IsNaN(ratio) {
			fmt.Print("\n\n", "*** ERROR in GetFileChangeRatio(). Ratio is NaN: ***", "\n", ratios, "\n\n\n")
			panic("Ratio is NaN in GetFileChangeRatio()!")
		}
	}

	return ratios
}

func (r *RepoConfig) GetCommonlyChangedFiles() []FileChangeCount {
	commits := r.getCurrYearGitCommits()
	fileChangeTracker := make(map[string]int)

	for _, commit := range commits {
		for _, changes := range commit.FileChanges {
			fileChangeTracker[changes.FilePath] += 1
		}
	}

	fileChangeSlice := []FileChangeCount{}
	for file, count := range fileChangeTracker {
		fileChangeSlice = append(fileChangeSlice, FileChangeCount{
			File:  file,
			Count: count,
		})
	}

	sort.Slice(fileChangeSlice, func(i int, j int) bool {
		return fileChangeSlice[i].Count > fileChangeSlice[j].Count
	})

	return fileChangeSlice[0:5]
}

func (r *RepoConfig) GetFileCountPrevYear() int {
	files := r.GetPrevYearFileList()

	return len(files)
}

func (r *RepoConfig) GetFileCountCurrYear() int {
	files := r.GetCurrYearFileList()

	return len(files)
}

func (r *RepoConfig) GetLargestFilesCurrYear() []FileSize {
	if !r.HasCurrYearFileBlames() {
		return []FileSize{}
	}

	fileBlames := r.GetCurrYearFileBlames()

	sort.Slice(fileBlames, func(i int, j int) bool {
		return fileBlames[i].LineCount > fileBlames[j].LineCount
	})

	fileRangeMax := math.Min(float64(len(fileBlames)), 10)

	fileSizes := []FileSize{}
	for _, fileBlame := range fileBlames[0:int(fileRangeMax)] {
		fileSizes = append(fileSizes, FileSize{
			File:      fileBlame.File,
			LineCount: fileBlame.LineCount,
		})
	}

	return fileSizes
}

func (r *RepoConfig) GetSmallestFilesCurrYear() []FileSize {
	if !r.HasCurrYearFileBlames() {
		return []FileSize{}
	}

	fileBlames := r.GetCurrYearFileBlames()

	sort.Slice(fileBlames, func(i int, j int) bool {
		return fileBlames[i].LineCount < fileBlames[j].LineCount
	})

	fileRangeMax := math.Min(float64(len(fileBlames)), 10)

	fileSizes := []FileSize{}
	for _, fileBlame := range fileBlames[0:int(fileRangeMax)] {
		fileSizes = append(fileSizes, FileSize{
			File:      fileBlame.File,
			LineCount: fileBlame.LineCount,
		})
	}

	return fileSizes
}

func (r *RepoConfig) GetTotalLinesOfCodePrevYear() int {
	if !r.HasPrevYearFileBlames() {
		return 0
	}

	fileBlames := r.GetPrevYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func (r *RepoConfig) GetTotalLinesOfCodeCurrYear() int {
	if !r.HasCurrYearFileBlames() {
		return 0
	}

	fileBlames := r.GetCurrYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func (r *RepoConfig) GetSizeOfRepoByWeekCurrYear() []RepoSizeTimeStamp {
	commits := r.getCurrYearGitCommits()
	weeks := utils.GetWeeksOfYear()

	weekCommitsMap := make(map[int][]GitCommit)
	for _, week := range weeks {
		weekCommitsMap[week] = nil
	}

	for _, commit := range commits {
		date := utils.GetDateFromISOString(commit.Date)
		_, week := date.ISOWeek()

		weekCommitsMap[week] = append(weekCommitsMap[week], commit)
	}

	totalLinesOfCodePrevYear := r.GetTotalLinesOfCodePrevYear()
	runningTotal := totalLinesOfCodePrevYear
	final := []RepoSizeTimeStamp{}

	for week := 1; week <= 52; week++ {
		weekCommits := weekCommitsMap[week]

		for _, commit := range weekCommits {
			for _, fileChanges := range commit.FileChanges {
				runningTotal += fileChanges.Insertions
				runningTotal -= fileChanges.Deletions
			}

		}

		final = append(final, RepoSizeTimeStamp{
			WeekNumber: week,
			LineCount:  runningTotal,
		})
	}

	return final
}
