package target

import (
	"fmt"
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

	common.LogInfo("Running setup")

	common.LogTrace("Creating mstk directory")

	baseDir, err := common.GetMSTKDir()
	common.PanicOnError(err, "Failed to get MSTK root directory")

	common.LogTrace("Creating %s", baseDir)

	common.PanicOnError(os.Mkdir(baseDir, os.ModePerm), "Failed to create mstk root directory")
	common.PanicOnError(os.Mkdir(baseDir+"bin", os.ModePerm), "Failed to create bin directory")
	common.PanicOnError(os.Mkdir(baseDir+"docker", os.ModePerm), "Failed to create docker directory")
	common.PanicOnError(os.Mkdir(baseDir+"k8s", os.ModePerm), "Failed to create k8s directory")

	common.LogTrace("Compiling services")
	services, err := GetMstkServicesList()
	common.PanicOnError(err, "Failed to get mstk services list")

	common.LogTrace("Found %v services %v", len(services), services)

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, service := range services {
		go compileService(service, &wg)
	}
	wg.Wait()

	common.CopyDir("kubes/", baseDir+"k8s/", []string{".md"})

	common.LogInfo("Setup done you can now create a template project using 'mstk template' command")
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Returns a list of services in mstk
*/
func GetMstkServicesList() ([]string, error) {
	directories, err := common.GetDirectories("gomods/")
	if err != nil {
		return nil, err
	}

	return directories, nil
}

/*
Compiles a single service
*/
func compileService(path string, wg *sync.WaitGroup) {
	defer common.TimeCurrentFn()
	defer wg.Done()

	common.LogTrace("Compiling %s", path)

	mstkRoot, err := common.GetMSTKDir()
	common.PanicOnError(err, "Failed to get mstk root directory")

	serviceName := filepath.Base(path)
	targetFile := fmt.Sprintf("%sbin/%s", mstkRoot, serviceName)
	sourceDir := fmt.Sprintf("./gomods/%s/", serviceName)

	builder := api.CreateBuilder()
	docker := api.CreateDocker(mstkRoot, "mstk")

	// Build sources
	common.PanicOnError(builder.Build(sourceDir, targetFile), "Failed to build")

	// Generate docker files
	common.PanicOnError(docker.WriteTemplate(serviceName, "bin/"), "Failed to create template")

	// Push to minikube
	common.PanicOnError(docker.BuildAndPush(serviceName), "Failed to push image")
}
