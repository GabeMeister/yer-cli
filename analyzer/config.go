package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func initConfig(repoDir string, includedFileExtensions []string, excludedDirs []string, duplicateEngineers map[string]string) Config {
	config := Config{
		Path:                  repoDir,
		Name:                  filepath.Base(repoDir),
		IncludeFileExtensions: includedFileExtensions,
		ExcludeDirectories:    excludedDirs,
		ExcludeFiles:          []string{},
		ExcludeEngineers:      []string{},
		DuplicateEngineers:    duplicateEngineers,
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("./config.json", data, 0644)

	return config
}

func getConfig(path string) Config {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var data Config
	jsonErr := json.Unmarshal(bytes, &data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return data
}

// A lot of times in repos somehow the same user has two different git usernames
// (for example, Gabe Jensen and GabeJensen). Could be because they changed
// laptops, decided to change their user name, etc. To make the stats more
// accurate, we should "bucket" the duplicate users into one, and hence this
// helper function.
func getRealUsername(userName string, config Config) string {
	if config.DuplicateEngineers[userName] != "" {
		return config.DuplicateEngineers[userName]
	} else {
		return userName
	}
}
