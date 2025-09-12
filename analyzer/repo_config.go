package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type RepoConfig struct {
	Id                    int                    `json:"id"`
	Path                  string                 `json:"path"`
	MasterBranchName      string                 `json:"master_branch_name"`
	IncludeFileExtensions []string               `json:"include_file_extensions"`
	ExcludeDirectories    []string               `json:"exclude_directories"`
	ExcludeFiles          []string               `json:"exclude_files"`
	ExcludeAuthors        []string               `json:"exclude_authors"`
	DuplicateAuthors      []DuplicateAuthorGroup `json:"duplicate_authors"`
	AllAuthors            []string               `json:"all_authors"`
	AnalyzeFileBlames     bool                   `json:"analyze_file_blames"`
}

func (r *RepoConfig) gatherMetrics() {
	stashRepo(r.Path)

	currYearErr := checkoutRepoToCommitOrBranchName(r.Path, r.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo to the latest commit")
		panic(currYearErr)
	}

	// We want the latest changes
	pullRepo(r.Path)

	commits := getCommitsFromGitLogs(r, false)
	commitsFileName := r.getCommitsFile()
	saveDataToFile(commits, commitsFileName)

	mergeCommits := getCommitsFromGitLogs(r, true)
	mergeCommitsFileName := r.getMergeCommitsFile()
	saveDataToFile(mergeCommits, mergeCommitsFileName)

	directPushToMasterCommits := getDirectPushToMasterCommitsCurrYear(r)
	directPushFileName := r.getDirectPushesFile()
	saveDataToFile(directPushToMasterCommits, directPushFileName)

	// Prev year files (if possible)
	if r.hasPrevYearCommits() {
		lastCommitPrevYear := r.getLastCommitPrevYear()
		fmt.Printf("Analyzing last year's repo for %s...\n", r.getName())
		prevYearErr := checkoutRepoToCommitOrBranchName(r.Path, lastCommitPrevYear.Commit)
		if prevYearErr != nil {
			fmt.Println("Unable to git checkout repo to last year's files")
			panic(prevYearErr)
		}

		prevYearFiles := getRepoFiles(r, lastCommitPrevYear.Commit)
		saveDataToFile(prevYearFiles, r.getPrevYearFileListFile())

		if r.AnalyzeFileBlames {
			prevYearBlames := getFileBlameSummary(r, prevYearFiles)
			saveDataToFile(prevYearBlames, r.getPrevYearFileBlamesFile())
		}
	}

	// Curr year files
	fmt.Printf("Analyzing this year's repo for %s...\n", r.getName())

	currYearErr = checkoutRepoToCommitOrBranchName(r.Path, r.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo back to the latest commit")
		panic(currYearErr)
	}

	currYearFiles := getRepoFiles(r, r.MasterBranchName)
	saveDataToFile(currYearFiles, r.getCurrYearFileListFile())

	if r.AnalyzeFileBlames {
		currYearBlames := getFileBlameSummary(r, currYearFiles)
		saveDataToFile(currYearBlames, r.getCurrYearFileBlamesFile())
	}
}

/*
 * FILE NAMING
 */

func (r *RepoConfig) getCommitsFile() string {
	return fmt.Sprintf(COMMITS_FILE_TEMPLATE, filepath.Base(r.Path))
}

func (r *RepoConfig) getMergeCommitsFile() string {
	return fmt.Sprintf(MERGE_COMMITS_FILE_TEMPLATE, filepath.Base(r.Path))
}

func (r *RepoConfig) getDirectPushesFile() string {
	return fmt.Sprintf(DIRECT_PUSH_ON_MASTER_COMMITS_FILE_TEMPLATE, filepath.Base(r.Path))
}

func (r *RepoConfig) getPrevYearFileListFile() string {
	return fmt.Sprintf(PREV_YEAR_FILE_LIST_FILE, filepath.Base(r.Path))
}

func (r *RepoConfig) getPrevYearFileBlamesFile() string {
	return fmt.Sprintf(PREV_YEAR_FILE_BLAMES_FILE, filepath.Base(r.Path))
}

