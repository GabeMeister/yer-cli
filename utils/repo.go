package utils

import (
	"errors"
	"os"
)

func HasRepoBeenAnalyzed() bool {
	_, fileErr := os.Stat(RECAP_FILE)

	return !errors.Is(fileErr, os.ErrNotExist)
}
