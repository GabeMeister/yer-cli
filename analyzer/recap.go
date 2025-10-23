package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
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
	MergeCommitsByMonthCurrYear     []CommitMonth                  `json:"merge_commits_by_month_curr_year"`
	MergeCommitsByWeekDayCurrYear   []CommitWeekDay                `json:"merge_commits_by_week_day_curr_year"`
	MergeCommitsByHourCurrYear      []CommitHour                   `json:"merge_commits_by_hour_curr_year"`
	MostInsertionsInCommitCurrYear  GitCommit                      `json:"most_insertions_in_commit_curr_year"`
	MostDeletionsInCommitCurrYear   GitCommit                      `json:"most_deletions_in_commit__curr_year"`
	LargestCommitMessageCurrYear    GitCommit                      `json:"largest_commit_message_curr_year"`
	SmallestCommitMessagesCurrYear  []GitCommit                    `json:"smallest_commit_messages_curr_year"`
	CommitMessageHistogramCurrYear  []CommitMessageLengthFrequency `json:"commit_message_histogram_curr_year"`
	MostMergesInOneDayCurrYear      MostMergesInOneDay             `json:"most_merges_in_one_day_curr_year"`
	AvgMergesToMasterPerDayCurrYear float64                        `json:"avg_merges_to_master_per_day_curr_year"`
	MergesToMasterCurrYear          int                            `json:"merges_to_master_curr_year"`
	MergesToMasterPrevYear          int                            `json:"merges_to_master_prev_year"`
	MergesToMasterAllTime           int                            `json:"merges_to_master_all_time"`
	CommonlyChangedFiles            []FileChangeCount              `json:"commonly_changed_files"`

	// Files
	FileCountPrevYear                 int              `json:"file_count_prev_year"`
	FileCountCurrYear                 int              `json:"file_count_curr_year"`
	FileCountPercentDifference        float64          `json:"file_count_percent_difference"`
	LargestFilesCurrYear              []FileSize       `json:"largest_files_curr_year"`
	SmallestFilesCurrYear             []FileSize       `json:"smallest_files_curr_year"`
	TotalLinesOfCodePrevYear          int              `json:"total_lines_of_code_prev_year"`
	TotalLinesOfCodeByFileExtPrevYear FileExtLineCount `json:"total_lines_of_code_by_file_ext_prev_year"`
	TotalLinesOfCodeCurrYear          int              `json:"total_lines_of_code_curr_year"`
	TotalLinesOfCodeByFileExtCurrYear FileExtLineCount `json:"total_lines_of_code_by_file_ext_curr_year"`
	SizeOfRepoByWeekCurrYear          []int            `json:"size_of_repo_by_week_curr_year"`

	// Team
	AllAuthors                           []string                      `json:"all_authors"`
	NewAuthorCommitsCurrYear             []GitCommit                   `json:"new_author_commits_curr_year"`
	NewAuthorCountCurrYear               int                           `json:"new_author_count_curr_year"`
	NewAuthorListCurrYear                []string                      `json:"new_author_list_curr_year"`
	AuthorCommitCountsCurrYear           map[string]int                `json:"author_commit_counts_curr_year"`
	AuthorCommitCountsPrevYear           map[string]int                `json:"author_commit_counts_prev_year"`
	AuthorCommitCountsAllTime            map[string]int                `json:"author_commit_counts_all_time"`
	AuthorCountCurrYear                  int                           `json:"author_count_curr_year"`
	AuthorCountPrevYear                  int                           `json:"author_count_prev_year"`
	AuthorCountAllTime                   int                           `json:"author_count_all_time"`
	AuthorTotalFileChangesPrevYear       map[string]int                `json:"author_total_file_changes_prev_year"`
	AuthorFileChangesOverTimeCurrYear    TotalFileChangeCount          `json:"author_file_changes_over_time_curr_year"`
	MostSingleDayCommitsByAuthorCurrYear MostSingleDayCommitsByAuthor  `json:"most_single_day_commits_by_author_curr_year"`
	DirectPushesOnMasterByAuthorCurrYear map[string]int                `json:"direct_pushes_on_master_by_author_curr_year"`
	MergesToMasterByAuthorCurrYear       map[string]int                `json:"merges_to_master_by_author_curr_year"`
	FileChangesByAuthorCurrYear          map[string]FileChangesSummary `json:"file_changes_by_author_curr_year"`
	FileChangeRatioByAuthorCurrYear      map[string]float64            `json:"file_change_ratio_by_author_curr_year"`
	TotalLinesOfCodeInRepoByAuthor       map[string]int                `json:"total_lines_of_code_in_repo_by_author"`
}

