package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"os"
	"time"
)

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
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR-1) {
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
