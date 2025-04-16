package analyzer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetAuthorsFromRepo(dir string, branch string) []string {
	gitCmd := exec.Command("git", "shortlog", branch, "-s")
	gitCmd.Dir = dir

	var stderr bytes.Buffer
	gitCmd.Stderr = &stderr

	rawOutput, err := gitCmd.Output()
	if err != nil {
		fmt.Println("Error running git command:", err)
		fmt.Println("Stderr:", stderr.String())
		return nil
	}

	text := string(rawOutput)
	lines := strings.Split(strings.TrimSpace(text), "\n")

	authors := []string{}
	for _, line := range lines {
		tokens := strings.Split(line, "\t")
		authors = append(authors, tokens[1])
	}

	return authors

}
