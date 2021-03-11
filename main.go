package main

import (
	"fmt"
	"os"
	"time"
)

// Removes all files and folders from windows appdata/local/temp folder that are
// older than 48 hours old (no specific considerations for daylight savings
// time transitions).
func main() {
	tempDir := os.TempDir()
	fmt.Println(tempDir)

	files, err := os.ReadDir(tempDir)
	if err != nil {
		panic(err.Error())
	}
	err = os.Chdir(tempDir)
	if err != nil {
		panic(err.Error())
	}

	cutoff := time.Now().Add(-time.Hour * 48)
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			panic(err.Error())
		}
		if info.ModTime().After(cutoff) {
			fmt.Println("keeping", info.Name())
			continue
		}
		fmt.Println("deleting", info.Name())
		if info.IsDir() {
			err = os.RemoveAll(f.Name())
			if err != nil {
				fmt.Println(fmt.Sprintf("error removing dir %s. error: %s", info.Name(), err.Error()))
				// Continue, try to remove as many as possible.
			}
		} else {
			err = os.Remove(f.Name())
			if err != nil {
				fmt.Println(fmt.Sprintf("error removing file %s. error: %s", info.Name(), err.Error()))
				// Continue, try to remove as many as possible.
			}
		}
	}
}
