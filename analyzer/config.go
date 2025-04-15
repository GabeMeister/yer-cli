package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"os"
)

type ConfigFileOptions struct {
	RepoDir                string
	MasterBranchName       string
	IncludedFileExtensions []string
	ExcludedDirs           []string
	AllAuthors             []string
	DuplicateAuthors     []DuplicateAuthorGroup
	IncludeFileBlames      bool
}

func InitConfig(options ConfigFileOptions) ConfigFile {
	config := ConfigFile{
		Repos: []RepoConfig{
			{
				Version:               "0.0.1",
				Path:                  options.RepoDir,
				Name:                  "",
				MasterBranchName:      options.MasterBranchName,
				IncludeFileExtensions: options.IncludedFileExtensions,
				ExcludeDirectories:    options.ExcludedDirs,
				ExcludeFiles:          []string{},
				ExcludeAuthors:      []string{},
				DuplicateAuthors:    options.DuplicateAuthors,
				IncludeFileBlames:     options.IncludeFileBlames,
			},
		},
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile(utils.DEFAULT_CONFIG_FILE, data, 0644)

	return config
}

func updateDuplicateAuthors(path string, duplicateAuthors []DuplicateAuthorGroup) error {
	// Update config, cause we wanna remember this for later
	config := GetConfig(path)
	config.Repos[0].DuplicateAuthors = duplicateAuthors
	SaveDataToFile(config, path)

	// Also want to update the commits.json file, replacing the duplicate git
	// usernames with the real ones
	commits := getGitCommits()

	for i := range commits {
		realUsername := GetRealAuthorName(config.Repos[0], commits[i].Author)
		commits[i].Author = realUsername
	}
	SaveDataToFile(commits, utils.COMMITS_FILE)

	return nil
}

func GetConfig(path string) ConfigFile {
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

func UpdateConfig(config ConfigFile) {
	SaveDataToFile(config, utils.DEFAULT_CONFIG_FILE)
}

func DoesConfigExist(path string) bool {
	_, err := os.ReadFile(path)

	return err == nil
}
