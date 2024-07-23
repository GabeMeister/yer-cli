package presentation

import (
	"GabeMeister/yer-cli/analyzer"
	"encoding/json"
	"os"
)

func getSummary() analyzer.Recap {
	data, err := os.ReadFile("./tmp/summary.json")
	if err != nil {
		panic(err)
	}
	var repoSummary analyzer.Recap

	json.Unmarshal(data, &repoSummary)

	return repoSummary
}
