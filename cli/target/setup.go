package target

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/nfwGytautas/mstk/cli/api"
	"github.com/nfwGytautas/mstk/cli/common"
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
	defer common.TimeCurrentFn()

	if !common.IsMSTKRoot() {
		common.LogPanic("setup needs to be ran inside a mstk root directory")
	}
	// TODO: Find the MSTK installation path automatically

	log.Println("Running setup")

	log.Println("Creating mstk directory")

	baseDir, err := common.GetMSTKDir()
	common.PanicOnError(err, "Failed to get MSTK root directory")

	log.Printf("Creating %s", baseDir)

	common.PanicOnError(os.Mkdir(baseDir, os.ModePerm), "Failed to create mstk root directory")
	common.PanicOnError(os.Mkdir(baseDir+"bin", os.ModePerm), "Failed to create bin directory")
	common.PanicOnError(os.Mkdir(baseDir+"docker", os.ModePerm), "Failed to create docker directory")
	common.PanicOnError(os.Mkdir(baseDir+"k8s", os.ModePerm), "Failed to create k8s directory")

	log.Println("Compiling services")
	services := GetMstkServicesList()

	log.Printf("Found %v services %v", len(services), services)

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, service := range services {
		go compileService(service, &wg)
	}
	wg.Wait()

	common.CopyDir("kubes/", baseDir+"k8s/", []string{".md"})

	log.Println("Setup done you can now create a template project using 'mstk template' command")
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Returns a list of services in mstk
*/
func GetMstkServicesList() []string {
	directories, err := common.GetDirectories("gomods/services/")
	if err != nil {
		log.Panic(err)
	}

	return directories
}

/*
Compiles a single service
*/
func compileService(path string, wg *sync.WaitGroup) {
	defer common.TimeCurrentFn()
	defer wg.Done()

	log.Printf("Compiling %s", path)

	mstkRoot, err := common.GetMSTKDir()
	common.PanicOnError(err, "Failed to get mstk root directory")

	serviceName := filepath.Base(path)
	targetFile := fmt.Sprintf("%sbin/%s", mstkRoot, serviceName)
	sourceDir := fmt.Sprintf("./gomods/services/%s/", serviceName)

	builder := api.CreateBuilder()
	docker := api.CreateDocker(mstkRoot, "mstk")

	// Build sources
	common.PanicOnError(builder.Build(sourceDir, targetFile), "Failed to build")

	// Generate docker files
	common.PanicOnError(docker.WriteTemplate(serviceName, "bin/"), "Failed to create template")

	// Push to minikube
	common.PanicOnError(docker.BuildAndPush(serviceName), "Failed to push image")
}
