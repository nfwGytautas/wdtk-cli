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
	EnsureMSTKRoot()
	// TODO: Find the MSTK installation path automatically

	log.Println("Running setup")

	log.Println("Creating mstk directory")
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Creating in %s", dirname)
	baseDir := dirname + "/mstk/"
	os.Mkdir(baseDir, os.ModePerm)
	os.Mkdir(baseDir+"bin", os.ModePerm)
	os.Mkdir(baseDir+"docker", os.ModePerm)
	os.Mkdir(baseDir+"k8s", os.ModePerm)

	log.Println("Compiling services")
	services := GetMstkServicesList()

	log.Printf("Found %v services %v", len(services), services)

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, service := range services {
		go compileService(service, &wg)
	}
	wg.Wait()

	copyDir("kubes/", baseDir+"k8s/", []string{".md"})

	log.Println("Setup done you can now create a template project using 'mstk template' command")
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

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	base := fmt.Sprintf("%s/mstk/", dirname)
	serviceName := filepath.Base(path)
	targetFile := fmt.Sprintf("%sbin/%s", base, serviceName)
	sourceDir := fmt.Sprintf("./gomods/services/%s/", serviceName)

	// Build sources
	buildSourcesForDocker(targetFile, sourceDir)

	// Generate docker files
	writeDockerFile(base, serviceName)

	// Push to minikube
	cfg := setupServiceCfg{
		tag:        "mstk/",
		name:       serviceName,
		dockerPath: base + "docker/Dockerfile." + serviceName,
		runDir:     base,
	}
	setupService(cfg)
}
