package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"
)

func AnalyzeManually() bool {
	var dir string

	for isValid := false; !isValid; isValid = isValidGitRepo(dir) {
		dir = readDir()
	}

	var fileExtensions []string

	for isValid := false; !isValid; isValid = areFileExtensionsValid(fileExtensions) {
		fileExtensions = getFileExtensions()
	}

	var excludedDirs []string
	for isValid := false; !isValid; isValid = areExcludedDirsValid(excludedDirs) {
		excludedDirs = getExcludedDirs()
	}

	config := initConfig(dir, fileExtensions, excludedDirs)

	analyzeRepo(config)

	return true
}

func AnalyzeWithConfig(path string) bool {
	configValid := isValidConfig(path)
	if !configValid {
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

	if strings.Contains(dir, "~") {
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			panic("Could not get user home directory.")
		}

		dir = strings.ReplaceAll(dir, "~", homeDir)
	}

	return dir
}

func getFileExtensions() []string {
	fmt.Println()
	fmt.Println("What file extensions should be analyzed? \nType them comma separated. (For example, type \"ts,js,py,sh\")")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fileExtensions := strings.Split(strings.TrimSpace(text), ",")
	for i := range fileExtensions {
		fileExtensions[i] = strings.TrimSpace(fileExtensions[i])
	}

	return fileExtensions
}

func getExcludedDirs() []string {
	fmt.Println()
	fmt.Println("What directories should be ignored? \nType them comma separated. (For example, type \"node_modules,build\")")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic("Couldn't read in user input when getting excluded directories!")
	}

	excludedDirs := strings.Split(strings.TrimSpace(text), ",")
	for i := range excludedDirs {
		excludedDirs[i] = strings.TrimSpace(excludedDirs[i])
	}

	excludedDirs = utils.Filter(excludedDirs, func(s string) bool {
		return s != ""
	})

	return excludedDirs
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
	engineerCommitCountsAllTime := GetEngineerCommitCountAllTime()

	now := time.Now()
	isoDateString := now.Format(time.RFC3339)

	repoRecap := Recap{
		Name:                        config.Name,
		DateAnalyzed:                isoDateString,
		NumCommitsAllTime:           numCommitsAllTime,
		NumCommitsPrevYear:          numCommitsPrevYear,
		NumCommitsCurrYear:          numCommitsCurrYear,
		NumCommitsInPast:            numCommitsInPast,
		EngineerCommitCountsAllTime: engineerCommitCountsAllTime,
	}
	data, err := json.MarshalIndent(repoRecap, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile(utils.RECAP_FILE, data, 0644)
}

func isValidGitRepo(dir string) bool {
	_, fileErr := os.Stat(dir)

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

func areFileExtensionsValid(fileExtensions []string) bool {
	for _, ext := range fileExtensions {
		if ext == "" {
			fmt.Println("Please enter at least one type of file extension.")
			return false
		}
	}

	return true
}

func areExcludedDirsValid(_ []string) bool {
	// Technically it's always gonna be valid for now but I could easily see
	// something getting added in the future

	return true
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
	var configData Config
	jsonErr := json.Unmarshal(content, &configData)
	if jsonErr != nil {
		// TODO: print instructions on making valid config file
		fmt.Println("Unable to parse json within config file.")
		return false
	}

	// Does it have the right stuff in the json?
	if configData.Path == "" {
		fmt.Println("Missing `path` in the config file.")
		return false
	}

	if configData.Name == "" {
		fmt.Println("Missing `name` in the config file.")
		return false
	}

	if len(configData.IncludeFileExtensions) == 0 {
		fmt.Println("Missing files to include in the config file.")
		return false
	}

	return true
}

func isFileReadable(path string) bool {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	file.Close()

	return true
}
