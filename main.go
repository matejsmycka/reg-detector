package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var NO_COLOR = flag.Bool("no-color", false, "Disable color output")

func main() {
	var directory string
	var regexFile string
	flag.StringVar(&directory, "d", "", "Directory to search for files (recursively)")
	flag.StringVar(&regexFile, "f", "", "File containing regular expressions to search for")
	flag.Parse()

	if directory == "" || regexFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	regexList, err := readRegexFile(regexFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	findingCount, filesCrawled := fileIter(directory, regexList)
	fmt.Printf("Files crawled: %d\n", filesCrawled)
	fmt.Printf("Number of regexes: %d\n", len(regexList))
	fmt.Printf("Total matches: %d\n", findingCount)
}

func isText(path string, info os.FileInfo, err error) bool {
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	if isWrongExtension(path) {
		fmt.Printf("Skipping file: %s\n", path)
		return false
	}
	return true
}

func fileIter(directory string, regexList []string) (int, int) {
	var mu sync.Mutex
	fmt.Printf("Searching in directory: %s\n", directory)
	findingCount := 0
	filesCrawled := 0
	textFiles := []string{}
	_ = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if isText(path, info, err) {
			textFiles = append(textFiles, path)
		}
		return nil
	})

	var wg sync.WaitGroup

	for _, filePath := range textFiles {
		wg.Add(1)
		go func(fPath string) {
			defer wg.Done()
			filesCrawled++
			count, err := searchFileWithRegexes(fPath, regexList)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			mu.Lock()
			findingCount += count
			mu.Unlock()
		}(filePath)
	}
	wg.Wait()
	return findingCount, filesCrawled
}

func searchFileWithRegexes(filePath string, regexList []string) (int, error) {
	fileContent, err := os.ReadFile(filePath)
	findings := 0
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(fileContent), "\n")

	for lineNum, line := range lines {
		for _, regexPattern := range regexList {
			regex, err := regexp.Compile(regexPattern)
			if err != nil {
				return 0, err
			}

			matches := findRegexMatches(regex, line)
			if len(matches) > 0 {
				findings++
				printMatches(filePath, lineNum, line)

			}
		}
	}

	return findings, nil
}

func isWrongExtension(filePath string) bool {
	var suffixes = [...]string{
		".exe", ".dll", ".png", ".md", ".ico",
		".jpeg", ".zip", ".gz", ".7z", ".ttf",
		".woff", ".woff2", ".eot", ".svg", ".gif",
		".jpg", ".pdf", ".doc", ".docx", ".xls",
		".xlsx", ".ppt", ".pptx", ".mp3", ".mp4",
		".avi", ".mov", ".wav", ".flac",
		".webm", ".webp", "lock.json",
		".lock", ".pack", ".tar", ".tar.gz",
		".tar.xz", ".git",
	}

	// if filepath.contains suffix, return true

	for _, suffix := range suffixes {
		if strings.Contains(filePath, suffix) {
			return true
		}
	}
	return false
}

func findRegexMatches(regex *regexp.Regexp, line string) [][]int {
	return regex.FindAllStringIndex(line, -1)
}

func printMatches(filePath string, lineNum int, line string) {
	if !*NO_COLOR {
		fmt.Printf("%s:%d Match: ", filePath, lineNum+1)
		color.Red("%.100s", line)
	} else {
		fmt.Printf("%s:%d Match: %.100s\n", filePath, lineNum+1, line)
	}
}

func readRegexFile(filename string) ([]string, error) {
	var regexList []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			regexList = append(regexList, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return regexList, nil
}
