package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func getCommitsFromGitLogs(config RepoConfig, mergeCommits bool) []GitCommit {
	s := GetSpinner()

	fmt.Println()
	s.Suffix = " Retrieving git logs..."
	s.Start()

	path := config.Path
	args := []string{
		"git",
		"log",
	}

	if mergeCommits {
		args = append(args, "--merges")
		args = append(args, "--first-parent")
		args = append(args, config.MasterBranchName)
	} else {
		args = append(args, "--no-merges")
	}

	args = append(args, "--reverse")
	args = append(args, "--format=-- Begin --%n-- Commit --%n%H%n-- Author --%n%aN%n-- Email --%n%aE%n-- Date --%n%ad%n-- Message --%n%B%n-- End --")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = path

	rawOutput, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	output := string(rawOutput)

	s.Suffix = " Analyzing git logs..."

	lines := strings.Split(output, "\n")

	var currentState = "BEGIN"
	var currentCommit GitCommit = GitCommit{
		Commit:  "",
		Author:  "",
		Email:   "",
		Message: "",
		Date:    "",
	}
	var commits []GitCommit

	for _, line := range lines {
		if line == "" {
			continue
		}

		switch line {
		case "-- Begin --":
			currentCommit = GitCommit{
				Commit:  "",
				Author:  "",
				Email:   "",
				Message: "",
				Date:    "",
			}
		case "-- Commit --":
			currentState = "COMMIT"
		case "-- Author --":
			currentState = "AUTHOR"
		case "-- Email --":
			currentState = "EMAIL"
		case "-- Date --":
			currentState = "DATE"
		case "-- Message --":
			currentState = "MESSAGE"
		case "-- End --":
			commits = append(commits, currentCommit)
		default:
			switch currentState {
			case "COMMIT":
				currentCommit.Commit = line
			case "AUTHOR":
				currentCommit.Author = line
			case "DATE":
				currentCommit.Date = line
			case "MESSAGE":
				if currentCommit.Message == "" {
					currentCommit.Message += line
				} else {
					currentCommit.Message += "|||" + line
				}
			case "EMAIL":
				currentCommit.Email = line
			default:
				panic(fmt.Sprintf("Unrecognized state: %s", currentState))
			}
		}
	}

	s.Stop()

	if !mergeCommits {
		fileChangeSummary := getFileChangeSummary(config)

		for i := range commits {
			commits[i].FileChanges = fileChangeSummary[commits[i].Commit]
		}
	}

	return commits
}

