package language

import (
	"os"

	"github.com/nfwGytautas/gdev/file"
)

// Template for main.go
const mainGoTemplate = `
package main

func main() {
	println("Running service")
}
`

// Template for go.mod file
const goModTemplate = `
module {{.Root}}services/{{.ServiceName}}

go 1.20
`

func writeGolangTemplate(data LanguageTemplateData) error {
	var err error

	// main.go
	err = os.WriteFile(data.Directory+"/main.go", []byte(mainGoTemplate), os.ModePerm)
	if err != nil {
		return err
	}

	// Create go.mod
	err = file.WriteTemplate(data.Directory+"/go.mod", goModTemplate, data)
	if err != nil {
		return nil
	}

	return nil
}
