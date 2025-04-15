package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"
)

func AnalyzeManually() bool {
	var dir string

	for isValid := false; !isValid; isValid = IsValidGitRepo(dir) {
		dir = readDir()
	}

	masterBranch := GetMasterBranchName(dir)

	// Check if repo is "clean" (on master branch, and no unstaged changes)
	if !isRepoClean(dir, masterBranch) {
		fmt.Println(`
This tool will inspect your git repo at various commits.
Please make sure your repo is on master (or main), 
and there are no unstaged changes before continuing.

Press enter to continue...`)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	}

	var fileExtensions []string

	for isValid := false; !isValid; isValid = areFileExtensionsValid(fileExtensions) {
		fileExtensions = getFileExtensions()
	}

	includeFileBlames := getShouldIncludeFileBlames()

	var excludedDirs []string
	for isValid := false; !isValid; isValid = areExcludedDirsValid(excludedDirs) {
		excludedDirs = getExcludedDirs()
	}

	config := InitConfig(ConfigFileOptions{
		RepoDir:                dir,
		MasterBranchName:       masterBranch,
		IncludedFileExtensions: fileExtensions,
		ExcludedDirs:           excludedDirs,
		DuplicateAuthors:     []DuplicateAuthorGroup{},
		IncludeFileBlames:      includeFileBlames,
	})
	// For now, we're just handling 1, we can handle multiple repos in a
	// concurrent way later
	repoConfig := config.Repos[0]

	gatherMetrics(config.Repos[0])

	duplicateAuthors := getDuplicateUsers()

	err := updateDuplicateAuthors(utils.DEFAULT_CONFIG_FILE, duplicateAuthors)
	if err != nil {
		panic(err)
	}

	calculateRecap(repoConfig)

	return true
}

func AnalyzeWithConfig(path string) bool {
	configValid := isValidConfig(path)
	if !configValid {
		return false
	}

	config := GetConfig(path)

	// For now, we're just handling 1 repo at a time
	repoConfig := config.Repos[0]

	// Check if repo is "clean" (on master branch, and no unstaged changes)
	if !isRepoClean(repoConfig.Path, repoConfig.MasterBranchName) {
		fmt.Println(`
This tool will inspect your git repo at various commits.
Please make sure your repo is on master (or main), 
and there are no unstaged changes before continuing.

Press enter to continue...`)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	}

	gatherMetrics(repoConfig)
	updateDuplicateAuthors(path, repoConfig.DuplicateAuthors)
	calculateRecap(repoConfig)

	return true
}

/*
 * PRIVATE
 */

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

	if dir == "." {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory")
			panic(err)
		}

		dir = currentDir
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

func getShouldIncludeFileBlames() bool {
	fmt.Println()
	fmt.Println("Do you want to include advanced stats? (Y/n)\nNote: this is usually fine for most repos, but answer no for repos that have EXTREMELY large commit histories (>100,000 commits)")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic("Couldn't read in user input when getting answer for including file blames!")
	}

	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	return strings.HasPrefix(text, "y") || text == ""
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

func GetMasterBranchName(dir string) string {
	masterBranchCmd := exec.Command("git", "branch", "-r")
	masterBranchCmd.Dir = dir
	rawOutput, _ := masterBranchCmd.Output()
	text := string(rawOutput)
	lines := strings.Split(text, "\n")
	idx := slices.IndexFunc(lines, func(line string) bool {
		return strings.Contains(line, "origin/HEAD")
	})
	masterBranchLine := lines[idx]

	masterBranchName := strings.ReplaceAll(masterBranchLine, "origin/HEAD -> origin/", "")
	masterBranchName = strings.TrimSpace(masterBranchName)

	return masterBranchName
}

func GetFileExtensionsInRepo(dir string) []string {
	lsFilesCmd := exec.Command("git", "ls-files")
	lsFilesCmd.Dir = dir
	rawOutput, _ := lsFilesCmd.Output()
	text := string(rawOutput)
	lines := strings.Split(text, "\n")

	extMap := make(map[string]bool)
	for _, line := range lines {
		fileExt := strings.ReplaceAll(filepath.Ext(line), ".", "")
		if fileExt != "" {
			extMap[fileExt] = true
		}
	}

	fileExtensions := []string{}
	for fileExt := range extMap {
		if slices.Contains(SUPPORTED_FILE_EXTENSIONS, fileExt) {
			fileExtensions = append(fileExtensions, fileExt)
		}
	}

	return fileExtensions
}

