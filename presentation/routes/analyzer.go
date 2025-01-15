package routes

import (
	"GabeMeister/yer-cli/analyzer"

	"GabeMeister/yer-cli/presentation/views/components/AnalyzeManuallyPage"
	"GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var InitialEngineers = []string{"Kenny", "Kenny1", "Kenny2", "Isaac Neace", "Gabe Jensen", "ktrotter", "Kaleb Trotter", "Stephen Bremer", "Kenny Kline", "Ezra Youngren", "Isaac", "Steve Bremer"}

func addAnalyzerRoutes(e *echo.Echo) {
	e.GET("/analyze-manually", func(c echo.Context) error {
		analyzer.InitConfig(analyzer.ConfigFileOptions{
			RepoDir:                "/home/gabe/dev/rb-frontend",
			MasterBranchName:       "master",
			IncludedFileExtensions: []string{"ts", "tsx", "js", "jsx"},
			ExcludedDirs:           []string{"node_modules", "build"},
			DuplicateEngineers: []analyzer.DuplicateEngineerGroup{
				{Real: "ktrotter", Duplicates: []string{"Kaleb Trotter"}},
			},
			IncludeFileBlames: true,
		})

		config := analyzer.GetConfig("./config.json")

		content := t.Render(t.RenderParams{
			C: c,
			Component: pages.AnalyzeManually(
				InitialEngineers,
				[]string{},
				config.Repos[0].DuplicateEngineers,
			),
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/search-engineers", func(c echo.Context) error {
		text := c.FormValue("filter-text")
		text = strings.ToLower(text)

		matches := []string{}
		for _, engineer := range InitialEngineers {
			lowerCaseEngineer := strings.ToLower(engineer)

			if strings.Contains(lowerCaseEngineer, text) {
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
		data, err := c.FormParams()
		if err != nil {
			panic(err)
		}

		leftItemsStr := data["all-engineers"][0]
		leftItems := strings.Split(leftItemsStr, ",")

		rightItemsStr := data["duplicate-engineers"][0]
		rightItems := strings.Split(rightItemsStr, ",")

		allEngineers := []string{}
		for _, s := range leftItems {
			if s == "" {
				continue
			}

			allEngineers = append(allEngineers, s)
		}

		selectedEngineers := []string{}
		for _, s := range rightItems {
			if s == "" {
				continue
			}

			selectedEngineers = append(selectedEngineers, s)
		}

		config := analyzer.GetConfig("./config.json")

		component := pages.AnalyzeManually(allEngineers, selectedEngineers, config.Repos[0].DuplicateEngineers)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/duplicate-group", func(c echo.Context) error {
		allEngineersRaw := c.FormValue("all-engineers")
		duplicatesListRaw := c.FormValue("duplicate-engineers")

		allEngineersLeft := strings.Split(allEngineersRaw, ",")
		duplicateEngineers := strings.Split(duplicatesListRaw, ",")

		config := analyzer.GetConfig("./config.json")
		config.Repos[0].DuplicateEngineers = append(config.Repos[0].DuplicateEngineers, analyzer.DuplicateEngineerGroup{
			Real:       duplicateEngineers[0],
			Duplicates: duplicateEngineers[1:],
		})

		analyzer.SaveDataToFile(config, "./config.json")

		component := pages.AnalyzeManually(allEngineersLeft, []string{}, config.Repos[0].DuplicateEngineers)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})
}
