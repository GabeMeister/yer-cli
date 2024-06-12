package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"
)

/*
 * TYPES
 */

type Config struct {
	Path                  string   `json:"path"`
	IncludeFileExtensions []string `json:"include_file_extensions"`
	ExcludeDirectories    []string `json:"exclude_directories"`
	ExcludeFiles          []string `json:"exclude_files"`
	ExcludeEngineers      []string `json:"exclude_engineers"`
}

type RepoMeta struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
}

type GitCommit struct {
	Commit  string `json:"commit"`
	Author  string `json:"author"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

type GitMergeCommit struct {
	Commit             string
	FirstParentCommit  string // Typically master when you branched off
	SecondParentCommit string // Typically the final commit of the MR
	Message            string
	Author             string
	Email              string
	Date               string
}

type RepoSummary struct {
	PastYearNumCommits int
}

/*
 * PRIVATE
 */
func parseMergeCommits() []GitMergeCommit {
	return nil
}

func parseCommits() []GitCommit {
	return nil
}

func getGitLogs(path string) []GitCommit {
	cmd := exec.Command(
		"git",
		"log",
		"--no-merges",
		"--format=-- Begin --%n-- Commit --%n%H%n-- Author --%n%aN%n-- Email --%n%aE%n-- Date --%n%ad%n-- Message --%n%B%n-- End --")
	cmd.Dir = path

	rawOutput, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	output := string(rawOutput)
	lines := strings.Split(output, "\n")

	var currentState = "BEGIN"
	var currentCommit GitCommit = GitCommit{
		Commit:  "",
		Author:  "",
		Email:   "",
		Message: "",
		Date:    "",
	}
	var commits []GitCommit

	for _, line := range lines {
		if line == "" {
			continue
		}

		switch line {
		case "-- Begin --":
			currentCommit = GitCommit{
				Commit:  "",
				Author:  "",
				Email:   "",
				Message: "",
				Date:    "",
			}
			break
		case "-- Commit --":
			currentState = "COMMIT"
			break
		case "-- Author --":
			currentState = "AUTHOR"
			break
		case "-- Email --":
			currentState = "EMAIL"
			break
		case "-- Date --":
			currentState = "DATE"
			break
		case "-- Message --":
			currentState = "MESSAGE"
			break
		case "-- End --":
			commits = append(commits, currentCommit)
			break
		default:
			switch currentState {
			case "COMMIT":
				currentCommit.Commit = line
				break
			case "AUTHOR":
				currentCommit.Author = line
				break
			case "DATE":
				currentCommit.Date = line
				break
			case "MESSAGE":
				if currentCommit.Message == "" {
					currentCommit.Message += line
				} else {
					currentCommit.Message += "|||" + line
				}
				break
			case "EMAIL":
				currentCommit.Email = line
				break
			default:
				panic(fmt.Sprintf("Unrecognized state: %s", currentState))
			}
		}
	}

	return commits
}

func getRepoMetaData(path string) RepoMeta {
	name := filepath.Base(path)

	return RepoMeta{Name: name, Directory: path}
}

func analyzeRepo(config Config) {
	metaData := getRepoMetaData(config.Path)
	SaveDataToFile(metaData, "./tmp/meta.json")
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

func initConfig(repoPath string) Config {
	config := Config{
		Path:                  repoPath,
		IncludeFileExtensions: []string{},
		ExcludeDirectories:    []string{},
		ExcludeFiles:          []string{},
		ExcludeEngineers:      []string{},
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("./config.json", data, 0644)

	return config
}

func getConfig(path string) Config {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var data Config
	jsonErr := json.Unmarshal(bytes, &data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return data
}

/*
 * PUBLIC
 */

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
