package analyzer

import (
	"encoding/json"
	"fmt"
	"os"

	input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"
)

func AnalyzeManually() {
	fmt.Println()
	fmt.Println("What directory is your repo is in?")
	repoDir, err := input_autocomplete.Read("> ")
	if err != nil {
		panic(err)
	}

	config := initConfig(repoDir)

	analyzeRepo(config)
}

func AnalyzeWithConfig(path string) {
	config := getConfig(path)
	analyzeRepo(config)
}

func analyzeRepo(config Config) {
	gatherMetrics(config)
	calculateRecap()
}

func gatherMetrics(config Config) {
	commits := getGitLogs(config.Path)
	SaveDataToFile(commits, "./tmp/commits.json")
}

func calculateRecap() {
	numCommitsAllTime := GetNumCommitsAllTime()
	numCommitsPrevYear := GetNumCommitsPrevYear()
	numCommitsCurrYear := GetNumCommitsCurrYear()
	numCommitsInPast := GetNumCommitsInPast()

	repoSummary := Recap{
		NumCommitsAllTime:  numCommitsAllTime,
		NumCommitsPrevYear: numCommitsPrevYear,
		NumCommitsCurrYear: numCommitsCurrYear,
		NumCommitsInPast:   numCommitsInPast,
	}
	data, err := json.MarshalIndent(repoSummary, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("./tmp/summary.json", data, 0644)
}
