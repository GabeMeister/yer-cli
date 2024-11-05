package analyzer

import (
	"encoding/json"
	"os"
	"sort"
)

func SaveDataToFile(data any, path string) {
	rawData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fileErr := os.WriteFile(path, rawData, 0644)
	if fileErr != nil {
		panic(fileErr)
	}
}

func GetFileChangeRatio(insertionsByEngineer map[string]int, deletionsByEngineer map[string]int) map[string]float64 {
	ratios := make(map[string]float64)

	// Assume all engineers have done insertions and deletions
	for engineer, insertionCount := range insertionsByEngineer {
		insertions := float64(insertionCount)
		deletions := float64(deletionsByEngineer[engineer])

		ratios[engineer] = insertions / deletions
	}

	return ratios
}

func GetCommonlyChangedFiles() []FileChangeCount {
	commits := getCurrYearGitCommits()
	fileChangeTracker := make(map[string]int)

	for _, commit := range commits {
		for _, changes := range commit.FileChanges {
			fileChangeTracker[changes.FilePath] += 1
		}
	}

	fileChangeSlice := []FileChangeCount{}
	for file, count := range fileChangeTracker {
		fileChangeSlice = append(fileChangeSlice, FileChangeCount{
			File:  file,
			Count: count,
		})
	}

	sort.Slice(fileChangeSlice, func(i int, j int) bool {
		return fileChangeSlice[i].Count > fileChangeSlice[j].Count
	})

	return fileChangeSlice[0:5]
}
