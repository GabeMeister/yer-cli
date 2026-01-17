package analyzer

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func HasRecapBeenRan() bool {
	_, err := os.ReadFile("tmp/multi_repo_recap.json")

	return err == nil
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

func OpenBrowser(url string) error {
	if !strings.HasPrefix(url, "http") {
		return errors.New("must have https at the start of the url")
	}

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux, freebsd, openbsd, netbsd
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

func getNumChangesInGitCommit(commit GitCommit) int {
	numChanges := 0
	for _, change := range commit.FileChanges {
		numChanges += change.Insertions
		numChanges += change.Deletions
	}

	return numChanges
}
