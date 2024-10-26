package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"os"
	"path/filepath"
)

func initConfig(repoDir string, includedFileExtensions []string, excludedDirs []string, duplicateEngineers map[string]string) ConfigFile {
	config := ConfigFile{
		Repos: []RepoConfig{
			RepoConfig{
				Path:                  repoDir,
				Name:                  filepath.Base(repoDir),
				IncludeFileExtensions: includedFileExtensions,
				ExcludeDirectories:    excludedDirs,
				ExcludeFiles:          []string{},
				ExcludeEngineers:      []string{},
				DuplicateEngineers:    duplicateEngineers,
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
	config := getConfig(path)
	config.Repos[0].DuplicateEngineers = duplicateEngineers
	SaveDataToFile(config, path)

	// Also want to update the commits.json file, replacing the duplicate git
	// usernames with the real ones
	commits := getGitCommits()

	for i := range commits {
		if realUsername, ok := duplicateEngineers[commits[i].Author]; ok {
			commits[i].Author = realUsername
		}
	}
	SaveDataToFile(commits, utils.COMMITS_FILE)

	return nil
}

func getConfig(path string) ConfigFile {
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
