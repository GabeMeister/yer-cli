package analyzer

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func HasRecapBeenRan() bool {
	files, err := os.ReadDir("tmp")
	if err != nil {
		return false
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "_recap.json") {
			return true
		}
	}

	return false
}

func saveDataToFile(data any, path string) {
	rawData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fileErr := os.WriteFile(path, rawData, 0644)
	if fileErr != nil {
		panic(fileErr)
	}
}

func isFileReadable(path string) bool {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	file.Close()

	return true
}

func getSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[14], 100*time.Millisecond)
}

func getNumWorkDaysInCurrYear() int {
	total := 0
	firstDayOfYear := time.Date(CURR_YEAR, time.January, 1, 0, 0, 0, 0, time.UTC)
	currDay := time.Now()

	tmpDay := firstDayOfYear
	for tmpDay.Before(currDay) {
		weekDay := tmpDay.Weekday()

		if weekDay != time.Saturday && weekDay != time.Sunday {
			total += 1
		}
		tmpDay = tmpDay.AddDate(0, 0, 1)
	}

	return total
}
