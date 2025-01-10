package routes

import (
	"embed"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, static embed.FS) error {
	addAnalyzerRoutes(e)
	addPresentationRoutes(e)
	addResourceRoutes(e, static)
	addDebuggingRoutes(e)

	return nil
}
