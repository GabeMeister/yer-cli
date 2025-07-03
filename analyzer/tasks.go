package analyzer

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func AnalyzeRepos(calculateOnly bool) bool {
	configValid := isValidConfig(DEFAULT_CONFIG_FILE)
	if !configValid {
		return false
	}

	config := MustGetConfig(DEFAULT_CONFIG_FILE)

	for _, r := range config.Repos {
		// Check if repo is "clean" (on master branch, and no unstaged changes)
		if !isRepoClean(r.Path, r.MasterBranchName) {
			fmt.Println(`
This tool will inspect your git repo(s) at various commits.
Please make sure your repo is on master (or main), 
and there are no unstaged changes before continuing.

Press enter to continue...`)
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')

			break
		}
	}

	for _, r := range config.Repos {
		if !calculateOnly {
			r.gatherMetrics()
			config.updateDuplicateAuthors(&r)
		}
		r.calculateRecap()
	}

	err := config.calculateMultiRepoRecap()
	if err != nil {
		panic(err)
	}

	return err == nil
}

func isRepoClean(dir string, masterBranch string) bool {
	// Check if we're on master branch
	branchCmd := exec.Command("git", "branch", "--show-current")
	branchCmd.Dir = dir
	branchOutput, err := branchCmd.Output()
	if err != nil {
		return false
	}

	currentBranch := strings.TrimSpace(string(branchOutput))
	if currentBranch != masterBranch {
		return false
	}

	// Check for unstaged changes
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = dir
	statusOutput, err := statusCmd.Output()
	if err != nil {
		return false
	}

	return len(statusOutput) == 0
}

func isValidConfig(path string) bool {
	_, fileErr := os.Stat(path)

	// Does it even exist?
	if errors.Is(fileErr, os.ErrNotExist) {
		fmt.Println("Could not find config file. Double check that your config file is found at `" + path + "`")
		return false
	}

	// Is it even a .json file?
	if !strings.HasSuffix(strings.ToLower(path), ".json") {
		fmt.Println("File found at " + path + " is not a json file")
		return false
	}

	// Can you read the file?
	if !isFileReadable(path) {
		fmt.Println("File found at " + path + " does not have read permissions.")
		return false
	}

	// Does it have the right schema?
	content, fileErr := os.ReadFile(path)
	if fileErr != nil {
		panic("Could not read file")
	}

	// Does it even contain json?
	var configData ConfigFile
	jsonErr := json.Unmarshal(content, &configData)
	if jsonErr != nil {
		// TODO: print instructions on making valid config file
		fmt.Println("Unable to parse json within config file.")
		return false
	}

	if len(configData.Repos) == 0 {
		fmt.Println("Must have at least one repo config specified within config file.")
		return false
	}

	repoConfig := configData.Repos[0]

	// Does it have the right stuff in the json?
	if repoConfig.Path == "" {
		fmt.Println("Missing `path` in the config file.")
		return false
	}

	if repoConfig.Path == "" {
		fmt.Println("Missing `path` in the config file.")
		return false
	}

	if len(repoConfig.IncludeFileExtensions) == 0 {
		fmt.Println("Missing files to include in the config file.")
		return false
	}

	return true
}
