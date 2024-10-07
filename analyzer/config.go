package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func initConfig(repoDir string, includedFileExtensions []string, excludedDirs []string) Config {
	config := Config{
		Path:                  repoDir,
		Name:                  filepath.Base(repoDir),
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
