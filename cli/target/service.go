package target

import (
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
	log.Printf("Creating service %s", serviceName)

	pc := project.ProjectConfig{}
	pc.Read()

	gw := project.GoWorkConfig{}
	gw.Read()

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
	writeGoMod(serviceName+"/"+"balancer/", &pc, &gw)
	writeGoMod(serviceName+"/"+"service/", &pc, &gw)

	gw.Write()

	writeTemplateMain(serviceRoot+"balancer/", balancerTemplate)
	writeTemplateMain(serviceRoot+"service/", serviceTemplate)

	// Write k8s deployment yml file
	// TODO: K8S deploy file

	pc.Services = append(pc.Services, project.ServiceEntry{Name: serviceName})
	pc.Write()

	log.Println("Done.")
}

/*
Write a go.mod file in the directory
*/
func writeGoMod(path string, pc *project.ProjectConfig, gw *project.GoWorkConfig) {
	f, err := os.OpenFile("services/"+path+"go.mod", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create go.mod")
		panic(50)
	}

	defer f.Close()

	f.WriteString("package " + pc.PackageLocation + path)
	f.WriteString("\n\n")

	f.WriteString(gw.GoVersion)

	gw.UseDirectives = append(gw.UseDirectives, "services/"+path)
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

// ========================================================================
// PRIVATE
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
