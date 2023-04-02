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
// ========================================================================]

/*
Action for create service target
*/
func CreateServiceAction(ctx *cli.Context) {
	defer TimeFn("Create service")()

	serviceName := ctx.Args().First()
	if serviceName == "" {
		log.Println("Service name not given")
		panic(50)
	}

	log.Printf("Creating service %s", serviceName)

	pc := project.ProjectConfig{}
	pc.Read()

	// Create directory structure
	serviceRoot := "services/" + serviceName + "/"
	err := os.Mkdir(serviceRoot, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create folder %v", err.Error())
		panic(50)
	}

	err = os.Mkdir(serviceRoot+"balancer", os.ModePerm)
	if err != nil {
		log.Printf("Failed to create folder %v", err.Error())
		panic(50)
	}

	err = os.Mkdir(serviceRoot+"service", os.ModePerm)
	if err != nil {
		log.Printf("Failed to create folder %v", err.Error())
		panic(50)
	}

	// Write files
	writeGoMod(serviceName+"/"+"balancer", &pc)
	writeGoMod(serviceName+"/"+"service", &pc)

	writeTemplateMain(serviceRoot+"balancer/", balancerTemplate)
	writeTemplateMain(serviceRoot+"service/", serviceTemplate)

	// Write k8s deployment yml file
	writeK8S(serviceRoot, serviceName, pc.Project)

	pc.Services = append(pc.Services, project.ServiceEntry{Name: serviceName})
	pc.Write()

	log.Println("Done.")
}

/*
Action for remove service target
*/
func RemoveServiceAction(ctx *cli.Context) {
	defer TimeFn("Remove service")()

	serviceName := ctx.Args().First()
	if serviceName == "" {
		log.Println("Service not specified")
		panic(50)
	}

	pc := project.ProjectConfig{}
	pc.Read()

	// Check if we have service in the project
	for i, service := range pc.Services {
		if service.Name == serviceName {
			log.Println("Found... Deleting")

			pc.Services[i] = pc.Services[len(pc.Services)-1]
			pc.Services = pc.Services[:len(pc.Services)-1]

			err := os.RemoveAll(fmt.Sprintf("services/%s/", serviceName))
			if err != nil {
				log.Printf("Failed to delete service folder %v", err.Error())
				panic(51)
			}

			break
		}
	}

	pc.Write()
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Write a go.mod file in the directory
*/
func writeGoMod(path string, pc *project.ProjectConfig) {
	f, err := os.OpenFile("services/"+path+"/go.mod", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create go.mod")
		panic(50)
	}

	defer f.Close()

	f.WriteString("module " + pc.PackageLocation + path)
	f.WriteString("\n\n")

	f.WriteString(pc.GoVersion)
}

/*
Write a template main.go file in the directory
*/
func writeTemplateMain(path, template string) {
	f, err := os.OpenFile(path+"main.go", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create main.go")
		panic(50)
	}

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
	if err != nil {
		log.Println("Failed to create a k8s template")
		panic(50)
	}

	buf := &bytes.Buffer{}
	err = template.Execute(buf, templateData)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.Create(fmt.Sprintf("%sdeployment-%s.yml", path, service))
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
