package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"GabeMeister/yer-cli/presentation/helpers"
	"GabeMeister/yer-cli/presentation/views/components"
	"GabeMeister/yer-cli/presentation/views/components/ConfigSetupPage"
	"GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

var InitialAuthors = []string{"Kenny", "Kenny1", "Kenny2", "Isaac Neace", "Gabe Jensen", "ktrotter", "Kaleb Trotter", "Stephen Bremer", "Kenny Kline", "Ezra Youngren", "Isaac", "Steve Bremer"}

func addAnalyzerRoutes(e *echo.Echo) {

	e.GET("/create-recap", func(c echo.Context) error {
		if !analyzer.DoesConfigExist(analyzer.DEFAULT_CONFIG_FILE) {
			analyzer.InitConfig(analyzer.ConfigFileOptions{
				RepoDir:                "",
				MasterBranchName:       "",
				IncludedFileExtensions: []string{},
				ExcludedDirs:           []string{},
				ExcludedFiles:          []string{},
				ExcludedAuthors:        []string{},
				AllAuthors:             []string{},
				DuplicateAuthors:       []analyzer.DuplicateAuthorGroup{},
				IncludeFileBlames:      true,
			})

			content := t.Render(t.RenderParams{
				C:         c,
				Component: pages.CreateRecap(),
			})

			return c.HTML(http.StatusOK, content)
		} else {
			config := analyzer.MustGetConfig(analyzer.DEFAULT_CONFIG_FILE)

			url := fmt.Sprintf("/add-repo?id=%d", config.Repos[0].Id)
			c.Redirect(301, url)

			return nil
		}
	})

	e.POST("/create-recap", func(c echo.Context) error {
		configFile := analyzer.MustGetConfig(analyzer.DEFAULT_CONFIG_FILE)
		recapName := helpers.MustGetFormValue(c, "recap-name")

		configFile.Name = recapName
		configFile.Save()

		url := fmt.Sprintf("/add-repo?id=%d", configFile.Repos[0].Id)
		c.Response().Header().Set("HX-Redirect", url)

		return c.HTML(http.StatusOK, "")
	})

	e.GET("/add-repo", func(c echo.Context) error {
		if !analyzer.DoesConfigExist(analyzer.DEFAULT_CONFIG_FILE) {
			// If there's no config, redirect
			c.Redirect(301, "/create-recap")
			return nil
		}

		config := analyzer.MustGetConfig(analyzer.DEFAULT_CONFIG_FILE)

		if len(config.Repos) == 0 {
			// If there's no repos, then the easiest thing to do is just to restart
			// the whole process
			c.Redirect(301, "/create-recap")
			return nil
		}

		newParam := c.QueryParam("new")

		if newParam == "true" {
			newRepo := config.AddNewRepoConfig()
			config.Save()

			c.Redirect(301, fmt.Sprintf("/add-repo?id=%d", newRepo.Id))
			return nil
		}

		id, err := helpers.GetIntQueryParam(c, "id")
		if err != nil {
			return RenderErrorMessage(c, err)
		}

		repoIdx := config.GetRepoIndex(id)
		if repoIdx == -1 {
			return RenderMessage(c, "Couldn't find correct repo to edit. Please restart setup process by visiting `/create-recap")
		}

		repo := config.Repos[repoIdx]

		var ungroupedAuthors []string
		if repo.Path != "" && repo.MasterBranchName != "" {
			duplicateAuthors := analyzer.GetDuplicateAuthorList(repo)
			ungroupedAuthors = analyzer.GetAuthorsFromRepo(
				repo.Path,
				repo.MasterBranchName,
				duplicateAuthors,
			)
		}

		year := time.Now().Year()

		content := t.Render(t.RenderParams{
			C: c,
			Component: pages.ConfigSetup(pages.ConfigSetupProps{
				Id:                    repo.Id,
				RecapName:             config.Name,
				RepoPath:              repo.Path,
				Year:                  year,
				MasterBranch:          repo.MasterBranchName,
				UngroupedAuthors:      ungroupedAuthors,
				DuplicateAuthorGroups: repo.DuplicateAuthors,
				IncludeFileExtensions: strings.Join(repo.IncludeFileExtensions, ","),
				ExcludeDirs:           strings.Join(repo.ExcludeDirectories, ","),
				ExcludeFiles:          strings.Join(repo.ExcludeFiles, ","),
				ExcludeAuthors:        strings.Join(repo.ExcludeAuthors, ","),
				RepoConfigList:        config.Repos,
			}),
		})

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/config-file", func(c echo.Context) error {
		repoIdParam := c.FormValue("id")
		repoId, err := strconv.Atoi(repoIdParam)
		if err != nil {
			panic(fmt.Sprintf("Repo ID param is not a number: %s", repoIdParam))
		}
		repoPath := c.FormValue("repo-path")
		masterBranchName := c.FormValue("master-branch-name")
		includeFileExtensions := c.FormValue("include-file-extensions")
		excludeDirs := c.FormValue("exclude-dirs")
		excludeFiles := c.FormValue("exclude-files")
		excludeAuthors := c.FormValue("exclude-authors")
		formParams, _ := c.FormParams()
		marshaledDupGroups := formParams["dup-group"]
		dupGroups := []analyzer.DuplicateAuthorGroup{}
		for _, g := range marshaledDupGroups {
			dupGroups = append(dupGroups, helpers.UnmarshalDuplicateGroup(g))
		}
		ungroupedAuthors := formParams["ungrouped-author"]

		config := analyzer.MustGetConfig(analyzer.DEFAULT_CONFIG_FILE)
		repo := analyzer.MustGetRepoConfig(config, repoId)
		repoIdx := config.GetRepoIndex(repo.Id)

		repo.Path = repoPath
		repo.MasterBranchName = masterBranchName
		repo.IncludeFileExtensions = strings.Split(includeFileExtensions, ",")
		repo.ExcludeDirectories = strings.Split(excludeDirs, ",")
		repo.ExcludeFiles = strings.Split(excludeFiles, ",")
		repo.ExcludeAuthors = strings.Split(excludeAuthors, ",")
		repo.DuplicateAuthors = dupGroups

		config.Repos[repoIdx] = repo

		config.Save()

		year := time.Now().Year()

		component := pages.ConfigSetup(pages.ConfigSetupProps{
			Id:                    repo.Id,
			RepoConfigList:        config.Repos,
			RecapName:             config.Name,
			RepoPath:              repo.Path,
			Toast:                 "Saved!",
			Year:                  year,
			MasterBranch:          masterBranchName,
			IncludeFileExtensions: includeFileExtensions,
			ExcludeDirs:           excludeDirs,
			ExcludeFiles:          excludeFiles,
			ExcludeAuthors:        excludeAuthors,
			UngroupedAuthors:      ungroupedAuthors,
			DuplicateAuthorGroups: dupGroups,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/repo-config/delete", func(c echo.Context) error {
		repoIdParam := c.FormValue("id")
		repoId, err := strconv.Atoi(repoIdParam)
		if err != nil {
			panic(fmt.Sprintf("Repo ID param is not a number: %s", repoIdParam))
		}

		config := analyzer.MustGetConfig(analyzer.DEFAULT_CONFIG_FILE)
		config = analyzer.RemoveRepoFromConfig(config, repoId)
		config.Save()

		redirectUrl := fmt.Sprintf("/add-repo?id=%d", config.Repos[0].Id)

		c.Response().Header().Set("HX-Redirect", redirectUrl)

		return c.NoContent(http.StatusOK)
	})

	e.GET("/dir-list-modal", func(c echo.Context) error {
		baseDir := c.FormValue("base-dir")
		if baseDir == "" {
			homeDir, _ := os.UserHomeDir()
			baseDir = homeDir
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
		baseDir := c.FormValue("base-dir")
		isGitRepo := analyzer.IsValidGitRepo(baseDir)

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
				RepoPath:  baseDir,
				OutOfBand: true,
			})

			// Updates the repo path form input in the original form
			content += t.Render(t.RenderParams{
				C:         c,
				Component: component,
			})

			masterBranchName := analyzer.GetMasterBranchName(baseDir)
			masterBranchInput := ConfigSetupPage.MasterBranchInput(ConfigSetupPage.MasterBranchInputProps{
				Name:      masterBranchName,
				OutOfBand: true,
			})

			// Updates the master branch input
			content += t.Render(t.RenderParams{
				C:         c,
				Component: masterBranchInput,
			})

			fileExtensions := analyzer.GetFileExtensionsInRepo(baseDir)
			fileExtInput := ConfigSetupPage.IncludeFileExtensions(ConfigSetupPage.IncludeFileExtensionsProps{
				IncludeFileExtensions: strings.Join(fileExtensions, ","),
				OutOfBand:             true,
			})

			// Updates their file extensions with what is in the repo
			content += t.Render(t.RenderParams{
				C:         c,
				Component: fileExtInput,
			})

			authors := analyzer.GetAuthorsFromRepo(baseDir, masterBranchName, []string{})
			dupGroupsBtn := ConfigSetupPage.DuplicateGroupBtn(ConfigSetupPage.DuplicateGroupBtnProps{
				UngroupedAuthors: authors,
				OutOfBand:        true,
			})

			// Updates their file extensions with what is in the repo
			content += t.Render(t.RenderParams{
				C:         c,
				Component: dupGroupsBtn,
			})

		} else {
			// If the repo isn't valid, display the Directory List form with an error
			searchTerm := c.FormValue("search-term")
			filteredDirs := utils.GetFilteredDirs(baseDir, searchTerm)
			component := ConfigSetupPage.DirectoryListForm(ConfigSetupPage.DirectoryListFormProps{
				Dirs:       filteredDirs,
				BaseDir:    baseDir,
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
		baseDir := c.FormValue("base-dir")
		if baseDir == "" {
			homeDir, _ := os.UserHomeDir()
			baseDir = homeDir
		}
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

	e.POST("/duplicate-authors-modal", func(c echo.Context) error {
		formValues, _ := c.FormParams()
		ungroupedAuthors := formValues["ungrouped-author"]
		existingDupGroupsRaw := formValues["dup-group"]
		existingDupGroups := []analyzer.DuplicateAuthorGroup{}
		for _, dupGroupRaw := range existingDupGroupsRaw {
			existingDupGroups = append(existingDupGroups, helpers.UnmarshalDuplicateGroup(dupGroupRaw))
		}

		component := ConfigSetupPage.DuplicateAuthorModal(ConfigSetupPage.DuplicateAuthorModalProps{
			UngroupedAuthors:  ungroupedAuthors,
			ExistingDupGroups: existingDupGroups,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/duplicate-author-grouping", func(c echo.Context) error {
		formValues, _ := c.FormParams()
		authorsMarkedAsDuplicate := formValues["author-marked-as-duplicate"]
		existingDupGroupsRaw := formValues["existing-dup-group"]
		existingDupGroups := []analyzer.DuplicateAuthorGroup{}
		for _, dupGroupRaw := range existingDupGroupsRaw {
			existingDupGroups = append(existingDupGroups, helpers.UnmarshalDuplicateGroup(dupGroupRaw))
		}
		realName := c.FormValue("real-name")
		ungroupedAuthors := formValues["ungrouped-author"]

		if realName == "" {
			component := ConfigSetupPage.DuplicateAuthorForm(ConfigSetupPage.DuplicateAuthorFormProps{
				UngroupedAuthors:  ungroupedAuthors,
				ExistingDupGroups: existingDupGroups,
				SelectedAuthors:   authorsMarkedAsDuplicate,
				Errors: map[string]string{
					"real-name": "Please enter the real name to use",
				},
			})
			content := t.Render(t.RenderParams{
				C:         c,
				Component: component,
			})

			return c.HTML(http.StatusOK, content)
		}

		if len(authorsMarkedAsDuplicate) <= 1 {
			component := ConfigSetupPage.DuplicateAuthorModal(ConfigSetupPage.DuplicateAuthorModalProps{
				UngroupedAuthors:  ungroupedAuthors,
				ExistingDupGroups: existingDupGroups,
				Errors: map[string]string{
					"author-marked-as-duplicate": "Please select at least 2 authors to group together",
				},
			})
			content := t.Render(t.RenderParams{
				C:         c,
				Component: component,
			})

			return c.HTML(http.StatusOK, content)
		}

		component := components.EmptyDiv()
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		// Update the duplicate authors input
		filteredUngroupedAuthors := []string{}
		for _, ungroupedAuthor := range ungroupedAuthors {
			if !slices.Contains(authorsMarkedAsDuplicate, ungroupedAuthor) {
				filteredUngroupedAuthors = append(filteredUngroupedAuthors, ungroupedAuthor)
			}
		}

		dupGroup := analyzer.DuplicateAuthorGroup{
			Real:       realName,
			Duplicates: authorsMarkedAsDuplicate,
		}
		existingDupGroups = append(existingDupGroups, dupGroup)
		component = ConfigSetupPage.DuplicateGroupBtn(ConfigSetupPage.DuplicateGroupBtnProps{
			OutOfBand:        true,
			UngroupedAuthors: filteredUngroupedAuthors,
			DuplicateAuthors: existingDupGroups,
		})
		content += t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		content += "<div id='modal-root' hx-swap-oob='true'></div>"

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/duplicate-author-grouping", func(c echo.Context) error {
		dupGroupToDelete := c.FormValue("dup-group-real-name-to-delete")
		formParams, _ := c.FormParams()
		marshalledDupGroups := formParams["dup-group"]
		ungroupedAuthors := formParams["ungrouped-author"]

		dupGroups := []analyzer.DuplicateAuthorGroup{}
		for _, marshaledDupGroup := range marshalledDupGroups {
			dupGroup := helpers.UnmarshalDuplicateGroup(marshaledDupGroup)
			if dupGroup.Real != dupGroupToDelete {
				dupGroups = append(dupGroups, dupGroup)
			} else {
				// The authors of the deleted duplicate group now become ungrouped
				ungroupedAuthors = append(ungroupedAuthors, dupGroup.Duplicates...)
			}
		}

		component := ConfigSetupPage.DuplicateGroupBtn(ConfigSetupPage.DuplicateGroupBtnProps{
			UngroupedAuthors: ungroupedAuthors,
			DuplicateAuthors: dupGroups,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.GET("/recap-name-textbox", func(c echo.Context) error {
		recapName := c.FormValue("recap-name")
		component := ConfigSetupPage.RecapNameTextbox(ConfigSetupPage.RecapNameTextboxProps{
			RecapName: recapName,
		})
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/recap-name", func(c echo.Context) error {
		recapName := c.FormValue("recap-name")

		config := analyzer.MustGetConfig(analyzer.DEFAULT_CONFIG_FILE)
		config.Name = recapName
		config.Save()

		component := ConfigSetupPage.RecapNameDisplay(ConfigSetupPage.RecapNameTextboxProps{
			RecapName: recapName,
		})
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
