package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

func GetIsMultiYearRepo() bool {
	commits := getGitCommits()
	firstCommit := commits[0]

	return utils.GetYearFromDateStr(firstCommit.Date) < CURR_YEAR
}

func GetNumCommitsAllTime() int {
	commits := getGitCommits()

	return len(commits)
}

func GetNumCommitsPrevYear() int {
	commits := getGitCommits()
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() == PREV_YEAR
	})

	return len(prevYearCommits)
}

func GetNumCommitsCurrYear() int {
	commits := getGitCommits()
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() == CURR_YEAR
	})

	return len(prevYearCommits)
}

func GetNumCommitsInPast() int {
	commits := getGitCommits()
	prevYearCommits := utils.Filter(commits, func(c GitCommit) bool {
		d, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", c.Date)
		if err != nil {
			panic(err)
		}

		return d.Year() < PREV_YEAR
	})

	return len(prevYearCommits)
}

func GetCommitsByMonthCurrYear() []CommitMonth {
	commits := getCurrYearGitCommits()

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

func GetCommitsByWeekDayCurrYear() []CommitWeekDay {
	commits := getCurrYearGitCommits()

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

func GetCommitsByHourCurrYear() []CommitHour {
	commits := getCurrYearGitCommits()

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

func getGitCommits() []GitCommit {
	bytes, err := os.ReadFile(utils.COMMITS_FILE)
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

func getMergeGitCommits() []GitCommit {
	bytes, err := os.ReadFile(utils.MERGE_COMMITS_FILE)
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

func getDirectPushOnMasterCommits() []GitCommit {
	bytes, err := os.ReadFile(utils.DIRECT_PUSH_ON_MASTER_COMMITS_FILE)
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

// Any commit that has come before the current year
func getPastGitCommits() []GitCommit {
	commits := getGitCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrBeforeYear(commit.Date, CURR_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func getPrevYearGitCommits() []GitCommit {
	commits := getGitCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, PREV_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func getCurrYearGitCommits() []GitCommit {
	commits := getGitCommits()
	final := []GitCommit{}

	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func getCurrYearMergeGitCommits() []GitCommit {
	commits := getMergeGitCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func getCurrYearDirectPushOnMasterCommits() []GitCommit {
	commits := getDirectPushOnMasterCommits()

	final := []GitCommit{}
	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			final = append(final, commit)
		}
	}

	return final
}

func GetMostInsertionsInCommitCurrYear() GitCommit {
	commits := getCurrYearGitCommits()

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

func GetMostDeletionsInCommitCurrYear() GitCommit {
	commits := getCurrYearGitCommits()

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

func GetLargestCommitMessageCurrYear() GitCommit {
	commits := getCurrYearGitCommits()

	largestLengthCommit := commits[0]

	for _, commit := range commits {
		if len(largestLengthCommit.Message) < len(commit.Message) {
			largestLengthCommit = commit
		}
	}

	return largestLengthCommit
}

func GetSmallestCommitMessagesCurrYear() []GitCommit {
	commits := getCurrYearGitCommits()

	sort.Slice(commits, func(i int, j int) bool {
		return len(commits[i].Message) < len(commits[j].Message)
	})

	return commits[0:5]
}

func GetCommitMessageHistogramCurrYear() []CommitMessageLengthFrequency {
	commits := getCurrYearGitCommits()

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

func GetFileChangesByAuthorCurrYear() map[string]int {
	commits := getCurrYearGitCommits()

	authorInsertionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorInsertionsMap[commit.Author] += change.Insertions
			authorInsertionsMap[commit.Author] += change.Deletions
		}
	}

	return authorInsertionsMap
}

func GetCodeInsertionsByAuthorCurrYear() map[string]int {
	commits := getCurrYearGitCommits()

	authorDeletionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorDeletionsMap[commit.Author] += change.Insertions
		}
	}

	return authorDeletionsMap
}

func GetCodeDeletionsByAuthorCurrYear() map[string]int {
	commits := getCurrYearGitCommits()

	authorDeletionsMap := make(map[string]int)

	for _, commit := range commits {
		for _, change := range commit.FileChanges {
			authorDeletionsMap[commit.Author] += change.Deletions
		}
	}

	return authorDeletionsMap
}