// A lot of times in repos, somehow the same user has two different git usernames
// (for example, Gabe Jensen and GabeJensen). It could be because they changed
// laptops, decided to change their user name randomly, etc. To make the stats
// more accurate, we "bucket" duplicate usernames into one.
func getDuplicateUsers() []DuplicateAuthorGroup {
	// commits := getGitCommits()
	// // Username -> int
	// userMap := make(map[string]int)

	// for _, commit := range commits {
	// 	userMap[commit.Author] = 1
	// }

	// fmt.Println()
	// fmt.Println("The list of git usernames are:")
	// fmt.Println()

	// userNames := []string{}
	// for userName := range userMap {
	// 	userNames = append(userNames, userName)
	// }

	// // Use this instead of strings.Sort() because you want lowercase and uppercase
	// // usernames to be next to each other. (For example, "Kaleb Trotter" and
	// // "ktrotter")
	// slices.SortFunc(userNames, func(userName1 string, userName2 string) int {
	// 	s1 := strings.ToLower(userName1)
	// 	s2 := strings.ToLower(userName2)

	// 	if s1 < s2 {
	// 		return -1
	// 	} else if s1 > s2 {
	// 		return 1
	// 	} else {
	// 		return 0
	// 	}
	// })

	// for _, userName := range userNames {
	// 	fmt.Println(userName)
	// }

	// fmt.Println()
	// fmt.Print("Are there any duplicates? (y/N) ")

	// reader := bufio.NewReader(os.Stdin)
	// text, err := reader.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }
	// answer := strings.TrimSpace(text)

	// if len(answer) > 0 && strings.ToLower(string(answer[0])) == "y" {
	// 	duplicateAuthorMap := []DuplicateAuthorGroup{}

	// 	for i := 0; i < 1000; i++ {
	// 		fmt.Println()

	// 		fillerWord := "a"
	// 		if i >= 1 {
	// 			fillerWord = "another"
	// 		}

	// 		fmt.Println("Type " + fillerWord + " duplicate username (or type \"exit\" when done):")
	// 		fmt.Print("> ")

	// 		reader := bufio.NewReader(os.Stdin)
	// 		text, err := reader.ReadString('\n')
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		duplicateUsername := strings.TrimSpace(text)

	// 		if duplicateUsername == "exit" {
	// 			break
	// 		}

	// 		fmt.Println()
	// 		fmt.Println("Type the real username for " + duplicateUsername + ":")
	// 		fmt.Print("> ")

	// 		reader = bufio.NewReader(os.Stdin)
	// 		text, err = reader.ReadString('\n')
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		realUsername := strings.TrimSpace(text)

	// 		duplicateAuthorMap[duplicateUsername] = realUsername

	// 		userNames = utils.Delete(userNames, func(item string) bool { return item == duplicateUsername })

	// 		fmt.Println()
	// 		fmt.Println("The remaining git usernames are:")
	// 		fmt.Println()

	// 		for _, userName := range userNames {
	// 			fmt.Println(userName)
	// 		}
	// 	}

	// 	return duplicateAuthorMap
	// } else {
	// 	return map[string]string{}
	// }

	// This function is going away so just put this here so things build
	return []DuplicateAuthorGroup{}
}

func gatherMetrics(config RepoConfig) {
	stashRepo(config.Path)

	currYearErr := checkoutRepoToCommitOrBranchName(config, config.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo to the latest commit")
		panic(currYearErr)
	}

	// We want the latest changes
	pullRepo(config.Path)

	commits := getCommitsFromGitLogs(config, false)
	SaveDataToFile(commits, utils.COMMITS_FILE)

	mergeCommits := getCommitsFromGitLogs(config, true)
	SaveDataToFile(mergeCommits, utils.MERGE_COMMITS_FILE)

	directPushToMasterCommits := getDirectPushToMasterCommitsCurrYear(config)
	SaveDataToFile(directPushToMasterCommits, utils.DIRECT_PUSH_ON_MASTER_COMMITS_FILE)

	// Prev year files (if possible)
	if hasPrevYearCommits() {
		lastCommitPrevYear := getLastCommitPrevYear()
		fmt.Println("Analyzing last year's repo...")
		prevYearErr := checkoutRepoToCommitOrBranchName(config, lastCommitPrevYear.Commit)
		if prevYearErr != nil {
			fmt.Println("Unable to git checkout repo to last year's files")
			panic(prevYearErr)
		}

		prevYearFiles := getRepoFiles(config, lastCommitPrevYear.Commit)
		SaveDataToFile(prevYearFiles, utils.PREV_YEAR_FILE_LIST_FILE)

		if config.IncludeFileBlames {
			prevYearBlames := GetFileBlameSummary(config, prevYearFiles)
			SaveDataToFile(prevYearBlames, utils.PREV_YEAR_FILE_BLAMES_FILE)
		}
	}

	// Curr year files
	fmt.Println("Analyzing this year's repo...")

	currYearErr = checkoutRepoToCommitOrBranchName(config, config.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo back to the latest commit")
		panic(currYearErr)
	}

	currYearFiles := getRepoFiles(config, config.MasterBranchName)
	SaveDataToFile(currYearFiles, utils.CURR_YEAR_FILE_LIST_FILE)

	if config.IncludeFileBlames {
		currYearBlames := GetFileBlameSummary(config, currYearFiles)
		SaveDataToFile(currYearBlames, utils.CURR_YEAR_FILE_BLAMES_FILE)
	}
}

