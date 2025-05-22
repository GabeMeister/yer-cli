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

func RunCreateRecapPage() {
	godotenv.Load()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	routes.Init(e, static)

	fmt.Println("\nTo setup your Year End Recap, browse to http://localhost:4000/create-recap")
	e.Logger.Fatal(e.Start(":4000"))
}

func RunPresentationPage() {
	godotenv.Load()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	routes.Init(e, static)

	fmt.Println("\nDone! Browse to http://localhost:4000/")
	e.Logger.Fatal(e.Start(":4000"))
}
