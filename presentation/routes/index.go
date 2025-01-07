package routes

import "github.com/labstack/echo/v4"

func Init(e *echo.Echo) error {
	AddAnalyzerRoutes(e)

	return nil
}
