package utils

import (
	"errors"
	"fmt"
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

func GetDirs(baseDir string) []string {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return []string{}
	}

	dirs := []string{".."}
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs
}

func GetFilteredDirs(baseDir string, searchTerm string) []string {
	searchTerm = strings.ToLower(searchTerm)
	fmt.Print("\n\n", "*** searchTerm ***", "\n", searchTerm, "\n\n\n")
	fmt.Print("\n\n", "*** baseDir ***", "\n", baseDir, "\n\n\n")
	dirs := GetDirs(baseDir)
	fmt.Print("\n\n", "*** dirs ***", "\n", dirs, "\n\n\n")
	filteredDirs := []string{}

	for _, dir := range dirs {
		if strings.Contains(strings.ToLower(dir), searchTerm) {
			filteredDirs = append(filteredDirs, dir)
		}
	}

	return filteredDirs
}
