package commands

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nfwGytautas/mstk/lib/gdev/array"
	"github.com/nfwGytautas/mstk/lib/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/urfave/cli/v2"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

func BuildCommand() *cli.Command {
	return &cli.Command{
		Flags:     []cli.Flag{},
		Name:      "build",
		Usage:     "Build services",
		ArgsUsage: "[all|services...]",
		Action:    runBuild,
	}
}

// PRIVATE FUNCTIONS
// ========================================================================

func runBuild(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		println("❌  Build command expects a either 'all', a service name or a list of services that you want to build")
		return nil
	}

	// Read wdtk.yml
	cfg := types.WDTKConfig{}
	err := cfg.Read()
	if err != nil {
		return err
	}

	println("🔨  Building...")

	// Create build log file
	logFile := fmt.Sprintf("deploy/logs/%s.build.log", time.Now().Format("2006-01-02 15:04:05"))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	servicesToBuild := ctx.Args().Slice()[1:]
	numBuilt := 0
	numFailed := 0

	buildAll := ctx.Args().Get(0) == "all"

	for _, service := range cfg.Services {
		if buildAll {
			err := runBuildScript(&cfg, &service, logFile)
			if err != nil {
				log.Println(err)
				numFailed++
			} else {
				numBuilt++
			}
		} else {
			if array.IsElementInArray(servicesToBuild, service.Name) {
				err := runBuildScript(&cfg, &service, logFile)
				if err != nil {
					log.Println(err)
					numFailed++
				} else {
					numBuilt++
				}
			}
		}
	}

	println(fmt.Sprintf("--- %d built, %d failed ---", numBuilt, numFailed))

	if numFailed != 0 {
		return errors.New("one or more builds failed")
	}

	return nil
}

func runBuildScript(cfg *types.WDTKConfig, service *types.ServiceDescriptionConfig, logFile string) error {
	abs, err := filepath.Abs("deploy/unix/")
	if err != nil {
		return err
	}

	println("Building " + service.Name)

	// Run the deployment script
	var outb, errb bytes.Buffer

	cmd := exec.Command("bash", fmt.Sprintf("./%s_BUILD_UNIX.sh", service.Name))
	cmd.Dir = abs
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()

	if err != nil {
		file.Append(logFile, outb.String())
		file.Append(logFile, errb.String())
		file.Append(logFile, err.Error())

		return err
	} else {
		file.Append(logFile, outb.String())
		file.Append(logFile, errb.String())

		return err
	}
}
