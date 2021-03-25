package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const separator = string(os.PathSeparator)

// Flags.
var isDryRun *bool
var isVerbose *bool
var hours *int64

// State.
var numFilesRemoved uint64

func removeFile(name string) {
	var err error
	if !*isDryRun {
		err = os.Remove(name)
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("error removing file %s. error: %s", name, err.Error()))
		return
	}
	atomic.AddUint64(&numFilesRemoved, 1)
	if *isVerbose {
		fmt.Println("removed file", name)
	}
}

func cleanFiles(wg *sync.WaitGroup, cutoff time.Time, files []fs.DirEntry, path string) {
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			panic(err.Error())
		}
		filename := info.Name()

		// Handle file.
		if !info.IsDir() {
			if info.ModTime().After(cutoff) {
				if *isVerbose {
					fmt.Println("keeping", filename)
				}
				continue
			}
			removeFile(path + separator + filename)
			continue
		}

		// Clean sub directory.
		subPath := path + separator + filename
		fmt.Println("cleaning sub directory", subPath)
		wg.Add(1)
		go cleanDir(wg, cutoff, subPath)

		// Check if sub directory is now empty and remove it if so.
		dirFiles, err := os.ReadDir(subPath)
		if err != nil {
			fmt.Println(fmt.Sprintf("error removing dir %s. error: %s", subPath, err.Error()))
			continue
		}
		if len(dirFiles) != 0 {
			continue
		}
		if !*isDryRun {
			err = os.Remove(subPath)
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("error removing dir %s. error: %s", subPath, err.Error()))
			continue // Clean as much as possible.
		}
		fmt.Println("removed dir", subPath)
	} // Files.
}

func cleanDir(wg *sync.WaitGroup, cutoff time.Time, path string) {
	defer wg.Done()
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("unabled to read dir", path, err.Error())
	}
	cleanFiles(wg, cutoff, files, path)
}

// Removes all files and folders from windows appdata/local/temp folder that are
// older than 48 hours old (no specific considerations for daylight savings
// time transitions).
func main() {
	// Allow user to see output before final exit.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("--------------------------------------------------------------")
			fmt.Println("recovered panic")
			fmt.Println(r)
			fmt.Println(string(debug.Stack()))
		}
		time.Sleep(time.Second * 2)
	}()

	// Parse flags.
	isDryRun = flag.Bool("dry", false, "A dry run will only print out file removal actions, not perform them.")
	isVerbose = flag.Bool("verbose", false, "Prints out keep and remove actions.")
	hours = flag.Int64("hours", 72, "Files older than this, in hours, will be removed.")
	flag.Parse()
	if *isDryRun {
		fmt.Println("--------------------------------------------------------------")
		fmt.Println("dry run enabled")
	}

	// Navigate to temp folder.
	tempDir := os.TempDir()
	err := os.Chdir(tempDir)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("cleaning", tempDir)

	startTime := time.Now()
	cutoff := startTime.Add(-time.Hour * time.Duration(*hours))
	var wg sync.WaitGroup
	wg.Add(1)
	go cleanDir(&wg, cutoff, ".")
	wg.Wait()

	fmt.Println("--------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("removed %v files", numFilesRemoved))
	fmt.Println("completed in", time.Now().Sub(startTime))
}
