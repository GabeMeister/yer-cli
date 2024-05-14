package main

import (
	"flag"
	"fmt"
)

var help = flag.Bool("h", false, "Print help menu")
var analyze = flag.Bool("a", false, "Analyze repo and gather stats")
var config = flag.String("c", "", "Specify path to config file. (see https://yearendrecap.com/help#config)")
var view = flag.Bool("v", false, "View stats in a local presentation")
var upload = flag.Bool("u", false, "Upload stats to the cloud, to be viewed anywhere")

func init() {
	flag.Parse()
}

func printHelp() {
	fmt.Println("Year End Recap CLI")
	fmt.Println()
	flag.PrintDefaults()
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
	} else if *upload {
		fmt.Println("Uploading stats to the cloud...")
	} else {
		printHelp()
	}
}
