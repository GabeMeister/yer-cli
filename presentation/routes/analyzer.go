package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/utils"
	"os"
	"strings"
	"time"

	"GabeMeister/yer-cli/presentation/views/components/ConfigSetupPage"
	"GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

var InitialEngineers = []string{"Kenny", "Kenny1", "Kenny2", "Isaac Neace", "Gabe Jensen", "ktrotter", "Kaleb Trotter", "Stephen Bremer", "Kenny Kline", "Ezra Youngren", "Isaac", "Steve Bremer"}

func addAnalyzerRoutes(e *echo.Echo) {

	e.GET("/create-recap", func(c echo.Context) error {
		if !analyzer.DoesConfigExist(utils.DEFAULT_CONFIG_FILE) {
			analyzer.InitConfig(analyzer.ConfigFileOptions{
				RepoDir:                "",
				MasterBranchName:       "",
				IncludedFileExtensions: []string{},
				ExcludedDirs:           []string{},
				DuplicateEngineers:     []analyzer.DuplicateEngineerGroup{},
				IncludeFileBlames:      true,
			})
		}

		config := analyzer.GetConfig(utils.DEFAULT_CONFIG_FILE)

		content := t.Render(t.RenderParams{
			C: c,
			Component: pages.ConfigSetup(pages.ConfigSetupProps{
				RecapName: config.Repos[0].Name,
				RepoPath:  config.Repos[0].Path,
			}),
		})

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/config-file", func(c echo.Context) error {
		recapName := c.FormValue("recap-name")
		repoPath := c.FormValue("repo-path")

		config := analyzer.GetConfig(utils.DEFAULT_CONFIG_FILE)
		config.Repos[0].Name = recapName
		config.Repos[0].Path = repoPath

		analyzer.UpdateConfig(config)

		component := pages.ConfigSetup(pages.ConfigSetupProps{
			RecapName: config.Repos[0].Name,
			RepoPath:  config.Repos[0].Path,
			Toast:     time.Now().Format("2006-01-02 15:04:05"),
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
		} else {
			// If the repo isn't valid, display the Directory List form with an error
			dirs := utils.GetDirs(repoPath)
			component := ConfigSetupPage.DirectoryListForm(ConfigSetupPage.DirectoryListFormProps{
				Dirs:    dirs,
				BaseDir: repoPath,
				Error:   "This is not a Git repo!",
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
		searchTerm = strings.ToLower(searchTerm)
		baseDir := c.FormValue("repo-path")

		dirs := utils.GetDirs(baseDir)
		filteredDirs := []string{}
		for _, dir := range dirs {
			if strings.Contains(strings.ToLower(dir), searchTerm) {
				filteredDirs = append(filteredDirs, dir)
			}
		}

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

	e.GET("/clear-modal", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<div id='modal-root'></div>")
	})

	// e.POST("/search-engineers", func(c echo.Context) error {
	// 	text := c.FormValue("filter-text")
	// 	text = strings.ToLower(text)

	// 	dupEngineersRaw := c.FormValue("duplicate-engineers")
	// 	tempDupEngineers := strings.Split(dupEngineersRaw, ",")

	// 	matches := []string{}
	// 	for _, engineer := range InitialEngineers {
	// 		lowerCaseEngineer := strings.ToLower(engineer)

	// 		if strings.Contains(lowerCaseEngineer, text) && !slices.Contains(tempDupEngineers, engineer) {
	// 			matches = append(matches, engineer)
	// 		}
	// 	}
	// 	component := ConfigSetupPage.AllEngineersList(matches)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)
	// })

	// e.PATCH("/temp-duplicate-group", func(c echo.Context) error {

	// 	leftItemsStr := c.FormValue("all-engineers")
	// 	leftItems := strings.Split(leftItemsStr, ",")

	// 	rightItemsStr := c.FormValue("duplicate-engineers")
	// 	rightItems := strings.Split(rightItemsStr, ",")

	// 	filterText := c.FormValue("filter-text")

	// 	allEngineers := lo.Filter(leftItems, func(engineer string, _ int) bool {
	// 		return engineer != ""
	// 	})

	// 	selectedEngineers := lo.Filter(rightItems, func(engineer string, _ int) bool {
	// 		return engineer != ""
	// 	})

	// 	config := analyzer.GetConfig("./config.json")

	// 	component := pages.(allEngineers, selectedEngineers, config.Repos[0].DuplicateEngineers, filterText)
	// 	content := t.Render(t.RenderParams{
	// 		C:         c,
	// 		Component: component,
	// 	})

	// 	return c.HTML(http.StatusOK, content)
	// })

	// 	e.POST("/duplicate-group", func(c echo.Context) error {
	// 		duplicatesListRaw := c.FormValue("duplicate-engineers")

	// 		duplicateEngineers := strings.Split(duplicatesListRaw, ",")

	// 		config := analyzer.GetConfig("./config.json")
	// 		config.Repos[0].DuplicateEngineers = append(config.Repos[0].DuplicateEngineers, analyzer.DuplicateEngineerGroup{
	// 			Real:       duplicateEngineers[0],
	// 			Duplicates: duplicateEngineers[1:],
	// 		})

	// 		allDups := make(map[string]bool)
	// 		for _, dupGroup := range config.Repos[0].DuplicateEngineers {
	// 			allDups[dupGroup.Real] = true

	// 			for _, dup := range dupGroup.Duplicates {
	// 				allDups[dup] = true
	// 			}
	// 		}

	// 		remainingEngineers := []string{}
	// 		for _, engineer := range InitialEngineers {
	// 			if _, found := allDups[engineer]; !found {
	// 				remainingEngineers = append(remainingEngineers, engineer)
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
