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

func getGitCommits() []GitCommit {
	bytes, err := os.ReadFile("./tmp/commits.json")
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
