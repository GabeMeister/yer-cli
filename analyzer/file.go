package analyzer

import (
	"encoding/json"
	"os"
)

type FileHandler struct{}

func (f *FileHandler) SaveDataToFile(data any, path string) {
	rawData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fileErr := os.WriteFile(path, rawData, 0644)
	if fileErr != nil {
		panic(fileErr)
	}
}
