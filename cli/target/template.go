package target

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/nfwGytautas/mstk/cli/common"
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
	defer common.TimeCurrentFn()

	packageName := cli.Args().Get(0)
	projectName := cli.Args().Get(1)

	if projectName == "" {
		common.LogPanic("Project name not provided")
	}

	if packageName == "" {
		common.LogPanic("Package name not provided")
	}

	common.LogInfo("Generating template project '%s'", projectName)

	// Create project root
	common.PanicOnError(os.Mkdir(projectName, os.ModePerm), "Failed to create project root directory")

	// Create subdirectories
	common.PanicOnError(os.Mkdir(projectName+"/services", os.ModePerm), "Failed to create services directory inside project root")
	common.PanicOnError(os.Mkdir(projectName+"/bin", os.ModePerm), "Failed to create bin directory inside project root")
	common.PanicOnError(os.Mkdir(projectName+"/k8s", os.ModePerm), "Failed to create k8s directory inside project root")
	common.PanicOnError(os.Mkdir(projectName+"/docker", os.ModePerm), "Failed to create docker directory inside project root")

	// CD to the new directory
	common.PanicOnError(os.Chdir(projectName), "Failed to cd into project root")

	// Write project toml
	writeProjectToml(projectName, packageName)

	pc := project.ProjectConfig{}
	pc.Read()

	common.LogTrace("Creating k8s namespace")
	pc.Kubernetes.CreateNamespace()

	// Create project config
	common.PanicOnError(writeSecret(projectName), "Failed to write secret file")

	// TODO: Basic flutter environment

	common.LogInfo("Done.")
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
	common.LogTrace("Writing mstk_project.toml")

	pc := project.ProjectConfig{}
	pc.PSD.Project = projectName
	pc.PSD.PackageLocation = packageName + projectName + "/services/"
	pc.Write()
}

/*
Write a template secret file for kubernetes
*/
func writeSecret(projectName string) error {
	var templateData struct {
		Project string
	}

	templateData.Project = projectName

	template, err := template.New("secret").Parse(templateSecret)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = template.Execute(buf, templateData)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s-secret.yml", projectName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	file.Sync()

	return nil
}
