package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"encoding/json"
	"os"
)

func GetRecap() (Recap, error) {
	if !utils.HasRepoBeenAnalyzed() {
		return Recap{}, os.ErrNotExist
	}

	data, err := os.ReadFile(utils.RECAP_FILE)
	if err != nil {
		panic(err)
	}
	var repoRecap Recap

	json.Unmarshal(data, &repoRecap)

	return repoRecap, nil
}