type MultiRepoRecap struct {
	Version                     string                     `json:"version"`
	Name                        string                     `json:"name"`
	DateAnalyzed                string                     `json:"date_analyzed"`
	RepoNames                   []string                   `json:"repo_names"`
	ActiveAuthorsCountByRepo    map[Repo]YearComparison    `json:"active_authors_count_by_repo"`
	FileCountByRepo             map[Repo]YearComparison    `json:"file_count_by_repo"`
	TotalLinesOfCodeByRepo      map[Repo]YearComparison    `json:"total_lines_of_code_by_repo"`
	SizeOfRepoWeeklyByRepo      map[Repo][]int             `json:"size_of_repo_weekly_by_repo"`
	CommitsMadeByRepo           map[Repo]YearComparison    `json:"commits_made_by_repo"`
	CommitsMadeByAuthor         map[Author]*YearComparison `json:"commits_made_by_author"`
	FileChangesMadeByAuthor     []AuthorFileChangesSummary `json:"file_changes_made_by_author"`
	LinesOfCodeOwnedByAuthor    map[Author]int             `json:"lines_of_code_owned_by_author"`
	MergeCommitsByMonth         []int                      `json:"aggregate_commits_by_month"`
	MergeCommitsByWeekDay       []int                      `json:"aggregate_commits_by_week_day"`
	MergeCommitsByHour          []int                      `json:"aggregate_commits_by_hour"`
	AvgMergesPerDayByRepo       map[Repo]float64           `json:"avg_merges_per_day_by_repo"`
	MergesToMasterByRepo        map[Repo]YearComparison    `json:"merges_to_master_by_repo"`
	MergesToMasterAllTimeByRepo map[Repo]int               `json:"merges_to_master_all_time_by_repo"`
}

type Repo string
type Author string
type AuthorList []string
type NewAuthorByRepo map[Repo]AuthorList
type Year string

const (
	PREV Year = "prev"
	CURR Year = "curr"
)

type YearComparison struct {
	Prev int `json:"prev"`
	Curr int `json:"curr"`
}

type FileChangesSummary struct {
	Insertions int `json:"insertions"`
	Deletions  int `json:"deletions"`
}

type AuthorFileChangesSummary struct {
	Author     Author `json:"author"`
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
}

