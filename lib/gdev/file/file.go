package file

import (
	"bytes"
	"html/template"
	"os"
)

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

// Write a template to file
func WriteTemplate(path, templateString string, data any) error {
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	t, err := template.New("template").Parse(templateString)
	if err != nil {
		return err
	}

	out := &bytes.Buffer{}
	err = t.Execute(out, data)
	if err != nil {
		return err
	}

	_, err = file.Write(out.Bytes())
	if err != nil {
		return err
	}
	file.Sync()

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================
