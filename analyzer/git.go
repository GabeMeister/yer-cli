package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// commit 37466879ef473db98830d3695ae4cc8bcf48ac33
// Author: Gabe Jensen <thegabejensen@gmail.com>
// Date:   Fri Oct 14 07:45:09 2022 -0700

//     Styled up the results for easier viewing

// commit 5e1dda9198545d27bc50edb1fcbfbbd85bab0c72
// Author: Gabe Jensen <thegabejensen@gmail.com>
// Date:   Wed Oct 12 22:03:42 2022 -0700

//     Added mobile styling and video embeds

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
	return RepoMeta{Name: "blah", Directory: path}
}

func analyzeRepo(config Config) {
	metaData := getRepoMetaData(config.Path)
	fmt.Println(metaData)
	SaveDataToFile(metaData, "./tmp/meta.json")
	commits := getGitLogs(config.Path)
	SaveDataToFile(commits, "./tmp/commits.json")
}

func initConfig(repoPath string) Config {
	// TODO
}

func AnalyzeManually() {
	fmt.Println()
	fmt.Println("What directory is your repo is in?")
	fmt.Print("> ")

	var repoPath string
	fmt.Scanln(&repoPath)

	initConfig(repoPath)
	config := getConfig("./config.json")

	analyzeRepo(config)
}

func AnalyzeWithConfig(configPath string) {
	fmt.Println()
	config := getConfig(configPath)
	analyzeRepo(config)
}

type Config struct {
	Path                  string   `json:"path"`
	IncludeFileExtensions []string `json:"include_file_extensions"`
	ExcludeDirectories    []string `json:"exclude_directories"`
	ExcludeFiles          []string `json:"exclude_files"`
	ExcludeEngineers      []string `json:"exclude_engineers"`
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
