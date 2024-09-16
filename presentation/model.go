package presentation

import (
	"GabeMeister/yer-cli/analyzer"
	"encoding/json"
	"os"
)

func getRecap() analyzer.Recap {
	data, err := os.ReadFile("./tmp/recap.json")
	if err != nil {
		panic(err)
	}
	var repoRecap analyzer.Recap

	json.Unmarshal(data, &repoRecap)

	return repoRecap
}
