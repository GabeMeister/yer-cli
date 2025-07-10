package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Recap struct {
	Version            string `json:"version"`
	Name               string `json:"name"`
	Directory          string `json:"directory"`
	DateAnalyzed       string `json:"date_analyzed"`
	IsMultiYearRepo    bool   `json:"is_multi_year_repo"`
	IncludesFileBlames bool   `json:"includes_file_blames"`

	// Commits
	NumCommitsAllTime               int                            `json:"num_commits_all_time"`
	NumCommitsPrevYear              int                            `json:"num_commits_prev_year"`
	NumCommitsCurrYear              int                            `json:"num_commits_curr_year"`
	CommitsByMonthCurrYear          []CommitMonth                  `json:"commits_by_month_curr_year"`
	CommitsByWeekDayCurrYear        []CommitWeekDay                `json:"commits_by_week_day_curr_year"`
	CommitsByHourCurrYear           []CommitHour                   `json:"commits_by_hour_curr_year"`
	MostInsertionsInCommitCurrYear  GitCommit                      `json:"most_insertions_in_commit_curr_year"`
	MostDeletionsInCommitCurrYear   GitCommit                      `json:"most_deletions_in_commit__curr_year"`
	LargestCommitMessageCurrYear    GitCommit                      `json:"largest_commit_message_curr_year"`
	SmallestCommitMessagesCurrYear  []GitCommit                    `json:"smallest_commit_messages_curr_year"`
	CommitMessageHistogramCurrYear  []CommitMessageLengthFrequency `json:"commit_message_histogram_curr_year"`
	MostMergesInOneDayCurrYear      MostMergesInOneDay             `json:"most_merges_in_one_day_curr_year"`
	AvgMergesToMasterPerDayCurrYear float64                        `json:"avg_merges_to_master_per_day_curr_year"`
	CommonlyChangedFiles            []FileChangeCount              `json:"commonly_changed_files"`

	// Files
	FileCountPrevYear          int                 `json:"file_count_prev_year"`
	FileCountCurrYear          int                 `json:"file_count_curr_year"`
	FileCountPercentDifference float64             `json:"file_count_percent_difference"`
	LargestFilesCurrYear       []FileSize          `json:"largest_files_curr_year"`
	SmallestFilesCurrYear      []FileSize          `json:"smallest_files_curr_year"`
	TotalLinesOfCodePrevYear   int                 `json:"total_lines_of_code_prev_year"`
	TotalLinesOfCodeCurrYear   int                 `json:"total_lines_of_code_curr_year"`
	SizeOfRepoByWeekCurrYear   []RepoSizeTimeStamp `json:"size_of_repo_by_week_curr_year"`

	// Team
	AllAuthors                           []string                     `json:"all_authors"`
	NewAuthorCommitsCurrYear             []GitCommit                  `json:"new_author_commits_curr_year"`
	NewAuthorCountCurrYear               int                          `json:"new_author_count_curr_year"`
	NewAuthorListCurrYear                []string                     `json:"new_author_list_curr_year"`
	AuthorCommitCountsCurrYear           map[string]int               `json:"author_commit_counts_curr_year"`
	AuthorCommitCountsAllTime            map[string]int               `json:"author_commit_counts_all_time"`
	AuthorCountCurrYear                  int                          `json:"author_count_curr_year"`
	AuthorCountAllTime                   int                          `json:"author_count_all_time"`
	AuthorTotalFileChangesPrevYear       map[string]int               `json:"author_total_file_changes_prev_year"`
	AuthorFileChangesOverTimeCurrYear    TotalFileChangeCount         `json:"author_file_changes_over_time_curr_year"`
	MostSingleDayCommitsByAuthorCurrYear MostSingleDayCommitsByAuthor `json:"most_single_day_commits_by_author_curr_year"`
	DirectPushesOnMasterByAuthorCurrYear map[string]int               `json:"direct_pushes_on_master_by_author_curr_year"`
	MergesToMasterByAuthorCurrYear       map[string]int               `json:"merges_to_master_by_author_curr_year"`
	FileChangesByAuthorCurrYear          map[string]int               `json:"file_changes_by_author_curr_year"`
	FileChangeRatioByAuthorCurrYear      map[string]float64           `json:"file_change_ratio_by_author_curr_year"`
	TotalLinesOfCodeInRepoByAuthor       map[string]int               `json:"total_lines_of_code_in_repo_by_author"`
}

