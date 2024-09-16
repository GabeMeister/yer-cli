package analyzer

import (
	"encoding/json"
	"errors"
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

func AnalyzeWithConfig(path string) bool {
	validConfig := isValidConfig(path)
	if !validConfig {
		return false
	}

	config := getConfig(path)
	analyzeRepo(config)

	return true
}

/*
 * PRIVATE
 */

func analyzeRepo(config Config) {
	gatherMetrics(config)
	calculateRecap(config)
}

func gatherMetrics(config Config) {
	commits := getGitLogs(config.Path)
	SaveDataToFile(commits, "./tmp/commits.json")
}

func calculateRecap(config Config) {
	numCommitsAllTime := GetNumCommitsAllTime()
	numCommitsPrevYear := GetNumCommitsPrevYear()
	numCommitsCurrYear := GetNumCommitsCurrYear()
	numCommitsInPast := GetNumCommitsInPast()

	repoRecap := Recap{
		Name:               config.Name,
		NumCommitsAllTime:  numCommitsAllTime,
		NumCommitsPrevYear: numCommitsPrevYear,
		NumCommitsCurrYear: numCommitsCurrYear,
		NumCommitsInPast:   numCommitsInPast,
	}
	data, err := json.MarshalIndent(repoRecap, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("./tmp/recap.json", data, 0644)
}

func isValidConfig(path string) bool {
	_, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Could not find config file. Double check that your config file is found at `" + path + "`")
		return false
	}

	// Is the file a json file?
	// Does it have the correct file permissions?

	return true
}
