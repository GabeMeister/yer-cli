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
var setupConfig = flag.Bool("s", false, "Setup a new Year End Recap configuration")
var analyzeRepo = flag.Bool("a", false, "Analyze repo(s) to gather highly amusing Git stats")
var view = flag.Bool("v", false, "View your highly amusing Git stats")

// var upload = flag.Bool("u", false, "Upload stats to the cloud, to be viewed anywhere")
var test = flag.Bool("t", false, "Run test")

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
	} else if *setupConfig {
		presentation.RunCreateRecapPage()
		fmt.Print("\n\nComplete! Now run `./year-end-recap -a` to analyze your repos\n\n")
	} else if *analyzeRepo {
		fmt.Println("Analyzing...")
		result := analyzer.AnalyzeRepos()
		if result {
			fmt.Printf("\nDone! View stats by running the following command:\n\n./year-end-recap -v\n\n")
		} else {
			fmt.Println("Failed to analyze repo. Please try again!")
		}
	} else if *view {
		presentation.RunPresentationPage()
	} else if *test {
		runTest()
	} else {
		printHelp()
	}
}
