package analyzer

import (
	"GabeMeister/yer-cli/utils"
)

func GetNewEngineerCommitsCurrYear(config Config) []GitCommit {
	pastCommits := getPastGitCommits()

	// userName -> throwaway
	engineersFromPast := make(map[string]int)
	for _, commit := range pastCommits {
		userName := getRealUsername(commit.Author, config)
		engineersFromPast[userName] = 1
	}

	currYearCommits := getCurrYearGitCommits()

	// username -> bool on whether they have been processed or not
	newEngineers := make(map[string]bool)
	for _, commit := range currYearCommits {
		userName := getRealUsername(commit.Author, config)

		if _, ok := engineersFromPast[userName]; !ok {
			newEngineers[userName] = false
		}
	}

	newEngineerCommits := []GitCommit{}

	for _, commit := range currYearCommits {
		userName := getRealUsername(commit.Author, config)
		processed, ok := newEngineers[userName]

		if ok && !processed {
			newEngineerCommits = append(newEngineerCommits, commit)
			newEngineers[userName] = true
		}
	}

	return newEngineerCommits
}

func GetEngineerCommitCountCurrYear(config Config) map[string]int {
	commits := getGitCommits()
	engineers := make(map[string]int)

	for _, commit := range commits {
		if utils.IsDateStrInYear(commit.Date, CURR_YEAR) {
			userName := getRealUsername(commit.Author, config)
			engineers[userName] += 1
		}
	}

	return engineers
}

func GetEngineerCommitCountAllTime(config Config) map[string]int {
	commits := getGitCommits()
	engineers := make(map[string]int)

	for _, commit := range commits {
		userName := getRealUsername(commit.Author, config)
		engineers[userName] += 1
	}

	return engineers
}

func GetEngineerCountCurrYear(config Config) int {
	engineers := GetEngineerCommitCountCurrYear(config)

	return len(engineers)
}
