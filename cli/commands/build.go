package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nfwGytautas/gdev/array"
	"github.com/nfwGytautas/webdev-tk/cli/build"
	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/nfwGytautas/webdev-tk/cli/util"
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
		println("❌  Build command expects a either 'all', a name or a list of names that you want to build")
		return nil
	}

	// Read wdtk.yml
	cfg := types.WDTKConfig{}
	err := cfg.Read()
	if err != nil {
		return err
	}

	// Create build log file
	logFile := fmt.Sprintf(".wdtk/logs/%s.build.log", time.Now().Format("2006-01-02 15:04:05"))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	f.Close()

	println("🔨  Services...")
	err = buildServices(ctx, &cfg, logFile)
	if err != nil {
		return err
	}

	println("🔨  Frontends...")
	err = buildFrontends(ctx, &cfg, logFile)
	if err != nil {
		return err
	}

	return nil
}

func buildServices(ctx *cli.Context, cfg *types.WDTKConfig, logFile string) error {
	servicesToBuild := ctx.Args().Slice()[1:]
	numBuilt := 0
	numFailed := 0

	buildAll := ctx.Args().Get(0) == "all"

	// Local
	for _, service := range cfg.GetServicesOfType(types.SERVICE_TYPE_LOCAL) {
		build := buildAll || array.IsElementInArray(servicesToBuild, service.Name)

		if build {
			err := buildLocalService(cfg, &service, logFile)
			if err != nil {
				fmt.Println(err)
				numFailed++
			} else {
				numBuilt++
			}
		}
	}

	// Git
	for _, service := range cfg.GetServicesOfType(types.SERVICE_TYPE_GIT) {
		build := buildAll || array.IsElementInArray(servicesToBuild, service.Name)

		if build {
			err := buildGitService(cfg, &service, logFile)
			if err != nil {
				fmt.Println(err)
				numFailed++
			} else {
				numBuilt++
			}
		}
	}

	println(fmt.Sprintf("--- %d built, %d failed ---", numBuilt, numFailed))

	if numFailed != 0 {
		return errors.New("one or more builds failed")
	}

	return nil
}

func buildLocalService(cfg *types.WDTKConfig, service *types.ServiceDescriptionConfig, logFile string) error {
	println(util.SPACING_1 + "- " + service.Name)
	abs, err := filepath.Abs(".wdtk/bin/services/")
	if err != nil {
		return err
	}

	data := build.ServiceBuildData{
		SourceDir:   "services/" + service.Name,
		OutDir:      abs + "/",
		ServiceName: service.Name,
	}
	return build.BuildService(data, *service.Source.Language)
}

func buildGitService(cfg *types.WDTKConfig, service *types.ServiceDescriptionConfig, logFile string) error {
	println(util.SPACING_1 + "- " + service.Name)
	abs, err := filepath.Abs(".wdtk/bin/services/")
	if err != nil {
		return err
	}

	source, err := service.GitLocalDestination()
	if err != nil {
		return err
	}

	data := build.ServiceBuildData{
		SourceDir:   source,
		OutDir:      abs + "/",
		ServiceName: service.Name,
	}
	return build.BuildService(data, *service.Source.Language)
}

func buildFrontends(ctx *cli.Context, cfg *types.WDTKConfig, logFile string) error {
	if cfg.Frontend == nil {
		return errors.New("tried to build frontends without any defined frontends")
	}

	frontendsToBuild := ctx.Args().Slice()[1:]
	buildAll := ctx.Args().Get(0) == "all"

	numBuilt := 0
	numFailed := 0

	abs, err := filepath.Abs(".wdtk/bin/frontends/")
	if err != nil {
		return err
	}

	for _, platform := range cfg.Frontend.Platforms {
		println(util.SPACING_1 + "- '" + platform.Type + "' with " + platform.Toolchain)
		toBuild := buildAll || array.IsElementInArray(frontendsToBuild, platform.Type)

		if toBuild {
			data := build.FrontendBuildData{
				Type:   platform.Type,
				OutDir: abs + "/",
			}
			err := build.BuildFrontend(data, platform.Toolchain)
			if err != nil {
				log.Println(err)
				numFailed++
			} else {
				numBuilt++
			}
		}
	}

	println(fmt.Sprintf("--- %d built, %d failed ---", numBuilt, numFailed))

	if numFailed != 0 {
		return errors.New("one or more builds failed")
	}

	return nil
}
