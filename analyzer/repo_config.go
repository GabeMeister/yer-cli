package analyzer

import (
	"fmt"
	"path/filepath"
)

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
	AnalyzeFileBlames     bool                   `json:"analyze_file_blames"`
}

func (r *RepoConfig) GetCommitsFile() string {
	return fmt.Sprintf(COMMITS_FILE_TEMPLATE, filepath.Base(r.Path))
}

func (r *RepoConfig) GetMergeCommitsFile() string {
	return fmt.Sprintf(MERGE_COMMITS_FILE_TEMPLATE, filepath.Base(r.Path))
}

func (r *RepoConfig) GetDirectPushesFile() string {
	return fmt.Sprintf(DIRECT_PUSH_ON_MASTER_COMMITS_FILE_TEMPLATE, filepath.Base(r.Path))
}

func (r *RepoConfig) GetPrevYearFileListFile() string {
	return fmt.Sprintf(PREV_YEAR_FILE_LIST_FILE, filepath.Base(r.Path))
}

func (r *RepoConfig) GetPrevYearFileBlamesFile() string {
	return fmt.Sprintf(PREV_YEAR_FILE_BLAMES_FILE, filepath.Base(r.Path))
}

func (r *RepoConfig) GetCurrYearFileListFile() string {
	return fmt.Sprintf(CURR_YEAR_FILE_LIST_FILE, filepath.Base(r.Path))
}

func (r *RepoConfig) GetCurrYearFileBlamesFile() string {
	return fmt.Sprintf(CURR_YEAR_FILE_BLAMES_FILE, filepath.Base(r.Path))
}

func (r *RepoConfig) GetRecapFile() string {
	return fmt.Sprintf(RECAP_FILE_TEMPLATE, filepath.Base(r.Path))
}

func (r *RepoConfig) GetName() string {
	return filepath.Base(r.Path)
}
