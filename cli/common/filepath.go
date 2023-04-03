package common

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Returns true if the current directory is MSTK root directory
*/
func IsMSTKRoot() bool {
	// Create bin directory
	cwd, err := os.Getwd()
	if err != nil {
		LogError(err.Error())
		return false
	}

	// Check that we are in MSTK
	base := filepath.Base(cwd)
	if base != "MSTK" && base != "mstk" {
		return false
	}

	return true
}

/*
Get all directories inside the root directory
*/
func GetDirectories(root string) ([]string, error) {
	var files []string

	f, err := os.Open(root)

	if err != nil {
		return files, err
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, file.Name())
		}
	}

	return files, nil
}

/*
Get the mstk setup directory
*/
func GetMSTKDir() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return dirname + "/mstk/", nil
}

/*
Copies all files from one directory to another
*/
func CopyDir(from, to string, ignoreExtensions []string) {
	f, err := os.Open(from)
	if err != nil {
		log.Panic(err)
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Copying directory '%s' to '%s'", from, to)

fileLoop:
	for _, file := range fileInfo {
		if file.IsDir() {
			// TODO: No symlinks etc.
			os.Mkdir(to+file.Name(), os.ModePerm)
			CopyDir(from+file.Name(), to+file.Name(), ignoreExtensions)
		} else {
			for _, ext := range ignoreExtensions {
				if ext == filepath.Ext(file.Name()) {
					continue fileLoop
				}
			}

			log.Printf("\t%s", file.Name())

			// Read all content of src to data, may cause OOM for a large file.
			data, err := os.ReadFile(from + file.Name())
			if err != nil {
				log.Panic(err)
			}

			// Write data to dst
			err = os.WriteFile(to+file.Name(), data, fs.ModePerm)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

/*
Get all files that have an extension that is in the provided array
*/
func GetFilesByExtension(root string, extensions []string) ([]string, error) {
	f, err := os.Open(root)
	if err != nil {
		return nil, err
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	var result []string

fileLoop:
	for _, file := range fileInfo {
		if !file.IsDir() {
			for _, ext := range extensions {
				if ext == filepath.Ext(file.Name()) {
					continue fileLoop
				}
			}

			log.Printf("\t%s", file.Name())
			result = append(result, root+file.Name())
		}
	}

	return result, nil
}

// PRIVATE FUNCTIONS
// ========================================================================
