package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"os"
	"path/filepath"
)

type ConfigFileOptions struct {
	RepoDir                string
	MasterBranchName       string
	IncludedFileExtensions []string
	ExcludedDirs           []string
	DuplicateEngineers     map[string]string
	IncludeFileBlames      bool
}

func InitConfig(options ConfigFileOptions) ConfigFile {
	config := ConfigFile{
		Repos: []RepoConfig{
			{
				Version:               "0.0.1",
				Path:                  options.RepoDir,
				Name:                  filepath.Base(options.RepoDir),
				MasterBranchName:      options.MasterBranchName,
				IncludeFileExtensions: options.IncludedFileExtensions,
				ExcludeDirectories:    options.ExcludedDirs,
				ExcludeFiles:          []string{},
				ExcludeEngineers:      []string{},
				DuplicateEngineers:    options.DuplicateEngineers,
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

func updateDuplicateEngineers(path string, duplicateEngineers map[string]string) error {
	// Update config, cause we wanna remember this for later
	config := GetConfig(path)
	config.Repos[0].DuplicateEngineers = duplicateEngineers
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
