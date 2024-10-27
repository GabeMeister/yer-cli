package main

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/presentation"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
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

func myFunc(args ...string) string {
	cmd := exec.Command(args[0], args[1:]...)
	output, _ := cmd.Output()
	fmt.Println(string(output))

	return string(output)
}

type Gabe struct {
	Age int
}

func runTest() {
	gabes := make(map[string]Gabe)

	gabe1 := Gabe{Age: 30}
	gabes["First"] = gabe1

	gabe1.Age = 100

	fmt.Println(gabes)

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
