package routes

import (
	"embed"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, static embed.FS) error {
	AddAnalyzerRoutes(e)
	AddPresentationRoutes(e)
	AddResourceRoutes(e, static)
	AddDebuggingRoutes(e)

	return nil
}
