package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"GabeMeister/yer-cli/presentation/helpers"
	"GabeMeister/yer-cli/presentation/views/components"
	"GabeMeister/yer-cli/presentation/views/components/ConfigSetupPage"
	"GabeMeister/yer-cli/presentation/views/components/errormessages"
	"GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
			c.Redirect(307, url)

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

		id, err := helpers.GetIntQueryParam(c, "id")
		if err != nil {
			return RenderErrorMessage(c, err)
		}

		repoIdx := config.GetRepoIndex(id)
		if repoIdx == -1 {
			component := errormessages.RepoNotFoundError(errormessages.RepoNotFoundErrorProps{
				RepoId:   id,
				RepoList: config.Repos,
			})
			content := t.Render(t.RenderParams{
				C:         c,
				Component: component,
			})

			return c.HTML(http.StatusOK, content)
		}

		// Clear any "empty" repos (only if the user isn't actively looking at it)
		config.RemoveEmptyReposAndSave(id)

		repo := config.Repos[repoIdx]

		var ungroupedAuthors []string
		if repo.Path != "" && repo.MasterBranchName != "" {
			duplicateAuthors := repo.GetDuplicateAuthorList()
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
				AnalyzeFileBlames:     repo.AnalyzeFileBlames,
			}),
		})

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/config-setup", func(c echo.Context) error {
		repoIdParam := c.FormValue("id")
		repoId, err := strconv.Atoi(repoIdParam)
		if err != nil {
			panic(fmt.Sprintf("Repo ID param is not a number: %s", repoIdParam))
		}
		repoPath := c.FormValue("repo-path")
		masterBranchName := c.FormValue("master-branch-name")
		includeFileExtensions := helpers.UnmarshalStrSlice(c.FormValue("include-file-extensions"))
		excludeDirs := helpers.UnmarshalStrSlice(c.FormValue("exclude-dirs"))
		excludeFiles := helpers.UnmarshalStrSlice(c.FormValue("exclude-files"))
		excludeAuthors := helpers.UnmarshalStrSlice(c.FormValue("exclude-authors"))
		analyzeFileBlames := true
		formParams, _ := c.FormParams()
		marshaledDupGroups := formParams["dup-group"]
		dupGroups := []analyzer.DuplicateAuthorGroup{}
		for _, g := range marshaledDupGroups {
			dupGroups = append(dupGroups, helpers.UnmarshalDuplicateGroup(g))
		}
		// ungroupedAuthors := formParams["ungrouped-author"]

		config := analyzer.MustGetConfig(analyzer.DEFAULT_CONFIG_FILE)
		repo := analyzer.MustGetRepoConfig(config, repoId)
		repoIdx := config.GetRepoIndex(repo.Id)

		repo.Path = repoPath
		repo.MasterBranchName = masterBranchName
		repo.IncludeFileExtensions = includeFileExtensions
		repo.ExcludeDirectories = excludeDirs
		repo.ExcludeFiles = excludeFiles
		repo.ExcludeAuthors = excludeAuthors
		repo.DuplicateAuthors = dupGroups
		repo.AnalyzeFileBlames = analyzeFileBlames

		config.Repos[repoIdx] = repo
		config.Save()

		// Check for existing "empty" repo, redirect to that one. Otherwise, we just
		// redirect to create a new repo
		emptyRepoIdx := utils.FindIndex(config.Repos, func(repo analyzer.RepoConfig) bool {
			return repo.Path == ""
		})

		emptyRepoId := -1
		if emptyRepoIdx == -1 {
			newRepo := config.AddNewRepoConfig()
			config.Save()
			emptyRepoId = newRepo.Id
		} else {
			emptyRepoId = config.Repos[emptyRepoIdx].Id
		}

		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/add-repo?id=%d", emptyRepoId))

		return c.NoContent(http.StatusOK)
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
			// The situation where they were searching and just hit enter to select
			// the first directory in the list
			baseDir = c.FormValue("first_dir")
		}

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

		if len(authorsMarkedAsDuplicate) == 0 {
			component := ConfigSetupPage.DuplicateAuthorModal(ConfigSetupPage.DuplicateAuthorModalProps{
				UngroupedAuthors:  ungroupedAuthors,
				ExistingDupGroups: existingDupGroups,
				Errors: map[string]string{
					"author-marked-as-duplicate": "Please select at least 1 author",
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

	e.GET("/finish-setup", func(c echo.Context) error {
		config := analyzer.MustGetConfig(analyzer.DEFAULT_CONFIG_FILE)
		config.RemoveEmptyReposAndSave(-1)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10)
			defer cancel()

			yellow := "\033[1;33m"
			reset := "\033[0m"

			fmt.Printf("%s┌──────────────────────────────────────┐%s\n", yellow, reset)
			fmt.Printf("%s│ Great! Now run the following command │%s\n", yellow, reset)
			fmt.Printf("%s│ to analyze your stats:               │%s\n", yellow, reset)
			fmt.Printf("%s│                                      │%s\n", yellow, reset)
			fmt.Printf("%s│ yer -a                               │%s\n", yellow, reset)
			fmt.Printf("%s└──────────────────────────────────────┘%s\n", yellow, reset)
			fmt.Println()

			time.Sleep(100 * time.Millisecond)
			e.Shutdown(ctx)
		}()

		content := t.Render(t.RenderParams{
			C:         c,
			Component: pages.FinishSetup(),
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/close-window", func(c echo.Context) error {
		os.Exit(0)

		return c.HTML(http.StatusOK, "")
	})

	e.GET("/commonly-changed-files", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Hey</h1>")
	})

}
