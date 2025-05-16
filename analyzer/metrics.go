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
 * COMMITS
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
 * COMMIT METRICS
 */

func (r *RepoConfig) GetIsMultiYearRepo() bool {
	commits := r.getGitCommits()
	firstCommit := commits[0]

	return utils.GetYearFromDateStr(firstCommit.Date) < CURR_YEAR
}

func GetNumCommitsAllTime(config RepoConfig) int {
	commits := getGitCommits(config)

	return len(commits)
}

func GetNumCommitsPrevYear(config RepoConfig) int {
	commits := getGitCommits(config)
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() == PREV_YEAR
	})

	return len(prevYearCommits)
}

func GetNumCommitsCurrYear(config RepoConfig) int {
	commits := getGitCommits(config)
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() == CURR_YEAR
	})

	return len(prevYearCommits)
}

func GetNumCommitsInPast(config RepoConfig) int {
	commits := getGitCommits(config)
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() < PREV_YEAR
	})

	return len(prevYearCommits)
}

func GetCommitsByMonthCurrYear(config RepoConfig) []CommitMonth {
	commits := getCurrYearGitCommits(config)

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

func GetCommitsByWeekDayCurrYear(config RepoConfig) []CommitWeekDay {
	commits := getCurrYearGitCommits(config)

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

func GetCommitsByHourCurrYear(config RepoConfig) []CommitHour {
	commits := getCurrYearGitCommits(config)

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

func GetMostInsertionsInCommitCurrYear(config RepoConfig) GitCommit {
	commits := getCurrYearGitCommits(config)

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

func GetMostDeletionsInCommitCurrYear(config RepoConfig) GitCommit {
	commits := getCurrYearGitCommits(config)

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

func GetLargestCommitMessageCurrYear(config RepoConfig) GitCommit {
	commits := getCurrYearGitCommits(config)

	largestLengthCommit := commits[0]

	for _, commit := range commits {
		if len(largestLengthCommit.Message) < len(commit.Message) {
			largestLengthCommit = commit
		}
	}

	return largestLengthCommit
}

func GetSmallestCommitMessagesCurrYear(config RepoConfig) []GitCommit {
	commits := getCurrYearGitCommits(config)

	sort.Slice(commits, func(i int, j int) bool {
		return len(commits[i].Message) < len(commits[j].Message)
	})

	return commits[0:5]
}

func GetCommitMessageHistogramCurrYear(config RepoConfig) []CommitMessageLengthFrequency {
	commits := getCurrYearGitCommits(config)

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

	commitMessageLengths := make([]CommitMessageLengthFrequency, maxLength+1)
	for i := range maxLength + 1 {
		commitMessageLengths[i].Length = i
	}

	for length, frequency := range lengthFrequencyMap {
		commitMessageLengths[length].Length = length
		commitMessageLengths[length].Frequency = frequency
	}

	return commitMessageLengths
}

func GetDirectPushesOnMasterByAuthorCurrYear() map[string]int {
	commits := getCurrYearDirectPushOnMasterCommits()

	authorToCommitMap := make(map[string]int)

	for _, commit := range commits {
		authorToCommitMap[commit.Author] += 1
	}

	return authorToCommitMap
}

func GetMergesToMasterByAuthorCurrYear() map[string]int {
	commits := getCurrYearMergeGitCommits()

	authorToCommitMap := make(map[string]int)

	for _, commit := range commits {
		authorToCommitMap[commit.Author] += 1
	}

	return authorToCommitMap
}

func GetMostMergesInOneDayCurrYear() MostMergesInOneDay {
	commits := getCurrYearMergeGitCommits()

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

func GetAvgMergesToMasterPerDayCurrYear() float64 {
	commits := getCurrYearMergeGitCommits()
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

func GetFileChangesByAuthorCurrYear(config RepoConfig) map[string]int {
	commits := getCurrYearGitCommits(config)

	authorInsertionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorInsertionsMap[commit.Author] += change.Insertions
			authorInsertionsMap[commit.Author] += change.Deletions
		}
	}

	return authorInsertionsMap
}

func GetCodeInsertionsByAuthorCurrYear(config RepoConfig) map[string]int {
	commits := getCurrYearGitCommits(config)

	authorDeletionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorDeletionsMap[commit.Author] += change.Insertions
		}
	}

	return authorDeletionsMap
}

func GetCodeDeletionsByAuthorCurrYear(config RepoConfig) map[string]int {
	commits := getCurrYearGitCommits(config)

	authorDeletionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorDeletionsMap[commit.Author] += change.Deletions
		}
	}

	return authorDeletionsMap
}

/*
 * TEAM METRICS
 */

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

func GetAuthorCommitsOverTimeCurrYear() []TotalCommitCount {
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

func GetAuthorFileChangesOverTimeCurrYear(config RepoConfig) []TotalFileChangeCount {
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
	commits := getCurrYearGitCommits(config)
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

/*
 * FILE METRICS
 */

func GetPrevYearFileList(r *RepoConfig) []string {
	bytes, err := os.ReadFile(utils.PREV_YEAR_FILE_LIST_FILE)
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

func GetCurrYearFileList() []string {
	bytes, err := os.ReadFile(utils.CURR_YEAR_FILE_LIST_FILE)
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

func HasPrevYearFileBlames() bool {
	_, err := os.ReadFile(utils.PREV_YEAR_FILE_BLAMES_FILE)

	return err == nil
}

func HasCurrYearFileBlames() bool {
	_, err := os.ReadFile(utils.CURR_YEAR_FILE_BLAMES_FILE)

	return err == nil
}

func GetPrevYearFileBlames() []FileBlame {
	bytes, err := os.ReadFile(utils.PREV_YEAR_FILE_BLAMES_FILE)
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

func GetCurrYearFileBlames() []FileBlame {
	bytes, err := os.ReadFile(utils.CURR_YEAR_FILE_BLAMES_FILE)
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

func GetFileChangeRatio(insertionsByAuthor map[string]int, deletionsByAuthor map[string]int) map[string]float64 {
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

func GetCommonlyChangedFiles(config RepoConfig) []FileChangeCount {
	commits := getCurrYearGitCommits(config)
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

func GetFileCountPrevYear() int {
	files := GetPrevYearFileList()

	return len(files)
}

func GetFileCountCurrYear() int {
	files := GetCurrYearFileList()

	return len(files)
}

func GetLargestFilesCurrYear() []FileSize {
	if !HasCurrYearFileBlames() {
		return []FileSize{}
	}

	fileBlames := GetCurrYearFileBlames()

	sort.Slice(fileBlames, func(i int, j int) bool {
		return fileBlames[i].LineCount > fileBlames[j].LineCount
	})

	fileSizes := []FileSize{}
	for _, fileBlame := range fileBlames[0:10] {
		fileSizes = append(fileSizes, FileSize{
			File:      fileBlame.File,
			LineCount: fileBlame.LineCount,
		})
	}

	return fileSizes
}

func GetSmallestFilesCurrYear() []FileSize {
	if !HasCurrYearFileBlames() {
		return []FileSize{}
	}

	fileBlames := GetCurrYearFileBlames()

	sort.Slice(fileBlames, func(i int, j int) bool {
		return fileBlames[i].LineCount < fileBlames[j].LineCount
	})

	fileSizes := []FileSize{}
	for _, fileBlame := range fileBlames[0:10] {
		fileSizes = append(fileSizes, FileSize{
			File:      fileBlame.File,
			LineCount: fileBlame.LineCount,
		})
	}

	return fileSizes
}

func GetTotalLinesOfCodePrevYear() int {
	if !HasPrevYearFileBlames() {
		return 0
	}

	fileBlames := GetPrevYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func GetTotalLinesOfCodeCurrYear() int {
	if !HasCurrYearFileBlames() {
		return 0
	}

	fileBlames := GetCurrYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func GetSizeOfRepoByWeekCurrYear(config RepoConfig) []RepoSizeTimeStamp {
	commits := getCurrYearGitCommits(config)
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

	totalLinesOfCodePrevYear := GetTotalLinesOfCodePrevYear()
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
