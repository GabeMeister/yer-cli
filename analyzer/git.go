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

func GetGitLogs() string {
	cmd := exec.Command("git", "log", "--no-merges", "--format=-- Begin --%n-- Commit --%n%H%n-- Author --%n%aN%n-- Email --%n%aE%n-- Date --%n%ad%n-- Message --%n%B%n-- End --")
	cmd.Dir = "/home/gabe/dev/rb-frontend"

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

	commitsJsonRaw, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		panic(err)
	}

	fileErr := os.WriteFile("commits.json", commitsJsonRaw, 0644)
	if fileErr != nil {
		panic(fileErr)
	}

	return "Done."
}
