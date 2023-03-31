package target

import (
	"log"
	"os"

	"github.com/nfwGytautas/mstk/cli/project"
	"github.com/urfave/cli"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Action to run on template target
*/
func TemplateAction(cli *cli.Context) {
	defer TimeFn("Generating template")()

	packageName := cli.Args().Get(0)
	projectName := cli.Args().Get(1)

	if projectName == "" {
		log.Println("Empty project name. Aborting")
		return
	}

	if packageName == "" {
		log.Println("Empty package name. Aborting")
		return
	}

	log.Printf("Generating template project '%s'", projectName)

	// Create project root
	err := os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		log.Printf("Directory with the name of '%s' already exists", projectName)
		return
	}

	// Create subdirectories
	err = os.Mkdir(projectName+"/services", os.ModePerm)
	if err != nil {
		log.Printf("Failed to create services directory %v", err.Error())
		return
	}

	// CD to the new directory
	err = os.Chdir(projectName)
	if err != nil {
		log.Printf("Failed to cd into project folder %v", err.Error())
		return
	}

	// Create go.work
	writeGoWork(projectName)
	writeProjectToml(projectName, packageName)

	// TODO: Basic react environment

	log.Println("Done.")
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Writes a template go.work file
*/
func writeGoWork(projectName string) {
	log.Println("Writing go.work")

	gw := project.GoWorkConfig{}
	gw.Write()
}

/*
Writes a template mstk_project.toml
*/
func writeProjectToml(projectName, packageName string) {
	log.Println("Writing mstk_project.toml")

	pc := project.ProjectConfig{}
	pc.Project = projectName
	pc.PackageLocation = packageName + projectName + "/services/"
	pc.Write()
}
