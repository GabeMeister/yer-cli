package analyzer

import (
	"fmt"
	"os/exec"
	"strings"
)

func getGitLogs(path string) []GitCommit {
	cmd := exec.Command(
		"git",
		"log",
		"--no-merges",
		"--reverse",
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
		case "-- Commit --":
			currentState = "COMMIT"
		case "-- Author --":
			currentState = "AUTHOR"
		case "-- Email --":
			currentState = "EMAIL"
		case "-- Date --":
			currentState = "DATE"
		case "-- Message --":
			currentState = "MESSAGE"
		case "-- End --":
			commits = append(commits, currentCommit)
		default:
			switch currentState {
			case "COMMIT":
				currentCommit.Commit = line
			case "AUTHOR":
				currentCommit.Author = line
			case "DATE":
				currentCommit.Date = line
			case "MESSAGE":
				if currentCommit.Message == "" {
					currentCommit.Message += line
				} else {
					currentCommit.Message += "|||" + line
				}
			case "EMAIL":
				currentCommit.Email = line
			default:
				panic(fmt.Sprintf("Unrecognized state: %s", currentState))
			}
		}
	}

	return commits
}
