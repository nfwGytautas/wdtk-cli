package target

import (
	"bytes"
	"fmt"
	"html/template"
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

	log.Println("Creating k8s namespace")
	createNamespace(projectName)

	// Create project config
	writeProjectToml(projectName, packageName)
	writeSecret(projectName)

	// TODO: Basic flutter environment

	log.Println("Done.")
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Template for service main function
*/
const templateSecret = `
apiVersion: v1
kind: Secret
metadata:
  name: mstk-project-secret # DO NOT CHANGE THIS
type: kubernetes.io/Opaque
stringData:
  MstkUser: USER
  MstkPsw: PSW
  Secret: API_SECRET_KEY
  Lifespan: "60"
`

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

/*
Write a template secret file for kubernetes
*/
func writeSecret(projectName string) {
	var templateData struct {
		Project string
	}

	templateData.Project = projectName

	template, err := template.New("secret").Parse(templateSecret)
	if err != nil {
		log.Println("Failed to create a secret template")
		panic(50)
	}

	buf := &bytes.Buffer{}
	err = template.Execute(buf, templateData)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.Create(fmt.Sprintf("%s-secret.yml", projectName))
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Panic(err)
	}
	file.Sync()
}
