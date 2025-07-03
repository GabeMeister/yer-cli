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
	"strings"
	"time"
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
			gatherMetrics(&r)
			config.updateDuplicateAuthors(&r)
		}
		r.calculateRecap()
	}

	err := config.CalculateMultiRepoRecap()
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

func gatherMetrics(r *RepoConfig) {
	stashRepo(r.Path)

	currYearErr := checkoutRepoToCommitOrBranchName(r.Path, r.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo to the latest commit")
		panic(currYearErr)
	}

	// We want the latest changes
	pullRepo(r.Path)

	commits := getCommitsFromGitLogs(r, false)
	commitsFileName := r.getCommitsFile()
	SaveDataToFile(commits, commitsFileName)

	mergeCommits := getCommitsFromGitLogs(r, true)
	mergeCommitsFileName := r.getMergeCommitsFile()
	SaveDataToFile(mergeCommits, mergeCommitsFileName)

	directPushToMasterCommits := getDirectPushToMasterCommitsCurrYear(r)
	directPushFileName := r.getDirectPushesFile()
	SaveDataToFile(directPushToMasterCommits, directPushFileName)

	// Prev year files (if possible)
	if r.hasPrevYearCommits() {
		lastCommitPrevYear := r.getLastCommitPrevYear()
		fmt.Printf("Analyzing last year's repo for %s...\n", r.getName())
		prevYearErr := checkoutRepoToCommitOrBranchName(r.Path, lastCommitPrevYear.Commit)
		if prevYearErr != nil {
			fmt.Println("Unable to git checkout repo to last year's files")
			panic(prevYearErr)
		}

		prevYearFiles := getRepoFiles(r, lastCommitPrevYear.Commit)
		SaveDataToFile(prevYearFiles, r.getPrevYearFileListFile())

		if r.AnalyzeFileBlames {
			prevYearBlames := getFileBlameSummary(r, prevYearFiles)
			SaveDataToFile(prevYearBlames, r.getPrevYearFileBlamesFile())
		}
	}

	// Curr year files
	fmt.Printf("Analyzing this year's repo for %s...\n", r.getName())

	currYearErr = checkoutRepoToCommitOrBranchName(r.Path, r.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo back to the latest commit")
		panic(currYearErr)
	}

	currYearFiles := getRepoFiles(r, r.MasterBranchName)
	SaveDataToFile(currYearFiles, r.getCurrYearFileListFile())

	if r.AnalyzeFileBlames {
		currYearBlames := getFileBlameSummary(r, currYearFiles)
		SaveDataToFile(currYearBlames, r.getCurrYearFileBlamesFile())
	}
}

