package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

var help = flag.Bool("h", false, "Print help menu")
var analyze = flag.Bool("a", false, "Analyze repo and gather stats")
var config = flag.String("c", "", "Specify path to config file. (see https://yearendrecap.com/help#config)")
var view = flag.Bool("v", false, "View stats in a local presentation")
var upload = flag.Bool("u", false, "Upload stats to the cloud, to be viewed anywhere")
var test = flag.Bool("t", false, "Test something out")

func init() {
	flag.Parse()
}

func printHelp() {
	fmt.Println("Year End Recap CLI")
	fmt.Println()
	flag.PrintDefaults()
}

func runLocalServer() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func main() {
	if *help {
		printHelp()
	} else if *analyze {
		if *config == "" {
			fmt.Println("Analyzing with manual prompts...")
		} else {
			fmt.Println("Analyzing using config:", *config)
		}
	} else if *view {
		fmt.Println("Viewing stats...")
		runLocalServer()
	} else if *upload {
		fmt.Println("Uploading stats to the cloud...")
	} else if *test {
		data, err := os.ReadFile("views/index.html")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))
	} else {
		printHelp()
	}
}
