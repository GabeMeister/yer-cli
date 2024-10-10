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

type GitCommit struct {
	Commit  string `json:"commit"`
	Author  string `json:"author"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Date    string `json:"date"`
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
	NumCommitsInPast   int `json:"num_commits_in_past"`

	// Team
	EngineerCommitCountsCurrYear map[string]int `json:"engineer_commit_counts_curr_year"`
	EngineerCommitCountsAllTime  map[string]int `json:"engineer_commit_counts_all_time"`
}
