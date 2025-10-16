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
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func printHelp() {
	fmt.Println("Year End Recap CLI")
	fmt.Println()
	customUsage()
}

func customUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])

	// Define custom order
	order := []string{"s", "a", "v"}

	for _, name := range order {
		f := flag.Lookup(name)
		if f != nil {
			fmt.Fprintf(os.Stderr, "  -%s:   %s\n", f.Name, f.Usage)
		}
	}
}

func runTest() {
	start := time.Now()

	repoPath := "/home/gabe/dev/rb-frontend"
	// Get all tracked files
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = repoPath
	output, _ := cmd.Output()
	allFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	files := []string{}
	for _, file := range allFiles {
		if strings.HasSuffix(file, ".ts") {
			files = append(files, file)
		}
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // limit concurrent blames
	var mu sync.Mutex              // mutex to protect the counter
	var totalLines int64           // shared counter

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			cmd := exec.Command("git", "blame", f)
			cmd.Dir = repoPath
			output, _ := cmd.Output()
			text := string(output)
			length := len(strings.Split(text, "\n"))

			fmt.Printf("%s | %d lines\n", file, length)

			// Add to the overall number of lines (thread-safe)
			mu.Lock()
			totalLines += int64(length)
			mu.Unlock()
		}(file)

		// // Synchronous way
		// cmd := exec.Command("git", "blame", file)
		// cmd.Dir = repoPath
		// output, _ := cmd.Output()
		// text := string(output)
		// length := len(strings.Split(text, "\n"))
		// fmt.Printf("%d lines | %s \n", length, file)
		// totalLines += int64(length)

	}
	wg.Wait()

	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Time taken: %f sec\n", time.Since(start).Seconds())
}

func main() {
	godotenv.Load()

	var help = flag.Bool("h", false, "Print help menu")
	var setupConfig = flag.Bool("s", false, "(Step 1) Setup a new Year End Recap configuration")
	var analyzeRepo = flag.Bool("a", false, "(Step 2) Analyze your Git repo(s) to gather highly amusing Git stats")
	var view = flag.Bool("v", false, "(Step 3) View your highly amusing Git stats")

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
	} else if *setupConfig {
		presentation.RunCreateRecapPage()
		fmt.Print("\n\nComplete! Now run `./year-end-recap -a` to analyze your repos\n\n")
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
			fmt.Printf("%s│ ./year-end-recap -v                  │%s\n", yellow, reset)
			fmt.Printf("%s└──────────────────────────────────────┘%s\n", yellow, reset)
			fmt.Println()
		} else {
			fmt.Println("\nPlease run `./year-end-recap -s` to setup your recap configuration, then try analyzing.")
		}
	} else if *view {
		presentation.RunPresentationPage()
	} else if *test {
		runTest()
	} else {
		printHelp()
	}
}
