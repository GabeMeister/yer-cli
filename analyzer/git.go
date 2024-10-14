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

// TODO
func parseMergeCommits() []GitMergeCommit {
	return nil
}

// TODO
func parseCommits() []GitCommit {
	return nil
}
