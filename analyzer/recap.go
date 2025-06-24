package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Recap struct {
	Version            string `json:"version"`
	Name               string `json:"name"`
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

func GetRepoRecap() (Recap, error) {
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

	data, err := os.ReadFile(fmt.Sprintf("./tmp/%s", recapFile))
	if err != nil {
		panic(err)
	}
	var repoRecap Recap

	json.Unmarshal(data, &repoRecap)

	return repoRecap, nil
}
