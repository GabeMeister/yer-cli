package presentation

import (
	"GabeMeister/yer-cli/presentation/routes"
	"embed"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

//go:embed static/*
var static embed.FS

func RunLocalServer() {
	godotenv.Load()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	routes.Init(e, static)

	fmt.Println("\nDone! Browse to http://localhost:4000/")
	e.Logger.Fatal(e.Start(":4000"))
}
