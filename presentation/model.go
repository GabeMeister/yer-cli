package presentation

import (
	"GabeMeister/yer-cli/analyzer"
	"encoding/json"
	"os"
)

func getSummary() analyzer.RepoSummary {
	data, err := os.ReadFile("./tmp/summary.json")
	if err != nil {
		panic(err)
	}
	var repoSummary analyzer.RepoSummary

	json.Unmarshal(data, &repoSummary)

	return repoSummary
}
