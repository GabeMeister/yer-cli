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
	"strings"
	"time"
)

func AnalyzeRepos() bool {
	configValid := isValidConfig(DEFAULT_CONFIG_FILE)
	if !configValid {
		return false
	}

	config := MustGetConfig(DEFAULT_CONFIG_FILE)

	for _, r := range config.Repos {
		// Check if repo is "clean" (on master branch, and no unstaged changes)
		if !isRepoClean(r.Path, r.MasterBranchName) {
			fmt.Println(`
This tool will inspect your git repo at various commits.
Please make sure your repo is on master (or main), 
and there are no unstaged changes before continuing.

Press enter to continue...`)
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}

		gatherMetrics(r)
		updateDuplicateAuthors(DEFAULT_CONFIG_FILE, r.DuplicateAuthors)
		calculateRecap(r)
	}

	return true
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

func gatherMetrics(r RepoConfig) {
	stashRepo(r.Path)

	currYearErr := r.checkoutRepoToCommitOrBranchName(r.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo to the latest commit")
		panic(currYearErr)
	}

	// We want the latest changes
	pullRepo(r.Path)

	commits := getCommitsFromGitLogs(r, false)
	commitsFileName := r.GetCommitsFile()
	SaveDataToFile(commits, commitsFileName)

	mergeCommits := getCommitsFromGitLogs(r, true)
	mergeCommitsFileName := r.GetMergeCommitsFile()
	SaveDataToFile(mergeCommits, mergeCommitsFileName)

	directPushToMasterCommits := getDirectPushToMasterCommitsCurrYear(r)
	directPushFileName := r.GetDirectPushesFile()
	SaveDataToFile(directPushToMasterCommits, directPushFileName)

	utils.Pause()

	// Prev year files (if possible)
	if r.hasPrevYearCommits() {
		lastCommitPrevYear := r.getLastCommitPrevYear()
		fmt.Println("Analyzing last year's repo...")
		prevYearErr := r.checkoutRepoToCommitOrBranchName(lastCommitPrevYear.Commit)
		if prevYearErr != nil {
			fmt.Println("Unable to git checkout repo to last year's files")
			panic(prevYearErr)
		}

		prevYearFiles := r.getRepoFiles(lastCommitPrevYear.Commit)
		SaveDataToFile(prevYearFiles, r.GetPrevYearFileListFile())

		if r.IncludeFileBlames {
			prevYearBlames := GetFileBlameSummary(r, prevYearFiles)
			SaveDataToFile(prevYearBlames, r.GetPrevYearFileBlamesFile())
		}
	}

	// Curr year files
	fmt.Println("Analyzing this year's repo...")

	currYearErr = r.checkoutRepoToCommitOrBranchName(r.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo back to the latest commit")
		panic(currYearErr)
	}

	currYearFiles := r.getRepoFiles(r.MasterBranchName)
	SaveDataToFile(currYearFiles, r.GetCurrYearFileListFile())

	if r.IncludeFileBlames {
		currYearBlames := GetFileBlameSummary(r, currYearFiles)
		SaveDataToFile(currYearBlames, r.GetCurrYearFileBlamesFile())
	}
}

func calculateRecap(r RepoConfig) {
	s := GetSpinner()

	fmt.Println()
	s.Suffix = " Calculating repo stats..."
	s.Start()

	isMultiYearRepo := r.GetIsMultiYearRepo()
	numCommitsAllTime := GetNumCommitsAllTime(r)
	numCommitsPrevYear := GetNumCommitsPrevYear(r)
	numCommitsCurrYear := GetNumCommitsCurrYear(r)
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
	commitsByMonthCurrYear := GetCommitsByMonthCurrYear(r)
	commitsByWeekDayCurrYear := GetCommitsByWeekDayCurrYear(r)
	commitsByHourCurrYear := GetCommitsByHourCurrYear(r)
	mostSingleDayCommitsByAuthorCurrYear := GetMostCommitsByAuthorCurrYear()
	mostInsertionsInCommitCurrYear := GetMostInsertionsInCommitCurrYear(r)
	mostDeletionsInCommitCurrYear := GetMostDeletionsInCommitCurrYear(r)
	largestCommitMessageCurrYear := GetLargestCommitMessageCurrYear(r)
	smallestCommitMessagesCurrYear := GetSmallestCommitMessagesCurrYear(r)
	commitMessageHistogramCurrYear := GetCommitMessageHistogramCurrYear(r)
	directPushesOnMasterByAuthorCurrYear := GetDirectPushesOnMasterByAuthorCurrYear()
	mergesToMasterByAuthorCurrYear := GetMergesToMasterByAuthorCurrYear()
	mostMergesInOneDayCurrYear := GetMostMergesInOneDayCurrYear()
	avgMergesToMasterPerDayCurrYear := GetAvgMergesToMasterPerDayCurrYear()
	fileChangesByAuthorCurrYear := GetFileChangesByAuthorCurrYear(r)
	codeInsertionsByAuthorCurrYear := GetCodeInsertionsByAuthorCurrYear(r)
	codeDeletionsByAuthorCurrYear := GetCodeDeletionsByAuthorCurrYear(r)
	fileChangeRatioCurrYear := GetFileChangeRatio(codeInsertionsByAuthorCurrYear, codeDeletionsByAuthorCurrYear)
	commonlyChangedFiles := GetCommonlyChangedFiles(r)
	fileCountPrevYear := GetFileCountPrevYear()
	fileCountCurrYear := GetFileCountCurrYear()
	largestFilesCurrYear := GetLargestFilesCurrYear()
	smallestFilesCurrYear := GetSmallestFilesCurrYear()
	totalLinesOfCodePrevYear := GetTotalLinesOfCodePrevYear()
	totalLinesOfCodeCurrYear := GetTotalLinesOfCodeCurrYear()
	totalLinesOfCodeInRepoByAuthor := GetTotalLinesOfCodeInRepoByAuthor()
	sizeOfRepoByWeekCurrYear := GetSizeOfRepoByWeekCurrYear(r)

	now := time.Now()
	isoDateString := now.Format(time.RFC3339)

	fileCountPercentDifference := (float64(fileCountCurrYear) - float64(fileCountPrevYear)) / float64(fileCountPrevYear)
	if math.IsNaN(fileCountPercentDifference) {
		panic("File count percent difference is NaN!")
	}

	repoRecap := Recap{
		// Metadata
		Version:            "0.0.1",
		Name:               r.Path,
		DateAnalyzed:       isoDateString,
		IsMultiYearRepo:    isMultiYearRepo,
		IncludesFileBlames: r.IncludeFileBlames,

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
	if !IsFileReadable(path) {
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
