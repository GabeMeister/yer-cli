package analyzer

type Config struct {
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

type Recap struct {
	Name         string `json:"name"`
	DateAnalyzed string `json:"date_analyzed"`

	// Commits
	NumCommitsAllTime  int `json:"num_commits_all_time"`
	NumCommitsPrevYear int `json:"num_commits_prev_year"`
	NumCommitsCurrYear int `json:"num_commits_curr_year"`

	// Team
	NewEngineerCommitsCurrYear      []GitCommit        `json:"new_engineer_commits_curr_year"`
	NewEngineerCountCurrYear        int                `json:"new_engineer_count_curr_year"`
	EngineerCommitCountsCurrYear    map[string]int     `json:"engineer_commit_counts_curr_year"`
	EngineerCommitCountsAllTime     map[string]int     `json:"engineer_commit_counts_all_time"`
	EngineerCountCurrYear           int                `json:"engineer_count_curr_year"`
	EngineerCountAllTime            int                `json:"engineer_count_all_time"`
	EngineerCommitsOverTimeCurrYear []TotalCommitCount `json:"engineer_commits_over_time_curr_year"`
}

// Example: { date: '2023-01-03T08:00:00.000Z', name: 'Steve Bremer', value: 24 },
// Used for Engineer Commits Over Time racing bar chart
type TotalCommitCount struct {
	// ISO Date string
	Date  string `json:"date"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}
