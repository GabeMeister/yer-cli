package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
)

func GetCommitsFile(config RepoConfig) string {
	return fmt.Sprintf(utils.COMMITS_FILE_TEMPLATE, filepath.Base(config.Path))
}

func GetMergeCommitsFile(config RepoConfig) string {
	return fmt.Sprintf(utils.MERGE_COMMITS_FILE_TEMPLATE, filepath.Base(config.Path))
}

func GetDirectPushesFile(config RepoConfig) string {
	return fmt.Sprintf(utils.DIRECT_PUSH_ON_MASTER_COMMITS_FILE_TEMPLATE, filepath.Base(config.Path))
}

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

func GetPrevYearFileList() []string {
	bytes, err := os.ReadFile(utils.PREV_YEAR_FILE_LIST_FILE)
	if err != nil {
		panic(err)
	}

	var files []string
	jsonErr := json.Unmarshal(bytes, &files)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return files

}

func GetCurrYearFileList() []string {
	bytes, err := os.ReadFile(utils.CURR_YEAR_FILE_LIST_FILE)
	if err != nil {
		panic(err)
	}

	var files []string
	jsonErr := json.Unmarshal(bytes, &files)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return files

}

func HasPrevYearFileBlames() bool {
	_, err := os.ReadFile(utils.PREV_YEAR_FILE_BLAMES_FILE)

	return err == nil
}

func HasCurrYearFileBlames() bool {
	_, err := os.ReadFile(utils.CURR_YEAR_FILE_BLAMES_FILE)

	return err == nil
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

func GetFileChangeRatio(insertionsByAuthor map[string]int, deletionsByAuthor map[string]int) map[string]float64 {
	ratios := make(map[string]float64)

	for author, insertionCount := range insertionsByAuthor {
		insertions := float64(insertionCount)
		deletions := float64(deletionsByAuthor[author])

		if deletions == 0 {
			ratios[author] = 1
		} else {
			ratios[author] = insertions / deletions
		}
	}

	for _, ratio := range ratios {
		if math.IsNaN(ratio) {
			fmt.Print("\n\n", "*** ERROR in GetFileChangeRatio(). Ratio is NaN: ***", "\n", ratios, "\n\n\n")
			panic("Ratio is NaN in GetFileChangeRatio()!")
		}
	}

	return ratios
}

func GetCommonlyChangedFiles(config RepoConfig) []FileChangeCount {
	commits := getCurrYearGitCommits(config)
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
	files := GetPrevYearFileList()

	return len(files)
}

func GetFileCountCurrYear() int {
	files := GetCurrYearFileList()

	return len(files)
}

func GetLargestFilesCurrYear() []FileSize {
	if !HasCurrYearFileBlames() {
		return []FileSize{}
	}

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
	if !HasCurrYearFileBlames() {
		return []FileSize{}
	}

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
	if !HasPrevYearFileBlames() {
		return 0
	}

	fileBlames := GetPrevYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func GetTotalLinesOfCodeCurrYear() int {
	if !HasCurrYearFileBlames() {
		return 0
	}

	fileBlames := GetCurrYearFileBlames()
	total := 0
	for _, fileBlame := range fileBlames {
		total += fileBlame.LineCount
	}

	return total
}

func GetSizeOfRepoByWeekCurrYear(config RepoConfig) []RepoSizeTimeStamp {
	commits := getCurrYearGitCommits(config)
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
