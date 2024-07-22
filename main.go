package main

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/presentation"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
)

var help = flag.Bool("h", false, "Print help menu")
var analyzeManually = flag.Bool("a", false, "Analyze repo and gather stats")
var analyzeWithConfig = flag.Bool("c", false, "Analyze repo and gather stats using config file. (see https://yearendrecap.com/help#config)")
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

func runTest() {

}

func main() {
	err := os.Mkdir("tmp", 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic("Could not create tmp directory. Please run your Year End Recap with the correct permissions.")
	}

	if *help {
		printHelp()
	} else if *analyzeManually {
		fmt.Println("Analyzing with manual prompts...")
		analyzer.AnalyzeManually()
		fmt.Printf("\nDone! View stats by running the following command:\n\n./year-end-recap -v\n\n")
	} else if *analyzeWithConfig {
		fmt.Println("Analyzing using config...")
		analyzer.AnalyzeWithConfig("./config.json")
		fmt.Printf("\nDone! View stats by running the following command:\n\n./year-end-recap -v\n\n")
	} else if *view {
		fmt.Println("Setting up local web server...")
		presentation.RunLocalServer()
	} else if *upload {
		fmt.Println("Uploading stats to the cloud...")
	} else if *test {
		runTest()
	} else {
		printHelp()
	}
}
