package analyzer

import (
	"time"

	"github.com/briandowns/spinner"
)

func GetSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[14], 100*time.Millisecond)
}
