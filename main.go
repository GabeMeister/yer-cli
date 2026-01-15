package main

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/presentation"
	"GabeMeister/yer-cli/utils"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/joho/godotenv"
)

var VERSION = "0.1"

func printHelp() {
	fmt.Println("Year End Recap CLI")
	fmt.Println()
	customUsage()
}

func customUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])

	// Define custom order
	order := []string{"s", "a", "p", "v", "h"}

	for _, name := range order {
		f := flag.Lookup(name)
		if f != nil {
			fmt.Fprintf(os.Stderr, "  -%s:   %s\n", f.Name, f.Usage)
		}
	}
}

func addOneToSlice(s []int) {
	for i := range s {
		s[i]++
	}
}

func addOneToGabe(g *Gabe) {
	g.Age += 1
}

type Gabe struct {
	Age int
}

func runTest() {
	// mySlice := []int{1, 2, 3, 4}
	// fmt.Print("\n\n", "*** mySlice ***", "\n", mySlice, "\n\n\n")
	// addOneToSlice(mySlice)
	// fmt.Print("\n\n", "*** mySlice ***", "\n", mySlice, "\n\n\n")
	gabe1 := Gabe{
		Age: 10,
	}
	fmt.Print("\n\n", "*** gabe1 ***", "\n", gabe1, "\n\n\n")
	addOneToGabe(&gabe1)
	fmt.Print("\n\n", "*** gabe1 ***", "\n", gabe1, "\n\n\n")

}

func main() {
	godotenv.Load()

	var help = flag.Bool("h", false, "Print help menu")
	var version = flag.Bool("v", false, "Print version")
	var setupConfig = flag.Bool("s", false, "[Step 1] Setup a new Year End Recap configuration")
	var analyzeRepo = flag.Bool("a", false, "[Step 2] Analyze your Git repo(s) to gather highly amusing Git stats")
	var view = flag.Bool("p", false, "[Step 3] View your highly amusing Git stats in a presentation")

	var test *bool
	var calculateOnly *bool
	prodCalculateOnly := false
	if utils.IsDevMode() {
		test = flag.Bool("t", false, "Run test")
		calculateOnly = flag.Bool("c", false, "Just run calculations while analyzing, and skip gathering metrics step")
	} else {
		calculateOnly = &prodCalculateOnly
	}

	flag.Usage = customUsage
	flag.Parse()

	err := os.Mkdir("tmp", 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic("Could not create tmp directory. Please run your Year End Recap with the correct permissions")
	}

	if *help {
		printHelp()
	} else if *version {
		fmt.Printf("Year End Recap v%s\n", VERSION)
	} else if *setupConfig {
		presentation.RunCreateRecapPage()
		fmt.Print("\n\nComplete! Now run `yer -a` to analyze your repos\n\n")
	} else if *analyzeRepo {
		result := analyzer.AnalyzeRepos(*calculateOnly)
		if result {
			yellow := "\033[1;33m"
			reset := "\033[0m"

			fmt.Println()
			fmt.Printf("%s┌──────────────────────────────────────┐%s\n", yellow, reset)
			fmt.Printf("%s│ Done! Now run the following command  │%s\n", yellow, reset)
			fmt.Printf("%s│ to view your stats:                  │%s\n", yellow, reset)
			fmt.Printf("%s│                                      │%s\n", yellow, reset)
			fmt.Printf("%s│ yer -p                               │%s\n", yellow, reset)
			fmt.Printf("%s└──────────────────────────────────────┘%s\n", yellow, reset)
			fmt.Println()
		} else {
			fmt.Println("\nPlease run `yer -s` to setup your recap configuration, then try analyzing.")
		}
	} else if *view {
		presentation.RunPresentationPage()
	} else if utils.IsDevMode() && *test {
		runTest()
	} else {
		printHelp()
	}
}
