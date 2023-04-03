package target

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/nfwGytautas/mstk/cli/common"
	"github.com/nfwGytautas/mstk/cli/project"
	"github.com/urfave/cli"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

/*
Template for balancer main function
*/
const balancerTemplate = `
package main

import "log"

func main() {
	log.Println("MSTK template balancer")
}

`

/*
Template for service main function
*/
const serviceTemplate = `
package main

import "log"

func main() {
	log.Println("MSTK template service")
}

`

/*
Template for deployment file
*/
const k8sTemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}-{{.Service}}-service
spec:
  selector:
    matchLabels:
      app: {{.ProjectName}}-{{.Service}}-service
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: {{.ProjectName}}-{{.Service}}-service
    spec:
      containers:
      - image: {{.ProjectName}}/{{.Service}}-service:0.0.0
        name: {{.ProjectName}}-{{.Service}}-service
        imagePullPolicy: Never
        resources:
          limits:
            memory: "500M"
            cpu: "50m"
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: API_SECRET
          valueFrom:
            secretKeyRef:
              name: mstk-project-secret
              key: Secret
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}-{{.Service}}-balancer
spec:
  selector:
    matchLabels:
      app: {{.ProjectName}}-{{.Service}}-balancer
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: {{.ProjectName}}-{{.Service}}-balancer
    spec:
      containers:
      - image: {{.ProjectName}}/{{.Service}}-balancer:0.0.0
        name: {{.ProjectName}}-{{.Service}}-balancer
        imagePullPolicy: Never
        resources:
          limits:
            memory: "500M"
            cpu: "50m"
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: API_SECRET
          valueFrom:
            secretKeyRef:
              name: mstk-project-secret
              key: Secret
---
apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}-{{.Service}}-service
spec:
  ports:
    - port: 8080
      targetPort: 8080
  selector:
    app: {{.ProjectName}}-{{.Service}}-service
---
apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}-{{.Service}}-balancer
spec:
  ports:
    - port: 8080
      targetPort: 8080
  selector:
    app: {{.ProjectName}}-{{.Service}}-balancer
`

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

	log.Printf("Creating service %s", serviceName)

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

	writeTemplateMain(serviceRoot+"balancer/", balancerTemplate)
	writeTemplateMain(serviceRoot+"service/", serviceTemplate)

	// Write docker files
	pc.Docker.WriteTemplate(serviceName+"-balancer", "bin/")
	pc.Docker.WriteTemplate(serviceName+"-service", "bin/")

	// Write k8s deployment yml file
	writeK8S("k8s/", serviceName, pc.PSD.Project)

	pc.PSD.Services = append(pc.PSD.Services, project.ServiceEntry{Name: serviceName})
	common.PanicOnError(pc.Write(), "Failed to write project config")

	log.Println("Done.")
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
			log.Println("Found... Deleting")

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
			filesToRemove = append(filesToRemove, fmt.Sprintf("bin/%s_service", serviceName))
			filesToRemove = append(filesToRemove, fmt.Sprintf("bin/%s_balancer", serviceName))

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
func writeTemplateMain(path, template string) {
	f, err := os.OpenFile(path+"main.go", os.O_CREATE|os.O_WRONLY, 0644)
	common.PanicOnError(err, "Failed to create main.go file")
	defer f.Close()

	f.WriteString(template)
}

/*
Writes a k8s deployment file
*/
func writeK8S(path, service, projectName string) {
	var templateData struct {
		ProjectName string
		Service     string
	}

	templateData.ProjectName = projectName
	templateData.Service = service

	template, err := template.New("k8s").Parse(k8sTemplate)
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
