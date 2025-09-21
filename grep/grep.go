package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type grepOptions struct {
	recursive   bool
	lineNumbers bool
	ignoreCase  bool
	pattern     string
	files       []string
}

func main() {
	var opts grepOptions

	flag.BoolVar(&opts.recursive, "r", false, "recursively search directories")
	flag.BoolVar(&opts.lineNumbers, "n", false, "show line numbers")
	flag.BoolVar(&opts.ignoreCase, "i", false, "case-insensitive matching")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] PATTERN [FILE...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Search for PATTERN in each FILE or standard input.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: pattern is required\n")
		flag.Usage()
		os.Exit(1)
	}

	opts.pattern = args[0]
	if len(args) > 1 {
		opts.files = args[1:]
	}

	regex, err := compilePattern(opts.pattern, opts.ignoreCase)
	if err != nil {
		log.Fatalf("Error compiling pattern: %v", err)
	}

	if len(opts.files) == 0 {
		if opts.recursive {
			if err := processPath(".", regex, opts); err != nil {
				fmt.Fprintf(os.Stderr, "Error processing current directory: %v\n", err)
			}
		} else {
			searchReader(os.Stdin, "", regex, opts)
		}
		return
	}

	for _, file := range opts.files {
		if err := processPath(file, regex, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file, err)
		}
	}
}

func compilePattern(pattern string, ignoreCase bool) (*regexp.Regexp, error) {
	if ignoreCase {
		pattern = "(?i)" + pattern
	}
	return regexp.Compile(pattern)
}

func processPath(path string, regex *regexp.Regexp, opts grepOptions) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if opts.recursive {
			return searchDirectory(path, regex, opts)
		} else {
			fmt.Fprintf(os.Stderr, "%s: Is a directory\n", path)
			return nil
		}
	}

	return searchFile(path, regex, opts)
}

func searchDirectory(root string, regex *regexp.Regexp, opts grepOptions) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if isBinaryFile(path) {
			return nil
		}

		return searchFile(path, regex, opts)
	})
}

func searchFile(filename string, regex *regexp.Regexp, opts grepOptions) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return searchReader(file, filename, regex, opts)
}

func searchReader(reader interface{ Read([]byte) (int, error) }, filename string, regex *regexp.Regexp, opts grepOptions) error {
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if regex.MatchString(line) {
			printMatch(filename, line, lineNum, opts)
		}
	}

	return scanner.Err()
}

func printMatch(filename, line string, lineNum int, opts grepOptions) {
	var output strings.Builder

	if filename != "" && (len(opts.files) > 1 || opts.recursive) {
		output.WriteString(filename)
		output.WriteString(":")
	}

	if opts.lineNumbers {
		fmt.Fprintf(&output, "%d:", lineNum)
	}

	output.WriteString(line)
	fmt.Println(output.String())
}

// kinda stupid i know
func isBinaryFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	binaryExts := map[string]bool{
		".exe": true, ".bin": true, ".so": true, ".dll": true,
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".pdf": true, ".zip": true, ".tar": true, ".gz": true,
		".mp3": true, ".mp4": true, ".avi": true, ".mov": true,
		".o": true, ".a": true, ".pyc": true,
	}

	return binaryExts[ext]
}
