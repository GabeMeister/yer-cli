package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

func initConfig(repoDir string, includedFileExtensions []string, excludedDirs []string) Config {
	now := time.Now()
	isoDateString := now.Format(time.RFC3339)

	config := Config{
		Path:                  repoDir,
		Name:                  filepath.Base(repoDir),
		DateAnalyzed:          isoDateString,
		IncludeFileExtensions: includedFileExtensions,
		ExcludeDirectories:    excludedDirs,
		ExcludeFiles:          []string{},
		ExcludeEngineers:      []string{},
		DuplicateEngineers:    make(map[string]string),
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
