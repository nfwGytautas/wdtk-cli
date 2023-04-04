package target

import (
	"fmt"
	"sync"

	"github.com/nfwGytautas/mstk/cli/common"
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
	defer common.TimeCurrentFn()

	common.LogInfo("Deploying")

	serviceName := ctx.Args().First()

	pc := project.ProjectConfig{}
	err := pc.Read()
	common.PanicOnError(err, "Failed to read mstk_project.toml")

	// Teardown first
	TeardownAction(ctx)

	if serviceName == "" {
		// Apply secret
		pc.Kubernetes.Apply(fmt.Sprintf("%s-secret.yml", pc.PSD.Project))
		pc.Kubernetes.ApplyMSTK()

		// All services
		var wg sync.WaitGroup
		wg.Add(len(pc.PSD.Services))
		for _, service := range pc.PSD.Services {
			go deployServiceMt(service.Name, &pc, &wg)
		}
		wg.Wait()
	} else {
		// Check if we have the service
		found := false
		for _, service := range pc.PSD.Services {
			if serviceName == service.Name {
				found = true
			}
		}

		if found {
			// Specific service
			deployService(serviceName, &pc)
		} else {
			common.LogPanic("Service %s not found inside project", serviceName)
		}
	}

	common.LogInfo("Done.")
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Multithreaded variant of deployService
*/
func deployServiceMt(service string, pc *project.ProjectConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	deployService(service, pc)
}

/*
Deploy a service to k8s
*/
func deployService(service string, pc *project.ProjectConfig) {
	common.LogTrace("Deploying %s", service)

	serviceRoot := fmt.Sprintf("./services/%s/", service)

	// Build
	common.PanicOnError(pc.Builder.Build(serviceRoot+"balancer/", fmt.Sprintf("bin/%s-balancer", service)), "Failed to build balancer")
	common.PanicOnError(pc.Builder.Build(serviceRoot+"service/", fmt.Sprintf("bin/%s-service", service)), "Failed to build service")

	// Setup services
	common.PanicOnError(pc.Docker.BuildAndPush(service+"-balancer"), "Failed to push balancer")
	common.PanicOnError(pc.Docker.BuildAndPush(service+"-service"), "Failed to push service")

	// Apply kubectl
	common.PanicOnError(pc.Kubernetes.Apply(fmt.Sprintf("k8s/deployment-%s.yml", service)), "Failed to apply deployment")
}
