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
	ExcludedFiles          []string
	ExcludedAuthors        []string
	AllAuthors             []string
	DuplicateAuthors       []DuplicateAuthorGroup
	IncludeFileBlames      bool
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

func AddRepoConfig(config ConfigFile) ConfigFile {
	// Find the largest id
	largestId := 1
	for _, repoConfig := range config.Repos {
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
		IncludeFileBlames:     false,
	}
	config.Repos = append(config.Repos, newRepoConfig)

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

func SaveConfig(config ConfigFile) {
	SaveDataToFile(config, utils.DEFAULT_CONFIG_FILE)
}

func DoesConfigExist(path string) bool {
	data, err := os.ReadFile(path)

	return err == nil && len(data) > 0
}
