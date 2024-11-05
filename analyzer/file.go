package analyzer

import (
	"GabeMeister/yer-cli/utils"
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

func GetPrevYearFileBlames() []FileBlame {
	bytes, err := os.ReadFile(utils.PREV_YEAR_FILE_BLAMES_FILE)
	if err != nil {
		panic(err)
	}

	var fileBlames []FileBlame
	jsonErr := json.Unmarshal(bytes, &fileBlames)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return fileBlames
}

func GetCurrYearFileBlames() []FileBlame {
	bytes, err := os.ReadFile(utils.CURR_YEAR_FILE_BLAMES_FILE)
	if err != nil {
		panic(err)
	}

	var fileBlames []FileBlame
	jsonErr := json.Unmarshal(bytes, &fileBlames)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return fileBlames
}

// func GetFilesCurrYear() []FileBlame {
// }

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

func GetFileCountPrevYear() int {
	fileBlames := GetPrevYearFileBlames()

	return len(fileBlames)
}

func GetFileCountCurrYear() int {
	fileBlames := GetCurrYearFileBlames()

	return len(fileBlames)
}
