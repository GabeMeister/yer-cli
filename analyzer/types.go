package analyzer

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

type CommitList []string

type AuthorCommitList map[string]CommitList

type DayCommitListByAuthor map[string]AuthorCommitList
