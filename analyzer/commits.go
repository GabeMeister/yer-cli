package analyzer

import (
	"encoding/json"
	"os"
)

func GetGitCommits() []GitCommit {
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

func GetTotalNumberOfCommits() int {
	commits := GetGitCommits()
	return len(commits)
}

// func GetNumberOfCommitsMadePastYear() int {
// 	commits := GetGitCommits()
// 	return len(commits)
// }
