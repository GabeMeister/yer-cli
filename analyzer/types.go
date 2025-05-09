package analyzer

type ConfigFile struct {
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Repos   []RepoConfig `json:"repos"`
}

type DuplicateAuthorGroup struct {
	Real       string   `json:"real"`
	Duplicates []string `json:"duplicates"`
}

type RepoConfig struct {
	Id                    int                    `json:"id"`
	Path                  string                 `json:"path"`
	MasterBranchName      string                 `json:"master_branch_name"`
	IncludeFileExtensions []string               `json:"include_file_extensions"`
	ExcludeDirectories    []string               `json:"exclude_directories"`
	ExcludeFiles          []string               `json:"exclude_files"`
	ExcludeAuthors        []string               `json:"exclude_authors"`
	DuplicateAuthors      []DuplicateAuthorGroup `json:"duplicate_authors"`
	AllAuthors            []string               `json:"all_authors"`
	IncludeFileBlames     bool                   `json:"include_file_blames"`
}

type FileChange struct {
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
	FilePath   string `json:"file_path"`
}

type FileChangeCount struct {
	File  string `json:"file"`
	Count int    `json:"count"`
}

type GitCommit struct {
	Commit      string       `json:"commit"`
	Author      string       `json:"author"`
	Email       string       `json:"email"`
	Message     string       `json:"message"`
	Date        string       `json:"date"`
	FileChanges []FileChange `json:"file_changes"`
}

type CommitMonth struct {
	Month   string `json:"month"`
	Commits int    `json:"commits"`
}

type CommitWeekDay struct {
	Day     string `json:"day"`
	Commits int    `json:"commits"`
}

type CommitHour struct {
	Hour    string `json:"hour"`
	Commits int    `json:"commits"`
}

type MostSingleDayCommitsByAuthor struct {
	Username string   `json:"username"`
	Date     string   `json:"date"`
	Count    int      `json:"count"`
	Commits  []string `json:"commits"`
}

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
	NewAuthorCommitsCurrYear             []GitCommit                  `json:"new_author_commits_curr_year"`
	NewAuthorCountCurrYear               int                          `json:"new_author_count_curr_year"`
	NewAuthorListCurrYear                []string                     `json:"new_author_list_curr_year"`
	AuthorCommitCountsCurrYear           map[string]int               `json:"author_commit_counts_curr_year"`
	AuthorCommitCountsAllTime            map[string]int               `json:"author_commit_counts_all_time"`
	AuthorCountCurrYear                  int                          `json:"author_count_curr_year"`
	AuthorCountAllTime                   int                          `json:"author_count_all_time"`
	AuthorCommitsOverTimeCurrYear        []TotalCommitCount           `json:"author_commits_over_time_curr_year"`
	AuthorFileChangesOverTimeCurrYear    []TotalFileChangeCount       `json:"author_file_changes_over_time_curr_year"`
	MostSingleDayCommitsByAuthorCurrYear MostSingleDayCommitsByAuthor `json:"most_single_day_commits_by_author_curr_year"`
	DirectPushesOnMasterByAuthorCurrYear map[string]int               `json:"direct_pushes_on_master_by_author_curr_year"`
	MergesToMasterByAuthorCurrYear       map[string]int               `json:"merges_to_master_by_author_curr_year"`
	FileChangesByAuthorCurrYear          map[string]int               `json:"file_changes_by_author_curr_year"`
	FileChangeRatioByAuthorCurrYear      map[string]float64           `json:"file_change_ratio_by_author_curr_year"`
	TotalLinesOfCodeInRepoByAuthor       map[string]int               `json:"total_lines_of_code_in_repo_by_author"`
}

// Example: { date: '2023-01-03T08:00:00.000Z', name: 'Steve Bremer', value: 24 },
// Used for Author Commits Over Time racing bar chart
type TotalCommitCount struct {
	// ISO Date string
	Date  string `json:"date"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Example: { date: '2023-01-03T08:00:00.000Z', name: 'Steve Bremer', value: 2400 },
// Used for Author Commits Over Time racing bar chart
type TotalFileChangeCount struct {
	// ISO Date string
	Date  string `json:"date"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type CommitMessageLengthFrequency struct {
	Length    int `json:"length"`
	Frequency int `json:"frequency"`
}

type MostMergesInOneDay struct {
	Count   int         `json:"count"`
	Date    string      `json:"date"`
	Commits []GitCommit `json:"commits"`
}

type FileBlame struct {
	File      string         `json:"file"`
	LineCount int            `json:"line_count"`
	GitBlame  map[string]int `json:"git_blame"`
}

type FileSize struct {
	File      string `json:"file"`
	LineCount int    `json:"line_count"`
}

type RepoSizeTimeStamp struct {
	WeekNumber int `json:"week_number"`
	LineCount  int `json:"line_count"`
}
