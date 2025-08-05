package presentation

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/presentation/routes"
	"GabeMeister/yer-cli/utils"
	"embed"
	"fmt"
	"time"

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

	fmt.Print("\n", "💻 Opening browser...", "\n\n")

	go func() {
		e.Logger.Fatal(e.Start(":4000"))
	}()

	// Small delay to ensure server starts
	time.Sleep(100 * time.Millisecond)

	analyzer.OpenBrowser("http://localhost:4000/create-recap")

	// Keep main goroutine alive
	select {}
}

func RunPresentationPage() {
	godotenv.Load()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	routes.Init(e, static)

	fmt.Print("\n", "💻 Opening browser...", "\n\n")

	go func() {
		e.Logger.Fatal(e.Start(":4000"))
	}()

	// Small delay to ensure server starts
	time.Sleep(100 * time.Millisecond)

	if !utils.IsDevMode() {
		analyzer.OpenBrowser("http://localhost:4000")
	}

	// Keep main goroutine alive
	select {}
}
