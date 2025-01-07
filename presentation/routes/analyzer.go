package routes

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/presentation/views/components/AnalyzeManuallyPage"
	presentation_views_pages "GabeMeister/yer-cli/presentation/views/pages"
	t "GabeMeister/yer-cli/presentation/views/template"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AddAnalyzerRoutes(e *echo.Echo) {
	initialEngineers := []string{"Kenny", "Isaac Neace", "Gabe Jensen", "ktrotter", "Kaleb Trotter", "Stephen Bremer", "Kenny Kline", "Ezra Youngren", "Isaac", "Steve Bremer"}

	e.GET("/sortable", func(c echo.Context) error {

		analyzer.InitConfig(analyzer.ConfigFileOptions{
			RepoDir:                "/home/gabe/dev/rb-frontend",
			MasterBranchName:       "master",
			IncludedFileExtensions: []string{"ts", "tsx", "js", "jsx"},
			ExcludedDirs:           []string{"node_modules", "build"},
			DuplicateEngineers:     make(map[string]string),
			IncludeFileBlames:      true,
		})

		content := t.Render(t.RenderParams{
			C:         c,
			Component: presentation_views_pages.Sortable(initialEngineers, []string{}, make(map[string]string)),
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/search-engineers", func(c echo.Context) error {
		text := c.FormValue("filter-text")
		text = strings.ToLower(text)

		matches := []string{}
		for _, engineer := range initialEngineers {
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

	e.POST("/items", func(c echo.Context) error {
		data, err := c.FormParams()
		if err != nil {
			panic(err)
		}

		leftItemsStr := data["left-form-items"][0]
		leftItems := strings.Split(leftItemsStr, ",")

		rightItemsStr := data["right-form-items"][0]
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

		component := presentation_views_pages.Sortable(allEngineers, selectedEngineers, config.Repos[0].DuplicateEngineers)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})

	e.POST("/submit-duplicate", func(c echo.Context) error {
		duplicatesList := c.FormValue("duplicate-engineers")
		userNames := strings.Split(duplicatesList, ",")

		config := analyzer.GetConfig("./config.json")
		config.Repos[0].DuplicateEngineers[userNames[0]] = userNames[1]
		analyzer.SaveDataToFile(config, "./config.json")

		component := presentation_views_pages.Sortable(initialEngineers, []string{}, config.Repos[0].DuplicateEngineers)
		content := t.Render(t.RenderParams{
			C:         c,
			Component: component,
		})

		return c.HTML(http.StatusOK, content)
	})
}
