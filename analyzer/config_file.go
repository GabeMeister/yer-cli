package analyzer

import (
	"encoding/json"
	"os"
)

type ConfigFileOptions struct {
	RepoDir                string
	MasterBranchName       string
	IncludedFileExtensions []string
	ExcludedDirs           []string
	ExcludedFiles          []string
	ExcludedAuthors        []string
	AllAuthors             []string
	DuplicateAuthors       []DuplicateAuthorGroup
	IncludeFileBlames      bool
}

type ConfigFile struct {
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Repos   []RepoConfig `json:"repos"`
}

func InitConfig(options ConfigFileOptions) ConfigFile {
	config := ConfigFile{
		Name:    "",
		Version: "0.0.1",
		Repos: []RepoConfig{
			{
				Id:                    1,
				Path:                  options.RepoDir,
				MasterBranchName:      options.MasterBranchName,
				IncludeFileExtensions: options.IncludedFileExtensions,
				ExcludeDirectories:    options.ExcludedDirs,
				ExcludeFiles:          []string{},
				ExcludeAuthors:        []string{},
				DuplicateAuthors:      options.DuplicateAuthors,
				AnalyzeFileBlames:     options.IncludeFileBlames,
			},
		},
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile(DEFAULT_CONFIG_FILE, data, 0644)

	return config
}

func DoesConfigExist(path string) bool {
	data, err := os.ReadFile(path)

	if err == nil && len(data) > 0 {
		config := MustGetConfig(path)
		if config.Name != "" {
			return true
		}
	}

	return false
}

func MustGetConfig(path string) ConfigFile {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var data ConfigFile
	jsonErr := json.Unmarshal(bytes, &data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return data
}

func RemoveRepoFromConfig(c ConfigFile, repoId int) ConfigFile {
	repoIdx := c.GetRepoIndex(repoId)
	c.Repos = append(c.Repos[:repoIdx], c.Repos[repoIdx+1:]...)

	return c
}

func MustGetRepoConfig(config ConfigFile, repoId int) RepoConfig {
	var repo RepoConfig
	for _, r := range config.Repos {

		if r.Id == repoId {
			repo = r
			break
		}
	}

	if repo.Id == 0 {
		panic("Could not find correct repo to patch in config file")
	}

	return repo
}

func (c *ConfigFile) AddNewRepoConfig() *RepoConfig {
	// Find the largest id
	largestId := 1
	for _, repoConfig := range c.Repos {
		if repoConfig.Id > largestId {
			largestId = repoConfig.Id
		}
	}

	newRepoConfig := RepoConfig{
		Id:                    largestId + 1,
		Path:                  "",
		MasterBranchName:      "",
		IncludeFileExtensions: []string{},
		ExcludeDirectories:    []string{},
		ExcludeFiles:          []string{},
		ExcludeAuthors:        []string{},
		DuplicateAuthors:      []DuplicateAuthorGroup{},
		AnalyzeFileBlames:     false,
	}
	c.Repos = append(c.Repos, newRepoConfig)

	return &newRepoConfig
}

func (c *ConfigFile) GetRepoIndex(repoId int) int {
	index := -1

	for idx, r := range c.Repos {

		if r.Id == repoId {
			index = idx
			break
		}
	}

	return index
}

func (c *ConfigFile) Save() {
	saveDataToFile(c, DEFAULT_CONFIG_FILE)
}

func (c *ConfigFile) updateDuplicateAuthors(r *RepoConfig) error {
	repoIdx := c.GetRepoIndex(r.Id)

	// Update config, cause we wanna remember this for later
	c.Repos[repoIdx].DuplicateAuthors = r.DuplicateAuthors
	saveDataToFile(c, DEFAULT_CONFIG_FILE)

	// Also want to update the commits.json file, replacing the duplicate git
	// usernames with the real ones
	commits := r.getGitCommits()

	for i := range commits {
		realUsername := c.Repos[repoIdx].getRealAuthorName(commits[i].Author)
		commits[i].Author = realUsername
	}
	saveDataToFile(commits, r.getCommitsFile())

	return nil
}
