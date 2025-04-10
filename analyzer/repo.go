package analyzer

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetAuthorsFromRepo(dir string) []string {
	gitCmd := exec.Command("git", "shortlog", "-s")
	gitCmd.Dir = dir
	rawOutput, _ := gitCmd.Output()

	text := string(rawOutput)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		fmt.Print("\n\n", "*** line ***", "\n", line, "\n\n\n")
	}

	return lines
}
