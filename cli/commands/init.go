package commands

import (
	"bytes"
	"fmt"
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

	projectName = ctx.String("name")
	if ctx.String("name") == "" {
		projectName, err = os.Getwd()
		if err != nil {
			return err
		}
		projectName = filepath.Base(projectName)
	}

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

func writeConfigFile(projectName string) error {
	println("✏️  Writing 'wdtk.yml'")

	// Write wdtk.yml template
	data := templates.WDTKTemplateData{}
	data.ProjectName = projectName

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

	err = os.Mkdir("services/ExampleService", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("services/ExampleService/service/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("services/ExampleService/balancer/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("frontend", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("tools", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile("tools/UnixUpdateGoMods.sh", []byte(templates.UnixUpdateGoMods), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("deploy", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("deploy/logs/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("deploy/generated/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("deploy/bin/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("deploy/bin/unix/", os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Mkdir("deploy/unix/", os.ModePerm)
	if err != nil {
		return err
	}

	err = templates.WriteServiceTemplate("services/ExampleService/service/main.go")
	if err != nil {
		return err
	}

	err = templates.WriteServiceTemplate("services/ExampleService/balancer/main.go")
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
