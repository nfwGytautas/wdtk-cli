package target

import (
	"fmt"
	"log"
	"sync"

	"github.com/nfwGytautas/mstk/cli/common"
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
	defer common.TimeCurrentFn()
	log.Println("Tearing down")

	serviceName := ctx.Args().First()

	pc := project.ProjectConfig{}
	pc.Read()

	if serviceName == "" {
		pc.Kubernetes.DeleteMSTK()

		// All services
		var wg sync.WaitGroup
		wg.Add(len(pc.PSD.Services))
		for _, service := range pc.PSD.Services {
			go teardownServiceMt(service.Name, &pc, &wg)
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
			teardownService(serviceName, &pc)
		} else {
			common.LogPanic("Service %s not found in project", serviceName)
		}
	}
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Multithreaded version of teardown services
*/
func teardownServiceMt(service string, pc *project.ProjectConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	teardownService(service, pc)
}

/*
Teardown services
*/
func teardownService(service string, pc *project.ProjectConfig) {
	serviceRoot := fmt.Sprintf("k8s/deployment-%s.yml", service)
	pc.Kubernetes.Delete(serviceRoot)
}
