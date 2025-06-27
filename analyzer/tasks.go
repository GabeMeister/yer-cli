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
		calculateRecap(&r)
	}

	err := config.CalculateMultiRepoRecap()

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

	currYearErr := r.checkoutRepoToCommitOrBranchName(r.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo to the latest commit")
		panic(currYearErr)
	}

	// We want the latest changes
	pullRepo(r.Path)

	commits := r.getCommitsFromGitLogs(false)
	commitsFileName := r.GetCommitsFile()
	SaveDataToFile(commits, commitsFileName)

	mergeCommits := r.getCommitsFromGitLogs(true)
	mergeCommitsFileName := r.GetMergeCommitsFile()
	SaveDataToFile(mergeCommits, mergeCommitsFileName)

	directPushToMasterCommits := r.getDirectPushToMasterCommitsCurrYear()
	directPushFileName := r.GetDirectPushesFile()
	SaveDataToFile(directPushToMasterCommits, directPushFileName)

	// Prev year files (if possible)
	if r.hasPrevYearCommits() {
		lastCommitPrevYear := r.getLastCommitPrevYear()
		fmt.Printf("Analyzing last year's repo for %s...\n", r.GetName())
		prevYearErr := r.checkoutRepoToCommitOrBranchName(lastCommitPrevYear.Commit)
		if prevYearErr != nil {
			fmt.Println("Unable to git checkout repo to last year's files")
			panic(prevYearErr)
		}

		prevYearFiles := r.getRepoFiles(lastCommitPrevYear.Commit)
		SaveDataToFile(prevYearFiles, r.GetPrevYearFileListFile())

		if r.AnalyzeFileBlames {
			prevYearBlames := r.GetFileBlameSummary(prevYearFiles)
			SaveDataToFile(prevYearBlames, r.GetPrevYearFileBlamesFile())
		}
	}

	// Curr year files
	fmt.Printf("Analyzing this year's repo for %s...\n", r.GetName())

	currYearErr = r.checkoutRepoToCommitOrBranchName(r.MasterBranchName)
	if currYearErr != nil {
		fmt.Println("Unable to git checkout repo back to the latest commit")
		panic(currYearErr)
	}

	currYearFiles := r.getRepoFiles(r.MasterBranchName)
	SaveDataToFile(currYearFiles, r.GetCurrYearFileListFile())

	if r.AnalyzeFileBlames {
		currYearBlames := r.GetFileBlameSummary(currYearFiles)
		SaveDataToFile(currYearBlames, r.GetCurrYearFileBlamesFile())
	}
}

func calculateRecap(r *RepoConfig) {
	s := GetSpinner()

	fmt.Println()
	utils.PrintProgress(s, fmt.Sprintf("Calculating repo stats for %s...", r.GetName()))

	if !utils.IsDevMode() {
		s.Start()
	}

	isMultiYearRepo := r.GetIsMultiYearRepo()
	numCommitsAllTime := r.GetNumCommitsAllTime()
	numCommitsPrevYear := r.GetNumCommitsPrevYear()
	numCommitsCurrYear := r.GetNumCommitsCurrYear()
	allAuthors := r.GetAllAuthorsList()
	newAuthorCommitsCurrYear := r.GetNewAuthorCommitsCurrYear()
	newAuthorCountCurrYear := len(newAuthorCommitsCurrYear)
	newAuthorListCurrYear := utils.Map(newAuthorCommitsCurrYear, func(commit GitCommit) string {
		return commit.Author
	})
	authorCommitCountsCurrYear := r.GetAuthorCommitCountCurrYear()
	authorCommitCountsAllTime := r.GetAuthorCommitCountAllTime()
	authorCountCurrYear := r.GetAuthorCountCurrYear()
	authorCountAllTime := r.GetAuthorCountAllTime()
	authorTotalFileChangesPrevYear := r.GetAuthorTotalFileChangesPrevYear()
	authorFileChangesOverTimeCurrYear := r.GetAuthorFileChangesOverTimeCurrYear()
	commitsByMonthCurrYear := r.GetCommitsByMonthCurrYear()
	commitsByWeekDayCurrYear := r.GetCommitsByWeekDayCurrYear()
	commitsByHourCurrYear := r.GetCommitsByHourCurrYear()
	mostSingleDayCommitsByAuthorCurrYear := r.GetMostCommitsByAuthorCurrYear()
	mostInsertionsInCommitCurrYear := r.GetMostInsertionsInCommitCurrYear()
	mostDeletionsInCommitCurrYear := r.GetMostDeletionsInCommitCurrYear()
	largestCommitMessageCurrYear := r.GetLargestCommitMessageCurrYear()
	smallestCommitMessagesCurrYear := r.GetSmallestCommitMessagesCurrYear()
	commitMessageHistogramCurrYear := r.GetCommitMessageHistogramCurrYear()
	directPushesOnMasterByAuthorCurrYear := r.GetDirectPushesOnMasterByAuthorCurrYear()
	mergesToMasterByAuthorCurrYear := r.GetMergesToMasterByAuthorCurrYear()
	mostMergesInOneDayCurrYear := r.GetMostMergesInOneDayCurrYear()
	avgMergesToMasterPerDayCurrYear := r.GetAvgMergesToMasterPerDayCurrYear()
	fileChangesByAuthorCurrYear := r.GetFileChangesByAuthorCurrYear()
	codeInsertionsByAuthorCurrYear := r.GetCodeInsertionsByAuthorCurrYear()
	codeDeletionsByAuthorCurrYear := r.GetCodeDeletionsByAuthorCurrYear()
	fileChangeRatioCurrYear := r.GetFileChangeRatio(codeInsertionsByAuthorCurrYear, codeDeletionsByAuthorCurrYear)
	commonlyChangedFiles := r.GetCommonlyChangedFiles()
	fileCountPrevYear := r.GetFileCountPrevYear()
	fileCountCurrYear := r.GetFileCountCurrYear()
	largestFilesCurrYear := r.GetLargestFilesCurrYear()
	smallestFilesCurrYear := r.GetSmallestFilesCurrYear()
	totalLinesOfCodePrevYear := r.GetTotalLinesOfCodePrevYear()
	totalLinesOfCodeCurrYear := r.GetTotalLinesOfCodeCurrYear()
	totalLinesOfCodeInRepoByAuthor := r.GetTotalLinesOfCodeInRepoByAuthor()
	sizeOfRepoByWeekCurrYear := r.GetSizeOfRepoByWeekCurrYear()

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

	repoRecapFile := r.GetRecapFile()
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
