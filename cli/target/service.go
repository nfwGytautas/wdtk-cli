package target

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/nfwGytautas/mstk/cli/common"
	"github.com/nfwGytautas/mstk/cli/project"
	"github.com/nfwGytautas/mstk/cli/templates"
	"github.com/urfave/cli"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Action for create service target
*/
func CreateServiceAction(ctx *cli.Context) {
	defer common.TimeCurrentFn()

	serviceName := ctx.Args().First()
	if serviceName == "" {
		common.LogPanic("Service name not provided")
	}

	common.LogInfo("Creating service %s", serviceName)

	serviceRoot := "services/" + serviceName + "/"

	pc := project.ProjectConfig{}
	common.PanicOnError(pc.Read(), "Failed to read mstk_project.toml")

	// Create directory structure
	err := os.Mkdir(serviceRoot, os.ModePerm)
	common.PanicOnError(err, "Failed to create service directory")

	err = os.Mkdir(serviceRoot+"balancer", os.ModePerm)
	common.PanicOnError(err, "Failed to create balancer directory")

	err = os.Mkdir(serviceRoot+"service", os.ModePerm)
	common.PanicOnError(err, "Failed to create service directory")

	// Write files
	writeGoMod(serviceName+"/"+"balancer", &pc)
	writeGoMod(serviceName+"/"+"service", &pc)

	writeTemplateMain(serviceRoot+"balancer/", templates.BalancerTemplate, templates.BalancerTemplateData{
		ServiceName: serviceName,
	})

	writeTemplateMain(serviceRoot+"service/", templates.ServiceTemplate, templates.ServiceTemplateData{})

	// Write docker files
	pc.Docker.WriteTemplate(serviceName+"-balancer", "bin/")
	pc.Docker.WriteTemplate(serviceName+"-service", "bin/")

	// Write k8s deployment yml file
	writeK8S("k8s/", serviceName, pc.PSD.Project)

	pc.PSD.Services = append(pc.PSD.Services, project.ServiceEntry{Name: serviceName})
	common.PanicOnError(pc.Write(), "Failed to write project config")

	common.LogInfo("Done.")
}

/*
Action for remove service target
*/
func RemoveServiceAction(ctx *cli.Context) {
	defer common.TimeCurrentFn()

	serviceName := ctx.Args().First()
	if serviceName == "" {
		common.LogPanic("Service name not provided")
	}

	pc := project.ProjectConfig{}
	pc.Read()

	// Check if we have service in the project
	for i, service := range pc.PSD.Services {
		if service.Name == serviceName {
			common.LogDebug("Found... Deleting")

			// Teardown
			teardownService(service.Name, &pc)

			// Delete
			pc.PSD.Services[i] = pc.PSD.Services[len(pc.PSD.Services)-1]
			pc.PSD.Services = pc.PSD.Services[:len(pc.PSD.Services)-1]

			err := os.RemoveAll(fmt.Sprintf("services/%s/", serviceName))
			common.PanicOnError(err, "Failed to delete service directory")

			// Remove any artifacts from docker/ k8s/ bin/ directories
			filesToRemove, err := common.GetFilesByExtension(
				"docker/",
				[]string{
					fmt.Sprintf(".%s_balancer", serviceName),
					fmt.Sprintf(".%s_service", serviceName),
				})
			common.PanicOnError(err, "Failed to query docker files")

			filesToRemove = append(filesToRemove, fmt.Sprintf("k8s/deployment-%s.yml", serviceName))
			filesToRemove = append(filesToRemove, fmt.Sprintf("bin/%s-service", serviceName))
			filesToRemove = append(filesToRemove, fmt.Sprintf("bin/%s-balancer", serviceName))

			for _, file := range filesToRemove {
				err := os.Remove(file)
				common.PanicOnError(err, "Failed to remove file")
			}

			break
		}
	}

	pc.Write()
}

// PRIVATE FUNCTIONS
// ========================================================================

/*
Write a go.mod file in the directory
*/
func writeGoMod(path string, pc *project.ProjectConfig) {
	f, err := os.OpenFile("services/"+path+"/go.mod", os.O_CREATE|os.O_WRONLY, 0644)
	common.PanicOnError(err, "Failed to create go.mod file")
	defer f.Close()

	f.WriteString("module " + pc.PSD.PackageLocation + path)
	f.WriteString("\n\n")
	f.WriteString(pc.PSD.GoVersion)
}

/*
Write a template main.go file in the directory
*/
func writeTemplateMain(path, templateString string, data any) {
	t, err := template.New("main").Parse(templateString)
	common.PanicOnError(err, "Failed to create a main template")

	buf := &bytes.Buffer{}
	err = t.Execute(buf, data)
	common.PanicOnError(err, "Failed to execute main template")

	file, err := os.Create(fmt.Sprintf("%smain.go", path))
	common.PanicOnError(err, "Failed to create main file")
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	common.PanicOnError(err, "Failed to write main file")
	file.Sync()
}

/*
Writes a k8s deployment file
*/
func writeK8S(path, service, projectName string) {
	templateData := templates.K8STemplateData{
		ProjectName: projectName,
		Service:     service,
	}

	template, err := template.New("k8s").Parse(templates.K8STemplate)
	common.PanicOnError(err, "Failed to create a k8s template")

	buf := &bytes.Buffer{}
	err = template.Execute(buf, templateData)
	common.PanicOnError(err, "Failed to execute k8s template")

	file, err := os.Create(fmt.Sprintf("%sdeployment-%s.yml", path, service))
	common.PanicOnError(err, "Failed to create deployment.yml file")
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	common.PanicOnError(err, "Failed to write deployment.yml file")
	file.Sync()
}
