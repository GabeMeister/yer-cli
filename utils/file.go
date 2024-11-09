package utils

import (
	"errors"
	"os"
	"path"
	"strings"
)

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	return !info.IsDir()

}

func GetFileExtension(filePath string) string {
	ext := path.Ext(filePath)

	if strings.HasPrefix(ext, ".") {
		return ext[1:]
	} else {
		return ext
	}
}