type MultiRepoRecap struct {
	Version                  string   `json:"version"`
	Name                     string   `json:"name"`
	DateAnalyzed             string   `json:"date_analyzed"`
	RepoNames                []string `json:"repo_names"`
	ActiveAuthorsCountByRepo map[Repo]int
	FileCountByRepoCurrYear  map[Repo]int
	TotalLinesOfCodeByRepo   map[Repo]int
}

type Repo string
type Author string
type AuthorList []string
type NewAuthorByRepo map[Repo]AuthorList

func GetRepoRecapFromTmpDir() (Recap, error) {
	if !HasRecapBeenRan() {
		return Recap{}, os.ErrNotExist
	}

	files, err := os.ReadDir("tmp")
	if err != nil {
		fmt.Println("Unable to read tmp directory to get repo recap", err)
		os.Exit(1)
	}

	var recapFile string

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "_recap.json") {
			recapFile = file.Name()
		}
	}

	// Set to temporary recap files
	// recapFile = "demo-stack_recap.json"

	data, err := os.ReadFile(fmt.Sprintf("./tmp/%s", recapFile))
	if err != nil {
		panic(err)
	}
	var repoRecap Recap

	json.Unmarshal(data, &repoRecap)

	return repoRecap, nil
}

func calculateRepoRecap(r *RepoConfig) {
	s := getSpinner()

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
	saveDataToFile(repoRecap, repoRecapFile)

	s.Stop()
}

func getAllRecaps(c *ConfigFile) ([]Recap, error) {
	recaps := []Recap{}
	valid := true

	// Verify that all recap files exist, according to what's in the config
	for _, r := range c.Repos {
		if !r.hasRecapFile() {
			valid = false
			break
		} else {
			recap, err := r.getRepoRecap()
			if err != nil {
				return []Recap{}, err
			}

			recaps = append(recaps, recap)
		}
	}

	if !valid {
		return []Recap{}, os.ErrInvalid
	}

	return recaps, nil
}

func calculateMultiRepoRecap(c *ConfigFile) error {
	recaps, err := getAllRecaps(c)
	if err != nil {
		return err
	}

	now := time.Now()

	// Combine all metrics from the separate recaps
	repoNames := getRepoNames(recaps)
	activeAuthorsCountByRepo := getActiveAuthorsCountByRepo(recaps)
	fileCountByRepoCurrYear := getFileCountByRepoCurrYear(recaps)
	totalLinesOfCodeByRepo := getTotalLinesOfCodeByRepo(recaps)

	// Combine stats
	multiRepoRecap := MultiRepoRecap{
		DateAnalyzed:             now.Format(time.RFC3339),
		Name:                     c.Name,
		RepoNames:                repoNames,
		ActiveAuthorsCountByRepo: activeAuthorsCountByRepo,
		FileCountByRepoCurrYear:  fileCountByRepoCurrYear,
		TotalLinesOfCodeByRepo:   totalLinesOfCodeByRepo,
	}

	saveDataToFile(multiRepoRecap, MULTI_REPO_RECAP_FILE)

	return nil
}

/*
 * MULTI REPO METRICS
 */

func getRepoNames(recaps []Recap) []string {
	repoNames := []string{}
	for _, recap := range recaps {
		repoNames = append(repoNames, recap.Name)
	}

	return repoNames
}

func getActiveAuthorsCountByRepo(recaps []Recap) map[Repo]int {
	activeAuthorsMap := make(map[Repo]int)

	for _, recap := range recaps {
		activeAuthorsMap[Repo(recap.Name)] = recap.AuthorCountCurrYear
	}

	return activeAuthorsMap
}

func getFileCountByRepoCurrYear(recaps []Recap) map[Repo]int {
	fileCountMap := make(map[Repo]int)

	for _, recap := range recaps {
		fileCountMap[Repo(recap.Name)] = recap.FileCountCurrYear
	}

	return fileCountMap
}

func getTotalLinesOfCodeByRepo(recaps []Recap) map[Repo]int {
	fileCountMap := make(map[Repo]int)

	for _, recap := range recaps {
		fileCountMap[Repo(recap.Name)] = recap.TotalLinesOfCodeCurrYear
	}

	return fileCountMap
}
