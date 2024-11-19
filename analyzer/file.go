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

		if deletions == 0 {
			ratios[engineer] = 1
		} else {
			ratios[engineer] = insertions / deletions
		}
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

func GetLargestFilesCurrYear() []FileSize {
	fileBlames := GetCurrYearFileBlames()

	sort.Slice(fileBlames, func(i int, j int) bool {
		return fileBlames[i].LineCount > fileBlames[j].LineCount
	})

	fileSizes := []FileSize{}
	for _, fileBlame := range fileBlames[0:10] {
		fileSizes = append(fileSizes, FileSize{
			File:      fileBlame.File,
			LineCount: fileBlame.LineCount,
		})
	}

	return fileSizes
}

func GetSmallestFilesCurrYear() []FileSize {
	fileBlames := GetCurrYearFileBlames()

	sort.Slice(fileBlames, func(i int, j int) bool {
		return fileBlames[i].LineCount < fileBlames[j].LineCount
	})

	fileSizes := []FileSize{}
	for _, fileBlame := range fileBlames[0:10] {
		fileSizes = append(fileSizes, FileSize{
			File:      fileBlame.File,
			LineCount: fileBlame.LineCount,
		})
	}

	return fileSizes
}

func GetTotalLinesOfCodePrevYear() int {
	fileBlames := GetPrevYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func GetTotalLinesOfCodeCurrYear() int {
	fileBlames := GetCurrYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func GetSizeOfRepoByWeekCurrYear() []RepoSizeTimeStamp {
	commits := getCurrYearGitCommits()
	weeks := utils.GetWeeksOfYear()

	weekCommitsMap := make(map[int][]GitCommit)
	for _, week := range weeks {
		weekCommitsMap[week] = nil
	}

	for _, commit := range commits {
		date := utils.GetDateFromISOString(commit.Date)
		_, week := date.ISOWeek()

		weekCommitsMap[week] = append(weekCommitsMap[week], commit)
	}

	totalLinesOfCodePrevYear := GetTotalLinesOfCodePrevYear()
	runningTotal := totalLinesOfCodePrevYear
	final := []RepoSizeTimeStamp{}

	for week := 1; week <= 52; week++ {
		weekCommits := weekCommitsMap[week]

		for _, commit := range weekCommits {
			for _, fileChanges := range commit.FileChanges {
				runningTotal += fileChanges.Insertions
				runningTotal -= fileChanges.Deletions
			}

		}

		final = append(final, RepoSizeTimeStamp{
			WeekNumber: week,
			LineCount:  runningTotal,
		})
	}

	return final
}