func calculateRecap(config RepoConfig) {
	s := GetSpinner()

	fmt.Println()
	s.Suffix = " Calculating repo stats..."
	s.Start()

	isMultiYearRepo := GetIsMultiYearRepo()
	numCommitsAllTime := GetNumCommitsAllTime()
	numCommitsPrevYear := GetNumCommitsPrevYear()
	numCommitsCurrYear := GetNumCommitsCurrYear()
	newAuthorCommitsCurrYear := GetNewAuthorCommitsCurrYear()
	newAuthorCountCurrYear := len(newAuthorCommitsCurrYear)
	newAuthorListCurrYear := utils.Map(newAuthorCommitsCurrYear, func(commit GitCommit) string {
		return commit.Author
	})
	authorCommitCountsCurrYear := GetAuthorCommitCountCurrYear()
	authorCommitCountsAllTime := GetAuthorCommitCountAllTime()
	authorCountCurrYear := GetAuthorCountCurrYear()
	authorCountAllTime := GetAuthorCountAllTime()
	authorCommitsOverTimeCurrYear := GetAuthorCommitsOverTimeCurrYear()
	authorFileChangesOverTimeCurrYear := GetAuthorFileChangesOverTimeCurrYear()
	commitsByMonthCurrYear := GetCommitsByMonthCurrYear()
	commitsByWeekDayCurrYear := GetCommitsByWeekDayCurrYear()
	commitsByHourCurrYear := GetCommitsByHourCurrYear()
	mostSingleDayCommitsByAuthorCurrYear := GetMostCommitsByAuthorCurrYear()
	mostInsertionsInCommitCurrYear := GetMostInsertionsInCommitCurrYear()
	mostDeletionsInCommitCurrYear := GetMostDeletionsInCommitCurrYear()
	largestCommitMessageCurrYear := GetLargestCommitMessageCurrYear()
	smallestCommitMessagesCurrYear := GetSmallestCommitMessagesCurrYear()
	commitMessageHistogramCurrYear := GetCommitMessageHistogramCurrYear()
	directPushesOnMasterByAuthorCurrYear := GetDirectPushesOnMasterByAuthorCurrYear()
	mergesToMasterByAuthorCurrYear := GetMergesToMasterByAuthorCurrYear()
	mostMergesInOneDayCurrYear := GetMostMergesInOneDayCurrYear()
	avgMergesToMasterPerDayCurrYear := GetAvgMergesToMasterPerDayCurrYear()
	fileChangesByAuthorCurrYear := GetFileChangesByAuthorCurrYear()
	codeInsertionsByAuthorCurrYear := GetCodeInsertionsByAuthorCurrYear()
	codeDeletionsByAuthorCurrYear := GetCodeDeletionsByAuthorCurrYear()
	fileChangeRatioCurrYear := GetFileChangeRatio(codeInsertionsByAuthorCurrYear, codeDeletionsByAuthorCurrYear)
	commonlyChangedFiles := GetCommonlyChangedFiles()
	fileCountPrevYear := GetFileCountPrevYear()
	fileCountCurrYear := GetFileCountCurrYear()
	largestFilesCurrYear := GetLargestFilesCurrYear()
	smallestFilesCurrYear := GetSmallestFilesCurrYear()
	totalLinesOfCodePrevYear := GetTotalLinesOfCodePrevYear()
	totalLinesOfCodeCurrYear := GetTotalLinesOfCodeCurrYear()
	totalLinesOfCodeInRepoByAuthor := GetTotalLinesOfCodeInRepoByAuthor()
	sizeOfRepoByWeekCurrYear := GetSizeOfRepoByWeekCurrYear()

	now := time.Now()
	isoDateString := now.Format(time.RFC3339)

	fileCountPercentDifference := (float64(fileCountCurrYear) - float64(fileCountPrevYear)) / float64(fileCountPrevYear)
	if math.IsNaN(fileCountPercentDifference) {
		panic("File count percent difference is NaN!")
	}

	repoRecap := Recap{
		// Metadata
		Version:            "0.0.1",
		Name:               config.Name,
		DateAnalyzed:       isoDateString,
		IsMultiYearRepo:    isMultiYearRepo,
		IncludesFileBlames: config.IncludeFileBlames,

		// Commits
		NumCommitsAllTime:               numCommitsAllTime,
		NumCommitsPrevYear:              numCommitsPrevYear,
		NumCommitsCurrYear:              numCommitsCurrYear,
		CommitsByMonthCurrYear:          commitsByMonthCurrYear,
		CommitsByWeekDayCurrYear:        commitsByWeekDayCurrYear,
		CommitsByHourCurrYear:           commitsByHourCurrYear,
		MostInsertionsInCommitCurrYear:  mostInsertionsInCommitCurrYear,
		MostDeletionsInCommitCurrYear:   mostDeletionsInCommitCurrYear,
		LargestCommitMessageCurrYear:    largestCommitMessageCurrYear,
		SmallestCommitMessagesCurrYear:  smallestCommitMessagesCurrYear,
		CommitMessageHistogramCurrYear:  commitMessageHistogramCurrYear,
		MostMergesInOneDayCurrYear:      mostMergesInOneDayCurrYear,
		AvgMergesToMasterPerDayCurrYear: avgMergesToMasterPerDayCurrYear,
		CommonlyChangedFiles:            commonlyChangedFiles,

		// Files
		FileCountPrevYear:          fileCountPrevYear,
		FileCountCurrYear:          fileCountCurrYear,
		FileCountPercentDifference: fileCountPercentDifference,
		LargestFilesCurrYear:       largestFilesCurrYear,
		SmallestFilesCurrYear:      smallestFilesCurrYear,
		TotalLinesOfCodePrevYear:   totalLinesOfCodePrevYear,
		TotalLinesOfCodeCurrYear:   totalLinesOfCodeCurrYear,
		SizeOfRepoByWeekCurrYear:   sizeOfRepoByWeekCurrYear,

		// Team
		NewAuthorCommitsCurrYear:             newAuthorCommitsCurrYear,
		NewAuthorCountCurrYear:               newAuthorCountCurrYear,
		NewAuthorListCurrYear:                newAuthorListCurrYear,
		AuthorCommitCountsCurrYear:           authorCommitCountsCurrYear,
		AuthorCommitCountsAllTime:            authorCommitCountsAllTime,
		AuthorCountCurrYear:                  authorCountCurrYear,
		AuthorCountAllTime:                   authorCountAllTime,
		AuthorCommitsOverTimeCurrYear:        authorCommitsOverTimeCurrYear,
		AuthorFileChangesOverTimeCurrYear:    authorFileChangesOverTimeCurrYear,
		MostSingleDayCommitsByAuthorCurrYear: mostSingleDayCommitsByAuthorCurrYear,
		DirectPushesOnMasterByAuthorCurrYear: directPushesOnMasterByAuthorCurrYear,
		MergesToMasterByAuthorCurrYear:       mergesToMasterByAuthorCurrYear,
		FileChangesByAuthorCurrYear:          fileChangesByAuthorCurrYear,
		FileChangeRatioByAuthorCurrYear:      fileChangeRatioCurrYear,
		TotalLinesOfCodeInRepoByAuthor:       totalLinesOfCodeInRepoByAuthor,
	}
	data, err := json.MarshalIndent(repoRecap, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile(utils.RECAP_FILE, data, 0644)

	s.Stop()
}

func IsValidGitRepo(dir string) bool {
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

	if repoConfig.Name == "" {
		fmt.Println("Missing `name` in the config file.")
		return false
	}

	if len(repoConfig.IncludeFileExtensions) == 0 {
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
