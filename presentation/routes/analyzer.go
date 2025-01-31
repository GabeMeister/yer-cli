package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"slices"

	"github.com/samber/lo"

	"GabeMeister/yer-cli/presentation/views/components/AnalyzeManuallyPage"
	"GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var InitialEngineers = []string{"Kenny", "Kenny1", "Kenny2", "Isaac Neace", "Gabe Jensen", "ktrotter", "Kaleb Trotter", "Stephen Bremer", "Kenny Kline", "Ezra Youngren", "Isaac", "Steve Bremer"}

func addAnalyzerRoutes(e *echo.Echo) {
	e.GET("/create-recap", func(c echo.Context) error {
		analyzer.InitConfig(analyzer.ConfigFileOptions{
			RepoDir:                "",
			MasterBranchName:       "",
			IncludedFileExtensions: []string{},
			ExcludedDirs:           []string{},
			DuplicateEngineers:     []analyzer.DuplicateEngineerGroup{},
			IncludeFileBlames:      true,
		})

		config := analyzer.GetConfig("./config.json")

		content := t.Render(t.RenderParams{
			C: c,
			Component: pages.AnalyzeManually(
				InitialEngineers,
				[]string{},
				config.Repos[0].DuplicateEngineers,
				"",
			),
		})

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/config-file", func(c echo.Context) error {
		recapName := c.FormValue("recap-name")
		recapName = strings.ToLower(recapName)

		updatedRecap := analyzer.ConfigFile{
			Repos: []analyzer.RepoConfig{
				{
					Name: recapName,
				}},
		}

		analyzer.UpdateConfig(updatedRecap)

		component := AnalyzeManuallyPage.SuccessMessage()
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/search-engineers", func(c echo.Context) error {
		text := c.FormValue("filter-text")
		text = strings.ToLower(text)

		dupEngineersRaw := c.FormValue("duplicate-engineers")
		tempDupEngineers := strings.Split(dupEngineersRaw, ",")

		matches := []string{}
		for _, engineer := range InitialEngineers {
			lowerCaseEngineer := strings.ToLower(engineer)

			if strings.Contains(lowerCaseEngineer, text) && !slices.Contains(tempDupEngineers, engineer) {
				matches = append(matches, engineer)
			}
		}
		component := AnalyzeManuallyPage.AllEngineersList(matches)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.PATCH("/temp-duplicate-group", func(c echo.Context) error {

		leftItemsStr := c.FormValue("all-engineers")
		leftItems := strings.Split(leftItemsStr, ",")

		rightItemsStr := c.FormValue("duplicate-engineers")
		rightItems := strings.Split(rightItemsStr, ",")

		filterText := c.FormValue("filter-text")

		allEngineers := lo.Filter(leftItems, func(engineer string, _ int) bool {
			return engineer != ""
		})

		selectedEngineers := lo.Filter(rightItems, func(engineer string, _ int) bool {
			return engineer != ""
		})

		config := analyzer.GetConfig("./config.json")

		component := pages.AnalyzeManually(allEngineers, selectedEngineers, config.Repos[0].DuplicateEngineers, filterText)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/duplicate-group", func(c echo.Context) error {
		duplicatesListRaw := c.FormValue("duplicate-engineers")

		duplicateEngineers := strings.Split(duplicatesListRaw, ",")

		config := analyzer.GetConfig("./config.json")
		config.Repos[0].DuplicateEngineers = append(config.Repos[0].DuplicateEngineers, analyzer.DuplicateEngineerGroup{
			Real:       duplicateEngineers[0],
			Duplicates: duplicateEngineers[1:],
		})

		allDups := make(map[string]bool)
		for _, dupGroup := range config.Repos[0].DuplicateEngineers {
			allDups[dupGroup.Real] = true

			for _, dup := range dupGroup.Duplicates {
				allDups[dup] = true
			}
		}

		remainingEngineers := []string{}
		for _, engineer := range InitialEngineers {
			if _, found := allDups[engineer]; !found {
				remainingEngineers = append(remainingEngineers, engineer)
			}
		}

		analyzer.SaveDataToFile(config, "./config.json")

		component := pages.AnalyzeManually(
			remainingEngineers,
			[]string{},
			config.Repos[0].DuplicateEngineers,
			"",
		)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})
}
