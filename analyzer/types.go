package analyzer

type ConfigFile struct {
	Repos []RepoConfig `json:"repos"`
}

type RepoConfig struct {
	Name                  string            `json:"name"`
	Path                  string            `json:"path"`
	IncludeFileExtensions []string          `json:"include_file_extensions"`
	ExcludeDirectories    []string          `json:"exclude_directories"`
	ExcludeFiles          []string          `json:"exclude_files"`
	ExcludeEngineers      []string          `json:"exclude_engineers"`
	DuplicateEngineers    map[string]string `json:"duplicate_engineers"`
}

type FileChange struct {
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
	FilePath   string `json:"file_path"`
}

type GitCommit struct {
	Commit      string       `json:"commit"`
	Author      string       `json:"author"`
	Email       string       `json:"email"`
	Message     string       `json:"message"`
	Date        string       `json:"date"`
	FileChanges []FileChange `json:"file_changes"`
}

type GitMergeCommit struct {
	Commit             string
	FirstParentCommit  string // Typically master when you branched off
	SecondParentCommit string // Typically the final commit of the MR
	Message            string
	Author             string
	Email              string
	Date               string
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

type MostSingleDayCommitsByEngineer struct {
	Username string   `json:"username"`
	Date     string   `json:"date"`
	Count    int      `json:"count"`
	Commits  []string `json:"commits"`
}

type Recap struct {
	Name         string `json:"name"`
	DateAnalyzed string `json:"date_analyzed"`

	// Commits
	NumCommitsAllTime              int                            `json:"num_commits_all_time"`
	NumCommitsPrevYear             int                            `json:"num_commits_prev_year"`
	NumCommitsCurrYear             int                            `json:"num_commits_curr_year"`
	CommitsByMonthCurrYear         []CommitMonth                  `json:"commits_by_month_curr_year"`
	CommitsByWeekDayCurrYear       []CommitWeekDay                `json:"commits_by_week_day_curr_year"`
	CommitsByHourCurrYear          []CommitHour                   `json:"commits_by_hour_curr_year"`
	MostInsertionsInCommitCurrYear GitCommit                      `json:"most_insertions_in_commit_curr_year"`
	MostDeletionsInCommitCurrYear  GitCommit                      `json:"most_deletions_in_commit__curr_year"`
	LargestCommitMessageCurrYear   GitCommit                      `json:"largest_commit_message_curr_year"`
	SmallestCommitMessagesCurrYear []GitCommit                    `json:"smallest_commit_messages_curr_year"`
	CommitMessageHistogramCurrYear []CommitMessageLengthFrequency `json:"commit_message_histogram_curr_year"`

	// Team
	NewEngineerCommitsCurrYear             []GitCommit                    `json:"new_engineer_commits_curr_year"`
	NewEngineerCountCurrYear               int                            `json:"new_engineer_count_curr_year"`
	EngineerCommitCountsCurrYear           map[string]int                 `json:"engineer_commit_counts_curr_year"`
	EngineerCommitCountsAllTime            map[string]int                 `json:"engineer_commit_counts_all_time"`
	EngineerCountCurrYear                  int                            `json:"engineer_count_curr_year"`
	EngineerCountAllTime                   int                            `json:"engineer_count_all_time"`
	EngineerCommitsOverTimeCurrYear        []TotalCommitCount             `json:"engineer_commits_over_time_curr_year"`
	MostSingleDayCommitsByEngineerCurrYear MostSingleDayCommitsByEngineer `json:"most_single_day_commits_by_engineer_curr_year"`
}

// Example: { date: '2023-01-03T08:00:00.000Z', name: 'Steve Bremer', value: 24 },
// Used for Engineer Commits Over Time racing bar chart
type TotalCommitCount struct {
	// ISO Date string
	Date  string `json:"date"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type CommitMessageLengthFrequency struct {
	Length    int `json:"length"`
	Frequency int `json:"frequency"`
}
