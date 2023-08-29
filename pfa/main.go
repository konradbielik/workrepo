package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"sync"
	"flag"
)

func main() {
	// Parse command-line arguments
	startDirPtr := flag.String("dir", ".", "Starting directory for search")
	patternPtr := flag.String("pattern", "", "Regex pattern to search for")
	flag.Parse()

	startDir := *startDirPtr
	pattern := *patternPtr

	// Validate input
	if pattern == "" {
		fmt.Println("Please provide a valid regex pattern.")
		return
	}

	// Initialize variables
	var totalFiles, totalHits int
	var totalSize int64
	startTime := time.Now()
	var wg sync.WaitGroup

	// Create a regular expression
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error compiling regex: %v\n", err)
		return
	}

	// Walk through directories and files recursively
	err = filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return nil
		}

		if !info.IsDir() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Read the file content
				content, err := os.ReadFile(path)
				if err != nil {
					fmt.Printf("Error reading %s: %v\n", path, err)
					return
				}

				// Search for pattern in the content
				if regex.Match(content) {
					fmt.Printf("Found at %s\n", path)
					totalHits++
				}
			}()
			
			totalFiles++
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", startDir, err)
		return
	}

	wg.Wait()

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	fmt.Printf("Total files searched: %d\n", totalFiles)
	fmt.Printf("Total size: %.2f MB\n", float64(totalSize)/1024/1024)
	fmt.Printf("Total hits: %d\n", totalHits)
	fmt.Printf("Time taken: %s\n", elapsedTime)
}
