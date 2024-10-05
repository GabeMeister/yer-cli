package presentation

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"os"
)

func getRecap() (analyzer.Recap, error) {
	if !utils.HasRepoBeenAnalyzed() {
		return analyzer.Recap{}, os.ErrNotExist
	}

	data, err := os.ReadFile(utils.RECAP_FILE)
	if err != nil {
		panic(err)
	}
	var repoRecap analyzer.Recap

	json.Unmarshal(data, &repoRecap)

	return repoRecap, nil
}
