package analyzer

import (
	"encoding/json"
	"fmt"
	"os"

	input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"
)

func analyzeRepo(config Config) {
	commits := getGitLogs(config.Path)
	SaveDataToFile(commits, "./tmp/commits.json")

	numCommits := GetTotalNumberOfCommits()
	fmt.Println(numCommits)
	repoSummary := RepoSummary{
		PastYearNumCommits: numCommits,
	}
	data, err := json.MarshalIndent(repoSummary, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("./tmp/summary.json", data, 0644)
}

func AnalyzeManually() {
	fmt.Println()
	fmt.Println("What directory is your repo is in?")
	repoPath, err := input_autocomplete.Read("> ")
	if err != nil {
		panic(err)
	}

	config := initConfig(repoPath)

	analyzeRepo(config)
}

func AnalyzeWithConfig(configPath string) {
	fmt.Println()
	config := getConfig(configPath)
	analyzeRepo(config)
}
