package file

import "os"

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// Check if file exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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

// PRIVATE FUNCTIONS
// ========================================================================
