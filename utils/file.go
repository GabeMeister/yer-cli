package utils

import (
	"path"
	"strings"
)

func GetFileExtension(filePath string) string {
	ext := path.Ext(filePath)

	if strings.HasPrefix(ext, ".") {
		return ext[1:]
	} else {
		return ext
	}
}
