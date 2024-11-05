package analyzer

import (
	"encoding/json"
	"os"
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

func GetFileChangeRatio(insertionsByEngineer map[string]int, deletionsByEngineer map[string]int) map[string]float64 {
	ratios := make(map[string]float64)

	// Assume all engineers have done insertions and deletions
	for engineer, insertionCount := range insertionsByEngineer {
		insertions := float64(insertionCount)
		deletions := float64(deletionsByEngineer[engineer])

		ratios[engineer] = insertions / deletions
	}

	return ratios
}
