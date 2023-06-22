package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/nfwGytautas/webdev-tk/cli/templates"
	"github.com/urfave/cli/v2"
)

// PUBLIC FUNCTIONS
// ========================================================================

// Create cli.Command for init
func InitCommand() *cli.Command {
	return &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Name of the project, if empty will use directory name",
			},
		},
		Name:   "init",
		Usage:  "Initialize an empty project inside the current directory",
		Action: runInit,
	}
}

// PRIVATE FUNCTIONS
// ========================================================================

// Target for init command
func runInit(ctx *cli.Context) error {
	var (
		projectName string
		err         error
	)

	err = checkEmpty()
	if err != nil {
		println("❌  Directory not empty")
		return nil
	}

	projectName = ctx.String("name")

	fmt.Printf("🛠️  Initializing new project '%s'\n", projectName)

	err = writeConfigFile(projectName)
	if err != nil {
		return err
	}

	err = createDirectoryStructure()
	if err != nil {
		return err
	}

	println("👏 Done")

	return nil
}

func checkEmpty() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		return err
	}

	if len(files) != 0 {
		return errors.New("directory not empty")
	}

	return nil
}

func writeConfigFile(projectName string) error {
	println("✏️  Writing 'wdtk.yml'")

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Write wdtk.yml template
	data := templates.WDTKTemplateData{}
	data.CurrentDir = currentDir
	data.ProjectName = projectName

	if data.ProjectName == "" {
		data.ProjectName = filepath.Base(currentDir)
	}

	file, err := os.OpenFile("wdtk.yml", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	t, err := template.New("template").Parse(templates.WDTKTemplate)
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

func createDirectoryStructure() error {
	println("🗂️  Creating standard directory structure")
	err := os.Mkdir("services", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("frontend", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(".wdtk/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(".wdtk/logs/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(".wdtk/generated/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(".wdtk/bin/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(".wdtk/bin/services/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(".wdtk/bin/frontends/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir(".wdtk/remotes/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile("services/README.md", []byte(templates.ServicesReadME), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile("frontend/README.md", []byte(templates.FrontendReadME), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile("README.md", []byte(templates.RootReadME), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(".gitignore", []byte(templates.GitIgnore), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile("go.work", []byte(""), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
