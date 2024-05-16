package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"net/http"

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

//go:embed views/*
var f embed.FS

type Greeting struct {
	Name string
}
type Repo struct {
	RepoName string
}

func runTest() {
	indexOutput := renderTemplate("views/index.html", Greeting{Name: "Zach"})
	fmt.Println(indexOutput)
	repoOutput := renderTemplate("views/repo.html", Repo{RepoName: "Next.js"})
	fmt.Println(repoOutput)
}

func renderTemplate(path string, data interface{}) string {
	htmlStr, _ := f.ReadFile(path)
	t := template.Must(template.New(path).Parse(string(htmlStr)))

	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func runLocalServer() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, renderTemplate("views/index.html", Greeting{Name: "Josh"}))
	})
	e.GET("/repo", func(c echo.Context) error {
		return c.HTML(http.StatusOK, renderTemplate("views/repo.html", Repo{RepoName: "RB Frontend"}))
	})
	e.Logger.Fatal(e.Start(":4000"))
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
		runTest()
	} else {
		printHelp()
	}
}