func getDirectPushToMasterCommitsCurrYear(config RepoConfig) []GitCommit {
	path := config.Path
	commits := getCurrYearGitCommits()

	args := []string{
		"git",
		"log",
		"--no-merges",
		"--reverse",
		"--first-parent",
		config.MasterBranchName,
		"--format=%H",
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = path

	rawOutput, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	output := string(rawOutput)

	commitHashes := strings.Split(output, "\n")
	commitHashMap := make(map[string]bool)

	for _, hash := range commitHashes {
		commitHashMap[hash] = true
	}

	directPushToMasterCommits := []GitCommit{}
	for _, commit := range commits {
		if commitHashMap[commit.Commit] {
			directPushToMasterCommits = append(directPushToMasterCommits, commit)
		}
	}

	sort.Slice(directPushToMasterCommits, func(i int, j int) bool {
		date1 := utils.GetDateFromISOString(directPushToMasterCommits[i].Date)
		date2 := utils.GetDateFromISOString(directPushToMasterCommits[j].Date)

		return date1.UnixMicro() < date2.UnixMicro()
	})

	return directPushToMasterCommits
}

func isFileChangeLine(line string) bool {
	// Regex for matching email addresses
	var emailRegex = regexp.MustCompile(`^\d+\s+\d+\s+.+$`)
	return emailRegex.MatchString(line)
}

func getFileChangeSummary(config RepoConfig) map[string][]FileChange {
	s := GetSpinner()
	s.Suffix = " Retrieving file changes..."
	s.Start()

	path := config.Path
	cmd := exec.Command(
		"git",
		"log",
		"--no-merges",
		"--reverse",
		"--after",
		fmt.Sprintf("%d-01-01", CURR_YEAR),
		"--numstat")
	cmd.Dir = path

	rawOutput, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	output := string(rawOutput)

	s.Suffix = " Analyzing file changes..."

	lines := strings.Split(output, "\n")
	fileChangeMap := make(map[string][]FileChange)

	currHash := ""
	currFileChanges := []FileChange{}

	for _, line := range lines {
		if strings.HasPrefix(line, "commit") {
			// We found a new commit, so we need to add the previous commit in and
			// reset the temp variables
			if currHash != "" {
				currFileChanges = filterToOnlyIncludedFiles(config, currFileChanges)

				fileChangeMap[currHash] = currFileChanges
				currHash = ""
				currFileChanges = []FileChange{}
			}

			// Initialize a "new" commit
			tokens := strings.Split(line, " ")
			currHash = tokens[1]
		} else if isFileChangeLine(line) {
			// Regex to match any whitespace
			whitespace := regexp.MustCompile(`\s+`)

			// Split the string by any whitespace
			parts := whitespace.Split(line, -1)
			insertions, _ := strconv.Atoi(parts[0])
			deletions, _ := strconv.Atoi(parts[1])
			filePath := parts[2]

			currFileChanges = append(currFileChanges, FileChange{
				Insertions: insertions,
				Deletions:  deletions,
				FilePath:   filePath,
			})
		}
	}

	// Add in the final commit
	fileChangeMap[currHash] = currFileChanges

	s.Stop()

	return fileChangeMap

}

func filterToOnlyIncludedFiles(config RepoConfig, fileChanges []FileChange) []FileChange {
	filteredFileChanges := utils.Filter(fileChanges, func(c FileChange) bool {
		fileExt := utils.GetFileExtension(c.FilePath)

		return utils.Includes(config.IncludeFileExtensions, func(ext string) bool {
			return fileExt == ext
		})
	})

	return filteredFileChanges
}

func getRepoFiles(config RepoConfig, commitOrBranchName string) []string {
	args := []string{
		"git",
		"ls-tree",
		"-r",
		commitOrBranchName,
		"--name-only",
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = config.Path

	rawOutput, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	output := string(rawOutput)
	files := strings.Split(output, "\n")
	includedFiles := []string{}

	// Filter to just files that we want to include and filter out files we want
	// to exclude
	for _, file := range files {
		fileExtension := utils.GetFileExtension(file)
		validFileExtension := utils.Includes(config.IncludeFileExtensions, func(ext string) bool {
			return fileExtension == ext
		})

		if !validFileExtension {
			continue
		}

		isExcludedFile := utils.Includes(config.ExcludeDirectories, func(dir string) bool {
			return strings.HasPrefix(file, dir)
		})

		if isExcludedFile {
			continue
		}

		includedFiles = append(includedFiles, file)
	}

	return includedFiles
}

func GetFileBlameSummary(config RepoConfig, files []string) []FileBlame {
	s := GetSpinner()
	fmt.Println()
	s.Suffix = " Analyzing Git blames..."
	s.Start()

	fileBlames := []FileBlame{}
	totalFiles := len(files)

	for idx, file := range files {
		s.Suffix = fmt.Sprintf(" Processed %d/%d files. (currently on %s)...", idx, totalFiles, file)
		args := []string{
			"git",
			"blame",
			file,
			"--line-porcelain",
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = config.Path

		rawOutput, err := cmd.Output()
		if err != nil {
			panic(err)
		}

		output := string(rawOutput)
		lines := strings.Split(output, "\n")
		lines = utils.Filter(lines, func(line string) bool {
			return strings.HasPrefix(line, "committer ")
		})
		authors := utils.Map(lines, func(line string) string {
			authorName := strings.ReplaceAll(line, "committer ", "")

			return GetRealAuthorName(config, authorName)
		})

		authorLineCountMap := make(map[string]int)
		for _, author := range authors {
			authorLineCountMap[author] += 1
		}

		fileBlames = append(fileBlames, FileBlame{
			File:      file,
			LineCount: len(authors),
			GitBlame:  authorLineCountMap,
		})
	}

	s.Stop()

	return fileBlames
}

func getLastCommitPrevYear() GitCommit {
	commits := getPrevYearGitCommits()
	lastIdx := len(commits) - 1

	return commits[lastIdx]
}

func checkoutRepoToCommitOrBranchName(config RepoConfig, commitOrBranchName string) error {
	args := []string{
		"git",
		"checkout",
		commitOrBranchName,
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = config.Path

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func GetRealAuthorName(config RepoConfig, authorName string) string {
	name := authorName

	for {
		realAuthorName, ok := config.DuplicateEngineers[name]
		if ok {
			name = realAuthorName

			// Loop again in case there's a "chain" of duplicate user names. For
			// example one could have, "ktrotter" -> "kaleb.trotter",
			// but then also have "kaleb.trotter" -> "Kaleb Trotter"
			continue
		}

		break
	}

	return name
}

func stashRepo(dir string) {
	stashCmd := exec.Command("git", "stash")
	stashCmd.Dir = dir
	stashCmd.Output()
}

func pullRepo(dir string) {
	pullcmd := exec.Command("git", "pull")
	pullcmd.Dir = dir
	pullcmd.Output()
}

func hasPrevYearCommits() bool {
	commits := getPrevYearGitCommits()

	return len(commits) > 0
}
