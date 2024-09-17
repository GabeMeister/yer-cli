package analyzer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"
)

func AnalyzeManually() {

	var dir string

	for isValid := false; !isValid; isValid = isValidGitRepo(dir) {
		dir = readDir()
	}

	fmt.Println("Valid!")

	config := initConfig(dir)

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

func readDir() string {
	fmt.Println()
	fmt.Println("What directory is your repo is in?")

	dir, err := input_autocomplete.Read("> ")
	if err != nil {
		fmt.Println("Error reading manual input. Please try again.")
		panic(err)
	}

	return dir
}

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

func isValidGitRepo(dir string) bool {
	_, fileErr := os.Stat(dir)

	// TODO: check to make sure read privileges are allowed on directory

	if errors.Is(fileErr, os.ErrNotExist) {
		fmt.Println("Directory not found, please try again.")
		return false
	} else if errors.Is(fileErr, os.ErrPermission) {
		fmt.Println("Unable to access directory, make sure it has proper permissions and try again.")
		return false
	} else {
		gitDirPath := filepath.Join(dir, ".git")
		_, gitDirErr := os.Stat(gitDirPath)
		if errors.Is(gitDirErr, os.ErrNotExist) {
			fmt.Println("No Git repo found in specified directory. Please try again.")
			return false
		}
	}

	return true
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
