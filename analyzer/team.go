package analyzer

import "GabeMeister/yer-cli/utils"

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
