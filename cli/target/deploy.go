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
Action for deploy target
*/
func DeployAction(ctx *cli.Context) {
	defer TimeFn("Deploy")()

	log.Println("Deploying")

	serviceName := ctx.Args().First()

	pc := project.ProjectConfig{}
	pc.Read()

	// Teardown first
	TeardownAction(ctx)

	if serviceName == "" {
		// Apply secret
		applyKubectl(fmt.Sprintf("%s-secret.yml", pc.Project), pc.Project)

		// All services
		for _, service := range pc.Services {
			// TODO: Goroutines
			deployService(service.Name, &pc)
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
			deployService(serviceName, &pc)
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
Deploy a service to k8s
*/
func deployService(service string, pc *project.ProjectConfig) {
	log.Printf("Deploying %s", service)

	serviceRoot := fmt.Sprintf("./services/%s/", service)

	// Create bin directory
	err := os.Mkdir(serviceRoot+"bin/", os.ModePerm)
	if err != nil {
		if !strings.Contains(err.Error(), "file exists") {
			log.Printf("Failed to create bin folder %v", err.Error())
			panic(60)
		}
	}

	binDir := serviceRoot + "bin/"

	// Build
	buildSourcesForDocker(binDir+"balancer", serviceRoot+"balancer/")
	buildSourcesForDocker(binDir+"service", serviceRoot+"balancer/")

	// Docker
	writeDockerFile(binDir, "balancer")
	writeDockerFile(binDir, "service")

	// Setup services
	balancerCfg := setupServiceCfg{
		tag:        pc.Project + "/",
		name:       fmt.Sprintf("%s-balancer", service),
		dockerPath: binDir + "Dockerfile.balancer",
	}
	serviceCfg := setupServiceCfg{
		tag:        pc.Project + "/",
		name:       fmt.Sprintf("%s-service", service),
		dockerPath: binDir + "Dockerfile.service",
	}
	setupService(balancerCfg)
	setupService(serviceCfg)

	// Apply kubectl
	applyKubectl(serviceRoot, pc.Project)
}

/*
Executes a rolling restart
*/
func restartKubernetes() {
	defer TimeFn("Restart")()

	applyCmd := exec.Command(
		"kubectl", "rollout", "restart",
	)
	log.Printf("Running %s", applyCmd.String())

	err := applyCmd.Run()
	if err != nil {
		log.Panic(err)
	}
}
