package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"GabeMeister/yer-cli/presentation/views/components/ConfigSetupPage"
	"GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

var InitialAuthors = []string{"Kenny", "Kenny1", "Kenny2", "Isaac Neace", "Gabe Jensen", "ktrotter", "Kaleb Trotter", "Stephen Bremer", "Kenny Kline", "Ezra Youngren", "Isaac", "Steve Bremer"}

func addAnalyzerRoutes(e *echo.Echo) {

	e.GET("/create-recap", func(c echo.Context) error {
		if !analyzer.DoesConfigExist(utils.DEFAULT_CONFIG_FILE) {
			analyzer.InitConfig(analyzer.ConfigFileOptions{
				RepoDir:                "",
				MasterBranchName:       "",
				IncludedFileExtensions: []string{},
				ExcludedDirs:           []string{},
				AllAuthors:             []string{},
				DuplicateAuthors:       []analyzer.DuplicateAuthorGroup{},
				IncludeFileBlames:      true,
			})
		}

		config := analyzer.GetConfig(utils.DEFAULT_CONFIG_FILE)

		year := time.Now().Year()

		content := t.Render(t.RenderParams{
			C: c,
			Component: pages.ConfigSetup(pages.ConfigSetupProps{
				RecapName:    config.Repos[0].Name,
				RepoPath:     config.Repos[0].Path,
				Year:         year,
				MasterBranch: config.Repos[0].MasterBranchName,
				AllAuthors:   []string{"Ezra", "Isaac", "Steve"},
				// AllAuthors:            config.Repos[0].AllAuthors,
				IncludeFileExtensions: strings.Join(config.Repos[0].IncludeFileExtensions, ","),
				ExcludeDirs:           strings.Join(config.Repos[0].ExcludeDirectories, ","),
				ExcludeFiles:          strings.Join(config.Repos[0].ExcludeFiles, ","),
			}),
		})

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/config-file", func(c echo.Context) error {
		recapName := c.FormValue("recap-name")
		repoPath := c.FormValue("repo-path")
		masterBranchName := c.FormValue("master-branch-name")
		includeFileExtensions := c.FormValue("include-file-extensions")
		excludeDirs := c.FormValue("exclude-dirs")
		excludeFiles := c.FormValue("exclude-files")

		config := analyzer.GetConfig(utils.DEFAULT_CONFIG_FILE)
		config.Repos[0].Name = recapName
		config.Repos[0].Path = repoPath
		config.Repos[0].MasterBranchName = masterBranchName
		config.Repos[0].IncludeFileExtensions = strings.Split(includeFileExtensions, ",")
		config.Repos[0].ExcludeDirectories = strings.Split(excludeDirs, ",")
		config.Repos[0].ExcludeFiles = strings.Split(excludeFiles, ",")

		analyzer.UpdateConfig(config)

		year := time.Now().Year()

		component := pages.ConfigSetup(pages.ConfigSetupProps{
			RecapName:             config.Repos[0].Name,
			RepoPath:              config.Repos[0].Path,
			Toast:                 "Saved!",
			Year:                  year,
			MasterBranch:          masterBranchName,
			IncludeFileExtensions: includeFileExtensions,
			ExcludeDirs:           excludeDirs,
			ExcludeFiles:          excludeFiles,
			AllAuthors:            []string{"David"},
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/dir-list-modal", func(c echo.Context) error {
		config := analyzer.GetConfig(utils.DEFAULT_CONFIG_FILE)
		baseDir, _ := os.UserHomeDir()
		if len(config.Repos) > 0 && config.Repos[0].Path != "" {
			baseDir = config.Repos[0].Path
		}
		dirs := utils.GetDirs(baseDir)

		component := ConfigSetupPage.DirectoryListModal(ConfigSetupPage.DirectoryListModalProps{
			BaseDir: baseDir,
			Dirs:    dirs,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/dir-list", func(c echo.Context) error {
		baseDir := c.FormValue("dir")
		if baseDir == "" {
			panic("Using PATCH /dir-list wrong: need to include base dir")
		}

		dirs := utils.GetDirs(baseDir)

		component := ConfigSetupPage.DirectoryListForm(ConfigSetupPage.DirectoryListFormProps{
			Dirs:    dirs,
			BaseDir: baseDir,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/repo-path", func(c echo.Context) error {
		repoPath := c.FormValue("repo-path")
		isGitRepo := analyzer.IsValidGitRepo(repoPath)

		content := ""

		if isGitRepo {
			// Since we're updating part of the modal with this endpoint, but we
			// actually want to clear the modal, we just return bogus html for now,
			// and then update the modal and the repo path input out of band
			content += `<div></div>`

			// Clears out the modal
			content += `
				<div id='modal-root' hx-swap-oob="true"></div>
			`

			component := ConfigSetupPage.RepoPath(ConfigSetupPage.RepoPathProps{
				RepoPath:  repoPath,
				OutOfBand: true,
			})

			// Updates the repo path form input in the original form
			content += t.Render(t.RenderParams{
				C:         c,
				Component: component,
			})

			masterBranchName := analyzer.GetMasterBranchName(repoPath)
			masterBranchInput := ConfigSetupPage.MasterBranchInput(ConfigSetupPage.MasterBranchInputProps{
				Name:      masterBranchName,
				OutOfBand: true,
			})

			// Updates the master branch input
			content += t.Render(t.RenderParams{
				C:         c,
				Component: masterBranchInput,
			})

			fileExtensions := analyzer.GetFileExtensionsInRepo(repoPath)
			fileExtInput := ConfigSetupPage.IncludeFileExtensions(ConfigSetupPage.IncludeFileExtensionsProps{
				IncludeFileExtensions: strings.Join(fileExtensions, ","),
				OutOfBand:             true,
			})

			// Updates their file extensions with what is in the repo
			content += t.Render(t.RenderParams{
				C:         c,
				Component: fileExtInput,
			})

			authors := analyzer.GetAuthorsFromRepo(repoPath, masterBranchName)
			fmt.Print("\n\n", "*** authors ***", "\n", authors, "\n\n\n")
			// allAuthorsComponent := ConfigSetupPage.AllAuthorsList(ConfigSetupPage.AllAuthorsListProps{
			// 	AllAuthors: authors,
			// 	OutOfBand:  true,
			// })

			// // Updates their file extensions with what is in the repo
			// content += t.Render(t.RenderParams{
			// 	C:         c,
			// 	Component: allAuthorsComponent,
			// })

		} else {
			// If the repo isn't valid, display the Directory List form with an error
			searchTerm := c.FormValue("search-term")
			filteredDirs := utils.GetFilteredDirs(repoPath, searchTerm)
			component := ConfigSetupPage.DirectoryListForm(ConfigSetupPage.DirectoryListFormProps{
				Dirs:       filteredDirs,
				BaseDir:    repoPath,
				Error:      "This is not a Git repo!",
				SearchTerm: searchTerm,
			})
			content += t.Render(t.RenderParams{
				C:         c,
				Component: component,
			})
		}

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/filtered-dir-contents", func(c echo.Context) error {
		searchTerm := c.FormValue("search-term")
		baseDir := c.FormValue("repo-path")
		filteredDirs := utils.GetFilteredDirs(baseDir, searchTerm)

		component := ConfigSetupPage.DirectoryList(ConfigSetupPage.DirectoryListProps{
			BaseDir: baseDir,
			Dirs:    filteredDirs,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/duplicate-authors-modal", func(c echo.Context) error {
		ungroupedAuthors := c.FormValue("ungrouped-authors")
		fmt.Print("\n\n", "*** ungroupedAuthors ***", "\n", ungroupedAuthors, "\n\n\n")
		duplicateAuthors := c.FormValue("duplicate-authors")
		fmt.Print("\n\n", "*** duplicateAuthors ***", "\n", duplicateAuthors, "\n\n\n")

		component := ConfigSetupPage.DuplicateAuthorModal(ConfigSetupPage.DuplicateAuthorModalProps{})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/clear-modal", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<div id='modal-root'></div>")
	})

	// e.POST("/search-authors", func(c echo.Context) error {
	// 	text := c.FormValue("filter-text")
	// 	text = strings.ToLower(text)

	// 	dupAuthorsRaw := c.FormValue("duplicate-authors")
	// 	tempDupAuthors := strings.Split(dupAuthorsRaw, ",")

	// 	matches := []string{}
	// 	for _, author := range InitialAuthors {
	// 		lowerCaseAuthor := strings.ToLower(author)

	// 		if strings.Contains(lowerCaseAuthor, text) && !slices.Contains(tempDupAuthors, author) {
	// 			matches = append(matches, author)
	// 		}
	// 	}
	// 	component := ConfigSetupPage.AllAuthorsList(matches)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)
	// })

	// e.PATCH("/temp-duplicate-group", func(c echo.Context) error {

	// 	leftItemsStr := c.FormValue("all-authors")
	// 	leftItems := strings.Split(leftItemsStr, ",")

	// 	rightItemsStr := c.FormValue("duplicate-authors")
	// 	rightItems := strings.Split(rightItemsStr, ",")

	// 	filterText := c.FormValue("filter-text")

	// 	allAuthors := lo.Filter(leftItems, func(author string, _ int) bool {
	// 		return author != ""
	// 	})

	// 	selectedAuthors := lo.Filter(rightItems, func(author string, _ int) bool {
	// 		return author != ""
	// 	})

	// 	config := analyzer.GetConfig("./config.json")

	// 	component := pages.(allAuthors, selectedAuthors, config.Repos[0].DuplicateAuthors, filterText)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)
	// })

	// 	e.POST("/duplicate-group", func(c echo.Context) error {
	// 		duplicatesListRaw := c.FormValue("duplicate-authors")

	// 		duplicateAuthors := strings.Split(duplicatesListRaw, ",")

	// 		config := analyzer.GetConfig("./config.json")
	// 		config.Repos[0].DuplicateAuthors = append(config.Repos[0].DuplicateAuthors, analyzer.DuplicateAuthorGroup{
	// 			Real:       duplicateAuthors[0],
	// 			Duplicates: duplicateAuthors[1:],
	// 		})

	// 		allDups := make(map[string]bool)
	// 		for _, dupGroup := range config.Repos[0].DuplicateAuthors {
	// 			allDups[dupGroup.Real] = true

	// 			for _, dup := range dupGroup.Duplicates {
	// 				allDups[dup] = true
	// 			}
	// 		}

	// 		remainingAuthors := []string{}
	// 		for _, author := range InitialAuthors {
	// 			if _, found := allDups[author]; !found {
	// 				remainingAuthors = append(remainingAuthors, author)
	// 			}
	// 		}

	// 		analyzer.SaveDataToFile(config, "./config.json")

	// 		component := pages.ConfigSetup(pages.ConfigSetupProps{
	// 			RecapName: config.Repos[0].Name,
	// 		})
	// 		content := t.Render(t.RenderParams{
	// 			C:         c,
	// 			Component: component,
	// 		})

	//		return c.HTML(http.StatusOK, content)
	//	})
}
