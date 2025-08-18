package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wget <url>")
		os.Exit(1)
	}

	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("Invalid URL: %v\n", err)
		os.Exit(1)
	}

	filename := path.Base(parsedURL.Path)
	if filename == "" || filename == "/" {
		filename = "index.html"
	}

	resp, err := http.Get(rawURL)
	if err != nil {
		fmt.Printf("Error downloading: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP error: %s\n", resp.Status)
		os.Exit(1)
	}

	out, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Downloaded %s\n", filename)
}
