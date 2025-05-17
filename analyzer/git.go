package analyzer

import (
	"GabeMeister/yer-cli/utils"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func getCommitsFromGitLogs(r RepoConfig, mergeCommits bool) []GitCommit {
	s := GetSpinner()

	fmt.Println()
	s.Suffix = " Retrieving git logs..."
	s.Start()

	path := r.Path
	args := []string{
		"git",
		"log",
	}

	if mergeCommits {
		args = append(args, "--merges")
		args = append(args, "--first-parent")
		args = append(args, r.MasterBranchName)
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
		fileChangeSummary := getFileChangeSummary(r)

		for i := range commits {
			commits[i].FileChanges = fileChangeSummary[commits[i].Commit]
		}
	}

	return commits
}

func getDirectPushToMasterCommitsCurrYear(r RepoConfig) []GitCommit {
	path := r.Path
	commits := r.getCurrYearGitCommits()

	args := []string{
		"git",
		"log",
		"--no-merges",
		"--reverse",
		"--first-parent",
		r.MasterBranchName,
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

func getFileChangeSummary(r RepoConfig) map[string][]FileChange {
	s := GetSpinner()
	s.Suffix = " Retrieving line changes..."
	s.Start()

	path := r.Path
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

	s.Suffix = " Analyzing line changes..."

	lines := strings.Split(output, "\n")
	fileChangeMap := make(map[string][]FileChange)

	currHash := ""
	currFileChanges := []FileChange{}

	for _, line := range lines {
		if strings.HasPrefix(line, "commit") {
			// We found a new commit, so we need to add the previous commit in and
			// reset the temp variables
			if currHash != "" {
				currFileChanges = filterToOnlyIncludedFiles(r, currFileChanges)

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

func filterToOnlyIncludedFiles(r RepoConfig, fileChanges []FileChange) []FileChange {
	filteredFileChanges := utils.Filter(fileChanges, func(c FileChange) bool {
		fileExt := utils.GetFileExtension(c.FilePath)

		return utils.Includes(r.IncludeFileExtensions, func(ext string) bool {
			return fileExt == ext
		})
	})

	return filteredFileChanges
}

func (r *RepoConfig) getRepoFiles(commitOrBranchName string) []string {
	args := []string{
		"git",
		"ls-tree",
		"-r",
		commitOrBranchName,
		"--name-only",
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = r.Path

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
		validFileExtension := utils.Includes(r.IncludeFileExtensions, func(ext string) bool {
			return fileExtension == ext
		})

		if !validFileExtension {
			continue
		}

		isExcludedFile := utils.Includes(r.ExcludeDirectories, func(dir string) bool {
			return strings.HasPrefix(file, dir)
		})

		if isExcludedFile {
			continue
		}

		includedFiles = append(includedFiles, file)
	}

	return includedFiles
}

func GetFileBlameSummary(r RepoConfig, files []string) []FileBlame {
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
		cmd.Dir = r.Path

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

			return GetRealAuthorName(r, authorName)
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

func (r *RepoConfig) getLastCommitPrevYear() GitCommit {
	commits := r.getPrevYearGitCommits()
	lastIdx := len(commits) - 1

	return commits[lastIdx]
}

func (r *RepoConfig) checkoutRepoToCommitOrBranchName(commitOrBranch string) error {
	args := []string{
		"git",
		"checkout",
		commitOrBranch,
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = r.Path

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func GetRealAuthorName(r RepoConfig, authorName string) string {
	for _, dupGroup := range r.DuplicateAuthors {
		for _, dup := range dupGroup.Duplicates {
			if authorName == dup {
				return dup
			}
		}
	}

	return authorName
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

func (r *RepoConfig) hasPrevYearCommits() bool {
	commits := r.getPrevYearGitCommits()

	return len(commits) > 0
}

func GetAuthorsFromRepo(dir string, branch string, ignoreAuthors []string) []string {
	gitCmd := exec.Command("git", "shortlog", branch, "-s")
	gitCmd.Dir = dir

	var stderr bytes.Buffer
	gitCmd.Stderr = &stderr

	rawOutput, err := gitCmd.Output()
	if err != nil {
		fmt.Println("Error running git command:", err)
		fmt.Println("Stderr:", stderr.String())
		return nil
	}

	text := string(rawOutput)
	lines := strings.Split(strings.TrimSpace(text), "\n")

	authors := []string{}
	for _, line := range lines {
		tokens := strings.Split(line, "\t")
		author := tokens[1]

		if !slices.Contains(ignoreAuthors, author) {
			authors = append(authors, tokens[1])
		}
	}

	return authors

}

func GetDuplicateAuthorList(repo RepoConfig) []string {
	duplicateAuthors := []string{}

	for _, dupGroup := range repo.DuplicateAuthors {
		duplicateAuthors = append(duplicateAuthors, dupGroup.Duplicates...)
	}

	return duplicateAuthors
}

func IsValidGitRepo(dir string) bool {
	_, fileErr := os.Stat(dir)

	if errors.Is(fileErr, os.ErrNotExist) {
		fmt.Println("Directory not found, please try again.")
		return false
	} else if errors.Is(fileErr, os.ErrPermission) {
		fmt.Println("Unable to access directory, make sure it has proper permissions and try again.")
		return false
	} else {
		gitDirPath := filepath.Join(dir, ".git")
		_, gitDirErr := os.Stat(gitDirPath)
		if errors.Is(gitDirErr, os.ErrNotExist) {
			fmt.Println("No Git repo found in specified directory. Please try again.")
			return false
		}
	}

	return true
}

func HasRecapBeenRan() bool {
	_, fileErr := os.Stat(RECAP_FILE)

	return !errors.Is(fileErr, os.ErrNotExist)
}

func GetMasterBranchName(dir string) string {
	masterBranchCmd := exec.Command("git", "branch", "-r")
	masterBranchCmd.Dir = dir
	rawOutput, _ := masterBranchCmd.Output()
	text := string(rawOutput)
	lines := strings.Split(text, "\n")
	idx := slices.IndexFunc(lines, func(line string) bool {
		return strings.Contains(line, "origin/HEAD")
	})
	masterBranchLine := lines[idx]

	masterBranchName := strings.ReplaceAll(masterBranchLine, "origin/HEAD -> origin/", "")
	masterBranchName = strings.TrimSpace(masterBranchName)

	return masterBranchName
}

func GetFileExtensionsInRepo(dir string) []string {
	lsFilesCmd := exec.Command("git", "ls-files")
	lsFilesCmd.Dir = dir
	rawOutput, _ := lsFilesCmd.Output()
	text := string(rawOutput)
	lines := strings.Split(text, "\n")

	extMap := make(map[string]bool)
	for _, line := range lines {
		fileExt := strings.ReplaceAll(filepath.Ext(line), ".", "")
		if fileExt != "" {
			extMap[fileExt] = true
		}
	}

	fileExtensions := []string{}
	for fileExt := range extMap {
		if slices.Contains(SUPPORTED_FILE_EXTENSIONS, fileExt) {
			fileExtensions = append(fileExtensions, fileExt)
		}
	}

	return fileExtensions
}
