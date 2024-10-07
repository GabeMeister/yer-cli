package analyzer

func GetEngineerCommitCountAllTime() map[string]int {
	commits := getGitCommits()
	engineers := make(map[string]int)

	for _, commit := range commits {
		// TODO: check duplicate authors
		engineers[commit.Author] += 1
	}

	return engineers
}
