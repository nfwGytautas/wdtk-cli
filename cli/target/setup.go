package target

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/urfave/cli"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Flags for setup target
*/
var SetupFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "platform",
		Usage: "platform for docker",
	},
}

/*
Execute setup target
*/
func SetupAction(ctx *cli.Context) {
	defer TimeFn("Setup")()
	// TODO: Find the MSTK installation path automatically

	log.Println("Running setup")
	EnsureMSTKRoot()

	log.Println("Creating bin")
	os.Mkdir("bin/", os.ModePerm)

	log.Println("Compiling services")
	services := GetMstkServicesList()

	log.Printf("Found %v services %v", len(services), services)

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, service := range services {
		go compileService(service, &wg)
	}
	wg.Wait()

	log.Println("Setup done, your minikube environment should have mstk microservices up and running")
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Version for images
*/
const version = "0.0.0"

/*
Compiles a single service
*/
func compileService(path string, wg *sync.WaitGroup) {
	defer TimeFn(fmt.Sprintf("Preparing '%s'", path))()
	defer wg.Done()

	log.Printf("Compiling %s", path)

	serviceName := filepath.Base(path)
	targetFile := fmt.Sprintf("./bin/%s", serviceName)
	sourceDir := fmt.Sprintf("./gomods/%s/", serviceName)

	// Build sources
	buildSourcesForDocker(targetFile, sourceDir)

	// Generate docker files
	writeDockerFile("./bin", serviceName)

	// Push to minikube
	cfg := setupServiceCfg{
		tag:        "mstk/",
		name:       serviceName,
		dockerPath: "./bin/Dockerfile." + serviceName,
	}
	setupService(cfg)

	// Apply kubectl commands
	applyKubectl("kubes/" + serviceName)
}
