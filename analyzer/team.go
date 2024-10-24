package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"fmt"
)

func GetNewEngineerCommitsCurrYear(config Config) []GitCommit {
	pastCommits := getPastGitCommits()

	// userName -> throwaway
	engineersFromPast := make(map[string]int)
	for _, commit := range pastCommits {
		engineersFromPast[commit.Author] = 1
	}

	currYearCommits := getCurrYearGitCommits()

	// username -> bool on whether they have been processed or not
	newEngineers := make(map[string]bool)
	for _, commit := range currYearCommits {
		if _, ok := engineersFromPast[commit.Author]; !ok {
			newEngineers[commit.Author] = false
		}
	}

	newEngineerCommits := []GitCommit{}

	for _, commit := range currYearCommits {
		userName := commit.Author
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
			userName := commit.Author
			engineers[userName] += 1
		}
	}

	return engineers
}

func GetAllUsernames() []string {
	engineers := GetEngineerCommitCountAllTime()

	usernames := []string{}
	for username := range engineers {
		usernames = append(usernames, username)
	}

	return usernames
}

func GetEngineerCommitCountAllTime() map[string]int {
	commits := getGitCommits()
	engineers := make(map[string]int)

	for _, commit := range commits {
		userName := commit.Author
		engineers[userName] += 1
	}

	return engineers
}

func GetEngineerCountCurrYear(config Config) int {
	engineers := GetEngineerCommitCountCurrYear(config)

	return len(engineers)
}

func GetEngineerCountAllTime(config Config) int {
	engineers := GetEngineerCommitCountAllTime()

	return len(engineers)
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
func GetEngineerCommitsOverTimeCurrYear(config Config) []TotalCommitCount {

	// Get list of engineers
	usernames := GetAllUsernames()
	fmt.Println(usernames)

	// Bucket commit counts for all enginers in past
	pastCommits := getPastGitCommits()
	fmt.Println(pastCommits)

	// Get current year commits
	currCommits := getCurrYearGitCommits()
	fmt.Println(currCommits)

	// Bucket commits into days they fall on
	// TODO

	// Create an array of 365 days
	// Iterate through array, adding in "snapshot" of commit counts for each engineer that day, copying previous days into the next one
	// Iterate through array, and add individual TotalCommitCount structs into final array

	return []TotalCommitCount{
		{Date: "2023-01-03T08:00:00.000Z", Name: "Steve Bremer", Value: 24},
		{Date: "2023-01-03T08:00:00.000Z", Name: "Gabe Jensen", Value: 340},
	}
}
