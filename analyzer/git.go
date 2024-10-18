package analyzer

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func getGitLogs(path string) []GitCommit {
	s := GetSpinner()

	fmt.Println()
	s.Suffix = " Retrieving git logs..."
	s.Start()

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

	s.Suffix = " Analyzing git logs..."

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

	s.Stop()

	fileChangeSummary := getFileChangeSummary(path)

	for i := range commits {
		commits[i].FileChanges = fileChangeSummary[commits[i].Commit]
	}

	return commits
}

func isFileChangeLine(line string) bool {
	// Regex for matching email addresses
	var emailRegex = regexp.MustCompile(`^\d+\s+\d+\s+.+$`)
	return emailRegex.MatchString(line)
}

func getFileChangeSummary(path string) map[string][]FileChange {
	s := GetSpinner()
	s.Suffix = " Retrieving file changes..."
	s.Start()

	cmd := exec.Command(
		"git",
		"log",
		"--no-merges",
		"--reverse",
		"--numstat")
	cmd.Dir = path

	rawOutput, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	output := string(rawOutput)

	s.Suffix = " Analyzing file changes..."

	lines := strings.Split(output, "\n")
	fileChangeMap := make(map[string][]FileChange)

	currHash := ""
	currFileChanges := []FileChange{}

	for _, line := range lines {
		if strings.HasPrefix(line, "commit") {
			// We found a new commit, so we need to add the previous commit in and
			// reset the temp variables
			if currHash != "" {
				fileChangeMap[currHash] = currFileChanges
				currHash = ""
				currFileChanges = []FileChange{}
			}

			// Initialize a "new" commit
			tokens := strings.Split(line, " ")
			currHash = tokens[1]
		} else if isFileChangeLine(line) {
			// Regex to match any whitespace
			whitespace := regexp.MustCompile(`\s+`)

			// Split the string by any whitespace
			parts := whitespace.Split(line, -1)
			insertions, _ := strconv.Atoi(parts[0])
			deletions, _ := strconv.Atoi(parts[1])
			filePath := parts[2]

			currFileChanges = append(currFileChanges, FileChange{
				Insertions: insertions,
				Deletions:  deletions,
				FilePath:   filePath,
			})
		}
	}

	// Add in the final commit
	fileChangeMap[currHash] = currFileChanges

	s.Stop()

	return fileChangeMap

}
