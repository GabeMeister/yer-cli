package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"fmt"
)

func GetNewEngineerCommitsCurrYear() []GitCommit {
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

func GetEngineerCommitCountCurrYear() map[string]int {
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

func GetEngineerCommitCountPrevYear() map[string]int {
	commits := getGitCommits()
	engineers := make(map[string]int)

	for _, commit := range commits {
		commitYear := utils.GetYearFromDateStr(commit.Date)

		if commitYear < CURR_YEAR {
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

func GetEngineerCountCurrYear() int {
	engineers := GetEngineerCommitCountCurrYear()

	return len(engineers)
}

func GetEngineerCountAllTime() int {
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
func GetEngineerCommitsOverTimeCurrYear() []TotalCommitCount {

	// Get list of engineers
	// usernames := GetAllUsernames()

	// Bucket commit counts for all enginers in past
	// pastCommits := getPastGitCommits()
	// fmt.Println(pastCommits)

	// Get current year commits
	currCommits := getCurrYearGitCommits()
	for _, commit := range currCommits {
		fmt.Println()
		fmt.Println(commit.Commit, commit.Author, commit.Date)
		fmt.Println()
	}

	// Bucket commits into days they fall on
	engineerCommitCountPrevYear := GetEngineerCommitCountPrevYear()
	fmt.Println(engineerCommitCountPrevYear)
	// TODO

	// Create map of daily array

	// Iterate through array, adding in "snapshot" of commit counts for each engineer that day, copying previous days into the next one
	// Iterate through array, and add individual TotalCommitCount structs into final array

	return []TotalCommitCount{
		{Date: "2023-01-03T08:00:00.000Z", Name: "Steve Bremer", Value: 24},
		{Date: "2023-01-03T08:00:00.000Z", Name: "Gabe Jensen", Value: 340},
	}
}
