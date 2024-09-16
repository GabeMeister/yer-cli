package analyzer

type Config struct {
	Name                  string   `json:"name"`
	Path                  string   `json:"path"`
	IncludeFileExtensions []string `json:"include_file_extensions"`
	ExcludeDirectories    []string `json:"exclude_directories"`
	ExcludeFiles          []string `json:"exclude_files"`
	ExcludeEngineers      []string `json:"exclude_engineers"`
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
	Name               string
	NumCommitsAllTime  int
	NumCommitsPrevYear int
	NumCommitsCurrYear int
	NumCommitsInPast   int
}