func GetMultiRepoRecapFromTmpDir() (MultiRepoRecap, error) {
	if !HasRecapBeenRan() {
		return MultiRepoRecap{}, os.ErrNotExist
	}

	filePath := "./tmp/multi_repo_recap.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var repoRecap MultiRepoRecap

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
	authorCommitCountsPrevYear := r.getAuthorCommitCountPrevYear()
	authorCommitCountsAllTime := r.getAuthorCommitCountAllTime()
	authorCountCurrYear := r.getAuthorCountCurrYear()
	authorCountPrevYear := r.getAuthorCountPrevYear()
	authorCountAllTime := r.getAuthorCountAllTime()
	authorTotalFileChangesPrevYear := r.getAuthorTotalFileChangesPrevYear()
	authorFileChangesOverTimeCurrYear := r.getAuthorFileChangesOverTimeCurrYear()
	commitsByMonthCurrYear := r.getCommitsByMonthCurrYear()
	commitsByWeekDayCurrYear := r.getCommitsByWeekDayCurrYear()
	commitsByHourCurrYear := r.getCommitsByHourCurrYear()
	mergeCommitsByMonthCurrYear := r.getMergeCommitsByMonthCurrYear()
	mergeCommitsByWeekDayCurrYear := r.getMergeCommitsByWeekDayCurrYear()
	mergeCommitsByHourCurrYear := r.getMergeCommitsByHourCurrYear()
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
	mergesToMasterPrevYear := r.getMergesToMasterPrevYear()
	mergesToMasterCurrYear := r.getMergesToMasterCurrYear()
	mergesToMasterAllTime := r.getMergesToMasterAllTime()
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
	totalLinesOfCodeByFileExtPrevYear := r.getTotalLinesOfCodeByFileExtPrevYear()
	totalLinesOfCodeCurrYear := r.getTotalLinesOfCodeCurrYear()
	totalLinesOfCodeByFileExtCurrYear := r.getTotalLinesOfCodeByFileExtCurrYear()
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
		MergeCommitsByMonthCurrYear:     mergeCommitsByMonthCurrYear,
		MergeCommitsByWeekDayCurrYear:   mergeCommitsByWeekDayCurrYear,
		MergeCommitsByHourCurrYear:      mergeCommitsByHourCurrYear,
		MostInsertionsInCommitCurrYear:  mostInsertionsInCommitCurrYear,
		MostDeletionsInCommitCurrYear:   mostDeletionsInCommitCurrYear,
		LargestCommitMessageCurrYear:    largestCommitMessageCurrYear,
		SmallestCommitMessagesCurrYear:  smallestCommitMessagesCurrYear,
		CommitMessageHistogramCurrYear:  commitMessageHistogramCurrYear,
		MostMergesInOneDayCurrYear:      mostMergesInOneDayCurrYear,
		AvgMergesToMasterPerDayCurrYear: avgMergesToMasterPerDayCurrYear,
		MergesToMasterPrevYear:          mergesToMasterPrevYear,
		MergesToMasterCurrYear:          mergesToMasterCurrYear,
		MergesToMasterAllTime:           mergesToMasterAllTime,
		CommonlyChangedFiles:            commonlyChangedFiles,

		// Files
		FileCountPrevYear:                 fileCountPrevYear,
		FileCountCurrYear:                 fileCountCurrYear,
		FileCountPercentDifference:        fileCountPercentDifference,
		LargestFilesCurrYear:              largestFilesCurrYear,
		SmallestFilesCurrYear:             smallestFilesCurrYear,
		TotalLinesOfCodePrevYear:          totalLinesOfCodePrevYear,
		TotalLinesOfCodeByFileExtPrevYear: totalLinesOfCodeByFileExtPrevYear,
		TotalLinesOfCodeCurrYear:          totalLinesOfCodeCurrYear,
		TotalLinesOfCodeByFileExtCurrYear: totalLinesOfCodeByFileExtCurrYear,
		SizeOfRepoByWeekCurrYear:          sizeOfRepoByWeekCurrYear,

		// Team
		AllAuthors:                           allAuthors,
		NewAuthorCommitsCurrYear:             newAuthorCommitsCurrYear,
		NewAuthorCountCurrYear:               newAuthorCountCurrYear,
		NewAuthorListCurrYear:                newAuthorListCurrYear,
		AuthorCommitCountsCurrYear:           authorCommitCountsCurrYear,
		AuthorCommitCountsPrevYear:           authorCommitCountsPrevYear,
		AuthorCommitCountsAllTime:            authorCommitCountsAllTime,
		AuthorCountCurrYear:                  authorCountCurrYear,
		AuthorCountPrevYear:                  authorCountPrevYear,
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
	fileCountByRepoCurrYear := getFileCountByRepo(recaps)
	totalLinesOfCodeByRepo := getTotalLinesOfCodeByRepo(recaps)
	sizeOfRepoWeeklyByRepo := getSizeOfRepoWeeklyByRepo(recaps)
	commitsMadeByRepo := getCommitsMadeByRepo(recaps)
	commitsMadeByAuthor := getCommitsMadeByAuthor(recaps)
	fileChangesMadeByAuthor := getFileChangesMadeByAuthor(recaps)
	linesOfCodeOwnedByAuthor := getLinesOfCodeOwnedByAuthor(recaps)
	mergeCommitsByMonth := getMergeCommitsByMonth(recaps)
	mergeCommitsByWeekDay := getMergeCommitsByWeekDay(recaps)
	mergeCommitsByHour := getMergeCommitsByHour(recaps)
	avgMergesPerDayByRepo := getAvgMergesPerDayByRepo(recaps)
	mergesToMasterByRepo := getMergesToMasterByRepo(recaps)
	mergesToMasterAllTimeByRepo := getMergesToMasterAllTimeByRepo(recaps)

	// Combine stats
	multiRepoRecap := MultiRepoRecap{
		DateAnalyzed:                now.Format(time.RFC3339),
		Name:                        c.Name,
		RepoNames:                   repoNames,
		ActiveAuthorsCountByRepo:    activeAuthorsCountByRepo,
		FileCountByRepo:             fileCountByRepoCurrYear,
		TotalLinesOfCodeByRepo:      totalLinesOfCodeByRepo,
		SizeOfRepoWeeklyByRepo:      sizeOfRepoWeeklyByRepo,
		CommitsMadeByRepo:           commitsMadeByRepo,
		CommitsMadeByAuthor:         commitsMadeByAuthor,
		FileChangesMadeByAuthor:     fileChangesMadeByAuthor,
		LinesOfCodeOwnedByAuthor:    linesOfCodeOwnedByAuthor,
		MergeCommitsByMonth:         mergeCommitsByMonth,
		MergeCommitsByWeekDay:       mergeCommitsByWeekDay,
		MergeCommitsByHour:          mergeCommitsByHour,
		AvgMergesPerDayByRepo:       avgMergesPerDayByRepo,
		MergesToMasterByRepo:        mergesToMasterByRepo,
		MergesToMasterAllTimeByRepo: mergesToMasterAllTimeByRepo,
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

func getActiveAuthorsCountByRepo(recaps []Recap) map[Repo]YearComparison {
	activeAuthorsMap := make(map[Repo]YearComparison)

	for _, recap := range recaps {
		activeAuthorsMap[Repo(recap.Name)] = YearComparison{
			Prev: recap.AuthorCountPrevYear,
			Curr: recap.AuthorCountCurrYear,
		}
	}

	return activeAuthorsMap

}

func getFileCountByRepo(recaps []Recap) map[Repo]YearComparison {
	fileCountMap := make(map[Repo]YearComparison)

	for _, recap := range recaps {
		fileCountMap[Repo(recap.Name)] = YearComparison{
			Prev: recap.FileCountPrevYear,
			Curr: recap.FileCountCurrYear,
		}
	}

	return fileCountMap

}

func getTotalLinesOfCodeByRepo(recaps []Recap) map[Repo]YearComparison {
	totalLinesMap := make(map[Repo]YearComparison)

	for _, recap := range recaps {
		totalLinesMap[Repo(recap.Name)] = YearComparison{
			Prev: recap.TotalLinesOfCodePrevYear,
			Curr: recap.TotalLinesOfCodeCurrYear,
		}
	}

	return totalLinesMap
}

func getSizeOfRepoWeeklyByRepo(recaps []Recap) map[Repo][]int {
	sizeOfRepoMap := make(map[Repo][]int)

	for _, recap := range recaps {
		sizeOfRepoMap[Repo(recap.Name)] = recap.SizeOfRepoByWeekCurrYear
	}

	return sizeOfRepoMap
}

func getCommitsMadeByRepo(recaps []Recap) map[Repo]YearComparison {
	commitsMap := make(map[Repo]YearComparison)

	for _, recap := range recaps {
		commitsMap[Repo(recap.Name)] = YearComparison{
			Prev: recap.NumCommitsPrevYear,
			Curr: recap.NumCommitsCurrYear,
		}
	}

	return commitsMap
}

func getAuthorsFromRecaps(recaps []Recap) []Author {
	authorMap := make(map[Author]bool)

	for _, r := range recaps {
		for _, author := range r.AllAuthors {
			authorMap[Author(author)] = true

		}
	}

	return utils.MapKeysToSlice(authorMap)
}

func getCommitsMadeByAuthor(recaps []Recap) map[Author]*YearComparison {
	commitsMap := make(map[Author]*YearComparison)

	// Get list of all authors
	allAuthors := getAuthorsFromRecaps(recaps)

	for _, author := range allAuthors {
		commitsMap[author] = &YearComparison{
			Prev: 0,
			Curr: 0,
		}
	}

	for _, recap := range recaps {
		for author, commits := range recap.AuthorCommitCountsPrevYear {
			commitsMap[Author(author)].Prev += commits
		}

		for author, commits := range recap.AuthorCommitCountsCurrYear {
			commitsMap[Author(author)].Curr += commits
		}
	}

	return commitsMap
}

// Insertions

func getFileChangesMadeByAuthor(recaps []Recap) []AuthorFileChangesSummary {
	authorInsertionsMap := make(map[Author]int)
	authorDeletionsMap := make(map[Author]int)

	for _, recap := range recaps {
		for author, fileChangesSummary := range recap.FileChangesByAuthorCurrYear {
			authorInsertionsMap[Author(author)] += fileChangesSummary.Insertions
			authorDeletionsMap[Author(author)] += fileChangesSummary.Deletions
		}
	}

	authorFileChangesSummary := []AuthorFileChangesSummary{}

	for author := range authorInsertionsMap {
		authorFileChangesSummary = append(authorFileChangesSummary, AuthorFileChangesSummary{
			Author:     author,
			Insertions: authorInsertionsMap[author],
			Deletions:  authorDeletionsMap[author],
		})
	}

	return authorFileChangesSummary
}

func getLinesOfCodeOwnedByAuthor(recaps []Recap) map[Author]int {
	linesOfCodeMap := make(map[Author]int)

	// Get list of all authors
	allAuthors := getAuthorsFromRecaps(recaps)

	for _, author := range allAuthors {
		linesOfCodeMap[author] = 0
	}

	for _, recap := range recaps {
		for author, lineCount := range recap.TotalLinesOfCodeInRepoByAuthor {
			linesOfCodeMap[Author(author)] += lineCount
		}
	}

	return linesOfCodeMap
}

func getMergeCommitsByMonth(recaps []Recap) []int {
	commits := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for _, recap := range recaps {
		for idx, commitMonth := range recap.MergeCommitsByMonthCurrYear {
			commits[idx] += commitMonth.Commits
		}
	}

	return commits
}

func getMergeCommitsByWeekDay(recaps []Recap) []int {
	commits := []int{0, 0, 0, 0, 0, 0, 0}

	for _, recap := range recaps {
		for idx, commitMonth := range recap.MergeCommitsByWeekDayCurrYear {
			commits[idx] += commitMonth.Commits
		}
	}

	return commits
}

func getMergeCommitsByHour(recaps []Recap) []int {
	commits := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for _, recap := range recaps {
		for idx, commitMonth := range recap.MergeCommitsByHourCurrYear {
			commits[idx] += commitMonth.Commits
		}
	}

	return commits
}

func getAvgMergesPerDayByRepo(recaps []Recap) map[Repo]float64 {
	mergesMap := make(map[Repo]float64)

	for _, recap := range recaps {
		mergesMap[Repo(recap.Name)] = recap.AvgMergesToMasterPerDayCurrYear
	}

	return mergesMap
}

func getMergesToMasterByRepo(recaps []Recap) map[Repo]YearComparison {
	mergesMap := make(map[Repo]YearComparison)

	for _, recap := range recaps {
		mergesMap[Repo(recap.Name)] = YearComparison{
			Prev: recap.MergesToMasterPrevYear,
			Curr: recap.MergesToMasterCurrYear,
		}
	}

	return mergesMap
}

func getMergesToMasterAllTimeByRepo(recaps []Recap) map[Repo]int {
	mergesMap := make(map[Repo]int)

	for _, recap := range recaps {
		mergesMap[Repo(recap.Name)] = recap.MergesToMasterAllTime
	}

	return mergesMap
}
