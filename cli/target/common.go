package target

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Ensure we are running mstk in mstk root
*/
func EnsureMSTKRoot() {
	// Create bin directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	// Check that we are in MSTK
	base := filepath.Base(cwd)
	if base != "MSTK" && base != "mstk" {
		log.Panicf("Target is only allowed in root directory, but was ran inside %s", base)
	}
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
Time the execution time of a function
*/
func TimeFn(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}

/*
Returns a list of services in mstk
*/
func GetMstkServicesList() []string {
	directories, err := GetDirectories("gomods/")
	if err != nil {
		log.Panic(err)
	}

	// Filter out only services
	var services []string
	for _, dir := range directories {
		if !strings.HasSuffix(dir, "-api") {
			services = append(services, dir)
		}
	}

	return services
}
