package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func randomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}
	return string(bytes), nil
}

func shredFile(path string, passes int) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s is not a valid file: %v", path, err)
	}
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", path)
	}
	fileSize := fileInfo.Size()

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", path, err)
	}
	defer f.Close()

	for i := 0; i < passes; i++ {
		fmt.Printf("Pass %d of %d...\n", i+1, passes)
		randomData := make([]byte, fileSize)
		if _, err := rand.Read(randomData); err != nil {
			return fmt.Errorf("error generating random data for pass %d: %v", i+1, err)
		}

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("error seeking to start of file: %v", err)
		}

		if _, err := f.Write(randomData); err != nil {
			return fmt.Errorf("error writing pass %d: %v", i+1, err)
		}

		if err := f.Sync(); err != nil {
			return fmt.Errorf("error syncing pass %d: %v", i+1, err)
		}
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("error closing file %s: %v", path, err)
	}

	dir := filepath.Dir(path)
	newName, err := randomString(8)
	if err != nil {
		return fmt.Errorf("error generating random name: %v", err)
	}
	newPath := filepath.Join(dir, newName)

	if err := os.Rename(path, newPath); err != nil {
		return fmt.Errorf("error renaming file %s to %s: %v", path, newPath, err)
	}

	if err := os.Remove(newPath); err != nil {
		return fmt.Errorf("error removing file %s: %v", newPath, err)
	}

	fmt.Printf("File %s has been securely shredded.\n", path)
	return nil
}

func shredDirectory(path string, passes int) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error stating directory: %v", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	var dirs []string

	err = filepath.WalkDir(path, func(walkPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			dirs = append(dirs, walkPath)
			return nil
		}

		fmt.Printf("Shredding file %s...\n", walkPath)
		if err := shredFile(walkPath, passes); err != nil {
			fmt.Printf("Error shredding file %s: %v\n", walkPath, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %v", err)
	}

	for i := len(dirs)/2 - 1; i >= 0; i-- {
		opp := len(dirs) - 1 - i
		dirs[i], dirs[opp] = dirs[opp], dirs[i]
	}

	for _, dir := range dirs {
		if err := os.Remove(dir); err != nil {
			fmt.Printf("Error removing directory %s: %v\n", dir, err)
		} else {
			fmt.Printf("Directory %s has been removed.\n", dir)
		}
	}

	fmt.Printf("Directory %s has been securely shredded and removed.\n", path)
	return nil
}

func main() {
	var passes int
	var recursive bool

	flag.IntVar(&passes, "p", 3, "Number of overwrite passes")
	flag.BoolVar(&recursive, "r", false, "Recursively shred directories")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Error: path argument is required.")
		os.Exit(1)
	}
	path := args[0]

	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	switch {
	case fileInfo.Mode().IsRegular():
		if err := shredFile(path, passes); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case fileInfo.IsDir():
		if !recursive {
			fmt.Printf("Error: %s is a directory. Use -r to shred directories recursively.\n", path)
			os.Exit(1)
		}
		if err := shredDirectory(path, passes); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Error: %s is not a valid file or directory.\n", path)
		os.Exit(1)
	}
}