func (r *RepoConfig) calculateRecap() {
	s := GetSpinner()

	fmt.Println()
	utils.PrintProgress(s, fmt.Sprintf("Calculating repo stats for %s...", r.getName()))

	if !utils.IsDevMode() {
		s.Start()
	}

	isMultiYearRepo := r.getIsMultiYearRepo()
	numCommitsAllTime := r.getNumCommitsAllTime()
	numCommitsPrevYear := r.getNumCommitsPrevYear()
	numCommitsCurrYear := r.getNumCommitsCurrYear()
	allAuthors := r.getAllAuthorsList()
	newAuthorCommitsCurrYear := r.getNewAuthorCommitsCurrYear()
	newAuthorCountCurrYear := len(newAuthorCommitsCurrYear)
	newAuthorListCurrYear := utils.Map(newAuthorCommitsCurrYear, func(commit GitCommit) string {
		return commit.Author
	})
	authorCommitCountsCurrYear := r.getAuthorCommitCountCurrYear()
	authorCommitCountsAllTime := r.getAuthorCommitCountAllTime()
	authorCountCurrYear := r.getAuthorCountCurrYear()
	authorCountAllTime := r.getAuthorCountAllTime()
	authorTotalFileChangesPrevYear := r.getAuthorTotalFileChangesPrevYear()
	authorFileChangesOverTimeCurrYear := r.getAuthorFileChangesOverTimeCurrYear()
	commitsByMonthCurrYear := r.getCommitsByMonthCurrYear()
	commitsByWeekDayCurrYear := r.getCommitsByWeekDayCurrYear()
	commitsByHourCurrYear := r.getCommitsByHourCurrYear()
	mostSingleDayCommitsByAuthorCurrYear := r.getMostCommitsByAuthorCurrYear()
	mostInsertionsInCommitCurrYear := r.getMostInsertionsInCommitCurrYear()
	mostDeletionsInCommitCurrYear := r.getMostDeletionsInCommitCurrYear()
	largestCommitMessageCurrYear := r.getLargestCommitMessageCurrYear()
	smallestCommitMessagesCurrYear := r.getSmallestCommitMessagesCurrYear()
	commitMessageHistogramCurrYear := r.getCommitMessageHistogramCurrYear()
	directPushesOnMasterByAuthorCurrYear := r.getDirectPushesOnMasterByAuthorCurrYear()
	mergesToMasterByAuthorCurrYear := r.getMergesToMasterByAuthorCurrYear()
	mostMergesInOneDayCurrYear := r.getMostMergesInOneDayCurrYear()
	avgMergesToMasterPerDayCurrYear := r.getAvgMergesToMasterPerDayCurrYear()
	fileChangesByAuthorCurrYear := r.getFileChangesByAuthorCurrYear()
	codeInsertionsByAuthorCurrYear := r.getCodeInsertionsByAuthorCurrYear()
	codeDeletionsByAuthorCurrYear := r.getCodeDeletionsByAuthorCurrYear()
	fileChangeRatioCurrYear := r.getFileChangeRatio(codeInsertionsByAuthorCurrYear, codeDeletionsByAuthorCurrYear)
	commonlyChangedFiles := r.getCommonlyChangedFiles()
	fileCountPrevYear := r.getFileCountPrevYear()
	fileCountCurrYear := r.getFileCountCurrYear()
	largestFilesCurrYear := r.getLargestFilesCurrYear()
	smallestFilesCurrYear := r.getSmallestFilesCurrYear()
	totalLinesOfCodePrevYear := r.getTotalLinesOfCodePrevYear()
	totalLinesOfCodeCurrYear := r.getTotalLinesOfCodeCurrYear()
	totalLinesOfCodeInRepoByAuthor := r.getTotalLinesOfCodeInRepoByAuthor()
	sizeOfRepoByWeekCurrYear := r.getSizeOfRepoByWeekCurrYear()

	now := time.Now()
	isoDateString := now.Format(time.RFC3339)

	var fileCountPercentDifference float64
	if fileCountPrevYear != 0 {
		fileCountPercentDifference = (float64(fileCountCurrYear) - float64(fileCountPrevYear)) / float64(fileCountPrevYear)
	}
	if math.IsNaN(fileCountPercentDifference) {
		panic("File count percent difference is NaN!")
	}

	repoRecap := Recap{
		// Metadata
		Version:            "0.0.1",
		Name:               filepath.Base(r.Path),
		Directory:          r.Path,
		DateAnalyzed:       isoDateString,
		IsMultiYearRepo:    isMultiYearRepo,
		IncludesFileBlames: r.AnalyzeFileBlames,

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
		AllAuthors:                           allAuthors,
		NewAuthorCommitsCurrYear:             newAuthorCommitsCurrYear,
		NewAuthorCountCurrYear:               newAuthorCountCurrYear,
		NewAuthorListCurrYear:                newAuthorListCurrYear,
		AuthorCommitCountsCurrYear:           authorCommitCountsCurrYear,
		AuthorCommitCountsAllTime:            authorCommitCountsAllTime,
		AuthorCountCurrYear:                  authorCountCurrYear,
		AuthorCountAllTime:                   authorCountAllTime,
		AuthorTotalFileChangesPrevYear:       authorTotalFileChangesPrevYear,
		AuthorFileChangesOverTimeCurrYear:    authorFileChangesOverTimeCurrYear,
		MostSingleDayCommitsByAuthorCurrYear: mostSingleDayCommitsByAuthorCurrYear,
		DirectPushesOnMasterByAuthorCurrYear: directPushesOnMasterByAuthorCurrYear,
		MergesToMasterByAuthorCurrYear:       mergesToMasterByAuthorCurrYear,
		FileChangesByAuthorCurrYear:          fileChangesByAuthorCurrYear,
		FileChangeRatioByAuthorCurrYear:      fileChangeRatioCurrYear,
		TotalLinesOfCodeInRepoByAuthor:       totalLinesOfCodeInRepoByAuthor,
	}

	repoRecapFile := r.getRecapFilePath()
	SaveDataToFile(repoRecap, repoRecapFile)

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
