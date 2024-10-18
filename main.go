package main

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/presentation"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

var help = flag.Bool("h", false, "Print help menu")
var analyzeRepo = flag.Bool("a", false, "Analyze repo and gather stats")
var configFile = flag.String("c", "", "Specify config file to analyze with. (see https://yearendrecap.com/help#config)")
var view = flag.Bool("v", false, "View stats in a local presentation")

// var upload = flag.Bool("u", false, "Upload stats to the cloud, to be viewed anywhere")
var test = flag.Bool("t", false, "pls ignore")

func init() {
	flag.Parse()
}

func printHelp() {
	fmt.Println("Year End Recap CLI")
	fmt.Println()
	flag.PrintDefaults()
}

func runTest() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Suffix = " Hello 1 "
	s.Start() // Start the spinner
	time.Sleep(2 * time.Second)
	s.Suffix = " Hello 2 "
	time.Sleep(2 * time.Second)
	s.Suffix = " Hello 3 "
	time.Sleep(2 * time.Second)
	s.Stop()
}

func main() {
	err := os.Mkdir("tmp", 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic("Could not create tmp directory. Please run your Year End Recap with the correct permissions.")
	}

	if *help {
		printHelp()
	} else if *analyzeRepo {
		fmt.Println("Analyzing with manual prompts...")
		result := analyzer.AnalyzeManually()
		if result {
			fmt.Printf("\nDone! View stats by running the following command:\n\n./year-end-recap -v\n\n")
		} else {
			fmt.Println("Failed to analyze repo. Please try again!")

		}
	} else if *configFile != "" {
		fmt.Println("Analyzing using config...")
		result := analyzer.AnalyzeWithConfig(*configFile)
		if result {
			fmt.Printf("\nDone! View stats by running the following command:\n\n./year-end-recap -v\n\n")
		} else {
			fmt.Println("Failed to analyze repo. Please try again!")
		}
	} else if *view {
		fmt.Println("Setting up local web server...")
		presentation.RunLocalServer()
	} else if *test {
		runTest()
	} else {
		printHelp()
	}
}
