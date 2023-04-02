package target

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nfwGytautas/mstk/cli/project"
	"github.com/urfave/cli"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Action for teardown target
*/
func TeardownAction(ctx *cli.Context) {
	defer TimeFn("Teardown")()
	log.Println("Tearing down")

	serviceName := ctx.Args().First()

	pc := project.ProjectConfig{}
	pc.Read()

	if serviceName == "" {
		teardownMstkK8s(pc.Project)

		// All services
		for _, service := range pc.Services {
			// TODO: Goroutines
			teardownService(service.Name, &pc)
		}
	} else {
		// Check if we have the service
		found := false
		for _, service := range pc.Services {
			if serviceName == service.Name {
				found = true
			}
		}

		if found {
			// Specific service
			teardownService(serviceName, &pc)
		} else {
			log.Printf("Service %s not found inside project", serviceName)
			panic(50)
		}
	}
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Teardown services
*/
func teardownService(service string, pc *project.ProjectConfig) {
	defer TimeFn(fmt.Sprintf("Cleaning up %s", service))()

	log.Printf("Cleaning up %s", service)

	serviceRoot := fmt.Sprintf("./services/%s/", service)

	cleanupCmd := exec.Command("kubectl", "delete", "-f", serviceRoot, "-n", pc.Project)
	log.Printf("Running %s", cleanupCmd.String())

	_, err := cleanupCmd.Output()
	if err != nil {
		// Not found is not an actual error, just it doesn't exist which is fine since we are cleaning up anyway
		if !strings.Contains(
			string((err.(*exec.ExitError).Stderr)),
			"not found",
		) {
			log.Println(string((err.(*exec.ExitError).Stderr)))
			log.Panic(err)
		}
	}
}

/*
Teardown mstk services
*/
func teardownMstkK8s(namespace string) {
	defer TimeFn("Cleaning mstk k8s")()

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Creating in %s", dirname)
	baseDir := dirname + "/mstk/k8s/"

	cleanupCmd := exec.Command("kubectl", "delete", "-f", baseDir, "-n", namespace)
	log.Printf("Running %s", cleanupCmd.String())

	_, err = cleanupCmd.Output()
	if err != nil {
		// Not found is not an actual error, just it doesn't exist which is fine since we are cleaning up anyway
		if !strings.Contains(
			string((err.(*exec.ExitError).Stderr)),
			"not found",
		) {
			log.Println(string((err.(*exec.ExitError).Stderr)))
			log.Panic(err)
		}
	}
}
