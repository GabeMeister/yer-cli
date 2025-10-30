package main

import (
	"GabeMeister/yer-cli/analyzer"
	"GabeMeister/yer-cli/presentation"
	"GabeMeister/yer-cli/utils"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

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

func runTest() {
	filepath.WalkDir("/home/gabe/dev", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			gitDir := filepath.Join(path, ".git")
			_, gitErr := os.Stat(gitDir)
			isGitDir := gitErr == nil
			isNodeModules := d.Name() == "node_modules"

			if isGitDir {
				fmt.Println(path, d.Name())
				return fs.SkipDir
			} else if isNodeModules {
				return fs.SkipDir
			} else {
				// Open the file in append mode (os.O_APPEND), write-only (os.O_WRONLY),
				// and create it if it doesn't exist (os.O_CREATE).
				// The 0644 permission grants read/write to owner, read-only to group and others.
				f, err := os.OpenFile("temp.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					log.Fatalf("failed to open file: %v", err)
				}
				defer f.Close() // Ensure the file is closed when the function exits

				// Write the string to the file
				if _, err := f.WriteString(fmt.Sprintf("%s | %s\n", path, d.Name())); err != nil {
					log.Fatalf("failed to write to file: %v", err)
				}
			}

		}

		if err != nil {
			return err
		}

		return nil
	})

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
