package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"
)

var dryRun = false

func removeFile(name string) {
	var err error
	if !dryRun {
		err = os.Remove(name)
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("error removing file %s. error: %s", name, err.Error()))
	} else {
		fmt.Println("removed file", name)
	}
}

func cleanFiles(cutoff time.Time, files []fs.DirEntry) {
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			panic(err.Error())
		}
		path := info.Name()

		if info.ModTime().After(cutoff) {
			fmt.Println("keeping", path)
			continue
		}

		// If it's a file, remove it.
		if !info.IsDir() {
			removeFile(path)
			continue
		}

		// Enter directory, clean it, and cd back up to cwd.
		err = os.Chdir(path)
		if err != nil {
			fmt.Println(fmt.Sprintf("error entering dir %s. error: %s", path, err.Error()))
			continue
		}
		cleanCurrentDir(cutoff)
		err = os.Chdir("./..")
		if err != nil {
			panic(err.Error())
		}

		// Check if directory is now empty and remove it if so.
		dirFiles, err := os.ReadDir(path)
		if err != nil {
			fmt.Println(fmt.Sprintf("error removing dir %s. error: %s", path, err.Error()))
			continue
		}
		if len(dirFiles) != 0 {
			continue
		}
		if !dryRun {
			err = os.Remove(path)
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("error removing dir %s. error: %s", path, err.Error()))
			continue // Clean as much as possible.
		}
		fmt.Println("removed dir", path)
	} // Files.
}

func cleanCurrentDir(cutoff time.Time) {
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err.Error())
	}
	cleanFiles(cutoff, files)
}

// Removes all files and folders from windows appdata/local/temp folder that are
// older than 48 hours old (no specific considerations for daylight savings
// time transitions).
func main() {
	tempDir := os.TempDir()
	fmt.Println(tempDir)

	err := os.Chdir(tempDir)
	if err != nil {
		panic(err.Error())
	}

	cutoff := time.Now().Add(-time.Hour * 48)
	cleanCurrentDir(cutoff)
}