func (r *RepoConfig) getCurrYearFileListFile() string {
	return fmt.Sprintf(CURR_YEAR_FILE_LIST_FILE, filepath.Base(r.Path))
}

func (r *RepoConfig) getCurrYearFileBlamesFile() string {
	return fmt.Sprintf(CURR_YEAR_FILE_BLAMES_FILE, filepath.Base(r.Path))
}

func (r *RepoConfig) getRecapFilePath() string {
	return fmt.Sprintf(RECAP_FILE_TEMPLATE, filepath.Base(r.Path))
}

func (r *RepoConfig) hasRecapFile() bool {
	filePath := r.getRecapFilePath()
	_, fileErr := os.Stat(filePath)

	return !errors.Is(fileErr, os.ErrNotExist)
}

func (r *RepoConfig) getName() string {
	return filepath.Base(r.Path)
}

/*
 * COMMIT GETTERS
 */

func (r *RepoConfig) getLastCommitPrevYear() GitCommit {
	commits := r.getPrevYearGitCommits()
	lastIdx := len(commits) - 1

	return commits[lastIdx]
}

func (r *RepoConfig) hasPrevYearCommits() bool {
	commits := r.getPrevYearGitCommits()

	return len(commits) > 0
}

func (r *RepoConfig) getGitCommits() []GitCommit {
	bytes, err := os.ReadFile(r.getCommitsFile())
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
	bytes, err := os.ReadFile(r.getMergeCommitsFile())
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
	bytes, err := os.ReadFile(r.getDirectPushesFile())
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

func (r *RepoConfig) getPrevYearMergeGitCommits() []GitCommit {
	commits := r.getMergeGitCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, PREV_YEAR) {
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

func (r *RepoConfig) getIsMultiYearRepo() bool {
	commits := r.getGitCommits()
	firstCommit := commits[0]

	return utils.GetYearFromDateStr(firstCommit.Date) < CURR_YEAR
}

func (r *RepoConfig) getNumCommitsAllTime() int {
	commits := r.getGitCommits()

	return len(commits)
}

func (r *RepoConfig) getNumCommitsPrevYear() int {
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

func (r *RepoConfig) getNumCommitsCurrYear() int {
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

func (r *RepoConfig) getNumCommitsInPast() int {
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

func (r *RepoConfig) getCommitsByMonthCurrYear() []CommitMonth {
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

func (r *RepoConfig) getCommitsByWeekDayCurrYear() []CommitWeekDay {
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

func (r *RepoConfig) getCommitsByHourCurrYear() []CommitHour {
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

func (r *RepoConfig) getMergeCommitsByMonthCurrYear() []CommitMonth {
	commits := r.getCurrYearMergeGitCommits()

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

func (r *RepoConfig) getMergeCommitsByWeekDayCurrYear() []CommitWeekDay {
	commits := r.getCurrYearMergeGitCommits()

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

func (r *RepoConfig) getMergeCommitsByHourCurrYear() []CommitHour {
	commits := r.getCurrYearMergeGitCommits()

	hourMap := make(map[int]int)

	for _, commit := range commits {
		// Convert into the proper timezone that the user's in
		name, offset := time.Now().Zone()
		currentTZ := time.FixedZone(name, offset)

		t, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", commit.Date)
		if err != nil {
			panic(err)
		}

		userTime := t.In(currentTZ)
		currentTZHour := userTime.Hour()
		hourMap[currentTZHour] += 1
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

func (r *RepoConfig) getMostInsertionsInCommitCurrYear() GitCommit {
	commits := r.getCurrYearGitCommits()

	if len(commits) == 0 {
		return GitCommit{}
	}

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

func (r *RepoConfig) getMostDeletionsInCommitCurrYear() GitCommit {
	commits := r.getCurrYearGitCommits()

	if len(commits) == 0 {
		return GitCommit{}
	}

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

func (r *RepoConfig) getLargestCommitMessageCurrYear() GitCommit {
	commits := r.getCurrYearGitCommits()

	if len(commits) == 0 {
		return GitCommit{}
	}

	largestLengthCommit := commits[0]

	for _, commit := range commits {
		if len(largestLengthCommit.Message) < len(commit.Message) {
			largestLengthCommit = commit
		}
	}

	return largestLengthCommit
}

func (r *RepoConfig) getSmallestCommitMessagesCurrYear() []GitCommit {
	commits := r.getCurrYearGitCommits()

	if len(commits) == 0 {
		return []GitCommit{}
	}

	sort.Slice(commits, func(i int, j int) bool {
		return len(commits[i].Message) < len(commits[j].Message)
	})

	return commits[0:5]
}

func (r *RepoConfig) getCommitMessageHistogramCurrYear() []CommitMessageLengthFrequency {
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

func (r *RepoConfig) getDirectPushesOnMasterByAuthorCurrYear() map[string]int {
	commits := r.getCurrYearDirectPushOnMasterCommits()

	authorToCommitMap := make(map[string]int)

	for _, commit := range commits {
		authorToCommitMap[commit.Author] += 1
	}

	return authorToCommitMap
}

func (r *RepoConfig) getMergesToMasterByAuthorCurrYear() map[string]int {
	commits := r.getCurrYearMergeGitCommits()

	authorToCommitMap := make(map[string]int)

	for _, commit := range commits {
		authorToCommitMap[commit.Author] += 1
	}

	return authorToCommitMap
}

func (r *RepoConfig) getMostMergesInOneDayCurrYear() MostMergesInOneDay {
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

func (r *RepoConfig) getMergesToMasterCurrYear() int {
	commits := r.getCurrYearMergeGitCommits()
	return len(commits)
}

func (r *RepoConfig) getMergesToMasterPrevYear() int {
	commits := r.getPrevYearMergeGitCommits()
	return len(commits)
}

func (r *RepoConfig) getAvgMergesToMasterPerDayCurrYear() float64 {
	commits := r.getCurrYearMergeGitCommits()
	if len(commits) == 0 {
		return 0.0
	}

	numWorkDays := getNumWorkDaysInCurrYear()
	final := float64(len(commits)) / float64(numWorkDays)

	if math.IsNaN(final) {
		panic("GetAvgMergesToMasterPerDayCurrYear() is NaN!")
	}

	return final
}

func (r *RepoConfig) getMergesToMasterAllTime() int {
	commits := r.getMergeGitCommits()
	return len(commits)
}

func (r *RepoConfig) getFileChangesByAuthorCurrYear() map[string]FileChangesSummary {
	authorInsertionsMap := make(map[string]int)
	authorDeletionsMap := make(map[string]int)

	commits := r.getCurrYearGitCommits()
	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorInsertionsMap[commit.Author] += change.Insertions
			authorDeletionsMap[commit.Author] += change.Deletions
		}
	}

	authorFileChangesMap := make(map[string]FileChangesSummary)

	for author := range authorInsertionsMap {
		authorFileChangesMap[author] = FileChangesSummary{
			Insertions: authorInsertionsMap[author],
			Deletions:  authorDeletionsMap[author],
		}
	}

	return authorFileChangesMap
}

func (r *RepoConfig) getCodeInsertionsByAuthorCurrYear() map[string]int {
	commits := r.getCurrYearGitCommits()

	authorDeletionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorDeletionsMap[commit.Author] += change.Insertions
		}
	}

	return authorDeletionsMap
}

func (r *RepoConfig) getCodeDeletionsByAuthorCurrYear() map[string]int {
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

func (repo *RepoConfig) GetDuplicateAuthorList() []string {
	duplicateAuthors := []string{}

	for _, dupGroup := range repo.DuplicateAuthors {
		duplicateAuthors = append(duplicateAuthors, dupGroup.Duplicates...)
	}

	return duplicateAuthors
}

func (r *RepoConfig) getRealAuthorName(authorName string) string {
	for _, dupGroup := range r.DuplicateAuthors {
		for _, dup := range dupGroup.Duplicates {
			if authorName == dup {
				return dupGroup.Real
			}
		}
	}

	return authorName
}

func (r *RepoConfig) getNewAuthorCommitsCurrYear() []GitCommit {
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

func (r *RepoConfig) getAuthorCommitCountCurrYear() map[string]int {
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

func (r *RepoConfig) getAuthorCommitCountPrevYear() map[string]int {
	commits := r.getGitCommits()
	authors := make(map[string]int)

	for _, commit := range commits {
		commitYear := utils.GetYearFromDateStr(commit.Date)

		if commitYear == PREV_YEAR {
			userName := commit.Author
			authors[userName] += 1
		}
	}

	return authors
}

func (r *RepoConfig) getAuthorCommitCountAllTime() map[string]int {
	commits := r.getGitCommits()
	authors := make(map[string]int)

	for _, commit := range commits {
		userName := commit.Author
		authors[userName] += 1
	}

	return authors
}

func (r *RepoConfig) getAuthorCountCurrYear() int {
	authors := r.getAuthorCommitCountCurrYear()

	return len(authors)
}

func (r *RepoConfig) getAuthorCountPrevYear() int {
	authors := r.getAuthorCommitCountPrevYear()

	return len(authors)
}

func (r *RepoConfig) getAuthorCountAllTime() int {
	authors := r.getAuthorCommitCountAllTime()

	return len(authors)
}

func (r *RepoConfig) getAuthorCommitsOverTimeCurrYear() []TotalCommitCount {
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

func (r *RepoConfig) getAllAuthorsList() []string {
	allAuthorsMap := make(map[string]bool)

	for _, commit := range r.getGitCommits() {
		allAuthorsMap[commit.Author] = true
	}

	return utils.MapKeysToSlice(allAuthorsMap)
}

func (r *RepoConfig) getInitialAuthorFileChangeMap() map[string]int {
	allAuthors := r.getAllAuthorsList()
	authorsMap := make(map[string]int)

	for _, a := range allAuthors {
		authorsMap[a] = 0
	}

	return authorsMap
}

// Used mainly as initial starting data for author file changes over time
func (r *RepoConfig) getAuthorTotalFileChangesPrevYear() map[string]int {
	if !r.gasPrevYearFileBlames() || !r.hasCurrYearFileBlames() {
		return make(map[string]int)
	}

	// author => line change count
	fileChangeTracker := make(map[string]int)

	// Bucket file changes for all enginers in past
	prevFileBlames := r.getPrevYearFileBlames()
	for _, fileBlame := range prevFileBlames {
		for author, lineCount := range fileBlame.GitBlame {
			fileChangeTracker[author] += lineCount
		}
	}

	return fileChangeTracker
}

func (r *RepoConfig) getAuthorFileChangesOverTimeCurrYear() TotalFileChangeCount {
	// Make map out of all authors
	// author => line change count
	fileChangeTracker := r.getInitialAuthorFileChangeMap()

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

func (r *RepoConfig) getMostCommitsByAuthorCurrYear() MostSingleDayCommitsByAuthor {
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

func (r *RepoConfig) getTotalLinesOfCodeInRepoByAuthor() map[string]int {
	if !r.hasCurrYearFileBlames() {
		return make(map[string]int)
	}

	authorLineCountMap := make(map[string]int)

	fileBlames := r.getCurrYearFileBlames()
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

func (r *RepoConfig) filterToOnlyIncludedFiles(fileChanges []FileChange) []FileChange {
	filteredFileChanges := utils.Filter(fileChanges, func(c FileChange) bool {
		fileExt := utils.GetFileExtension(c.FilePath)

		isValidFileExt := utils.Includes(r.IncludeFileExtensions, func(ext string) bool {
			return fileExt == ext
		})
		isExcludedFile := utils.Includes(r.ExcludeFiles, func(filePath string) bool {
			return c.FilePath == filePath
		})

		return isValidFileExt && !isExcludedFile
	})

	return filteredFileChanges
}

func (r *RepoConfig) getPrevYearFileList() []string {
	bytes, err := os.ReadFile(r.getPrevYearFileListFile())
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

func (r *RepoConfig) getCurrYearFileList() []string {
	bytes, err := os.ReadFile(r.getCurrYearFileListFile())
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

func (r *RepoConfig) gasPrevYearFileBlames() bool {
	_, err := os.ReadFile(r.getPrevYearFileBlamesFile())

	return err == nil
}

func (r *RepoConfig) hasCurrYearFileBlames() bool {
	_, err := os.ReadFile(r.getCurrYearFileBlamesFile())

	return err == nil
}

func (r *RepoConfig) getPrevYearFileBlames() []FileBlame {
	bytes, err := os.ReadFile(r.getPrevYearFileBlamesFile())
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

func (r *RepoConfig) getCurrYearFileBlames() []FileBlame {
	bytes, err := os.ReadFile(r.getCurrYearFileBlamesFile())
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

func (r *RepoConfig) getFileChangeRatio(insertionsByAuthor map[string]int, deletionsByAuthor map[string]int) map[string]float64 {
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

func (r *RepoConfig) getCommonlyChangedFiles() []FileChangeCount {
	commits := r.getCurrYearGitCommits()

	if len(commits) == 0 {
		return []FileChangeCount{}
	}

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

func (r *RepoConfig) getFileCountPrevYear() int {
	if !r.hasPrevYearCommits() {
		return 0
	}

	files := r.getPrevYearFileList()

	return len(files)
}

func (r *RepoConfig) getFileCountCurrYear() int {
	files := r.getCurrYearFileList()

	return len(files)
}

func (r *RepoConfig) getLargestFilesCurrYear() []FileSize {
	if !r.hasCurrYearFileBlames() {
		return []FileSize{}
	}

	fileBlames := r.getCurrYearFileBlames()

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

func (r *RepoConfig) getSmallestFilesCurrYear() []FileSize {
	if !r.hasCurrYearFileBlames() {
		return []FileSize{}
	}

	fileBlames := r.getCurrYearFileBlames()

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

func (r *RepoConfig) getTotalLinesOfCodePrevYear() int {
	if !r.gasPrevYearFileBlames() {
		return 0
	}

	fileBlames := r.getPrevYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func (r *RepoConfig) getTotalLinesOfCodeCurrYear() int {
	if !r.hasCurrYearFileBlames() {
		return 0
	}

	fileBlames := r.getCurrYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func (r *RepoConfig) getSizeOfRepoByWeekCurrYear() []int {
	commits := r.getCurrYearGitCommits()
	weeks := utils.GetWeeksOfYear()

	weekCommitsMap := make(map[int][]GitCommit)
	for _, week := range weeks {
		weekCommitsMap[week] = nil
	}

	for _, commit := range commits {
		date := utils.GetDateFromISOString(commit.Date)
		_, week := date.ISOWeek()

		// We want to display weeks starting on Sundays, so if a year doesn't fall
		// cleanly on a Sunday we just bucket those first few days together with the
		// first week of the year that has a Sunday. (for example, if Jan 1 was a
		// Wednesday, bucket any commits done from Jan 1-4 into the week starting on
		// Jan 5
		if date.Month() == time.January && week >= 52 {
			week = 1
		} else if date.Month() == time.December && (week == 53 || week == 1) {
			week = 52
		}

		weekCommitsMap[week] = append(weekCommitsMap[week], commit)
	}

	totalLinesOfCodePrevYear := r.getTotalLinesOfCodePrevYear()
	runningTotal := totalLinesOfCodePrevYear
	final := []int{}

	for week := 1; week <= 52; week++ {
		weekCommits := weekCommitsMap[week]

		for _, commit := range weekCommits {
			for _, fileChanges := range commit.FileChanges {
				runningTotal += fileChanges.Insertions
				runningTotal -= fileChanges.Deletions
			}

		}

		final = append(final, runningTotal)
	}

	return final
}

/*
 * RECAP
 */

func (r *RepoConfig) getRepoRecap() (Recap, error) {
	data, err := os.ReadFile(r.getRecapFilePath())
	if err != nil {
		return Recap{}, err
	}
	var repoRecap Recap

	err = json.Unmarshal(data, &repoRecap)
	if err != nil {
		return Recap{}, err
	}

	return repoRecap, nil
}
